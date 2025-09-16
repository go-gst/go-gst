// This example demonstrates the use of GStreamer's ToC API.
//
// This API is used to manage a table of contents contained in the handled media stream.
// Chapters within a matroska file would be an example of a scenario for using
// this API. Elements that can parse ToCs from a stream (such as matroskademux)
// notify all elements in the pipeline when they encountered a ToC.
// For this, the example operates the following pipeline:
//
//	                         /-{queue} - {fakesink}
//	{filesrc} - {decodebin} - {queue} - {fakesink}
//	                         \- ...
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

func tagsetter() error {
	gst.Init()

	if len(os.Args) < 2 {
		return errors.New("usage: toc <file>")
	}

	pipeline := gst.NewPipeline("")

	src := gst.ElementFactoryMake("filesrc", "")

	decodebin := gst.ElementFactoryMake("decodebin", "")

	src.SetObjectProperty("location", os.Args[1])

	pipeline.AddMany(src, decodebin)
	gst.LinkMany(src, decodebin)

	// Connect to decodebin's pad-added signal, that is emitted whenever it found another stream
	// from the input file and found a way to decode it to its raw format.
	decodebin.ConnectPadAdded(func(srcPad *gst.Pad) {

		// In this example, we are only interested about parsing the ToC, so
		// we simply pipe every encountered stream into a fakesink, essentially
		// throwing away the data.
		queue := gst.ElementFactoryMake("queue", "")
		fakesink := gst.ElementFactoryMake("fakesink", "")

		pipeline.AddMany(queue, fakesink)
		gst.LinkMany(queue, fakesink)
		queue.SyncStateWithParent()
		fakesink.SyncStateWithParent()

		sinkPad := queue.StaticPad("sink")
		if sinkPad == nil {
			fmt.Println("Could not get static pad from sink")
			return
		}

		srcPad.
			Link(sinkPad)
	})

	if ret := pipeline.BlockSetState(gst.StatePaused, gst.ClockTime(time.Second)); ret != gst.StateChangeSuccess {
		return fmt.Errorf("could not change state")
	}

	// Instead of using the main loop, we manually iterate over GStreamer's bus messages
	// in this example. We don't need any special functionality like timeouts or GLib socket
	// notifications, so this is sufficient. The bus is manually operated by repeatedly calling
	// timed_pop on the bus with the desired timeout for when to stop waiting for new messages.
	// (-1 = Wait forever)
	for {
		msg := pipeline.Bus().TimedPop(gst.ClockTimeNone)
		switch msg.Type() {

		// When we use this method of popping from the bus (instead of a Watch), we own a
		// reference to every message received (this may be abstracted later).
		default:
			// fmt.Println(msg)

		// End of stream
		case gst.MessageEos:
			// Errors from any elements
		case gst.MessageError:
			gerr, debug := msg.ParseError()
			if debug != "" {
				fmt.Println("go-gst-debug:", debug)
			}
			return gerr

		// Some element found a ToC in the current media stream and told
		// us by posting a message to GStreamer's bus.
		case gst.MessageToc:
			// Parse the toc from the message
			toc, updated := msg.ParseToc()
			fmt.Printf("Received toc: %s - updated %v\n", toc.Scope().String(), updated)
			// Get a list of tags that are ToC specific.
			if tags := toc.Tags(); tags != nil {
				fmt.Println("- tags:", tags)
			}
			// ToCs do not have a fixed structure. Depending on the format that
			// they were parsed from, they might have different tree-like structures,
			// so applications that want to support ToCs (for example in the form
			// of jumping between chapters in a video) have to try parsing  and
			// interpreting the ToC manually.
			// In this example, we simply want to print the ToC structure, so
			// we iterate everything and don't try to interpret anything.
			for _, entry := range toc.Entries() {
				// Every entry in a ToC has its own type. One type could for
				// example be Chapter.
				fmt.Printf("\t%s - %s\n", entry.EntryType().String(), entry.Uid())

				// Every ToC entry can have a set of timestamps (start, stop).
				if start, stop, ok := entry.StartStopTimes(); ok {
					startDur := time.Duration(start) * time.Nanosecond
					stopDur := time.Duration(stop) * time.Nanosecond
					fmt.Printf("\t- start: %s, stop: %s\n", startDur, stopDur)
				}

				// Every ToC entry can have tags to it.
				if tags := entry.Tags(); tags != nil {
					fmt.Println("\t- tags:", tags)
				}

				// Every ToC entry can have a set of child entries.
				// With this structure, you can create trees of arbitrary depth.
				for _, subEntry := range entry.SubEntries() {
					fmt.Printf("\n\t\t%s - %s\n", subEntry.EntryType().String(), subEntry.Uid())
					if start, stop, ok := entry.StartStopTimes(); ok {
						startDur := time.Duration(start) * time.Nanosecond
						stopDur := time.Duration(stop) * time.Nanosecond
						fmt.Printf("\t\t- start: %s, stop: %s\n", startDur, stopDur)
					}
					if tags := entry.Tags(); tags != nil {
						fmt.Println("\t\t- tags:", tags)
					}
				}
			}
		}
	}
}

func main() {
	tagsetter()
}
