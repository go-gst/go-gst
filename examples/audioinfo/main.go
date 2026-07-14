package main

import (
	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstaudio"
)

func main() {
	// this example is mostly a test to check fixed size array conversions
	gst.Init()

	audioInfo := gstaudio.NewAudioInfo()

	audioInfo.SetFormat(gstaudio.AudioFormatS16le, 44100, 2, [64]gstaudio.AudioChannelPosition{gstaudio.AudioChannelPositionFrontLeft, gstaudio.AudioChannelPositionFrontRight})

	caps := audioInfo.ToCaps()
	println(caps.String())
}
