package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Stream is a Go representation of a GstStream.
type Stream struct{ *Object }

// FromGstStreamUnsafeNone captures a pointer with a ref and finalizer.
func FromGstStreamUnsafeNone(stream unsafe.Pointer) *Stream {
	return &Stream{wrapObject(glib.TransferNone(stream))}
}

// FromGstStreamUnsafeFull captures a pointer with just a finalizer.
func FromGstStreamUnsafeFull(stream unsafe.Pointer) *Stream {
	return &Stream{wrapObject(glib.TransferNone(stream))}
}

// NewStream returns a new Stream with the given ID, caps, type, and flags.
func NewStream(id string, caps *Caps, sType StreamType, flags StreamFlags) *Stream {
	cID := C.CString(id)
	defer C.free(unsafe.Pointer(cID))
	stream := C.gst_stream_new(cID, caps.Instance(), C.GstStreamType(sType), C.GstStreamFlags(flags))
	return FromGstStreamUnsafeFull(unsafe.Pointer(stream))
}

// Instance returns the underlying GstStream.
func (s *Stream) Instance() *C.GstStream {
	return C.toGstStream(s.Unsafe())
}

// Caps returns the caps for this stream.
func (s *Stream) Caps() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_stream_get_caps(s.Instance())))
}

// StreamFlags returns the flags for this stream.
func (s *Stream) StreamFlags() StreamFlags {
	return StreamFlags(C.gst_stream_get_stream_flags(s.Instance()))
}

// StreamID returns the id of this stream.
func (s *Stream) StreamID() string {
	return C.GoString(C.gst_stream_get_stream_id(s.Instance()))
}

// StreamType returns the type of this stream.
func (s *Stream) StreamType() StreamType {
	return StreamType(C.gst_stream_get_stream_type(s.Instance()))
}

// Tags returns the tag list for this stream.
func (s *Stream) Tags() *TagList {
	tags := C.gst_stream_get_tags(s.Instance())
	if tags == nil {
		return nil
	}
	return FromGstTagListUnsafeFull(unsafe.Pointer(tags))
}

// SetCaps sets the caps for this stream.
func (s *Stream) SetCaps(caps *Caps) {
	C.gst_stream_set_caps(s.Instance(), caps.Instance())
}

// SetStreamFlags sets the flags for this stream.
func (s *Stream) SetStreamFlags(flags StreamFlags) {
	C.gst_stream_set_stream_flags(s.Instance(), C.GstStreamFlags(flags))
}

// SetStreamType sets the type of this stream.
func (s *Stream) SetStreamType(sType StreamType) {
	C.gst_stream_set_stream_type(s.Instance(), C.GstStreamType(sType))
}

// SetTags sets the tags for this stream.
func (s *Stream) SetTags(tags *TagList) {
	C.gst_stream_set_tags(s.Instance(), tags.Instance())
}
