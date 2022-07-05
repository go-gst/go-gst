package audio

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

type AudioMeta struct {
	ptr *C.GstAudioMeta
}

func wrapMetaFull(ptr *C.GstAudioMeta) *AudioMeta {
	meta := &AudioMeta{ptr}
	return meta
}

func BufferAddAudioMeta(buffer *gst.Buffer, info *Info, samples int /*, offsets *[]int*/) *AudioMeta {
	/*gSizeOffsets := C.gsize(unsafe.Sizeof(unsafe.Pointer(offsets)))*/
	// if you pass in NULL as the last param then gstreamer assumes things are tightly packed...
	// and they are for me so this makes things work... but obviously would be nice for it to work
	// properly

	return wrapMetaFull(C.gst_buffer_add_audio_meta(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
		info.ptr,
		C.gsize(samples),
		/*&gSizeOffsets,*/
		nil,
	))
}
