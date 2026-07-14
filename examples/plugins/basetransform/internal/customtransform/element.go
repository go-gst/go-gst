package customtransform

import (
	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstbase"
)

type customBaseTransform struct {
	gstbase.BaseTransformInstance
}

func classInit(class *gstbase.BaseTransformClass) {
	class.ParentClass().SetStaticMetadata(
		"custom base transform",
		"Transform/demo",
		"custom base transform",
		"Wilhelm Bartel <bartel.wilhelm@gmail.com>",
	)

	class.ParentClass().AddPadTemplate(gst.NewPadTemplate(
		"src",
		gst.PadSrc,
		gst.PadAlways,
		gst.CapsFromString("audio/x-raw,channels=2,rate=48000"),
	))
	class.ParentClass().AddPadTemplate(gst.NewPadTemplate(
		"sink",
		gst.PadSink,
		gst.PadAlways,
		gst.CapsFromString("audio/x-raw,channels=2,rate=48000"),
	))
}
