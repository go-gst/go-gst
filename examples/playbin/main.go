// This example demonstrates GStreamer's playbin element.
//
// This element takes an arbitrary URI as parameter, and if there is a source
// element within gstreamer, that supports this uri, the playbin will try
// to automatically create a pipeline that properly plays this media source.
// For this, the playbin internally relies on more bin elements, like the
// autovideosink and the decodebin.
// Essentially, this element is a single-element pipeline able to play
// any format from any uri-addressable source that gstreamer supports.
// Much of the playbin's behavior can be controlled by so-called flags, as well
// as the playbin's properties and signals.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

func playbin() error {
	gst.Init()

	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	if len(os.Args) < 2 {
		return errors.New("usage: playbin <uri>")
	}

	gst.Init()

	// Create a new playbin and set the URI on it
	ret := gst.ElementFactoryMake("playbin", "")
	if ret != nil {
		return fmt.Errorf("could not create playbin")
	}

	playbin := ret.(*gst.Pipeline)

	playbin.SetObjectProperty("uri", os.Args[1])

	// The playbin element itself is a pipeline, so it can be used as one, despite being
	// created from an element factory.
	bus := playbin.Bus()

	playbin.SetState(gst.StatePlaying)

	bus.AddWatch(0, func(bus *gst.Bus, msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEos:
			mainLoop.Quit()
			return false
		case gst.MessageError:
			err, debug := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug != "" {
				fmt.Println("DEBUG")
				fmt.Println(debug)
			}
			mainLoop.Quit()
			return false
		// Watch state change events
		case gst.MessageStateChanged:
			if _, newState, _ := msg.ParseStateChanged(); newState == gst.StatePlaying {
				// Generate a dot graph of the pipeline to GST_DEBUG_DUMP_DOT_DIR if defined
				gst.DebugBinToDotFile(&playbin.Bin, gst.DebugGraphShowAll, "PLAYING")
			}

		// Tag messages contain changes to tags on the stream. This can include metadata about
		// the stream such as codecs, artists, albums, etc.
		case gst.MessageTag:
			tags := msg.ParseTag()
			fmt.Println("Tags:")
			if artist, ok := tags.String(gst.TAG_ARTIST); ok {
				fmt.Println("  Artist:", artist)
			}
			if album, ok := tags.String(gst.TAG_ALBUM); ok {
				fmt.Println("  Album:", album)
			}
			if title, ok := tags.String(gst.TAG_TITLE); ok {
				fmt.Println("  Title:", title)
			}
		}
		return true
	})

	mainLoop.Run()

	return nil
}

func main() {
	if err := playbin(); err != nil {
		fmt.Println(err)
	}
}
