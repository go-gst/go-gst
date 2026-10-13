// This example shows how to use the appsink callbacks.
//
// Instead of connecting to the appsink signals, callbacks are installed with SetCallbacks. This
// avoids the overhead of the signal emission and allows the propose_allocation query to be
// modified in place.
//
// Also see: https://gstreamer.freedesktop.org/documentation/app/gstappsink.html#gst_app_sink_set_callbacks
package main

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstapp"
	"github.com/go-gst/go-gst/pkg/gstvideo"
)

func createPipeline() (gst.Pipeline, error) {
	gst.Init()

	ret, err := gst.ParseLaunch("videotestsrc num-buffers=50 ! video/x-raw,format=RGBA,width=320,height=240 ! appsink name=sink")
	if err != nil {
		return nil, err
	}

	pipeline := ret.(gst.Pipeline)

	sink := pipeline.GetByName("sink").(gstapp.AppSink)

	// Count the frames we receive. The callbacks are called from the streaming thread,
	// so we use an atomic in case it is read from somewhere else.
	var frames atomic.Uint64

	sink.SetCallbacks(&gstapp.AppSinkCallbacks{
		// ProposeAllocation is called when upstream asks us how buffers should be allocated.
		// Here we tell upstream that we support video meta on the buffers.
		//
		// The query is only valid during this callback and must not be retained.
		ProposeAllocation: func(sink gstapp.AppSink, query *gst.Query) bool {
			caps, _ := query.ParseAllocation()
			fmt.Println("propose allocation for caps:", caps.String())

			query.AddAllocationMeta(gstvideo.VideoMetaAPIGetType(), nil)

			return true
		},

		// NewPreroll is called when the pipeline prerolls, before going to PLAYING.
		NewPreroll: func(sink gstapp.AppSink) gst.FlowReturn {
			sample := sink.PullPreroll()
			if sample == nil {
				return gst.FlowEOS
			}

			fmt.Println("preroll with caps:", sample.GetCaps().String())

			return gst.FlowOK
		},

		// NewSample is called for every sample that arrives at the appsink.
		NewSample: func(sink gstapp.AppSink) gst.FlowReturn {
			sample := sink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}

			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}

			n := frames.Add(1)
			fmt.Printf("frame %d: pts=%v, size=%d\n", n, buffer.PTS(), buffer.GetSize())

			return gst.FlowOK
		},

		// EOS is called when the appsink received the end of stream.
		EOS: func(sink gstapp.AppSink) {
			fmt.Println("appsink reached EOS after", frames.Load(), "frames")
		},
	})

	return pipeline, nil
}

func mainLoop(pipeline gst.Pipeline) error {
	pipeline.SetState(gst.StatePlaying)
	defer pipeline.SetState(gst.StateNull)

	for msg := range pipeline.GetBus().Messages(context.Background()) {
		switch msg.Type() {
		case gst.MessageEOS:
			return nil
		case gst.MessageError:
			debug, gerr := msg.ParseError()
			if debug != "" {
				fmt.Println(gerr.Error(), debug)
			}
			return gerr
		}
	}

	return fmt.Errorf("unexpected end of messages without EOS")
}

func main() {
	pipeline, err := createPipeline()
	if err != nil {
		fmt.Println("Error creating pipeline:", err)
		return
	}

	if err := mainLoop(pipeline); err != nil {
		fmt.Println("Error running pipeline:", err)
	}
}
