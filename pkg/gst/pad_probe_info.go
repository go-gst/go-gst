package gst

import (
	"runtime"
	"unsafe"
)

// #cgo pkg-config: gstreamer-1.0
// #include <gst/gst.h>
import "C"

// BorrowBuffer returns the buffer the probe was called with, or nil if the probe is not
// for a buffer. Unlike [PadProbeInfo.GetBuffer] it takes no reference, so a buffer the
// pad holds the only reference to stays writable.
//
// The buffer is only valid until the probe callback returns.
func (info *PadProbeInfo) BorrowBuffer() *Buffer {
	defer runtime.KeepAlive(info)
	if info.native._type&C.GST_PAD_PROBE_TYPE_BUFFER == 0 {
		return nil
	}
	return UnsafeBufferFromGlibBorrow(unsafe.Pointer(C.gst_pad_probe_info_get_buffer(info.native)))
}

// MakeBufferWritable makes the probe's buffer writable, copying it if it is shared, and
// stores the result in the probe info so it is what the pad passes on. It returns the
// writable buffer, or nil if the probe is not for a buffer.
//
// The buffer is only valid until the probe callback returns.
func (info *PadProbeInfo) MakeBufferWritable() *Buffer {
	defer runtime.KeepAlive(info)
	if info.native._type&C.GST_PAD_PROBE_TYPE_BUFFER == 0 {
		return nil
	}
	b := C.gst_pad_probe_info_get_buffer(info.native)
	if b == nil {
		return nil
	}
	w := C.gst_mini_object_make_writable((*C.GstMiniObject)(unsafe.Pointer(b)))
	info.native.data = C.gpointer(w)
	return UnsafeBufferFromGlibBorrow(unsafe.Pointer(w))
}

// BorrowEvent returns the event the probe was called with, or nil if the probe is not
// for an event. Unlike [PadProbeInfo.GetEvent] it takes no reference.
//
// The event is only valid until the probe callback returns.
func (info *PadProbeInfo) BorrowEvent() *Event {
	defer runtime.KeepAlive(info)
	if info.native._type&C.GST_PAD_PROBE_TYPE_EVENT_BOTH == 0 {
		return nil
	}
	return UnsafeEventFromGlibBorrow(unsafe.Pointer(C.gst_pad_probe_info_get_event(info.native)))
}
