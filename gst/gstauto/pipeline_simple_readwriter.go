package gstauto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tinyzimmer/go-gst/gst"
)

// PipelineReadWriterSimple implements a ReadWritePipeliner that configures gstreamer
// to read from the internal write-buffer via an fdsrc and write to the internal read-buffer
// via an fdsink.
type PipelineReadWriterSimple struct {
	*PipelineReadWriter
}

// NewPipelineReadWriterSimpleFromString returns a new PipelineReadWriterSimple from
// the given launch string. An fdsrc listening on the write buffer and an fdsink to the read buffer
// are formatted into the provided string.
func NewPipelineReadWriterSimpleFromString(launchStr string) (*PipelineReadWriterSimple, error) {
	pipelineReadWriter, err := NewPipelineReadWriterFromString(addFdSrcToStr(addFdSinkToStr(launchStr)))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if destroyErr := pipelineReadWriter.Pipeline().Destroy(); destroyErr != nil {
				fmt.Println("[go-gst] Error while destroying failed pipeline instance:", destroyErr.Error())
			}
		}
	}()

	// Retrieve the sinks in the pipeline, most of the time there is just one
	var sinks []*gst.Element
	sinks, err = pipelineReadWriter.Pipeline().GetSinkElements()
	if err != nil {
		return nil, err
	}

	// Fetch the fdsink and reconfigure it to point to the read buffer.
	for _, sink := range sinks {
		if strings.Contains(sink.Name(), "fdsink") {
			if err = sink.Set("fd", pipelineReadWriter.ReaderFd()); err != nil {
				return nil, err
			}
		}
	}

	// Retrieve the sources in the pipeline, most of the time there is just one
	var sources []*gst.Element
	sources, err = pipelineReadWriter.Pipeline().GetSourceElements()
	if err != nil {
		return nil, err
	}

	// Fetch the fdsrc and reconfigure it to point to the write buffer.
	for _, source := range sources {
		if strings.Contains(source.Name(), "fdsrc") {
			if err = source.Set("fd", pipelineReadWriter.WriterFd()); err != nil {
				return nil, err
			}
		}
	}

	// Return the pipeline
	return &PipelineReadWriterSimple{pipelineReadWriter}, nil
}

// NewPipelineReadWriterSimpleFromConfig returns a new PipelineReadWriterSimple populated from
// the given launch config. An fdsrc is added to the start of the launch config and tied
// to the write buffer, and an fdsink is added to the end tied to the read-buffer.
func NewPipelineReadWriterSimpleFromConfig(cfg *PipelineConfig) (*PipelineReadWriterSimple, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Elements cannot be nil in the config")
	}
	pipelineReadWriter, err := NewPipelineReadWriter("")
	if err != nil {
		return nil, err
	}
	cfg.pushPluginToTop(&PipelineElement{
		Name: "fdsrc",
		Data: map[string]interface{}{
			"fd": pipelineReadWriter.WriterFd(),
		},
	})
	cfg.Elements = append(cfg.Elements, &PipelineElement{
		Name: "fdsink",
		Data: map[string]interface{}{
			"fd": pipelineReadWriter.ReaderFd(),
		},
	})
	if err := cfg.Apply(pipelineReadWriter.Pipeline()); err != nil {
		if destroyErr := pipelineReadWriter.Pipeline().Destroy(); destroyErr != nil {
			fmt.Println("[go-gst] Error while destroying failed pipeline instance:", destroyErr.Error())
		}
		return nil, err
	}
	return &PipelineReadWriterSimple{pipelineReadWriter}, nil
}
