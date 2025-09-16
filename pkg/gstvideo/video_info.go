package gstvideo

// #cgo pkg-config: gstreamer-video-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/video/video.h>
import "C"

// SetFramerate sets the framerate of the video info as a fraction of
// num/denom in frames per second.
func (info *VideoInfo) SetFramerate(num, denom int) {
	info.videoInfo.native.fps_n = C.gint(num)
	info.videoInfo.native.fps_d = C.gint(denom)
}

// SetFramerate sets the framerate of the video info as a fraction of
// denom/num in frames per second.
func (info *VideoInfo) GetSize() int {
	return int(info.videoInfo.native.size)
}
