// This example shows how to use the [gstapp.AppSrcWriter] to feed a pipeline with
// the io interfaces of the standard library.
//
// Where the appsrc example implements the need-data callback by hand, the writer
// turns the same element into an [io.WriteCloser]: it honours need-data and
// enough-data internally, so [io.Copy] blocks instead of letting the queue grow, and
// Close sends the EOS event.
//
// No caps are set on the appsrc, so decodebin has to typefind the bytes we push. The
// pipeline decodes them to a WAV file, which is the shape of any "my program has the
// bytes, let GStreamer deal with them" task.
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
	outputFile   = "out.wav"
)

func createPipeline() (gst.Pipeline, gstapp.AppSrc, error) {
	gst.Init()

	// The elements we want to talk to are named so they can be looked up below, and
	// decodebin's dynamic pads are linked by the parser instead of by a pad-added
	// handler.
	ret, err := gst.ParseLaunch(
		"appsrc name=src ! decodebin ! audioconvert ! wavenc ! filesink name=sink",
	)

	if err != nil {
		return nil, nil, err
	}

	pipeline := ret.(gst.Pipeline)

	pipeline.GetByName("sink").SetObjectProperty("location", outputFile)

	src := pipeline.GetByName("src").(gstapp.AppSrc)

	return pipeline, src, nil
}

func writeAudio(src gstapp.AppSrc, input string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()

	writer := src.Writer()

	n, err := io.Copy(writer, in)
	if err != nil {
		return err
	}

	fmt.Println("Pushed bytes:", n)

	// Closing the writer emits EOS, which is what makes wavenc finalize its header
	// and the pipeline post the EOS message the main loop waits for.
	return writer.Close()
}

func mainLoop(pipeline gst.Pipeline) error {
	for msg := range pipeline.GetBus().Messages(context.Background()) {
		switch msg.Type() {
		case gst.MessageEOS:
			return nil
		case gst.MessageError:
			debug, gerr := msg.ParseError()
			if debug != "" {
				fmt.Println(gerr.Error(), debug)
			}
			return gerr
		}
	}

	return fmt.Errorf("unexpected end of messages without EOS")
}

func main() {
	input := defaultInput
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	pipeline, src, err := createPipeline()

	if err != nil {
		fmt.Println("Error creating pipeline:", err)
		return
	}

	defer pipeline.SetState(gst.StateNull)

	// The appsrc only accepts buffers once the pipeline has been started, so push
	// from a goroutine while the main loop watches the bus.
	pipeline.SetState(gst.StatePlaying)

	writeErr := make(chan error, 1)

	go func() {
		writeErr <- writeAudio(src, input)
	}()

	if err := mainLoop(pipeline); err != nil {
		fmt.Println("Error running pipeline:", err)
		return
	}

	if err := <-writeErr; err != nil {
		fmt.Println("Error writing audio:", err)
		return
	}

	fmt.Println("Wrote", outputFile)
}
