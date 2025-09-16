package custombin

import (
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

func classInit(class *gst.BinClass) {
	class.ParentClass().SetStaticMetadata(
		"custom test source",
		"Src/Test",
		"Demo source bin with volume",
		"Wilhelm Bartel <bartel.wilhelm@gmail.com>",
	)

	class.ParentClass().AddPadTemplate(gst.NewPadTemplate(
		"src",
		gst.PadSrc,
		gst.PadAlways,
		gst.CapsFromString("audio/x-raw,channels=2,rate=48000"),
	))
}

type customBin struct {
	gst.BinInstance // parent object must be first embedded field
	source1         gst.Element
	source2         gst.Element
	mixer           gst.Element
}

// constructed is the method we use to override the GOBject.constructed method.
func (bin *customBin) constructed() {
	bin.source1 = gst.ElementFactoryMakeWithProperties("gocustomsrc", map[string]any{
		"duration": int64(5 * time.Second),
	})

	bin.source2 = gst.ElementFactoryMakeWithProperties("gocustomsrc", map[string]any{
		"duration": int64(10 * time.Second),
	})

	bin.mixer = gst.ElementFactoryMake("audiomixer", "")

	bin.AddMany(
		bin.source1,
		bin.source2,
		bin.mixer,
	)

	srcpad := bin.mixer.GetStaticPad("src")

	ghostpad := gst.NewGhostPadFromTemplate("src", srcpad, bin.GetPadTemplate("src"))

	bin.source1.Link(bin.mixer)
	bin.source2.Link(bin.mixer)

	bin.AddPad(ghostpad)
}
