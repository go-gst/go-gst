package gst

import (
	"unsafe"
)

// #cgo pkg-config: gstreamer-1.0
// #include <gst/gst.h>
import "C"

// TakeBuffer removes the buffer from the probe info and returns it. The caller owns the returned
// buffer, so unlike the one returned by GetBuffer it stays valid after the probe callback returns.
// The buffer is dropped from the data flow unless a new one is put back with SetBuffer.
//
// Returns nil if the probe info does not hold a buffer.
func (info *PadProbeInfo) TakeBuffer() *Buffer {
	return UnsafeBufferFromGlibFull(info.take(C.GST_PAD_PROBE_TYPE_BUFFER))
}

// TakeBufferList removes the buffer list from the probe info and returns it. The caller owns the
// returned list, so unlike the one returned by GetBufferList it stays valid after the probe callback
// returns. The list is dropped from the data flow unless a new one is put back with SetBufferList.
//
// Returns nil if the probe info does not hold a buffer list.
func (info *PadProbeInfo) TakeBufferList() *BufferList {
	return UnsafeBufferListFromGlibFull(info.take(C.GST_PAD_PROBE_TYPE_BUFFER_LIST))
}

// TakeEvent removes the event from the probe info and returns it. The caller owns the returned
// event, so unlike the one returned by GetEvent it stays valid after the probe callback returns.
// The event is dropped from the data flow unless a new one is put back with SetEvent.
//
// Returns nil if the probe info does not hold an event.
func (info *PadProbeInfo) TakeEvent() *Event {
	return UnsafeEventFromGlibFull(info.take(C.GST_PAD_PROBE_TYPE_EVENT_BOTH))
}

// SetBuffer replaces the buffer in the probe info and takes ownership of buffer, so buffer must not
// be used afterwards. Passing nil drops the buffer from the data flow. The probe info, and any data
// returned by its getters, is only valid until the probe callback returns.
//
// Panics if the probe info is not for a buffer.
func (info *PadProbeInfo) SetBuffer(buffer *Buffer) {
	info.checkType(C.GST_PAD_PROBE_TYPE_BUFFER, "buffer")
	info.set(UnsafeBufferToGlibFull(buffer))
}

// SetBufferList replaces the buffer list in the probe info and takes ownership of list, so list
// must not be used afterwards. Passing nil drops the list from the data flow. The probe info, and
// any data returned by its getters, is only valid until the probe callback returns.
//
// Panics if the probe info is not for a buffer list.
func (info *PadProbeInfo) SetBufferList(list *BufferList) {
	info.checkType(C.GST_PAD_PROBE_TYPE_BUFFER_LIST, "buffer list")
	info.set(UnsafeBufferListToGlibFull(list))
}

// SetEvent replaces the event in the probe info and takes ownership of event, so event must not be
// used afterwards. Passing nil drops the event from the data flow. The probe info, and any data
// returned by its getters, is only valid until the probe callback returns.
//
// Panics if the probe info is not for an event.
func (info *PadProbeInfo) SetEvent(event *Event) {
	info.checkType(C.GST_PAD_PROBE_TYPE_EVENT_BOTH, "event")
	info.set(UnsafeEventToGlibFull(event))
}

func (info *PadProbeInfo) take(mask C.GstPadProbeType) unsafe.Pointer {
	if info.native._type&mask == 0 {
		return nil
	}
	data := info.native.data
	info.native.data = nil
	return unsafe.Pointer(data)
}

func (info *PadProbeInfo) checkType(mask C.GstPadProbeType, name string) {
	if info.native._type&mask == 0 {
		panic("gst: PadProbeInfo does not hold a " + name)
	}
}

// set stores data, which must carry a reference owned by the probe info.
func (info *PadProbeInfo) set(data unsafe.Pointer) {
	// A borrowed getter result passed back in is the same object and carries no reference of its own.
	if old := unsafe.Pointer(info.native.data); old != nil && old != data {
		C.gst_mini_object_unref((*C.GstMiniObject)(old))
	}
	info.native.data = C.gpointer(data)
}
