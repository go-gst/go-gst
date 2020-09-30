package gst

// #include "gst.go.h"
import "C"

import "unsafe"

// ClockTime is a go representation of a GstClockTime. Most of the time these are casted
// to time.Duration objects. It represents a time value in nanoseconds.
type ClockTime uint64

// ClockTimeDiff is a datatype to hold a time difference, measured in nanoseconds.
type ClockTimeDiff int64

const (
	// ClockFormat is the string used when formatting clock strings
	ClockFormat string = "u:%02u:%02u.%09u"
	// ClockTimeNone means infinite timeout (unsigned representation of -1) or an otherwise unknown value.
	ClockTimeNone ClockTime = C.GST_CLOCK_TIME_NONE
	// BufferOffsetNone is a constant for no-offset return results.
	BufferOffsetNone ClockTime = C.GST_BUFFER_OFFSET_NONE
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

// PadDirection is a cast of GstPadDirection to a go type.
type PadDirection int

// Type casting of pad directions
const (
	PadUnknown PadDirection = C.GST_PAD_UNKNOWN // (0) - the direction is unknown
	PadSource  PadDirection = C.GST_PAD_SRC     // (1) - the pad is a source pad
	PadSink    PadDirection = C.GST_PAD_SINK    // (2) - the pad is a sink pad
)

// String implements a Stringer on PadDirection.
func (p PadDirection) String() string {
	switch p {
	case PadUnknown:
		return "Unknown"
	case PadSource:
		return "Src"
	case PadSink:
		return "Sink"
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

// PadPresence is a cast of GstPadPresence to a go type.
type PadPresence int

// Type casting of pad presences
const (
	PadAlways    PadPresence = C.GST_PAD_ALWAYS    // (0) - the pad is always available
	PadSometimes PadPresence = C.GST_PAD_SOMETIMES // (1) - the pad will become available depending on the media stream
	PadRequest   PadPresence = C.GST_PAD_REQUEST   // (2) - the pad is only available on request with gst_element_request_pad.
)

// String implements a stringer on PadPresence.
func (p PadPresence) String() string {
	switch p {
	case PadAlways:
		return "Always"
	case PadSometimes:
		return "Sometimes"
	case PadRequest:
		return "Request"
	}
	return ""
}

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

// StateChangeReturn is a representation of GstStateChangeReturn.
type StateChangeReturn int

// Type casts of state change returns
const (
	StateChangeFailure   StateChangeReturn = C.GST_STATE_CHANGE_FAILURE
	StateChangeSuccess   StateChangeReturn = C.GST_STATE_CHANGE_SUCCESS
	StateChangeAsync     StateChangeReturn = C.GST_STATE_CHANGE_ASYNC
	StateChangeNoPreroll StateChangeReturn = C.GST_STATE_CHANGE_NO_PREROLL
)

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
