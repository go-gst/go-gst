package gst

// #include "gst.go.h"
import "C"
import "unsafe"

// Segment is a go wrapper around a GstSegment.
// See: https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#GstSegment
type Segment struct {
	ptr *C.GstSegment
}

// FromGstSegmentUnsafe wraps the GstSegment pointer.
func FromGstSegmentUnsafe(segment unsafe.Pointer) *Segment {
	return wrapSegment((*C.GstSegment)(segment))
}

// ToGstSegment converts the given pointer into a Segment without affecting the ref count or
// placing finalizers.
func ToGstSegment(segment unsafe.Pointer) *Segment {
	return wrapSegment((*C.GstSegment)(segment))
}

// NewSegment allocates and initializes a new Segment. Free when you are finished.
func NewSegment() *Segment {
	return wrapSegment(C.gst_segment_new())
}

// NewFormattedSegment returns a new Segment with the given format.
func NewFormattedSegment(f Format) *Segment {
	segment := NewSegment()
	segment.SetFormat(f)
	return segment
}

// Instance returns the underlying GstSegment instance.
func (s *Segment) Instance() *C.GstSegment { return s.ptr }

// GetFlags returns the flags on this segment.
func (s *Segment) GetFlags() SegmentFlags { return SegmentFlags(s.Instance().flags) }

// SetFlags sets the flags on this segment.
func (s *Segment) SetFlags(flags SegmentFlags) { s.Instance().flags = C.GstSegmentFlags(flags) }

// GetRate returns the rate for this segment.
func (s *Segment) GetRate() float64 { return float64(s.Instance().rate) }

// SetRate sets the rate for this segment.
func (s *Segment) SetRate(rate float64) { s.Instance().rate = C.gdouble(rate) }

// GetAppliedRate returns the applied rate for this segment.
func (s *Segment) GetAppliedRate() float64 { return float64(s.Instance().applied_rate) }

// SetAppliedRate sets the applied rate for this segment
func (s *Segment) SetAppliedRate(rate float64) { s.Instance().applied_rate = C.gdouble(rate) }

// GetFormat returns the format for this segment.
func (s *Segment) GetFormat() Format { return Format(s.Instance().format) }

// SetFormat sets the format on this segment.
func (s *Segment) SetFormat(f Format) { s.Instance().format = C.GstFormat(f) }

// GetBase returns the base for this segment.
func (s *Segment) GetBase() uint64 { return uint64(s.Instance().base) }

// GetOffset returns the offset for this segment.
func (s *Segment) GetOffset() uint64 { return uint64(s.Instance().offset) }

// GetStart returns the start of this segment.
func (s *Segment) GetStart() uint64 { return uint64(s.Instance().start) }

// GetStop returns the stop of this segment.
func (s *Segment) GetStop() uint64 { return uint64(s.Instance().stop) }

// GetTime returns the time of this segment.
func (s *Segment) GetTime() uint64 { return uint64(s.Instance().time) }

// GetPosition returns the position of this segment.
func (s *Segment) GetPosition() uint64 { return uint64(s.Instance().position) }

// GetDuration gets the duration of this segment.
func (s *Segment) GetDuration() uint64 { return uint64(s.Instance().duration) }

// Clip clips the given start and stop values to the segment boundaries given in segment. start and stop are compared and clipped
// to segment start and stop values.
//
// If the function returns FALSE, start and stop are known to fall outside of segment and clip_start and clip_stop are not updated.
//
// When the function returns TRUE, clip_start and clip_stop will be updated. If clip_start or clip_stop are different from start or stop
// respectively, the region fell partially in the segment.
//
// Note that when stop is -1, clip_stop will be set to the end of the segment. Depending on the use case, this may or may not be what you want.
func (s *Segment) Clip(format Format, start, stop uint64) (ok bool, clipStart, clipStop uint64) {
	var gclipStart, gclipStop C.guint64
	gok := C.gst_segment_clip(
		s.Instance(),
		C.GstFormat(format),
		C.guint64(start),
		C.guint64(stop),
		&gclipStart, &gclipStop,
	)
	return gobool(gok), uint64(gclipStart), uint64(gclipStop)
}

// Copy creates a copy of this segment.
func (s *Segment) Copy() *Segment { return wrapSegment(C.gst_segment_copy(s.Instance())) }

// CopyInto copies the contents of this segment into the given one.
func (s *Segment) CopyInto(segment *Segment) {
	C.gst_segment_copy_into(s.Instance(), segment.Instance())
}

// DoSeek updates the segment structure with the field values of a seek event (see NewSeekEvent).
//
// After calling this method, the segment field position and time will contain the requested new position in the segment.
// The new requested position in the segment depends on rate and start_type and stop_type.
//
// For positive rate, the new position in the segment is the new segment start field when it was updated with a start_type
// different from SeekTypeNone. If no update was performed on segment start position (#SeekTypeNone), start is ignored and
// segment position is unmodified.
//
// For negative rate, the new position in the segment is the new segment stop field when it was updated with a stop_type different
// from SeekTypeNone. If no stop was previously configured in the segment, the duration of the segment will be used to update the
// stop position. If no update was performed on segment stop position (#SeekTypeNone), stop is ignored and segment position
// is unmodified.
//
// The applied rate of the segment will be set to 1.0 by default. If the caller can apply a rate change, it should update segment
// rate and applied_rate after calling this function.
//
// update will be set to TRUE if a seek should be performed to the segment position field. This field can be FALSE if, for example,
// only the rate has been changed but not the playback position.
func (s *Segment) DoSeek(rate float64, format Format, flags SeekFlags, startType SeekType, start uint64, stopType SeekType, stop uint64) (ok, update bool) {
	var gupdate C.gboolean
	gok := C.gst_segment_do_seek(
		s.Instance(),
		C.gdouble(rate),
		C.GstFormat(format),
		C.GstSeekFlags(flags),
		C.GstSeekType(startType),
		C.guint64(start),
		C.GstSeekType(stopType),
		C.guint64(stop),
		&gupdate,
	)
	return gobool(gok), gobool(gupdate)
}

// Free frees the allocated segment.
func (s *Segment) Free() { C.gst_segment_free(s.Instance()) }

// Init reinitializes a segment to its default values.
func (s *Segment) Init(format Format) {
	C.gst_segment_init(s.Instance(), C.GstFormat(format))
}

// IsEqual checks for two segments being equal. Equality here is defined as perfect equality, including floating point values.
func (s *Segment) IsEqual(segment *Segment) bool {
	return gobool(C.gst_segment_is_equal(
		s.Instance(), segment.Instance(),
	))
}

// OffsetRunningTime adjusts the values in segment so that offset is applied to all future running-time calculations.
func (s *Segment) OffsetRunningTime(format Format, offset int64) bool {
	return gobool(C.gst_segment_offset_running_time(
		s.Instance(),
		C.GstFormat(format),
		C.gint64(offset),
	))
}

// PositionFromRunningTime converts running_time into a position in the segment so that ToRunningTime with that position returns
// running_time. The position in the segment for runningTime is returned.
func (s *Segment) PositionFromRunningTime(format Format, runningTime uint64) uint64 {
	return uint64(C.gst_segment_position_from_running_time(s.Instance(), C.GstFormat(format), C.guint64(runningTime)))
}

// PositionFromRunningTimeFull translates running_time to the segment position using the currently configured segment. Compared to
// PositionFromRunningTime this function can return negative segment position.
//
// This function is typically used by elements that need to synchronize buffers against the clock or each other.
//
// running_time can be any value and the result of this function for values outside of the segment is extrapolated.
func (s *Segment) PositionFromRunningTimeFull(format Format, runningTime uint64) int64 {
	var position C.guint64
	ret := C.gst_segment_position_from_running_time_full(s.Instance(), C.GstFormat(format), C.guint64(runningTime), &position)
	if int(ret) > 0 {
		return int64(position)
	}
	return int64(position) * -1
}

// PositionFromStreamTime converts stream_time into a position in the segment so that ToStreamTime with that position returns stream_time.
func (s *Segment) PositionFromStreamTime(format Format, streamTime uint64) uint64 {
	return uint64(C.gst_segment_position_from_stream_time(s.Instance(), C.GstFormat(format), C.guint64(streamTime)))
}

// PositionFromStreamTimeFull translates stream_time to the segment position using the currently configured segment. Compared to PositionFromStreamTime
// this function can return negative segment position.
//
// This function is typically used by elements that need to synchronize buffers against the clock or each other.
//
// stream_time can be any value and the result of this function for values outside of the segment is extrapolated.
func (s *Segment) PositionFromStreamTimeFull(format Format, streamTime uint64) int64 {
	var position C.guint64
	ret := C.gst_segment_position_from_stream_time_full(s.Instance(), C.GstFormat(format), C.guint64(streamTime), &position)
	if int(ret) > 0 {
		return int64(position)
	}
	return int64(position) * -1
}

// SetRunningTime adjusts the start/stop and base values of segment such that the next valid buffer will be one with running_time.
func (s *Segment) SetRunningTime(format Format, runningTime uint64) bool {
	return gobool(C.gst_segment_set_running_time(
		s.Instance(),
		C.GstFormat(format),
		C.guint64(runningTime),
	))
}

// ToRunningTime translates position to the total running time using the currently configured segment. Position is a value between
// segment start and stop time.
//
// This function is typically used by elements that need to synchronize to the global clock in a pipeline. The running time is a
// constantly increasing value starting from 0. When segment Init is called, this value will reset to 0.
//
// This function returns -1 if the position is outside of segment start and stop.
func (s *Segment) ToRunningTime(format Format, position uint64) uint64 {
	return uint64(C.gst_segment_to_running_time(s.Instance(), C.GstFormat(format), C.guint64(position)))
}

// ToRunningTimeFull translates position to the total running time using the currently configured segment. Compared to ToRunningTime
// this function can return negative running-time.
//
// This function is typically used by elements that need to synchronize buffers against the clock or each other.
//
// position can be any value and the result of this function for values outside of the segment is extrapolated.
func (s *Segment) ToRunningTimeFull(format Format, position uint64) int64 {
	var runningTime C.guint64
	ret := C.gst_segment_to_running_time_full(s.Instance(), C.GstFormat(format), C.guint64(position), &runningTime)
	if int(ret) > 0 {
		return int64(runningTime)
	}
	return int64(runningTime) * -1
}

// ToStreamTime translates position to stream time using the currently configured segment. The position value must be between segment
// start and stop value.
//
// This function is typically used by elements that need to operate on the stream time of the buffers it receives, such as effect
// plugins. In those use cases, position is typically the buffer timestamp or clock time that one wants to convert to the stream time.
// The stream time is always between 0 and the total duration of the media stream.
func (s *Segment) ToStreamTime(format Format, position uint64) uint64 {
	return uint64(C.gst_segment_to_stream_time(s.Instance(), C.GstFormat(format), C.guint64(position)))
}

// ToStreamTimeFull translates position to the total stream time using the currently configured segment. Compared to ToStreamTime this
// function can return negative stream-time.
//
// This function is typically used by elements that need to synchronize buffers against the clock or each other.
//
// position can be any value and the result of this function for values outside of the segment is extrapolated.
func (s *Segment) ToStreamTimeFull(format Format, position uint64) int64 {
	var streamTime C.guint64
	ret := C.gst_segment_to_running_time_full(s.Instance(), C.GstFormat(format), C.guint64(position), &streamTime)
	if int(ret) > 0 {
		return int64(streamTime)
	}
	return int64(streamTime) * -1
}
