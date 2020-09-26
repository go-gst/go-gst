package main

import (
	"errors"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tinyzimmer/go-gst-launch/gst"

	"github.com/spf13/cobra"
)

var framesPerSecond int
var imageFormat string

func init() {
	gifCmd.PersistentFlags().IntVarP(&framesPerSecond, "frame-rate", "r", 10, "The number of frames per-second to encode into the GIF")
	gifCmd.PersistentFlags().StringVarP(&imageFormat, "format", "f", "png", "The image format to encode frames to")

	rootCmd.AddCommand(gifCmd)
}

var gifCmd = &cobra.Command{
	Use:   "gif",
	Short: "Encodes the given video to GIF format",
	Long: `Look at the available options to change the compression levels and format.
	
Requires libav be installed.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if srcFile == "" && !fromStdin {
			return errors.New("No input provided")
		}
		if destFile == "" && !toStdout {
			return errors.New("No output provided")
		}
		return nil
	},
	RunE: gifEncode,
}

func gifEncode(cmd *cobra.Command, args []string) error {

	var imageEncoder string
	var decoder func(io.Reader) (image.Image, error)
	switch strings.ToLower(imageFormat) {
	case "png":
		imageEncoder = "pngenc"
		decoder = png.Decode
	case "jpg":
		imageEncoder = "jpegenc"
		decoder = jpeg.Decode
	case "jpeg":
		imageEncoder = "jpegenc"
		decoder = jpeg.Decode
	default:
		return fmt.Errorf("Invalid image format %s: Valid options [ png | jpg ]", strings.ToLower(imageFormat))
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	launchStr := fmt.Sprintf(
		`filesrc location=%s ! decodebin ! videoconvert ! videoscale ! videorate ! video/x-raw,framerate=%d/1 ! %s ! multifilesink location="%s/%%04d.%s"`,
		srcFile, framesPerSecond, imageEncoder, tmpDir, strings.ToLower(imageFormat),
	)

	logInfo("gif", "Converting video to image frames")

	gstPipeline, err := gst.NewPipelineFromString(launchStr)
	if err != nil {
		return err
	}
	defer gstPipeline.Destroy()

	if err := gstPipeline.Start(); err != nil {
		return err
	}

	if verbose {
		setupVerbosePipelineListeners(gstPipeline, "gif")
	}

	gst.Wait(gstPipeline)

	logInfo("gif", "Building output gif:", destFile)

	dest, err := getDestFile()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return err
	}

	outGif := &gif.GIF{
		Image: make([]*image.Paletted, 0),
		Delay: make([]int, 0),
	}

	numFrames := len(files)
	for idx, file := range files {
		if !toStdout {
			fmt.Printf("\rEncoding frame %d/%d", idx, numFrames)
		}
		f, err := os.Open(filepath.Join(tmpDir, file.Name()))
		if err != nil {
			return err
		}
		img, err := decoder(f)
		if err != nil {
			return err
		}
		frame := image.NewPaletted(img.Bounds(), palette.Plan9)
		for x := 1; x <= img.Bounds().Dx(); x++ {
			for y := 1; y <= img.Bounds().Dy(); y++ {
				frame.Set(x, y, img.At(x, y))
			}
		}
		outGif.Image = append(outGif.Image, frame)
		outGif.Delay = append(outGif.Delay, 0)
	}

	if err := gif.EncodeAll(dest, outGif); err != nil {
		return err
	}

	if !toStdout {
		fmt.Println()
	}

	logInfo("gif", "If the command reached this state and you see a GStreamer-CRITICAL error, you can ignore it")
	return nil
}
