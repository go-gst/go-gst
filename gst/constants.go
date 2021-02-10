package gst

// #include "gst.go.h"
import "C"

import (
	"time"
	"unsafe"
)

// Version represents information about the current GST version.
type Version int

const (
	// VersionMajor is the major version number of the GStreamer core.
	VersionMajor Version = C.GST_VERSION_MAJOR
	// VersionMinor is the minor version number of the GStreamer core.
	VersionMinor Version = C.GST_VERSION_MINOR
)

// License represents a type of license used on a plugin.
type License string

// Types of licenses
const (
	LicenseLGPL        License = "LGPL"
	LicenseGPL         License = "GPL"
	LicenseQPL         License = "QPL"
	LicenseGPLQPL      License = "GPL/QPL"
	LicenseMPL         License = "MPL"
	LicenseBSD         License = "BSD"
	LicenseMIT         License = "MIT/X11"
	LicenseProprietary License = "Proprietary"
	LicenseUnknown     License = "unknown"
)

// ClockTimeDiff is a datatype to hold a time difference, measured in nanoseconds.
type ClockTimeDiff int64

// ClockTimeNone means infinite timeout or an empty value
var ClockTimeNone time.Duration = time.Duration(-1)

// BufferOffsetNone is a var for no-offset return results.
var BufferOffsetNone time.Duration = time.Duration(-1)

var (
	// ClockTimeNone means infinite timeout (unsigned representation of -1) or an otherwise unknown value.
	gstClockTimeNone C.GstClockTime = 0xffffffffffffffff
	// // BufferOffsetNone is a constant for no-offset return results.
	// gstBufferOffsetNone C.GstClockTime = 0xffffffffffffffff
)

// ClockEntryType wraps GstClockEntryType
type ClockEntryType int

// Type castings of ClockEntryTypes
const (
	ClockEntrySingle   ClockEntryType = C.GST_CLOCK_ENTRY_SINGLE   // (0) – a single shot timeout
	ClockEntryPeriodic ClockEntryType = C.GST_CLOCK_ENTRY_PERIODIC // (1) – a periodic timeout request
)

// ClockFlags wraps GstClockFlags
type ClockFlags int

// Type castings of ClockFlags
const (
	ClockFlagCanDoSingleSync    ClockFlags = C.GST_CLOCK_FLAG_CAN_DO_SINGLE_SYNC    // (16) – clock can do a single sync timeout request
	ClockFlagCanDoSingleAsync   ClockFlags = C.GST_CLOCK_FLAG_CAN_DO_SINGLE_ASYNC   // (32) – clock can do a single async timeout request
	ClockFlagCanDoPeriodicSync  ClockFlags = C.GST_CLOCK_FLAG_CAN_DO_PERIODIC_SYNC  // (64) – clock can do sync periodic timeout requests
	ClockFlagCanDoPeriodicAsync ClockFlags = C.GST_CLOCK_FLAG_CAN_DO_PERIODIC_ASYNC // (128) – clock can do async periodic timeout callbacks
	ClockFlagCanSetResolution   ClockFlags = C.GST_CLOCK_FLAG_CAN_SET_RESOLUTION    // (256) – clock's resolution can be changed
	ClockFlagCanSetMaster       ClockFlags = C.GST_CLOCK_FLAG_CAN_SET_MASTER        // (512) – clock can be slaved to a master clock
	ClockFlagNeedsStartupSync   ClockFlags = C.GST_CLOCK_FLAG_NEEDS_STARTUP_SYNC    // (1024) – clock needs to be synced before it can be used (Since: 1.6)
	ClockFlagLast               ClockFlags = C.GST_CLOCK_FLAG_LAST                  // (4096) – subclasses can add additional flags starting from this flag
)

// ClockReturn wraps a GstClockReturn
type ClockReturn int

// Type castings of clock returns
const (
	ClockOK          ClockReturn = C.GST_CLOCK_OK          // (0) – The operation succeeded.
	ClockEarly       ClockReturn = C.GST_CLOCK_EARLY       // (1) – The operation was scheduled too late.
	ClockUnscheduled ClockReturn = C.GST_CLOCK_UNSCHEDULED // (2) – The clockID was unscheduled
	ClockBusy        ClockReturn = C.GST_CLOCK_BUSY        // (3) – The ClockID is busy
	ClockBadTime     ClockReturn = C.GST_CLOCK_BADTIME     // (4) – A bad time was provided to a function.
	ClockError       ClockReturn = C.GST_CLOCK_ERROR       // (5) – An error occurred
	ClockUnsupported ClockReturn = C.GST_CLOCK_UNSUPPORTED // (6) – Operation is not supported
	ClockDone        ClockReturn = C.GST_CLOCK_DONE        // (7) – The ClockID is done waiting
)

// BusSyncReply casts GstBusSyncReply to a go type
type BusSyncReply int

// Type castings of SyncReplies
const (
	BusDrop  BusSyncReply = C.GST_BUS_DROP  // (0) – drop the message
	BusPass  BusSyncReply = C.GST_BUS_PASS  // (1) – pass the message to the async queue
	BusAsync BusSyncReply = C.GST_BUS_ASYNC // (2) – pass message to async queue, continue if message is handled
)

// BufferFlags casts GstBufferFlags to a go type.
type BufferFlags int

// Type castings of BufferFlags
const (
	BufferFlagLive         BufferFlags = C.GST_BUFFER_FLAG_LIVE          // (16) – the buffer is live data and should be discarded in the PAUSED state.
	BufferFlagDecodeOnly   BufferFlags = C.GST_BUFFER_FLAG_DECODE_ONLY   // (32) – the buffer contains data that should be dropped because it will be clipped against the segment boundaries or because it does not contain data that should be shown to the user.
	BufferFlagDiscont      BufferFlags = C.GST_BUFFER_FLAG_DISCONT       // (64) – the buffer marks a data discontinuity in the stream. This typically occurs after a seek or a dropped buffer from a live or network source.
	BufferFlagResync       BufferFlags = C.GST_BUFFER_FLAG_RESYNC        // (128) – the buffer timestamps might have a discontinuity and this buffer is a good point to resynchronize.
	BufferFlagCorrupted    BufferFlags = C.GST_BUFFER_FLAG_CORRUPTED     // (256) – the buffer data is corrupted.
	BufferFlagMarker       BufferFlags = C.GST_BUFFER_FLAG_MARKER        // (512) – the buffer contains a media specific marker. for video this is the end of a frame boundary, for audio this is the start of a talkspurt.
	BufferFlagHeader       BufferFlags = C.GST_BUFFER_FLAG_HEADER        // (1024) – the buffer contains header information that is needed to decode the following data.
	BufferFlagGap          BufferFlags = C.GST_BUFFER_FLAG_GAP           // (2048) – the buffer has been created to fill a gap in the stream and contains media neutral data (elements can switch to optimized code path that ignores the buffer content).
	BufferFlagDroppable    BufferFlags = C.GST_BUFFER_FLAG_DROPPABLE     // (4096) – the buffer can be dropped without breaking the stream, for example to reduce bandwidth.
	BufferFlagDeltaUnit    BufferFlags = C.GST_BUFFER_FLAG_DELTA_UNIT    // (8192) – this unit cannot be decoded independently.
	BufferFlagSyncAfter    BufferFlags = C.GST_BUFFER_FLAG_SYNC_AFTER    // (32768) – Elements which write to disk or permanent storage should ensure the data is synced after writing the contents of this buffer. (Since: 1.6)
	BufferFlagNonDroppable BufferFlags = C.GST_BUFFER_FLAG_NON_DROPPABLE // (65536) – This buffer is important and should not be dropped. This can be used to mark important buffers, e.g. to flag RTP packets carrying keyframes or codec setup data for RTP Forward Error Correction purposes, or to prevent still video frames from being dropped by elements due to QoS. (Since: 1.14)
	BufferFlagLast         BufferFlags = C.GST_BUFFER_FLAG_LAST          // (1048576) – additional media specific flags can be added starting from this flag.
)

// BufferCopyFlags casts GstBufferCopyFlags to a go type.
type BufferCopyFlags int

// Type castings of BufferCopyFlags
const (
	BufferCopyNone        BufferCopyFlags = C.GST_BUFFER_COPY_NONE       // (0) – copy nothing
	BufferCopyBufferFlags BufferCopyFlags = C.GST_BUFFER_COPY_FLAGS      // (1) – flag indicating that buffer flags should be copied
	BufferCopyTimestamps  BufferCopyFlags = C.GST_BUFFER_COPY_TIMESTAMPS // (2) – flag indicating that buffer pts, dts, duration, offset and offset_end should be copied
	BufferCopyMeta        BufferCopyFlags = C.GST_BUFFER_COPY_META       // (4) – flag indicating that buffer meta should be copied
	BufferCopyMemory      BufferCopyFlags = C.GST_BUFFER_COPY_MEMORY     // (8) – flag indicating that buffer memory should be reffed and appended to already existing memory. Unless the memory is marked as NO_SHARE, no actual copy of the memory is made but it is simply reffed. Add GST_BUFFER_COPY_DEEP to force a real copy.
	BufferCopyMerge       BufferCopyFlags = C.GST_BUFFER_COPY_MERGE      // (16) – flag indicating that buffer memory should be merged
	BufferCopyDeep        BufferCopyFlags = C.GST_BUFFER_COPY_DEEP       // (32) – flag indicating that memory should always be copied instead of reffed (Since: 1.2)
)

// BufferPoolAcquireFlags casts GstBufferPoolAcquireFlags to a go type.
type BufferPoolAcquireFlags int

// Type castings of BufferPoolAcquireFlags
const (
	BufferPoolAcquireFlagNone     BufferPoolAcquireFlags = C.GST_BUFFER_POOL_ACQUIRE_FLAG_NONE     // (0) – no flags
	BufferPoolAcquireFlagKeyUnit  BufferPoolAcquireFlags = C.GST_BUFFER_POOL_ACQUIRE_FLAG_KEY_UNIT // (1) – buffer is keyframe
	BufferPoolAcquireFlagDontWait BufferPoolAcquireFlags = C.GST_BUFFER_POOL_ACQUIRE_FLAG_DONTWAIT // (2) – when the bufferpool is empty, acquire_buffer will by default block until a buffer is released into the pool again. Setting this flag makes acquire_buffer return GST_FLOW_EOS instead of blocking.
	BufferPoolAcquireFlagDiscont  BufferPoolAcquireFlags = C.GST_BUFFER_POOL_ACQUIRE_FLAG_DISCONT  // (4) – buffer is discont
	BufferPoolAcquireFlagLast     BufferPoolAcquireFlags = C.GST_BUFFER_POOL_ACQUIRE_FLAG_LAST     // (65536) – last flag, subclasses can use private flags starting from this value.
)

// BufferingMode is a representation of GstBufferingMode
type BufferingMode int

// Type casts of buffering modes
const (
	BufferingStream    BufferingMode = C.GST_BUFFERING_STREAM    // (0) – a small amount of data is buffered
	BufferingDownload  BufferingMode = C.GST_BUFFERING_DOWNLOAD  // (1) – the stream is being downloaded
	BufferingTimeshift BufferingMode = C.GST_BUFFERING_TIMESHIFT //  (2) – the stream is being downloaded in a ringbuffer
	BufferingLive      BufferingMode = C.GST_BUFFERING_LIVE      // (3) – the stream is a live stream
)

// String implements a stringer on a BufferingMode.
func (b BufferingMode) String() string {
	switch b {
	case BufferingStream:
		return "A small amount of data is buffered"
	case BufferingDownload:
		return "The stream is being downloaded"
	case BufferingTimeshift:
		return "The stream is being downloaded in a ringbuffer"
	case BufferingLive:
		return "The stream is live"
	}
	return ""
}

// SegmentFlags casts GstSegmentFlags
type SegmentFlags int

// Type castings
const (
	SegmentFlagNone      SegmentFlags = C.GST_SEGMENT_FLAG_NONE                // (0) – no flags
	SegmentFlagReset     SegmentFlags = C.GST_SEGMENT_FLAG_RESET               // (1) – reset the pipeline running_time to the segment running_time
	SegmentFlagTrickMode SegmentFlags = C.GST_SEGMENT_FLAG_TRICKMODE           // (16) – perform skip playback (Since: 1.6)
	SegmentFlagSkip      SegmentFlags = C.GST_SEGMENT_FLAG_SKIP                // (16) – Deprecated backward compatibility flag, replaced by GST_SEGMENT_FLAG_TRICKMODE
	SegmentFlagSegment   SegmentFlags = C.GST_SEGMENT_FLAG_SEGMENT             // (8) – send SEGMENT_DONE instead of EOS
	SegmentFlagKeyUnits  SegmentFlags = C.GST_SEGMENT_FLAG_TRICKMODE_KEY_UNITS // (128) – Decode only keyframes, where possible (Since: 1.6)
	// SegmentFlagTrickModeForwardPredicted SegmentFlags = C.GST_SEGMENT_FLAG_TRICKMODE_FORWARD_PREDICTED // (512) – Decode only keyframes or forward predicted frames, where possible (Since: 1.18)
	SegmentFlagTrickModeNoAudio SegmentFlags = C.GST_SEGMENT_FLAG_TRICKMODE_NO_AUDIO // (256) – Do not decode any audio, where possible (Since: 1.6)
)

// EventType is a go cast for a GstEventType
type EventType int

// Type casts for EventTypes
const (
	EventTypeUnknown          EventType = C.GST_EVENT_UNKNOWN           //(0) – unknown event.
	EventTypeFlushStart       EventType = C.GST_EVENT_FLUSH_START       // (2563) – Start a flush operation. This event clears all data from the pipeline and unblock all streaming threads.
	EventTypeFlushStop        EventType = C.GST_EVENT_FLUSH_STOP        // (5127) – Stop a flush operation. This event resets the running-time of the pipeline.
	EventTypeStreamStart      EventType = C.GST_EVENT_STREAM_START      // (10254) – Event to mark the start of a new stream. Sent before any other serialized event and only sent at the start of a new stream, not after flushing seeks.
	EventTypeCaps             EventType = C.GST_EVENT_CAPS              // (12814) – GstCaps event. Notify the pad of a new media type.
	EventTypeSegment          EventType = C.GST_EVENT_SEGMENT           // (17934) – A new media segment follows in the dataflow. The segment events contains information for clipping buffers and converting buffer timestamps to running-time and stream-time.
	EventTypeStreamCollection EventType = C.GST_EVENT_STREAM_COLLECTION // (19230) – A new GstStreamCollection is available (Since: 1.10)
	EventTypeTag              EventType = C.GST_EVENT_TAG               // (20510) – A new set of metadata tags has been found in the stream.
	EventTypeBufferSize       EventType = C.GST_EVENT_BUFFERSIZE        // (23054) – Notification of buffering requirements. Currently not used yet.
	EventTypeSinkMessage      EventType = C.GST_EVENT_SINK_MESSAGE      // (25630) – An event that sinks turn into a message. Used to send messages that should be emitted in sync with rendering.
	EventTypeStreamGroupDone  EventType = C.GST_EVENT_STREAM_GROUP_DONE // (26894) – Indicates that there is no more data for the stream group ID in the message. Sent before EOS in some instances and should be handled mostly the same. (Since: 1.10)
	EventTypeEOS              EventType = C.GST_EVENT_EOS               // (28174) – End-Of-Stream. No more data is to be expected to follow without either a STREAM_START event, or a FLUSH_STOP and a SEGMENT event.
	EventTypeTOC              EventType = C.GST_EVENT_TOC               // (30750) – An event which indicates that a new table of contents (TOC) was found or updated.
	EventTypeProtection       EventType = C.GST_EVENT_PROTECTION        // (33310) – An event which indicates that new or updated encryption information has been found in the stream.
	EventTypeSegmentDone      EventType = C.GST_EVENT_SEGMENT_DONE      // (38406) – Marks the end of a segment playback.
	EventTypeGap              EventType = C.GST_EVENT_GAP               // (40966) – Marks a gap in the datastream.
	// EventTypeInstantRateChange      EventType = C.GST_EVENT_INSTANT_RATE_CHANGE      // (46090) – Notify downstream that a playback rate override should be applied as soon as possible. (Since: 1.18)
	EventTypeQOS           EventType = C.GST_EVENT_QOS            // (48641) – A quality message. Used to indicate to upstream elements that the downstream elements should adjust their processing rate.
	EventTypeSeek          EventType = C.GST_EVENT_SEEK           // (51201) – A request for a new playback position and rate.
	EventTypeNavigation    EventType = C.GST_EVENT_NAVIGATION     // (53761) – Navigation events are usually used for communicating user requests, such as mouse or keyboard movements, to upstream elements.
	EventTypeLatency       EventType = C.GST_EVENT_LATENCY        // (56321) – Notification of new latency adjustment. Sinks will use the latency information to adjust their synchronisation.
	EventTypeStep          EventType = C.GST_EVENT_STEP           // (58881) – A request for stepping through the media. Sinks will usually execute the step operation.
	EventTypeReconfigure   EventType = C.GST_EVENT_RECONFIGURE    // (61441) – A request for upstream renegotiating caps and reconfiguring.
	EventTypeTOCSelect     EventType = C.GST_EVENT_TOC_SELECT     // (64001) – A request for a new playback position based on TOC entry's UID.
	EventTypeSelectStreams EventType = C.GST_EVENT_SELECT_STREAMS // (66561) – A request to select one or more streams (Since: 1.10)
	// EventTypeInstantRateSyncTime    EventType = C.GST_EVENT_INSTANT_RATE_SYNC_TIME   // (66817) – Sent by the pipeline to notify elements that handle the instant-rate-change event about the running-time when the rate multiplier should be applied (or was applied). (Since: 1.18)
	EventTypeCustomUpstream         EventType = C.GST_EVENT_CUSTOM_UPSTREAM          // (69121) – Upstream custom event
	EventTypeCustomDownstream       EventType = C.GST_EVENT_CUSTOM_DOWNSTREAM        // (71686) – Downstream custom event that travels in the data flow.
	EventTypeCustomOOB              EventType = C.GST_EVENT_CUSTOM_DOWNSTREAM_OOB    // (74242) – Custom out-of-band downstream event.
	EventTypeCustomDownstreamSticky EventType = C.GST_EVENT_CUSTOM_DOWNSTREAM_STICKY // (76830) – Custom sticky downstream event.
	EventTypeCustomBoth             EventType = C.GST_EVENT_CUSTOM_BOTH              // (79367) – Custom upstream or downstream event. In-band when travelling downstream.
	EventTypeCustomBothOOB          EventType = C.GST_EVENT_CUSTOM_BOTH_OOB          // (81923) – Custom upstream or downstream out-of-band event.
)

// String implements a stringer on event types
func (e EventType) String() string { return C.GoString(C.gst_event_type_get_name(C.GstEventType(e))) }

// EventTypeFlags casts GstEventTypeFlags
type EventTypeFlags int

// Type castings
const (
	EventTypeFlagUpstream    EventTypeFlags = C.GST_EVENT_TYPE_UPSTREAM     // (1) – Set if the event can travel upstream.
	EventTypeFlagDownstream  EventTypeFlags = C.GST_EVENT_TYPE_DOWNSTREAM   // (2) – Set if the event can travel downstream.
	EventTypeFlagSerialized  EventTypeFlags = C.GST_EVENT_TYPE_SERIALIZED   // (4) – Set if the event should be serialized with data flow.
	EventTypeFlagSticky      EventTypeFlags = C.GST_EVENT_TYPE_STICKY       // (8) – Set if the event is sticky on the pads.
	EventTypeFlagStickyMulti EventTypeFlags = C.GST_EVENT_TYPE_STICKY_MULTI // (16) – Multiple sticky events can be on a pad, each identified by the event name.
)

// GapFlags casts GstGapFlags
type GapFlags int

// Type castings
const (
	GapFlagMissingData GapFlags = 1
)

// QOSType casts GstQOSType
type QOSType int

// Type castings
const (
	QOSTypeOverflow  QOSType = C.GST_QOS_TYPE_OVERFLOW  // (0) – The QoS event type that is produced when upstream elements are producing data too quickly and the element can't keep up processing the data. Upstream should reduce their production rate. This type is also used when buffers arrive early or in time.
	QOSTypeUnderflow QOSType = C.GST_QOS_TYPE_UNDERFLOW // (1) – The QoS event type that is produced when upstream elements are producing data too slowly and need to speed up their production rate.
	QOSTypeThrottle  QOSType = C.GST_QOS_TYPE_THROTTLE  // (2) – The QoS event type that is produced when the application enabled throttling to limit the data rate.
)

// Format is a representation of GstFormat.
type Format int

// Type casts of formats
const (
	FormatUndefined Format = C.GST_FORMAT_UNDEFINED // (0) – undefined format
	FormatDefault   Format = C.GST_FORMAT_DEFAULT   // (1) – the default format of the pad/element. This can be samples for raw audio, or frames/fields for raw video.
	FormatBytes     Format = C.GST_FORMAT_BYTES     // (2) - bytes
	FormatTime      Format = C.GST_FORMAT_TIME      // (3) – time in nanoseconds
)

// String implements a stringer on GstFormat types
func (f Format) String() string {
	switch f {
	case FormatUndefined:
		return "undefined"
	case FormatDefault:
		return "default"
	case FormatBytes:
		return "bytes"
	case FormatTime:
		return "time"
	}
	return ""
}

// MessageType is an alias to the C equivalent of GstMessageType.
// See the official documentation for definitions of the messages:
// https://gstreamer.freedesktop.org/documentation/gstreamer/gstmessage.html?gi-language=c#GstMessageType
type MessageType int

// Type casting of GstMessageTypes
// See the official documentation for definitions of the messages:
// https://gstreamer.freedesktop.org/documentation/gstreamer/gstmessage.html?gi-language=c#GstMessageType
const (
	MessageUnknown          MessageType = C.GST_MESSAGE_UNKNOWN
	MessageEOS              MessageType = C.GST_MESSAGE_EOS
	MessageError            MessageType = C.GST_MESSAGE_ERROR
	MessageWarning          MessageType = C.GST_MESSAGE_WARNING
	MessageInfo             MessageType = C.GST_MESSAGE_INFO
	MessageTag              MessageType = C.GST_MESSAGE_TAG
	MessageBuffering        MessageType = C.GST_MESSAGE_BUFFERING
	MessageStateChanged     MessageType = C.GST_MESSAGE_STATE_CHANGED
	MessageStateDirty       MessageType = C.GST_MESSAGE_STATE_DIRTY
	MessageStepDone         MessageType = C.GST_MESSAGE_STEP_DONE
	MessageClockProvide     MessageType = C.GST_MESSAGE_CLOCK_PROVIDE
	MessageClockLost        MessageType = C.GST_MESSAGE_CLOCK_LOST
	MessageNewClock         MessageType = C.GST_MESSAGE_NEW_CLOCK
	MessageStructureChange  MessageType = C.GST_MESSAGE_STRUCTURE_CHANGE
	MessageStreamStatus     MessageType = C.GST_MESSAGE_STREAM_STATUS
	MessageApplication      MessageType = C.GST_MESSAGE_APPLICATION
	MessageElement          MessageType = C.GST_MESSAGE_ELEMENT
	MessageSegmentStart     MessageType = C.GST_MESSAGE_SEGMENT_START
	MessageSegmentDone      MessageType = C.GST_MESSAGE_SEGMENT_DONE
	MessageDurationChanged  MessageType = C.GST_MESSAGE_DURATION_CHANGED
	MessageLatency          MessageType = C.GST_MESSAGE_LATENCY
	MessageAsyncStart       MessageType = C.GST_MESSAGE_ASYNC_START
	MessageAsyncDone        MessageType = C.GST_MESSAGE_ASYNC_DONE
	MessageRequestState     MessageType = C.GST_MESSAGE_REQUEST_STATE
	MessageStepStart        MessageType = C.GST_MESSAGE_STEP_START
	MessageQoS              MessageType = C.GST_MESSAGE_QOS
	MessageProgress         MessageType = C.GST_MESSAGE_PROGRESS
	MessageTOC              MessageType = C.GST_MESSAGE_TOC
	MessageResetTime        MessageType = C.GST_MESSAGE_RESET_TIME
	MessageStreamStart      MessageType = C.GST_MESSAGE_STREAM_START
	MessageNeedContext      MessageType = C.GST_MESSAGE_NEED_CONTEXT
	MessageHaveContext      MessageType = C.GST_MESSAGE_HAVE_CONTEXT
	MessageExtended         MessageType = C.GST_MESSAGE_EXTENDED
	MessageDeviceAdded      MessageType = C.GST_MESSAGE_DEVICE_ADDED
	MessageDeviceRemoved    MessageType = C.GST_MESSAGE_DEVICE_REMOVED
	MessagePropertyNotify   MessageType = C.GST_MESSAGE_PROPERTY_NOTIFY
	MessageStreamCollection MessageType = C.GST_MESSAGE_STREAM_COLLECTION
	MessageStreamsSelected  MessageType = C.GST_MESSAGE_STREAMS_SELECTED
	MessageRedirect         MessageType = C.GST_MESSAGE_REDIRECT
	MessageDeviceChanged    MessageType = C.GST_MESSAGE_DEVICE_CHANGED
	MessageAny              MessageType = C.GST_MESSAGE_ANY
)

// String implements a stringer on MessageTypes
func (m MessageType) String() string {
	return C.GoString(C.gst_message_type_get_name((C.GstMessageType)(m)))
}

// PadFlags is a go cast of GstPadFlags
type PadFlags int

// Type casts of PadFlags
const (
	PadFlagBlocked         PadFlags = C.GST_PAD_FLAG_BLOCKED          // (16) – is dataflow on a pad blocked
	PadFlagFlushing        PadFlags = C.GST_PAD_FLAG_FLUSHING         // (32) – is pad flushing
	PadFlagEOS             PadFlags = C.GST_PAD_FLAG_EOS              // (64) – is pad in EOS state
	PadFlagBlocking        PadFlags = C.GST_PAD_FLAG_BLOCKING         // (128) – is pad currently blocking on a buffer or event
	PadFlagParent          PadFlags = C.GST_PAD_FLAG_NEED_PARENT      // (256) – ensure that there is a parent object before calling into the pad callbacks.
	PadFlagReconfigure     PadFlags = C.GST_PAD_FLAG_NEED_RECONFIGURE // (512) – the pad should be reconfigured/renegotiated. The flag has to be unset manually after reconfiguration happened.
	PadFlagPendingEvents   PadFlags = C.GST_PAD_FLAG_PENDING_EVENTS   // (1024) – the pad has pending events
	PadFlagFixedCaps       PadFlags = C.GST_PAD_FLAG_FIXED_CAPS       // (2048) – the pad is using fixed caps. This means that once the caps are set on the pad, the default caps query function will only return those caps.
	PadFlagProxyCaps       PadFlags = C.GST_PAD_FLAG_PROXY_CAPS       // (4096) – the default event and query handler will forward all events and queries to the internally linked pads instead of discarding them.
	PadFlagProxyAllocation PadFlags = C.GST_PAD_FLAG_PROXY_ALLOCATION // (8192) – the default query handler will forward allocation queries to the internally linked pads instead of discarding them.
	PadFlagProxyScheduling PadFlags = C.GST_PAD_FLAG_PROXY_SCHEDULING // (16384) – the default query handler will forward scheduling queries to the internally linked pads instead of discarding them.
	PadFlagAcceptIntersect PadFlags = C.GST_PAD_FLAG_ACCEPT_INTERSECT // (32768) – the default accept-caps handler will check it the caps intersect the query-caps result instead of checking for a subset. This is interesting for parsers that can accept incompletely specified caps.
	PadFlagAcceptTemplate  PadFlags = C.GST_PAD_FLAG_ACCEPT_TEMPLATE  // (65536) – the default accept-caps handler will use the template pad caps instead of query caps to compare with the accept caps. Use this in combination with GST_PAD_FLAG_ACCEPT_INTERSECT. (Since: 1.6)
	PadFlagLast            PadFlags = C.GST_PAD_FLAG_LAST             // (1048576) – offset to define more flags
)

// PadLinkCheck is a go cast of GstPadLinkCheck
type PadLinkCheck int

// Type casts of PadLinkChecks
const (
	PadLinkCheckNothing       PadLinkCheck = C.GST_PAD_LINK_CHECK_NOTHING        // (0) – Don't check hierarchy or caps compatibility.
	PadLinkCheckHierarchy     PadLinkCheck = C.GST_PAD_LINK_CHECK_HIERARCHY      // (1) – Check the pads have same parents/grandparents. Could be omitted if it is already known that the two elements that own the pads are in the same bin.
	PadLinkCheckTemplateCaps  PadLinkCheck = C.GST_PAD_LINK_CHECK_TEMPLATE_CAPS  // (2) – Check if the pads are compatible by using their template caps. This is much faster than GST_PAD_LINK_CHECK_CAPS, but would be unsafe e.g. if one pad has GST_CAPS_ANY.
	PadLinkCheckCaps          PadLinkCheck = C.GST_PAD_LINK_CHECK_CAPS           // (4) – Check if the pads are compatible by comparing the caps returned by gst_pad_query_caps.
	PadLinkCheckNoReconfigure PadLinkCheck = C.GST_PAD_LINK_CHECK_NO_RECONFIGURE // (8) – Disables pushing a reconfigure event when pads are linked.
	PadLinkCheckDefault       PadLinkCheck = C.GST_PAD_LINK_CHECK_DEFAULT        // (5) – The default checks done when linking pads (i.e. the ones used by gst_pad_link).
)

// PadMode is a cast of GstPadMode.
type PadMode int

// Type casts of PadModes
const (
	PadModeNone PadMode = C.GST_PAD_MODE_NONE // (0) – Pad will not handle dataflow
	PadModePush PadMode = C.GST_PAD_MODE_PUSH // (1) – Pad handles dataflow in downstream push mode
	PadModePull PadMode = C.GST_PAD_MODE_PULL // (2) – Pad handles dataflow in upstream pull mode
)

// String implements a stringer on PadMode
func (p PadMode) String() string {
	return C.GoString(C.gst_pad_mode_get_name(C.GstPadMode(p)))
}

// PadDirection is a cast of GstPadDirection to a go type.
type PadDirection int

// Type casting of pad directions
const (
	PadDirectionUnknown PadDirection = C.GST_PAD_UNKNOWN // (0) - the direction is unknown
	PadDirectionSource  PadDirection = C.GST_PAD_SRC     // (1) - the pad is a source pad
	PadDirectionSink    PadDirection = C.GST_PAD_SINK    // (2) - the pad is a sink pad
)

// String implements a Stringer on PadDirection.
func (p PadDirection) String() string {
	switch p {
	case PadDirectionUnknown:
		return "unknown"
	case PadDirectionSource:
		return "src"
	case PadDirectionSink:
		return "sink"
	}
	return ""
}

// PadLinkReturn os a representation of GstPadLinkReturn.
type PadLinkReturn int

// Type casts for PadLinkReturns.
const (
	PadLinkOK             PadLinkReturn = C.GST_PAD_LINK_OK
	PadLinkWrongHierarchy PadLinkReturn = C.GST_PAD_LINK_WRONG_HIERARCHY
	PadLinkWasLinked      PadLinkReturn = C.GST_PAD_LINK_WAS_LINKED
	PadLinkWrongDirection PadLinkReturn = C.GST_PAD_LINK_WRONG_DIRECTION
	PadLinkNoFormat       PadLinkReturn = C.GST_PAD_LINK_NOFORMAT
	PadLinkNoSched        PadLinkReturn = C.GST_PAD_LINK_NOSCHED
	PadLinkRefused        PadLinkReturn = C.GST_PAD_LINK_REFUSED
)

// String implemeents a stringer on PadLinkReturn
func (p PadLinkReturn) String() string {
	return C.GoString(C.gst_pad_link_get_name(C.GstPadLinkReturn(p)))
}

// PadPresence is a cast of GstPadPresence to a go type.
type PadPresence int

// Type casting of pad presences
const (
	PadPresenceAlways    PadPresence = C.GST_PAD_ALWAYS    // (0) - the pad is always available
	PadPresenceSometimes PadPresence = C.GST_PAD_SOMETIMES // (1) - the pad will become available depending on the media stream
	PadPresenceRequest   PadPresence = C.GST_PAD_REQUEST   // (2) - the pad is only available on request with gst_element_request_pad.
)

// String implements a stringer on PadPresence.
func (p PadPresence) String() string {
	switch p {
	case PadPresenceAlways:
		return "always"
	case PadPresenceSometimes:
		return "sometimes"
	case PadPresenceRequest:
		return "request"
	}
	return ""
}

// PadProbeReturn casts GstPadProbeReturn
type PadProbeReturn int

// Type castings of ProbeReturns
const (
	PadProbeDrop      PadProbeReturn = C.GST_PAD_PROBE_DROP    // (0) – drop data in data probes. For push mode this means that the data item is not sent downstream. For pull mode, it means that the data item is not passed upstream. In both cases, no other probes are called for this item and GST_FLOW_OK or TRUE is returned to the caller.
	PadProbeOK        PadProbeReturn = C.GST_PAD_PROBE_OK      // (1) – normal probe return value. This leaves the probe in place, and defers decisions about dropping or passing data to other probes, if any. If there are no other probes, the default behaviour for the probe type applies ('block' for blocking probes, and 'pass' for non-blocking probes).
	PadProbeRemove    PadProbeReturn = C.GST_PAD_PROBE_REMOVE  // (2) – remove this probe.
	PadProbePass      PadProbeReturn = C.GST_PAD_PROBE_PASS    // (3) – pass the data item in the block probe and block on the next item.
	PadProbeUnhandled PadProbeReturn = C.GST_PAD_PROBE_HANDLED // (4) – Data has been handled in the probe and will not be forwarded further. For events and buffers this is the same behaviour as GST_PAD_PROBE_DROP (except that in this case you need to unref the buffer or event yourself). For queries it will also return TRUE to the caller. The probe can also modify the GstFlowReturn value by using the GST_PAD_PROBE_INFO_FLOW_RETURN() accessor. Note that the resulting query must contain valid entries. Since: 1.6
)

// PadProbeType casts GstPadProbeType
type PadProbeType int

// Type castings of PadProbeTypes
const (
	PadProbeTypeInvalid         PadProbeType = C.GST_PAD_PROBE_TYPE_INVALID          // (0) – invalid probe type
	PadProbeTypeIdle            PadProbeType = C.GST_PAD_PROBE_TYPE_IDLE             // (1) – probe idle pads and block while the callback is called
	PadProbeTypeBlock           PadProbeType = C.GST_PAD_PROBE_TYPE_BLOCK            // (2) – probe and block pads
	PadProbeTypeBuffer          PadProbeType = C.GST_PAD_PROBE_TYPE_BUFFER           // (16) – probe buffers
	PadProbeTypeBufferList      PadProbeType = C.GST_PAD_PROBE_TYPE_BUFFER_LIST      // (32) – probe buffer lists
	PadProbeTypeEventDownstream PadProbeType = C.GST_PAD_PROBE_TYPE_EVENT_DOWNSTREAM // (64) – probe downstream events
	PadProbeTypeEventUpstream   PadProbeType = C.GST_PAD_PROBE_TYPE_EVENT_UPSTREAM   // (128) – probe upstream events
	PadProbeTypeEventFlush      PadProbeType = C.GST_PAD_PROBE_TYPE_EVENT_FLUSH      // (256) – probe flush events. This probe has to be explicitly enabled and is not included in the @GST_PAD_PROBE_TYPE_EVENT_DOWNSTREAM or @GST_PAD_PROBE_TYPE_EVENT_UPSTREAM probe types.
	PadProbeTypeQueryDownstream PadProbeType = C.GST_PAD_PROBE_TYPE_QUERY_DOWNSTREAM // (512) – probe downstream queries
	PadProbeTypeQueryUpstream   PadProbeType = C.GST_PAD_PROBE_TYPE_QUERY_UPSTREAM   // (1024) – probe upstream queries
	PadProbeTypePush            PadProbeType = C.GST_PAD_PROBE_TYPE_PUSH             // (4096) – probe push
	PadProbeTypePull            PadProbeType = C.GST_PAD_PROBE_TYPE_PULL             // (8192) – probe pull
	PadProbeTypeBlocking        PadProbeType = C.GST_PAD_PROBE_TYPE_BLOCKING         // (3) – probe and block at the next opportunity, at data flow or when idle
	PadProbeTypeDataDownstream  PadProbeType = C.GST_PAD_PROBE_TYPE_DATA_DOWNSTREAM  // (112) – probe downstream data (buffers, buffer lists, and events)
	PadProbeTypeDataUpstream    PadProbeType = C.GST_PAD_PROBE_TYPE_DATA_UPSTREAM    // (128) – probe upstream data (events)
	PadProbeTypeDataBoth        PadProbeType = C.GST_PAD_PROBE_TYPE_DATA_BOTH        // (240) – probe upstream and downstream data (buffers, buffer lists, and events)
	PadProbeTypeBlockDownstream PadProbeType = C.GST_PAD_PROBE_TYPE_BLOCK_DOWNSTREAM // (114) – probe and block downstream data (buffers, buffer lists, and events)
	PadProbeTypeBlockUpstream   PadProbeType = C.GST_PAD_PROBE_TYPE_BLOCK_UPSTREAM   // (130) – probe and block upstream data (events)
	PadProbeTypeEventBoth       PadProbeType = C.GST_PAD_PROBE_TYPE_EVENT_BOTH       // (192) – probe upstream and downstream events
	PadProbeTypeQueryBoth       PadProbeType = C.GST_PAD_PROBE_TYPE_QUERY_BOTH       // (1536) – probe upstream and downstream queries
	PadProbeTypeAllBoth         PadProbeType = C.GST_PAD_PROBE_TYPE_ALL_BOTH         // (1776) – probe upstream events and queries and downstream buffers, buffer lists, events and queries
	PadProbeTypeScheduling      PadProbeType = C.GST_PAD_PROBE_TYPE_SCHEDULING       // (12288) – probe push and pull
)

// SchedulingFlags casts GstSchedulingFlags
type SchedulingFlags int

// Type casts
const (
	SchedulingFlagSeekable         SchedulingFlags = C.GST_SCHEDULING_FLAG_SEEKABLE          // (1) – if seeking is possible
	SchedulingFlagSequential       SchedulingFlags = C.GST_SCHEDULING_FLAG_SEQUENTIAL        // (2) – if sequential access is recommended
	SchedulingFlagBandwidthLimited SchedulingFlags = C.GST_SCHEDULING_FLAG_BANDWIDTH_LIMITED // (4) – if bandwidth is limited and buffering possible (since 1.2)
)

// State is a type cast of the C GstState
type State int

// Type casting for GstStates
const (
	VoidPending  State = C.GST_STATE_VOID_PENDING // (0) – no pending state.
	StateNull    State = C.GST_STATE_NULL         // (1) – the NULL state or initial state of an element.
	StateReady   State = C.GST_STATE_READY        // (2) – the element is ready to go to PAUSED.
	StatePaused  State = C.GST_STATE_PAUSED       // (3) – the element is PAUSED, it is ready to accept and process data. Sink elements however only accept one buffer and then block.
	StatePlaying State = C.GST_STATE_PLAYING      // (4) – the element is PLAYING, the GstClock is running and the data is flowing.
)

// String returns the string representation of this state.
func (s State) String() string {
	return C.GoString(C.gst_element_state_get_name((C.GstState)(s)))
}

// SeekFlags is a representation of GstSeekFlags.
type SeekFlags int

// Type casts of SeekFlags
const (
	SeekFlagNone        SeekFlags = C.GST_SEEK_FLAG_NONE
	SeekFlagFlush       SeekFlags = C.GST_SEEK_FLAG_FLUSH
	SeekFlagAccurate    SeekFlags = C.GST_SEEK_FLAG_ACCURATE
	SeekFlagKeyUnit     SeekFlags = C.GST_SEEK_FLAG_KEY_UNIT
	SeekFlagSegment     SeekFlags = C.GST_SEEK_FLAG_SEGMENT
	SeekFlagSkip        SeekFlags = C.GST_SEEK_FLAG_SKIP
	SeekFlagSnapBefore  SeekFlags = C.GST_SEEK_FLAG_SNAP_BEFORE
	SeekFlagSnapAfter   SeekFlags = C.GST_SEEK_FLAG_SNAP_AFTER
	SeekFlagSnapNearest SeekFlags = C.GST_SEEK_FLAG_SNAP_NEAREST
)

// SeekType is a representation of GstSeekType.
type SeekType int

// Type casts of seek types
const (
	SeekTypeNone SeekType = C.GST_SEEK_TYPE_NONE
	SeekTypeSet  SeekType = C.GST_SEEK_TYPE_SET
	SeekTypeEnd  SeekType = C.GST_SEEK_TYPE_END
)

// StateChange is the different state changes an element goes through. StateNull ⇒ StatePlaying is called
// an upwards state change and StatePlaying ⇒ StateNull a downwards state change.
//
// See https://gstreamer.freedesktop.org/documentation/gstreamer/gstelement.html?gi-language=c#GstStateChange
// for more information on the responsibiltiies of elements during each transition.
type StateChange int

// StateChange castings
const (
	StateChangeNullToReady      StateChange = C.GST_STATE_CHANGE_NULL_TO_READY
	StateChangeReadyToPaused    StateChange = C.GST_STATE_CHANGE_READY_TO_PAUSED
	StateChangePausedToPlaying  StateChange = C.GST_STATE_CHANGE_PAUSED_TO_PLAYING
	StateChangePlayingToPaused  StateChange = C.GST_STATE_CHANGE_PLAYING_TO_PAUSED
	StateChangePausedToReady    StateChange = C.GST_STATE_CHANGE_PAUSED_TO_READY
	StateChangeReadyToNull      StateChange = C.GST_STATE_CHANGE_READY_TO_NULL
	StateChangeNullToNull       StateChange = C.GST_STATE_CHANGE_NULL_TO_NULL
	StateChangeReadyToReady     StateChange = C.GST_STATE_CHANGE_READY_TO_READY
	StateChangePausedToPaused   StateChange = C.GST_STATE_CHANGE_PAUSED_TO_PAUSED
	StateChangePlayingToPlaying StateChange = C.GST_STATE_CHANGE_PLAYING_TO_PLAYING
)

// String returns the string representation of a StateChange
func (s StateChange) String() string {
	return C.GoString(C.gst_state_change_get_name(C.GstStateChange(s)))
}

// StateChangeReturn is a representation of GstStateChangeReturn.
type StateChangeReturn int

// Type casts of state change returns
const (
	StateChangeFailure   StateChangeReturn = C.GST_STATE_CHANGE_FAILURE
	StateChangeSuccess   StateChangeReturn = C.GST_STATE_CHANGE_SUCCESS
	StateChangeAsync     StateChangeReturn = C.GST_STATE_CHANGE_ASYNC
	StateChangeNoPreroll StateChangeReturn = C.GST_STATE_CHANGE_NO_PREROLL
)

func (s StateChangeReturn) String() string {
	return C.GoString(C.gst_element_state_change_return_get_name(C.GstStateChangeReturn(s)))
}

// ElementFlags casts C GstElementFlags to a go type
type ElementFlags int

// Type casting of element flags
const (
	ElementFlagLockedState  ElementFlags = C.GST_ELEMENT_FLAG_LOCKED_STATE  // (16) – ignore state changes from parent
	ElementFlagSink         ElementFlags = C.GST_ELEMENT_FLAG_SINK          // (32) – the element is a sink
	ElementFlagSource       ElementFlags = C.GST_ELEMENT_FLAG_SOURCE        // (64) – the element is a source.
	ElementFlagProvideClock ElementFlags = C.GST_ELEMENT_FLAG_PROVIDE_CLOCK // (128) – the element can provide a clock
	ElementFlagRequireClock ElementFlags = C.GST_ELEMENT_FLAG_REQUIRE_CLOCK // (256) – the element requires a clock
	ElementFlagIndexable    ElementFlags = C.GST_ELEMENT_FLAG_INDEXABLE     // (512) – the element can use an index
	ElementFlagLast         ElementFlags = C.GST_ELEMENT_FLAG_LAST          // (16384) – offset to define more flags
)

// MiniObjectFlags casts GstMiniObjectFlags to a go type.
type MiniObjectFlags int

// Type casting of mini-object flags
const (
	MiniObjectFlagLockable     MiniObjectFlags = C.GST_MINI_OBJECT_FLAG_LOCKABLE      // (1) – the object can be locked and unlocked with gst_mini_object_lock and gst_mini_object_unlock.
	MiniObjectFlagLockReadOnly MiniObjectFlags = C.GST_MINI_OBJECT_FLAG_LOCK_READONLY // (2) – the object is permanently locked in READONLY mode. Only read locks can be performed on the object.
	MiniObjectFlagMayBeLeaked  MiniObjectFlags = C.GST_MINI_OBJECT_FLAG_MAY_BE_LEAKED // (4) – the object is expected to stay alive even after gst_deinit has been called and so should be ignored by leak detection tools. (Since: 1.10)
	MiniObjectFlagLast         MiniObjectFlags = C.GST_MINI_OBJECT_FLAG_LAST          // (16) – first flag that can be used by subclasses.
)

// FlowReturn is go type casting for GstFlowReturn.
type FlowReturn int

// Type casting of the GstFlowReturn types. Custom ones are omitted for now.
const (
	FlowOK            FlowReturn = C.GST_FLOW_OK             // Data passing was ok
	FlowNotLinked     FlowReturn = C.GST_FLOW_NOT_LINKED     // Pad is not linked
	FlowFlushing      FlowReturn = C.GST_FLOW_FLUSHING       // Pad is flushing
	FlowEOS           FlowReturn = C.GST_FLOW_EOS            // Pad is EOS
	FlowNotNegotiated FlowReturn = C.GST_FLOW_NOT_NEGOTIATED // Pad is not negotiated
	FlowError         FlowReturn = C.GST_FLOW_ERROR          // Some (fatal) error occurred
	FlowNotSupported  FlowReturn = C.GST_FLOW_NOT_SUPPORTED  // The operation is not supported.
)

// String impelements a stringer on FlowReturn
func (f FlowReturn) String() string {
	return C.GoString(C.gst_flow_get_name(C.GstFlowReturn(f)))
}

// MapFlags is a go casting of GstMapFlags
type MapFlags int

// Type casting of the map flag types
const (
	MapRead     MapFlags = C.GST_MAP_READ      //  (1) – map for read access
	MapWrite    MapFlags = C.GST_MAP_WRITE     // (2) - map for write access
	MapFlagLast MapFlags = C.GST_MAP_FLAG_LAST // (65536) – first flag that can be used for custom purposes
)

// StreamStatusType represents a type of change in a stream's status
type StreamStatusType int

// Type castings of the stream status types
const (
	StreamStatusCreate  StreamStatusType = C.GST_STREAM_STATUS_TYPE_CREATE  // (0) – A new thread need to be created.
	StreamStatusEnter   StreamStatusType = C.GST_STREAM_STATUS_TYPE_ENTER   // (1) – a thread entered its loop function
	StreamStatusLeave   StreamStatusType = C.GST_STREAM_STATUS_TYPE_LEAVE   // (2) – a thread left its loop function
	StreamStatusDestroy StreamStatusType = C.GST_STREAM_STATUS_TYPE_DESTROY // (3) – a thread is destroyed
	StreamStatusStart   StreamStatusType = C.GST_STREAM_STATUS_TYPE_START   // (8) – a thread is started
	StreamStatusPause   StreamStatusType = C.GST_STREAM_STATUS_TYPE_PAUSE   // (9) – a thread is paused
	StreamStatusStop    StreamStatusType = C.GST_STREAM_STATUS_TYPE_STOP    // (10) – a thread is stopped
)

func (s StreamStatusType) String() string {
	switch s {
	case StreamStatusCreate:
		return "A new thread needs to be created"
	case StreamStatusEnter:
		return "A thread has entered its loop function"
	case StreamStatusLeave:
		return "A thread has left its loop function"
	case StreamStatusDestroy:
		return "A thread has been destroyed"
	case StreamStatusStart:
		return "A thread has started"
	case StreamStatusPause:
		return "A thread has paused"
	case StreamStatusStop:
		return "A thread has stopped"
	}
	return ""
}

// StructureChangeType is a go representation of a GstStructureChangeType
type StructureChangeType int

// Type castings of StructureChangeTypes
const (
	StructureChangePadLink   StructureChangeType = C.GST_STRUCTURE_CHANGE_TYPE_PAD_LINK   // (0) – Pad linking is starting or done.
	StructureChangePadUnlink StructureChangeType = C.GST_STRUCTURE_CHANGE_TYPE_PAD_UNLINK // (1) – Pad unlinking is starting or done.
)

// String implements a stringer on StructureChangeTypes.
func (s StructureChangeType) String() string {
	switch s {
	case StructureChangePadLink:
		return "pad link"
	case StructureChangePadUnlink:
		return "pad unlink"
	}
	return ""
}

// ProgressType is a go representation of a GstProgressType
type ProgressType int

// Type castings of ProgressTypes
const (
	ProgressTypeStart     ProgressType = C.GST_PROGRESS_TYPE_START    // (0) – A new task started.
	ProgressTypeContinue  ProgressType = C.GST_PROGRESS_TYPE_CONTINUE // (1) – A task completed and a new one continues.
	ProgressTypeComplete  ProgressType = C.GST_PROGRESS_TYPE_COMPLETE // (2) – A task completed.
	ProgressTypeCancelled ProgressType = C.GST_PROGRESS_TYPE_CANCELED // (3) – A task was canceled.
	ProgressTypeError     ProgressType = C.GST_PROGRESS_TYPE_ERROR    // (4) – A task caused an error. An error message is also posted on the bus.
)

// String implements a stringer on ProgressTypes
func (p ProgressType) String() string {
	switch p {
	case ProgressTypeStart:
		return "started"
	case ProgressTypeContinue:
		return "continuing"
	case ProgressTypeComplete:
		return "completed"
	case ProgressTypeCancelled:
		return "cancelled"
	case ProgressTypeError:
		return "error"
	}
	return ""
}

// StreamType is a go representation of a GstStreamType
type StreamType int

// Type castings of stream types
const (
	StreamTypeUnknown   StreamType = C.GST_STREAM_TYPE_UNKNOWN   // (1) – The stream is of unknown (unclassified) type.
	StreamTypeAudio     StreamType = C.GST_STREAM_TYPE_AUDIO     // (2) – The stream is of audio data
	StreamTypeVideo     StreamType = C.GST_STREAM_TYPE_VIDEO     // (4) – The stream carries video data
	StreamTypeContainer StreamType = C.GST_STREAM_TYPE_CONTAINER // (8) – The stream is a muxed container type
	StreamTypeText      StreamType = C.GST_STREAM_TYPE_TEXT      // (16) – The stream contains subtitle / subpicture data.
)

// String implements a stringer on StreamTypes.
func (s StreamType) String() string {
	name := C.gst_stream_type_get_name((C.GstStreamType)(s))
	defer C.free(unsafe.Pointer(name))
	return C.GoString(name)
}

// StreamFlags represent configuration options for a new stream.
type StreamFlags int

// Type castings of StreamFlags
const (
	StreamFlagNone     StreamFlags = C.GST_STREAM_FLAG_NONE     // (0) – This stream has no special attributes
	StreamFlagSparse   StreamFlags = C.GST_STREAM_FLAG_SPARSE   // (1) – This stream is a sparse stream (e.g. a subtitle stream), data may flow only in irregular intervals with large gaps in between.
	StreamFlagSelect   StreamFlags = C.GST_STREAM_FLAG_SELECT   // (2) – This stream should be selected by default. This flag may be used by demuxers to signal that a stream should be selected by default in a playback scenario.
	StreamFlagUnselect StreamFlags = C.GST_STREAM_FLAG_UNSELECT // (4) – This stream should not be selected by default. This flag may be used by demuxers to signal that a stream should not be selected by default in a playback scenario, but only if explicitly selected by the user (e.g. an audio track for the hard of hearing or a director's commentary track).
)

// MemoryFlags represent flags for wrapped memory
type MemoryFlags int

// Type castins of MemoryFlags
const (
	MemoryFlagReadOnly             MemoryFlags = C.GST_MEMORY_FLAG_READONLY              // (2) – memory is readonly. It is not allowed to map the memory with GST_MAP_WRITE.
	MemoryFlagNoShare              MemoryFlags = C.GST_MEMORY_FLAG_NO_SHARE              // (16) – memory must not be shared. Copies will have to be made when this memory needs to be shared between buffers. (DEPRECATED: do not use in new code, instead you should create a custom GstAllocator for memory pooling instead of relying on the GstBuffer they were originally attached to.)
	MemoryFlagZeroPrefixed         MemoryFlags = C.GST_MEMORY_FLAG_ZERO_PREFIXED         // (32) – the memory prefix is filled with 0 bytes
	MemoryFlagZeroPadded           MemoryFlags = C.GST_MEMORY_FLAG_ZERO_PADDED           // (64) – the memory padding is filled with 0 bytes
	MemoryFlagPhysicallyContiguous MemoryFlags = C.GST_MEMORY_FLAG_PHYSICALLY_CONTIGUOUS // (128) – the memory is physically contiguous. (Since: 1.2)
	MemoryFlagNotMappable          MemoryFlags = C.GST_MEMORY_FLAG_NOT_MAPPABLE          // (256) – the memory can't be mapped via gst_memory_map without any preconditions. (Since: 1.2)
	MemoryFlagLast                 MemoryFlags = 1048576                                 // first flag that can be used for custom purposes
)

// URIType casts C GstURIType to a go type
type URIType int

// Type cast URI types
const (
	URIUnknown URIType = C.GST_URI_UNKNOWN // (0) – The URI direction is unknown
	URISink    URIType = C.GST_URI_SINK    // (1) – The URI is a consumer.
	URISource  URIType = C.GST_URI_SRC     // (2) - The URI is a producer.
)

func (u URIType) String() string {
	switch u {
	case URIUnknown:
		return "Unknown"
	case URISink:
		return "Sink"
	case URISource:
		return "Source"
	}
	return ""
}

// MetaFlags casts C GstMetaFlags to a go type.
type MetaFlags int

// Type casts of GstMetaFlags
const (
	MetaFlagNone     MetaFlags = C.GST_META_FLAG_NONE     // (0) – no flags
	MetaFlagReadOnly MetaFlags = C.GST_META_FLAG_READONLY // (1) – metadata should not be modified
	MetaFlagPooled   MetaFlags = C.GST_META_FLAG_POOLED   // (2) – metadata is managed by a bufferpool
	MetaFlagLocked   MetaFlags = C.GST_META_FLAG_LOCKED   // (4) – metadata should not be removed
	MetaFlagLast     MetaFlags = C.GST_META_FLAG_LAST     // (65536) – additional flags can be added starting from this flag.
)

// QueryType casts GstQueryType
type QueryType int

// Type casts
const (
	QueryUnknown    QueryType = C.GST_QUERY_UNKNOWN     // (0) – unknown query type
	QueryPosition   QueryType = C.GST_QUERY_POSITION    // (2563) – current position in stream
	QueryDuration   QueryType = C.GST_QUERY_DURATION    // (5123) – total duration of the stream
	QueryLatency    QueryType = C.GST_QUERY_LATENCY     // (7683) – latency of stream
	QueryJitter     QueryType = C.GST_QUERY_JITTER      // (10243) – current jitter of stream
	QueryRate       QueryType = C.GST_QUERY_RATE        // (12803) – current rate of the stream
	QuerySeeking    QueryType = C.GST_QUERY_SEEKING     // (15363) – seeking capabilities
	QuerySegment    QueryType = C.GST_QUERY_SEGMENT     // (17923) – segment start/stop positions
	QueryConvert    QueryType = C.GST_QUERY_CONVERT     // (20483) – convert values between formats
	QueryFormats    QueryType = C.GST_QUERY_FORMATS     // (23043) – query supported formats for convert
	QueryBuffering  QueryType = C.GST_QUERY_BUFFERING   // (28163) – query available media for efficient seeking.
	QueryCustom     QueryType = C.GST_QUERY_CUSTOM      // (30723) – a custom application or element defined que	QueryType = C.ry.
	QueryURI        QueryType = C.GST_QUERY_URI         // (33283) – query the URI of the source or sink.
	QueryAllocation QueryType = C.GST_QUERY_ALLOCATION  // (35846) – the buffer allocation properties
	QueryScheduling QueryType = C.GST_QUERY_SCHEDULING  // (38401) – the scheduling properties
	QueryAcceptCaps QueryType = C.GST_QUERY_ACCEPT_CAPS // (40963) – the accept caps query
	QueryCaps       QueryType = C.GST_QUERY_CAPS        // (43523) – the caps query
	QueryDrain      QueryType = C.GST_QUERY_DRAIN       // (46086) – wait till all serialized data is consumed downstream
	QueryContext    QueryType = C.GST_QUERY_CONTEXT     // (48643) – query the pipeline-local context from downstream or upstream (since 1.2)
	QueryBitrate    QueryType = C.GST_QUERY_BITRATE     // (51202) – the bitrate query (since 1.16)
)

func (q QueryType) String() string { return C.GoString(C.gst_query_type_get_name(C.GstQueryType(q))) }

// QueryTypeFlags casts GstQueryTypeFlags
type QueryTypeFlags int

// Type casts
const (
	QueryTypeUpstream   QueryTypeFlags = C.GST_QUERY_TYPE_UPSTREAM   // (1) – Set if the query can travel upstream.
	QueryTypeDownstream QueryTypeFlags = C.GST_QUERY_TYPE_DOWNSTREAM // (2) – Set if the query can travel downstream.
	QueryTypeSerialized QueryTypeFlags = C.GST_QUERY_TYPE_SERIALIZED // (4) – Set if the query should be serialized with data flow.
)

// TaskState casts GstTaskState
type TaskState int

// Type castings
const (
	TaskStarted TaskState = C.GST_TASK_STARTED // (0) – the task is started and running
	TaskStopped TaskState = C.GST_TASK_STOPPED // (1) – the task is stopped
	TaskPaused  TaskState = C.GST_TASK_PAUSED  // (2) – the task is paused
)

// TOCScope represents the scope of a TOC.
type TOCScope int

// Type castings of TOCScopes.
const (
	// (1) – global TOC representing all selectable options (this is what applications are usually interested in)
	TOCScopeGlobal TOCScope = C.GST_TOC_SCOPE_GLOBAL
	// (2) – TOC for the currently active/selected stream (this is a TOC representing the current stream from start
	// to EOS, and is what a TOC writer / muxer is usually interested in; it will usually be a subset of the global
	// TOC, e.g. just the chapters of the current title, or the chapters selected for playback from the current title)
	TOCScopeCurrent TOCScope = C.GST_TOC_SCOPE_CURRENT
)

// String implements a stringer on a TOCScope
func (t TOCScope) String() string {
	switch t {
	case TOCScopeGlobal:
		return "global"
	case TOCScopeCurrent:
		return "current"
	}
	return ""
}

// TOCLoopType represents a GstTocLoopType
type TOCLoopType int

// Type castings of TOCLoopTypes
const (
	TOCLoopNone     TOCLoopType = C.GST_TOC_LOOP_NONE      // (0) – single forward playback
	TOCLoopForward  TOCLoopType = C.GST_TOC_LOOP_FORWARD   // (1) – repeat forward
	TOCLoopReverse  TOCLoopType = C.GST_TOC_LOOP_REVERSE   // (2) – repeat backward
	TOCLoopPingPong TOCLoopType = C.GST_TOC_LOOP_PING_PONG // (3) – repeat forward and backward
)

// TOCEntryType represents a GstTocEntryType.
type TOCEntryType int

// Type castings of TOCEntryTypes
const (
	TOCEntryTypeAngle   TOCEntryType = C.GST_TOC_ENTRY_TYPE_ANGLE   // (-3) – entry is an angle (i.e. an alternative)
	TOCEntryTypeVersion TOCEntryType = C.GST_TOC_ENTRY_TYPE_VERSION // (-2) – entry is a version (i.e. alternative)
	TOCEntryTypeEdition TOCEntryType = C.GST_TOC_ENTRY_TYPE_EDITION // (-1) – entry is an edition (i.e. alternative)
	TOCEntryTypeInvalid TOCEntryType = C.GST_TOC_ENTRY_TYPE_INVALID // (0) – invalid entry type value
	TOCEntryTypeTitle   TOCEntryType = C.GST_TOC_ENTRY_TYPE_TITLE   // (1) – entry is a title (i.e. a part of a sequence)
	TOCEntryTypeTrack   TOCEntryType = C.GST_TOC_ENTRY_TYPE_TRACK   // (2) – entry is a track (i.e. a part of a sequence)
	TOCEntryTypeChapter TOCEntryType = C.GST_TOC_ENTRY_TYPE_CHAPTER // (3) – entry is a chapter (i.e. a part of a sequence)
)

// TagFlag represents a GstTagFlag
type TagFlag int

// Type castins of TagFlags
const (
	TagFlagUndefined TagFlag = C.GST_TAG_FLAG_UNDEFINED // (0) – undefined flag
	TagFlagMeta      TagFlag = C.GST_TAG_FLAG_META      // (1) – tag is meta data
	TagFlagEncoded   TagFlag = C.GST_TAG_FLAG_ENCODED   // (2) – tag is encoded
	TagFlagDecoded   TagFlag = C.GST_TAG_FLAG_DECODED   // (3) – tag is decoded
	TagFlagCount     TagFlag = C.GST_TAG_FLAG_COUNT     // (4) – number of tag flags
)

// TagMergeMode represents a GstTagMergeMode.
// See: https://gstreamer.freedesktop.org/documentation/gstreamer/gsttaglist.html#GstTagMergeMode
type TagMergeMode int

// Type castings of TagMergeModes
const (
	TagMergeUndefined  TagMergeMode = C.GST_TAG_MERGE_UNDEFINED   // (0) – undefined merge mode
	TagMergeReplaceAll TagMergeMode = C.GST_TAG_MERGE_REPLACE_ALL // (1) – replace all tags (clear list and append)
	TagMergeReplace    TagMergeMode = C.GST_TAG_MERGE_REPLACE     // (2) – replace tags
	TagMergeAppend     TagMergeMode = C.GST_TAG_MERGE_APPEND      // (3) – append tags
	TagMergePrepend    TagMergeMode = C.GST_TAG_MERGE_PREPEND     // (4) – prepend tags
	TagMergeKeep       TagMergeMode = C.GST_TAG_MERGE_KEEP        // (5) – keep existing tags
	TagMergeKeepAll    TagMergeMode = C.GST_TAG_MERGE_KEEP_ALL    // (6) – keep all existing tags
	TagMergeCount      TagMergeMode = C.GST_TAG_MERGE_COUNT       // (7) – the number of merge modes
)

// TagScope represents a GstTagScope
type TagScope int

// Type castings of tag scopes
const (
	TagScopeStream TagScope = C.GST_TAG_SCOPE_STREAM // (0) – tags specific to this single stream
	TagScopeGlobal TagScope = C.GST_TAG_SCOPE_GLOBAL // (1) – global tags for the complete medium
)

// Tag wraps the builtin gstreamer tags
type Tag string

// Type castings of Tags
// For more information see: https://gstreamer.freedesktop.org/documentation/gstreamer/gsttaglist.html?gi-language=c#constants
const (
	TagAlbum                        Tag = C.GST_TAG_ALBUM
	TagAlbumArtist                  Tag = C.GST_TAG_ALBUM_ARTIST
	TagAlbumArtistSortName          Tag = C.GST_TAG_ALBUM_ARTIST_SORTNAME
	TagAlbumGain                    Tag = C.GST_TAG_ALBUM_GAIN
	TagAlbumPeak                    Tag = C.GST_TAG_ALBUM_PEAK
	TagAlbumSortName                Tag = C.GST_TAG_ALBUM_SORTNAME
	TagAlbumVolumeCount             Tag = C.GST_TAG_ALBUM_VOLUME_COUNT
	TagAlbumVolumeNumber            Tag = C.GST_TAG_ALBUM_VOLUME_NUMBER
	TagApplicationData              Tag = C.GST_TAG_APPLICATION_DATA
	TagApplicationName              Tag = C.GST_TAG_APPLICATION_NAME
	TagArtist                       Tag = C.GST_TAG_ARTIST
	TagArtistSortName               Tag = C.GST_TAG_ARTIST_SORTNAME
	TagAttachment                   Tag = C.GST_TAG_ATTACHMENT
	TagAudioCodec                   Tag = C.GST_TAG_AUDIO_CODEC
	TagBeatsPerMinute               Tag = C.GST_TAG_BEATS_PER_MINUTE
	TagBitrate                      Tag = C.GST_TAG_BITRATE
	TagCodec                        Tag = C.GST_TAG_CODEC
	TagComment                      Tag = C.GST_TAG_COMMENT
	TagComposer                     Tag = C.GST_TAG_COMPOSER
	TagComposerSortName             Tag = C.GST_TAG_COMPOSER_SORTNAME
	TagConductor                    Tag = C.GST_TAG_CONDUCTOR
	TagContact                      Tag = C.GST_TAG_CONTACT
	TagContainerFormat              Tag = C.GST_TAG_CONTAINER_FORMAT
	TagCopyright                    Tag = C.GST_TAG_COPYRIGHT
	TagCopyrightURI                 Tag = C.GST_TAG_COPYRIGHT_URI
	TagDate                         Tag = C.GST_TAG_DATE
	TagDateTime                     Tag = C.GST_TAG_DATE_TIME
	TagDescription                  Tag = C.GST_TAG_DESCRIPTION
	TagDeviceManufacturer           Tag = C.GST_TAG_DEVICE_MANUFACTURER
	TagDeviceModel                  Tag = C.GST_TAG_DEVICE_MODEL
	TagDuration                     Tag = C.GST_TAG_DURATION
	TagEncodedBy                    Tag = C.GST_TAG_ENCODED_BY
	TagEncoder                      Tag = C.GST_TAG_ENCODER
	TagEncoderVersion               Tag = C.GST_TAG_ENCODER_VERSION
	TagExtendedComment              Tag = C.GST_TAG_EXTENDED_COMMENT
	TagGenre                        Tag = C.GST_TAG_GENRE
	TagGeoLocationCaptureDirection  Tag = C.GST_TAG_GEO_LOCATION_CAPTURE_DIRECTION
	TagGeoLocationCity              Tag = C.GST_TAG_GEO_LOCATION_CITY
	TagGeoLocationCountry           Tag = C.GST_TAG_GEO_LOCATION_COUNTRY
	TagGeoLocationElevation         Tag = C.GST_TAG_GEO_LOCATION_ELEVATION
	TagGeoLocationHoriozontalError  Tag = C.GST_TAG_GEO_LOCATION_HORIZONTAL_ERROR
	TagGeoLocationLatitude          Tag = C.GST_TAG_GEO_LOCATION_LATITUDE
	TagGeoLocationLongitude         Tag = C.GST_TAG_GEO_LOCATION_LONGITUDE
	TagGeoLocationMovementDirection Tag = C.GST_TAG_GEO_LOCATION_MOVEMENT_DIRECTION
	TagGeoLocationMovementSpeed     Tag = C.GST_TAG_GEO_LOCATION_MOVEMENT_SPEED
	TagGeoLocationName              Tag = C.GST_TAG_GEO_LOCATION_NAME
	TagGeoLocationSubLocation       Tag = C.GST_TAG_GEO_LOCATION_SUBLOCATION
	TagGrouping                     Tag = C.GST_TAG_GROUPING
	TagHomepage                     Tag = C.GST_TAG_HOMEPAGE
	TagImage                        Tag = C.GST_TAG_IMAGE
	TagImageOrientation             Tag = C.GST_TAG_IMAGE_ORIENTATION
	TagInterpretedBy                Tag = C.GST_TAG_INTERPRETED_BY
	TagISRC                         Tag = C.GST_TAG_ISRC
	TagKeywords                     Tag = C.GST_TAG_KEYWORDS
	TagLanguageCode                 Tag = C.GST_TAG_LANGUAGE_CODE
	TagLanguageName                 Tag = C.GST_TAG_LANGUAGE_NAME
	TagLicense                      Tag = C.GST_TAG_LICENSE
	TagLicenseURI                   Tag = C.GST_TAG_LICENSE_URI
	TagLocation                     Tag = C.GST_TAG_LOCATION
	TagLyrics                       Tag = C.GST_TAG_LYRICS
	TagMaximumBitrate               Tag = C.GST_TAG_MAXIMUM_BITRATE
	TagMIDIBaseNote                 Tag = C.GST_TAG_MIDI_BASE_NOTE
	TagMinimumBitrate               Tag = C.GST_TAG_MINIMUM_BITRATE
	TagNominalBitrate               Tag = C.GST_TAG_NOMINAL_BITRATE
	TagOrganization                 Tag = C.GST_TAG_ORGANIZATION
	TagPerformer                    Tag = C.GST_TAG_PERFORMER
	TagPreviewImage                 Tag = C.GST_TAG_PREVIEW_IMAGE
	TagPrivateData                  Tag = C.GST_TAG_PRIVATE_DATA
	TagPublisher                    Tag = C.GST_TAG_PUBLISHER
	TagReferenceLevel               Tag = C.GST_TAG_REFERENCE_LEVEL
	TagSerial                       Tag = C.GST_TAG_SERIAL
	TagShowEpisodeNumber            Tag = C.GST_TAG_SHOW_EPISODE_NUMBER
	TagShowName                     Tag = C.GST_TAG_SHOW_NAME
	TagShowSeasonNumber             Tag = C.GST_TAG_SHOW_SEASON_NUMBER
	TagShowSortName                 Tag = C.GST_TAG_SHOW_SORTNAME
	TagSubtitleCodec                Tag = C.GST_TAG_SUBTITLE_CODEC
	TagTitle                        Tag = C.GST_TAG_TITLE
	TagTitleSortName                Tag = C.GST_TAG_TITLE_SORTNAME
	TagTrackCount                   Tag = C.GST_TAG_TRACK_COUNT
	TagTrackGain                    Tag = C.GST_TAG_TRACK_GAIN
	TagTrackNumber                  Tag = C.GST_TAG_TRACK_NUMBER
	TagTrackPeak                    Tag = C.GST_TAG_TRACK_PEAK
	TagUserRating                   Tag = C.GST_TAG_USER_RATING
	TagVersion                      Tag = C.GST_TAG_VERSION
	TagVideoCodec                   Tag = C.GST_TAG_VIDEO_CODEC
)

// TypeFindProbability represents a probability for type find functions. Higher values
// reflect higher certainty.
type TypeFindProbability int

// Type castings
const (
	TypeFindNone          TypeFindProbability = C.GST_TYPE_FIND_NONE           // (0) – type undetected.
	TypeFindMinimum       TypeFindProbability = C.GST_TYPE_FIND_MINIMUM        // (1) – unlikely typefind.
	TypeFindPossible      TypeFindProbability = C.GST_TYPE_FIND_POSSIBLE       // (50) – possible type detected.
	TypeFindLikely        TypeFindProbability = C.GST_TYPE_FIND_LIKELY         // (80) – likely a type was detected.
	TypeFindNearlyCertain TypeFindProbability = C.GST_TYPE_FIND_NEARLY_CERTAIN // (99) – nearly certain that a type was detected.
	TypeFindMaximum       TypeFindProbability = C.GST_TYPE_FIND_MAXIMUM        // (100) – very certain a type was detected.
)
