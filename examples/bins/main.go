package main

import (
	"fmt"
	"os"

	"github.com/go-gst/go-gst/examples"
	"github.com/go-gst/go-gst/gst"
	"github.com/gotk3/gotk3/glib"
)

func runPipeline(mainLoop *glib.MainLoop) error {
	gst.Init(&os.Args)

	bin, err := gst.NewBinFromString("fakesrc num-buffers=5 ! fakesink", true)
	if err != nil {
		return err
	}

	pipeline, err := gst.NewPipeline("pipeline")
	if err != nil {
		return err
	}

	pipeline.Add(bin.Element)
	pipeline.GetBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received stop the main loop
			bin.BlockSetState(gst.StateNull)
			mainLoop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
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

	return nil
}

func main() {
	examples.RunLoop(func(loop *glib.MainLoop) error {
		return runPipeline(loop)
	})
}
