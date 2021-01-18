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

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
)

var srcURI string

func playbin(mainLoop *glib.MainLoop) error {
	if len(os.Args) < 2 {
		return errors.New("Usage: playbin <uri>")
	}

	gst.Init(nil)

	// Create a new playbin and set the URI on it
	playbin, err := gst.NewElement("playbin")
	if err != nil {
		return err
	}
	playbin.Set("uri", os.Args[1])

	// The playbin element itself is a pipeline, so it can be used as one, despite being
	// created from an element factory.
	bus := playbin.GetBus()

	playbin.SetState(gst.StatePlaying)

	bus.AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS:
			mainLoop.Quit()
			return false
		case gst.MessageError:
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG")
				fmt.Println(debug)
			}
			mainLoop.Quit()
			return false
		// Watch state change events
		case gst.MessageStateChanged:
			if _, newState := msg.ParseStateChanged(); newState == gst.StatePlaying {
				bin := gst.ToGstBin(playbin)
				// Generate a dot graph of the pipeline to GST_DEBUG_DUMP_DOT_DIR if defined
				bin.DebugBinToDotFile(gst.DebugGraphShowAll, "PLAYING")
			}

		// Tag messages contain changes to tags on the stream. This can include metadata about
		// the stream such as codecs, artists, albums, etc.
		case gst.MessageTag:
			tags := msg.ParseTags()
			fmt.Println("Tags:")
			if artist, ok := tags.GetString(gst.TagArtist); ok {
				fmt.Println("  Artist:", artist)
			}
			if album, ok := tags.GetString(gst.TagAlbum); ok {
				fmt.Println("  Album:", album)
			}
			if title, ok := tags.GetString(gst.TagTitle); ok {
				fmt.Println("  Title:", title)
			}
		}
		return true
	})

	return mainLoop.RunError()
}

func main() {
	examples.RunLoop(playbin)
}
