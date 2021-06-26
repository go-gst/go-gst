package main

import (
	"errors"
	"fmt"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"os"
	"time"
)

type workflow struct {
	*gst.Pipeline
}

func (w *workflow) newSrc() {
	src, err := gst.NewElementWithName("videotestsrc", "src2")
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	src.Set("is-live", true)
	w.Add(src)

	caps, err := gst.NewElementWithName("capsfilter", "caps2")
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	caps.Set("caps", gst.NewCapsFromString("video/x-raw , width=640, height=360"))
	w.Add(caps)

	src.Link(caps)

	// Get a sink pad on compositor
	mixer, err := w.GetElementByName("mixer")
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	pad := mixer.GetRequestPad("sink_%u")
	pad.SetProperty("xpos", 640)
	pad.SetProperty("ypos", 0)

	caps.GetStaticPad("src").Link(pad)
	caps.SyncStateWithParent()
	src.SyncStateWithParent()

}
func (w *workflow) delSrc() {

	mixer, err := w.GetElementByName("mixer")
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}

	src, err := w.GetElementByName("src2")
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	caps, err := w.GetElementByName("caps2")
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
	pad := mixer.GetStaticPad("sink_1")
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
	gst.Init(nil)
	var err error
	var w workflow
	w.Pipeline, err = gst.NewPipeline("")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	elements, err := gst.NewElementMany("videotestsrc", "capsfilter", "compositor", "autovideosink")
	caps := elements[1]
	caps.SetProperty("caps", gst.NewCapsFromString("video/x-raw , width=640, height=360"))
	caps.SetProperty("name", "caps1")
	mixer := elements[2]
	mixer.SetProperty("name", "mixer")
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, err
	}
	w.AddMany(elements...)
	gst.ElementLinkMany(elements...)

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
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		var err error

		// If the stream has ended or any element posts an error to the
		// bus, populate error.
		switch msg.Type() {
		case gst.MessageEOS:
			err = errors.New("end-of-stream")
		case gst.MessageError:
			// The parsed error implements the error interface, but also
			// contains additional debug information.
			gerr := msg.ParseError()
			fmt.Println("go-gst-debug:", gerr.DebugString())
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
	examples.RunLoop(func(loop *glib.MainLoop) error {
		pipeline, err := createPipeline()
		if err != nil {
			return err
		}
		return runPipeline(loop, pipeline)
	})
}
