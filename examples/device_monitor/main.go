// This example uses gstreamer's discoverer api.
//
// https://gstreamer.freedesktop.org/data/doc/gstreamer/head/gst-plugins-base-libs/html/GstDiscoverer.html
// To detect as much information from a given URI.
// The amount of time that the discoverer is allowed to use is limited by a timeout.
// This allows to handle e.g. network problems gracefully. When the timeout hits before
// discoverer was able to detect anything, discoverer will report an error.
// In this example, we catch this error and stop the application.
// Discovered information could for example contain the stream's duration or whether it is
// seekable (filesystem) or not (some http servers).
package main

import (
	"fmt"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
)

func runPipeline(loop *glib.MainLoop) error {

	gst.Init(nil)
	fmt.Println("Running device monitor")
	// if len(os.Args) < 2 {
	// 	fmt.Printf("USAGE: %s <uri>\n", os.Args[0])
	// 	os.Exit(1)
	// }

	// uri := os.Args[1]
	fmt.Println("Creating device monitor")

	monitor := gst.NewDeviceMonitor()
	fmt.Println("Created device monitor", monitor)

	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// 	os.Exit(2)
	// }
	caps := gst.NewCapsFromString("video/x-raw")

	monitor.AddFilter("Video/Source", caps)

	fmt.Println("Getting device monitor bus")
	bus := monitor.GetBus()
	fmt.Println("Got device monitor bus", bus)

	bus.AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageDeviceAdded:
			message := msg.ParseDeviceAdded().GetDisplayName()
			fmt.Println("Added: ", message)
		case gst.MessageDeviceRemoved:
			message := msg.ParseDeviceRemoved().GetDisplayName()
			fmt.Println("Removed: ", message)
		default:
			// All messages implement a Stringer. However, this is
			// typically an expensive thing to do and should be avoided.
			fmt.Println("Type: ", msg.Type())
			fmt.Println("Message: ", msg)
		}
		return true
	})

	monitor.Start()
	fmt.Println("Started device monitor")
	devices := monitor.GetDevices()
	for i, v := range devices {
		fmt.Printf("Device: %d %s\n", i, v.GetDisplayName())
	}

	loop.Run()

	return nil
}

func main() {
	examples.RunLoop(func(loop *glib.MainLoop) error {
		return runPipeline(loop)
	})
}
