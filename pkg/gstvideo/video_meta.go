package gstvideo

// #cgo pkg-config: gstreamer-video-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/video/video.h>
import "C"

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
