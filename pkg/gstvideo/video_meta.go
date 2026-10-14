package gstvideo

// #cgo pkg-config: gstreamer-video-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/video/video.h>
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-gst/pkg/gst"
)

// VideoMeta wraps GstVideoMeta.
//
// see also https://gstreamer.freedesktop.org/documentation/video/gstvideometa.html#GstVideoMeta
type VideoMeta struct {
	*videoMeta
}

// videoMeta is the underlying C struct wrapper.
type videoMeta struct {
	native *C.GstVideoMeta
}

// instance returns the underlying C pointer. This is used by the bindings internally.
func (v *VideoMeta) instance() *C.GstVideoMeta {
	if v == nil {
		return nil
	}
	return v.native
}

// UnsafeVideoMetaFromGlibBorrow is used to convert raw C.GstVideoMeta pointers to go. This is used by the bindings internally.
func UnsafeVideoMetaFromGlibBorrow(p unsafe.Pointer) *VideoMeta {
	if p == nil {
		return nil
	}
	return &VideoMeta{&videoMeta{(*C.GstVideoMeta)(p)}}
}

// UnsafeVideoMetaFromGlibNone is used to convert raw C.GstVideoMeta pointers to go without transferring ownership. This is used by the bindings internally.
// VideoMeta is owned by the underlying GstBuffer and cannot be referenced or freed independently.
func UnsafeVideoMetaFromGlibNone(p unsafe.Pointer) *VideoMeta {
	return UnsafeVideoMetaFromGlibBorrow(p)
}

// UnsafeVideoMetaFromGlibFull is used to convert raw C.GstVideoMeta pointers to go while taking ownership. This is used by the bindings internally.
func UnsafeVideoMetaFromGlibFull(p unsafe.Pointer) *VideoMeta {
	return UnsafeVideoMetaFromGlibBorrow(p)
}

// UnsafeVideoMetaFree unrefs/frees the underlying resource.
func UnsafeVideoMetaFree(v *VideoMeta) {
	if v != nil {
		v.native = nil
	}
}

// UnsafeVideoMetaToGlibNone returns the underlying C pointer. This is used by the bindings internally.
func UnsafeVideoMetaToGlibNone(v *VideoMeta) unsafe.Pointer {
	if v == nil {
		return nil
	}
	return unsafe.Pointer(v.native)
}

// UnsafeVideoMetaToGlibFull returns the underlying C pointer and gives up ownership.
// This is used by the bindings internally.
func UnsafeVideoMetaToGlibFull(v *VideoMeta) unsafe.Pointer {
	if v == nil {
		return nil
	}
	_p := unsafe.Pointer(v.native)
	v.native = nil
	return _p
}

// GetWidth returns the video width.
func (meta *VideoMeta) GetWidth() uint32 {
	return uint32(meta.videoMeta.native.width)
}

// GetHeight returns the video height.
func (meta *VideoMeta) GetHeight() uint32 {
	return uint32(meta.videoMeta.native.height)
}

// GetNPlanes returns the number of planes in the image.
func (meta *VideoMeta) GetNPlanes() int {
	return int(meta.videoMeta.native.n_planes)
}

// GetOffset returns the offsets for the planes.
func (meta *VideoMeta) GetOffset() []uint64 {
	result := make([]uint64, len(meta.videoMeta.native.offset))
	for i := range meta.GetNPlanes() {
		result[i] = uint64(meta.videoMeta.native.offset[i])
	}
	return result
}

// GetStride returns the strides for the planes.
func (meta *VideoMeta) GetStride() []uint64 {
	result := make([]uint64, len(meta.videoMeta.native.stride))
	for i := range meta.GetNPlanes() {
		result[i] = uint64(meta.videoMeta.native.stride[i])
	}
	return result
}

// VideoMetaGetInfo wraps gst_video_meta_get_info.
//
// see also https://gstreamer.freedesktop.org/documentation/video/gstvideometa.html#gst_video_meta_get_info
func VideoMetaGetInfo() *gst.MetaInfo {
	var cret *C.GstMetaInfo
	cret = C.gst_video_meta_get_info()
	return gst.UnsafeMetaInfoFromGlibNone(unsafe.Pointer(cret))
}

