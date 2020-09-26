package gstauto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tinyzimmer/go-gst-launch/gst"
)

// PipelineReaderApp implements a ReadPipeliner that configures gstreamer
// with an appsink. The appsink allows for more granular control over the data
// at the end of the pipeline.
type PipelineReaderApp struct {
	*PipelineReader

	appSink *gst.AppSink
}

// NewPipelineReaderAppFromString returns a new PipelineReaderApp populated from
// the given launch string. An appsink is added to the end of the launch string and made
// available via the GetAppSink method.
func NewPipelineReaderAppFromString(launchStr string) (*PipelineReaderApp, error) {
	fmt.Println(addAppSinkToStr(launchStr))
	pipelineReader, err := NewPipelineReaderFromString(addAppSinkToStr(launchStr))
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

	appPipeline := &PipelineReaderApp{PipelineReader: pipelineReader}

	// Retrieve the sinks in the pipeline, most of the time there is just one
	var sinks []*gst.Element
	sinks, err = pipelineReader.Pipeline().GetSinkElements()
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

// NewPipelineReaderAppFromConfig returns a new PipelineReaderApp populated from
// the given launch config. An appsink is added to the end of the launch config and
// made available via the GetAppSink method.
func NewPipelineReaderAppFromConfig(cfg *PipelineConfig) (*PipelineReaderApp, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Elements cannot be nil in the config")
	}
	pipelineReader, err := NewPipelineReader("")
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

	cfg.Elements = append(cfg.Elements, &PipelineElement{Name: "appsink"})

	if err = cfg.Apply(pipelineReader.Pipeline()); err != nil {
		return nil, err
	}

	appPipeline := &PipelineReaderApp{PipelineReader: pipelineReader}

	// Retrieve the sinks in the pipeline, most of the time there is just one
	var sinks []*gst.Element
	sinks, err = pipelineReader.Pipeline().GetSinkElements()
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

func addAppSinkToStr(pstr string) string {
	if pstr == "" {
		return "appsink"
	}
	return fmt.Sprintf("%s ! appsink", pstr)
}

// GetAppSink returns the app sink for this pipeline.
func (p *PipelineReaderApp) GetAppSink() *gst.AppSink { return p.appSink }
