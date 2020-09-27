package app

import (
	"errors"
	"strings"

	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
	"github.com/tinyzimmer/go-gst/gst/gstauto"
)

// PipelineReadWriterApp implements a ReadPipeliner that configures gstreamer
// with an appsink. The appsink allows for more granular control over the data
// at the end of the pipeline, and the appsrc allows for control over the data
// at the start.
type PipelineReadWriterApp struct {
	*gstauto.PipelineReadWriter

	appSrc  *app.Source
	appSink *app.Sink
}

// NewPipelineReadWriterAppFromString returns a new PipelineReadWriterApp populated from
// the given launch string. An appsink is added to the end of the launch string and made
// available via the GetAppSink method, and an appsrc is added to the end and made
// available via the GetAppSource method.
func NewPipelineReadWriterAppFromString(launchStr string) (*PipelineReadWriterApp, error) {
	pipelineReadWriter, err := gstauto.NewPipelineReadWriterFromString(addAppSourceToStr(addAppSinkToStr(launchStr)))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			runOrPrintErr(pipelineReadWriter.Pipeline().Destroy)
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
			appPipeline.appSrc = &app.Source{Element: src}
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
			appPipeline.appSink = &app.Sink{Element: sink}
		}
	}

	// Return the pipeline
	return appPipeline, nil
}

// NewPipelineReadWriterAppFromConfig returns a new PipelineReadWriterApp populated from
// the given launch config. An appsink is added to the end of the launch config and
// made available via the GetAppSink method, and an appsrc is added at the front and made
// available via the GetAppSource method.
func NewPipelineReadWriterAppFromConfig(cfg *gstauto.PipelineConfig) (*PipelineReadWriterApp, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Elements cannot be nil in the config")
	}
	pipelineReadWriter, err := gstauto.NewPipelineReadWriter("")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			runOrPrintErr(pipelineReadWriter.Pipeline().Destroy)
		}
	}()

	cfg.Elements = append(cfg.Elements, &gstauto.PipelineElement{Name: "appsink"})

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
			appPipeline.appSrc = &app.Source{Element: src}
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
			appPipeline.appSink = &app.Sink{Element: sink}
		}
	}

	return appPipeline, nil
}

// GetAppSink returns the app sink for this pipeline.
func (p *PipelineReadWriterApp) GetAppSink() *app.Sink { return p.appSink }

// GetAppSource returns the app src for this pipeline.
func (p *PipelineReadWriterApp) GetAppSource() *app.Source { return p.appSrc }
