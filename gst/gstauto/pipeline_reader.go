package gstauto

import (
	"fmt"

	"github.com/tinyzimmer/go-gst/gst"
)

// Empty assignment to ensure PipelineReader satisfies the ReadPipeliner interface.
var _ ReadPipeliner = &PipelineReader{}

// PipelineReader is the base struct to be used to implement ReadPipeliners.
type PipelineReader struct {
	*readCloser
	pipeline *gst.Pipeline
}

// NewPipelineReader returns a new PipelineReader with an empty pipeline. Use an empty name
// to have gstreamer auto-generate one. This method is intended for use in the construction
// of other interfaces.
func NewPipelineReader(name string) (*PipelineReader, error) {
	pipeline, err := gst.NewPipeline(name)
	if err != nil {
		return nil, err
	}
	rCloser, err := newReadCloser()
	if err != nil {
		if closeErr := pipeline.Destroy(); closeErr != nil {
			fmt.Println("[gst-auto] Failed to destroy errored pipeline:", closeErr.Error())
		}
		return nil, err
	}
	return &PipelineReader{
		readCloser: rCloser,
		pipeline:   pipeline,
	}, nil
}

// NewPipelineReaderFromString returns a new PipelineReader with a pipeline populated
// by the provided gstreamer launch string. If you are looking to build a simple
// ReadPipeliner you probably want to use NewPipelineReaderSimpleFromString.
func NewPipelineReaderFromString(launchStr string) (*PipelineReader, error) {
	pipeline, err := gst.NewPipelineFromString(launchStr)
	if err != nil {
		return nil, err
	}
	rCloser, err := newReadCloser()
	if err != nil {
		if closeErr := pipeline.Destroy(); closeErr != nil {
			fmt.Println("[gst-auto] Failed to destroy errored pipeline:", closeErr.Error())
		}
		return nil, err
	}
	return &PipelineReader{
		readCloser: rCloser,
		pipeline:   pipeline,
	}, nil
}

// Pipeline returns the underlying Pipeline instance for this pipeliner. It implements the
// Pipeliner interface.
func (r *PipelineReader) Pipeline() *gst.Pipeline { return r.pipeline }

// ReaderFd returns the file descriptor that can be written to for the read-buffer. This value
// is used when wanting to allow an underlying pipeline to write to the internal buffer (e.g. when using a fdsink).
func (r *PipelineReader) ReaderFd() int { return int(r.readCloser.rWriter.Fd()) }

// Close will stop and unref the underlying pipeline.
func (r *PipelineReader) Close() error {
	if err := r.Pipeline().Destroy(); err != nil {
		return err
	}
	return r.readCloser.Close()
}

// CloseAsync will close the underlying pipeline asynchronously. It is the caller's
// responsibility to call Unref on the pipeline and close buffers once it is no longer being used.
// This can be accomplished via calling a regular Close (which is idempotent).
func (r *PipelineReader) CloseAsync() error { return r.pipeline.SetState(gst.StateNull) }
