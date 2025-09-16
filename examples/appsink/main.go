// This example shows how to use the appsink element.
//
// Also see: https://gstreamer.freedesktop.org/documentation/tutorials/basic/short-cutting-the-pipeline.html?gi-language=c
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/signal"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstapp"
)

func createPipeline() (gst.Pipeline, error) {
	gst.Init()

	pipeline := gst.NewPipeline("").(gst.Pipeline)

	src := gst.ElementFactoryMake("audiotestsrc", "")
	sink := gst.ElementFactoryMake("appsink", "").(gstapp.AppSink)

	pipeline.AddMany(src, sink)
	src.Link(sink)

	// Tell the appsink what format we want. It will then be the audiotestsrc's job to
	// provide the format we request.
	// This can be set after linking the two objects, because format negotiation between
	// both elements will happen during pre-rolling of the pipeline.
	sink.SetCaps(gst.CapsFromString(
		"audio/x-raw, format=S16LE, layout=interleaved, channels=1",
	))

	sink.ConnectNewSample(func(sink gstapp.AppSink) gst.FlowReturn {
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
		mapInfo, ok := buffer.Map(gst.MapRead)
		if !ok {
			return gst.FlowError
		}
		defer mapInfo.Unmap()

		// Calculate the root mean square for the buffer
		// (https://en.wikipedia.org/wiki/Root_mean_square)

		samples := mapInfo.Int16Data(binary.LittleEndian)
		var square float64
		for _, i := range samples {
			square += float64(i * i)
		}
		rms := math.Sqrt(square / float64(len(samples)))
		fmt.Println("rms:", rms)

		return gst.FlowOK
	})

	return pipeline, nil
}

func handleMessage(msg *gst.Message) error {
	switch msg.Type() {
	case gst.MessageEOS:
		return fmt.Errorf("end of stream")
	case gst.MessageError:
		debug, gerr := msg.ParseError()
		if debug != "" {
			fmt.Println(debug)
		}
		return gerr
	}
	return nil
}

func runPipeline(pipeline gst.Pipeline) error {

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Retrieve the bus from the pipeline
	bus := pipeline.GetBus()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pipeline.SendEvent(gst.NewEventEOS())
	}()

	// Loop over messsages from the pipeline
	for {
		msg := bus.TimedPop(gst.ClockTimeNone)
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
	pipeline, err := createPipeline()
	if err != nil {
		println(err)
		return
	}

	err = runPipeline(pipeline)

	println(err)
}
