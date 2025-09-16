// This example demonstrates using gstreamer to convert a video stream into image frames
// and then encoding those frames to a gif.
package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstapp"
	"github.com/go-gst/go-gst/pkg/gstvideo"
)

var srcFile string
var outFile string

const width = 320
const height = 240

func encodeGif() error {
	gst.Init()

	// Initialize an empty buffer for the encoded gif images.
	outGif := &gif.GIF{
		Image: make([]*image.Paletted, 0),
		Delay: make([]int, 0),
	}

	// Create a new pipeline instance
	pipeline := gst.NewPipeline("").(gst.Pipeline)

	filesrc := gst.ElementFactoryMake("filesrc", "")
	decodebin := gst.ElementFactoryMake("decodebin", "")

	// Add the elements to the pipeline.
	pipeline.AddMany(filesrc, decodebin)

	// Set the location of the source file the filesrc element and link it to the
	// decodebin.
	filesrc.SetObjectProperty("location", srcFile)
	gst.LinkMany(filesrc, decodebin)

	// Conncet to decodebin's pad-added signal to build the rest of the pipeline
	// dynamically. For more information on why this is needed, see the decodebin
	// example.
	decodebin.ConnectPadAdded(func(self gst.Element, srcPad gst.Pad) {
		// Build out the rest of the elements for the pipeline pipeline.

		queue := gst.ElementFactoryMake("queue", "")
		videoconvert := gst.ElementFactoryMake("videoconvert", "")
		videoscale := gst.ElementFactoryMake("videoscale", "")
		videorate := gst.ElementFactoryMake("videorate", "")
		jpegenc := gst.ElementFactoryMake("jpegenc", "")

		// Add the elements to the pipeline and sync their state with the pipeline
		pipeline.AddMany(queue, videoconvert, videoscale, videorate, jpegenc)

		queue.SyncStateWithParent()
		videoconvert.SyncStateWithParent()
		videoscale.SyncStateWithParent()
		videorate.SyncStateWithParent()
		jpegenc.SyncStateWithParent()

		// Start linking elements

		queue.Link(videoconvert)

		// We need to tell the pipeline the output format we want. Here we are going to request
		// RGBx color with predefined boundaries and 5 frames per second.
		videoInfo := gstvideo.NewVideoInfo()
		videoInfo.SetFormat(gstvideo.VideoFormatRgbx, width, height)
		videoInfo.SetFramerate(5, 1)

		// videoconvert.LinkFiltered(videoscale, videoInfo.ToCaps())
		gst.LinkMany(videoconvert, videoscale, videorate)

		videorate.LinkFiltered(jpegenc, videoInfo.ToCaps())

		// Create an app sink that we are going to use to pull images from the pipeline
		// one at a time. (An error can happen here too, but for the sake of brevity...)
		appSink := gst.ElementFactoryMake("appsink", "").(gstapp.AppSink)
		pipeline.Add(appSink)
		jpegenc.Link(appSink)
		appSink.SyncStateWithParent()
		appSink.SetWaitOnEos(false)

		// We can query the decodebin for the duration of the video it received. We can then
		// use this value to calculate the total number of frames we expect to produce.
		query := gst.NewQueryDuration(gst.FormatTime)
		if ok := self.Query(query); !ok {
			self.MessageError(0, int32(gst.LibraryErrorFailed), "Failed to query video duration from decodebin", "")
			return
		}

		// Fetch the result from the query.
		_, duration := query.ParseDuration()

		// This value is in nanoseconds. Since we told the videorate element to produce 5 frames
		// per second, we know the total frames will be (duration / 1e+9) * 5.
		totalFrames := int((time.Duration(duration) * time.Nanosecond).Seconds()) * 5

		// Getting data out of the sink is done by setting callbacks. Each new sample
		// will be a new jpeg image from the pipeline.
		var frameNum int

		appSink.ConnectEos(func(self gstapp.AppSink) {
			fmt.Println("\nWriting the results of the gif to", outFile)
			file, err := os.Create(outFile)
			if err != nil {
				fmt.Println("Could not create output file:", err)
				return
			}
			defer file.Close()
			if err := gif.EncodeAll(file, outGif); err != nil {
				fmt.Println("Could not encode images to gif format!", err)
			}
			// Signal the pipeline that we've completed EOS.
			// (this should not be required, need to investigate)
			pipeline.GetBus().Post(gst.NewMessageEos(appSink))
		})

		appSink.ConnectNewSample(func(sink gstapp.AppSink) gst.FlowReturn {
			// Increment the frame number counter
			frameNum++

			if frameNum > totalFrames {
				// If we've reached the total number of frames we are expecting. We can
				// signal the main loop to quit.
				// This needs to be done from a goroutine to not block the app sink
				// callback.
				return gst.FlowEos
			}

			// Pull the sample from the sink
			sample := sink.PullSample()
			if sample == nil {
				return gst.FlowOK
			}

			fmt.Printf("\033[2K\r")
			fmt.Printf("Processing image frame %d/%d", frameNum, totalFrames)

			// Retrieve the buffer from the sample.
			buffer := sample.GetBuffer()

			mapped, ok := buffer.Map(gst.MapRead)
			if !ok {
				panic("Failed to map buffer")
			}

			// mapped buffers implement io.Reader
			img, err := jpeg.Decode(mapped)
			if err != nil {
				self.MessageError(gst.LibraryErrorQuark(), int32(gst.LibraryErrorFailed), "Error decoding jpeg frame", err.Error())
				return gst.FlowError
			}

			// Create a new paletted image with the same bounds as the pulled one
			frame := image.NewPaletted(img.Bounds(), gstvideo.VideoFormatGetPalette(gstvideo.VideoFormatRGB8P))

			// Iterate the bounds of the image and set the pixels in their correct place.
			for x := 1; x <= img.Bounds().Dx(); x++ {
				for y := 1; y <= img.Bounds().Dy(); y++ {
					frame.Set(x, y, img.At(x, y))
				}
			}

			// Append the image data to the gif
			outGif.Image = append(outGif.Image, frame)
			outGif.Delay = append(outGif.Delay, 0)
			return gst.FlowOK
		})

		// Link the src pad to the queue
		srcPad.Link(queue.GetStaticPad("sink"))
	})

	fmt.Println("Encoding video to gif")

	// Now that the pipeline is all set up we can start it.
	pipeline.SetState(gst.StatePlaying)

	for msg := range pipeline.GetBus().Messages(context.Background()) {
		switch msg.Type() {
		case gst.MessageEos:
			return nil
		case gst.MessageError:
			debug, gerr := msg.ParseError()
			fmt.Println("ERROR:", gerr.Error())
			if debug != "" {
				fmt.Println("DEBUG")
				fmt.Println(debug)
			}

			return gerr
		}
	}

	return fmt.Errorf("unexpected end of messages without EOS")
}

func main() {
	// Add flag arguments
	flag.StringVar(&srcFile, "i", "", "The video to encode to gif. This argument is required.")
	flag.StringVar(&outFile, "o", "", "The file to output the gif to. By default a file is created in this directory with the same name as the input.")

	// Parse the command line
	flag.Parse()

	// Make sure the user provided a source file
	if srcFile == "" {
		flag.Usage()
		fmt.Println("The input file cannot be empty!")
		os.Exit(1)
	}

	// If the user did not provide a destination file, generate one.
	if outFile == "" {
		base := path.Base(srcFile)
		spl := strings.Split(base, ".")
		if len(spl) < 3 {
			outFile = spl[0]
		} else {
			outFile = strings.Join(spl[:len(spl)-2], ".")
		}
		outFile = outFile + ".gif"
	}

	encodeGif()
}
