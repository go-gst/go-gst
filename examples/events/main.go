// This example demonstrates how events can be created and sent to the pipeline.
// What this example does is scheduling a timeout in a goroutine, and
// sending an EOS message on the bus from there - telling the pipeline
// to shut down. Once that event is processed by everything, the EOS message
// is going to be sent and we catch that one to shut down everything.
//
// GStreamer's bus is an abstraction layer above an arbitrary main loop.
// This makes sure that GStreamer can be used in conjunction with any existing
// other framework (GUI frameworks, mostly) that operate their own main loops.
// Main idea behind the bus is the simplification between the application and
// GStreamer, because GStreamer is heavily threaded underneath.
//
// Any thread can post messages to the bus, which is essentially a thread-safe
// queue of messages to process. When a new message was sent to the bus, it
// will wake up the main loop implementation underneath it (which will then
// process the pending messages from the main loop thread).
//
// An application itself can post messages to the bus aswell.
// This makes it possible, e.g., to schedule an arbitrary piece of code
// to run in the main loop thread - avoiding potential threading issues.
package main

import (
	"fmt"
	"time"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
)

func runPipeline(loop *glib.MainLoop) error {
	gst.Init(nil)

	// Build a pipeline with fake audio data going to a fakesink
	pipeline, err := gst.NewPipelineFromString("audiotestsrc ! fakesink")
	if err != nil {
		return err
	}

	// Retrieve the message bus for the pipeline
	bus := pipeline.GetPipelineBus()

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// This sets the bus's signal handler (don't be mislead by the "add", there can only be one).
	// Every message from the bus is passed through this function. Its return value determines
	// whether the handler wants to be called again.
	bus.AddWatch(func(msg *gst.Message) (cont bool) {
		// Assume we are continuing
		cont = true

		switch msg.Type() {
		case gst.MessageEOS:
			fmt.Println("Received EOS")
			// An EndOfStream event was sent to the pipeline, so we tell our main loop
			// to stop execution here.
			loop.Quit()
		case gst.MessageError:
			err := msg.ParseError()
			fmt.Println("ERROR:", err)
			fmt.Println("DEBUG:", err.DebugString())
			loop.Quit()
		}

		return
	})

	// Kick off a goroutine that after 5 seconds will send an eos event to the pipeline.
	go func() {
		for range time.NewTicker(time.Second * 5).C {
			fmt.Println("Sending EOS")
			// We create an EndOfStream event here, that tells all elements to drain
			// their internal buffers to their following elements, essentially draining the
			// whole pipeline (front to back). It ensuring that no data is left unhandled and potentially
			// headers were rewritten (e.g. when using something like an MP4 or Matroska muxer).
			// The EOS event is handled directly from this very goroutine until the first
			// queue element is reached during pipeline-traversal, where it is then queued
			// up and later handled from the queue's streaming thread for the elements
			// following that queue.
			// Once all sinks are done handling the EOS event (and all buffers that were before the
			// EOS event in the pipeline already), the pipeline would post an EOS message on the bus,
			// essentially telling the application that the pipeline is completely drained.
			pipeline.SendEvent(gst.NewEOSEvent())
			return
		}
	}()

	// Operate GStreamer's bus, facilliating GLib's mainloop here.
	// This function call will block until you tell the mainloop to quit
	// (see above for how to do this).
	loop.Run()

	// Stop the pipeline
	if err := pipeline.SetState(gst.StateNull); err != nil {
		fmt.Println("Error stopping pipeline:", err)
	}

	// Remove the watch function from the bus.
	// Again: There can always only be one watch function.
	// Thus we don't have to tell it which function to remove.
	bus.RemoveWatch()

	return nil
}

func main() {
	examples.RunLoop(func(loop *glib.MainLoop) error {
		return runPipeline(loop)
	})
}
