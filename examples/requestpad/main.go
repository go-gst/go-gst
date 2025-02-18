package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/pkg/gst"
)

type workflow struct {
	*gst.Pipeline
}

func (w *workflow) newSrc() {
	src := gst.ElementFactoryMake("videotestsrc", "src2")

	src.SetObjectProperty("is-live", true)
	w.Add(src)

	caps := gst.ElementFactoryMake("capsfilter", "caps2")

	caps.SetObjectProperty("caps", gst.CapsFromString("video/x-raw , width=640, height=360"))
	w.Add(caps)

	src.Link(caps)

	// Get a sink pad on compositor
	mixer := w.ByName("mixer")

	pad := mixer.GetRequestPad("sink_%u")
	pad.SetObjectProperty("xpos", 640)
	pad.SetObjectProperty("ypos", 0)

	caps.StaticPad("src").Link(pad)
	caps.SyncStateWithParent()
	src.SyncStateWithParent()

}
func (w *workflow) delSrc() {

	mixer := w.ByName("mixer")

	src := w.ByName("src2")

	caps := w.ByName("caps2")

	pad := mixer.StaticPad("sink_1")
	if pad == nil {
		fmt.Printf("pad is null\n")
		return
	}

	src.SetState(gst.StateNull)
	caps.SetState(gst.StateNull)

	w.Remove(src)
	w.Remove(caps)
	mixer.ReleaseRequestPad(pad)
}

func createPipeline() (*gst.Pipeline, error) {
	gst.Init()
	ret, err := gst.ParseLaunch("videotestsrc ! video/x-raw , capsfilter caps=width=640,height=360 name=caps1 ! compositor name=mixer ! autovideosink")

	if err != nil {
		os.Exit(2)
	}

	var w workflow

	w.Pipeline = ret.(*gst.Pipeline)

	go func() {
		time.Sleep(time.Second)
		w.newSrc()
		time.Sleep(time.Second)
		w.delSrc()
		//runtime.GC()
	}()

	return w.Pipeline, nil
}

func runPipeline(loop *glib.MainLoop, pipeline *gst.Pipeline) error {
	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Add a message watch to the bus to quit on any error
	pipeline.Bus().AddWatch(0, func(bus *gst.Bus, msg *gst.Message) bool {
		var err error

		// If the stream has ended or any element posts an error to the
		// bus, populate error.
		switch msg.Type() {
		case gst.MessageEos:
			err = errors.New("end-of-stream")
		case gst.MessageError:
			// The parsed error implements the error interface, but also
			// contains additional debug information.
			gerr, debug := msg.ParseError()
			fmt.Println("go-gst-debug:", debug)
			err = gerr
		}

		// If either condition triggered an error, log and quit
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			loop.Quit()
			return false
		}

		return true
	})

	// Block on the main loop
	return loop.RunError()
}

func main() {
	loop := glib.NewMainLoop(glib.MainContextDefault(), false)

	pipeline, err := createPipeline()
	if err != nil {
		os.Exit(2)
	}

	runPipeline(loop, pipeline)
}
