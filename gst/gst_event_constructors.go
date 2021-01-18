package gst

// #include "gst.go.h"
import "C"

import (
	"time"
	"unsafe"
)

// NewBufferSizeEvent creates a new buffersize event. The event is sent downstream and notifies
// elements that they should provide a buffer of the specified dimensions.
//
// When the async flag is set, a thread boundary is preferred.
func NewBufferSizeEvent(format Format, minSize, maxSize int64, async bool) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_buffer_size(
		C.GstFormat(format),
		C.gint64(minSize),
		C.gint64(maxSize),
		gboolean(async),
	)))
}

// NewCapsEvent creates a new CAPS event for caps. The caps event can only travel downstream synchronized with
// the buffer flow and contains the format of the buffers that will follow after the event.
func NewCapsEvent(caps *Caps) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_caps(
		caps.Instance(),
	)))
}

// NewEOSEvent creates a new EOS event. The eos event can only travel downstream synchronized with the buffer flow.
// Elements that receive the EOS event on a pad can return FlowEOS as a FlowReturn when data after the EOS event arrives.
//
// The EOS event will travel down to the sink elements in the pipeline which will then post the MessageEOS on the bus
// after they have finished playing any buffered data.
//
// When all sinks have posted an EOS message, an EOS message is forwarded to the application.
//
// The EOS event itself will not cause any state transitions of the pipeline.
func NewEOSEvent() *Event { return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_eos())) }

// NewFlushStartEvent allocates a new flush start event. The flush start event can be sent upstream and downstream and
// travels out-of-bounds with the dataflow.
//
// It marks pads as being flushing and will make them return FlowFlushing when used for data flow with gst_pad_push,
// gst_pad_chain, gst_pad_get_range and gst_pad_pull_range. Any event (except a EventFlushSStop) received on a flushing
// pad will return FALSE immediately.
//
// Elements should unlock any blocking functions and exit their streaming functions as fast as possible when this event is received.
//
// This event is typically generated after a seek to flush out all queued data in the pipeline so that the new media is played as soon as possible.
func NewFlushStartEvent() *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_flush_start()))
}

// NewFlushStopEvent allocates a new flush stop event. The flush stop event can be sent upstream and downstream and travels serialized with the
// dataflow. It is typically sent after sending a FLUSH_START event to make the pads accept data again.
//
// Elements can process this event synchronized with the dataflow since the preceding FLUSH_START event stopped the dataflow.
//
// This event is typically generated to complete a seek and to resume dataflow.
func NewFlushStopEvent(resetTime bool) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_flush_stop(gboolean(resetTime))))
}

// NewGapEvent creates a new GAP event. A gap event can be thought of as conceptually equivalent to a buffer to signal that there is no data for a
//certain amount of time. This is useful to signal a gap to downstream elements which may wait for data, such as muxers or mixers or overlays,
// especially for sparse streams such as subtitle streams.
func NewGapEvent(timestamp, duration time.Duration) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_gap(
		C.GstClockTime(timestamp.Nanoseconds()),
		C.GstClockTime(duration.Nanoseconds()),
	)))
}

// NewInstantRateChangeEvent creates a new instant-rate-change event. This event is sent by seek handlers (e.g. demuxers) when receiving a seek with the
// GST_SEEK_FLAG_INSTANT_RATE_CHANGE and signals to downstream elements that the playback rate in the existing segment should be immediately multiplied
// by the rate_multiplier factor.
//
// The flags provided replace any flags in the existing segment, for the flags within the GST_SEGMENT_INSTANT_FLAGS set. Other GstSegmentFlags are ignored
// and not transferred in the event.
// func NewInstantRateChangeEvent(rateMultiplier float64, newFlags SegmentFlags) *Event {
// 	return wrapEvent(C.gst_event_new_instant_rate_change(
// 		C.gdouble(rateMultiplier),
// 		C.GstSegmentFlags(newFlags),
// 	))
// }

// NewInstantRateSyncTimeEvent creates a new instant-rate-sync-time event. This event is sent by the pipeline to notify elements handling the instant-rate-change
// event about the running-time when the new rate should be applied. The running time may be in the past when elements handle this event, which can lead to switching
// artifacts. The magnitude of those depends on the exact timing of event delivery to each element and the magnitude of the change in playback rate being applied.
//
// The running_time and upstream_running_time are the same if this is the first instant-rate adjustment, but will differ for later ones to compensate for the
// accumulated offset due to playing at a rate different to the one indicated in the playback segments.
// func NewInstantRateSyncTimeEvent(rateMultiplier float64, runningTime, upstreamRunningTime time.Duration) *Event {
// 	return wrapEvent(C.gst_event_new_instant_rate_sync_time(
// 		C.gdouble(rateMultiplier),
// 		C.GstClockTime(durationToClockTime(runningTime)),
// 		C.GstClockTime(durationToClockTime(upstreamRunningTime)),
// 	))
// }

// NewLatencyEvent creates a new latency event. The event is sent upstream from the sinks and notifies elements that they should add an additional latency to the
// running time before synchronising against the clock.
//
// The latency is mostly used in live sinks and is always expressed in the time format.
func NewLatencyEvent(latency time.Duration) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_latency(
		C.GstClockTime(latency.Nanoseconds()),
	)))
}

// NewNavigationEvent creates a new navigation event from the given description. The event will take ownership of the structure.
func NewNavigationEvent(structure *Structure) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_navigation(
		structure.Instance(),
	)))
}

// NewProtectionEvent creates a new event containing information specific to a particular protection system (uniquely identified by system_id), by which that protection system
// can acquire key(s) to decrypt a protected stream.
//
// In order for a decryption element to decrypt media protected using a specific system, it first needs all the protection system specific information necessary to acquire the
// decryption key(s) for that stream. The functions defined here enable this information to be passed in events from elements that extract it (e.g., ISOBMFF demuxers, MPEG DASH
// demuxers) to protection decrypter elements that use it.
//
// Events containing protection system specific information are created using NewProtectionEvent, and they can be parsed by downstream elements using ParseProtection.
//
// In Common Encryption, protection system specific information may be located within ISOBMFF files, both in movie (moov) boxes and movie fragment (moof) boxes; it may also be
// contained in ContentProtection elements within MPEG DASH MPDs. The events created by gst_event_new_protection contain data identifying from which of these locations the encapsulated
// protection system specific information originated. This origin information is required as some protection systems use different encodings depending upon where the information originates.
//
// The events returned by NewProtectionEvent are implemented in such a way as to ensure that the most recently-pushed protection info event of a particular origin and system_id will be
// stuck to the output pad of the sending element.
func NewProtectionEvent(systemID string, buffer *Buffer, origin string) *Event {
	cSystemID := C.CString(systemID)
	cOrigin := C.CString(origin)
	defer C.free(unsafe.Pointer(cSystemID))
	defer C.free(unsafe.Pointer(cOrigin))
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_protection(
		(*C.gchar)(unsafe.Pointer(cSystemID)),
		buffer.Instance(),
		(*C.gchar)(unsafe.Pointer(cOrigin)),
	)))
}

// NewQOSEvent allocates a new qos event with the given values. The QOS event is generated in an element that wants an upstream element to either reduce or increase its rate because of
// high/low CPU load or other resource usage such as network performance or throttling. Typically sinks generate these events for each buffer they receive.
//
// Type indicates the reason for the QoS event. GST_QOS_TYPE_OVERFLOW is used when a buffer arrived in time or when the sink cannot keep up with the upstream datarate. GST_QOS_TYPE_UNDERFLOW
// is when the sink is not receiving buffers fast enough and thus has to drop late buffers. GST_QOS_TYPE_THROTTLE is used when the datarate is artificially limited by the application, for
// example to reduce power consumption.
//
// Proportion indicates the real-time performance of the streaming in the element that generated the QoS event (usually the sink). The value is generally computed based on more long term
// statistics about the streams timestamps compared to the clock. A value < 1.0 indicates that the upstream element is producing data faster than real-time. A value > 1.0 indicates that the
// upstream element is not producing data fast enough. 1.0 is the ideal proportion value. The proportion value can safely be used to lower or increase the quality of the element.
//
// Diff is the difference against the clock in running time of the last buffer that caused the element to generate the QOS event. A negative value means that the buffer with timestamp arrived
// in time. A positive value indicates how late the buffer with timestamp was. When throttling is enabled, diff will be set to the requested throttling interval.
//
// Timestamp is the timestamp of the last buffer that cause the element to generate the QOS event. It is expressed in running time and thus an ever increasing value.
//
// The upstream element can use the diff and timestamp values to decide whether to process more buffers. For positive diff, all buffers with timestamp <= timestamp + diff will certainly arrive
// late in the sink as well. A (negative) diff value so that timestamp + diff would yield a result smaller than 0 is not allowed.
//
// The application can use general event probes to intercept the QoS event and implement custom application specific QoS handling.
func NewQOSEvent(qType QOSType, proportion float64, diff ClockTimeDiff, timestamp time.Duration) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_qos(
		C.GstQOSType(qType),
		C.gdouble(proportion),
		C.GstClockTimeDiff(diff),
		C.GstClockTime(timestamp.Nanoseconds()),
	)))
}

// NewReconfigureEvent creates a new reconfigure event. The purpose of the reconfigure event is to travel upstream and make elements renegotiate their caps or reconfigure their buffer pools.
// This is useful when changing properties on elements or changing the topology of the pipeline.
func NewReconfigureEvent() *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_reconfigure()))
}

// NewSeekEvent allocates a new seek event with the given parameters.
//
// The seek event configures playback of the pipeline between start to stop at the speed given in rate, also called a playback segment. The start and stop values are expressed in format.
//
// A rate of 1.0 means normal playback rate, 2.0 means double speed. Negatives values means backwards playback. A value of 0.0 for the rate is not allowed and should be accomplished instead
// by PAUSING the pipeline.
//
// A pipeline has a default playback segment configured with a start position of 0, a stop position of -1 and a rate of 1.0. The currently configured playback segment can be queried with
// GST_QUERY_SEGMENT.
//
// start_type and stop_type specify how to adjust the currently configured start and stop fields in playback segment. Adjustments can be made relative or absolute to the last configured values.
// A type of GST_SEEK_TYPE_NONE means that the position should not be updated.
//
// When the rate is positive and start has been updated, playback will start from the newly configured start position.
//
// For negative rates, playback will start from the newly configured stop position (if any). If the stop position is updated, it must be different from -1 (#GST_CLOCK_TIME_NONE) for negative rates.
//
// It is not possible to seek relative to the current playback position, to do this, PAUSE the pipeline, query the current playback position with GST_QUERY_POSITION and update the playback segment
// current position with a GST_SEEK_TYPE_SET to the desired position.
func NewSeekEvent(rate float64, format Format, flags SeekFlags, startType SeekType, start int64, stopType SeekType, stop int64) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_seek(
		C.gdouble(rate),
		C.GstFormat(format),
		C.GstSeekFlags(flags),
		C.GstSeekType(startType),
		C.gint64(start),
		C.GstSeekType(stopType),
		C.gint64(stop),
	)))
}

// NewSegmentEvent creates a new SEGMENT event for segment. The segment event can only travel downstream synchronized with the buffer flow and contains timing information and playback properties
// for the buffers that will follow.
//
// The segment event marks the range of buffers to be processed. All data not within the segment range is not to be processed. This can be used intelligently by plugins to apply more efficient
// methods of skipping unneeded data. The valid range is expressed with the start and stop values.
//
// The time value of the segment is used in conjunction with the start value to convert the buffer timestamps into the stream time. This is usually done in sinks to report the current stream_time.
// time represents the stream_time of a buffer carrying a timestamp of start. time cannot be -1.
//
// start cannot be -1, stop can be -1. If there is a valid stop given, it must be greater or equal the start, including when the indicated playback rate is < 0.
//
// The applied_rate value provides information about any rate adjustment that has already been made to the timestamps and content on the buffers of the stream. (@rate * applied_rate) should always
// equal the rate that has been requested for playback. For example, if an element has an input segment with intended playback rate of 2.0 and applied_rate of 1.0, it can adjust incoming timestamps
// and buffer content by half and output a segment event with rate of 1.0 and applied_rate of 2.0
//
// After a segment event, the buffer stream time is calculated with:
//
// time + (TIMESTAMP(buf) - start) * ABS (rate * applied_rate)
func NewSegmentEvent(segment *Segment) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_segment(
		segment.Instance(),
	)))
}

// NewSegmentDoneEvent creates a new segment-done event. This event is sent by elements that finish playback of a segment as a result of a segment seek.
func NewSegmentDoneEvent(format Format, position int64) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_segment_done(
		C.GstFormat(format), C.gint64(position),
	)))
}

// NewSelectStreamsEvent allocates a new select-streams event.
//
// The select-streams event requests the specified streams to be activated.
//
// The list of streams corresponds to the "Stream ID" of each stream to be activated. Those ID can be obtained via the GstStream objects present in EventStreamStart,
// EventStreamCollection or MessageStreamCollection.
//
// Note: The list of streams can not be empty.
func NewSelectStreamsEvent(streams []*Stream) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_select_streams(
		streamSliceToGlist(streams),
	)))
}

// NewSinkMessageEvent creates a new sink-message event. The purpose of the sink-message event is to instruct a sink to post the message contained in the event
// synchronized with the stream.
//
// name is used to store multiple sticky events on one pad.
func NewSinkMessageEvent(name string, msg *Message) *Event {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_sink_message(
		(*C.gchar)(unsafe.Pointer(cName)),
		msg.Instance(),
	)))
}

// NewStepEvent createss a new step event. The purpose of the step event is to instruct a sink to skip amount (expressed in format) of media. It can be used to
// implement stepping through the video frame by frame or for doing fast trick modes.
//
// A rate of <= 0.0 is not allowed. Pause the pipeline, for the effect of rate = 0.0 or first reverse the direction of playback using a seek event to get the same
// effect as rate < 0.0.
//
// The flush flag will clear any pending data in the pipeline before starting the step operation.
//
// The intermediate flag instructs the pipeline that this step operation is part of a larger step operation.
func NewStepEvent(format Format, amount uint64, rate float64, flush, intermediate bool) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_step(
		C.GstFormat(format),
		C.guint64(amount),
		C.gdouble(rate),
		gboolean(flush),
		gboolean(intermediate),
	)))
}

// NewStreamCollectionEvent creates a new STREAM_COLLECTION event. The stream collection event can only travel downstream synchronized with the buffer flow.
//
// Source elements, demuxers and other elements that manage collections of streams and post GstStreamCollection messages on the bus also send this event downstream
// on each pad involved in the collection, so that activation of a new collection can be tracked through the downstream data flow.
func NewStreamCollectionEvent(collection *StreamCollection) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_stream_collection(
		collection.Instance(),
	)))
}

// NewStreamGroupDoneEvent creates a new Stream Group Done event. The stream-group-done event can only travel downstream synchronized with the buffer flow. Elements
// that receive the event on a pad should handle it mostly like EOS, and emit any data or pending buffers that would depend on more data arriving and unblock, since
// there won't be any more data.
//
// This event is followed by EOS at some point in the future, and is generally used when switching pads - to unblock downstream so that new pads can be exposed before
// sending EOS on the existing pads.
func NewStreamGroupDoneEvent(groupID uint) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_stream_group_done(C.guint(groupID))))
}

// NewStreamStartEvent creates a new STREAM_START event. The stream start event can only travel downstream synchronized with the buffer flow. It is expected to be the
// first event that is sent for a new stream.
//
// Source elements, demuxers and other elements that create new streams are supposed to send this event as the first event of a new stream. It should not be sent after a
// flushing seek or in similar situations and is used to mark the beginning of a new logical stream. Elements combining multiple streams must ensure that this event is only
// forwarded downstream once and not for every single input stream.
//
// The stream_id should be a unique string that consists of the upstream stream-id, / as separator and a unique stream-id for this specific stream. A new stream-id should
// only be created for a stream if the upstream stream is split into (potentially) multiple new streams, e.g. in a demuxer, but not for every single element in the pipeline.
// Pad CreateStreamID can be used to create a stream-id. There are no particular semantics for the stream-id, though it should be deterministic (to support stream matching)
// and it might be used to order streams (besides any information conveyed by stream flags).
func NewStreamStartEvent(streamID string) *Event {
	cName := C.CString(streamID)
	defer C.free(unsafe.Pointer(cName))
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_stream_start(
		(*C.gchar)(unsafe.Pointer(cName)),
	)))
}

// NewTagEvent generates a metadata tag event from the given taglist.
//
// The scope of the taglist specifies if the taglist applies to the complete medium or only to this specific stream. As the tag event is a sticky event, elements should merge
// tags received from upstream with a given scope with their own tags with the same scope and create a new tag event from it.
func NewTagEvent(tagList *TagList) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_tag(
		tagList.Ref().Instance(),
	)))
}

// NewTOCEvent generates a TOC event from the given toc. The purpose of the TOC event is to inform elements that some kind of the TOC was found.
func NewTOCEvent(toc *TOC, updated bool) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_toc(
		toc.Instance(), gboolean(updated),
	)))
}

// NewTOCSelectEvent generates a TOC select event with the given uid. The purpose of the TOC select event is to start playback based on the TOC's
// entry with the given uid.
func NewTOCSelectEvent(uid string) *Event {
	cUID := C.CString(uid)
	defer C.free(unsafe.Pointer(cUID))
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_toc_select(
		(*C.gchar)(unsafe.Pointer(cUID)),
	)))
}

// NewCustomEvent creates a new custom-typed event. This can be used for anything not handled by other event-specific functions to pass an event
// to another element.
func NewCustomEvent(eventType EventType, structure *Structure) *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_new_custom(
		C.GstEventType(eventType), structure.Instance(),
	)))
}
