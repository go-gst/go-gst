package gst

/*
#include "gst.go.h"
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

func init() {
	tm := []glib.TypeMarshaler{
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
		{
			T: glib.Type(C.gst_buffer_get_type()),
			F: marshalBuffer,
		},
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
			T: glib.Type(C.gst_proxy_pad_get_type()),
			F: marshalProxyPad,
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
			T: glib.Type(C.bufferListType()),
			F: marshalBufferList,
		},
		{
			T: glib.Type(C.getCapsType()),
			F: marshalCaps,
		},
		{
			T: glib.Type(C.GST_TYPE_CAPS_FEATURES),
			F: marshalCapsFeatures,
		},
		{
			T: glib.Type(C.GST_TYPE_CONTEXT),
			F: marshalContext,
		},
		{
			T: glib.Type(C.GST_TYPE_TOC_ENTRY),
			F: marshalTOCEntry,
		},
		{
			T: glib.Type(C.GST_TYPE_TOC),
			F: marshalTOC,
		},
		{
			T: glib.Type(C.GST_TYPE_TAG_LIST),
			F: marsalTagList,
		},
		{
			T: glib.Type(C.GST_TYPE_EVENT),
			F: marshalEvent,
		},
		{
			T: glib.Type(C.GST_TYPE_SEGMENT),
			F: marshalSegment,
		},
		{
			T: glib.Type(C.GST_TYPE_QUERY),
			F: marshalQuery,
		},
		{
			T: glib.Type(C.gst_message_get_type()),
			F: marshalMessage,
		},
	}

	glib.RegisterGValueMarshalers(tm)
}

func uintptrToGVal(p uintptr) *C.GValue {
	return (*C.GValue)(unsafe.Pointer(p)) // vet thinks this is unsafe and there is no way around it for now.
	// but the given ptr is an address to a C object so go's concerns are misplaced.
}

// Object wrappers

func wrapAllocator(obj *glib.Object) *Allocator           { return &Allocator{wrapObject(obj)} }
func wrapBin(obj *glib.Object) *Bin                       { return &Bin{wrapElement(obj)} }
func wrapBuffer(buf *C.GstBuffer) *Buffer                 { return &Buffer{ptr: buf} }
func wrapBufferList(bufList *C.GstBufferList) *BufferList { return &BufferList{ptr: bufList} }
func wrapBufferPool(obj *glib.Object) *BufferPool         { return &BufferPool{wrapObject(obj)} }
func wrapBus(obj *glib.Object) *Bus                       { return &Bus{Object: wrapObject(obj)} }
func wrapCaps(caps *C.GstCaps) *Caps                      { return &Caps{native: caps} }
func wrapClock(obj *glib.Object) *Clock                   { return &Clock{wrapObject(obj)} }
func wrapContext(ctx *C.GstContext) *Context              { return &Context{ptr: ctx} }
func wrapDevice(obj *glib.Object) *Device                 { return &Device{wrapObject(obj)} }
func wrapElement(obj *glib.Object) *Element               { return &Element{wrapObject(obj)} }
func wrapEvent(ev *C.GstEvent) *Event                     { return &Event{ptr: ev} }
func wrapGhostPad(obj *glib.Object) *GhostPad             { return &GhostPad{wrapProxyPad(obj)} }
func wrapMainContext(ctx *C.GMainContext) *MainContext    { return &MainContext{ptr: ctx} }
func wrapMainLoop(loop *C.GMainLoop) *MainLoop            { return &MainLoop{ptr: loop} }
func wrapMapInfo(mapInfo *C.GstMapInfo) *MapInfo          { return &MapInfo{ptr: mapInfo} }
func wrapMemory(mem *C.GstMemory) *Memory                 { return &Memory{ptr: mem} }
func wrapMessage(msg *C.GstMessage) *Message              { return &Message{msg: msg} }
func wrapMeta(meta *C.GstMeta) *Meta                      { return &Meta{ptr: meta} }
func wrapMetaInfo(info *C.GstMetaInfo) *MetaInfo          { return &MetaInfo{ptr: info} }
func wrapPad(obj *glib.Object) *Pad                       { return &Pad{wrapObject(obj)} }
func wrapPadTemplate(obj *glib.Object) *PadTemplate       { return &PadTemplate{wrapObject(obj)} }
func wrapPipeline(obj *glib.Object) *Pipeline             { return &Pipeline{Bin: wrapBin(obj)} }
func wrapPluginFeature(obj *glib.Object) *PluginFeature   { return &PluginFeature{wrapObject(obj)} }
func wrapPlugin(obj *glib.Object) *Plugin                 { return &Plugin{wrapObject(obj)} }
func wrapProxyPad(obj *glib.Object) *ProxyPad             { return &ProxyPad{wrapPad(obj)} }
func wrapQuery(query *C.GstQuery) *Query                  { return &Query{ptr: query} }
func wrapRegistry(obj *glib.Object) *Registry             { return &Registry{wrapObject(obj)} }
func wrapSample(sample *C.GstSample) *Sample              { return &Sample{sample: sample} }
func wrapSegment(segment *C.GstSegment) *Segment          { return &Segment{ptr: segment} }
func wrapStream(obj *glib.Object) *Stream                 { return &Stream{wrapObject(obj)} }
func wrapTagList(tagList *C.GstTagList) *TagList          { return &TagList{ptr: tagList} }
func wrapTOC(toc *C.GstToc) *TOC                          { return &TOC{ptr: toc} }
func wrapTOCEntry(toc *C.GstTocEntry) *TOCEntry           { return &TOCEntry{ptr: toc} }

func wrapCapsFeatures(features *C.GstCapsFeatures) *CapsFeatures {
	return &CapsFeatures{native: features}
}

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
	if n == ClockTimeNone {
		return time.Duration(-1)
	}
	return time.Duration(uint64(n)) * time.Nanosecond
}

func guint64ToDuration(n C.guint64) time.Duration { return clockTimeToDuration(ClockTime(n)) }

func durationToClockTime(dur time.Duration) ClockTime {
	if dur.Nanoseconds() < 0 {
		return ClockTimeNone
	}
	return ClockTime(dur.Nanoseconds())
}

// Marshallers

func marshalBufferingMode(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return BufferingMode(c), nil
}

func marshalFormat(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return Format(c), nil
}

func marshalMessageType(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return MessageType(c), nil
}

func marshalPadLinkReturn(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return PadLinkReturn(c), nil
}

func marshalState(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return State(c), nil
}

func marshalSeekFlags(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return SeekFlags(c), nil
}

func marshalSeekType(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return SeekType(c), nil
}

func marshalStateChangeReturn(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(uintptrToGVal(p))
	return StateChangeReturn(c), nil
}

func marshalGhostPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapGhostPad(obj), nil
}

func marshalProxyPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapProxyPad(obj), nil
}

func marshalPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPad(obj), nil
}

func marshalMessage(p uintptr) (interface{}, error) {
	c := C.g_value_get_boxed(uintptrToGVal(p))
	return &Message{(*C.GstMessage)(unsafe.Pointer(c))}, nil
}

func marshalObject(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapObject(obj), nil
}

func marshalBus(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapBus(obj), nil
}

func marshalElementFactory(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapElementFactory(obj), nil
}

func marshalPipeline(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPipeline(obj), nil
}

func marshalPluginFeature(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPluginFeature(obj), nil
}

func marshalElement(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapElement(obj), nil
}

func marshalBin(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapBin(obj), nil
}

func marshalAllocationParams(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstAllocationParams)(unsafe.Pointer(c))
	return wrapAllocationParams(obj), nil
}

func marshalMemory(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstMemory)(unsafe.Pointer(c))
	return wrapMemory(obj), nil
}

func marshalBuffer(p uintptr) (interface{}, error) {
	c := C.getBufferValue(uintptrToGVal(p))
	return wrapBuffer(c), nil
}

func marshalBufferList(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstBufferList)(unsafe.Pointer(c))
	return wrapBufferList(obj), nil
}

func marshalCaps(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstCaps)(unsafe.Pointer(c))
	return wrapCaps(obj), nil
}

func marshalCapsFeatures(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstCapsFeatures)(unsafe.Pointer(c))
	return wrapCapsFeatures(obj), nil
}

func marshalContext(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstContext)(unsafe.Pointer(c))
	return wrapContext(obj), nil
}

func marshalTOC(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstToc)(unsafe.Pointer(c))
	return wrapTOC(obj), nil
}

func marshalTOCEntry(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstTocEntry)(unsafe.Pointer(c))
	return wrapTOCEntry(obj), nil
}

func marsalTagList(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstTagList)(unsafe.Pointer(c))
	return wrapTagList(obj), nil
}

func marshalEvent(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstEvent)(unsafe.Pointer(c))
	return wrapEvent(obj), nil
}

func marshalSegment(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstSegment)(unsafe.Pointer(c))
	return wrapSegment(obj), nil
}

func marshalQuery(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(uintptrToGVal(p))
	obj := (*C.GstQuery)(unsafe.Pointer(c))
	return wrapQuery(obj), nil
}
