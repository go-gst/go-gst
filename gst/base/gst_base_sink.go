package base

/*
#include "gst.go.h"

gboolean baseSinkParentEvent (GstBaseSink * sink, GstEvent * event)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(sink));
	GstBaseSinkClass * parent = toGstBaseSinkClass(g_type_class_peek_parent(this_class));
	return parent->event(sink, event);
}
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

// GstBaseSink represents a GstBaseSink.
type GstBaseSink struct{ *gst.Element }

// ToGstBaseSink returns a GstBaseSink object for the given object. It will work on either gst.Object
// or glib.Object interfaces.
func ToGstBaseSink(obj interface{}) *GstBaseSink {
	switch obj := obj.(type) {
	case *gst.Object:
		return &GstBaseSink{&gst.Element{Object: obj}}
	case *glib.Object:
		return &GstBaseSink{&gst.Element{Object: &gst.Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: obj}}}}
	}
	return nil
}

// Instance returns the underlying C GstBaseSrc instance
func (g *GstBaseSink) Instance() *C.GstBaseSink {
	return C.toGstBaseSink(g.Unsafe())
}

// DoPreroll is for if the sink spawns its own thread for pulling buffers from upstream.
// It should call this method after it has pulled a buffer. If the element needed to preroll,
// this function will perform the preroll and will then block until the element state is changed.
//
// This function should be called with the PREROLL_LOCK held and the object that caused the preroll.
//
// Since the object will always be a gst.MiniObject (which is not implemented properly), this method will check
// against the provided types for structs known to be used in this context. The currently known options
// are events, messages, queries, structures, and buffers. If you come across a need to use this function with
// an unsupported type, feel free to raise an Issue or open a PR.
func (g *GstBaseSink) DoPreroll(obj interface{}) gst.FlowReturn {
	miniobj := getPrerollObj(obj)
	if miniobj == nil {
		return gst.FlowError
	}
	return gst.FlowReturn(C.gst_base_sink_do_preroll(g.Instance(), miniobj))
}

func getPrerollObj(obj interface{}) *C.GstMiniObject {
	switch obj := obj.(type) {
	case *gst.Event:
		return (*C.GstMiniObject)(unsafe.Pointer(obj.Instance()))
	case *gst.Buffer:
		return (*C.GstMiniObject)(unsafe.Pointer(obj.Instance()))
	case *gst.Message:
		return (*C.GstMiniObject)(unsafe.Pointer(obj.Instance()))
	case *gst.Query:
		return (*C.GstMiniObject)(unsafe.Pointer(obj.Instance()))
	case *gst.Structure:
		return (*C.GstMiniObject)(unsafe.Pointer(obj.Instance()))
	default:
		return nil
	}
}

// GetBlocksize gets the number of bytes that the sink will pull when it is operating in pull mode.
func (g *GstBaseSink) GetBlocksize() uint { return uint(C.gst_base_sink_get_blocksize(g.Instance())) }

// GetDropOutOfSegment checks if sink is currently configured to drop buffers which are outside the current segment
func (g *GstBaseSink) GetDropOutOfSegment() bool {
	return gobool(C.gst_base_sink_get_drop_out_of_segment(g.Instance()))
}

// GetLastSample gets the last sample that arrived in the sink and was used for preroll or for rendering.
// This property can be used to generate thumbnails.
//
// The GstCaps on the sample can be used to determine the type of the buffer. Unref after usage. Sample will
// be nil if no buffer has arrived yet.
func (g *GstBaseSink) GetLastSample() *gst.Sample {
	sample := C.gst_base_sink_get_last_sample(g.Instance())
	if sample == nil {
		return nil
	}
	return gst.FromGstSampleUnsafeFull(unsafe.Pointer(sample))
}

// GetLatency gets the currently configured latency.
func (g *GstBaseSink) GetLatency() time.Duration {
	return time.Duration(C.gst_base_sink_get_latency(g.Instance()))
}

// GetMaxBitrate gets the maximum amount of bits per second the sink will render.
func (g *GstBaseSink) GetMaxBitrate() uint64 {
	return uint64(C.gst_base_sink_get_max_bitrate(g.Instance()))
}

// GetMaxLateness gets the max lateness value.
func (g *GstBaseSink) GetMaxLateness() int64 {
	return int64(C.gst_base_sink_get_max_lateness(g.Instance()))
}

// GetProcessingDeadline gets the processing deadline of the sink.
func (g *GstBaseSink) GetProcessingDeadline() time.Duration {
	return time.Duration(C.gst_base_sink_get_processing_deadline(g.Instance()))
}

// GetRenderDelay gets the render delay for the sink.
func (g *GstBaseSink) GetRenderDelay() time.Duration {
	return time.Duration(C.gst_base_sink_get_render_delay(g.Instance()))
}

// SINCE 1.18
// // SinkStats represents the current statistics on a GstBaseSink.
// type SinkStats struct {
// 	AverageRate float64
// 	Dropped     uint64
// 	Rendered    uint64
// }

// // GetSinkStats returns various GstBaseSink statistics.
// func (g *GstBaseSink) GetSinkStats() *SinkStats {
// 	st := gst.FromGstStructureUnsafe(unsafe.Pointer(C.gst_base_sink_get_stats(g.Instance())))
// 	stats := &SinkStats{}
// 	if avgRate, err := st.GetValue("average-rate"); err == nil {
// 		stats.AverageRate = avgRate.(float64)
// 	}
// 	if dropped, err := st.GetValue("dropped"); err == nil {
// 		stats.Dropped = dropped.(uint64)
// 	}
// 	if rendered, err := st.GetValue("rendered"); err == nil {
// 		stats.Rendered = rendered.(uint64)
// 	}
// 	return stats
// }

// GetSync checks if the sink is currently configured to synchronize on the clock.
func (g *GstBaseSink) GetSync() bool { return gobool(C.gst_base_sink_get_sync(g.Instance())) }

// GetThrottleTime gets the time that will be inserted between frames to control maximum buffers
// per second.
func (g *GstBaseSink) GetThrottleTime() uint64 {
	return uint64(C.gst_base_sink_get_throttle_time(g.Instance()))
}

// GetTsOffset gets the synchronization offset of sink.
func (g *GstBaseSink) GetTsOffset() time.Duration {
	return time.Duration(C.gst_base_sink_get_ts_offset(g.Instance()))
}

// IsAsyncEnabled checks if the sink is currently configured to perform asynchronous state changes to PAUSED.
func (g *GstBaseSink) IsAsyncEnabled() bool {
	return gobool(C.gst_base_sink_is_async_enabled(g.Instance()))
}

// IsLastSampleEnabled checks if the sink is currently configured to store the last received sample.
func (g *GstBaseSink) IsLastSampleEnabled() bool {
	return gobool(C.gst_base_sink_is_last_sample_enabled(g.Instance()))
}

// IsQoSEnabled checks if sink is currently configured to send QoS events upstream.
func (g *GstBaseSink) IsQoSEnabled() bool {
	return gobool(C.gst_base_sink_is_qos_enabled(g.Instance()))
}

// ParentEvent calls the parent class's event handler for the given event.
func (g *GstBaseSink) ParentEvent(ev *gst.Event) bool {
	return gobool(C.baseSinkParentEvent(
		g.Instance(),
		(*C.GstEvent)(unsafe.Pointer(ev.Instance())),
	))
}

// QueryLatency queries the sink for the latency parameters. The latency will be queried from the
// upstream elements. live will be TRUE if sink is configured to synchronize against the clock.
// upstreamLive will be TRUE if an upstream element is live.
//
// If both live and upstreamLive are TRUE, the sink will want to compensate for the latency introduced
// by the upstream elements by setting the minLatency to a strictly positive value.
//
// This function is mostly used by subclasses.
func (g *GstBaseSink) QueryLatency() (ok, live, upstreamLive bool, minLatency, maxLatency time.Duration) {
	var glive, gupLive C.gboolean
	var gmin, gmax C.GstClockTime
	ret := C.gst_base_sink_query_latency(g.Instance(), &glive, &gupLive, &gmin, &gmax)
	return gobool(ret), gobool(glive), gobool(gupLive), time.Duration(gmin), time.Duration(gmax)
}

// SetAsyncEnabled configures sink to perform all state changes asynchronously. When async is disabled,
// the sink will immediately go to PAUSED instead of waiting for a preroll buffer. This feature is useful
// if the sink does not synchronize against the clock or when it is dealing with sparse streams.
func (g *GstBaseSink) SetAsyncEnabled(enabled bool) {
	C.gst_base_sink_set_async_enabled(g.Instance(), gboolean(enabled))
}

// SetBlocksize sets the number of bytes this sink will pull when operating in pull mode.
func (g *GstBaseSink) SetBlocksize(blocksize uint) {
	C.gst_base_sink_set_blocksize(g.Instance(), C.guint(blocksize))
}

// SetDropOutOfSegment configures sink to drop buffers which are outside the current segment.
func (g *GstBaseSink) SetDropOutOfSegment(drop bool) {
	C.gst_base_sink_set_drop_out_of_segment(g.Instance(), gboolean(drop))
}

// SetLastSampleEnabled configures the sink to store the last received sample.
func (g *GstBaseSink) SetLastSampleEnabled(enabled bool) {
	C.gst_base_sink_set_last_sample_enabled(g.Instance(), gboolean(enabled))
}

// SetMaxBitrate sets the maximum amount of bits per second the sink will render.
func (g *GstBaseSink) SetMaxBitrate(bitrate uint64) {
	C.gst_base_sink_set_max_bitrate(g.Instance(), C.guint64(bitrate))
}

// SetMaxLateness sets the new max lateness value to max_lateness. This value is used to decide if
// a buffer should be dropped or not based on the buffer timestamp and the current clock time. A
// value of -1 means an unlimited time.
func (g *GstBaseSink) SetMaxLateness(maxLateness int64) {
	C.gst_base_sink_set_max_lateness(g.Instance(), C.gint64(maxLateness))
}

// SetProcessingDeadline sets the maximum amount of time (in nanoseconds) that the pipeline can take
// for processing the buffer. This is added to the latency of live pipelines.
//
// This function is usually called by subclasses.
func (g *GstBaseSink) SetProcessingDeadline(deadline time.Duration) {
	C.gst_base_sink_set_processing_deadline(g.Instance(), C.GstClockTime(deadline.Nanoseconds()))
}

// SetQoSEnabled configures sink to send Quality-of-Service events upstream.
func (g *GstBaseSink) SetQoSEnabled(enabled bool) {
	C.gst_base_sink_set_qos_enabled(g.Instance(), gboolean(enabled))
}

// SetRenderDelay sets the render delay in sink to delay. The render delay is the time between actual
// rendering of a buffer and its synchronisation time. Some devices might delay media rendering which
// can be compensated for with this function.
//
// After calling this function, this sink will report additional latency and other sinks will adjust
// their latency to delay the rendering of their media.
//
// This function is usually called by subclasses.
func (g *GstBaseSink) SetRenderDelay(delay time.Duration) {
	C.gst_base_sink_set_render_delay(g.Instance(), C.GstClockTime(delay.Nanoseconds()))
}

// SetSync configures sink to synchronize on the clock or not. When sync is FALSE, incoming samples will
// be played as fast as possible. If sync is TRUE, the timestamps of the incoming buffers will be used to
// schedule the exact render time of its contents.
func (g *GstBaseSink) SetSync(sync bool) { C.gst_base_sink_set_sync(g.Instance(), gboolean(sync)) }

// SetThrottleTime sets the time that will be inserted between rendered buffers. This can be used to control
// the maximum buffers per second that the sink will render.
func (g *GstBaseSink) SetThrottleTime(throttle uint64) {
	C.gst_base_sink_set_throttle_time(g.Instance(), C.guint64(throttle))
}

// SetTsOffset adjusts the synchronization of sink with offset. A negative value will render buffers earlier
// than their timestamp. A positive value will delay rendering. This function can be used to fix playback of
// badly timestamped buffers.
func (g *GstBaseSink) SetTsOffset(offset time.Duration) {
	C.gst_base_sink_set_ts_offset(g.Instance(), C.GstClockTimeDiff(offset.Nanoseconds()))
}

// Wait will wait for preroll to complete and will then block until timeout is reached. It is usually called by
// subclasses that use their own internal synchronization but want to let some synchronization (like EOS) be
// handled by the base class.
//
// This function should only be called with the PREROLL_LOCK held (like when receiving an EOS event in the
// ::event vmethod or when handling buffers in ::render).
//
// The timeout argument should be the running_time of when the timeout should happen and will be adjusted with any
// latency and offset configured in the sink.
func (g *GstBaseSink) Wait(timeout time.Duration) (ret gst.FlowReturn, jitter time.Duration) {
	var jit C.GstClockTimeDiff
	gret := C.gst_base_sink_wait(g.Instance(), C.GstClockTime(timeout.Nanoseconds()), &jit)
	return gst.FlowReturn(gret), time.Duration(jit)
}

// WaitClock will block until timeout is reached. It is usually called by subclasses that use their own
// internal synchronization.
//
// If time is not valid, no synchronisation is done and GST_CLOCK_BADTIME is returned. Likewise, if synchronization
// is disabled in the element or there is no clock, no synchronization is done and GST_CLOCK_BADTIME is returned.
//
// This function should only be called with the PREROLL_LOCK held, like when receiving an EOS event in the event()
// vmethod or when receiving a buffer in the render() vmethod.
//
// The timeout argument should be the running_time of when this method should return and is not adjusted with any
// latency or offset configured in the sink.
func (g *GstBaseSink) WaitClock(timeout time.Duration) (ret gst.ClockReturn, jitter time.Duration) {
	var jit C.GstClockTimeDiff
	gret := C.gst_base_sink_wait_clock(g.Instance(), C.GstClockTime(timeout.Nanoseconds()), &jit)
	return gst.ClockReturn(gret), time.Duration(jit)
}

// WaitPreroll will block until the preroll is complete.
//
// If the render() method performs its own synchronisation against the clock it must unblock when going from
// PLAYING to the PAUSED state and call this method before continuing to render the remaining data.
//
// If the render() method can block on something else than the clock, it must also be ready to unblock immediately
// on the unlock() method and cause the render() method to immediately call this function. In this case, the
// subclass must be prepared to continue rendering where it left off if this function returns GST_FLOW_OK.
//
// This function will block until a state change to PLAYING happens (in which case this function returns
// GST_FLOW_OK) or the processing must be stopped due to a state change to READY or a FLUSH event (in which case
// this function returns GST_FLOW_FLUSHING).
//
// This function should only be called with the PREROLL_LOCK held, like in the render function.
func (g *GstBaseSink) WaitPreroll() gst.FlowReturn {
	return gst.FlowReturn(C.gst_base_sink_wait_preroll(g.Instance()))
}
