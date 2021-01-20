package base

/*
#include "gst.go.h"

void
setGstBaseTransformPassthroughOnSameCaps (GstBaseTransform * obj, gboolean enabled)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(obj));
	GstBaseTransformClass * klass = toGstBaseTransformClass(g_type_class_peek_parent(this_class));
	klass->passthrough_on_same_caps = enabled;
}

void
setGstBaseTransformTransformIPOnPassthrough (GstBaseTransform * obj, gboolean enabled)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(obj));
	GstBaseTransformClass * klass = toGstBaseTransformClass(g_type_class_peek_parent(this_class));
	klass->transform_ip_on_passthrough = enabled;
}
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

// GstBaseTransformFlowDropped is a GstFlowReturn that can be returned from Transform() and TransformIP()
// to indicate that no output buffer was generated.
const GstBaseTransformFlowDropped gst.FlowReturn = C.GST_BASE_TRANSFORM_FLOW_DROPPED

// GstBaseTransform represents a GstBaseTransform.
type GstBaseTransform struct{ *gst.Element }

// ToGstBaseTransform returns a GstBaseTransform object for the given object. It will work on either gst.Object
// or glib.Object interfaces.
func ToGstBaseTransform(obj interface{}) *GstBaseTransform {
	switch obj := obj.(type) {
	case *gst.Object:
		return &GstBaseTransform{&gst.Element{Object: obj}}
	case *glib.Object:
		return &GstBaseTransform{&gst.Element{Object: &gst.Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: obj}}}}
	}
	return nil
}

// Instance returns the underlying C GstBaseTransform instance
func (g *GstBaseTransform) Instance() *C.GstBaseTransform {
	return C.toGstBaseTransform(g.Unsafe())
}

// GetAllocator retrieves the memory allocator used by this base transform. Unref after usage.
func (g *GstBaseTransform) GetAllocator() (*gst.Allocator, *gst.AllocationParams) {
	var allocParams C.GstAllocationParams
	var allocator *C.GstAllocator
	C.gst_base_transform_get_allocator(g.Instance(), &allocator, &allocParams)
	if allocator == nil {
		return nil, nil
	}
	return gst.FromGstAllocatorUnsafeFull(unsafe.Pointer(allocator)), gst.FromGstAllocationParamsUnsafe(unsafe.Pointer(&allocParams))
}

// GetBufferPool returns the BufferPool used by this transform. Unref after usage.
func (g *GstBaseTransform) GetBufferPool() *gst.BufferPool {
	pool := C.gst_base_transform_get_buffer_pool(g.Instance())
	if pool == nil {
		return nil
	}
	return gst.FromGstBufferPoolUnsafeFull(unsafe.Pointer(pool))
}

// IsInPlace returns if the transform is configured to do in-place transform.
func (g *GstBaseTransform) IsInPlace() bool {
	return gobool(C.gst_base_transform_is_in_place(g.Instance()))
}

// IsPassthrough returns if the transform is configured for passthrough.
func (g *GstBaseTransform) IsPassthrough() bool {
	return gobool(C.gst_base_transform_is_passthrough(g.Instance()))
}

// IsQoSEnabled queries if the transform will handle QoS.
func (g *GstBaseTransform) IsQoSEnabled() bool {
	return gobool(C.gst_base_transform_is_qos_enabled(g.Instance()))
}

// QueuedBuffer returns the currentl queued buffer.
func (g *GstBaseTransform) QueuedBuffer() *gst.Buffer {
	return gst.FromGstBufferUnsafeNone(unsafe.Pointer(g.Instance().queued_buf))
}

// SINCE 1.18
// // Reconfigure negotiates src pad caps with downstream elements if the source pad is marked as needing
// // reconfiguring. Unmarks GST_PAD_FLAG_NEED_RECONFIGURE in any case. But marks it again if negotiation fails.
// //
// // Do not call this in the Transform() or TransformIP() vmethod. Call this in SubmitInputBuffer(),
// // PrepareOutputBuffer() or in GenerateOutput() before any output buffer is allocated.
// //
// // It will by default be called when handling an ALLOCATION query or at the very beginning of the default
// // SubmitInputBuffer() implementation.
// func (g *GstBaseTransform) Reconfigure() bool {
// 	return gobool(C.gst_base_transform_reconfigure(g.Instance()))
// }

// ReconfigureSink instructs transform to request renegotiation upstream. This function is typically called
// after properties on the transform were set that influence the input format.
func (g *GstBaseTransform) ReconfigureSink() { C.gst_base_transform_reconfigure_sink(g.Instance()) }

// ReconfigureSrc instructs trans to renegotiate a new downstream transform on the next buffer. This function
// is typically called after properties on the transform were set that influence the output format.
func (g *GstBaseTransform) ReconfigureSrc() { C.gst_base_transform_reconfigure_src(g.Instance()) }

// SetGapAware configures how buffers are handled. If gapAware is FALSE (the default), output buffers will
// have the GST_BUFFER_FLAG_GAP flag unset.
//
// If set to TRUE, the element must handle output buffers with this flag set correctly, i.e. it can assume that
// the buffer contains neutral data but must unset the flag if the output is no neutral data.
func (g *GstBaseTransform) SetGapAware(gapAware bool) {
	C.gst_base_transform_set_gap_aware(g.Instance(), gboolean(gapAware))
}

// SetInPlace determines whether a non-writable buffer will be copied before passing to the TransformIP function.
// This is always true if no Transform() function is implemented, and always false if ONLY a Transform() function
// is implemented.
func (g *GstBaseTransform) SetInPlace(inPlace bool) {
	C.gst_base_transform_set_in_place(g.Instance(), gboolean(inPlace))
}

// SetPassthrough sets the default passthrough mode for this filter. The is mostly useful for filters that do not
// care about negotiation. This is always true for filters which don't implement either a Transform(), TransformIP(),
// or GenerateOutput() method.
func (g *GstBaseTransform) SetPassthrough(passthrough bool) {
	C.gst_base_transform_set_passthrough(g.Instance(), gboolean(passthrough))
}

// SetPassthroughOnSameCaps when set to true will automatically enable passthrough if caps are the same.
func (g *GstBaseTransform) SetPassthroughOnSameCaps(passthrough bool) {
	C.setGstBaseTransformPassthroughOnSameCaps(g.Instance(), gboolean(passthrough))
}

// SetPreferPassthrough sets whether passthrough is preferred. If preferPassthrough is TRUE (the default), trans
// will check and prefer passthrough caps from the list of caps returned by the TransformCaps() vmethod.
//
// If set to FALSE, the element must order the caps returned from the TransformCaps() function in such a way that
// the preferred format is first in the list. This can be interesting for transforms that can do passthrough
// transforms but prefer to do something else, like a capsfilter.
func (g *GstBaseTransform) SetPreferPassthrough(preferPassthrough bool) {
	C.gst_base_transform_set_prefer_passthrough(g.Instance(), gboolean(preferPassthrough))
}

// SetQoSEnabled enables or disables QoS handling in the filter.
func (g *GstBaseTransform) SetQoSEnabled(enabled bool) {
	C.gst_base_transform_set_qos_enabled(g.Instance(), gboolean(enabled))
}

// SetTransformIPOnPassthrough If set to TRUE, TransformIP() will be called in passthrough mode. The passed
// buffer might not be writable. When FALSE, neither Transform() nor TransformIP() will be called in passthrough
// mode. Set to TRUE by default.
func (g *GstBaseTransform) SetTransformIPOnPassthrough(enabled bool) {
	C.setGstBaseTransformTransformIPOnPassthrough(g.Instance(), gboolean(enabled))
}

// SinkPad returns the sink pad object for this element.
func (g *GstBaseTransform) SinkPad() *gst.Pad {
	return gst.FromGstPadUnsafeNone(unsafe.Pointer(g.Instance().sinkpad))
}

// SrcPad returns the src pad object for this element.
func (g *GstBaseTransform) SrcPad() *gst.Pad {
	return gst.FromGstPadUnsafeNone(unsafe.Pointer(g.Instance().srcpad))
}

// UpdateQoS sets the QoS parameters in the transform. This function is called internally when a QOS event is received
// but subclasses can provide custom information when needed.
//
// proportion is the proportion, diff is the diff against the clock, and timestamp is the timestamp of the buffer
// generating the QoS expressed in running_time.
func (g *GstBaseTransform) UpdateQoS(proportion float64, diff, timestamp time.Duration) {
	C.gst_base_transform_update_qos(
		g.Instance(),
		C.gdouble(proportion),
		C.GstClockTimeDiff(diff.Nanoseconds()),
		C.GstClockTime(timestamp.Nanoseconds()),
	)
}

// UpdateSrcCaps updates the srcpad caps and sends the caps downstream. This function can be used by subclasses
// when they have already negotiated their caps but found a change in them (or computed new information). This
// way, they can notify downstream about that change without losing any buffer.
func (g *GstBaseTransform) UpdateSrcCaps(caps *gst.Caps) {
	C.gst_base_transform_update_src_caps(
		g.Instance(),
		(*C.GstCaps)(unsafe.Pointer(caps.Instance())),
	)
}
