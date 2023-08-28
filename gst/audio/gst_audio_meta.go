package audio

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

// AudioMeta is a Go representation of a GstAudioMeta.
type AudioMeta struct{ *Object }

// FromGstAudioMetaUnsafeNone wraps the given audioMeta with a ref and finalizer.
func FromGstAudioMetaUnsafeNone(audioMeta unsafe.Pointer) *AudioMeta {
	return &AudioMeta{wrapObject(glib.TransferNone(audioMeta))}
}

func BufferAddAudioMeta(buffer *gst.Buffer, info *Info, samples int, offsets []int) *AudioMeta {
	// offsets is not yet implemented, always pass `nil` or this will panic
	if offsets != nil {
		panic("offsets is not implemented")
	}

	// gSizeOffsets := C.gsize(unsafe.Sizeof(unsafe.Pointer(offsets)))
	// if you pass in NULL as the last param then gstreamer assumes things are tightly packed...
	// so that's what we currently assume until we inplement offsets

	return FromGstAudioMetaUnsafeNone(unsafe.Pointer(C.gst_buffer_add_audio_meta(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
		info.ptr,
		C.gsize(samples),
		offsets,
	)))
}
