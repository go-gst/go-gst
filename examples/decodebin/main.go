// This example demonstrates the use of the decodebin element.
//
// The decodebin element tries to automatically detect the incoming
// format and to autoplug the appropriate demuxers / decoders to handle it.
// and decode it to raw audio, video or subtitles.
// Before the pipeline hasn't been prerolled, the decodebin can't possibly know what
// format it gets as its input. So at first, the pipeline looks like this:
//
//   {filesrc} - {decodebin}
//
// As soon as the decodebin has detected the stream format, it will try to decode every
// contained stream to its raw format.
// The application connects a signal-handler to decodebin's pad-added signal, which tells us
// whenever the decodebin provided us with another contained (raw) stream from the input file.
//
// This application supports audio and video streams. Video streams are
// displayed using an autovideosink, and audiostreams are played back using autoaudiosink.
// So for a file that contains one audio and one video stream,
// the pipeline looks like the following:
//
//                          /-[audio]-{audioconvert}-{audioresample}-{autoaudiosink}
//   {filesrc}-{decodebin}-|
//                          \-[video]-{videoconvert}-{videoscale}-{autovideosink}
//
// Both auto-sinks at the end automatically select the best available (actual) sink. Since the
// selection of available actual sinks is platform specific
// (like using pulseaudio for audio output on linux, e.g.),
// we need to add the audioconvert and audioresample elements before handing the stream to the
// autoaudiosink, because we need to make sure, that the stream is always supported by the actual sink.
// Especially Windows APIs tend to be quite picky about samplerate and sample-format.
// The same applies to videostreams.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
)

var srcFile string

func buildPipeline() (*gst.Pipeline, error) {
	gst.Init(nil)

	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	src, err := gst.NewElement("filesrc")
	if err != nil {
		return nil, err
	}

	decodebin, err := gst.NewElement("decodebin")
	if err != nil {
		return nil, err
	}

	src.Set("location", srcFile)

	pipeline.AddMany(src, decodebin)
	src.Link(decodebin)

	// Connect to decodebin's pad-added signal, that is emitted whenever
	// it found another stream from the input file and found a way to decode it to its raw format.
	// decodebin automatically adds a src-pad for this raw stream, which
	// we can use to build the follow-up pipeline.
	decodebin.Connect("pad-added", func(self *gst.Element, srcPad *gst.Pad) {

		// Try to detect whether this is video or audio
		var isAudio, isVideo bool
		caps := srcPad.GetCurrentCaps()
		for i := 0; i < caps.GetSize(); i++ {
			st := caps.GetStructureAt(i)
			if strings.HasPrefix(st.Name(), "audio/") {
				isAudio = true
			}
			if strings.HasPrefix(st.Name(), "video/") {
				isVideo = true
			}
		}

		fmt.Printf("New pad added, is_audio=%v, is_video=%v\n", isAudio, isVideo)

		if !isAudio && !isVideo {
			err := errors.New("Could not detect media stream type")
			// We can send errors directly to the pipeline bus if they occur.
			// These will be handled downstream.
			msg := gst.NewErrorMessage(self, gst.NewGError(1, err), fmt.Sprintf("Received caps: %s", caps.String()), nil)
			pipeline.GetPipelineBus().Post(msg)
			return
		}

		if isAudio {
			// decodebin found a raw audiostream, so we build the follow-up pipeline to
			// play it on the default audio playback device (using autoaudiosink).
			elements, err := gst.NewElementMany("queue", "audioconvert", "audioresample", "autoaudiosink")
			if err != nil {
				// We can create custom errors (with optional structures) and send them to the pipeline bus.
				// The first argument reflects the source of the error, the second is the error itself, followed by a debug string.
				msg := gst.NewErrorMessage(self, gst.NewGError(2, err), "Could not create elements for audio pipeline", nil)
				pipeline.GetPipelineBus().Post(msg)
				return
			}
			pipeline.AddMany(elements...)
			gst.ElementLinkMany(elements...)

			// !!ATTENTION!!:
			// This is quite important and people forget it often. Without making sure that
			// the new elements have the same state as the pipeline, things will fail later.
			// They would still be in Null state and can't process data.
			for _, e := range elements {
				e.SyncStateWithParent()
			}

			// The queue was the first element returned above
			queue := elements[0]
			// Get the queue element's sink pad and link the decodebin's newly created
			// src pad for the audio stream to it.
			sinkPad := queue.GetStaticPad("sink")
			srcPad.Link(sinkPad)

		} else if isVideo {
			// decodebin found a raw videostream, so we build the follow-up pipeline to
			// display it using the autovideosink.
			elements, err := gst.NewElementMany("queue", "videoconvert", "videoscale", "autovideosink")
			if err != nil {
				msg := gst.NewErrorMessage(self, gst.NewGError(2, err), "Could not create elements for video pipeline", nil)
				pipeline.GetPipelineBus().Post(msg)
				return
			}
			pipeline.AddMany(elements...)
			gst.ElementLinkMany(elements...)

			for _, e := range elements {
				e.SyncStateWithParent()
			}

			queue := elements[0]
			// Get the queue element's sink pad and link the decodebin's newly created
			// src pad for the video stream to it.
			sinkPad := queue.GetStaticPad("sink")
			srcPad.Link(sinkPad)
		}
	})
	return pipeline, nil
}

func runPipeline(loop *gst.MainLoop, pipeline *gst.Pipeline) error {
	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)
	// Stop and cleanup the pipeline when we exit
	defer pipeline.Destroy()

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
	flag.StringVar(&srcFile, "f", "", "The file to decode")
	flag.Parse()
	if srcFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	examples.RunLoop(func(loop *gst.MainLoop) error {
		pipeline, err := buildPipeline()
		if err != nil {
			return err
		}
		return runPipeline(loop, pipeline)
	})
}
