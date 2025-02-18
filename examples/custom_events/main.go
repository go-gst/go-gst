// This example demonstrates the use of custom events in a pipeline.
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

// ExampleCustomEvent demonstrates a custom event structue. Currerntly nested structs
// are not supported.
type ExampleCustomEvent struct {
	Count   int
	SendEOS bool
}

func createPipeline() (*gst.Pipeline, error) {
	gst.Init()

	// Create a new pipeline from a launch string
	ret, err := gst.ParseLaunch(
		"audiotestsrc name=src ! queue max-size-time=2000000000 ! fakesink name=sink sync=true",
	)

	if err != nil {
		return nil, err
	}

	pipeline := ret.(*gst.Pipeline)

	var sink *gst.Element
	var sinkpad *gst.Pad

	// Retrieve the sink pad
	for v := range pipeline.IterateSinks().Values() {
		sink = v.(*gst.Element)
		sinkpad = sink.StaticPad("sink")
		break
	}

	if sink == nil || sinkpad == nil {
		return nil, fmt.Errorf("could not find sink")
	}

	// Add a probe for out custom event
	sinkpad.AddProbe(gst.PadProbeTypeEventDownstream, func(self *gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		// Retrieve the event from the probe
		ev := info.Event()

		// Extra check to make sure it is the right type.
		if ev.Type() != gst.EventCustomDownstream {
			return gst.PadProbeHandled
		}

		// Unmarshal the event into our custom one
		var customEvent ExampleCustomEvent
		if err := ev.Structure().UnmarshalInto(&customEvent); err != nil {
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
			sink.CallAsync(func(el gst.Elementer) {
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

func runPipeline(loop *glib.MainLoop, pipeline *gst.Pipeline) {
	// Create a watch on the pipeline to kill the main loop when EOS is received
	pipeline.Bus().AddWatch(0, func(bus *gst.Bus, msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEos:
			fmt.Println("Got EOS message")
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

			if !pipeline.SendEvent(gst.NewEventCustom(gst.EventCustomDownstream, st)) {
				fmt.Println("Warning: failed to send custom event")
			}
			if count == 3 {
				break
			}
			count++
		}
	}()

	// When passing an object created by the bindings between scopes, there is a posibility
	// the finalizer will leak and destroy your object before you are done with it.  One way
	// of dealing with this is by taking an additional Ref and disposing of it when you are
	// done with the new scope. An alternative is to declare Keep() *after* where you know
	// you will be done with the object. This instructs the runtime to defer the finalizer
	// until after this point is passed in the code execution.

	loop.Run()

	runtime.KeepAlive(pipeline)
}

func main() {
	pipeline, err := createPipeline()

	if err != nil {
		panic(err)
	}

	mainloop := glib.NewMainLoop(glib.MainContextDefault(), false)

	runPipeline(mainloop, pipeline)
}
