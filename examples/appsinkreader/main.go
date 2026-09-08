// This example shows how to use the [gstapp.AppSinkReader] to drain a pipeline with
// the io interfaces of the standard library.
//
// Where the appsink example implements the new-sample callback by hand, the reader
// turns the same element into an [io.ReadCloser]: it pulls samples itself and
// presents their buffers as one continuous byte stream, so the whole output of the
// pipeline can be consumed with a single [io.Copy].
//
// The pipeline transcodes an MP3 to Ogg Vorbis and this program collects the result,
// which is the shape of any "let GStreamer produce bytes for my program" task.
//
// Also see: https://gstreamer.freedesktop.org/documentation/tutorials/basic/short-cutting-the-pipeline.html?gi-language=c
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstapp"
)

const (
	defaultInput = "testdata/grand_project-wonders-of-the-earth-550792.mp3"
	outputFile   = "out.ogg"
)

func createPipeline(input string) (gst.Pipeline, gstapp.AppSink, error) {
	gst.Init()

	// Ogg is chosen because it needs no seeking, which an appsink cannot offer, so
	// the bytes we collect form a complete and playable file. The elements we want to
	// talk to are named so they can be looked up below, and decodebin's dynamic pads
	// are linked by the parser instead of by a pad-added handler.
	ret, err := gst.ParseLaunch(
		"filesrc name=src ! decodebin ! audioconvert ! vorbisenc ! oggmux ! appsink name=sink",
	)

	if err != nil {
		return nil, nil, err
	}

	pipeline := ret.(gst.Pipeline)

	pipeline.GetByName("src").SetObjectProperty("location", input)

	sink := pipeline.GetByName("sink").(gstapp.AppSink)

	return pipeline, sink, nil
}

func readAudio(sink gstapp.AppSink) (int64, error) {
	// Note that emit-signals does not have to be enabled: the reader pulls samples
	// itself and only listens for the eos signal, which the appsink emits either way.
	reader := sink.Reader()
	defer reader.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// io.Copy stops on io.EOF, which the reader returns only on a real end of stream.
	// A pipeline torn down before that yields io.ErrUnexpectedEOF instead, so a
	// truncated stream is never mistaken for a complete one.
	return io.Copy(out, reader)
}

// watchBus reports pipeline messages. On an error it stops the pipeline, which makes
// a blocked read fail instead of waiting for samples that will never arrive.
func watchBus(pipeline gst.Pipeline) {
	for msg := range pipeline.GetBus().Messages(context.Background()) {
		switch msg.Type() {
		case gst.MessageEOS:
			return
		case gst.MessageError:
			debug, gerr := msg.ParseError()
			fmt.Println("Error running pipeline:", gerr.Error(), debug)

			pipeline.SetState(gst.StateNull)

			return
		}
	}
}

func main() {
	input := defaultInput
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	pipeline, sink, err := createPipeline(input)

	if err != nil {
		fmt.Println("Error creating pipeline:", err)
		return
	}

	defer pipeline.SetState(gst.StateNull)

	pipeline.SetState(gst.StatePlaying)

	go watchBus(pipeline)

	n, err := readAudio(sink)

	if err != nil {
		fmt.Println("Error reading audio:", err)
		return
	}

	fmt.Println("Wrote", outputFile, "bytes:", n)
}
