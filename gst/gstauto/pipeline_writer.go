package gstauto

import (
	"github.com/tinyzimmer/go-gst/gst"
)

// Empty assignment to ensure PipelineWriter satisfies the WritePipeliner interface.
var _ WritePipeliner = &PipelineWriter{}

// PipelineWriter is the base struct to be used to implement WritePipeliners.
type PipelineWriter struct {
	*writeCloser
	pipeline *gst.Pipeline
}

// NewPipelineWriter returns a new PipelineWriter with an empty pipeline. Use an empty name
// to have gstreamer auto-generate one. This method is intended for use in the construction
// of other interfaces.
func NewPipelineWriter(name string) (*PipelineWriter, error) {
	pipeline, err := gst.NewPipeline(name)
	if err != nil {
		return nil, err
	}
	wCloser, err := newWriteCloser()
	if err != nil {
		runOrPrintErr(pipeline.Destroy)
		return nil, err
	}
	return &PipelineWriter{
		writeCloser: wCloser,
		pipeline:    pipeline,
	}, nil
}

// NewPipelineWriterFromString returns a new PipelineWriter with a pipeline populated
// by the provided gstreamer launch string. If you are looking to build a simple
// WritePipeliner you probably want to use NewPipelineWriterSimpleFromString.
func NewPipelineWriterFromString(launchStr string) (*PipelineWriter, error) {
	pipeline, err := gst.NewPipelineFromString(launchStr)
	if err != nil {
		return nil, err
	}
	wCloser, err := newWriteCloser()
	if err != nil {
		runOrPrintErr(pipeline.Destroy)
		return nil, err
	}
	return &PipelineWriter{
		writeCloser: wCloser,
		pipeline:    pipeline,
	}, nil
}

// Pipeline returns the underlying Pipeline instance for this pipeliner. It implements the
// Pipeliner interface.
func (w *PipelineWriter) Pipeline() *gst.Pipeline { return w.pipeline }

// Start sets the underlying Pipeline state to PLAYING.
func (w *PipelineWriter) Start() error { return w.Pipeline().Start() }

// WriterFd returns the file descriptor that can be used to read from the write-buffer. This value
// is used when wanting to allow an underlying pipeline the ability to read data written to
// the buffer (e.g. when using a fdsrc).
func (w *PipelineWriter) WriterFd() int { return int(w.writeCloser.wReader.Fd()) }

// Close will stop and unref the underlying pipeline.
func (w *PipelineWriter) Close() error {
	if err := w.Pipeline().Destroy(); err != nil {
		return err
	}
	return w.writeCloser.Close()
}

// CloseAsync will close the underlying pipeline asynchronously. It is the caller's
// responsibility to call Unref on the pipeline and close buffers once it is no longer being used.
// This can be accomplished via calling a regular Close (which is idempotent).
func (w *PipelineWriter) CloseAsync() error { return w.pipeline.SetState(gst.StateNull) }
