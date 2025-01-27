package video

// #include <gst/video/video.h>
import "C"

// FrameFlags represents flags for video frames.
type FrameFlags int

// Type castings
const (
	FrameFlagNone          FrameFlags = C.GST_VIDEO_FRAME_FLAG_NONE            // (0)  – no flags
	FrameFlagInterlaced    FrameFlags = C.GST_VIDEO_FRAME_FLAG_INTERLACED      // (1)  – The video frame is interlaced. In mixed interlace-mode, this flag specifies if the frame is interlaced or progressive.
	FrameFlagTTF           FrameFlags = C.GST_VIDEO_FRAME_FLAG_TFF             // (2)  – The video frame has the top field first
	FrameFlagRFF           FrameFlags = C.GST_VIDEO_FRAME_FLAG_RFF             // (4)  – The video frame has the repeat flag
	FrameFlagOneField      FrameFlags = C.GST_VIDEO_FRAME_FLAG_ONEFIELD        // (8)  – The video frame has one field
	FrameFlagMultipleView  FrameFlags = C.GST_VIDEO_FRAME_FLAG_MULTIPLE_VIEW   // (16) – The video contains one or more non-mono views
	FrameFlagFirstInBundle FrameFlags = C.GST_VIDEO_FRAME_FLAG_FIRST_IN_BUNDLE // (32) – The video frame is the first in a set of corresponding views provided as sequential frames.
	FrameFlagTopField      FrameFlags = C.GST_VIDEO_FRAME_FLAG_TOP_FIELD       // (10) – The video frame has the top field only. This is the same as GST_VIDEO_FRAME_FLAG_TFF | GST_VIDEO_FRAME_FLAG_ONEFIELD (Since: 1.16).
	FrameFlagBottonField   FrameFlags = C.GST_VIDEO_FRAME_FLAG_BOTTOM_FIELD    // (8)  – The video frame has the bottom field only. This is the same as GST_VIDEO_FRAME_FLAG_ONEFIELD (GST_VIDEO_FRAME_FLAG_TFF flag unset) (Since: 1.16).
)
