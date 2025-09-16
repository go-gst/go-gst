package customsrc

import (
	"math"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

// default: 1024, this value makes it easier to calculate num buffers with the sample rate
const samplesperbuffer = 4800

const samplerate = 48000

type customSrc struct {
	gst.Bin // parent must be embedded as the first field

	source gst.Elementer
	volume gst.Elementer

	Duration time.Duration `glib:"duration"`
}

// InstanceInit should initialize the element. Keep in mind that the properties are not yet present. When this is called.
func (bin *customSrc) init() {
	bin.source = gst.ElementFactoryMake("audiotestsrc", "")
	bin.volume = gst.ElementFactoryMake("volume", "")

	bin.AddMany(
		bin.source,
		bin.volume,
	)

	srcpad := bin.volume.StaticPad("src")

	ghostpad := gst.NewGhostPadFromTemplate("src", srcpad, bin.PadTemplate("src"))

	gst.LinkMany(
		bin.source,
		bin.volume,
	)

	bin.AddPad(&ghostpad.Pad)

	bin.updateSource()
}

// updateSource will get called to update the audiotestsrc when a property changes
func (s *customSrc) updateSource() {
	if s.source != nil {
		numBuffers := (float64(s.Duration / time.Second)) / (float64(samplesperbuffer) / float64(samplerate))

		s.source.SetObjectProperty("num-buffers", int(math.Ceil(numBuffers)))
	}
}
