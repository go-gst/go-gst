// This example demonstrates the use of custom events in a pipeline.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

// ExampleCustomEvent demonstrates a custom event structue. Currerntly nested structs
// are not supported.
type ExampleCustomEvent struct {
	Count   int
	SendEOS bool
}

func createPipeline() (gst.Pipeline, error) {
	gst.Init()

	// Create a new pipeline from a launch string
	ret, err := gst.ParseLaunch(
		"audiotestsrc name=src ! queue max-size-time=2000000000 ! fakesink name=sink sync=true",
	)

	if err != nil {
		return nil, err
	}

	pipeline := ret.(gst.Pipeline)

	var sink gst.Element
	var sinkpad gst.Pad

	// Retrieve the sink pad
	for v := range pipeline.IterateSinks().Values() {
		sink = v.(gst.Element)

		sinkpad = sink.GetStaticPad("sink")
		break
	}

	if sink == nil || sinkpad == nil {
		return nil, fmt.Errorf("could not find sink")
	}

	// Add a probe for out custom event
	sinkpad.AddProbe(gst.PadProbeTypeEventDownstream, func(self gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		// Retrieve the event from the probe
		ev := info.GetEvent()

		// Extra check to make sure it is the right type.
		if ev.GetType() != gst.EventCustomDownstream {
			return gst.PadProbeHandled
		}

		// Unmarshal the event into our custom one
		var customEvent ExampleCustomEvent
		if err := ev.GetStructure().UnmarshalInto(&customEvent); err != nil {
			fmt.Println("Could not parse the custom event!")
			return gst.PadProbeHandled
		}

		// Log and act accordingly
		fmt.Printf("Received custom event with count=%d send_eos=%v\n", customEvent.Count, customEvent.SendEOS)
		if customEvent.SendEOS {
			// We need to use the CallAsync method to send the signal.
			// This is becaues the SendEvent method blocks and this could cause a dead lock sending the
			// event directly from the probe. This is the near equivalent of using go func() { ... }(),
			// however displayed this way for demonstration purposes.
			sink.CallAsync(func(el gst.Element) {
				fmt.Println("Send EOS is true, sending eos")
				if !pipeline.SendEvent(gst.NewEventEos()) {
					fmt.Println("WARNING: Failed to send EOS to pipeline")
				}
				fmt.Println("Sent EOS")
			})
			return gst.PadProbeRemove
		}
		fmt.Println("Send EOS is false ignoring")
		return gst.PadProbeOK
	})

	return pipeline, nil
}

func runPipeline(pipeline gst.Pipeline) {

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)
	defer pipeline.SetState(gst.StateNull)

	go func() {
		// Loop and on the third iteration send the custom event.
		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()

		count := 0
		for range ticker.C {
			ev := ExampleCustomEvent{Count: count}
			if count == 3 {
				ev.SendEOS = true
			}
			st := gst.MarshalStructure(ev)

			if !pipeline.SendEvent(gst.NewEventCustom(gst.EventCustomDownstream, st)) {
				fmt.Println("Warning: failed to send custom event")
			}
			if count == 3 {
				break
			}
			count++
		}
	}()

	for msg := range pipeline.GetBus().Messages(context.Background()) {
		switch msg.Type() {
		case gst.MessageEos:
			fmt.Println("Got EOS message")
			return
		default:
			fmt.Println(msg)
		}
	}
}

func main() {
	pipeline, err := createPipeline()

	if err != nil {
		panic(err)
	}

	runPipeline(pipeline)
}
