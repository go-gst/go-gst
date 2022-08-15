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
	fmt.Println("Running device provider")
	// if len(os.Args) < 2 {
	// 	fmt.Printf("USAGE: %s <uri>\n", os.Args[0])
	// 	os.Exit(1)
	// }

	// uri := os.Args[1]
	fmt.Println("Creating device monitor")

	provider := gst.FindDeviceProviderByName("decklinkdeviceprovider")
	fmt.Println("Created device provider", provider)

	// if err != nil {
	// 	fmt.Println("ERROR:", err)
	// 	os.Exit(2)
	// }

	fmt.Println("listing devices from provider")
	devices := provider.GetDevices()
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
