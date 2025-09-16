package customsrc

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-gst/go-glib/pkg/gobject/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

// default: 1024, this value makes it easier to calculate num buffers with the sample rate
const samplesperbuffer = 4800

const samplerate = 48000

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
		gst.CapsFromString(fmt.Sprintf("audio/x-raw,channels=2,rate=%d", samplerate)),
	))

	class.ParentClass().ParentClass().ParentClass().ParentClass().InstallProperties([]*gobject.ParamSpec{
		gobject.ParamSpecInt64(
			"duration",
			"duration",
			"Duration of the source in nanoseconds",
			0,
			math.MaxInt64,
			0,
			gobject.ParamReadwrite|gobject.ParamConstruct),
	})
}

type customSrc struct {
	gst.BinInstance // parent must be embedded as the first field

	source gst.Element
	volume gst.Element

	Duration time.Duration
}

// InstanceInit should initialize the element. Keep in mind that the properties are not yet present. When this is called.
func (bin *customSrc) init() {
	bin.source = gst.ElementFactoryMakeWithProperties("audiotestsrc", map[string]any{
		"samplesperbuffer": samplesperbuffer,
	})
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
		bin.Duration = time.Duration(value.(int64)) // declared as int64 property in classInit
		log.Printf("set duration to %s", bin.Duration)
		bin.updateSource()
	default:
		panic("unknown property")
	}
}

func (bin *customSrc) getProperty(_ uint, pspec *gobject.ParamSpec) any {
	switch pspec.Name() {
	case "duration":
		return int64(bin.Duration) // declared as int64 property in classInit
	default:
		panic("unknown property")
	}
}

// updateSource will get called to update the audiotestsrc when a property changes
func (s *customSrc) updateSource() {
	if s.source == nil {
		// the construct param may be set before we initialized the source
		return
	}

	numBuffers := (float64(s.Duration / time.Second)) / (float64(samplesperbuffer) / float64(samplerate))

	s.source.SetObjectProperty("num-buffers", int32(math.Ceil(numBuffers)))

	log.Printf("set num-buffers to %d", int(math.Ceil(numBuffers)))
}
