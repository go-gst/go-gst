package gstvideo

// #cgo pkg-config: gstreamer-video-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/video/video.h>
import "C"

// SetFramerate sets the framerate of the video info as a fraction of
// denom/num in frames per second.
func (info *VideoInfo) SetFramerate(denom, num int) {
	info.videoInfo.native.fps_d = C.gint(denom)
	info.videoInfo.native.fps_n = C.gint(num)
}

// SetFramerate sets the framerate of the video info as a fraction of
// denom/num in frames per second.
func (info *VideoInfo) GetSize() int {
	return int(info.videoInfo.native.size)
}
