package customtransform

import (
	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/base"
)

type customBaseTransform struct{}

// ClassInit is the place where you define pads and properties
func (*customBaseTransform) ClassInit(klass *glib.ObjectClass) {
	class := gst.ToElementClass(klass)
	class.SetMetadata(
		"custom base transform",
		"Transform/demo",
		"custom base transform",
		"Wilhelm Bartel <bartel.wilhelm@gmail.com>",
	)
	class.AddPadTemplate(gst.NewPadTemplate(
		"src",
		gst.PadDirectionSource,
		gst.PadPresenceAlways,
		gst.NewCapsFromString("audio/x-raw,channels=2,rate=48000"),
	))
	class.AddPadTemplate(gst.NewPadTemplate(
		"sink",
		gst.PadDirectionSink,
		gst.PadPresenceAlways,
		gst.NewCapsFromString("audio/x-raw,channels=2,rate=48000"),
	))
}

// SetProperty gets called for every property. The id is the index in the slice defined above.
func (s *customBaseTransform) SetProperty(self *glib.Object, id uint, value *glib.Value) {}

// GetProperty is called to retrieve the value of the property at index `id` in the properties
// slice provided at ClassInit.
func (o *customBaseTransform) GetProperty(self *glib.Object, id uint) *glib.Value {
	return nil
}

// New is called by the bindings to create a new instance of your go element. Use this to initialize channels, maps, etc.
//
// Think of New like the constructor of your struct
func (*customBaseTransform) New() glib.GoObjectSubclass {
	return &customBaseTransform{}
}

// InstanceInit should initialize the element. Keep in mind that the properties are not yet present. When this is called.
func (s *customBaseTransform) InstanceInit(instance *glib.Object) {}

func (s *customBaseTransform) Constructed(o *glib.Object) {}

func (s *customBaseTransform) Finalize(o *glib.Object) {}

// see base.GstBaseTransformImpl interface for the method signatures of the virtual methods
//
// it is not required to implement all methods
var _ base.GstBaseTransformImpl = nil

func (s *customBaseTransform) SinkEvent(self *base.GstBaseTransform, event *gst.Event) bool {
	return self.ParentSinkEvent(event)
}

func (s *customBaseTransform) SrcEvent(self *base.GstBaseTransform, event *gst.Event) bool {
	return self.ParentSrcEvent(event)
}
