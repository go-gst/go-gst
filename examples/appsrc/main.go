// This example shows how to use the appsrc element.
package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
	"github.com/tinyzimmer/go-gst/gst/video"
)

const width = 320
const height = 240

func createPipeline() (*gst.Pipeline, error) {
	gst.Init(nil)

	// Create a pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	// Create the elements
	elems, err := gst.NewElementMany("appsrc", "videoconvert", "autovideosink")
	if err != nil {
		return nil, err
	}

	// Add the elements to the pipeline and link them
	pipeline.AddMany(elems...)
	gst.ElementLinkMany(elems...)

	// Get the app sourrce from the first element returned
	src := app.SrcFromElement(elems[0])

	// Specify the format we want to provide as application into the pipeline
	// by creating a video info with the given format and creating caps from it for the appsrc element.
	videoInfo := video.NewInfo().
		WithFormat(video.FormatRGBA, width, height).
		WithFPS(gst.Fraction(2, 1))

	src.SetCaps(videoInfo.ToCaps())
	src.SetProperty("format", gst.FormatTime)

	// Initialize a frame counter
	var i int

	// Get all 256 colors in the RGB8P palette.
	palette := video.FormatRGB8P.Palette()

	// Since our appsrc element operates in pull mode (it asks us to provide data),
	// we add a handler for the need-data callback and provide new data from there.
	// In our case, we told gstreamer that we do 2 frames per second. While the
	// buffers of all elements of the pipeline are still empty, this will be called
	// a couple of times until all of them are filled. After this initial period,
	// this handler will be called (on average) twice per second.
	src.SetCallbacks(&app.SourceCallbacks{
		NeedDataFunc: func(self *app.Source, _ uint) {

			// If we've reached the end of the palette, end the stream.
			if i == len(palette) {
				src.EndStream()
				return
			}

			fmt.Println("Producing frame:", i)

			// Create a buffer that can hold exactly one video RGBA frame.
			buffer := gst.NewBufferWithSize(videoInfo.Size())

			// For each frame we produce, we set the timestamp when it should be displayed
			// The autovideosink will use this information to display the frame at the right time.
			buffer.SetPresentationTimestamp(time.Duration(i) * 500 * time.Millisecond)

			// Produce an image frame for this iteration.
			pixels := produceImageFrame(palette[i])

			// At this point, buffer is only a reference to an existing memory region somewhere.
			// When we want to access its content, we have to map it while requesting the required
			// mode of access (read, read/write).
			// See: https://gstreamer.freedesktop.org/documentation/plugin-development/advanced/allocation.html
			//
			// There are convenience wrappers for building buffers directly from byte sequences as
			// well.
			buffer.Map(gst.MapWrite).WriteData(pixels)
			defer buffer.Unmap()

			// Push the buffer onto the pipeline.
			self.PushBuffer(buffer)

			i++
		},
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

func handleMessage(msg *gst.Message) error {
	defer msg.Unref() // Messages are a good candidate for trying out runtime finalizers

	switch msg.Type() {
	case gst.MessageEOS:
		return app.ErrEOS
	case gst.MessageError:
		gerr := msg.ParseError()
		if debug := gerr.DebugString(); debug != "" {
			fmt.Println(debug)
		}
		return gerr
	}

	return nil
}

func mainLoop(pipeline *gst.Pipeline) error {

	defer pipeline.Destroy() // Will stop and unref the pipeline when this function returns

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Retrieve the bus from the pipeline
	bus := pipeline.GetPipelineBus()

	// Loop over messsages from the pipeline
	for {
		msg := bus.TimedPop(time.Duration(-1))
		if msg == nil {
			break
		}
		if err := handleMessage(msg); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	examples.Run(func() error {
		var pipeline *gst.Pipeline
		var err error
		if pipeline, err = createPipeline(); err != nil {
			return err
		}
		return mainLoop(pipeline)
	})
}
