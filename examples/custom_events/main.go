// This example demonstrates the use of custom events in a pipeline.
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
)

// ExampleCustomEvent demonstrates a custom event structue. Currerntly nested structs
// are not supported.
type ExampleCustomEvent struct {
	Count   int
	SendEOS bool
}

func createPipeline() (*gst.Pipeline, error) {
	gst.Init(nil)

	// Create a new pipeline from a launch string
	pipeline, err := gst.NewPipelineFromString(
		"audiotestsrc name=src ! queue max-size-time=2000000000 ! fakesink name=sink sync=true",
	)

	// Retrieve the sink element
	sinks, err := pipeline.GetSinkElements()
	if err != nil {
		return nil, err
	} else if len(sinks) != 1 {
		return nil, errors.New("Expected one sink back")
	}
	sink := sinks[0]

	// Get the sink pad
	sinkpad := sink.GetStaticPad("sink")

	// Add a probe for out custom event
	sinkpad.AddProbe(gst.PadProbeTypeEventDownstream, func(self *gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		// Retrieve the event from the probe
		ev := info.GetEvent()

		// Extra check to make sure it is the right type.
		if ev.Type() != gst.EventTypeCustomDownstream {
			return gst.PadProbeUnhandled
		}

		// Unmarshal the event into our custom one
		var customEvent ExampleCustomEvent
		if err := ev.GetStructure().UnmarshalInto(&customEvent); err != nil {
			fmt.Println("Could not parse the custom event!")
			return gst.PadProbeUnhandled
		}

		// Log and act accordingly
		fmt.Printf("Received custom event with count=%d send_eos=%v\n", customEvent.Count, customEvent.SendEOS)
		if customEvent.SendEOS {
			fmt.Println("Send EOS is true, sending eos")
			if !pipeline.GetPipelineBus().Post(gst.NewEOSMessage(self)) {
				fmt.Println("WARNING: Failed to send EOS to pipeline")
			}
		} else {
			fmt.Println("Send EOS is false ignoring")
		}
		return gst.PadProbeOK
	})

	return pipeline, nil
}

func mainLoop(loop *gst.MainLoop, pipeline *gst.Pipeline) error {
	// Create a watch on the pipeline to kill the main loop when EOS is received
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS:
			fmt.Println("Got EOS message")
			pipeline.Destroy()
			loop.Quit()
		default:
			fmt.Println(msg)
		}
		return true
	})

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	go func() {
		// Loop and on the third iteration send the custom event.
		ticker := time.NewTicker(time.Second * 2)
		count := 0
		for range ticker.C {
			ev := ExampleCustomEvent{Count: count}
			if count == 3 {
				ev.SendEOS = true
			}
			st := gst.MarshalStructure(ev)
			if !pipeline.SendEvent(gst.NewCustomEvent(gst.EventTypeCustomDownstream, st)) {
				fmt.Println("Warning: failed to send custom event")
			}
			if count == 3 {
				break
			}
			count++
		}
	}()

	return loop.RunError()
}

func main() {
	examples.RunLoop(func(loop *gst.MainLoop) error {
		var pipeline *gst.Pipeline
		var err error
		if pipeline, err = createPipeline(); err != nil {
			return err
		}
		return mainLoop(loop, pipeline)
	})
}
