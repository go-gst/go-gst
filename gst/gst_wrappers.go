package gst

// #include "gst.go.h"
import "C"

import (
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

func init() {
	tm := []glib.TypeMarshaler{
		// Enums
		{
			T: glib.Type(C.gst_buffering_mode_get_type()),
			F: marshalBufferingMode,
		},
		{
			T: glib.Type(C.gst_format_get_type()),
			F: marshalFormat,
		},
		{
			T: glib.Type(C.gst_message_type_get_type()),
			F: marshalMessageType,
		},
		{
			T: glib.Type(C.gst_pad_link_return_get_type()),
			F: marshalPadLinkReturn,
		},
		{
			T: glib.Type(C.gst_state_get_type()),
			F: marshalState,
		},
		{
			T: glib.Type(C.gst_seek_flags_get_type()),
			F: marshalSeekFlags,
		},
		{
			T: glib.Type(C.gst_seek_type_get_type()),
			F: marshalSeekType,
		},
		{
			T: glib.Type(C.gst_state_change_return_get_type()),
			F: marshalStateChangeReturn,
		},

		// Objects/Interfaces
		{
			T: glib.Type(C.gst_pipeline_get_type()),
			F: marshalPipeline,
		},
		{
			T: glib.Type(C.gst_bin_get_type()),
			F: marshalBin,
		},
		{
			T: glib.Type(C.gst_bus_get_type()),
			F: marshalBus,
		},
		{
			T: glib.Type(C.gst_element_get_type()),
			F: marshalElement,
		},
		{
			T: glib.Type(C.gst_element_factory_get_type()),
			F: marshalElementFactory,
		},
		{
			T: glib.Type(C.gst_ghost_pad_get_type()),
			F: marshalGhostPad,
		},
		{
			T: glib.Type(C.gst_object_get_type()),
			F: marshalObject,
		},
		{
			T: glib.Type(C.gst_pad_get_type()),
			F: marshalPad,
		},
		{
			T: glib.Type(C.gst_plugin_feature_get_type()),
			F: marshalPluginFeature,
		},
		{
			T: glib.Type(C.gst_allocation_params_get_type()),
			F: marshalAllocationParams,
		},
		{
			T: glib.Type(C.GST_TYPE_MEMORY),
			F: marshalMemory,
		},
		{
			T: glib.Type(C.gst_atomic_queue_get_type()),
			F: marshalAtomicQueue,
		},

		// Boxed
		{T: glib.Type(C.gst_message_get_type()), F: marshalMessage},
	}
	glib.RegisterGValueMarshalers(tm)
}

// Object wrappers

func wrapAllocator(obj *glib.Object) *Allocator            { return &Allocator{wrapObject(obj)} }
func wrapAtomicQueue(queue *C.GstAtomicQueue) *AtomicQueue { return &AtomicQueue{ptr: queue} }
func wrapBin(obj *glib.Object) *Bin                        { return &Bin{wrapElement(obj)} }
func wrapBuffer(buf *C.GstBuffer) *Buffer                  { return &Buffer{ptr: buf} }
func wrapBus(obj *glib.Object) *Bus                        { return &Bus{Object: wrapObject(obj)} }
func wrapClock(obj *glib.Object) *Clock                    { return &Clock{wrapObject(obj)} }
func wrapDevice(obj *glib.Object) *Device                  { return &Device{wrapObject(obj)} }
func wrapElement(obj *glib.Object) *Element                { return &Element{wrapObject(obj)} }
func wrapGhostPad(obj *glib.Object) *GhostPad              { return &GhostPad{wrapPad(obj)} }
func wrapMainContext(ctx *C.GMainContext) *MainContext     { return &MainContext{ptr: ctx} }
func wrapMainLoop(loop *C.GMainLoop) *MainLoop             { return &MainLoop{ptr: loop} }
func wrapMemory(mem *C.GstMemory) *Memory                  { return &Memory{ptr: mem} }
func wrapMessage(msg *C.GstMessage) *Message               { return &Message{msg: msg} }
func wrapMeta(meta *C.GstMeta) *Meta                       { return &Meta{ptr: meta} }
func wrapMetaInfo(info *C.GstMetaInfo) *MetaInfo           { return &MetaInfo{ptr: info} }
func wrapPad(obj *glib.Object) *Pad                        { return &Pad{wrapObject(obj)} }
func wrapPadTemplate(obj *glib.Object) *PadTemplate        { return &PadTemplate{wrapObject(obj)} }
func wrapPipeline(obj *glib.Object) *Pipeline              { return &Pipeline{Bin: wrapBin(obj)} }
func wrapPluginFeature(obj *glib.Object) *PluginFeature    { return &PluginFeature{wrapObject(obj)} }
func wrapPlugin(obj *glib.Object) *Plugin                  { return &Plugin{wrapObject(obj)} }
func wrapRegistry(obj *glib.Object) *Registry              { return &Registry{wrapObject(obj)} }
func wrapSample(sample *C.GstSample) *Sample               { return &Sample{sample: sample} }
func wrapStream(obj *glib.Object) *Stream                  { return &Stream{wrapObject(obj)} }
func wrapTagList(tagList *C.GstTagList) *TagList           { return &TagList{ptr: tagList} }

func wrapObject(obj *glib.Object) *Object {
	return &Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: obj}}
}

func wrapElementFactory(obj *glib.Object) *ElementFactory {
	return &ElementFactory{wrapPluginFeature(obj)}
}

func wrapStreamCollection(obj *glib.Object) *StreamCollection {
	return &StreamCollection{wrapObject(obj)}
}

func wrapAllocationParams(obj *C.GstAllocationParams) *AllocationParams {
	return &AllocationParams{ptr: obj}
}

// Clock wrappers

func clockTimeToDuration(n ClockTime) time.Duration {
	return time.Duration(uint64(n)) * time.Nanosecond
}

func guint64ToDuration(n C.guint64) time.Duration { return clockTimeToDuration(ClockTime(n)) }

// Enums/Constants

func marshalBufferingMode(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return BufferingMode(c), nil
}

func marshalFormat(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return Format(c), nil
}

func marshalMessageType(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return MessageType(c), nil
}

func marshalPadLinkReturn(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return PadLinkReturn(c), nil
}

func marshalState(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return State(c), nil
}

func marshalSeekFlags(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return SeekFlags(c), nil
}

func marshalSeekType(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return SeekType(c), nil
}

func marshalStateChangeReturn(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))
	return StateChangeReturn(c), nil
}

func marshalGhostPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapGhostPad(obj), nil
}

func marshalPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPad(obj), nil
}

func marshalMessage(p uintptr) (interface{}, error) {
	c := C.g_value_get_boxed((*C.GValue)(unsafe.Pointer(p)))
	return &Message{(*C.GstMessage)(unsafe.Pointer(c))}, nil
}

func marshalObject(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapObject(obj), nil
}

func marshalBus(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapBus(obj), nil
}

func marshalElementFactory(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapElementFactory(obj), nil
}

func marshalPipeline(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPipeline(obj), nil
}

func marshalPluginFeature(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPluginFeature(obj), nil
}

func marshalElement(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapElement(obj), nil
}

func marshalBin(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapBin(obj), nil
}

func marshalAllocationParams(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := (*C.GstAllocationParams)(unsafe.Pointer(c))
	return wrapAllocationParams(obj), nil
}

func marshalMemory(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := (*C.GstMemory)(unsafe.Pointer(c))
	return wrapMemory(obj), nil
}

func marshalAtomicQueue(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := (*C.GstAtomicQueue)(unsafe.Pointer(c))
	return wrapAtomicQueue(obj), nil
}
