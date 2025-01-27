package video

/*
#include <gst/gst.h>
#include <gst/video/video.h>
*/
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-gst/gst"
)

// VideoMeta is a Go representation of a GstVideoMeta.
type VideoMeta struct {
	ptr *C.GstVideoMeta
}

// FromGstVideoMetaUnsafeNone wraps the given videoMeta in an VideoMeta instance.
func FromGstVideoMetaUnsafe(videoMeta unsafe.Pointer) *VideoMeta {
	return &VideoMeta{ptr: (*C.GstVideoMeta)(videoMeta)}
}

func BufferAddVideoMetaFull(buffer *gst.Buffer, flags FrameFlags, format Format, width, height uint, offset []uint64, stride []int) *VideoMeta {
	if len(offset) != len(stride) {
		panic("different num planes for offset and stride")
	}

	n_planes := len(offset)

	cOffset := [4]C.gsize{}
	cStride := [4]C.gint{}

	for i := range n_planes {
		cOffset[i] = C.gsize(offset[i])
		cStride[i] = C.gint(stride[i])
	}

	return FromGstVideoMetaUnsafe(unsafe.Pointer(C.gst_buffer_add_video_meta_full(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
		C.GstVideoFrameFlags(flags),
		C.GstVideoFormat(format),
		C.guint(width),
		C.guint(height),
		C.guint(n_planes),
		&cOffset[0],
		&cStride[0],
	)))
}

// CropMetaInfo contains extra buffer metadata describing image cropping.
type CropMetaInfo struct {
	ptr *C.GstVideoCropMeta
}

// GetCropMetaInfo returns the default CropMetaInfo.
func GetCropMetaInfo() *CropMetaInfo {
	meta := C.gst_video_crop_meta_get_info()
	return &CropMetaInfo{(*C.GstVideoCropMeta)(unsafe.Pointer(meta))}
}

// Instance returns the underlying C GstVideoCropMeta instance.
func (c *CropMetaInfo) Instance() *C.GstVideoCropMeta {
	return c.ptr
}

// Meta returns the parent Meta instance.
func (c *CropMetaInfo) Meta() *gst.Meta {
	meta := c.Instance().meta
	return gst.FromGstMetaUnsafe(unsafe.Pointer(&meta))
}

// X returns the horizontal offset.
func (c *CropMetaInfo) X() uint { return uint(c.Instance().x) }

// Y returns the vertical offset.
func (c *CropMetaInfo) Y() uint { return uint(c.Instance().y) }

// Width returns the cropped width.
func (c *CropMetaInfo) Width() uint { return uint(c.Instance().width) }

// Height returns the cropped height.
func (c *CropMetaInfo) Height() uint { return uint(c.Instance().height) }
