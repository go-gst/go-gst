package customsrc

import (
	"fmt"
	"math"
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/examples/plugins/registered_elements/internal/common"
	"github.com/go-gst/go-gst/gst"
)

// default: 1024, this value makes it easier to calculate num buffers with the sample rate
const samplesperbuffer = 4800

const samplerate = 48000

var properties = []*glib.ParamSpec{
	glib.NewInt64Param(
		"duration",
		"duration",
		"duration the source",
		0,
		math.MaxInt64,
		0,
		glib.ParameterReadWrite,
	),
}

type customSrc struct {
	// self   *gst.Bin
	source *gst.Element
	volume *gst.Element

	duration time.Duration
}

// ClassInit is the place where you define pads and properties
func (*customSrc) ClassInit(klass *glib.ObjectClass) {
	class := gst.ToElementClass(klass)
	class.SetMetadata(
		"custom test source",
		"Src/Test",
		"Demo source bin with volume",
		"Wilhelm Bartel <bartel.wilhelm@gmail.com>",
	)
	class.AddPadTemplate(gst.NewPadTemplate(
		"src",
		gst.PadDirectionSource,
		gst.PadPresenceAlways,
		gst.NewCapsFromString(fmt.Sprintf("audio/x-raw,channels=2,rate=%d", samplerate)),
	))
	class.InstallProperties(properties)
}

// SetProperty gets called for every property. The id is the index in the slice defined above.
func (s *customSrc) SetProperty(self *glib.Object, id uint, value *glib.Value) {
	param := properties[id]

	bin := gst.ToGstBin(self)

	switch param.Name() {
	case "duration":
		state := bin.GetCurrentState()
		if !(state == gst.StateNull || state != gst.StateReady) {
			return
		}

		gv, _ := value.GoValue()

		durI, _ := gv.(int64)

		s.duration = time.Duration(durI)

		s.updateSource()
	}
}

// GetProperty is called to retrieve the value of the property at index `id` in the properties
// slice provided at ClassInit.
func (o *customSrc) GetProperty(self *glib.Object, id uint) *glib.Value {
	param := properties[id]

	switch param.Name() {
	case "duration":
		v, _ := glib.GValue(int64(o.duration))
		return v
	}

	return nil
}

func (*customSrc) New() glib.GoObjectSubclass {
	return &customSrc{}
}

// InstanceInit should initialize the element. Keep in mind that the properties are not yet present. When this is called.
func (s *customSrc) InstanceInit(instance *glib.Object) {
	self := gst.ToGstBin(instance)

	s.source = common.Must(gst.NewElement("audiotestsrc"))
	s.volume = common.Must(gst.NewElement("volume"))

	klass := instance.Class()
	class := gst.ToElementClass(klass)

	self.AddMany(
		s.source,
		s.volume,
	)

	srcpad := s.volume.GetStaticPad("src")

	ghostpad := gst.NewGhostPadFromTemplate("src", srcpad, class.GetPadTemplate("src"))

	gst.ElementLinkMany(
		s.source,
		s.volume,
	)

	self.AddPad(ghostpad.Pad)

	s.updateSource()
}

func (s *customSrc) Constructed(o *glib.Object) {}

func (s *customSrc) Finalize(o *glib.Object) {
	common.FinalizersCalled++
}

// updateSource will get called to update the audiotestsrc when a property changes
func (s *customSrc) updateSource() {
	if s.source != nil {
		numBuffers := (float64(s.duration / time.Second)) / (float64(samplesperbuffer) / float64(samplerate))

		s.source.SetProperty("num-buffers", int(math.Ceil(numBuffers)))
	}
}
