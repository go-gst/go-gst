package main

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle messages
	for msg := range pipeline.GetBus().Messages(ctx) {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received stop
			fmt.Println("End-of-stream reached")
			bin.BlockSetState(gst.StateNull, gst.ClockTime(time.Second))
			cancel()
		case gst.MessageError: // Error messages are always fatal
			debug, err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			cancel()
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			fmt.Println(msg)
		}
	}
}
