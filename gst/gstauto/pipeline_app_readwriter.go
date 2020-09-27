package gstauto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tinyzimmer/go-gst/gst"
)

// PipelineReadWriterApp implements a ReadPipeliner that configures gstreamer
// with an appsink. The appsink allows for more granular control over the data
// at the end of the pipeline, and the appsrc allows for control over the data
// at the start.
type PipelineReadWriterApp struct {
	*PipelineReadWriter

	appSrc  *gst.AppSrc
	appSink *gst.AppSink
}

// NewPipelineReadWriterAppFromString returns a new PipelineReadWriterApp populated from
// the given launch string. An appsink is added to the end of the launch string and made
// available via the GetAppSink method, and an appsrc is added to the end and made
// available via the GetAppSource method.
func NewPipelineReadWriterAppFromString(launchStr string) (*PipelineReadWriterApp, error) {
	pipelineReadWriter, err := NewPipelineReadWriterFromString(addAppSourceToStr(addAppSinkToStr(launchStr)))
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

	appPipeline := &PipelineReadWriterApp{PipelineReadWriter: pipelineReadWriter}

	// Retrieve the sources in the pipeline, most of the time there is just one
	var sources []*gst.Element
	sources, err = pipelineReadWriter.Pipeline().GetSourceElements()
	if err != nil {
		return nil, err
	}

	// Fetch the appsrc and make a local reference to it
	for _, src := range sources {
		if strings.Contains(src.Name(), "appsrc") {
			appPipeline.appSrc = &gst.AppSrc{Element: src}
		}
	}

	// Retrieve the sinks in the pipeline, most of the time there is just one
	var sinks []*gst.Element
	sinks, err = pipelineReadWriter.Pipeline().GetSinkElements()
	if err != nil {
		return nil, err
	}

	// Fetch the appsink and make a local reference to it
	for _, sink := range sinks {
		if strings.Contains(sink.Name(), "appsink") {
			appPipeline.appSink = &gst.AppSink{Element: sink}
		}
	}

	// Return the pipeline
	return appPipeline, nil
}

// NewPipelineReadWriterAppFromConfig returns a new PipelineReadWriterApp populated from
// the given launch config. An appsink is added to the end of the launch config and
// made available via the GetAppSink method, and an appsrc is added at the front and made
// available via the GetAppSource method.
func NewPipelineReadWriterAppFromConfig(cfg *PipelineConfig) (*PipelineReadWriterApp, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Elements cannot be nil in the config")
	}
	pipelineReadWriter, err := NewPipelineReadWriter("")
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

	cfg.Elements = append(cfg.Elements, &PipelineElement{Name: "appsink"})

	if err = cfg.Apply(pipelineReadWriter.Pipeline()); err != nil {
		return nil, err
	}

	appPipeline := &PipelineReadWriterApp{PipelineReadWriter: pipelineReadWriter}

	// Retrieve the sources in the pipeline, most of the time there is just one
	var sources []*gst.Element
	sources, err = pipelineReadWriter.Pipeline().GetSourceElements()
	if err != nil {
		return nil, err
	}

	// Fetch the appsrc and make a local reference to it
	for _, src := range sources {
		if strings.Contains(src.Name(), "appsrc") {
			appPipeline.appSrc = &gst.AppSrc{Element: src}
		}
	}

	// Retrieve the sinks in the pipeline, most of the time there is just one
	var sinks []*gst.Element
	sinks, err = pipelineReadWriter.Pipeline().GetSinkElements()
	if err != nil {
		return nil, err
	}

	// Fetch the appsink and make a local reference to it
	for _, sink := range sinks {
		if strings.Contains(sink.Name(), "appsink") {
			appPipeline.appSink = &gst.AppSink{Element: sink}
		}
	}

	return appPipeline, nil
}

// GetAppSink returns the app sink for this pipeline.
func (p *PipelineReadWriterApp) GetAppSink() *gst.AppSink { return p.appSink }

// GetAppSource returns the app src for this pipeline.
func (p *PipelineReadWriterApp) GetAppSource() *gst.AppSrc { return p.appSrc }
