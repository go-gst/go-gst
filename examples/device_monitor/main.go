// This example uses gstreamer's device monitor api.
//
// https://gstreamer.freedesktop.org/documentation/gstreamer/gstdevicemonitor.html
package main

import (
	"context"
	"fmt"

	"github.com/go-gst/go-gst/pkg/gst"
)

func run() error {
	gst.Init()

	fmt.Println("Creating device monitor")

	monitor := gst.NewDeviceMonitor()
	fmt.Println("Created device monitor", monitor)

	caps := gst.CapsFromString("video/x-raw")

	monitor.AddFilter("Video/Source", caps)

	fmt.Println("Getting device monitor bus")
	bus := monitor.GetBus()
	fmt.Println("Got device monitor bus")

	fmt.Println("Starting device monitor")
	monitor.Start()
	fmt.Println("Started device monitor")
	devices := monitor.GetDevices()
	fmt.Printf("Got %d devices\n", len(devices))
	for i, v := range devices {
		fmt.Printf("Device: %d %s\n", i, v.GetDisplayName())
	}

	for msg := range bus.Messages(context.Background()) {
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
			fmt.Println("Message: ", msg.String())
		}
	}

	return nil
}

func main() {
	run()
}
