// This example demonstrates using gstreamer to convert a video stream into image frames
// and then encoding those frames to a gif.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

var srcFile string
var outFile string

func encodeGif(mainLoop *gst.MainLoop) error {
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
			err := gst.NewGError(2, err)
			pipeline.GetPipelineBus().
				Post(gst.NewErrorMessage(self, err, "", nil))
		}

		// Add the elements to the pipeline and sync their state with the pipeline
		pipeline.AddMany(elements...)
		for _, e := range elements {
			e.SyncStateWithParent()
		}

		// Retrieve direct references to some of the elements.
		queue := elements[0]
		videorate := elements[len(elements)-2]
		jpegenc := elements[len(elements)-1]

		// Link all elements up until the videorate. We are going to apply caps there and use
		// a filtered link.
		gst.ElementLinkMany(elements[:len(elements)-1]...)

		// We are going to filter images out all the way down to 5 frames per second
		rateCaps := gst.NewCapsFromString("video/x-raw, framerate=5/1")
		videorate.LinkFiltered(jpegenc, rateCaps)

		// Create an app sink that we are going to use to pull images from the pipeline
		// one at a time. (An error can happen here too, but for the sake of brevity...)
		appSink, _ := app.NewAppSink()
		pipeline.Add(appSink.Element)
		jpegenc.Link(appSink.Element)

		// Getting data out of the sink is done by setting callbacks. Each new sample
		// will be a new jpeg image from the pipeline.
		var frameNum int
		appSink.SetCallbacks(&app.SinkCallbacks{
			NewSampleFunc: func(sink *app.Sink) gst.FlowReturn {
				// We can retrieve a reader with the raw bytes of the image directly from the
				// sink.
				imgReader := sink.PullSample().GetBuffer().Reader()

				img, err := jpeg.Decode(imgReader)
				if err != nil {
					fmt.Println("Error decoding jpeg frame:", err)
					return gst.FlowError
				}

				frameNum++
				fmt.Printf("\033[2K\rProcessing image frame %d", frameNum)
				// fmt.Printf("Processing image frame %d\n", frameNum)

				// Create a new paletted image with the same bounds as the pulled one
				frame := image.NewPaletted(img.Bounds(), palette.Plan9)

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

	// Add a watch on the bus on the pipeline and wait for any errors
	// or the end of the stream.
	var isError bool
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		// Uncomment this for very debuggy output
		// default:
		// 	fmt.Println(msg)
		case gst.MessageEOS:
			mainLoop.Quit()
			return false
		case gst.MessageError:
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			mainLoop.Quit()
			isError = true
			return false
		}

		return true
	})

	fmt.Println("Encoding video to gif")

	// Now that the pipeline is all set up we can start it.
	pipeline.SetState(gst.StatePlaying)

	// Iterate on the main loop until the pipeline is finished.
	mainLoop.Run()

	fmt.Println()

	// If no error happened on the pipeline. Write the results of the gif
	// to the destination.
	if !isError {
		fmt.Println("Writing the results of the gif to", outFile)
		file, err := os.Create(outFile)
		if err != nil {
			return err
		}
		defer file.Close()
		if err := gif.EncodeAll(file, outGif); err != nil {
			return err
		}
	}

	return nil
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
	}

	examples.RunLoop(func(mainLoop *gst.MainLoop) error {
		return encodeGif(mainLoop)
	})
}
