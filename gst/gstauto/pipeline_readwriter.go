package gstauto

import (
	"fmt"

	"github.com/tinyzimmer/go-gst/gst"
)

// Empty assignment to ensure PipelineReadWriter satisfies the ReadWritePipeliner interface.
var _ ReadWritePipeliner = &PipelineReadWriter{}

// PipelineReadWriter is the base struct to be used to implement ReadWritePipeliners.
type PipelineReadWriter struct {
	*readWriteCloser
	pipeline *gst.Pipeline
}

// NewPipelineReadWriter returns a new PipelineReadWriter with an empty pipeline. Use an empty name
// to have gstreamer auto-generate one. This method is intended for use in the construction
// of other interfaces.
func NewPipelineReadWriter(name string) (*PipelineReadWriter, error) {
	pipeline, err := gst.NewPipeline(name)
	if err != nil {
		return nil, err
	}
	rwCloser, err := newReadWriteCloser()
	if err != nil {
		if closeErr := pipeline.Destroy(); closeErr != nil {
			fmt.Println("[gst-auto] Failed to destroy errored pipeline:", closeErr.Error())
		}
		return nil, err
	}
	return &PipelineReadWriter{
		readWriteCloser: rwCloser,
		pipeline:        pipeline,
	}, nil
}

// NewPipelineReadWriterFromString returns a new PipelineReadWriter with a pipeline populated
// by the provided gstreamer launch string. If you are looking to build a simple
// ReadWritePipeliner you probably want to use NewPipelineReadWriterSimpleFromString.
func NewPipelineReadWriterFromString(launchStr string) (*PipelineReadWriter, error) {
	pipeline, err := gst.NewPipelineFromString(launchStr)
	if err != nil {
		return nil, err
	}
	rwCloser, err := newReadWriteCloser()
	if err != nil {
		if closeErr := pipeline.Destroy(); closeErr != nil {
			fmt.Println("[gst-auto] Failed to destroy errored pipeline:", closeErr.Error())
		}
		return nil, err
	}
	return &PipelineReadWriter{
		readWriteCloser: rwCloser,
		pipeline:        pipeline,
	}, nil
}

// Pipeline returns the underlying Pipeline instance for this pipeliner. It implements the
// Pipeliner interface.
func (rw *PipelineReadWriter) Pipeline() *gst.Pipeline { return rw.pipeline }

// ReaderFd returns the file descriptor that can be written to for the read-buffer. This value
// is used when wanting to allow an underlying pipeline to write to the internal buffer
// (e.g. when using a fdsink).
func (rw *PipelineReadWriter) ReaderFd() uintptr { return rw.readWriteCloser.readCloser.rWriter.Fd() }

// WriterFd returns the file descriptor that can be used to read from the write-buffer. This value
// is used when wanting to allow an underlying pipeline the ability to read data written to the buffer
// (e.g. when using a fdsrc).
func (rw *PipelineReadWriter) WriterFd() uintptr { return rw.readWriteCloser.writeCloser.wReader.Fd() }

// Close will stop and unref the underlying pipeline and read/write buffers.
func (rw *PipelineReadWriter) Close() error {
	if err := rw.Pipeline().Destroy(); err != nil {
		return err
	}
	return rw.readWriteCloser.Close()
}

// CloseAsync will close the underlying pipeline asynchronously. It is the caller's
// responsibility to call Unref on the pipeline and close buffers once it is no longer being used.
// This can be accomplished via calling a regular Close (which is idempotent).
func (rw *PipelineReadWriter) CloseAsync() error { return rw.pipeline.SetState(gst.StateNull) }
