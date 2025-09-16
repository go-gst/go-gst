// This example uses gstreamer's device provider api.
//
// https://gstreamer.freedesktop.org/documentation/gstreamer/gstdeviceprovider.html
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-gst/go-gst/pkg/gst"
)

func runPipeline() error {

	gst.Init()
	fmt.Println("Running device provider")

	fmt.Println("Creating device monitor")

	provider := gst.DeviceProviderFactoryGetByName("avfdeviceprovider")
	fmt.Println("Created device provider")

	if provider == nil {
		fmt.Println("No provider found")
		os.Exit(2)
	}

	fmt.Println("Starting device monitor")
	provider.Start()
	fmt.Println("Started device monitor")

	fmt.Println("listing devices from provider")
	devices := provider.GetDevices()
	for i, v := range devices {
		fmt.Printf("Device: %d %s\n", i, v.GetDisplayName())
	}

	fmt.Println("Getting device provider bus")
	bus := provider.GetBus()
	fmt.Println("Got device provider bus")

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
	runPipeline()
}
