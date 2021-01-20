package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

// GstBaseSrc represents a GstBaseSrc.
type GstBaseSrc struct{ *gst.Element }

// ToGstBaseSrc returns a GstBaseSrc object for the given object. It will work on either gst.Object
// or glib.Object interfaces.
func ToGstBaseSrc(obj interface{}) *GstBaseSrc {
	switch obj := obj.(type) {
	case *gst.Object:
		return &GstBaseSrc{&gst.Element{Object: obj}}
	case *glib.Object:
		return &GstBaseSrc{&gst.Element{Object: &gst.Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: obj}}}}
	}
	return nil
}

// Instance returns the underlying C GstBaseSrc instance
func (g *GstBaseSrc) Instance() *C.GstBaseSrc {
	return C.toGstBaseSrc(g.Unsafe())
}

// GetAllocator retrieves the memory allocator used by this base src. Unref after usage.
func (g *GstBaseSrc) GetAllocator() (*gst.Allocator, *gst.AllocationParams) {
	var allocParams C.GstAllocationParams
	var allocator *C.GstAllocator
	C.gst_base_src_get_allocator(g.Instance(), &allocator, &allocParams)
	if allocator == nil {
		return nil, nil
	}
	return gst.FromGstAllocatorUnsafeFull(unsafe.Pointer(allocator)), gst.FromGstAllocationParamsUnsafe(unsafe.Pointer(&allocParams))
}

// GetBlocksize returns the number of bytes that the source will push out with each buffer.
func (g *GstBaseSrc) GetBlocksize() uint { return uint(C.gst_base_src_get_blocksize(g.Instance())) }

// GetBufferPool returns the BufferPool used by this source. Unref after usage.
func (g *GstBaseSrc) GetBufferPool() *gst.BufferPool {
	pool := C.gst_base_src_get_buffer_pool(g.Instance())
	if pool == nil {
		return nil
	}
	return gst.FromGstBufferPoolUnsafeFull(unsafe.Pointer(pool))
}

// DoTimestamp will query if the timestamps outgoing on this source's buffers are based on the current
// running time.
func (g *GstBaseSrc) DoTimestamp() bool {
	return gobool(C.gst_base_src_get_do_timestamp(g.Instance()))
}

// IsAsync retrieves the current async behavior of the source.
func (g *GstBaseSrc) IsAsync() bool { return gobool(C.gst_base_src_is_async(g.Instance())) }

// IsLive checks if this source is in live mode.
func (g *GstBaseSrc) IsLive() bool { return gobool(C.gst_base_src_is_live(g.Instance())) }

// SINCE 1.18
// // Negotiate negotiates this source's pad caps with downstream elements. Do not call this in the Fill()
// // vmethod. Call this in Create() or in Alloc(), before any buffer is allocated.
// func (g *GstBaseSrc) Negotiate() bool { return gobool(C.gst_base_src_negotiate(g.Instance())) }

// // NewSegment prepares a new segment for emission downstream. This function must only be called by derived
// // sub-classes, and only from the create function, as the stream-lock needs to be held.
// //
// // The format for the segment must be identical with the current format of the source, as configured with
// // SetFormat.
// //
// // The format of src must not be gst.FormatUndefined and the format should be configured via SetFormat before
// // calling this method.
// func (g *GstBaseSrc) NewSegment(segment *gst.Segment) bool {
// 	return gobool(C.gst_base_src_new_segment(g.Instance(), (*C.GstSegment)(unsafe.Pointer(segment.Instance()))))
// }

// QueryLatency queries the source for the latency parameters. live will be TRUE when src is configured as a
// live source. minLatency and maxLatency will be set to the difference between the running time and the timestamp
// of the first buffer.
//
// This function is mostly used by subclasses.
func (g *GstBaseSrc) QueryLatency() (ok, live bool, minLatency, maxLatency time.Duration) {
	var glive C.gboolean
	var gmin C.GstClockTime
	var gmax C.GstClockTime
	gok := C.gst_base_src_query_latency(g.Instance(), &glive, &gmin, &gmax)
	return gobool(gok), gobool(glive), time.Duration(gmin), time.Duration(gmax)
}

// SetAsync configures async behaviour in src, no state change will block. The open, close, start, stop, play and
// pause virtual methods will be executed in a different thread and are thus allowed to perform blocking operations.
// Any blocking operation should be unblocked with the unlock vmethod.
func (g *GstBaseSrc) SetAsync(async bool) { C.gst_base_src_set_async(g.Instance(), gboolean(async)) }

// SetAutomaticEOS sets whether EOS should be automatically emmitted.
//
// If automaticEOS is TRUE, src will automatically go EOS if a buffer after the total size is returned. By default
// this is TRUE but sources that can't return an authoritative size and only know that they're EOS when trying to
// read more should set this to FALSE.
//
// When src operates in gst.FormatTime, GstBaseSrc will send an EOS when a buffer outside of the currently configured
// segment is pushed if automaticEOS is TRUE. Since 1.16, if automatic_eos is FALSE an EOS will be pushed only when
// the Create() implementation returns gst.FlowEOS.
func (g *GstBaseSrc) SetAutomaticEOS(automaticEOS bool) {
	C.gst_base_src_set_automatic_eos(g.Instance(), gboolean(automaticEOS))
}

// SetBlocksize sets the number of bytes that src will push out with each buffer. When blocksize is set to -1, a
// default length will be used.
func (g *GstBaseSrc) SetBlocksize(size uint) {
	C.gst_base_src_set_blocksize(g.Instance(), C.guint(size))
}

// SetCaps sets new caps on the source pad.
func (g *GstBaseSrc) SetCaps(caps *gst.Caps) bool {
	return gobool(C.gst_base_src_set_caps(g.Instance(), (*C.GstCaps)(unsafe.Pointer(caps.Instance()))))
}

// SetDoTimestamp configures src to automatically timestamp outgoing buffers based on the current running_time of the pipeline.
// This property is mostly useful for live sources.
func (g *GstBaseSrc) SetDoTimestamp(doTimestamp bool) {
	C.gst_base_src_set_do_timestamp(g.Instance(), gboolean(doTimestamp))
}

// SetDynamicSize sets if the size is dynamic for this source.
//
// If not dynamic, size is only updated when needed, such as when trying to read past current tracked size. Otherwise, size is
// checked for upon each read.
func (g *GstBaseSrc) SetDynamicSize(dynamic bool) {
	C.gst_base_src_set_dynamic_size(g.Instance(), gboolean(dynamic))
}

// SetFormat sets the default format of the source. This will be the format used for sending
// SEGMENT events and for performing seeks.
//
// If a format of gst.FormatBytes is set, the element will be able to operate in pull mode if the
// IsSeekable returns TRUE.
//
// This function must only be called in when the element is paused.
func (g *GstBaseSrc) SetFormat(format gst.Format) {
	C.gst_base_src_set_format(g.Instance(), C.GstFormat(format))
}

// SetLive sets if the element listens to a live source.
//
// A live source will not produce data in the PAUSED state and will therefore not be able to participate in the
// PREROLL phase of a pipeline. To signal this fact to the application and the pipeline, the state change return
// value of the live source will be gst.StateChangeNoPreroll.
func (g *GstBaseSrc) SetLive(live bool) { C.gst_base_src_set_live(g.Instance(), gboolean(live)) }

// StartComplete completes an asynchronous start operation. When the subclass overrides the start method,
// it should call StartComplete when the start operation completes either from the same thread or from an
// asynchronous helper thread.
func (g *GstBaseSrc) StartComplete(ret gst.FlowReturn) {
	C.gst_base_src_start_complete(g.Instance(), C.GstFlowReturn(ret))
}

// StartWait waits until the start operation is complete.
func (g *GstBaseSrc) StartWait() gst.FlowReturn {
	return gst.FlowReturn(C.gst_base_src_start_wait(g.Instance()))
}

// SubmitBufferList submits a list of buffers to the source.
//
// Subclasses can call this from their create virtual method implementation to submit a buffer list to be pushed out
// later. This is useful in cases where the create function wants to produce multiple buffers to be pushed out in one
// go in form of a GstBufferList, which can reduce overhead drastically, especially for packetised inputs (for data
// streams where the packetisation/chunking is not important it is usually more efficient to return larger buffers instead).
//
// Subclasses that use this function from their create function must return GST_FLOW_OK and no buffer from their create
// virtual method implementation. If a buffer is returned after a buffer list has also been submitted via this function the
// behaviour is undefined.
//
// Subclasses must only call this function once per create function call and subclasses must only call this function when
// the source operates in push mode.
func (g *GstBaseSrc) SubmitBufferList(bufferList *gst.BufferList) {
	C.gst_base_src_submit_buffer_list(g.Instance(), (*C.GstBufferList)(unsafe.Pointer(bufferList.Instance())))
}

// WaitPlaying will block until a state change to PLAYING happens (in which case this function returns gst.FlowOK) or the
// processing must be stopped due to a state change to READY or a FLUSH event (in which case this function returns GST_FLOW_FLUSHING).
//
// If the Create() method performs its own synchronisation against the clock it must unblock when going from PLAYING to the PAUSED
// state and call this method before continuing to produce the remaining data.
//
// gst.FlowOK will be returned if the source is PLAYING and processing can continue. Any other return value should be
// returned from the Create() vmethod.
func (g *GstBaseSrc) WaitPlaying() gst.FlowReturn {
	return gst.FlowReturn(C.gst_base_src_wait_playing(g.Instance()))
}
