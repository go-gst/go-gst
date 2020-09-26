package gstauto

import (
	"io"

	"github.com/tinyzimmer/go-gst-launch/gst"
)

// Pipeliner is a the base interface for structs extending the functionality of
// the Pipeline object. It provides a single method which returns the underlying
// Pipeline object.
type Pipeliner interface {
	Pipeline() *gst.Pipeline
}

// ReadPipeliner is a Pipeliner that also implements a ReadCloser.
type ReadPipeliner interface {
	Pipeliner
	io.ReadCloser
}

// WritePipeliner is a Pipeliner that also implements a WriteCloser.
type WritePipeliner interface {
	Pipeliner
	io.WriteCloser
}

// ReadWritePipeliner is a Pipeliner that also implements a ReadWriteCloser.
type ReadWritePipeliner interface {
	ReadPipeliner
	WritePipeliner
}
