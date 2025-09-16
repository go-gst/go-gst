// This example demonstrates how to use GStreamer's query functionality.
//
// These are a way to query information from either elements or pads.
// Such information could for example be the current position within
// the stream (i.e. the playing time). Queries can traverse the pipeline
// (both up and downstream). This functionality is essential, since most
// queries can only answered by specific elements in a pipeline (such as the
// stream's duration, which often can only be answered by the demuxer).
// Since gstreamer has many elements that itself contain other elements that
// we don't know of, we can simply send a query for the duration into the
// pipeline and the query is passed along until an element feels capable
// of answering.
// For convenience, the API has a set of pre-defined queries, but also
// allows custom queries (which can be defined and used by your own elements).
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

func queries() error {

	if len(os.Args) < 2 {
		fmt.Println("USAGE: queries <pipeline>")
		os.Exit(1)
	}

	gst.Init()

	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	// Let GStreamer create a pipeline from the parsed launch syntax on the cli.
	pipelineStr := strings.Join(os.Args[1:], " ")
	ret, err := gst.ParseLaunch(pipelineStr)
	if err != nil {
		return err
	}

	pipeline := ret.(gst.Binner)

	// Get a reference to the pipeline bus
	bus := pipeline.Bus()

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Kick of a goroutine that will send a query to the pipeline
	// every second.
	go func() {
		for range time.NewTicker(time.Second).C {
			// Create a new position query and send it to the pipeline.
			// This will traverse all elements in the pipeline, until one feels
			// capable of answering the query.
			pos := gst.NewQueryPosition(gst.FormatTime)
			if ok := pipeline.Query(pos); !ok {
				fmt.Println("Failed to query position from pipeline")
			}
			// Create a new duration query and send it to the pipeline.
			// This will traverse all elements in the pipeline, until one feels
			// capable of answering the query.
			dur := gst.NewQueryDuration(gst.FormatTime)
			if ok := pipeline.Query(dur); !ok {
				fmt.Println("Failed to query duration from pipeline")
			}

			// The values from the queries above are both in nanoseconds (this may be abstracted later).
			// We can convert them to durations.
			_, posVal := pos.ParsePosition() //  If either of the above queries failed, these values
			_, durVal := dur.ParseDuration() //  will be 0.
			posDur := time.Duration(posVal) * time.Nanosecond
			durDur := time.Duration(durVal) * time.Nanosecond

			fmt.Println(posDur, "/", durDur)
		}
	}()

	bus.AddWatch(0, func(bus *gst.Bus, msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEos:
			mainLoop.Quit()
		case gst.MessageError:
			gstErr, debug := msg.ParseError()
			fmt.Printf("Error from %s: %s\n", msg.Src(), gstErr.Error())
			if debug != "" {
				fmt.Println("go-gst-debug:", debug)
			}
			mainLoop.Quit()
		}
		return true
	})

	mainLoop.Run()

	bus.RemoveWatch()

	return nil
}

func main() {
	if err := queries(); err != nil {
		fmt.Println(err)
	}
}
