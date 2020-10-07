package video

// #include <gst/video/video.h>
import "C"
import (
	"runtime"
	"unsafe"
)

// ChromaFlags are extra flags that influence the result from NewChromaResample.
type ChromaFlags int

// Type castings
const (
	ChromaFlagNone       ChromaFlags = C.GST_VIDEO_CHROMA_FLAG_NONE       // (0) – no flags
	ChromaFlagInterlaced ChromaFlags = C.GST_VIDEO_CHROMA_FLAG_INTERLACED // (1) – the input is interlaced
)

// ChromaMethod represents different subsampling and upsampling methods.
type ChromaMethod int

// Type castings
const (
	ChromaMethodNearest ChromaMethod = C.GST_VIDEO_CHROMA_METHOD_NEAREST // (0) – Duplicates the chroma samples when upsampling and drops when subsampling
	ChromaMethodLinear  ChromaMethod = C.GST_VIDEO_CHROMA_METHOD_LINEAR  // (1) – Uses linear interpolation to reconstruct missing chroma and averaging to subsample
)

// ChromaSite represents various Chroma sitings.
type ChromaSite int

// Type castings
const (
	ChromaSiteUnknown  ChromaSite = C.GST_VIDEO_CHROMA_SITE_UNKNOWN   // (0) – unknown cositing
	ChromaSiteNone     ChromaSite = C.GST_VIDEO_CHROMA_SITE_NONE      // (1) – no cositing
	ChromaSiteHCosited ChromaSite = C.GST_VIDEO_CHROMA_SITE_H_COSITED // (2) – chroma is horizontally cosited
	ChromaSiteVCosited ChromaSite = C.GST_VIDEO_CHROMA_SITE_V_COSITED // (4) – chroma is vertically cosited
	ChromaSiteAltLine  ChromaSite = C.GST_VIDEO_CHROMA_SITE_ALT_LINE  // (8) – choma samples are sited on alternate lines
	ChromaSiteCosited  ChromaSite = C.GST_VIDEO_CHROMA_SITE_COSITED   // (6) – chroma samples cosited with luma samples
	ChromaSiteJpeg     ChromaSite = C.GST_VIDEO_CHROMA_SITE_JPEG      // (1) – jpeg style cositing, also for mpeg1 and mjpeg
	ChromaSiteMpeg2    ChromaSite = C.GST_VIDEO_CHROMA_SITE_MPEG2     // (2) – mpeg2 style cositing
	ChromaSiteDV       ChromaSite = C.GST_VIDEO_CHROMA_SITE_DV        // (14) – DV style cositing
)

// String implements a stringer on ChromaSite.
func (c ChromaSite) String() string {
	out := C.gst_video_chroma_to_string(C.GstVideoChromaSite(c))
	defer C.g_free((C.gpointer)(unsafe.Pointer(out)))
	return C.GoString(out)
}

// ChromaResample is a utility object for resampling chroma planes and converting between different chroma sampling sitings.
type ChromaResample struct {
	ptr *C.GstVideoChromaResample
}

// NewChromaResample creates a new resampler object for the given parameters. When h_factor or v_factor is > 0,
// upsampling will be used, otherwise subsampling is performed.
func NewChromaResample(method ChromaMethod, site ChromaSite, flags ChromaFlags, format Format, hFactor, vFactor int) *ChromaResample {
	resample := C.gst_video_chroma_resample_new(
		C.GstVideoChromaMethod(method),
		C.GstVideoChromaSite(site),
		C.GstVideoChromaFlags(flags),
		C.GstVideoFormat(format),
		C.gint(hFactor), C.gint(vFactor),
	)
	if resample == nil {
		return nil
	}
	goResample := &ChromaResample{resample}
	runtime.SetFinalizer(goResample, func(c *ChromaResample) { C.gst_video_chroma_resample_free(c.instance()) })
	return goResample
}

func (c *ChromaResample) instance() *C.GstVideoChromaResample { return c.ptr }

// GetInfo returns the info about the Resample. The resampler must be fed n_lines at a time. The first line
// should be at offset.
func (c *ChromaResample) GetInfo() (nLines uint, offset int) {
	var lines C.guint
	var off C.gint
	C.gst_video_chroma_resample_get_info(c.instance(), &lines, &off)
	return uint(lines), int(off)
}
