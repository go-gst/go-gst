package gstauto

import (
	"io"

	"github.com/tinyzimmer/go-gst/gst"
)

// Pipeliner is a the base interface for structs extending the functionality of
// the Pipeline object. It provides a single method which returns the underlying
// Pipeline object.
type Pipeliner interface {
	io.Closer
	// Pipeline should return the underlying pipeline
	Pipeline() *gst.Pipeline
	// Start should start the underlying pipeline.
	Start() error
}

// ReadPipeliner is a Pipeliner that also implements a ReadCloser.
type ReadPipeliner interface {
	Pipeliner
	io.Reader
}

// WritePipeliner is a Pipeliner that also implements a WriteCloser.
type WritePipeliner interface {
	Pipeliner
	io.Writer
}

// ReadWritePipeliner is a Pipeliner that also implements a ReadWriteCloser.
type ReadWritePipeliner interface {
	Pipeliner
	io.ReadWriter
}
