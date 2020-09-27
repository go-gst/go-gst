package gstauto

import (
	"fmt"

	"github.com/tinyzimmer/go-gst/gst"
)

// Blank assignment to make sure PipelinerSimple satisfies the Pipeliner interface.
var _ Pipeliner = &PipelinerSimple{}

// PipelinerSimple is a simple struct that implements the Pipeliner interface.
// It doesn't provide any additional read/write capabilities. Its primary intention
// is for pipelines where the caller does not wish to personally read or write from
// either end of the buffer.
type PipelinerSimple struct {
	pipeline *gst.Pipeline
}

// Pipeline implements the Pipeliner interface.
func (s *PipelinerSimple) Pipeline() *gst.Pipeline { return s.pipeline }

// NewPipelinerSimple returns a new empty PipelinerSimple. Pass an empty string
// for name to use an auto-generated one.
func NewPipelinerSimple(name string) (*PipelinerSimple, error) {
	pipeline, err := gst.NewPipeline(name)
	if err != nil {
		return nil, err
	}
	return &PipelinerSimple{pipeline: pipeline}, nil
}

// NewPipelinerSimpleFromString returns a new PipelinerSimpler from the given
// launch string.
func NewPipelinerSimpleFromString(launchStr string) (*PipelinerSimple, error) {
	pipeline, err := gst.NewPipelineFromString(launchStr)
	if err != nil {
		return nil, err
	}
	return &PipelinerSimple{pipeline: pipeline}, nil
}

// NewPipelinerSimpleFromConfig returns a new PipelinerSimple from the given
// PipelineConfig.
func NewPipelinerSimpleFromConfig(cfg *PipelineConfig) (*PipelinerSimple, error) {
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}
	if err := cfg.Apply(pipeline); err != nil {
		if destroyErr := pipeline.Destroy(); destroyErr != nil {
			fmt.Println("[go-gst] Error while destroying failed pipeline instance:", destroyErr.Error())
		}
		return nil, err
	}
	return &PipelinerSimple{pipeline: pipeline}, nil
}
