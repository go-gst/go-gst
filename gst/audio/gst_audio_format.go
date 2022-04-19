package audio

/*
#include "gst.go.h"
*/
import "C"
import (
	"math"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

func init() {
	glib.RegisterGValueMarshalers([]glib.TypeMarshaler{
		{
			T: glib.Type(C.gst_audio_format_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_enum(uintptrToGVal(p))
				return Format(c), nil
			},
		},
	})
}

// FormatInfo is a structure containing information about an audio format.
type FormatInfo struct {
	ptr *C.GstAudioFormatInfo
}

// Format returns the format for this info.
func (f *FormatInfo) Format() Format { return Format(f.ptr.format) }

// Name returns the name for this info.
func (f *FormatInfo) Name() string { return C.GoString(f.ptr.name) }

// Description returns a user readable description for this info.
func (f *FormatInfo) Description() string { return C.GoString(f.ptr.description) }

// Flags returns the flags on this info.
func (f *FormatInfo) Flags() FormatFlags { return FormatFlags(f.ptr.flags) }

// Endianness returns the endianness for this info.
func (f *FormatInfo) Endianness() int { return int(f.ptr.endianness) }

// Width returns the amount of bits used for one sample.
func (f *FormatInfo) Width() int { return int(f.ptr.width) }

// Depth returns the amount of valid bits in width.
func (f *FormatInfo) Depth() int { return int(f.ptr.depth) }

// Silence returns the data for a single silent sample.
func (f *FormatInfo) Silence() []byte {
	return C.GoBytes(unsafe.Pointer(&f.ptr.silence), C.gint(f.Width()/8))
}

// UnpackFormat is the format of unpacked samples.
func (f *FormatInfo) UnpackFormat() Format { return Format(f.ptr.unpack_format) }

// Layout represents the layout of audio samples for different channels.
type Layout int

// Castings for Layouts
const (
	LayoutInterleaved    Layout = C.GST_AUDIO_LAYOUT_INTERLEAVED     // (0) – interleaved audio
	LayoutNonInterleaved Layout = C.GST_AUDIO_LAYOUT_NON_INTERLEAVED // (1) – non-interleaved audio
)

// FormatFlags are the different audio flags that a format can have.
type FormatFlags int

// Castings for FormatFlags
const (
	FormatFlagInteger FormatFlags = C.GST_AUDIO_FORMAT_FLAG_INTEGER // (1) – integer samples
	FormatFlagFloat   FormatFlags = C.GST_AUDIO_FORMAT_FLAG_FLOAT   // (2) – float samples
	FormatFlagSigned  FormatFlags = C.GST_AUDIO_FORMAT_FLAG_SIGNED  // (4) – signed samples
	FormatFlagComplex FormatFlags = C.GST_AUDIO_FORMAT_FLAG_COMPLEX // (16) – complex layout
	FormatFlagUnpack  FormatFlags = C.GST_AUDIO_FORMAT_FLAG_UNPACK  // (32) – the format can be used in GstAudioFormatUnpack and GstAudioFormatPack functions
)

// PackFlags are the different flags that can be used when packing and unpacking.
type PackFlags int

// Castings for PackFlags
const (
	PackFlagNone          PackFlags = C.GST_AUDIO_PACK_FLAG_NONE           // (0) – No flag
	PackFlagTruncateRange PackFlags = C.GST_AUDIO_PACK_FLAG_TRUNCATE_RANGE // (1) – When the source has a smaller depth than the target format, set the least significant bits of the target to 0. This is likely slightly faster but less accurate. When this flag is not specified, the most significant bits of the source are duplicated in the least significant bits of the destination.
)

// Format is an enum describing the most common audio formats
type Format int

// TypeFormat is the GType for a GstAudioFormat.
var TypeFormat = glib.Type(C.gst_audio_format_get_type())

// ToGValue implements a glib.ValueTransformer
func (f Format) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeFormat)
	if err != nil {
		return nil, err
	}
	val.SetEnum(int(f))
	return val, nil
}

func (f Format) String() string {
	return C.GoString(C.gst_audio_format_to_string(C.GstAudioFormat(f)))
}

// Info returns the info for this Format.
func (f Format) Info() *FormatInfo {
	return &FormatInfo{
		ptr: C.gst_audio_format_get_info(C.GstAudioFormat(f)),
	}
}

// FormatFromString returns the format for the given string representation.
func FormatFromString(format string) Format {
	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))
	return Format(C.gst_audio_format_from_string((*C.gchar)(unsafe.Pointer(f))))
}

// LittleEndian represents little-endian format
const LittleEndian int = C.G_LITTLE_ENDIAN

// BigEndian represents big-endian format
const BigEndian int = C.G_BIG_ENDIAN

// FormatFromInteger returns a Format with the given parameters. Signed is whether signed or unsigned format.
// Endianness can be either LittleEndian or BigEndian. Width is the amount of bits used per sample, and depth
// is the amount of used bits in width.
func FormatFromInteger(signed bool, endianness, width, depth int) Format {
	return Format(C.gst_audio_format_build_integer(
		gboolean(signed),
		C.gint(endianness),
		C.gint(width),
		C.gint(depth),
	))
}

// RawFormats returns all the raw formats supported by GStreamer.
func RawFormats() []Format {
	var l C.guint
	formats := C.gst_audio_formats_raw(&l)
	out := make([]Format, int(l))
	tmpslice := (*[(math.MaxInt32 - 1) / unsafe.Sizeof(C.GST_AUDIO_FORMAT_UNKNOWN)]C.GstAudioFormat)(unsafe.Pointer(formats))[:l:l]
	for i, s := range tmpslice {
		out[i] = Format(s)
	}
	return out
}

// MakeRawCaps returns a generic raw audio caps for the formats defined. If formats is nil, all supported
// formats are used.
func MakeRawCaps(formats []Format, layout Layout) *gst.Caps {
	var caps *C.GstCaps
	if formats == nil {
		caps = C.gst_audio_make_raw_caps(nil, C.guint(0), C.GstAudioLayout(layout))
	} else {
		caps = C.gst_audio_make_raw_caps((*C.GstAudioFormat)(unsafe.Pointer(&formats[0])), C.guint(len(formats)), C.GstAudioLayout(layout))
	}
	if caps == nil {
		return nil
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// DefaultRate is the default sampling rate used in consumer audio
const DefaultRate = 44100

// Castings for Formats
const (
	// DefaultFormat is the default format used in consumer audio
	DefaultFormat Format = FormatS16LE

	FormatUnknown Format = C.GST_AUDIO_FORMAT_UNKNOWN  // (0) – unknown or unset audio format
	FormatEncoded Format = C.GST_AUDIO_FORMAT_ENCODED  // (1) – encoded audio format
	FormatS8      Format = C.GST_AUDIO_FORMAT_S8       // (2) – 8 bits in 8 bits, signed
	FormatU8      Format = C.GST_AUDIO_FORMAT_U8       // (3) – 8 bits in 8 bits, unsigned
	FormatS16LE   Format = C.GST_AUDIO_FORMAT_S16LE    // (4) – 16 bits in 16 bits, signed, little endian
	FormatS16BE   Format = C.GST_AUDIO_FORMAT_S16BE    // (5) – 16 bits in 16 bits, signed, big endian
	FormatU16LE   Format = C.GST_AUDIO_FORMAT_U16LE    // (6) – 16 bits in 16 bits, unsigned, little endian
	FormatU16BE   Format = C.GST_AUDIO_FORMAT_U16BE    // (7) – 16 bits in 16 bits, unsigned, big endian
	FormatS2432LE Format = C.GST_AUDIO_FORMAT_S24_32LE // (8) – 24 bits in 32 bits, signed, little endian
	FormatS2332BE Format = C.GST_AUDIO_FORMAT_S24_32BE // (9) – 24 bits in 32 bits, signed, big endian
	FormatU2432LE Format = C.GST_AUDIO_FORMAT_U24_32LE // (10) – 24 bits in 32 bits, unsigned, little endian
	FormatU2432BE Format = C.GST_AUDIO_FORMAT_U24_32BE // (11) – 24 bits in 32 bits, unsigned, big endian
	FormatS32LE   Format = C.GST_AUDIO_FORMAT_S32LE    // (12) – 32 bits in 32 bits, signed, little endian
	FormatS32BE   Format = C.GST_AUDIO_FORMAT_S32BE    // (13) – 32 bits in 32 bits, signed, big endian
	FormatU32LE   Format = C.GST_AUDIO_FORMAT_U32LE    // (14) – 32 bits in 32 bits, unsigned, little endian
	FormatU32BE   Format = C.GST_AUDIO_FORMAT_U32BE    // (15) – 32 bits in 32 bits, unsigned, big endian
	FormatS24LE   Format = C.GST_AUDIO_FORMAT_S24LE    // (16) – 24 bits in 24 bits, signed, little endian
	FormatS24BE   Format = C.GST_AUDIO_FORMAT_S24BE    // (17) – 24 bits in 24 bits, signed, big endian
	FormatU24LE   Format = C.GST_AUDIO_FORMAT_U24LE    // (18) – 24 bits in 24 bits, unsigned, little endian
	FormatU24BE   Format = C.GST_AUDIO_FORMAT_U24BE    // (19) – 24 bits in 24 bits, unsigned, big endian
	FormatS20LE   Format = C.GST_AUDIO_FORMAT_S20LE    // (20) – 20 bits in 24 bits, signed, little endian
	FormatS20BE   Format = C.GST_AUDIO_FORMAT_S20BE    // (21) – 20 bits in 24 bits, signed, big endian
	FormatU20LE   Format = C.GST_AUDIO_FORMAT_U20LE    // (22) – 20 bits in 24 bits, unsigned, little endian
	FormatU20BE   Format = C.GST_AUDIO_FORMAT_U20BE    // (23) – 20 bits in 24 bits, unsigned, big endian
	FormatS18LE   Format = C.GST_AUDIO_FORMAT_S18LE    // (24) – 18 bits in 24 bits, signed, little endian
	FormatS18BE   Format = C.GST_AUDIO_FORMAT_S18BE    // (25) – 18 bits in 24 bits, signed, big endian
	FormatU18LE   Format = C.GST_AUDIO_FORMAT_U18LE    // (26) – 18 bits in 24 bits, unsigned, little endian
	FormatU18BE   Format = C.GST_AUDIO_FORMAT_U18BE    // (27) – 18 bits in 24 bits, unsigned, big endian
	FormatF32LE   Format = C.GST_AUDIO_FORMAT_F32LE    // (28) – 32-bit floating point samples, little endian
	FormatF32BE   Format = C.GST_AUDIO_FORMAT_F32BE    // (29) – 32-bit floating point samples, big endian
	FormatF64LE   Format = C.GST_AUDIO_FORMAT_F64LE    // (30) – 64-bit floating point samples, little endian
	FormatF64BE   Format = C.GST_AUDIO_FORMAT_F64BE    // (31) – 64-bit floating point samples, big endian
	FormatS16     Format = C.GST_AUDIO_FORMAT_S16      // (4) – 16 bits in 16 bits, signed, native endianness
	FormatU16     Format = C.GST_AUDIO_FORMAT_U16      // (6) – 16 bits in 16 bits, unsigned, native endianness
	FormatS2432   Format = C.GST_AUDIO_FORMAT_S24_32   // (8) – 24 bits in 32 bits, signed, native endianness
	FormatU2432   Format = C.GST_AUDIO_FORMAT_U24_32   // (10) – 24 bits in 32 bits, unsigned, native endianness
	FormatS32     Format = C.GST_AUDIO_FORMAT_S32      // (12) – 32 bits in 32 bits, signed, native endianness
	FormatU32     Format = C.GST_AUDIO_FORMAT_U32      // (14) – 32 bits in 32 bits, unsigned, native endianness
	FormatS24     Format = C.GST_AUDIO_FORMAT_S24      // (16) – 24 bits in 24 bits, signed, native endianness
	FormatU24     Format = C.GST_AUDIO_FORMAT_U24      // (18) – 24 bits in 24 bits, unsigned, native endianness
	FormatS20     Format = C.GST_AUDIO_FORMAT_S20      // (20) – 20 bits in 24 bits, signed, native endianness
	FormatU20     Format = C.GST_AUDIO_FORMAT_U20      // (22) – 20 bits in 24 bits, unsigned, native endianness
	FormatS18     Format = C.GST_AUDIO_FORMAT_S18      // (24) – 18 bits in 24 bits, signed, native endianness
	FormatU18     Format = C.GST_AUDIO_FORMAT_U18      // (26) – 18 bits in 24 bits, unsigned, native endianness
	FormatF32     Format = C.GST_AUDIO_FORMAT_F32      // (28) – 32-bit floating point samples, native endianness
	FormatF64     Format = C.GST_AUDIO_FORMAT_F64      // (30) – 64-bit floating point samples, native endianness
)
