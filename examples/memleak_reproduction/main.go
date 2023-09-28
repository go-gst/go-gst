package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: playbin <uri>")
	}

	fmt.Printf("opening file %s", os.Args[1])

	gst.Init(nil)

	pipeline, err := gst.NewPipeline("")

	if err != nil {
		panic(err)
	}

	// Create a new playbin and set the URI on it
	playbin, err := gst.NewElement("uridecodebin")
	if err != nil {
		panic(err)
	}
	playbin.Set("uri", os.Args[1])

	output, err := gst.NewElement("autoaudiosink")
	if err != nil {
		panic(err)
	}

	pipeline.AddMany(playbin, output)

	handle, _ := playbin.Connect("pad-added", func(e *gst.Element, p *gst.Pad) {
		err := e.Link(output)

		if err != nil {
			panic(fmt.Sprintf("got error during linking: %v", err))
		}
	})

	// The playbin element itself is a pipeline, so it can be used as one, despite being
	// created from an element factory.
	bus := pipeline.GetBus()

	pipeline.SetState(gst.StatePlaying)

	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

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

	mainLoop.RunError()

	playbin.HandlerDisconnect(handle)

	pipeline.BlockSetState(gst.StateNull)

	runtime.GC()
	runtime.GC()
	runtime.GC()
	runtime.GC()

	gst.Deinit()

	runtime.GC()
	runtime.GC()
	runtime.GC()
	runtime.GC()
}
