package main

import (
	"fmt"

	"github.com/tinyzimmer/go-gst/gst"
)

func capsWeirdness() {
	gst.Init(nil)

	caps := gst.NewCapsFromString("audio/x-raw")

	// caps.ForEach(func(features *gst.CapsFeatures, structure *gst.Structure) bool {
	// 	fmt.Println(features)
	// 	return true
	// })

	caps.FilterAndMapInPlace(func(features *gst.CapsFeatures, structure *gst.Structure) bool {
		fmt.Println(features)
		return true
	})

	caps.Unref()
}
