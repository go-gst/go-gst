package customtransform

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstbase"
)

// Register needs to be called after gst.Init() to make the gocustombin available in the standard
// gst element registry. After this call the element can be used like any other gstreamer element
func Register() bool {
	registered := glib.RegisterSubclassWithConstructor[*customBaseTransform](
		func() *customBaseTransform {
			return &customBaseTransform{}
		},
		glib.WithOverrides[*customBaseTransform, gstbase.BaseTransformOverrides](func(b *customBaseTransform) gstbase.BaseTransformOverrides {
			return gstbase.BaseTransformOverrides{}
		}),
		glib.WithClassInit[*gstbase.BaseTransformClass](func(class *gstbase.BaseTransformClass) {
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
		}),
	)

	return gst.ElementRegister(
		// no plugin:
		nil,
		// The name of the element
		"gocustomtransform",
		// The rank of the element
		uint(gst.RankNone),
		registered.Type(),
	)
}
