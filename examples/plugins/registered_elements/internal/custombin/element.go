package custombin

import (
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

type customBin struct {
	gst.Bin // parent object must be first embedded field
	source1 gst.Elementer
	source2 gst.Elementer
	mixer   gst.Elementer
}

// init should initialize the element. Keep in mind that the properties are not yet present. When this is called.
func (bin *customBin) init() {
	bin.source1 = gst.ElementFactoryMakeWithProperties("gocustomsrc", map[string]interface{}{
		"duration": int64(5 * time.Second),
	})

	bin.source2 = gst.ElementFactoryMakeWithProperties("gocustomsrc", map[string]interface{}{
		"duration": int64(10 * time.Second),
	})

	bin.mixer = gst.ElementFactoryMake("audiomixer", "")

	bin.AddMany(
		bin.source1,
		bin.source2,
		bin.mixer,
	)

	srcpad := bin.mixer.StaticPad("src")

	ghostpad := gst.NewGhostPadFromTemplate("src", srcpad, bin.PadTemplate("src"))

	bin.source1.Link(bin.mixer)
	bin.source2.Link(bin.mixer)

	bin.AddPad(&ghostpad.Pad)
}
