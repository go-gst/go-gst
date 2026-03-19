package gstvideo

import (
	"image/color"
	"unsafe"
)

// #cgo pkg-config: gstreamer-video-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/video/video.h>
import "C"

// VideoFormatGetPalette wraps gst_video_format_get_palette
func VideoFormatGetPalette(format VideoFormat) []color.Color {
	var size C.gsize
	ptr := C.gst_video_format_get_palette(C.GstVideoFormat(format), &size)

	paletteBytes := unsafe.Slice((*byte)(ptr), size)

	// Convert the byte slice to a slice of color.Color
	return bytesToColorPalette(paletteBytes)
}

// bytesToColorPalette converts a byte slice into a slice of color.Color.
// Each color is represented by 4 bytes (RGBA).
func bytesToColorPalette(paletteBytes []byte) []color.Color {
	const bytesPerColor = 4
	numColors := len(paletteBytes) / bytesPerColor
	palette := make([]color.Color, numColors)

	for i := 0; i < numColors; i++ {
		offset := i * bytesPerColor
		palette[i] = color.RGBA{
			R: paletteBytes[offset],
			G: paletteBytes[offset+1],
			B: paletteBytes[offset+2],
			A: paletteBytes[offset+3],
		}
	}

	return palette
}
