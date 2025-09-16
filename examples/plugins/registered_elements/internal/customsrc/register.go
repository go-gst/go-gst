package customsrc

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/go-gst/go-gst/pkg/gst"
)

// Register needs to be called after gst.Init() to make the gocustomsrc available in the standard
// gst element registry. After this call the element can be used like any other gstreamer element
func Register() bool {
	registered := glib.RegisterSubclassWithConstructor[*customSrc](
		func() *customSrc {
			return &customSrc{}
		},
		glib.WithOverrides[*customSrc, gst.BinOverrides](func(b *customSrc) gst.BinOverrides {
			return gst.BinOverrides{}
		}),
		glib.WithClassInit[*gst.BinClass](func(bc *gst.BinClass) {
			bc.ParentClass().SetStaticMetadata(
				"custom test source",
				"Src/Test",
				"Demo source bin with volume",
				"Wilhelm Bartel <bartel.wilhelm@gmail.com>",
			)

			bc.ParentClass().AddPadTemplate(gst.NewPadTemplate(
				"src",
				gst.PadSrc,
				gst.PadAlways,
				gst.CapsFromString("audio/x-raw,channels=2,rate=48000"),
			))
		}),
	)

	return gst.ElementRegister(
		// no plugin:
		nil,
		// The name of the element
		"gocustomsrc",
		// The rank of the element
		uint(gst.RankNone),
		registered.Type(),
	)
}
