package gstauto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tinyzimmer/go-gst/gst"
)

// PipelineWriterApp implements a WritePipeliner that configures gstreamer
// with an appsrc. The appsrc allows for more granular control over the data
// at the start of the pipeline.
type PipelineWriterApp struct {
	*PipelineWriter

	appSrc *gst.AppSrc
}

// NewPipelineWriterAppFromString returns a new PipelineWriterApp populated from
// the given launch string. An appsrc is added to the start of the launch string and made
// available via the GetAppSource method.
func NewPipelineWriterAppFromString(launchStr string) (*PipelineWriterApp, error) {
	pipelineWriter, err := NewPipelineWriterFromString(addAppSourceToStr(launchStr))
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

	appPipeline := &PipelineWriterApp{PipelineWriter: pipelineWriter}

	// Retrieve the sources in the pipeline, most of the time there is just one
	var sources []*gst.Element
	sources, err = pipelineWriter.Pipeline().GetSourceElements()
	if err != nil {
		return nil, err
	}

	// Fetch the appsrc and make a local reference to it
	for _, src := range sources {
		if strings.Contains(src.Name(), "appsrc") {
			appPipeline.appSrc = &gst.AppSrc{Element: src}
		}
	}

	// Return the pipeline
	return appPipeline, nil
}

// NewPipelineWriterAppFromConfig returns a new PipelineWriterApp populated from
// the given launch config. An appsrc is added to the start of the launch config and
// made available via the GetAppSource method.
func NewPipelineWriterAppFromConfig(cfg *PipelineConfig) (*PipelineWriterApp, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Elements cannot be nil in the config")
	}
	pipelineWriter, err := NewPipelineWriter("")
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

	cfg.pushPluginToTop(&PipelineElement{Name: "appsrc"})

	if err = cfg.Apply(pipelineWriter.Pipeline()); err != nil {
		return nil, err
	}

	appPipeline := &PipelineWriterApp{PipelineWriter: pipelineWriter}

	// Retrieve the sources in the pipeline, most of the time there is just one
	var sources []*gst.Element
	sources, err = pipelineWriter.Pipeline().GetSourceElements()
	if err != nil {
		return nil, err
	}

	// Fetch the appsrc and make a local reference to it
	for _, src := range sources {
		if strings.Contains(src.Name(), "appsrc") {
			appPipeline.appSrc = &gst.AppSrc{Element: src}
		}
	}

	return appPipeline, nil
}

func addAppSourceToStr(pstr string) string {
	if pstr == "" {
		return "appsrc"
	}
	return fmt.Sprintf("%s ! appsrc", pstr)
}

// GetAppSource returns the app src for this pipeline.
func (p *PipelineWriterApp) GetAppSource() *gst.AppSrc { return p.appSrc }
