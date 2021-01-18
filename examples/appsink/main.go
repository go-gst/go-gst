// This example shows how to use the appsink element.
package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

func createPipeline() (*gst.Pipeline, error) {
	gst.Init(nil)

	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	src, err := gst.NewElement("audiotestsrc")
	if err != nil {
		return nil, err
	}

	// Should this actually be a *gst.Element that produces an Appsink interface like the
	// rust examples?
	sink, err := app.NewAppSink()
	if err != nil {
		return nil, err
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
			// the map directly to signed 16-bit little-endian integers.
			samples := buffer.Map(gst.MapRead).AsInt16LESlice()
			defer buffer.Unmap()

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

	return pipeline, nil
}

func handleMessage(msg *gst.Message) error {

	switch msg.Type() {
	case gst.MessageEOS:
		return app.ErrEOS
	case gst.MessageError:
		return msg.ParseError()
	}

	return nil
}

func mainLoop(pipeline *gst.Pipeline) error {

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Retrieve the bus from the pipeline
	bus := pipeline.GetPipelineBus()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pipeline.SendEvent(gst.NewEOSEvent())
	}()

	// Loop over messsages from the pipeline
	for {
		msg := bus.TimedPop(time.Duration(-1))
		if msg == nil {
			break
		}
		if err := handleMessage(msg); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	examples.Run(func() error {
		var pipeline *gst.Pipeline
		var err error
		if pipeline, err = createPipeline(); err != nil {
			return err
		}
		return mainLoop(pipeline)
	})
}
