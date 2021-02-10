package video

/*
#include <gst/video/video.h>

GstVideoChromaSite     infoChromaSite        (GstVideoInfo * info)                     { return GST_VIDEO_INFO_CHROMA_SITE(info); }
GstVideoColorimetry    infoColorimetry       (GstVideoInfo * info)                     { return GST_VIDEO_INFO_COLORIMETRY(info); }
gint                   infoFieldHeight       (GstVideoInfo * info)                     { return GST_VIDEO_INFO_FIELD_HEIGHT(info); }
GstVideoFieldOrder     infoFieldOrder        (GstVideoInfo * info)                     { return GST_VIDEO_INFO_FIELD_ORDER(info); }
gint                   infoFieldRateN        (GstVideoInfo * info)                     { return GST_VIDEO_INFO_FIELD_RATE_N(info); }
GstVideoFlags          infoFlags             (GstVideoInfo * info)                     { return GST_VIDEO_INFO_FLAGS(info); }
gboolean               infoFlagIsSet         (GstVideoInfo * info, GstVideoFlags flag) { return GST_VIDEO_INFO_FLAG_IS_SET(info, flag); }
void                   infoFlagSet           (GstVideoInfo * info, GstVideoFlags flag) { GST_VIDEO_INFO_FLAG_SET(info, flag); }
void                   infoFlagUnset         (GstVideoInfo * info, GstVideoFlags flag) { GST_VIDEO_INFO_FLAG_UNSET(info, flag); }
GstVideoFormat         infoFormat            (GstVideoInfo * info)                     { return GST_VIDEO_INFO_FORMAT(info); }
gint                   infoFPSd              (GstVideoInfo * info)                     { return GST_VIDEO_INFO_FPS_D(info); }
gint                   infoFPSn              (GstVideoInfo * info)                     { return GST_VIDEO_INFO_FPS_N(info); }
gboolean               infoHasAlpha          (GstVideoInfo * info)                     { return GST_VIDEO_INFO_HAS_ALPHA(info); }
gint                   infoHeight            (GstVideoInfo * info)                     { return GST_VIDEO_INFO_HEIGHT(info); }
GstVideoInterlaceMode  infoInterlaceMode     (GstVideoInfo * info)                     { return GST_VIDEO_INFO_INTERLACE_MODE(info); }
gboolean               infoIsGray            (GstVideoInfo * info)                     { return GST_VIDEO_INFO_IS_GRAY(info); }
gboolean               infoIsInterlaced      (GstVideoInfo * info)                     { return GST_VIDEO_INFO_IS_INTERLACED(info); }
gboolean               infoIsRGB             (GstVideoInfo * info)                     { return GST_VIDEO_INFO_IS_RGB(info); }
gboolean               infoIsYUV             (GstVideoInfo * info)                     { return GST_VIDEO_INFO_IS_YUV(info); }
GstVideoMultiviewFlags infoMultiviewFlags    (GstVideoInfo * info)                     { return GST_VIDEO_INFO_MULTIVIEW_FLAGS(info); }
GstVideoMultiviewMode  infoMultiviewMode     (GstVideoInfo * info)                     { return GST_VIDEO_INFO_MULTIVIEW_MODE(info); }
const gchar *          infoName              (GstVideoInfo * info)                     { return GST_VIDEO_INFO_NAME(info); }
guint                  infoNComponents       (GstVideoInfo * info)                     { return GST_VIDEO_INFO_N_COMPONENTS(info); }
guint                  infoNPlanes           (GstVideoInfo * info)                     { return GST_VIDEO_INFO_N_PLANES(info); }
gint                   infoPARd              (GstVideoInfo * info)                     { return GST_VIDEO_INFO_PAR_D(info); }
gint                   infoPARn              (GstVideoInfo * info)                     { return GST_VIDEO_INFO_PAR_N(info); }
gsize                  infoSize              (GstVideoInfo * info)                     { return GST_VIDEO_INFO_SIZE(info); }
gint                   infoViews             (GstVideoInfo * info)                     { return GST_VIDEO_INFO_VIEWS(info); }
gint                   infoWidth             (GstVideoInfo * info)                     { return GST_VIDEO_INFO_WIDTH(info); }
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// CapsFeatureFormatInterlaced is the name of the caps feature indicating that the stream is interlaced.
//
// Currently it is only used for video with 'interlace-mode=alternate' to ensure backwards compatibility
// for this new mode. In this mode each buffer carries a single field of interlaced video. BufferFlagTopField
// and BufferFlagBottomField indicate whether the buffer carries a top or bottom field. The order of
// buffers/fields in the stream and the timestamps on the buffers indicate the temporal order of the fields.
// Top and bottom fields are expected to alternate in this mode. The frame rate in the caps still signals the
// frame rate, so the notional field rate will be twice the frame rate from the caps.
const CapsFeatureFormatInterlaced string = C.GST_CAPS_FEATURE_FORMAT_INTERLACED

// FieldOrder is the field order of interlaced content. This is only valid for interlace-mode=interleaved
// and not interlace-mode=mixed. In the case of mixed or FieldOrderrUnknown, the field order is signalled
// via buffer flags.
type FieldOrder int

// Type castings
const (
	FieldOrderUnknown          FieldOrder = C.GST_VIDEO_FIELD_ORDER_UNKNOWN            // (0) – unknown field order for interlaced content. The actual field order is signalled via buffer flags.
	FieldOrderTopFieldFirst    FieldOrder = C.GST_VIDEO_FIELD_ORDER_TOP_FIELD_FIRST    // (1) – top field is first
	FieldOrderBottomFieldFirst FieldOrder = C.GST_VIDEO_FIELD_ORDER_BOTTOM_FIELD_FIRST // (2) – bottom field is first
)

// String implements a stringer on FieldOrder
func (f FieldOrder) String() string {
	cStr := C.gst_video_field_order_to_string(C.GstVideoFieldOrder(f))
	return C.GoString(cStr)
}

// Flags represents extra video flags
type Flags int

// Type castings
const (
	FlagNone               Flags = C.GST_VIDEO_FLAG_NONE                // (0) – no flags
	FlagVariableFPS        Flags = C.GST_VIDEO_FLAG_VARIABLE_FPS        // (1) – a variable fps is selected, fps_n and fps_d denote the maximum fps of the video
	FlagPremultipliedAlpha Flags = C.GST_VIDEO_FLAG_PREMULTIPLIED_ALPHA // (2) – Each color has been scaled by the alpha value.
)

// InterlaceMode is the possible values describing the interlace mode of the stream.
type InterlaceMode int

// Type castings
const (
	InterlaceModeProgressive InterlaceMode = C.GST_VIDEO_INTERLACE_MODE_PROGRESSIVE // (0) – all frames are progressive
	InterlaceModeInterleaved InterlaceMode = C.GST_VIDEO_INTERLACE_MODE_INTERLEAVED // (1) – 2 fields are interleaved in one video frame. Extra buffer flags describe the field order.
	InterlaceModeMixed       InterlaceMode = C.GST_VIDEO_INTERLACE_MODE_MIXED       // (2) – frames contains both interlaced and progressive video, the buffer flags describe the frame and fields.
	InterlaceModeFields      InterlaceMode = C.GST_VIDEO_INTERLACE_MODE_FIELDS      // (3) – 2 fields are stored in one buffer, use the frame ID to get access to the required field. For multiview (the 'views' property > 1) the fields of view N can be found at frame ID (N * 2) and (N * 2) + 1. Each field has only half the amount of lines as noted in the height property. This mode requires multiple GstVideoMeta metadata to describe the fields.
	InterlaceModeAlternate   InterlaceMode = C.GST_VIDEO_INTERLACE_MODE_ALTERNATE   // (4) – 1 field is stored in one buffer, GST_VIDEO_BUFFER_FLAG_TF or GST_VIDEO_BUFFER_FLAG_BF indicates if the buffer is carrying the top or bottom field, respectively. The top and bottom buffers are expected to alternate in the pipeline, with this mode (Since: 1.16).
)

// String implements a stringer on interlace mode
func (i InterlaceMode) String() string {
	return C.GoString(C.gst_video_interlace_mode_to_string(C.GstVideoInterlaceMode(i)))
}

// MultiviewFlags are used to indicate extra properties of a stereo/multiview stream beyond the frame layout
// and buffer mapping that is conveyed in the MultiviewMode.
type MultiviewFlags int

// Type castings
const (
	MultiviewFlagsNone           MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_NONE             // (0) – No flags
	MultiviewFlagsRightViewFirst MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_RIGHT_VIEW_FIRST // (1) – For stereo streams, the normal arrangement of left and right views is reversed.
	MultiviewFlagsLeftFlipped    MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_LEFT_FLIPPED     // (2) – The left view is vertically mirrored.
	MultiviewFlagsLeftFlopped    MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_LEFT_FLOPPED     // (4) – The left view is horizontally mirrored.
	MultiviewFlagsRightFlipped   MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_RIGHT_FLIPPED    // (8) – The right view is vertically mirrored.
	MultiviewFlagsRightFlopped   MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_RIGHT_FLOPPED    // (16) – The right view is horizontally mirrored.
	MultiviewFlagsHalfAspect     MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_HALF_ASPECT      // (16384) – For frame-packed multiview modes, indicates that the individual views have been encoded with half the true width or height and should be scaled back up for display. This flag is used for overriding input layout interpretation by adjusting pixel-aspect-ratio. For side-by-side, column interleaved or checkerboard packings, the pixel width will be doubled. For row interleaved and top-bottom encodings, pixel height will be doubled.
	MultiviewFlagsMixedMono      MultiviewFlags = C.GST_VIDEO_MULTIVIEW_FLAGS_MIXED_MONO       // (32768) – The video stream contains both mono and multiview portions, signalled on each buffer by the absence or presence of the GST_VIDEO_BUFFER_FLAG_MULTIPLE_VIEW buffer flag.
)

// MultiviewFramePacking represents the subset of MultiviewMode values that can be applied to any video frame
// without needing extra metadata. It can be used by elements that provide a property to override the multiview
// interpretation of a video stream when the video doesn't contain any markers.
//
// This enum is used (for example) on playbin, to re-interpret a played video stream as a stereoscopic video.
// The individual enum values are equivalent to and have the same value as the matching MultiviewMode.
type MultiviewFramePacking int

// Type castings
const (
	MultiviewFramePackingNone               MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_NONE                  // (-1) – A special value indicating no frame packing info.
	MultiviewFramePackingMono               MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_MONO                  // (0) – All frames are monoscopic.
	MultiviewFramePackingLeft               MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_LEFT                  // (1) – All frames represent a left-eye view.
	MultiviewFramePackingRight              MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_RIGHT                 // (2) – All frames represent a right-eye view.
	MultiviewFramePackingSideBySide         MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_SIDE_BY_SIDE          // (3) – Left and right eye views are provided in the left and right half of the frame respectively.
	MultiviewFramePackingSideBySideQuincunx MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_SIDE_BY_SIDE_QUINCUNX // (4) – Left and right eye views are provided in the left and right half of the frame, but have been sampled using quincunx method, with half-pixel offset between the 2 views.
	MultiviewFramePackingColumnInterleaved  MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_COLUMN_INTERLEAVED    // (5) – Alternating vertical columns of pixels represent the left and right eye view respectively.
	MultiviewFramePackingRowInterleaved     MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_ROW_INTERLEAVED       // (6) – Alternating horizontal rows of pixels represent the left and right eye view respectively.
	MultiviewFramePackingTopBottom          MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_TOP_BOTTOM            // (7) – The top half of the frame contains the left eye, and the bottom half the right eye.
	MultiviewFramePackingCheckerboard       MultiviewFramePacking = C.GST_VIDEO_MULTIVIEW_FRAME_PACKING_CHECKERBOARD          // (8) – Pixels are arranged with alternating pixels representing left and right eye views in a checkerboard fashion.
)

// MultiviewMode represents all possible stereoscopic 3D and multiview representations. In conjunction with
// MultiviewFlags, describes how multiview content is being transported in the stream.
type MultiviewMode int

// Type castings
const (
	MultiviewModeNone                  MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_NONE                     // (-1) – A special value indicating no multiview information. Used in GstVideoInfo and other places to indicate that no specific multiview handling has been requested or provided. This value is never carried on caps.
	MultiviewModeMono                  MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_MONO                     // (0) – All frames are monoscopic.
	MultiviewModeLeft                  MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_LEFT                     // (1) – All frames represent a left-eye view.
	MultiviewModeRight                 MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_RIGHT                    // (2) – All frames represent a right-eye view.
	MultiviewModeSideBySide            MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_SIDE_BY_SIDE             // (3) – Left and right eye views are provided in the left and right half of the frame respectively.
	MultiviewModeSideBySideQuincunx    MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_SIDE_BY_SIDE_QUINCUNX    // (4) – Left and right eye views are provided in the left and right half of the frame, but have been sampled using quincunx method, with half-pixel offset between the 2 views.
	MultiviewModeColumnInterleaved     MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_COLUMN_INTERLEAVED       // (5) – Alternating vertical columns of pixels represent the left and right eye view respectively.
	MultiviewModeRowInterleaved        MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_ROW_INTERLEAVED          // (6) – Alternating horizontal rows of pixels represent the left and right eye view respectively.
	MultiviewModeTopBottom             MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_TOP_BOTTOM               // (7) – The top half of the frame contains the left eye, and the bottom half the right eye.
	MultiviewModeCheckerboard          MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_CHECKERBOARD             // (8) – Pixels are arranged with alternating pixels representing left and right eye views in a checkerboard fashion.
	MultiviewModeFrameByFrame          MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_FRAME_BY_FRAME           // (32) – Left and right eye views are provided in separate frames alternately.
	MultiviewModeMultiviewFrameByFrame MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_MULTIVIEW_FRAME_BY_FRAME // (33) – Multiple independent views are provided in separate frames in sequence. This method only applies to raw video buffers at the moment. Specific view identification is via the GstVideoMultiviewMeta and GstVideoMeta(s) on raw video buffers.
	MultiviewModeSeparated             MultiviewMode = C.GST_VIDEO_MULTIVIEW_MODE_SEPARATED                // (34) – Multiple views are provided as separate GstMemory framebuffers attached to each GstBuffer, described by the GstVideoMultiviewMeta and GstVideoMeta(s)
)

// Info describes image properties. This information can be filled in from GstCaps with
// InfoFromCaps. The information is also used to store the specific video info when mapping
// a video frame with FrameMap.
type Info struct {
	ptr *C.GstVideoInfo
}

func wrapInfo(vinfo *C.GstVideoInfo) *Info {
	info := &Info{vinfo}
	runtime.SetFinalizer(info, (*Info).Free)
	return info
}

// Free will free this video info
func (i *Info) Free() {
	C.gst_video_info_free(i.instance())
}

// instance returns the underlying GstVideoInfo instance.
func (i *Info) instance() *C.GstVideoInfo { return i.ptr }

// NewInfo returns a new Info instance. You can populate it by chaining builders
// to this constructor.
func NewInfo() *Info {
	return wrapInfo(C.gst_video_info_new())
}

// FromCaps parses the caps and updates this info.
func (i *Info) FromCaps(caps *gst.Caps) *Info {
	C.gst_video_info_from_caps(i.instance(), fromCoreCaps(caps))
	return i
}

// Convert converts among various gst.Format types. This function handles gst.FormatBytes, gst.FormatTime,
// and gst.FormatDefault. For raw video, gst.FormatDefault corresponds to video frames. This function can
// be used to handle pad queries of the type gst.QueryTypeConvert.
func (i *Info) Convert(srcFormat, destFormat gst.Format, srcValue int64) (out int64, ok bool) {
	var gout C.gint64
	gok := C.gst_video_info_convert(i.instance(), C.GstFormat(srcFormat), C.gint64(srcValue), C.GstFormat(destFormat), &gout)
	return int64(gout), gobool(gok)
}

// IsEqual compares two GstVideoInfo and returns whether they are equal or not.
func (i *Info) IsEqual(info *Info) bool {
	return gobool(C.gst_video_info_is_equal(i.instance(), info.instance()))
}

// ChromaSite returns the ChromaSite for this info.
func (i *Info) ChromaSite() ChromaSite {
	return ChromaSite(C.infoChromaSite(i.instance()))
}

// Colorimetry returns the colorimetry for this info.
func (i *Info) Colorimetry() *Colorimetry {
	return colorimetryFromInstance(C.infoColorimetry(i.instance()))
}

// FieldHeight returns the field height for this info.
func (i *Info) FieldHeight() int {
	return int(C.infoFieldHeight(i.instance()))
}

// FieldOrder returns the field order for this info.
func (i *Info) FieldOrder() FieldOrder {
	return FieldOrder(C.infoFieldOrder(i.instance()))
}

// FieldRateN returns the rate numerator depending on the interlace mode.
func (i *Info) FieldRateN() int {
	return int(C.infoFieldRateN(i.instance()))
}

// Flags returns the flags on this info.
func (i *Info) Flags() Flags {
	return Flags(C.infoFlags(i.instance()))
}

// FlagIsSet returns true if the given flag(s) are set on the info.
func (i *Info) FlagIsSet(f Flags) bool {
	return gobool(C.infoFlagIsSet(i.instance(), C.GstVideoFlags(f)))
}

// FlagSet sets the given flag(s) on the info. The underlying info is returned
// for chaining builders.
func (i *Info) FlagSet(f Flags) *Info {
	C.infoFlagSet(i.instance(), C.GstVideoFlags(f))
	return i
}

// FlagUnset unsets the given flag(s) on the info. The underlying info is returned
// for chaining builders.
func (i *Info) FlagUnset(f Flags) *Info {
	C.infoFlagUnset(i.instance(), C.GstVideoFlags(f))
	return i
}

// Format returns the format for the info. You can call Info() on the return value
// to inspect the properties further.
func (i *Info) Format() Format {
	return Format(C.infoFormat(i.instance()))
}

// FPS returns the frames-per-second value for the info.
func (i *Info) FPS() *gst.FractionValue {
	return gst.Fraction(
		int(C.infoFPSn(i.instance())),
		int(C.infoFPSd(i.instance())),
	)
}

// HasAlpha returns true if the alpha flag is set on the format info.
func (i *Info) HasAlpha() bool {
	return gobool(C.infoHasAlpha(i.instance()))
}

// Height returns the height of the video.
func (i *Info) Height() int { return int(C.infoHeight(i.instance())) }

// InterlaceMode returns the interlace mode of this Info.
func (i *Info) InterlaceMode() InterlaceMode {
	return InterlaceMode(C.infoInterlaceMode(i.instance()))
}

// IsInterlaced returns true if the interlace mode is not Progressive.
func (i *Info) IsInterlaced() bool {
	return gobool(C.infoIsInterlaced(i.instance()))
}

// IsGray returns if the format is grayscale.
func (i *Info) IsGray() bool { return gobool(C.infoIsGray(i.instance())) }

// IsRGB returns if the format is RGB.
func (i *Info) IsRGB() bool { return gobool(C.infoIsRGB(i.instance())) }

// IsYUV returns if the format is YUV.
func (i *Info) IsYUV() bool { return gobool(C.infoIsYUV(i.instance())) }

// MultiviewFlags returns the MultiviewFlags on the info.
func (i *Info) MultiviewFlags() MultiviewFlags {
	return MultiviewFlags(C.infoMultiviewFlags(i.instance()))
}

// MultiviewMode returns the MultiviewMode on thee info.
func (i *Info) MultiviewMode() MultiviewMode {
	return MultiviewMode(C.infoMultiviewMode(i.instance()))
}

// Name returns a human readable name forr the info.
func (i *Info) Name() string {
	return C.GoString(C.infoName(i.instance()))
}

// NumComponents returns the number of components in the info.
func (i *Info) NumComponents() uint {
	return uint(C.infoNComponents(i.instance()))
}

// NumPlanes returns the number of planes in the info.
func (i *Info) NumPlanes() uint {
	return uint(C.infoNPlanes(i.instance()))
}

// PAR returns the pixel-aspect-ration value for the info.
func (i *Info) PAR() *gst.FractionValue {
	return gst.Fraction(
		int(C.infoPARn(i.instance())),
		int(C.infoPARd(i.instance())),
	)
}

// Size returns the size of the info.
func (i *Info) Size() int64 {
	return int64(C.infoSize(i.instance()))
}

// Views returns the number of views.
func (i *Info) Views() int {
	return int(C.infoViews(i.instance()))
}

// Width returns the width of the video.
func (i *Info) Width() int { return int(C.infoWidth(i.instance())) }

// WithAlign adjusts the offset and stride fields in info so that the padding and stride alignment in
// align is respected.
//
// Extra padding will be added to the right side when stride alignment padding is required and align
// will be updated with the new padding values.
func (i *Info) WithAlign(align *Alignment) *Info {
	C.gst_video_info_align(i.instance(), align.instance())
	return i
}

// WithFormat sets the format on this info.
//
// Note: This initializes info first, no values are preserved. This function does not set the offsets
// correctly for interlaced vertically subsampled formats. If the format is invalid (e.g. because the
// size of a frame can't be represented as a 32 bit integer), nothing will happen. This is is for
// convenience in chaining, but may be changed in the future.
func (i *Info) WithFormat(format Format, width, height uint) *Info {
	C.gst_video_info_set_format(i.instance(), C.GstVideoFormat(format), C.guint(width), C.guint(height))
	return i
}

// WithInterlacedFormat is the same as WithFormat but also allows to set the interlaced mode.
func (i *Info) WithInterlacedFormat(format Format, interlaceMode InterlaceMode, width, height uint) *Info {
	C.gst_video_info_set_interlaced_format(
		i.instance(),
		C.GstVideoFormat(format),
		C.GstVideoInterlaceMode(interlaceMode),
		C.guint(width), C.guint(height),
	)
	return i
}

// WithFPS sets the FPS on this info.
func (i *Info) WithFPS(f *gst.FractionValue) *Info {
	i.instance().fps_d = C.gint(f.Denom())
	i.instance().fps_n = C.gint(f.Num())
	return i
}

// WithPAR sets the FPS on this info.
func (i *Info) WithPAR(f *gst.FractionValue) *Info {
	i.instance().par_d = C.gint(f.Denom())
	i.instance().par_n = C.gint(f.Num())
	return i
}

// ToCaps returns the caps representation of this video info.
func (i *Info) ToCaps() *gst.Caps {
	caps := C.gst_video_info_to_caps(i.instance())
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}
