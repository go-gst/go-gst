package gst

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

func init() { registerMarshalers() }

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
func wrapElement(obj *glib.Object) *Element               { return &Element{wrapObject(obj)} }
func wrapEvent(ev *C.GstEvent) *Event                     { return &Event{ptr: ev} }
func wrapGhostPad(obj *glib.Object) *GhostPad             { return &GhostPad{wrapProxyPad(obj)} }
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

func wrapAllocationParams(obj *C.GstAllocationParams) *AllocationParams {
	return &AllocationParams{ptr: obj}
}

// Marshallers

func registerMarshalers() {
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
			T: glib.Type(C.gst_memory_get_type()),
			F: marshalMemory,
		},
		{
			T: glib.Type(C.gst_buffer_list_get_type()),
			F: marshalBufferList,
		},
		{
			T: TypeCaps,
			F: marshalCaps,
		},
		{
			T: TypeCapsFeatures,
			F: marshalCapsFeatures,
		},
		{
			T: glib.Type(C.gst_context_get_type()),
			F: marshalContext,
		},
		{
			T: glib.Type(C.gst_toc_entry_get_type()),
			F: marshalTOCEntry,
		},
		{
			T: glib.Type(C.gst_toc_get_type()),
			F: marshalTOC,
		},
		{
			T: glib.Type(C.gst_tag_list_get_type()),
			F: marsalTagList,
		},
		{
			T: glib.Type(C.gst_event_get_type()),
			F: marshalEvent,
		},
		{
			T: glib.Type(C.gst_segment_get_type()),
			F: marshalSegment,
		},
		{
			T: glib.Type(C.gst_query_get_type()),
			F: marshalQuery,
		},
		{
			T: glib.Type(C.gst_message_get_type()),
			F: marshalMessage,
		},
		{
			T: TypeBitmask,
			F: marshalBitmask,
		},
		{
			T: TypeFraction,
			F: marshalFraction,
		},
		{
			T: TypeFractionRange,
			F: marshalFractionRange,
		},
		{
			T: TypeStructure,
			F: marshalStructure,
		},
		{
			T: TypeFloat64Range,
			F: marshalDoubleRange,
		},
		{
			T: TypeFlagset,
			F: marshalFlagset,
		},
		{
			T: TypeInt64Range,
			F: marshalInt64Range,
		},
		{
			T: TypeIntRange,
			F: marshalIntRange,
		},
		{
			T: TypeValueArray,
			F: marshalValueArray,
		},
		{
			T: TypeValueList,
			F: marshalValueList,
		},
		{
			T: glib.Type(C.gst_sample_get_type()),
			F: marshalSample,
		},
	}

	glib.RegisterGValueMarshalers(tm)
}

func toGValue(p uintptr) *C.GValue {
	return (*C.GValue)((unsafe.Pointer)(p))
}

func marshalValueArray(p uintptr) (interface{}, error) {
	val := toGValue(p)
	out := ValueArrayValue(*glib.ValueFromNative(unsafe.Pointer(val)))
	return &out, nil
}

func marshalValueList(p uintptr) (interface{}, error) {
	val := glib.ValueFromNative(unsafe.Pointer(toGValue(p)))
	out := ValueListValue(*glib.ValueFromNative(unsafe.Pointer(val)))
	return &out, nil
}

func marshalInt64Range(p uintptr) (interface{}, error) {
	v := toGValue(p)
	return &Int64RangeValue{
		start: int64(C.gst_value_get_int64_range_min(v)),
		end:   int64(C.gst_value_get_int64_range_max(v)),
		step:  int64(C.gst_value_get_int64_range_step(v)),
	}, nil
}

func marshalIntRange(p uintptr) (interface{}, error) {
	v := toGValue(p)
	return &IntRangeValue{
		start: int(C.gst_value_get_int_range_min(v)),
		end:   int(C.gst_value_get_int_range_max(v)),
		step:  int(C.gst_value_get_int_range_step(v)),
	}, nil
}

func marshalBitmask(p uintptr) (interface{}, error) {
	v := toGValue(p)
	return Bitmask(C.gst_value_get_bitmask(v)), nil
}

func marshalFlagset(p uintptr) (interface{}, error) {
	v := toGValue(p)
	return &FlagsetValue{
		flags: uint(C.gst_value_get_flagset_flags(v)),
		mask:  uint(C.gst_value_get_flagset_mask(v)),
	}, nil
}

func marshalDoubleRange(p uintptr) (interface{}, error) {
	v := toGValue(p)
	return &Float64RangeValue{
		start: float64(C.gst_value_get_double_range_min(v)),
		end:   float64(C.gst_value_get_double_range_max(v)),
	}, nil
}

func marshalFraction(p uintptr) (interface{}, error) {
	v := toGValue(p)
	out := &FractionValue{
		num:   int(C.gst_value_get_fraction_numerator(v)),
		denom: int(C.gst_value_get_fraction_denominator(v)),
	}
	return out, nil
}

func marshalFractionRange(p uintptr) (interface{}, error) {
	v := toGValue(p)
	start := C.gst_value_get_fraction_range_min(v)
	end := C.gst_value_get_fraction_range_max(v)
	return &FractionRangeValue{
		start: ValueGetFraction(glib.ValueFromNative(unsafe.Pointer(start))),
		end:   ValueGetFraction(glib.ValueFromNative(unsafe.Pointer(end))),
	}, nil
}

func marshalBufferingMode(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return BufferingMode(c), nil
}

func marshalFormat(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return Format(c), nil
}

func marshalMessageType(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return MessageType(c), nil
}

func marshalPadLinkReturn(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return PadLinkReturn(c), nil
}

func marshalState(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return State(c), nil
}

func marshalSeekFlags(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return SeekFlags(c), nil
}

func marshalSeekType(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return SeekType(c), nil
}

func marshalStateChangeReturn(p uintptr) (interface{}, error) {
	c := C.g_value_get_enum(toGValue(p))
	return StateChangeReturn(c), nil
}

func marshalGhostPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapGhostPad(obj), nil
}

func marshalProxyPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapProxyPad(obj), nil
}

func marshalPad(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPad(obj), nil
}

func marshalMessage(p uintptr) (interface{}, error) {
	c := C.g_value_get_boxed(toGValue(p))
	return &Message{(*C.GstMessage)(unsafe.Pointer(c))}, nil
}

func marshalObject(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapObject(obj), nil
}

func marshalBus(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapBus(obj), nil
}

func marshalElementFactory(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapElementFactory(obj), nil
}

func marshalPipeline(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPipeline(obj), nil
}

func marshalPluginFeature(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapPluginFeature(obj), nil
}

func marshalElement(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapElement(obj), nil
}

func marshalBin(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(c))}
	return wrapBin(obj), nil
}

func marshalAllocationParams(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstAllocationParams)(unsafe.Pointer(c))
	return wrapAllocationParams(obj), nil
}

func marshalMemory(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstMemory)(unsafe.Pointer(c))
	return wrapMemory(obj), nil
}

func marshalBuffer(p uintptr) (interface{}, error) {
	c := C.getBufferValue(toGValue(p))
	return wrapBuffer(c), nil
}

func marshalBufferList(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstBufferList)(unsafe.Pointer(c))
	return wrapBufferList(obj), nil
}

func marshalCaps(p uintptr) (interface{}, error) {
	c := C.gst_value_get_caps(toGValue(p))
	obj := (*C.GstCaps)(unsafe.Pointer(c))
	return wrapCaps(obj), nil
}

func marshalCapsFeatures(p uintptr) (interface{}, error) {
	c := C.gst_value_get_caps_features(toGValue(p))
	obj := (*C.GstCapsFeatures)(unsafe.Pointer(c))
	return wrapCapsFeatures(obj), nil
}

func marshalStructure(p uintptr) (interface{}, error) {
	c := C.gst_value_get_structure(toGValue(p))
	obj := (*C.GstStructure)(unsafe.Pointer(c))
	return wrapStructure(obj), nil
}

func marshalContext(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstContext)(unsafe.Pointer(c))
	return wrapContext(obj), nil
}

func marshalTOC(p uintptr) (interface{}, error) {
	c := C.gst_value_get_structure(toGValue(p))
	obj := (*C.GstToc)(unsafe.Pointer(c))
	return wrapTOC(obj), nil
}

func marshalTOCEntry(p uintptr) (interface{}, error) {
	c := C.gst_value_get_structure(toGValue(p))
	obj := (*C.GstTocEntry)(unsafe.Pointer(c))
	return wrapTOCEntry(obj), nil
}

func marsalTagList(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstTagList)(unsafe.Pointer(c))
	return wrapTagList(obj), nil
}

func marshalEvent(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstEvent)(unsafe.Pointer(c))
	return wrapEvent(obj), nil
}

func marshalSegment(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstSegment)(unsafe.Pointer(c))
	return wrapSegment(obj), nil
}

func marshalQuery(p uintptr) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstQuery)(unsafe.Pointer(c))
	return wrapQuery(obj), nil
}

func marshalSample(p uintptr) (interface{}, error) {
	c := C.getSampleValue(toGValue(p))
	return wrapSample(c), nil
}
