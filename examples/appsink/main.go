package main

import (
	"fmt"
	"math"
	"time"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

func createPipeline() *gst.Pipeline {
	gst.Init(nil)

	pipeline, err := gst.NewPipeline("")
	if err != nil {
		panic(err)
	}

	src, err := gst.NewElement("audiotestsrc")
	if err != nil {
		panic(err)
	}

	sink, err := app.NewAppSink()
	if err != nil {
		panic(err)
	}

	pipeline.AddMany(src, sink.Element)
	src.Link(sink.Element)

	// Tell the appsink what format we want. It will then be the audiotestsrc's job to
	// provide the format we request.
	// This can be set after linking the two objects, because format negotiation between
	// both elements will happen during pre-rolling of the pipeline.
	sink.SetCaps(gst.NewCapsFromString(
		"audio/x-raw, format=S16LE, layout=interleaved, channels=1",
	))

	// Getting data out of the appsink is done by setting callbacks on it.
	// The appsink will then call those handlers, as soon as data is available.
	sink.SetCallbacks(&app.SinkCallbacks{
		// Add a "new-sample" callback
		NewSampleFunc: func(sink *app.Sink) gst.FlowReturn {

			// Pull the sample that triggered this callback
			sample := sink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}

			// Retrieve the buffer from the sample
			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}

			// At this point, buffer is only a reference to an existing memory region somewhere.
			// When we want to access its content, we have to map it while requesting the required
			// mode of access (read, read/write).
			//
			// We also know what format to expect because we set it with the caps. So we convert
			// the map directly to signed 16-bit integers.
			samples := buffer.Map(gst.MapRead).AsInt16Slice()

			// Calculate the root mean square for the buffer
			// (https://en.wikipedia.org/wiki/Root_mean_square)
			var square float64
			for _, i := range samples {
				square += float64(i * i)
			}
			rms := math.Sqrt(square / float64(len(samples)))
			fmt.Println("rms:", rms)

			return gst.FlowOK
		},
	})

	return pipeline
}

func mainLoop(pipeline *gst.Pipeline) error {
	defer pipeline.Destroy()

	pipeline.SetState(gst.StatePlaying)

	bus := pipeline.GetPipelineBus()

	for {
		msg := bus.TimedPop(time.Duration(-1))
		if msg == nil {
			break
		}
		switch msg.Type() {
		case gst.MessageEOS:
			break
		case gst.MessageError:
			return msg.ParseError()
		}
	}

	return nil
}

func main() {
	examples.Run(func() error {
		pipeline := createPipeline()
		return mainLoop(pipeline)
	})

}
