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

func BufferAddAudioMeta(buffer *gst.Buffer, info *Info, samples int64, offsets []int) *AudioMeta {
	gSizeOffsets := C.gsize(unsafe.Sizeof(unsafe.Pointer(&offsets)))
	return wrapMetaFull(C.gst_buffer_add_audio_meta(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
		info.ptr,
		C.gsize(samples),
		&gSizeOffsets,
	))
}
