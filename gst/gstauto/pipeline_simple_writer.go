package gstauto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tinyzimmer/go-gst/gst"
)

// PipelineWriterSimple implements a WritePipeliner that configures gstreamer
// to read directly from the internal write-buffer via a fdsrc.
type PipelineWriterSimple struct {
	*PipelineWriter
}

// NewPipelineWriterSimpleFromString returns a new PipelineWriterSimple populated from
// the given launch string. An fdsrc is added to the beginning of the string and tied to
// the write buffer.
func NewPipelineWriterSimpleFromString(launchStr string) (*PipelineWriterSimple, error) {
	pipelineWriter, err := NewPipelineWriterFromString(addFdSrcToStr(launchStr))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if destroyErr := pipelineWriter.Pipeline().Destroy(); destroyErr != nil {
				fmt.Println("[go-gst] Error while destroying failed pipeline instance:", destroyErr.Error())
			}
		}
	}()

	// Retrieve the sources in the pipeline, most of the time there is just one
	var sources []*gst.Element
	sources, err = pipelineWriter.Pipeline().GetSourceElements()
	if err != nil {
		return nil, err
	}

	// Fetch the fdsrc and reconfigure it to point to the write buffer.
	for _, source := range sources {
		if strings.Contains(source.Name(), "fdsrc") {
			if err = source.Set("fd", pipelineWriter.WriterFd()); err != nil {
				return nil, err
			}
		}
	}

	// Return the pipeline
	return &PipelineWriterSimple{pipelineWriter}, nil
}

// NewPipelineWriterSimpleFromConfig returns a new PipelineWriterSimple populated from
// the given launch config. An fdsrc is added to the start of the launch config and tied
// to the write buffer.
func NewPipelineWriterSimpleFromConfig(cfg *PipelineConfig) (*PipelineWriterSimple, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Elements cannot be nil in the config")
	}
	pipelineWriter, err := NewPipelineWriter("")
	if err != nil {
		return nil, err
	}
	cfg.pushPluginToTop(&PipelineElement{
		Name: "fdsrc",
		Data: map[string]interface{}{
			"fd": pipelineWriter.WriterFd(),
		},
	})
	if err := cfg.Apply(pipelineWriter.Pipeline()); err != nil {
		if destroyErr := pipelineWriter.Pipeline().Destroy(); destroyErr != nil {
			fmt.Println("[go-gst] Error while destroying failed pipeline instance:", destroyErr.Error())
		}
		return nil, err
	}
	return &PipelineWriterSimple{pipelineWriter}, nil
}

func addFdSrcToStr(pstr string) string {
	if pstr == "" {
		return "fdsrc"
	}
	return fmt.Sprintf("fdsink ! %s", pstr)
}
