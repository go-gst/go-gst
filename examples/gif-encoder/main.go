// This example demonstrates using gstreamer to convert a video stream into image frames
// and then encoding those frames to a gif.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
	"github.com/tinyzimmer/go-gst/gst/video"
)

var srcFile string
var outFile string

const width = 320
const height = 240

func encodeGif(mainLoop *glib.MainLoop) error {
	gst.Init(nil)

	// Initialize an empty buffer for the encoded gif images.
	outGif := &gif.GIF{
		Image: make([]*image.Paletted, 0),
		Delay: make([]int, 0),
	}

	// Create a new pipeline instance
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return err
	}

	// Create a filesrc and a decodebin element for the pipeline.
	elements, err := gst.NewElementMany("filesrc", "decodebin")
	if err != nil {
		return nil
	}

	filesrc := elements[0]   // The filsrc is the first element returned.
	decodebin := elements[1] // The decodebin is the second element returned.

	// Add the elements to the pipeline.
	pipeline.AddMany(elements...)

	// Set the location of the source file the filesrc element and link it to the
	// decodebin.
	filesrc.Set("location", srcFile)
	gst.ElementLinkMany(filesrc, decodebin)

	// Conncet to decodebin's pad-added signal to build the rest of the pipeline
	// dynamically. For more information on why this is needed, see the decodebin
	// example.
	decodebin.Connect("pad-added", func(self *gst.Element, srcPad *gst.Pad) {
		// Build out the rest of the elements for the pipeline pipeline.
		elements, err := gst.NewElementMany("queue", "videoconvert", "videoscale", "videorate", "jpegenc")
		if err != nil {
			// The Bus PostError method is a convenience wrapper for building rich messages and sending them
			// down the pipeline. The below call will create a new error message, populate the debug info
			// with a stack trace from this goroutine, and add additional details from the provided error.
			self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, "Failed to build elements for the linked pipeline", err.Error())
			return
		}

		// Add the elements to the pipeline and sync their state with the pipeline
		pipeline.AddMany(elements...)
		for _, e := range elements {
			e.SyncStateWithParent()
		}

		// Retrieve direct references to the elements for clarity.
		queue := elements[0]
		videoconvert := elements[1]
		videoscale := elements[2]
		videorate := elements[3]
		jpegenc := elements[4]

		// Start linking elements

		queue.Link(videoconvert)

		// We need to tell the pipeline the output format we want. Here we are going to request
		// RGBx color with predefined boundaries and 5 frames per second.
		videoInfo := video.NewInfo().
			WithFormat(video.FormatRGBx, width, height).
			WithFPS(gst.Fraction(5, 1))

		// videoconvert.LinkFiltered(videoscale, videoInfo.ToCaps())
		gst.ElementLinkMany(videoconvert, videoscale, videorate)

		videorate.LinkFiltered(jpegenc, videoInfo.ToCaps())

		// Create an app sink that we are going to use to pull images from the pipeline
		// one at a time. (An error can happen here too, but for the sake of brevity...)
		appSink, _ := app.NewAppSink()
		pipeline.Add(appSink.Element)
		jpegenc.Link(appSink.Element)
		appSink.SyncStateWithParent()
		appSink.SetWaitOnEOS(false)

		// We can query the decodebin for the duration of the video it received. We can then
		// use this value to calculate the total number of frames we expect to produce.
		query := gst.NewDurationQuery(gst.FormatTime)
		if ok := self.Query(query); !ok {
			self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, "Failed to query video duration from decodebin", err.Error())
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
		appSink.SetCallbacks(&app.SinkCallbacks{
			// We need to define an EOS callback on the sink for when we receive an EOS
			// upstream. This gives us an opportunity to cleanup and then signal the pipeline
			// that we are ready to be shut down.
			EOSFunc: func(sink *app.Sink) {
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
				pipeline.GetPipelineBus().Post(gst.NewEOSMessage(appSink))
			},
			NewSampleFunc: func(sink *app.Sink) gst.FlowReturn {
				// Increment the frame number counter
				frameNum++

				if frameNum > totalFrames {
					// If we've reached the total number of frames we are expecting. We can
					// signal the main loop to quit.
					// This needs to be done from a goroutine to not block the app sink
					// callback.
					return gst.FlowEOS
				}

				// Pull the sample from the sink
				sample := sink.PullSample()
				if sample == nil {
					return gst.FlowOK
				}
				defer sample.Unref()

				fmt.Printf("\033[2K\r")
				fmt.Printf("Processing image frame %d/%d", frameNum, totalFrames)

				// Retrieve the buffer from the sample.
				buffer := sample.GetBuffer()

				// We can get an io.Reader directly from the buffer.
				img, err := jpeg.Decode(buffer.Reader())
				if err != nil {
					self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, "Error decoding jpeg frame", err.Error())
					return gst.FlowError
				}

				// Create a new paletted image with the same bounds as the pulled one
				frame := image.NewPaletted(img.Bounds(), video.FormatRGB8P.Palette())

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
			},
		})

		// Link the src pad to the queue
		srcPad.Link(queue.GetStaticPad("sink"))
	})

	fmt.Println("Encoding video to gif")

	// Now that the pipeline is all set up we can start it.
	pipeline.SetState(gst.StatePlaying)

	// Add a watch on the bus on the pipeline and catch any errors
	// that happen.
	var pipelineErr error
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS:
			mainLoop.Quit()
		case gst.MessageError:
			gerr := msg.ParseError()
			fmt.Println("ERROR:", gerr.Error())
			if debug := gerr.DebugString(); debug != "" {
				fmt.Println("DEBUG")
				fmt.Println(debug)
			}
			mainLoop.Quit()
			pipelineErr = gerr
			return false
		}

		return true
	})

	// Iterate on the main loop until the pipeline is finished.
	mainLoop.Run()

	return pipelineErr
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

	examples.RunLoop(encodeGif)
}
