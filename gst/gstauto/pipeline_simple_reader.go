package gstauto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tinyzimmer/go-gst/gst"
)

// PipelineReaderSimple implements a ReadPipeliner that configures gstreamer
// to write directly to the internal read-buffer via an fdsink.
type PipelineReaderSimple struct {
	*PipelineReader
}

// NewPipelineReaderSimpleFromString returns a new PipelineReaderSimple populated from
// the given launch string. An fdsink is added to the end of the launch string and tied
// to the read buffer.
func NewPipelineReaderSimpleFromString(launchStr string) (*PipelineReaderSimple, error) {
	pipelineReader, err := NewPipelineReaderFromString(addFdSinkToStr(launchStr))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if destroyErr := pipelineReader.Pipeline().Destroy(); destroyErr != nil {
				fmt.Println("[go-gst] Error while destroying failed pipeline instance:", destroyErr.Error())
			}
		}
	}()

	// Retrieve the sinks in the pipeline, most of the time there is just one
	var sinks []*gst.Element
	sinks, err = pipelineReader.Pipeline().GetSinkElements()
	if err != nil {
		return nil, err
	}

	// Fetch the fdsink and reconfigure it to point to the read buffer.
	for _, sink := range sinks {
		if strings.Contains(sink.Name(), "fdsink") {
			if err = sink.Set("fd", pipelineReader.ReaderFd()); err != nil {
				return nil, err
			}
		}
	}

	// Return the pipeline
	return &PipelineReaderSimple{pipelineReader}, nil
}

// NewPipelineReaderSimpleFromConfig returns a new PipelineReaderSimple populated from
// the given launch config. An fdsink is added to the end of the launch config and tied
// to the read buffer.
func NewPipelineReaderSimpleFromConfig(cfg *PipelineConfig) (*PipelineReaderSimple, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Elements cannot be nil in the config")
	}
	pipelineReader, err := NewPipelineReader("")
	if err != nil {
		return nil, err
	}
	cfg.Elements = append(cfg.Elements, &PipelineElement{
		Name: "fdsink",
		Data: map[string]interface{}{
			"fd": pipelineReader.ReaderFd(),
		},
	})
	if err := cfg.Apply(pipelineReader.Pipeline()); err != nil {
		if destroyErr := pipelineReader.Pipeline().Destroy(); destroyErr != nil {
			fmt.Println("[go-gst] Error while destroying failed pipeline instance:", destroyErr.Error())
		}
		return nil, err
	}
	return &PipelineReaderSimple{pipelineReader}, nil
}

func addFdSinkToStr(pstr string) string {
	if pstr == "" {
		return "fdsink"
	}
	return fmt.Sprintf("%s ! fdsink", pstr)
}
