package gst

// #include "gst.go.h"
import "C"

import (
	"runtime"
	"time"
	"unsafe"
)

// Event is a go wrapper around a GstEvent.
type Event struct {
	ptr *C.GstEvent
}

// // FromGstEventUnsafe is an alias to FromGstEventUnsafeNone.
// func FromGstEventUnsafe(ev unsafe.Pointer) *Event { return FromGstEventUnsafeNone(ev) }

// FromGstEventUnsafeNone wraps the pointer to the given C GstEvent with the go type.
// A ref is taken and finalizer applied.
func FromGstEventUnsafeNone(ev unsafe.Pointer) *Event {
	event := ToGstEvent(ev)
	event.Ref()
	runtime.SetFinalizer(event, (*Event).Unref)
	return event
}

// FromGstEventUnsafeFull wraps the pointer to the given C GstEvent without taking a ref.
// A finalizer is applied.
func FromGstEventUnsafeFull(ev unsafe.Pointer) *Event {
	event := ToGstEvent(ev)
	runtime.SetFinalizer(event, (*Event).Unref)
	return event
}

// ToGstEvent converts the given pointer into an Event without affecting the ref count or
// placing finalizers.
func ToGstEvent(ev unsafe.Pointer) *Event {
	return wrapEvent((*C.GstEvent)(ev))
}

// Instance returns the underlying GstEvent instance.
func (e *Event) Instance() *C.GstEvent { return C.toGstEvent(unsafe.Pointer(e.ptr)) }

// Type returns the type of the event
func (e *Event) Type() EventType { return EventType(e.Instance()._type) }

// Timestamp returns the timestamp of the event.
func (e *Event) Timestamp() time.Duration {
	ts := e.Instance().timestamp
	return time.Duration(uint64(ts)) * time.Nanosecond
}

// Seqnum returns the sequence number of the event.
func (e *Event) Seqnum() uint32 {
	return uint32(e.Instance().seqnum)
}

// Copy copies the event using the event specific copy function.
func (e *Event) Copy() *Event {
	return FromGstEventUnsafeFull(unsafe.Pointer(C.gst_event_copy(e.Instance())))
}

// CopySegment parses a segment event and copies the Segment into the location given by segment.
func (e *Event) CopySegment(segment *Segment) {
	C.gst_event_copy_segment(e.Instance(), segment.Instance())
}

// GetRunningTimeOffset retrieves the accumulated running time offset of the event.
//
// Events passing through GstPad that have a running time offset set via gst_pad_set_offset will get their
// offset adjusted according to the pad's offset.
//
// If the event contains any information that related to the running time, this information will need to be
// updated before usage with this offset.
func (e *Event) GetRunningTimeOffset() int64 {
	return int64(C.gst_event_get_running_time_offset(e.Instance()))
}

// GetSeqnum retrieves the sequence number of a event.
//
// Events have ever-incrementing sequence numbers, which may also be set explicitly via SetSeqnum. Sequence
// numbers are typically used to indicate that a event corresponds to some other set of events or messages,
// for example an EOS event corresponding to a SEEK event. It is considered good practice to make this
// correspondence when possible, though it is not required.
//
// Note that events and messages share the same sequence number incrementor; two events or messages will never
// have the same sequence number unless that correspondence was made explicitly.
func (e *Event) GetSeqnum() uint32 {
	return uint32(C.gst_event_get_seqnum(e.Instance()))
}

// GetStructure accesses the structure of the event.
func (e *Event) GetStructure() *Structure {
	return wrapStructure(C.gst_event_get_structure(e.Instance()))
}

// HasName checks if event has the given name. This function is usually used to check the name of a custom event.
func (e *Event) HasName(name string) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return gobool(C.gst_event_has_name(e.Instance(), (*C.gchar)(unsafe.Pointer(cName))))
}

// ParseBufferSize gets the format, minsize, maxsize and async-flag in the buffersize event.
func (e *Event) ParseBufferSize() (format Format, minSize, maxSize int64, async bool) {
	var gformat C.GstFormat
	var gmin, gmax C.gint64
	var gasync C.gboolean
	C.gst_event_parse_buffer_size(e.Instance(), &gformat, &gmin, &gmax, &gasync)
	return Format(gformat), int64(gmin), int64(gmax), gobool(gasync)
}

// ParseCaps gets the caps from event. The caps remains valid as long as event remains valid.
func (e *Event) ParseCaps() *Caps {
	var caps *C.GstCaps
	C.gst_event_parse_caps(e.Instance(), &caps)
	return FromGstCapsUnsafeNone(unsafe.Pointer(caps))
}

// ParseFlushStop parses the FLUSH_STOP event and retrieve the reset_time member. Value reflects whether
// time should be reset.
func (e *Event) ParseFlushStop() (resetTime bool) {
	var gbool C.gboolean
	C.gst_event_parse_flush_stop(e.Instance(), &gbool)
	return gobool(gbool)
}

// ParseGap extracts timestamp and duration from a new GAP event.
func (e *Event) ParseGap() (timestamp, duration time.Duration) {
	var ts, dur C.GstClockTime
	C.gst_event_parse_gap(e.Instance(), &ts, &dur)
	return time.Duration(ts), time.Duration(dur)
}

// ParseGapFlags retrieves the gap flags that may have been set on a gap event with SetGapFlags.
// func (e *Event) ParseGapFlags() GapFlags {
// 	var out C.GstGapFlags
// 	C.gst_event_parse_gap_flags(e.Instance(), &out)
// 	return GapFlags(out)
// }

// ParseGroupID returns a group ID if set on the event.
func (e *Event) ParseGroupID() (ok bool, gid uint) {
	var out C.guint
	gok := C.gst_event_parse_group_id(e.Instance(), &out)
	return gobool(gok), uint(out)
}

// ParseInstantRateChange extracts rate and flags from an instant-rate-change event.
// func (e *Event) ParseInstantRateChange() {}

// ParseInstantRateSyncTime extracts the rate multiplier and running times from an instant-rate-sync-time event.
// func (e *Event) ParseInstantRateChange() {}

// ParseLatency gets the latency in the latency event.
func (e *Event) ParseLatency() time.Duration {
	var out C.GstClockTime
	C.gst_event_parse_latency(e.Instance(), &out)
	return time.Duration(out)
}

// ParseProtection parses an event containing protection system specific information and stores the results in
// system_id, data and origin. The data stored in system_id, origin and data are valid until event is released.
func (e *Event) ParseProtection() (systemID string, data *Buffer, origin string) {
	idPtr := C.malloc(C.sizeof_char * 1024)
	originPtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(idPtr))
	defer C.free(unsafe.Pointer(originPtr))
	var buf *C.GstBuffer
	C.gst_event_parse_protection(
		e.Instance(),
		(**C.gchar)(unsafe.Pointer(&idPtr)),
		&buf,
		(**C.gchar)(unsafe.Pointer(&originPtr)),
	)
	return C.GoString((*C.char)(unsafe.Pointer(idPtr))), wrapBuffer(buf), C.GoString((*C.char)(unsafe.Pointer(originPtr)))
}

// ParseQOS gets the type, proportion, diff and timestamp in the qos event. See NewQOSEvent for more information about
// the different QoS values.
//
// timestamp will be adjusted for any pad offsets of pads it was passing through.
func (e *Event) ParseQOS() (qTtype QOSType, proportion float64, diff ClockTimeDiff, timestamp time.Duration) {
	var gtype C.GstQOSType
	var gprop C.gdouble
	var gdiff C.GstClockTimeDiff
	var gts C.GstClockTime
	C.gst_event_parse_qos(
		e.Instance(),
		&gtype, &gprop, &gdiff, &gts,
	)
	return QOSType(gtype), float64(gprop), ClockTimeDiff(gdiff), time.Duration(gts)
}

// ParseSeek parses a seek event.
func (e *Event) ParseSeek() (rate float64, format Format, flags SeekFlags, startType SeekType, start int64, stopType SeekType, stop int64) {
	var grate C.gdouble
	var gformat C.GstFormat
	var gflags C.GstSeekFlags
	var gstartType, gstopType C.GstSeekType
	var gstart, gstop C.gint64
	C.gst_event_parse_seek(e.Instance(), &grate, &gformat, &gflags, &gstartType, &gstart, &gstopType, &gstop)
	return float64(grate), Format(gformat), SeekFlags(gflags), SeekType(gstartType), int64(gstart), SeekType(gstopType), int64(gstop)
}

// ParseSeekTrickModeInterval retrieves the trickmode interval that may have been set on a seek event with
// SetSeekTrickModeInterval.
func (e *Event) ParseSeekTrickModeInterval() (interval time.Duration) {
	var out C.GstClockTime
	C.gst_event_parse_seek_trickmode_interval(e.Instance(), &out)
	return time.Duration(out)
}

// ParseSegment parses a segment event and stores the result in the given segment location. segment remains valid
// only until the event is freed. Don't modify the segment and make a copy if you want to modify it or store it for
// later use.
func (e *Event) ParseSegment() *Segment {
	var out *C.GstSegment
	C.gst_event_parse_segment(e.Instance(), &out)
	return wrapSegment(out)
}

// ParseSegmentDone extracts the position and format from the segment done message.
func (e *Event) ParseSegmentDone() (Format, int64) {
	var format C.GstFormat
	var pos C.gint64
	C.gst_event_parse_segment_done(e.Instance(), &format, &pos)
	return Format(format), int64(pos)
}

// ParseSelectStreams parses the SELECT_STREAMS event and retrieve the contained streams.
func (e *Event) ParseSelectStreams() []*Stream {
	var outList *C.GList
	C.gst_event_parse_select_streams(e.Instance(), &outList)
	return glistToStreamSlice(outList)
}

// ParseSinkMessage parses the sink-message event. Unref msg after usage.
func (e *Event) ParseSinkMessage() *Message {
	var msg *C.GstMessage
	C.gst_event_parse_sink_message(e.Instance(), &msg)
	return wrapMessage(msg)
}

// ParseStep parses a step message
func (e *Event) ParseStep() (format Format, amount uint64, rate float64, flush, intermediate bool) {
	var gformat C.GstFormat
	var gamount C.guint64
	var grate C.gdouble
	var gflush, gintermediate C.gboolean
	C.gst_event_parse_step(e.Instance(), &gformat, &gamount, &grate, &gflush, &gintermediate)
	return Format(gformat), uint64(gamount), float64(grate), gobool(gflush), gobool(gintermediate)
}

// ParseStream parses a stream-start event and extract the GstStream from it.
func (e *Event) ParseStream() *Stream {
	var stream *C.GstStream
	C.gst_event_parse_stream(e.Instance(), &stream)
	return FromGstStreamUnsafeFull(unsafe.Pointer(stream))
}

// ParseStreamCollection parses a stream collection from the event.
func (e *Event) ParseStreamCollection() *StreamCollection {
	stream := &C.GstStreamCollection{}
	C.gst_event_parse_stream_collection(e.Instance(), &stream)
	return FromGstStreamCollectionUnsafeFull(unsafe.Pointer(stream))
}

// ParseStreamFlags parses the stream flags from an event.
func (e *Event) ParseStreamFlags() StreamFlags {
	var out C.GstStreamFlags
	C.gst_event_parse_stream_flags(e.Instance(), &out)
	return StreamFlags(out)
}

// ParseStreamGroupDone parses a stream-group-done event and store the result in the given group_id location.
func (e *Event) ParseStreamGroupDone() uint {
	var out C.guint
	C.gst_event_parse_stream_group_done(e.Instance(), &out)
	return uint(out)
}

// ParseStreamStart parses a stream-id event and store the result in the given stream_id location.
func (e *Event) ParseStreamStart() string {
	idPtr := C.malloc(C.sizeof_char * 1024)
	C.gst_event_parse_stream_start(e.Instance(), (**C.gchar)(unsafe.Pointer(&idPtr)))
	return C.GoString((*C.char)(unsafe.Pointer(idPtr)))
}

// ParseTag parses a tag event and stores the results in the given taglist location. Do not modify or free the returned
// tag list.
func (e *Event) ParseTag() *TagList {
	var out *C.GstTagList
	C.gst_event_parse_tag(e.Instance(), &out)
	return FromGstTagListUnsafeNone(unsafe.Pointer(out))
}

// ParseTOC parses a TOC event and store the results in the given toc and updated locations.
func (e *Event) ParseTOC() (toc *TOC, updated bool) {
	var out *C.GstToc
	var gupdated C.gboolean
	C.gst_event_parse_toc(e.Instance(), &out, &gupdated)
	return FromGstTOCUnsafeFull(unsafe.Pointer(out)), gobool(gupdated)
}

// ParseTOCSelect parses a TOC select event and store the results in the given uid location.
func (e *Event) ParseTOCSelect() string {
	idPtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(idPtr))
	C.gst_event_parse_toc_select(e.Instance(), (**C.gchar)(unsafe.Pointer(&idPtr)))
	return C.GoString((*C.char)(unsafe.Pointer(idPtr)))
}

// Ref increases the ref count on the event by one.
func (e *Event) Ref() *Event {
	C.gst_event_ref(e.Instance())
	return e
}

// SetGapFlags sets flags on event to give additional information about the reason for the GST_EVENT_GAP.
// func (e *Event) SetGapFlags(flags GapFlags) {
// 	C.gst_event_set_gap_flags(e.Instance(), C.GstGapFlags(flags))
// }

// NextGroupID returns a new group id that can be used for an event.
func NextGroupID() uint {
	return uint(C.gst_util_group_id_next())
}

// SetGroupID sets the group id for the stream. All streams that have the same group id are supposed to be played
// together, i.e. all streams inside a container file should have the same group id but different stream ids. The
// group id should change each time the stream is started, resulting in different group ids each time a file is
// played for example.
//
// Use NextGroupID to get a new group id.
func (e *Event) SetGroupID(id uint) {
	C.gst_event_set_group_id(e.Instance(), C.guint(id))
}

// SetRunningTimeOffset sets the running time offset of a event. See GetRunningTimeOffset for more information.
func (e *Event) SetRunningTimeOffset(offset int64) {
	C.gst_event_set_running_time_offset(e.Instance(), C.gint64(offset))
}

// SetSeekTrickModeInterval sets a trickmode interval on a (writable) seek event. Elements that support TRICKMODE_KEY_UNITS
// seeks SHOULD use this as the minimal interval between each frame they may output.
func (e *Event) SetSeekTrickModeInterval(interval time.Duration) {
	C.gst_event_set_seek_trickmode_interval(e.Instance(), C.GstClockTime(interval.Nanoseconds()))
}

// SetSeqnum sets the sequence number of a event.
//
// This function might be called by the creator of a event to indicate that the event relates to other events or messages.
// See GetSeqnum for more information.
func (e *Event) SetSeqnum(seqnum uint32) {
	C.gst_event_set_seqnum(e.Instance(), C.guint32(seqnum))
}

// SetStream sets the stream on the stream-start event
func (e *Event) SetStream(stream *Stream) {
	C.gst_event_set_stream(e.Instance(), stream.Instance())
}

// SetStreamFlags sets the stream flags on the event.
func (e *Event) SetStreamFlags(flags StreamFlags) {
	C.gst_event_set_stream_flags(e.Instance(), C.GstStreamFlags(flags))
}

// Unref decreases the refcount of an event, freeing it if the refcount reaches 0.
func (e *Event) Unref() { C.gst_event_unref(e.Instance()) }

// WritableStructure returns a writable version of the structure.
func (e *Event) WritableStructure() *Structure {
	return wrapStructure(C.gst_event_writable_structure(e.Instance()))
}
