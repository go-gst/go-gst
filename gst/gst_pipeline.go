package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -Wno-incompatible-pointer-types -g
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"
)

// PipelineFlags represents arguments passed to a new Pipeline.
type PipelineFlags int

const (
	// PipelineInternalOnly signals that this pipeline only handles data internally.
	PipelineInternalOnly PipelineFlags = 1 << iota
	// PipelineRead signals that the Read() method can be used on the end of this pipeline.
	PipelineRead
	// PipelineWrite signals that the Write() method can be used on the start of this pipeline.
	PipelineWrite
	// PipelineUseGstApp signals the desire to use an AppSink or AppSrc instead of the default
	// os pipes, fdsrc, and fdsink.
	// When using this flag, you should interact with the pipeline using the GetAppSink and
	// GetAppSrc methods.
	PipelineUseGstApp
	// PipelineReadWrite signals that this pipeline can be both read and written to.
	PipelineReadWrite = PipelineRead | PipelineWrite
)

// has returns true if these flags contain the given flag.
func (p PipelineFlags) has(b PipelineFlags) bool { return p&b != 0 }

// State is a type cast of the C GstState
type State int

// Type casting for GstStates
const (
	VoidPending  State = C.GST_STATE_VOID_PENDING // (0) – no pending state.
	StateNull          = C.GST_STATE_NULL         // (1) – the NULL state or initial state of an element.
	StateReady         = C.GST_STATE_READY        // (2) – the element is ready to go to PAUSED.
	StatePaused        = C.GST_STATE_PAUSED       // (3) – the element is PAUSED, it is ready to accept and process data. Sink elements however only accept one buffer and then block.
	StatePlaying       = C.GST_STATE_PLAYING      // (4) – the element is PLAYING, the GstClock is running and the data is flowing.
)

func (s State) String() string {
	return C.GoString(C.gst_element_state_get_name((C.GstState)(s)))
}

// Pipeline is the base implementation of a GstPipeline using CGO to wrap
// gstreamer API calls. It provides methods to be inherited by the extending
// PlaybackPipeline and RecordingPipeline objects. The struct itself implements
// a ReadWriteCloser.
type Pipeline struct {
	*Bin

	// a local reference to the bus so duplicates aren't created
	// when retrieved by the user
	bus *Bus

	// The buffers backing the Read and Write methods
	destBuf *bufio.Reader
	srcBuf  *bufio.Writer

	// used with PipelineWrite
	srcReader, srcWriter *os.File
	// used with PipelineRead
	destReader, destWriter *os.File

	// used with PipelineWrite AND PipelineGstApp
	appSrc *AppSrc
	// used with PipelineRead AND PipelineGstApp
	appSink   *AppSink
	autoFlush bool // when set to true, the contents of the app sink are automatically flushed to the read buffer.

	// The element that represents the source/dest pipeline
	// and any caps to apply to it.
	srcElement  *Element
	srcCaps     Caps
	destElement *Element

	// whether or not the pipeline was built from a string. this is checked when
	// starting to see who is responsible for build and linking the buffers.
	pipelineFromHelper bool

	// A channel where a caller can listen for errors asynchronously.
	errCh chan error
	// A channel where a caller can listen for messages
	msgCh []chan *Message
}

func newEmptyPipeline() (*C.GstPipeline, error) {
	pipeline := C.gst_pipeline_new((*C.gchar)(nil))
	if pipeline == nil {
		return nil, errors.New("Could not create new pipeline")
	}
	return C.toGstPipeline(unsafe.Pointer(pipeline)), nil
}

func newPipelineFromString(launchv string) (*C.GstPipeline, error) {
	if len(strings.Split(launchv, "!")) < 2 {
		return nil, fmt.Errorf("Given string is too short for a pipeline: %s", launchv)
	}
	cLaunchv := C.CString(launchv)
	defer C.free(unsafe.Pointer(cLaunchv))
	var gerr *C.GError
	pipeline := C.gst_parse_launch((*C.gchar)(cLaunchv), (**C.GError)(&gerr))
	if gerr != nil {
		defer C.g_error_free((*C.GError)(gerr))
		errMsg := C.GoString(gerr.message)
		return nil, errors.New(errMsg)
	}
	return C.toGstPipeline(unsafe.Pointer(pipeline)), nil
}

// NewPipeline builds and returns a new empty Pipeline instance.
func NewPipeline(flags PipelineFlags) (*Pipeline, error) {
	pipelineElement, err := newEmptyPipeline()
	if err != nil {
		return nil, err
	}

	pipeline := wrapPipeline(pipelineElement)

	if err := applyFlags(pipeline, flags); err != nil {
		return nil, err
	}

	return pipeline, nil
}

func applyFlags(pipeline *Pipeline, flags PipelineFlags) error {
	// If the user wants to be able to write to the pipeline, set up the
	// write-buffers
	if flags.has(PipelineWrite) {
		// Set up a pipe
		if err := pipeline.setupWriters(); err != nil {
			return err
		}
	}

	// If the user wants to be able to read from the pipeline, setup the
	// read-buffers.
	if flags.has(PipelineRead) {
		if err := pipeline.setupReaders(); err != nil {
			return err
		}
	}

	return nil
}

func wrapPipeline(elem *C.GstPipeline) *Pipeline { return &Pipeline{Bin: wrapBin(&elem.bin)} }

func (p *Pipeline) setupWriters() error {
	var err error
	p.srcReader, p.srcWriter, err = os.Pipe()
	if err != nil {
		return err
	}
	p.srcBuf = bufio.NewWriter(p.srcWriter)
	return nil
}

func (p *Pipeline) setupReaders() error {
	var err error
	p.destReader, p.destWriter, err = os.Pipe()
	if err != nil {
		return err
	}
	p.destBuf = bufio.NewReader(p.destReader)
	return nil
}

// Instance returns the native GstPipeline instance.
func (p *Pipeline) Instance() *C.GstPipeline { return C.toGstPipeline(p.unsafe()) }

// Read implements a Reader and returns data from the read buffer.
func (p *Pipeline) Read(b []byte) (int, error) {
	if p.destBuf == nil {
		return 0, io.ErrClosedPipe
	}
	return p.destBuf.Read(b)
}

// readerFd returns the file descriptor for the read buffer, or 0 if
// there isn't one. It returns the file descriptor that can be written to
// by gstreamer.
func (p *Pipeline) readerFd() uintptr {
	if p.destWriter == nil {
		return 0
	}
	return p.destWriter.Fd()
}

// Write implements a Writer and places data in the write buffer.
func (p *Pipeline) Write(b []byte) (int, error) {
	if p.srcBuf == nil {
		return 0, io.ErrClosedPipe
	}
	return p.srcBuf.Write(b)
}

// writerFd returns the file descriptor for the write buffer, or 0 if
// there isn't one. It returns the file descriptor that can be read from
// by gstreamer.
func (p *Pipeline) writerFd() uintptr {
	if p.srcWriter == nil {
		return 0
	}
	return p.srcReader.Fd()
}

// SetWriterCaps sets the caps on the write-buffer. You will usually want to call this
// on a custom pipeline, unless you are using downstream elements that do dynamic pad
// linking.
func (p *Pipeline) SetWriterCaps(caps Caps) { p.srcCaps = caps }

// LinkWriterTo links the write buffer on this Pipeline to the given element. This must
// be called when the pipeline is constructed with PipelineWrite or PipelineReadWrite.
func (p *Pipeline) LinkWriterTo(elem *Element) { p.srcElement = elem }

// LinkReaderTo links the read buffer on this Pipeline to the given element. This must
// be called when the pipeline is constructed with PipelineRead or PipelineReadWrite.
func (p *Pipeline) LinkReaderTo(elem *Element) { p.destElement = elem }

// IsUsingGstApp returns true if the current pipeline is using GstApp instead of file descriptors.
func (p *Pipeline) IsUsingGstApp() bool {
	return p.appSrc != nil || p.appSink != nil
}

// GetAppSrc returns the AppSrc for this pipeline if created with PipelineUseGstApp.
// Unref after usage.
func (p *Pipeline) GetAppSrc() *AppSrc {
	if p.appSrc == nil {
		return nil
	}
	// increases the ref count on the element
	return wrapAppSrc(p.appSrc.Element)
}

// GetAppSink returns the AppSink for this pipeline if created with PipelineUseGstApp.
// Unref after usage.
func (p *Pipeline) GetAppSink() *AppSink {
	if p.appSink == nil {
		return nil
	}
	// increases the ref count
	return wrapAppSink(p.appSink.Element)
}

// GetBus returns the message bus for this pipeline.
func (p *Pipeline) GetBus() *Bus {
	if p.bus == nil {
		cBus := C.gst_pipeline_get_bus((*C.GstPipeline)(p.Instance()))
		p.bus = wrapBus(cBus)
	}
	return p.bus
}

// SetAutoFlush sets whether or not samples should be automatically flushed to the read-buffer
// (default for pipelines not built with PipelineUseGstApp) and if messages should be flushed
// on the bus when the pipeline is stopped.
func (p *Pipeline) SetAutoFlush(b bool) {
	p.Set("auto-flush-bus", b)
	p.autoFlush = b
}

// AutoFlush returns true if the pipeline is using a GstAppSink and is configured to autoflush to the
// read-buffer.
func (p *Pipeline) AutoFlush() bool { return p.IsUsingGstApp() && p.autoFlush }

// Flush flushes the app sink to the read buffer. It is usually more desirable to interface
// with the PullSample and BlockPullSample methods on the AppSink interface directly. Or
// to set autoflush to true.
func (p *Pipeline) Flush() error {
	sample, err := p.appSink.PullSample()
	if err != nil { // err signals end of stream
		return err
	}
	if sample == nil {
		return nil
	}
	defer sample.Unref()
	if _, err := io.Copy(p.destWriter, sample.GetBuffer()); err != nil {
		return err
	}
	return nil
}

// BlockFlush is like Flush but it blocks until a sample is available. This is intended for
// use with PipelineUseGstApp.
func (p *Pipeline) BlockFlush() error {
	sample, err := p.appSink.BlockPullSample()
	if err != nil { // err signals end of stream
		return err
	}
	if sample == nil {
		return nil
	}
	defer sample.Unref()
	if _, err := io.Copy(p.destWriter, sample.GetBuffer()); err != nil {
		return err
	}
	return nil
}

// setupSrc sets up a source element with the given configuration.
func (p *Pipeline) setupSrc(pluginName string, args map[string]interface{}) (*Element, error) {
	elem, err := NewElement(pluginName)
	if err != nil {
		return nil, err
	}
	for k, v := range args {
		if err := elem.Set(k, v); err != nil {
			return nil, err
		}
	}
	if err := p.Add(elem); err != nil {
		return nil, err
	}
	if p.srcCaps != nil {
		return elem, elem.LinkFiltered(p.srcElement, p.srcCaps)
	}
	return elem, elem.Link(p.srcElement)
}

// setupFdSrc will setup a fdsrc as the source of the pipeline.
func (p *Pipeline) setupFdSrc() error {
	_, err := p.setupSrc("fdsrc", map[string]interface{}{
		"fd": p.writerFd(),
	})
	return err
}

// setupAppSrc sets up an appsrc as the source of the pipeline
func (p *Pipeline) setupAppSrc() error {
	appSrc, err := p.setupSrc("appsrc", map[string]interface{}{
		"block":        true,  // TODO: make this configurable
		"emit-signals": false, // https://gstreamer.freedesktop.org/documentation/app/appsrc.html?gi-language=c
	})
	if err != nil {
		return err
	}
	p.appSrc = &AppSrc{appSrc}
	return nil
}

// setupSrcElement will setup the source element when the pipeline is constructed with
// PipelineWrite.
func (p *Pipeline) setupSrcElement() error {
	if p.srcElement == nil {
		return errors.New("Pipeline was constructed with PipelineWrite but LinkWriterTo was never called")
	}
	if p.IsUsingGstApp() {
		return p.setupAppSrc()
	}
	return p.setupFdSrc()
}

// setupSink sets up a sink element with the given congifuration.
func (p *Pipeline) setupSink(pluginName string, args map[string]interface{}) (*Element, error) {
	elem, err := NewElement(pluginName)
	if err != nil {
		return nil, err
	}
	for k, v := range args {
		if err := elem.Set(k, v); err != nil {
			return nil, err
		}
	}
	if err := p.Add(elem); err != nil {
		return nil, err
	}
	return elem, p.destElement.Link(elem)
}

// setupFdSink sets up a fdsink as the sink of the pipeline.
func (p *Pipeline) setupFdSink() error {
	_, err := p.setupSink("fdsink", map[string]interface{}{
		"fd": p.readerFd(),
	})
	return err
}

// setupAppSink sets up an appsink as the sink of the pipeline.
func (p *Pipeline) setupAppSink() error {
	appSink, err := p.setupSink("appsink", map[string]interface{}{
		"emit-signals": false,
	})
	if err != nil {
		return err
	}
	p.appSink = wrapAppSink(appSink)
	return nil
}

// setupDestElement will setup the destination (sink) element when the pipeline is constructed with
// PipelineRead.
func (p *Pipeline) setupDestElement() error {
	if p.destElement == nil {
		return errors.New("Pipeline was constructed with PipelineRead but LinkReaderTo was never called")
	}
	if p.IsUsingGstApp() {
		return p.setupAppSink()
	}
	return p.setupFdSink()
}

// Start will start the GstPipeline. It is asynchronous so it does not need to be
// called within a goroutine, however, it is still safe to do so.
func (p *Pipeline) Start() error {
	// If there is a write buffer on this pipeline, set up an fdsrc
	if p.srcBuf != nil && !p.pipelineFromHelper {
		if err := p.setupSrcElement(); err != nil {
			return err
		}
	}

	// If there is a read buffer on this pipeline, set up an fdsink
	if p.destBuf != nil && !p.pipelineFromHelper {
		if err := p.setupDestElement(); err != nil {
			return err
		}
	}

	return p.startPipeline()
}

func (p *Pipeline) closeBuffers() error {
	if p.srcBuf != nil && p.srcReader != nil && p.srcWriter != nil {
		if err := p.srcReader.Close(); err != nil {
			return err
		}
		if err := p.srcWriter.Close(); err != nil {
			return err
		}
		p.srcBuf = nil
	}
	if p.destBuf != nil && p.destReader != nil && p.destWriter != nil {
		if err := p.destReader.Close(); err != nil {
			return err
		}
		if err := p.destWriter.Close(); err != nil {
			return err
		}
		p.destBuf = nil
	}
	return nil
}

// ReadBufferSize returns the current size of the unread portion of the read-buffer.
func (p *Pipeline) ReadBufferSize() int {
	if p.destBuf == nil {
		return 0
	}
	return p.destBuf.Buffered()
}

// WriteBufferSize returns the current size of the unread portion of the write-buffer.
func (p *Pipeline) WriteBufferSize() int {
	if p.srcBuf == nil {
		return 0
	}
	return p.srcBuf.Buffered()
}

// TotalBufferSize returns the sum of the Read and Write buffer unread portions.
func (p *Pipeline) TotalBufferSize() int { return p.WriteBufferSize() + p.ReadBufferSize() }

// Close implements a Closer and closes all buffers.
func (p *Pipeline) Close() error {
	defer p.Unref()
	if err := p.closeBuffers(); err != nil {
		return err
	}
	return p.SetState(StateNull)
}

// startPipeline will set the GstPipeline to the PLAYING state.
func (p *Pipeline) startPipeline() error {
	if err := p.SetState(StatePlaying); err != nil {
		return err
	}
	// If using GstApp with autoflush
	if p.AutoFlush() {
		go func() {
			for {
				if err := p.BlockFlush(); err != nil {
					// err signals end of stream
					return
				}
			}
		}()
	}
	return nil
}

// Wait waits for the given pipeline to reach end of stream.
func Wait(p *Pipeline) {
	if p.Instance() == nil {
		return
	}
	msgCh := p.GetBus().MessageChan()
	for {
		select {
		default:
			if p.Instance() == nil || p.GetState() == StateNull {
				return
			}
		case msg := <-msgCh:
			defer msg.Unref()
			switch msg.Type() {
			case MessageEOS:
				return
			}
		}
	}
}
