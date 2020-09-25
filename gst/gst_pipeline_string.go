package gst

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// NewPipelineFromLaunchString returns a new GstPipeline from the given launch string. If flags
// contain PipelineRead or PipelineWrite, the launch string is further formatted accordingly.
//
// If using PipelineWrite, you should generally start your pipeline with the caps of the source.
func NewPipelineFromLaunchString(launchStr string, flags PipelineFlags) (*Pipeline, error) {
	// reformat the string to point at the writerFd
	if flags.has(PipelineWrite) {

		if flags.has(PipelineUseGstApp) {
			if launchStr == "" {
				launchStr = "appsrc"
			} else {
				launchStr = fmt.Sprintf("appsrc ! %s", launchStr)
			}
		} else {
			if launchStr == "" {
				launchStr = "fdsrc"
			} else {
				launchStr = fmt.Sprintf("fdsrc ! %s", launchStr)
			}
		}

	}

	if flags.has(PipelineRead) {

		if flags.has(PipelineUseGstApp) {
			if launchStr == "" {
				launchStr = "appsink emit-signals=false"
			} else {
				launchStr = fmt.Sprintf("%s ! appsink emit-signals=false", launchStr)
			}
		} else {
			if launchStr == "" {
				launchStr = "fdsink"
			} else {
				launchStr = fmt.Sprintf("%s ! fdsink", launchStr)
			}
		}

	}

	pipelineElement, err := newPipelineFromString(launchStr)
	if err != nil {
		return nil, err
	}

	pipeline := wrapPipeline(glib.Take(unsafe.Pointer(pipelineElement)))

	if err := applyFlags(pipeline, flags); err != nil {
		return nil, err
	}

	if flags.has(PipelineWrite) {

		sources, err := pipeline.GetSourceElements()
		if err != nil {
			return nil, err
		}

		var srcType string
		if flags.has(PipelineUseGstApp) {
			srcType = "appsrc"
		} else {
			srcType = "fdsrc"
		}

		var pipelineSrc *Element
		for _, src := range sources {
			if strings.Contains(src.Name(), srcType) {
				pipelineSrc = src
			} else {
				src.Unref()
			}
		}

		if pipelineSrc == nil {
			return nil, errors.New("Could not detect pipeline source")
		}

		defer pipelineSrc.Unref()

		if flags.has(PipelineUseGstApp) {
			pipeline.appSrc = wrapAppSrc(pipelineSrc)
		} else {
			if err := pipelineSrc.Set("fd", pipeline.writerFd()); err != nil {
				return nil, err
			}
		}
	}

	if flags.has(PipelineRead) {
		sinks, err := pipeline.GetSinkElements()
		if err != nil {
			return nil, err
		}

		var sinkType string
		if flags.has(PipelineUseGstApp) {
			sinkType = "appsink"
		} else {
			sinkType = "fdsink"
		}

		var pipelineSink *Element
		for _, sink := range sinks {
			if strings.Contains(sink.Name(), sinkType) {
				pipelineSink = sink
			} else {
				sink.Unref()
			}
		}

		if pipelineSink == nil {
			return nil, errors.New("Could not detect pipeline sink")
		}

		defer pipelineSink.Unref()

		if flags.has(PipelineUseGstApp) {
			pipeline.appSink = wrapAppSink(pipelineSink)
		} else {
			if err := pipelineSink.Set("fd", pipeline.readerFd()); err != nil {
				return nil, err
			}
		}

	}

	// signal that this pipeline was made from a string and therefore already linked
	pipeline.pipelineFromHelper = true

	return pipeline, err
}
