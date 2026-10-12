package gstapp

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/go-gst/go-glib/pkg/gobject/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

type AppSinkExtManual interface {
	// Reader is a convenience method that returns a [io.ReadCloser] for the [AppSink].
	// This can be used with [io.Copy] or [io.CopyBuffer] to read the data of the AppSink.
	// You must call [AppSinkReader.Close] when done reading to free the resources used by the reader.
	Reader() *AppSinkReader
	// WriteTo implements [io.WriterTo] for [AppSink]. It blocks until the AppSink reaches EOS or an error occurs.
	WriteTo(w io.Writer) (n int64, err error)
}

func (a *AppSinkInstance) Reader() *AppSinkReader {
	return newAppSinkReader(a)
}

var _ io.ReadCloser = (*AppSinkReader)(nil)

var ErrBufferNotMappable = errors.New("could not map buffer for reading")

// appSinkPollInterval is the maximum delay between a Close and a blocked read
// observing it. The AppSink API cannot interrupt a pending pull, so reads poll
// with TryPullSample instead of blocking in PullSample.
const appSinkPollInterval = 100 * gst.Millisecond

// appSinkEOSGrace is how long a read waits for the eos signal once the AppSink
// reports EOS. The signal races with the pull that returns no sample.
const appSinkEOSGrace = 100 * time.Millisecond

func newAppSinkReader(app AppSink) *AppSinkReader {
	reader := &AppSinkReader{
		app: app,
		eos: make(chan struct{}),
	}

	reader.eosHandle = app.ConnectEOS(func(app AppSink) {
		reader.eosOnce.Do(func() { close(reader.eos) })
	})

	return reader
}

type AppSinkReader struct {
	// synchronization, must not be held while pulling a sample, otherwise Close
	// could not interrupt a blocked reader
	lock   sync.Mutex
	closed bool
	// current is the mapped buffer that Read drains, pending holds the buffers of
	// the pulled sample that are not mapped yet
	current *gst.MapInfo
	pending []*gst.Buffer

	app AppSink

	// eos is closed when the AppSink emits the eos signal. Unlike IsEOS that only
	// happens on a real end of stream, never for a stopped element.
	eos       chan struct{}
	eosOnce   sync.Once
	eosHandle gobject.SignalHandle
}

// Close implements [io.Closer]. It disconnects the eos signal handler, unmaps any
// buffer still held and makes the current and all subsequent reads return
// [io.ErrClosedPipe].
//
// Close is idempotent. It does not change the state of the AppSink and does not
// consume the remaining samples, so the pipeline keeps running.
func (a *AppSinkReader) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.closed {
		return nil
	}

	a.closed = true

	a.app.HandlerDisconnect(a.eosHandle)

	a.unmapCurrent()
	a.pending = nil

	return nil
}

// Read implements [io.Reader]. It blocks until data is available and returns
// [io.EOF] once the AppSink is at EOS and its queue is drained, or
// [io.ErrUnexpectedEOF] if the stream ended without an EOS event.
//
// Consecutive reads resume in the middle of the sample the previous one left off
// at, so a reader must only be used by one goroutine at a time.
func (a *AppSinkReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	for {
		if a.closed {
			return 0, io.ErrClosedPipe
		}

		if a.current != nil {
			n, err = a.current.Read(p)

			if n > 0 {
				return n, nil
			}

			if err != nil && !errors.Is(err, io.EOF) {
				return 0, err
			}

			// the buffer is exhausted, move on to the next one
			a.unmapCurrent()

			continue
		}

		if len(a.pending) > 0 {
			buffer := a.pending[0]
			a.pending = a.pending[1:]

			// an empty buffer maps to a zeroed MapInfo that is not even readable, and
			// it has nothing to contribute to the stream anyway
			if buffer.GetSize() == 0 {
				continue
			}

			info, ok := buffer.Map(gst.MapRead)
			if !ok {
				return 0, ErrBufferNotMappable
			}

			a.current = info

			continue
		}

		// nextBuffers blocks, so no mutex must be held here. Every return path below
		// runs locked again, as the deferred unlock expects.
		a.lock.Unlock()
		buffers, err := a.nextBuffers()
		a.lock.Lock()

		if err != nil {
			return 0, err
		}

		a.pending = buffers
	}
}

// nextBuffers blocks until the AppSink has another sample and returns its buffers.
// It returns [io.EOF] at end of stream, [io.ErrUnexpectedEOF] if the AppSink was
// stopped or never started, and [io.ErrClosedPipe] once the reader is closed.
func (a *AppSinkReader) nextBuffers() ([]*gst.Buffer, error) {
	for {
		if a.isClosed() {
			return nil, io.ErrClosedPipe
		}

		sample := a.app.TryPullSample(appSinkPollInterval)

		if a.isClosed() {
			return nil, io.ErrClosedPipe
		}

		if sample != nil {
			return sampleBuffers(sample), nil
		}

		if !a.app.IsEOS() {
			// the poll interval expired
			continue
		}

		// IsEOS also reports true for an element that was stopped or never started,
		// so only the eos signal proves a clean end of stream.
		if a.awaitEOS() {
			return nil, io.EOF
		}

		return nil, io.ErrUnexpectedEOF
	}
}

// awaitEOS reports whether the AppSink emitted the eos signal, waiting up to
// appSinkEOSGrace for it to arrive.
func (a *AppSinkReader) awaitEOS() bool {
	select {
	case <-a.eos:
		return true
	case <-time.After(appSinkEOSGrace):
		return false
	}
}

func (a *AppSinkReader) isClosed() bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.closed
}

// unmapCurrent releases the buffer that Read drains. a.lock must be held.
func (a *AppSinkReader) unmapCurrent() {
	if a.current != nil {
		a.current.Unmap()
		a.current = nil
	}
}

var _ io.WriterTo = (*AppSinkInstance)(nil)

// WriteTo implements [io.WriterTo] for [AppSink]. Reaching EOS is a successful
// completion and reported with a nil error.
//
// It returns [io.ErrUnexpectedEOF] if the stream ended without an EOS event, which
// is the case for an element that was stopped or never started. Buffers are written
// as they are pulled, so w sees the sample boundaries of the AppSink.
func (a *AppSinkInstance) WriteTo(w io.Writer) (n int64, err error) {
	reader := a.Reader()
	defer reader.Close()

	for {
		buffers, err := reader.nextBuffers()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return n, nil
			}

			return n, err
		}

		for _, buffer := range buffers {
			written, err := writeBuffer(w, buffer)

			n += written

			if err != nil {
				return n, err
			}
		}
	}
}

func writeBuffer(w io.Writer, buffer *gst.Buffer) (n int64, err error) {
	if buffer.GetSize() == 0 {
		return 0, nil
	}

	info, ok := buffer.Map(gst.MapRead)
	if !ok {
		return 0, ErrBufferNotMappable
	}

	defer info.Unmap()

	data := info.Data()

	written, err := w.Write(data)
	if err != nil {
		return int64(written), err
	}

	if written < len(data) {
		return int64(written), io.ErrShortWrite
	}

	return int64(written), nil
}

// sampleBuffers returns the buffers carried by sample in read order. A sample
// holds a single buffer, or a list of them when buffer list support was enabled.
// The buffers borrow from sample, which the bindings keep alive for them.
func sampleBuffers(sample *gst.Sample) []*gst.Buffer {
	if buffer := sample.GetBuffer(); buffer != nil {
		return []*gst.Buffer{buffer}
	}

	list := sample.GetBufferList()
	if list == nil {
		return nil
	}

	buffers := make([]*gst.Buffer, 0, list.Length())

	for i := uint(0); i < list.Length(); i++ {
		if buffer := list.Get(i); buffer != nil {
			buffers = append(buffers, buffer)
		}
	}

	return buffers
}
