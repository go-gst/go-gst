// This is a simplified go-reimplementation of the gst-launch-<version> cli tool.
// It has no own parameters and simply parses the cli arguments as launch syntax.
// When the parsing succeeded, the pipeline is run until the stream ends or an error happens.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/examples"
	"github.com/go-gst/go-gst/gst"
)

func runPipeline(mainLoop *glib.MainLoop) error {
	if len(os.Args) == 1 {
		return errors.New("pipeline string cannot be empty")
	}

	gst.Init(&os.Args)

	// Let GStreamer create a pipeline from the parsed launch syntax on the cli.
	pipeline, err := gst.NewPipelineFromString(strings.Join(os.Args[1:], " "))
	if err != nil {
		return err
	}

	// Add a message handler to the pipeline bus, printing interesting information to the console.
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // When end-of-stream is received stop the main loop
			pipeline.BlockSetState(gst.StateNull)
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
	return mainLoop.RunError()
}

func main() {
	examples.RunLoop(func(loop *glib.MainLoop) error {
		return runPipeline(loop)
	})
}
