package main

import (
	"fmt"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

func main() {
	gst.Init()

	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	bin, err := gst.ParseBinFromDescription("fakesrc num-buffers=5 ! fakesink", true)
	if err != nil {
		panic(err)
	}

	pipeline := gst.NewPipeline("pipeline")

	pipeline.Add(bin)

	pipeline.Bus().AddWatch(0, func(bus *gst.Bus, msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEos: // When end-of-stream is received stop the main loop
			bin.BlockSetState(gst.StateNull, gst.ClockTime(time.Second))
			mainLoop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err, debug := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			mainLoop.Quit()
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			fmt.Println(msg)
		}
		return true
	})

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Block on the main loop
	mainLoop.Run()
}
