// This example demonstrates the use of the decodebin element.
//
// The decodebin element tries to automatically detect the incoming
// format and to autoplug the appropriate demuxers / decoders to handle it.
// and decode it to raw audio, video or subtitles.
// Before the pipeline hasn't been prerolled, the decodebin can't possibly know what
// format it gets as its input. So at first, the pipeline looks like this:
//
//	{filesrc} - {decodebin}
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
//	                       /-[audio]-{audioconvert}-{audioresample}-{autoaudiosink}
//	{filesrc}-{decodebin}-|
//	                       \-[video]-{videoconvert}-{videoscale}-{autovideosink}
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
	"weak"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

var srcFile string

func buildPipeline() (*gst.Pipeline, error) {
	gst.Init()

	pipeline := gst.NewPipeline("")

	src := gst.ElementFactoryMake("filesrc", "")

	decodebin, ok := gst.ElementFactoryMake("decodebin", "").(*gst.Bin) // must cast since we need a weak reference
	if !ok {
		return nil, fmt.Errorf("could not create decodebin")
	}

	src.SetObjectProperty("location", srcFile)

	pipeline.AddMany(src, decodebin)
	src.Link(decodebin)

	// prevent reference cycles with the connect handler:
	weakDecodeBin := weak.Make(decodebin)

	// Connect to decodebin's pad-added signal, that is emitted whenever
	// it found another stream from the input file and found a way to decode it to its raw format.
	// decodebin automatically adds a src-pad for this raw stream, which
	// we can use to build the follow-up pipeline.
	decodebin.ConnectPadAdded(func(srcPad *gst.Pad) {
		// Try to detect whether this is video or audio
		var isAudio, isVideo bool
		caps := srcPad.CurrentCaps()
		for i := 0; i < int(caps.Size()); i++ {
			st := caps.Structure(uint(i))
			if strings.HasPrefix(st.Name(), "audio/") {
				isAudio = true
			}
			if strings.HasPrefix(st.Name(), "video/") {
				isVideo = true
			}
		}

		fmt.Printf("New pad added, is_audio=%v, is_video=%v\n", isAudio, isVideo)

		if !isAudio && !isVideo {
			err := errors.New("could not detect media stream type")
			// We can send errors directly to the pipeline bus if they occur.
			// These will be handled downstream.
			msg := gst.NewMessageError(weakDecodeBin.Value(), err, fmt.Sprintf("Received caps: %s", caps.String()))
			pipeline.Bus().Post(msg)
			return
		}

		if isAudio {
			// decodebin found a raw audiostream, so we build the follow-up pipeline to
			// play it on the default audio playback device (using autoaudiosink).
			audiosink, err := gst.ParseBinFromDescription("queue ! audioconvert ! audioresample ! autoaudiosink", true)
			if err != nil {
				msg := gst.NewMessageError(weakDecodeBin.Value(), err, "Could not create elements for audio pipeline")
				pipeline.Bus().Post(msg)
				return
			}
			pipeline.Add(audiosink)

			// !!ATTENTION!!:
			// This is quite important and people forget it often. Without making sure that
			// the new elements have the same state as the pipeline, things will fail later.
			// They would still be in Null state and can't process data.

			audiosink.SyncStateWithParent()

			// Get the queue element's sink pad and link the decodebin's newly created
			// src pad for the audio stream to it.
			sinkPad := audiosink.StaticPad("sink")
			srcPad.Link(sinkPad)

		} else if isVideo {
			// decodebin found a raw videostream, so we build the follow-up pipeline to
			// display it using the autovideosink.
			videosink, err := gst.ParseBinFromDescription("queue ! videoconvert ! videoscale ! autovideosink", true)
			if err != nil {
				msg := gst.NewMessageError(weakDecodeBin.Value(), err, "Could not create elements for video pipeline")
				pipeline.Bus().Post(msg)
				return
			}
			pipeline.Add(videosink)

			videosink.SyncStateWithParent()

			// Get the queue element's sink pad and link the decodebin's newly created
			// src pad for the video stream to it.
			sinkPad := videosink.StaticPad("sink")
			srcPad.Link(sinkPad)
		}
	})
	return pipeline, nil
}

func runPipeline(loop *glib.MainLoop, pipeline *gst.Pipeline) {
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
	loop.Run()
}

func main() {
	flag.StringVar(&srcFile, "f", "", "The file to decode")
	flag.Parse()
	if srcFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	pipeline, err := buildPipeline()

	if err != nil {
		panic(err)
	}

	mainloop := glib.NewMainLoop(glib.MainContextDefault(), false)

	runPipeline(mainloop, pipeline)
}
