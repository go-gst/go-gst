// This example shows how to use the appsrc element.
//
// Also see: https://gstreamer.freedesktop.org/documentation/tutorials/basic/short-cutting-the-pipeline.html?gi-language=c
package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstapp"
	"github.com/go-gst/go-gst/pkg/gstvideo"
)

const width = 320
const height = 240

func createPipeline() (gst.Pipeline, error) {
	println("Creating pipeline")
	gst.Init()

	// Create a pipeline
	pipeline := gst.NewPipeline("").(gst.Pipeline)

	src := gst.ElementFactoryMake("appsrc", "").(gstapp.AppSrc)
	conv := gst.ElementFactoryMake("videoconvert", "")
	sink := gst.ElementFactoryMake("autovideosink", "")

	// Add the elements to the pipeline and link them
	pipeline.AddMany(src, conv, sink)
	gst.LinkMany(src, conv, sink)

	// Specify the format we want to provide as application into the pipeline
	// by creating a video info with the given format and creating caps from it for the appsrc element.
	videoInfo := gstvideo.NewVideoInfo()

	ok := videoInfo.SetFormat(gstvideo.VideoFormatRGBA, width, height)

	if !ok {
		return nil, fmt.Errorf("failed to set video format")
	}

	videoInfo.SetFramerate(2, 1)

	caps := videoInfo.ToCaps()

	fmt.Println("Caps:", caps.String())

	src.SetObjectProperty("caps", caps)
	src.SetObjectProperty("format", gst.FormatTime)

	// Initialize a frame counter
	var i int

	// Get all 256 colors in the RGB8P palette.
	palette := gstvideo.VideoFormatGetPalette(gstvideo.VideoFormatRGB8P)

	// Since our appsrc element operates in pull mode (it asks us to provide data),
	// we add a handler for the need-data callback and provide new data from there.
	// In our case, we told gstreamer that we do 2 frames per second. While the
	// buffers of all elements of the pipeline are still empty, this will be called
	// a couple of times until all of them are filled. After this initial period,
	// this handler will be called (on average) twice per second.
	src.ConnectNeedData(func(self gstapp.AppSrc, _ uint) {

		// If we've reached the end of the palette, end the stream.
		if i == len(palette) {
			src.EndOfStream()
			return
		}

		fmt.Println("Producing frame:", i)

		// Create a buffer that can hold exactly one video RGBA frame.
		buffer := gst.NewBufferAllocate(nil, uint(videoInfo.GetSize()), nil)

		// For each frame we produce, we set the timestamp when it should be displayed
		// The autovideosink will use this information to display the frame at the right time.
		buffer.SetPTS(gst.ClockTime(time.Duration(i) * 500 * time.Millisecond))

		// Produce an image frame for this iteration.
		pixels := produceImageFrame(palette[i])

		// At this point, buffer is only a reference to an existing memory region somewhere.
		// When we want to access its content, we have to map it while requesting the required
		// mode of access (read, read/write).
		// See: https://gstreamer.freedesktop.org/documentation/plugin-development/advanced/allocation.html
		mapped, ok := buffer.Map(gst.MapWrite)
		if !ok {
			panic("Failed to map buffer")
		}
		_, err := mapped.Write(pixels)
		if err != nil {
			println("Failed to write to buffer:", err)
			panic("Failed to write to buffer")
		}

		mapped.Unmap()

		// Push the buffer onto the pipeline.
		self.PushBuffer(buffer)

		i++
	})

	return pipeline, nil
}

func produceImageFrame(c color.Color) []uint8 {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, c)
		}
	}

	return img.Pix
}

func mainLoop(pipeline gst.Pipeline) error {
	// Start the pipeline

	pipeline.SetState(gst.StatePlaying)

	for msg := range pipeline.GetBus().Messages(context.Background()) {
		switch msg.Type() {
		case gst.MessageEos:
			return nil
		case gst.MessageError:
			debug, gerr := msg.ParseError()
			if debug != "" {
				fmt.Println(gerr.Error(), debug)
			}
			return gerr
		default:
			fmt.Println(msg)
		}

		pipeline.DebugBinToDotFileWithTs(gst.DebugGraphShowVerbose, "pipeline")
	}

	return fmt.Errorf("unexpected end of messages without EOS")
}

func main() {
	pipeline, err := createPipeline()

	if err != nil {
		fmt.Println("Error creating pipeline:", err)
		return
	}

	err = mainLoop(pipeline)

	if err != nil {
		fmt.Println("Error running pipeline:", err)
	}
}
