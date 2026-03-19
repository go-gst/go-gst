package gst

import (
	"runtime"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// Map wraps gst_buffer_map
//
// Users should call [MapInfo.Unmap] or [MapInfo.Close] when done with the buffer
func (buffer *Buffer) Map(flags MapFlags) (*MapInfo, bool) {
	var carg0 *C.GstBuffer  // in, none, converted
	var carg2 C.GstMapFlags // in, none, casted
	var carg1 C.GstMapInfo  // out, transfer: none, C Pointers: 0, Name: MapInfo, caller-allocates
	var cret C.gboolean     // return

	carg0 = (*C.GstBuffer)(UnsafeBufferToGlibNone(buffer))
	carg2 = C.GstMapFlags(flags)

	cret = C.gst_buffer_map(carg0, &carg1, carg2)
	runtime.KeepAlive(buffer)
	runtime.KeepAlive(flags)

	var info *MapInfo
	var goret bool

	info = &MapInfo{
		mapInfo: &mapInfo{
			native: &carg1,
			buffer: buffer,
		},
	}

	info.autoCleanup()

	if cret != 0 {
		goret = true
	}

	return info, goret
}

// PTS returns the presentation timestamp of the buffer.
// It can be GST_CLOCK_TIME_NONE when the PTS is not known or relevant.
func (buffer *Buffer) PTS() ClockTime {
	return ClockTime(buffer.buffer.native.pts)
}

// SetPTS sets the presentation timestamp of the buffer.
// Use GST_CLOCK_TIME_NONE if the PTS is not known or relevant.
func (buffer *Buffer) SetPTS(pts ClockTime) {
	buffer.buffer.native.pts = C.GstClockTime(pts)
}

// DTS returns the decoding timestamp of the buffer.
// It can be GST_CLOCK_TIME_NONE when the DTS is not known or relevant.
func (buffer *Buffer) DTS() ClockTime {
	return ClockTime(buffer.buffer.native.dts)
}

// SetDTS sets the decoding timestamp of the buffer.
// Use GST_CLOCK_TIME_NONE if the DTS is not known or relevant.
func (buffer *Buffer) SetDTS(dts ClockTime) {
	buffer.buffer.native.dts = C.GstClockTime(dts)
}

// Duration returns the duration in time of the buffer data.
// It can be GST_CLOCK_TIME_NONE when the duration is not known or relevant.
func (buffer *Buffer) Duration() ClockTime {
	return ClockTime(buffer.buffer.native.duration)
}

// SetDuration sets the duration in time of the buffer data.
// Use GST_CLOCK_TIME_NONE if the duration is not known or relevant.
func (buffer *Buffer) SetDuration(duration ClockTime) {
	buffer.buffer.native.duration = C.GstClockTime(duration)
}

// Offset returns the media-specific offset for the buffer data.
// For video frames, this is the frame number of this buffer.
// For audio samples, this is the offset of the first sample in this buffer.
// For file data or compressed data, this is the byte offset of the first byte in this buffer.
func (buffer *Buffer) Offset() uint64 {
	return uint64(buffer.buffer.native.offset)
}

// SetOffset sets the media-specific offset for the buffer data.
// For video frames, this is the frame number of this buffer.
// For audio samples, this is the offset of the first sample in this buffer.
// For file data or compressed data, this is the byte offset of the first byte in this buffer.
func (buffer *Buffer) SetOffset(offset uint64) {
	buffer.buffer.native.offset = C.guint64(offset)
}

// OffsetEnd returns the last offset contained in this buffer.
// It has the same format as Offset.
func (buffer *Buffer) OffsetEnd() uint64 {
	return uint64(buffer.buffer.native.offset_end)
}

// SetOffsetEnd sets the last offset contained in this buffer.
// It has the same format as Offset.
func (buffer *Buffer) SetOffsetEnd(offsetEnd uint64) {
	buffer.buffer.native.offset_end = C.guint64(offsetEnd)
}
