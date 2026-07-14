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
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-gst/go-gst/pkg/gst"
)

func playbin() error {
	gst.Init()

	if len(os.Args) < 2 {
		return errors.New("usage: playbin <uri>")
	}

	// Create a new playbin and set the URI on it
	ret := gst.ElementFactoryMake("playbin", "")
	if ret != nil {
		return fmt.Errorf("could not create playbin")
	}

	playbin := ret.(gst.Pipeline)

	playbin.SetObjectProperty("uri", os.Args[1])

	// The playbin element itself is a pipeline, so it can be used as one, despite being
	// created from an element factory.
	bus := playbin.GetBus()

	playbin.SetState(gst.StatePlaying)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for msg := range bus.Messages(ctx) {
		switch msg.Type() {
		case gst.MessageEOS:
			return nil
		case gst.MessageError:
			debug, err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug != "" {
				fmt.Println("DEBUG")
				fmt.Println(debug)
			}
			return nil
		// Watch state change events
		case gst.MessageStateChanged:
			if _, newState, _ := msg.ParseStateChanged(); newState == gst.StatePlaying {
				// Generate a dot graph of the pipeline to GST_DEBUG_DUMP_DOT_DIR if defined
				playbin.DebugBinToDotFile(gst.DebugGraphShowAll, "PLAYING")
			}

		// Tag messages contain changes to tags on the stream. This can include metadata about
		// the stream such as codecs, artists, albums, etc.
		case gst.MessageTag:
			tags := msg.ParseTag()
			fmt.Println("Tags:")
			if artist, ok := tags.GetString(gst.TAG_ARTIST); ok {
				fmt.Println("  Artist:", artist)
			}
			if album, ok := tags.GetString(gst.TAG_ALBUM); ok {
				fmt.Println("  Album:", album)
			}
			if title, ok := tags.GetString(gst.TAG_TITLE); ok {
				fmt.Println("  Title:", title)
			}
		}
		return nil
	}

	return nil
}

func main() {
	if err := playbin(); err != nil {
		fmt.Println(err)
	}
}
