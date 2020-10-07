package video

/*
#include <gst/gst.h>
#include <gst/video/video.h>
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

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
