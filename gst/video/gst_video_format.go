package video

/*
#include <gst/video/video.h>

const gchar *         formatInfoName        (GstVideoFormatInfo * info)   { return GST_VIDEO_FORMAT_INFO_NAME(info); }

guint                 formatInfoBits        (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_BITS(info); }
guint                 formatInfoDepth       (GstVideoFormatInfo * info, guint c)  { return GST_VIDEO_FORMAT_INFO_DEPTH(info, c); }
GstVideoFormatFlags   formatInfoFlags       (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_FLAGS(info); }
GstVideoFormat        formatInfoFormat      (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_FORMAT(info); }
gboolean              formatInfoHasAlpha    (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_HAS_ALPHA(info); }
gboolean              formatInfoHasPalette  (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_HAS_PALETTE(info); }
guint                 formatInfoHSub        (GstVideoFormatInfo * info, guint c)  { return GST_VIDEO_FORMAT_INFO_H_SUB(info, c); }
gboolean              formatInfoIsComplex   (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_IS_COMPLEX(info); }
gboolean              formatInfoIsGray      (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_IS_GRAY(info); }
gboolean              formatInfoIsLE        (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_IS_LE(info); }
gboolean              formatInfoIsRGB       (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_IS_RGB(info); }
gboolean              formatInfoIsTiled     (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_IS_TILED(info); }
gboolean              formatInfoIsYUV       (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_IS_YUV(info); }
guint                 formatInfoNComponent  (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_N_COMPONENTS(info); }
guint                 formatInfoNPlanes     (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_N_PLANES(info); }
guint                 formatInfoPlane       (GstVideoFormatInfo * info, guint c)  { return GST_VIDEO_FORMAT_INFO_PLANE(info, c); }
guint                 formatInfoPOffset     (GstVideoFormatInfo * info, guint c)  { return GST_VIDEO_FORMAT_INFO_POFFSET(info, c); }
guint                 formatInfoPStride     (GstVideoFormatInfo * info, guint c)  { return GST_VIDEO_FORMAT_INFO_PSTRIDE(info, c); }
guint                 formatInfoTileHS      (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_TILE_HS(info); }
GstVideoTileMode      formatInfoTileMode    (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_TILE_MODE(info); }
guint                 formatInfoTileWS      (GstVideoFormatInfo * info)           { return GST_VIDEO_FORMAT_INFO_TILE_WS(info); }
guint                 formatInfoWSub        (GstVideoFormatInfo * info, guint c)  { return GST_VIDEO_FORMAT_INFO_W_SUB(info, c); }

*/
import "C"
import (
	"image/color"
	"math"
	"runtime"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

func init() {
	glib.RegisterGValueMarshalers([]glib.TypeMarshaler{
		{
			T: glib.Type(C.gst_video_format_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_enum(uintptrToGVal(p))
				return Format(c), nil
			},
		},
	})
}

// Format is an enum value describing the most common video formats.
type Format int

// Type castings
const (
	FormatUnknown    Format = C.GST_VIDEO_FORMAT_UNKNOWN     // (0) – Unknown or unset video format id
	FormatEncoded    Format = C.GST_VIDEO_FORMAT_ENCODED     // (1) – Encoded video format. Only ever use that in caps for special video formats in combination with non-system memory GstCapsFeatures where it does not make sense to specify a real video format.
	FormatI420       Format = C.GST_VIDEO_FORMAT_I420        // (2) – planar 4:2:0 YUV
	FormatYV12       Format = C.GST_VIDEO_FORMAT_YV12        // (3) – planar 4:2:0 YVU (like I420 but UV planes swapped)
	FormatYUY2       Format = C.GST_VIDEO_FORMAT_YUY2        // (4) – packed 4:2:2 YUV (Y0-U0-Y1-V0 Y2-U2-Y3-V2 Y4 ...)
	FormatUYVY       Format = C.GST_VIDEO_FORMAT_UYVY        // (5) – packed 4:2:2 YUV (U0-Y0-V0-Y1 U2-Y2-V2-Y3 U4 ...)
	FormatAYUV       Format = C.GST_VIDEO_FORMAT_AYUV        // (6) – packed 4:4:4 YUV with alpha channel (A0-Y0-U0-V0 ...)
	FormatRGBx       Format = C.GST_VIDEO_FORMAT_RGBx        // (7) – sparse rgb packed into 32 bit, space last
	FormatBGRx       Format = C.GST_VIDEO_FORMAT_BGRx        // (8) – sparse reverse rgb packed into 32 bit, space last
	FormatxRGB       Format = C.GST_VIDEO_FORMAT_xRGB        // (9) – sparse rgb packed into 32 bit, space first
	FormatxBGR       Format = C.GST_VIDEO_FORMAT_xBGR        // (10) – sparse reverse rgb packed into 32 bit, space first
	FormatRGBA       Format = C.GST_VIDEO_FORMAT_RGBA        // (11) – rgb with alpha channel last
	FormatBGRA       Format = C.GST_VIDEO_FORMAT_BGRA        // (12) – reverse rgb with alpha channel last
	FormatARGB       Format = C.GST_VIDEO_FORMAT_ARGB        // (13) – rgb with alpha channel first
	FormatABGR       Format = C.GST_VIDEO_FORMAT_ABGR        // (14) – reverse rgb with alpha channel first
	FormatRGB        Format = C.GST_VIDEO_FORMAT_RGB         // (15) – RGB packed into 24 bits without padding (R-G-B-R-G-B)
	FormatBGR        Format = C.GST_VIDEO_FORMAT_BGR         // (16) – reverse RGB packed into 24 bits without padding (B-G-R-B-G-R)
	FormatY41B       Format = C.GST_VIDEO_FORMAT_Y41B        // (17) – planar 4:1:1 YUV
	FormatY42B       Format = C.GST_VIDEO_FORMAT_Y42B        // (18) – planar 4:2:2 YUV
	FormatYVYU       Format = C.GST_VIDEO_FORMAT_YVYU        // (19) – packed 4:2:2 YUV (Y0-V0-Y1-U0 Y2-V2-Y3-U2 Y4 ...)
	FormatY444       Format = C.GST_VIDEO_FORMAT_Y444        // (20) – planar 4:4:4 YUV
	Formatv210       Format = C.GST_VIDEO_FORMAT_v210        // (21) – packed 4:2:2 10-bit YUV, complex format
	Formatv216       Format = C.GST_VIDEO_FORMAT_v216        // (22) – packed 4:2:2 16-bit YUV, Y0-U0-Y1-V1 order
	FormatNV12       Format = C.GST_VIDEO_FORMAT_NV12        // (23) – planar 4:2:0 YUV with interleaved UV plane
	FormatNV21       Format = C.GST_VIDEO_FORMAT_NV21        // (24) – planar 4:2:0 YUV with interleaved VU plane
	FormatGray8      Format = C.GST_VIDEO_FORMAT_GRAY8       // (25) – 8-bit grayscale
	FormatGray16BE   Format = C.GST_VIDEO_FORMAT_GRAY16_BE   // (26) – 16-bit grayscale, most significant byte first
	FormatGray16LE   Format = C.GST_VIDEO_FORMAT_GRAY16_LE   // (27) – 16-bit grayscale, least significant byte first
	Formatv308       Format = C.GST_VIDEO_FORMAT_v308        // (28) – packed 4:4:4 YUV (Y-U-V ...)
	FormatRGB16      Format = C.GST_VIDEO_FORMAT_RGB16       // (29) – rgb 5-6-5 bits per component
	FormatBGR16      Format = C.GST_VIDEO_FORMAT_BGR16       // (30) – reverse rgb 5-6-5 bits per component
	FormatRGB15      Format = C.GST_VIDEO_FORMAT_RGB15       // (31) – rgb 5-5-5 bits per component
	FormatBGR15      Format = C.GST_VIDEO_FORMAT_BGR15       // (32) – reverse rgb 5-5-5 bits per component
	FormatUYVP       Format = C.GST_VIDEO_FORMAT_UYVP        // (33) – packed 10-bit 4:2:2 YUV (U0-Y0-V0-Y1 U2-Y2-V2-Y3 U4 ...)
	FormatA420       Format = C.GST_VIDEO_FORMAT_A420        // (34) – planar 4:4:2:0 AYUV
	FormatRGB8P      Format = C.GST_VIDEO_FORMAT_RGB8P       // (35) – 8-bit paletted RGB
	FormatYUV9       Format = C.GST_VIDEO_FORMAT_YUV9        // (36) – planar 4:1:0 YUV
	FormatYVU9       Format = C.GST_VIDEO_FORMAT_YVU9        // (37) – planar 4:1:0 YUV (like YUV9 but UV planes swapped)
	FormatIYU1       Format = C.GST_VIDEO_FORMAT_IYU1        // (38) – packed 4:1:1 YUV (Cb-Y0-Y1-Cr-Y2-Y3 ...)
	FormatARGB64     Format = C.GST_VIDEO_FORMAT_ARGB64      // (39) – rgb with alpha channel first, 16 bits per channel
	FormatAYUV64     Format = C.GST_VIDEO_FORMAT_AYUV64      // (40) – packed 4:4:4 YUV with alpha channel, 16 bits per channel (A0-Y0-U0-V0 ...)
	Formatr210       Format = C.GST_VIDEO_FORMAT_r210        // (41) – packed 4:4:4 RGB, 10 bits per channel
	FormatI42010BE   Format = C.GST_VIDEO_FORMAT_I420_10BE   // (42) – planar 4:2:0 YUV, 10 bits per channel
	FormatI42010LE   Format = C.GST_VIDEO_FORMAT_I420_10LE   // (43) – planar 4:2:0 YUV, 10 bits per channel
	FormatI42210BE   Format = C.GST_VIDEO_FORMAT_I422_10BE   // (44) – planar 4:2:2 YUV, 10 bits per channel
	FormatI42210LE   Format = C.GST_VIDEO_FORMAT_I422_10LE   // (45) – planar 4:2:2 YUV, 10 bits per channel
	FormatY44410BE   Format = C.GST_VIDEO_FORMAT_Y444_10BE   // (46) – planar 4:4:4 YUV, 10 bits per channel (Since: 1.2)
	FormatY44410LE   Format = C.GST_VIDEO_FORMAT_Y444_10LE   // (47) – planar 4:4:4 YUV, 10 bits per channel (Since: 1.2)
	FormatGBR        Format = C.GST_VIDEO_FORMAT_GBR         // (48) – planar 4:4:4 RGB, 8 bits per channel (Since: 1.2)
	FormatGBR10BE    Format = C.GST_VIDEO_FORMAT_GBR_10BE    // (49) – planar 4:4:4 RGB, 10 bits per channel (Since: 1.2)
	FormatGBR10LE    Format = C.GST_VIDEO_FORMAT_GBR_10LE    // (50) – planar 4:4:4 RGB, 10 bits per channel (Since: 1.2)
	FormatNV16       Format = C.GST_VIDEO_FORMAT_NV16        // (51) – planar 4:2:2 YUV with interleaved UV plane (Since: 1.2)
	FormatNV24       Format = C.GST_VIDEO_FORMAT_NV24        // (52) – planar 4:4:4 YUV with interleaved UV plane (Since: 1.2)
	FormatNV1264Z32  Format = C.GST_VIDEO_FORMAT_NV12_64Z32  // (53) – NV12 with 64x32 tiling in zigzag pattern (Since: 1.4)
	FormatA42010BE   Format = C.GST_VIDEO_FORMAT_A420_10BE   // (54) – planar 4:4:2:0 YUV, 10 bits per channel (Since: 1.6)
	FormatA42010LE   Format = C.GST_VIDEO_FORMAT_A420_10LE   // (55) – planar 4:4:2:0 YUV, 10 bits per channel (Since: 1.6)
	FormatA42210BE   Format = C.GST_VIDEO_FORMAT_A422_10BE   // (56) – planar 4:4:2:2 YUV, 10 bits per channel (Since: 1.6)
	FormatA42210LE   Format = C.GST_VIDEO_FORMAT_A422_10LE   // (57) – planar 4:4:2:2 YUV, 10 bits per channel (Since: 1.6)
	FormatA44410BE   Format = C.GST_VIDEO_FORMAT_A444_10BE   // (58) – planar 4:4:4:4 YUV, 10 bits per channel (Since: 1.6)
	FormatA44410LE   Format = C.GST_VIDEO_FORMAT_A444_10LE   // (59) – planar 4:4:4:4 YUV, 10 bits per channel (Since: 1.6)
	FormatNV61       Format = C.GST_VIDEO_FORMAT_NV61        // (60) – planar 4:2:2 YUV with interleaved VU plane (Since: 1.6)
	FormatP01010BE   Format = C.GST_VIDEO_FORMAT_P010_10BE   // (61) – planar 4:2:0 YUV with interleaved UV plane, 10 bits per channel (Since: 1.10)
	FormatP01010LE   Format = C.GST_VIDEO_FORMAT_P010_10LE   // (62) – planar 4:2:0 YUV with interleaved UV plane, 10 bits per channel (Since: 1.10)
	FormatIYU2       Format = C.GST_VIDEO_FORMAT_IYU2        // (63) – packed 4:4:4 YUV (U-Y-V ...) (Since: 1.10)
	FormatVYUY       Format = C.GST_VIDEO_FORMAT_VYUY        // (64) – packed 4:2:2 YUV (V0-Y0-U0-Y1 V2-Y2-U2-Y3 V4 ...)
	FormatGBRA       Format = C.GST_VIDEO_FORMAT_GBRA        // (65) – planar 4:4:4:4 ARGB, 8 bits per channel (Since: 1.12)
	FormatGBRA10BE   Format = C.GST_VIDEO_FORMAT_GBRA_10BE   // (66) – planar 4:4:4:4 ARGB, 10 bits per channel (Since: 1.12)
	FormatGBRA10LE   Format = C.GST_VIDEO_FORMAT_GBRA_10LE   // (67) – planar 4:4:4:4 ARGB, 10 bits per channel (Since: 1.12)
	FormatGBR12BE    Format = C.GST_VIDEO_FORMAT_GBR_12BE    // (68) – planar 4:4:4 RGB, 12 bits per channel (Since: 1.12)
	FormatGBR12LE    Format = C.GST_VIDEO_FORMAT_GBR_12LE    // (69) – planar 4:4:4 RGB, 12 bits per channel (Since: 1.12)
	FormatGBRA12BE   Format = C.GST_VIDEO_FORMAT_GBRA_12BE   // (70) – planar 4:4:4:4 ARGB, 12 bits per channel (Since: 1.12)
	FormatGBRA12LE   Format = C.GST_VIDEO_FORMAT_GBRA_12LE   // (71) – planar 4:4:4:4 ARGB, 12 bits per channel (Since: 1.12)
	FormatI42012BE   Format = C.GST_VIDEO_FORMAT_I420_12BE   // (72) – planar 4:2:0 YUV, 12 bits per channel (Since: 1.12)
	FormatI42012LE   Format = C.GST_VIDEO_FORMAT_I420_12LE   // (73) – planar 4:2:0 YUV, 12 bits per channel (Since: 1.12)
	FormatI42212BE   Format = C.GST_VIDEO_FORMAT_I422_12BE   // (74) – planar 4:2:2 YUV, 12 bits per channel (Since: 1.12)
	FormatI42212LE   Format = C.GST_VIDEO_FORMAT_I422_12LE   // (75) – planar 4:2:2 YUV, 12 bits per channel (Since: 1.12)
	FormatY44412BE   Format = C.GST_VIDEO_FORMAT_Y444_12BE   // (76) – planar 4:4:4 YUV, 12 bits per channel (Since: 1.12)
	FormatY44412LE   Format = C.GST_VIDEO_FORMAT_Y444_12LE   // (77) – planar 4:4:4 YUV, 12 bits per channel (Since: 1.12)
	FormatGray10LE32 Format = C.GST_VIDEO_FORMAT_GRAY10_LE32 // (78) – 10-bit grayscale, packed into 32bit words (2 bits padding) (Since: 1.14)
	FormatNV1210LE32 Format = C.GST_VIDEO_FORMAT_NV12_10LE32 // (79) – 10-bit variant of GST_VIDEO_FORMAT_NV12, packed into 32bit words (MSB 2 bits padding) (Since: 1.14)
	FormatNV1610LE32 Format = C.GST_VIDEO_FORMAT_NV16_10LE32 // (80) – 10-bit variant of GST_VIDEO_FORMAT_NV16, packed into 32bit words (MSB 2 bits padding) (Since: 1.14)
	FormatNV1210LE40 Format = C.GST_VIDEO_FORMAT_NV12_10LE40 // (81) – Fully packed variant of NV12_10LE32 (Since: 1.16)
	FormatY210       Format = C.GST_VIDEO_FORMAT_Y210        // (82) – packed 4:2:2 YUV, 10 bits per channel (Since: 1.16)
	FormatY410       Format = C.GST_VIDEO_FORMAT_Y410        // (83) – packed 4:4:4 YUV, 10 bits per channel(A-V-Y-U...) (Since: 1.16)
	FormatVUYA       Format = C.GST_VIDEO_FORMAT_VUYA        // (84) – packed 4:4:4 YUV with alpha channel (V0-U0-Y0-A0...) (Since: 1.16)
	FormatBGR10A2LE  Format = C.GST_VIDEO_FORMAT_BGR10A2_LE  // (85) – packed 4:4:4 RGB with alpha channel(B-G-R-A), 10 bits for R/G/B channel and MSB 2 bits for alpha channel (Since: 1.16)
	FormatRGB10A2LE  Format = C.GST_VIDEO_FORMAT_RGB10A2_LE  // (86) – packed 4:4:4 RGB with alpha channel(R-G-B-A), 10 bits for R/G/B channel and MSB 2 bits for alpha channel (Since: 1.18)
	FormatY44416BE   Format = C.GST_VIDEO_FORMAT_Y444_16BE   // (87) – planar 4:4:4 YUV, 16 bits per channel (Since: 1.18)
	FormatY44416LE   Format = C.GST_VIDEO_FORMAT_Y444_16LE   // (88) – planar 4:4:4 YUV, 16 bits per channel (Since: 1.18)
	FormatP016BE     Format = C.GST_VIDEO_FORMAT_P016_BE     // (89) – planar 4:2:0 YUV with interleaved UV plane, 16 bits per channel (Since: 1.18)
	FormatP016LE     Format = C.GST_VIDEO_FORMAT_P016_LE     // (90) – planar 4:2:0 YUV with interleaved UV plane, 16 bits per channel (Since: 1.18)
	FormatP012BE     Format = C.GST_VIDEO_FORMAT_P012_BE     // (91) – planar 4:2:0 YUV with interleaved UV plane, 12 bits per channel (Since: 1.18)
	FormatP012LE     Format = C.GST_VIDEO_FORMAT_P012_LE     // (92) – planar 4:2:0 YUV with interleaved UV plane, 12 bits per channel (Since: 1.18)
	FormatY212BE     Format = C.GST_VIDEO_FORMAT_Y212_BE     // (93) – packed 4:2:2 YUV, 12 bits per channel (Y-U-Y-V) (Since: 1.18)
	FormatY212LE     Format = C.GST_VIDEO_FORMAT_Y212_LE     // (94) – packed 4:2:2 YUV, 12 bits per channel (Y-U-Y-V) (Since: 1.18)
	FormatY412BE     Format = C.GST_VIDEO_FORMAT_Y412_BE     // (95) – packed 4:4:4:4 YUV, 12 bits per channel(U-Y-V-A...) (Since: 1.18)
	FormatY412LE     Format = C.GST_VIDEO_FORMAT_Y412_LE     // (96) – packed 4:4:4:4 YUV, 12 bits per channel(U-Y-V-A...) (Since: 1.18)
	FormatNV124L4    Format = C.GST_VIDEO_FORMAT_NV12_4L4    // (97) – NV12 with 4x4 tiles in linear order.
	FormatNV1232L32  Format = C.GST_VIDEO_FORMAT_NV12_32L32  // (98) – NV12 with 32x32 tiles in linear order.
)

// AllFormats is a convenience function for retrieving all formats for inspection purposes.
// This is not really intended for use in an application, and moreso for debugging.
func AllFormats() []Format {
	return []Format{
		FormatI420,
		FormatYV12,
		FormatYUY2,
		FormatUYVY,
		FormatAYUV,
		FormatRGBx,
		FormatBGRx,
		FormatxRGB,
		FormatxBGR,
		FormatRGBA,
		FormatBGRA,
		FormatARGB,
		FormatABGR,
		FormatRGB,
		FormatBGR,
		FormatY41B,
		FormatY42B,
		FormatYVYU,
		FormatY444,
		Formatv210,
		Formatv216,
		FormatNV12,
		FormatNV21,
		FormatGray8,
		FormatGray16BE,
		FormatGray16LE,
		Formatv308,
		FormatRGB16,
		FormatBGR16,
		FormatRGB15,
		FormatBGR15,
		FormatUYVP,
		FormatA420,
		FormatRGB8P,
		FormatYUV9,
		FormatYVU9,
		FormatIYU1,
		FormatARGB64,
		FormatAYUV64,
		Formatr210,
		FormatI42010BE,
		FormatI42010LE,
		FormatI42210BE,
		FormatI42210LE,
		FormatY44410BE,
		FormatY44410LE,
		FormatGBR,
		FormatGBR10BE,
		FormatGBR10LE,
		FormatNV16,
		FormatNV24,
		FormatNV1264Z32,
		FormatA42010BE,
		FormatA42010LE,
		FormatA42210BE,
		FormatA42210LE,
		FormatA44410BE,
		FormatA44410LE,
		FormatNV61,
		FormatP01010BE,
		FormatP01010LE,
		FormatIYU2,
		FormatVYUY,
		FormatGBRA,
		FormatGBRA10BE,
		FormatGBRA10LE,
		FormatGBR12BE,
		FormatGBR12LE,
		FormatGBRA12BE,
		FormatGBRA12LE,
		FormatI42012BE,
		FormatI42012LE,
		FormatI42212BE,
		FormatI42212LE,
		FormatY44412BE,
		FormatY44412LE,
		FormatGray10LE32,
		FormatNV1210LE32,
		FormatNV1610LE32,
		FormatNV1210LE40,
		FormatY210,
		FormatY410,
		FormatVUYA,
		FormatBGR10A2LE,
		FormatRGB10A2LE,
		FormatY44416BE,
		FormatY44416LE,
		FormatP016BE,
		FormatP016LE,
		FormatP012BE,
		FormatP012LE,
		FormatY212BE,
		FormatY212LE,
		FormatY412BE,
		FormatY412LE,
		FormatNV124L4,
		FormatNV1232L32,
	}
}

// RawFormats returns a slice of all the raw video formats supported by GStreamer.
func RawFormats() []Format {
	var size C.guint
	formats := C.gst_video_formats_raw(&size)
	out := make([]Format, uint(size))
	for i, f := range (*[(math.MaxInt32 - 1) / unsafe.Sizeof(C.GST_VIDEO_FORMAT_UNKNOWN)]C.GstVideoFormat)(unsafe.Pointer(formats))[:size:size] {
		out[i] = Format(f)
	}
	return out
}

// MakeRawCaps returns a generic raw video caps for formats defined in formats. If formats is empty or nil, returns a caps for
// all the supported raw video formats, see RawFormats.
func MakeRawCaps(formats []Format) *gst.Caps {
	var caps *C.GstCaps
	if len(formats) == 0 {
		caps = C.gst_video_make_raw_caps(nil, C.guint(0))
	} else {
		caps = C.gst_video_make_raw_caps(
			(*C.GstVideoFormat)(unsafe.Pointer(&formats[0])),
			C.guint(len(formats)),
		)
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// MakeRawCapsWithFeatures returns a generic raw video caps for formats defined in formats with features. If formats is
// empty or nil, returns a caps for all the supported video formats, see RawFormats.
func MakeRawCapsWithFeatures(formats []Format, features *gst.CapsFeatures) *gst.Caps {
	var caps *C.GstCaps
	if len(formats) == 0 {
		caps = C.gst_video_make_raw_caps_with_features(nil, C.guint(0), fromCoreCapsFeatures(features))
	} else {
		caps = C.gst_video_make_raw_caps_with_features(
			(*C.GstVideoFormat)(unsafe.Pointer(&formats[0])),
			C.guint(len(formats)),
			fromCoreCapsFeatures(features),
		)
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// TypeFormat is the GType for a GstVideoFormat.
var TypeFormat = glib.Type(C.gst_video_format_get_type())

// ToGValue implements a glib.ValueTransformer
func (f Format) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeFormat)
	if err != nil {
		return nil, err
	}
	val.SetEnum(int(f))
	return val, nil
}

// Info returns the FormatInfo for this video format.
func (f Format) Info() *FormatInfo {
	finfo := C.gst_video_format_get_info(C.GstVideoFormat(f))
	info := &FormatInfo{ptr: finfo}
	runtime.SetFinalizer(info, func(_ *FormatInfo) { C.g_free((C.gpointer)(unsafe.Pointer(finfo))) })
	return info
}

// Palette returns the color palette for this format, or nil if the format does not have one.
// At time of writing, RGB8P appears to be the only format with it's own palette.
func (f Format) Palette() color.Palette {
	var size C.gsize
	ptr := C.gst_video_format_get_palette(C.GstVideoFormat(f), &size)
	if ptr == nil {
		return nil
	}
	paletteBytes := make([]uint8, int64(size))
	for i, t := range (*[(math.MaxInt32 - 1) / unsafe.Sizeof(uint8(0))]uint8)(ptr)[:int(size):int(size)] {
		paletteBytes[i] = t
	}
	return bytesToColorPalette(paletteBytes)
}

func bytesToColorPalette(in []uint8) color.Palette {
	palette := make([]color.Color, len(in)/4)
	for i := 0; i < len(in); i += 4 {
		palette[i/4] = color.RGBA{in[i], in[i+1], in[i+2], in[i+3]}
	}
	return color.Palette(palette)
}

// String implements a stringer on a Format.
func (f Format) String() string {
	return C.GoString(C.gst_video_format_to_string(C.GstVideoFormat(f)))
}

// FOURCC converts this format value into the corresponding FOURCC. Only a few YUV formats have corresponding
// FOURCC values. If format has no corresponding FOURCC value, 0 is returned.
func (f Format) FOURCC() uint32 {
	return uint32(C.gst_video_format_to_fourcc(C.GstVideoFormat(f)))
}

// FormatFlags are different video flags that a format info can have.
type FormatFlags int

// Type castings
const (
	FormatFlagYUV     FormatFlags = C.GST_VIDEO_FORMAT_FLAG_YUV     // (1) – The video format is YUV, components are numbered 0=Y, 1=U, 2=V.
	FormatFlagRGB     FormatFlags = C.GST_VIDEO_FORMAT_FLAG_RGB     // (2) – The video format is RGB, components are numbered 0=R, 1=G, 2=B.
	FormatFlagGray    FormatFlags = C.GST_VIDEO_FORMAT_FLAG_GRAY    // (4) – The video is gray, there is one gray component with index 0.
	FormatFlagAlpha   FormatFlags = C.GST_VIDEO_FORMAT_FLAG_ALPHA   // (8) – The video format has an alpha components with the number 3.
	FormatFlagLE      FormatFlags = C.GST_VIDEO_FORMAT_FLAG_LE      // (16) – The video format has data stored in little endianness.
	FormatFlagPalette FormatFlags = C.GST_VIDEO_FORMAT_FLAG_PALETTE // (32) – The video format has a palette. The palette is stored in the second plane and indexes are stored in the first plane.
	FormatFlagComplex FormatFlags = C.GST_VIDEO_FORMAT_FLAG_COMPLEX // (64) – The video format has a complex layout that can't be described with the usual information in the GstVideoFormatInfo.
	FormatFlagUnpack  FormatFlags = C.GST_VIDEO_FORMAT_FLAG_UNPACK  // (128) – This format can be used in a GstVideoFormatUnpack and GstVideoFormatPack function.
	FormatFlagTiled   FormatFlags = C.GST_VIDEO_FORMAT_FLAG_TILED   // (256) – The format is tiled, there is tiling information in the last plane.
)

// PackFlags are different flags that can be used when packing and unpacking.
type PackFlags int

// Type castings
const (
	PackFlagNone          PackFlags = C.GST_VIDEO_PACK_FLAG_NONE           // (0) – No flag
	PackFlagTruncateRange PackFlags = C.GST_VIDEO_PACK_FLAG_TRUNCATE_RANGE // (1) – When the source has a smaller depth than the target format, set the least significant bits of the target to 0. This is likely slightly faster but less accurate. When this flag is not specified, the most significant bits of the source are duplicated in the least significant bits of the destination.
	PackFlagInterlaced    PackFlags = C.GST_VIDEO_PACK_FLAG_INTERLACED     // (2) – The source is interlaced. The unpacked format will be interlaced as well with each line containing information from alternating fields. (Since: 1.2)
)

// FormatInfo contains information for a video format.
type FormatInfo struct {
	ptr *C.GstVideoFormatInfo
}

func (f *FormatInfo) instance() *C.GstVideoFormatInfo { return f.ptr }

// Bits returns the number of bits used to pack data items. This can be less than 8 when
// multiple pixels are stored in a byte. for values > 8 multiple bytes should be read
// according to the endianness flag before applying the shift and mask.
func (f *FormatInfo) Bits() uint { return uint(C.formatInfoBits(f.instance())) }

// ComponentDepth returns the depth in bits for the given component.
func (f *FormatInfo) ComponentDepth(component uint) uint {
	return uint(C.formatInfoDepth(f.instance(), C.guint(component)))
}

// ComponentHSub returns the subsampling factor of the height for the component.
func (f *FormatInfo) ComponentHSub(component uint) uint {
	return uint(C.formatInfoHSub(f.instance(), C.guint(component)))
}

// ComponentWSub returns the subsampling factor of the width for the component.
func (f *FormatInfo) ComponentWSub(n uint) uint {
	return uint(C.formatInfoWSub(f.instance(), C.guint(n)))
}

// Flags returns the flags on this info.
func (f *FormatInfo) Flags() FormatFlags { return FormatFlags(C.formatInfoFlags(f.instance())) }

// Format returns the format for this info.
func (f *FormatInfo) Format() Format { return Format(C.formatInfoFormat(f.instance())) }

// HasAlpha returns true if the alpha flag is set.
func (f *FormatInfo) HasAlpha() bool { return gobool(C.formatInfoHasAlpha(f.instance())) }

// HasPalette returns true if this info has a palette.
func (f *FormatInfo) HasPalette() bool { return gobool(C.formatInfoHasPalette(f.instance())) }

// IsComplex returns true if the complex flag is set.
func (f *FormatInfo) IsComplex() bool { return gobool(C.formatInfoIsComplex(f.instance())) }

// IsGray returns true if the gray flag is set.
func (f *FormatInfo) IsGray() bool { return gobool(C.formatInfoIsGray(f.instance())) }

// IsLE returns true if the LE flag is set.
func (f *FormatInfo) IsLE() bool { return gobool(C.formatInfoIsLE(f.instance())) }

// IsRGB returns true if the RGB flag is set.
func (f *FormatInfo) IsRGB() bool { return gobool(C.formatInfoIsRGB(f.instance())) }

// IsTiled returns true if the tiled flag is set.
func (f *FormatInfo) IsTiled() bool { return gobool(C.formatInfoIsTiled(f.instance())) }

// IsYUV returns true if the YUV flag is set.
func (f *FormatInfo) IsYUV() bool { return gobool(C.formatInfoIsYUV(f.instance())) }

// Name returns a human readable name for this info.
func (f *FormatInfo) Name() string { return C.GoString(C.formatInfoName(f.instance())) }

// NumComponents returns the number of components in this info.
func (f *FormatInfo) NumComponents() uint { return uint(C.formatInfoNComponent(f.instance())) }

// NumPlanes returns the number of planes in this info.
func (f *FormatInfo) NumPlanes() uint { return uint(C.formatInfoNPlanes(f.instance())) }

// Plane returns the given plane index.
func (f *FormatInfo) Plane(n uint) uint { return uint(C.formatInfoPlane(f.instance(), C.guint(n))) }

// PlaneOffset returns the offset for the given plane.
func (f *FormatInfo) PlaneOffset(n uint) uint {
	return uint(C.formatInfoPOffset(f.instance(), C.guint(n)))
}

// PlaneStride returns the stride for the given plane.
func (f *FormatInfo) PlaneStride(n uint) uint {
	return uint(C.formatInfoPStride(f.instance(), C.guint(n)))
}

// TileHS returns the height of a tile, in bytes, represented as a shift.
func (f *FormatInfo) TileHS() uint { return uint(C.formatInfoTileHS(f.instance())) }

// TileMode returns the tiling mode.
func (f *FormatInfo) TileMode() TileMode { return TileMode(C.formatInfoTileMode(f.instance())) }

// TileWS returns the width of a tile, in bytes, represented as a shift.
func (f *FormatInfo) TileWS() uint { return uint(C.formatInfoTileWS(f.instance())) }
