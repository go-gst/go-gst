package custombin

import (
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/examples/plugins/registered_elements/internal/common"
	"github.com/go-gst/go-gst/gst"
)

type customBin struct {
	// self    *gst.Bin
	source1 *gst.Element
	source2 *gst.Element
	mixer   *gst.Element
}

// ClassInit is the place where you define pads and properties
func (*customBin) ClassInit(klass *glib.ObjectClass) {
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
		gst.NewCapsFromString("audio/x-raw,channels=2,rate=48000"),
	))
}

// SetProperty gets called for every property. The id is the index in the slice defined above.
func (s *customBin) SetProperty(self *glib.Object, id uint, value *glib.Value) {}

// GetProperty is called to retrieve the value of the property at index `id` in the properties
// slice provided at ClassInit.
func (o *customBin) GetProperty(self *glib.Object, id uint) *glib.Value {
	return nil
}

// New is called by the bindings to create a new instance of your go element. Use this to initialize channels, maps, etc.
//
// Think of New like the constructor of your struct
func (*customBin) New() glib.GoObjectSubclass {
	return &customBin{}
}

// InstanceInit should initialize the element. Keep in mind that the properties are not yet present. When this is called.
func (s *customBin) InstanceInit(instance *glib.Object) {
	self := gst.ToGstBin(instance)

	s.source1 = common.Must(gst.NewElementWithProperties("gocustomsrc", map[string]interface{}{
		"duration": int64(5 * time.Second),
	}))
	s.source2 = common.Must(gst.NewElementWithProperties("gocustomsrc", map[string]interface{}{
		"duration": int64(10 * time.Second),
	}))

	s.mixer = common.Must(gst.NewElement("audiomixer"))

	klass := instance.Class()
	class := gst.ToElementClass(klass)

	self.AddMany(
		s.source1,
		s.source2,
		s.mixer,
	)

	srcpad := s.mixer.GetStaticPad("src")

	ghostpad := gst.NewGhostPadFromTemplate("src", srcpad, class.GetPadTemplate("src"))

	s.source1.Link(s.mixer)
	s.source2.Link(s.mixer)

	self.AddPad(ghostpad.Pad)
}

func (s *customBin) Constructed(o *glib.Object) {}

func (s *customBin) Finalize(o *glib.Object) {
	common.FinalizersCalled++
}
