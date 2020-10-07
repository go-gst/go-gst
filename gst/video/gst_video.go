package video

// #include <gst/video/video.h>
import "C"

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// OrientationMethod represents the different video orientation methods.
type OrientationMethod int

// Type castings
const (
	OrientationMethodIdentity OrientationMethod = C.GST_VIDEO_ORIENTATION_IDENTITY // (0) – Identity (no rotation)
	OrientationMethod90R      OrientationMethod = C.GST_VIDEO_ORIENTATION_90R      // (1) – Rotate clockwise 90 degrees
	OrientationMethod180      OrientationMethod = C.GST_VIDEO_ORIENTATION_180      // (2) – Rotate 180 degrees
	OrientationMethod90L      OrientationMethod = C.GST_VIDEO_ORIENTATION_90L      // (3) – Rotate counter-clockwise 90 degrees
	OrientationMethodHoriz    OrientationMethod = C.GST_VIDEO_ORIENTATION_HORIZ    // (4) – Flip horizontally
	OrientationMethodVert     OrientationMethod = C.GST_VIDEO_ORIENTATION_VERT     // (5) – Flip vertically
	OrientationMethodULLR     OrientationMethod = C.GST_VIDEO_ORIENTATION_UL_LR    // (6) – Flip across upper left/lower right diagonal
	OrientationMethodURLL     OrientationMethod = C.GST_VIDEO_ORIENTATION_UR_LL    // (7) – Flip across upper right/lower left diagonal
	OrientationMethodAuto     OrientationMethod = C.GST_VIDEO_ORIENTATION_AUTO     // (8) – Select flip method based on image-orientation tag
	OrientationMethodCustom   OrientationMethod = C.GST_VIDEO_ORIENTATION_CUSTOM   // (9) – Current status depends on plugin internal setup
)

// Additional video meta tags
const (
	TagVideoColorspage  gst.Tag = C.GST_META_TAG_VIDEO_COLORSPACE_STR
	TagVideoOrientation gst.Tag = C.GST_META_TAG_VIDEO_ORIENTATION_STR
	TagVideoSize        gst.Tag = C.GST_META_TAG_VIDEO_SIZE_STR
	TagVideo            gst.Tag = C.GST_META_TAG_VIDEO_STR
)

// Alignment represents parameters for the memory of video buffers. This structure is
// usually used to configure the bufferpool if it supports the BufferPoolOptionVideoAlignment.
type Alignment struct {
	// extra pixels on the top
	PaddingTop uint
	// extra pixels on bottom
	PaddingBottom uint
	// extra pixels on the left
	PaddingLeft uint
	// extra pixels on the right
	PaddingRight uint
}

func (a *Alignment) instance() *C.GstVideoAlignment {
	g := &C.GstVideoAlignment{
		padding_top:    C.guint(a.PaddingTop),
		padding_bottom: C.guint(a.PaddingBottom),
		padding_left:   C.guint(a.PaddingLeft),
		padding_right:  C.guint(a.PaddingRight),
	}
	runtime.SetFinalizer(a, func(_ *Alignment) { C.g_free((C.gpointer)(unsafe.Pointer(g))) })
	return g
}

// CalculateDisplayRatio will, given the Pixel Aspect Ratio and size of an input video frame, and
// the pixel aspect ratio of the intended display device, calculate the actual display ratio the
// video will be rendered with.
//
// See https://gstreamer.freedesktop.org/documentation/video/gstvideo.html?gi-language=c#gst_video_calculate_display_ratio
func CalculateDisplayRatio(videoWidth, videoHeight, videoParNum, videoParDenom, displayParNum, displayParDenom uint) (darNum, darDenom uint, ok bool) {
	var gNum, gDenom C.guint
	gok := C.gst_video_calculate_display_ratio(
		&gNum, &gDenom,
		C.guint(videoWidth), C.guint(videoHeight),
		C.guint(videoParNum), C.guint(videoParDenom),
		C.guint(displayParNum), C.guint(displayParDenom),
	)
	return uint(gNum), uint(gDenom), gobool(gok)
}

// GuessFramerate will, given the nominal duration of one video frame, check some standard framerates
// for a close match (within 0.1%) and return one if possible,
//
// It will calculate an arbitrary framerate if no close match was found, and return FALSE.
//
// It returns FALSE if a duration of 0 is passed.
func GuessFramerate(dur time.Duration) (destNum, destDenom int, ok bool) {
	var num, denom C.gint
	gok := C.gst_video_guess_framerate(durationToClockTime(dur), &num, &denom)
	return int(num), int(denom), gobool(gok)
}
