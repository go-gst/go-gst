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
	"os"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstpbutils"
)

func main() {

	gst.Init()

	if len(os.Args) < 2 {
		fmt.Printf("USAGE: %s <uri>\n", os.Args[0])
		os.Exit(1)
	}

	uri := os.Args[1]

	discoverer, err := gstpbutils.NewDiscoverer(gst.ClockTime(time.Second * 15))
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(2)
	}

	info, err := discoverer.DiscoverURI(uri)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(3)
	}

	printDiscovererInfo(info)
}

func printDiscovererInfo(info *gstpbutils.DiscovererInfo) {
	fmt.Println("URI:", info.URI())
	fmt.Println("Duration:", info.Duration())

	printTags(info)
	printStreamInfo(info.StreamInfo())

	children := info.StreamList()
	fmt.Println("Children streams:")
	for _, child := range children {
		printStreamInfo(child)
	}
}

func printTags(info *gstpbutils.DiscovererInfo) {
	fmt.Println("Tags:")
	tags := info.Tags()
	if tags != nil {
		fmt.Println("  ", tags)
		return
	}
	fmt.Println("  no tags")
}

func printStreamInfo(info *gstpbutils.DiscovererStreamInfo) {
	if info == nil {
		return
	}
	fmt.Println("Stream: ")
	fmt.Println("  Stream id:", info.StreamID())
	if caps := info.Caps(); caps != nil {
		fmt.Println("  Format:", caps)
	}
}
