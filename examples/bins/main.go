package main

import (
	"fmt"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

func main() {
	gst.Init()
	bin, err := gst.ParseBinFromDescription("fakesrc num-buffers=5 ! fakesink", true)
	if err != nil {
		panic(err)
	}

	pipeline := gst.NewPipeline("pipeline").(gst.Pipeline)

	pipeline.Add(bin)

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// handle messages
	messages:
	for msg := range pipeline.GetBus().Messages() {
		switch msg.Type() {
		case gst.MessageEos: // When end-of-stream is received stop
			fmt.Println("End-of-stream reached")
			bin.BlockSetState(gst.StateNull, gst.ClockTime(time.Second))
			break messages
		case gst.MessageError: // Error messages are always fatal
			debug, err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			break messages
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			fmt.Println(msg)
		}
	}
}
