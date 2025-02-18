// This is a simplified go-reimplementation of the gst-launch-<version> cli tool.
// It has no own parameters and simply parses the cli arguments as launch syntax.
// When the parsing succeeded, the pipeline is run until the stream ends or an error happens.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

func main() {
	gst.Init()

	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	// Let GStreamer create a pipeline from the parsed launch syntax on the cli.
	res, err := gst.ParseLaunch(strings.Join(os.Args[1:], " "))
	if err != nil {
		fmt.Printf("Parse error: %v", err)
		return
	}

	pipeline := res.(*gst.Pipeline)

	// Add a message handler to the pipeline bus, printing interesting information to the console.
	pipeline.Bus().AddWatch(0, func(_ *gst.Bus, msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEos: // When end-of-stream is received stop the main loop
			pipeline.BlockSetState(gst.StateNull, gst.ClockTime(time.Second))

			mainLoop.Quit()
		case gst.MessageError: // Error messages are always fatal
			err, debug := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			mainLoop.Quit()

		case gst.MessageAny:
		case gst.MessageApplication:
		case gst.MessageAsyncDone:
		case gst.MessageAsyncStart:
		case gst.MessageBuffering:
		case gst.MessageClockLost:
		case gst.MessageClockProvide:
		case gst.MessageDeviceAdded:
		case gst.MessageDeviceChanged:
		case gst.MessageDeviceRemoved:
		case gst.MessageDurationChanged:
		case gst.MessageElement:
		case gst.MessageExtended:
		case gst.MessageHaveContext:
		case gst.MessageInfo:
		case gst.MessageInstantRateRequest:
		case gst.MessageLatency:
		case gst.MessageNeedContext:
		case gst.MessageNewClock:
		case gst.MessageProgress:
		case gst.MessagePropertyNotify:
		case gst.MessageQos:
		case gst.MessageRedirect:
		case gst.MessageRequestState:
		case gst.MessageResetTime:
		case gst.MessageSegmentDone:
		case gst.MessageSegmentStart:
		case gst.MessageStateChanged:
			old, state, pending := msg.ParseStateChanged()

			fmt.Printf("State changed: %s => %s (%s)\n", old, state, pending)
		case gst.MessageStateDirty:
		case gst.MessageStepDone:
		case gst.MessageStepStart:
		case gst.MessageStreamCollection:
		case gst.MessageStreamStart:
		case gst.MessageStreamStatus:
		case gst.MessageStreamsSelected:
		case gst.MessageStructureChange:
		case gst.MessageTag:
		case gst.MessageToc:
		case gst.MessageUnknown:
		case gst.MessageWarning:
		default:
			panic("unexpected gst.MessageType")
		}
		return true
	})

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Block on the main loop

	mainLoop.Run()

	return
}
