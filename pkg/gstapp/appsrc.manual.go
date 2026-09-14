package gstapp

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/go-gst/go-glib/pkg/glib/v2"
	"github.com/go-gst/go-glib/pkg/gobject/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

type AppSrcExtManual interface {
	// Writer is a convenience method that returns a [io.WriteCloser] for the [AppSrc].
	// This can be used with [io.Copy] or [io.CopyBuffer] to write data into the AppSrc.
	// You must call [AppSrcWriter.Close] when done writing to signal the end of the stream and to free the resources used by the writer.
	Writer() *AppSrcWriter
	// ReadFrom implements [io.ReaderFrom] for [AppSrc]. It blocks until the reader returns EOF or an error.
	ReadFrom(r io.Reader) (n int64, err error)
}

func (a *AppSrcInstance) Writer() *AppSrcWriter {
	return newAppSrcWriter(a)
}

var _ io.WriteCloser = (*AppSrcWriter)(nil)

func newAppSrcWriter(app AppSrc) *AppSrcWriter {
	writer := &AppSrcWriter{
		app: app,
	}

	writer.dataCond = sync.NewCond(&writer.lock)

	writer.needDataHandle = app.ConnectNeedData(func(app AppSrc, length uint) {
		writer.setQueueFull(false)
	})

	writer.enoughDataHandle = app.ConnectEnoughData(func(app AppSrc) {
		writer.setQueueFull(true)
	})

	return writer
}

type AppSrcWriter struct {
	// synchronization, must not be held while pushing a buffer. AppSrc
	// has an internal queue, so pushing a buffer too many is not a problem, but
	// we should avoid pushing more data afer enough-data has been emitted
	lock      sync.Mutex
	dataCond  *sync.Cond
	queueFull bool
	closed    bool

	app AppSrc

	enoughDataHandle gobject.SignalHandle
	needDataHandle   gobject.SignalHandle
}

// Close implements [io.Closer]. It disconnects all set up signal handlers and emits an EOS event into the AppSrc.
//
// Close is idempotent. It returns an error only if the EOS event could not be
// sent; an AppSrc that is already flushing or at EOS is reported as success.
func (a *AppSrcWriter) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.closed {
		return nil
	}

	a.closed = true

	a.app.HandlerDisconnect(a.needDataHandle)
	a.app.HandlerDisconnect(a.enoughDataHandle)

	a.dataCond.Broadcast()

	return eosError(a.app.EndOfStream())
}

// Write implements [io.Writer].
func (a *AppSrcWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		// nothing to push, and an empty glib.Bytes has no data pointer, which
		// gst_buffer_new_wrapped_bytes rejects
		return 0, nil
	}

	if !a.waitForQueueSpace() {
		return 0, io.ErrClosedPipe
	}

	bytes := glib.NewBytes(p)

	buffer := gst.NewBufferWrappedBytes(bytes)

	// PushBuffer may synchronously emit the need data signal, so no mutex must be held here
	if err := pushError(a.app.PushBuffer(buffer)); err != nil {
		return 0, err
	}

	return len(p), nil
}

// waitForQueueSpace blocks until the appsrc is ready to accept more data or is closed.
// It returns false if the appsrc is closed, true otherwise.
func (a *AppSrcWriter) waitForQueueSpace() bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	for !a.closed && a.queueFull {
		a.dataCond.Wait()
	}

	return !a.closed
}

// setQueueFull sets the queueFull flag and broadcasts to all waiting goroutines that the state has changed.
func (a *AppSrcWriter) setQueueFull(full bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.queueFull = full
	a.dataCond.Broadcast()
}

var _ io.ReaderFrom = (*AppSrcInstance)(nil)

// ReadFrom implements [io.ReaderFrom] for [AppSrc].
func (a *AppSrcInstance) ReadFrom(r io.Reader) (n int64, err error) {
	var handler gobject.SignalHandle

	// buffer will be re-used and re-allocated as needed
	var buffer []byte

	// errchan channel will be used to signal when the reading is done
	errchan := make(chan error)

	done := false

	totalRead := int64(0)

	handler = a.ConnectNeedData(func(app AppSrc, length uint) {
		if done {
			return
		}

		if int(length) > cap(buffer) {
			buffer = make([]byte, length)
		}

		readBuffer := buffer[:length]

		nRead, readErr := r.Read(readBuffer)

		totalRead += int64(nRead)

		if nRead > 0 {
			// If we read some data, we push it to the appsrc
			bytes := glib.NewBytes(readBuffer[:nRead])
			gstBuffer := gst.NewBufferWrappedBytes(bytes)

			if err := pushError(app.PushBuffer(gstBuffer)); err != nil {
				done = true
				errchan <- err
				return
			}
		}

		if errors.Is(readErr, io.EOF) {
			// If we reach EOF, we signal that we're done and return
			done = true
			errchan <- nil
			return
		}

		if readErr != nil {
			done = true
			errchan <- readErr
			return
		}

	})

	// Wait for the reading to be done or for an error to occur
	err = <-errchan
	a.HandlerDisconnect(handler)

	// Signal that no more data will be sent. A read error takes precedence over
	// a failure to send EOS, since it is the more useful diagnostic.
	if eosErr := eosError(a.EndOfStream()); err == nil {
		err = eosErr
	}

	return totalRead, err
}

func pushError(flow gst.FlowReturn) error {
	err := flow.Error()
	if err == nil {
		return nil
	}

	if errors.Is(err, gst.ErrFlowFlushing) || errors.Is(err, gst.ErrFlowEOS) {
		return fmt.Errorf("cannot push buffer: %w: %w", err, io.ErrClosedPipe)
	}

	return fmt.Errorf("cannot push buffer: %w", err)
}

func eosError(flow gst.FlowReturn) error {
	err := flow.Error()
	if errors.Is(err, gst.ErrFlowFlushing) || errors.Is(err, gst.ErrFlowEOS) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("cannot send end of stream: %w", err)
	}

	return nil
}
