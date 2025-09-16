package customsrc

import (
	"math"
	"time"

	"github.com/diamondburned/gotk4/pkg/gobject/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

// default: 1024, this value makes it easier to calculate num buffers with the sample rate
const samplesperbuffer = 4800

const samplerate = 48000

type customSrc struct {
	gst.BinInstance // parent must be embedded as the first field

	source gst.Element
	volume gst.Element

	Duration time.Duration
}

// InstanceInit should initialize the element. Keep in mind that the properties are not yet present. When this is called.
func (bin *customSrc) init() {
	bin.source = gst.ElementFactoryMake("audiotestsrc", "")
	bin.volume = gst.ElementFactoryMake("volume", "")

	bin.AddMany(
		bin.source,
		bin.volume,
	)

	srcpad := bin.volume.GetStaticPad("src")

	ghostpad := gst.NewGhostPadFromTemplate("src", srcpad, bin.GetPadTemplate("src"))

	gst.LinkMany(
		bin.source,
		bin.volume,
	)

	bin.AddPad(ghostpad)

	bin.updateSource()
}

func (bin *customSrc) setProperty(_ uint, value any, pspec *gobject.ParamSpec) {
	switch pspec.Name() {
	case "duration":
		bin.Duration = value.(time.Duration)
		bin.updateSource()
	default:
		panic("unknown property")
	}
}

func (bin *customSrc) getProperty(_ uint, pspec *gobject.ParamSpec) any {
	switch pspec.Name() {
	case "duration":
		return bin.Duration
	default:
		panic("unknown property")
	}
}

// updateSource will get called to update the audiotestsrc when a property changes
func (s *customSrc) updateSource() {
	if s.source != nil {
		numBuffers := (float64(s.Duration / time.Second)) / (float64(samplesperbuffer) / float64(samplerate))

		s.source.SetObjectProperty("num-buffers", int(math.Ceil(numBuffers)))
	}
}
