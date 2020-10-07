package video

// #include <gst/video/video.h>
import "C"

// ColorMatrix is used to convert between Y'PbPr and non-linear RGB (R'G'B')
type ColorMatrix int

// Type castings
const (
	ColorMatrixUnknown   ColorMatrix = C.GST_VIDEO_COLOR_MATRIX_UNKNOWN   // (0) – unknown matrix
	ColorMatrixRGB       ColorMatrix = C.GST_VIDEO_COLOR_MATRIX_RGB       // (1) – identity matrix. Order of coefficients is actually GBR, also IEC 61966-2-1 (sRGB)
	ColorMatrixFCC       ColorMatrix = C.GST_VIDEO_COLOR_MATRIX_FCC       // (2) – FCC Title 47 Code of Federal Regulations 73.682 (a)(20)
	ColorMatrixBT709     ColorMatrix = C.GST_VIDEO_COLOR_MATRIX_BT709     // (3) – ITU-R BT.709 color matrix, also ITU-R BT1361 / IEC 61966-2-4 xvYCC709 / SMPTE RP177 Annex B
	ColorMatrixBT601     ColorMatrix = C.GST_VIDEO_COLOR_MATRIX_BT601     // (4) – ITU-R BT.601 color matrix, also SMPTE170M / ITU-R BT1358 525 / ITU-R BT1700 NTSC
	ColorMatrixSMPTE240M ColorMatrix = C.GST_VIDEO_COLOR_MATRIX_SMPTE240M // (5) – SMPTE 240M color matrix
	ColorMatrixBT2020    ColorMatrix = C.GST_VIDEO_COLOR_MATRIX_BT2020    // (6) – ITU-R BT.2020 color matrix. Since: 1.6
)

// ColorPrimaries define the how to transform linear RGB values to and from the
// CIE XYZ colorspace.
type ColorPrimaries int

// Type castings
const (
	ColorPrimariesUnknown    ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_UNKNOWN    // (0) – unknown color primaries
	ColorPrimariesBT709      ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_BT709      // (1) – BT709 primaries, also ITU-R BT1361 / IEC 61966-2-4 / SMPTE RP177 Annex B
	ColorPrimariesBT470M     ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_BT470M     // (2) – BT470M primaries, also FCC Title 47 Code of Federal Regulations 73.682 (a)(20)
	ColorPrimariesBT470BG    ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_BT470BG    // (3) – BT470BG primaries, also ITU-R BT601-6 625 / ITU-R BT1358 625 / ITU-R BT1700 625 PAL & SECAM
	ColorPrimariesSMPTE170M  ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_SMPTE170M  // (4) – SMPTE170M primaries, also ITU-R BT601-6 525 / ITU-R BT1358 525 / ITU-R BT1700 NTSC
	ColorPrimariesSMPTE240M  ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_SMPTE240M  // (5) – SMPTE240M primaries
	ColorPrimariesFilm       ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_FILM       // (6) – Generic film (colour filters using Illuminant C)
	ColorPrimariesBT2020     ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_BT2020     // (7) – ITU-R BT2020 primaries. Since: 1.6
	ColorPrimariesAdobeRGB   ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_ADOBERGB   // (8) – Adobe RGB primaries. Since: 1.8
	ColorPrimariesSMPTEST428 ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_SMPTEST428 // (9) – SMPTE ST 428 primaries (CIE 1931 XYZ). Since: 1.16
	ColorPrimariesSMPTERP431 ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_SMPTERP431 // (10) – SMPTE RP 431 primaries (ST 431-2 (2011) / DCI P3). Since: 1.16
	ColorPrimariesSMPTEEG432 ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_SMPTEEG432 // (11) – SMPTE EG 432 primaries (ST 432-1 (2010) / P3 D65). Since: 1.16
	ColorPrimariesEBU3213    ColorPrimaries = C.GST_VIDEO_COLOR_PRIMARIES_EBU3213    // (12) – EBU 3213 primaries (JEDEC P22 phosphors). Since: 1.16
)

// ColorRange represents possible color range values. These constants are defined for
// 8 bit color values and can be scaled for other bit depths.
type ColorRange int

// Type castings
const (
	ColorRangeUnknown ColorRange = C.GST_VIDEO_COLOR_RANGE_UNKNOWN // (0) – unknown range
	ColorRange0255    ColorRange = C.GST_VIDEO_COLOR_RANGE_0_255   // (1) – [0..255] for 8 bit components
	ColorRange16235   ColorRange = C.GST_VIDEO_COLOR_RANGE_16_235  // (2) – [16..235] for 8 bit components. Chroma has [16..240] range.
)

// TransferFunction defines the formula for converting between non-linear RGB (R'G'B')
// and linear RGB
type TransferFunction int

// Type castings
const (
	TransferUnknown    TransferFunction = C.GST_VIDEO_TRANSFER_UNKNOWN      // (0) – unknown transfer function
	TransferGamma10    TransferFunction = C.GST_VIDEO_TRANSFER_GAMMA10      // (1) – linear RGB, gamma 1.0 curve
	TransferGamma18    TransferFunction = C.GST_VIDEO_TRANSFER_GAMMA18      // (2) – Gamma 1.8 curve
	TransferGamma20    TransferFunction = C.GST_VIDEO_TRANSFER_GAMMA20      // (3) – Gamma 2.0 curve
	TransferGamma22    TransferFunction = C.GST_VIDEO_TRANSFER_GAMMA22      // (4) – Gamma 2.2 curve
	TransferBT709      TransferFunction = C.GST_VIDEO_TRANSFER_BT709        // (5) – Gamma 2.2 curve with a linear segment in the lower range, also ITU-R BT470M / ITU-R BT1700 625 PAL & SECAM / ITU-R BT1361
	TransferSMPTE240M  TransferFunction = C.GST_VIDEO_TRANSFER_SMPTE240M    // (6) – Gamma 2.2 curve with a linear segment in the lower range
	TransferSRGB       TransferFunction = C.GST_VIDEO_TRANSFER_SRGB         // (7) – Gamma 2.4 curve with a linear segment in the lower range. IEC 61966-2-1 (sRGB or sYCC)
	TransferGamma28    TransferFunction = C.GST_VIDEO_TRANSFER_GAMMA28      // (8) – Gamma 2.8 curve, also ITU-R BT470BG
	TransferLog100     TransferFunction = C.GST_VIDEO_TRANSFER_LOG100       // (9) – Logarithmic transfer characteristic 100:1 range
	TransferLog316     TransferFunction = C.GST_VIDEO_TRANSFER_LOG316       // (10) – Logarithmic transfer characteristic 316.22777:1 range (100 * sqrt(10) : 1)
	TransferBT202012   TransferFunction = C.GST_VIDEO_TRANSFER_BT2020_12    // (11) – Gamma 2.2 curve with a linear segment in the lower range. Used for BT.2020 with 12 bits per component. Since: 1.6
	TransferAdobeRGB   TransferFunction = C.GST_VIDEO_TRANSFER_ADOBERGB     // (12) – Gamma 2.19921875. Since: 1.8
	TransferBT202010   TransferFunction = C.GST_VIDEO_TRANSFER_BT2020_10    // (13) – Rec. ITU-R BT.2020-2 with 10 bits per component. (functionally the same as the values GST_VIDEO_TRANSFER_BT709 and GST_VIDEO_TRANSFER_BT601). Since: 1.18
	TransferSMPTE2084  TransferFunction = C.GST_VIDEO_TRANSFER_SMPTE2084    // (14) – SMPTE ST 2084 for 10, 12, 14, and 16-bit systems. Known as perceptual quantization (PQ) Since: 1.18
	TransferARIBSTDB67 TransferFunction = C.GST_VIDEO_TRANSFER_ARIB_STD_B67 // (15) – Association of Radio Industries and Businesses (ARIB) STD-B67 and Rec. ITU-R BT.2100-1 hybrid loggamma (HLG) system Since: 1.18
	TransferBT601      TransferFunction = C.GST_VIDEO_TRANSFER_BT601        // (16) – also known as SMPTE170M / ITU-R BT1358 525 or 625 / ITU-R BT1700 NTSC
)

// Pre-defined colorimetries
const (
	ColorimetryBT2020    string = C.GST_VIDEO_COLORIMETRY_BT2020
	ColorimetryBT202010  string = C.GST_VIDEO_COLORIMETRY_BT2020_10
	ColorimetryBT2100HLG string = C.GST_VIDEO_COLORIMETRY_BT2100_HLG
	ColorimetryBT2100PQ  string = C.GST_VIDEO_COLORIMETRY_BT2100_PQ
	ColorimetryBT601     string = C.GST_VIDEO_COLORIMETRY_BT601
	ColorimetryBT709     string = C.GST_VIDEO_COLORIMETRY_BT709
	ColorimetrySMPTE240M string = C.GST_VIDEO_COLORIMETRY_SMPTE240M
	ColorimetrySRRGB     string = C.GST_VIDEO_COLORIMETRY_SRGB
)

// ColorPrimariesInfo is a structure describing the chromaticity coordinates of an RGB system.
// These values can be used to construct a matrix to transform RGB to and from the XYZ colorspace.
type ColorPrimariesInfo struct {
	Primaries ColorPrimaries
	Wx, Wy    float64 // Reference white coordinates
	Rx, Ry    float64 // Red coordinates
	Gx, Gy    float64 // Green coordinates
	Bx, By    float64 // Blue coordinates
}

// func (c *ColorPrimariesInfo) instance() *C.GstVideoColorPrimariesInfo {
// 	i := &C.GstVideoColorPrimariesInfo{
// 		primaries: C.GstVideoColorPrimaries(c.Primaries),
// 		Wx:        C.gdouble(c.Wx),
// 		Wy:        C.gdouble(c.Wy),
// 		Rx:        C.gdouble(c.Rx),
// 		Ry:        C.gdouble(c.Ry),
// 		Gx:        C.gdouble(c.Gx),
// 		Gy:        C.gdouble(c.Gy),
// 		Bx:        C.gdouble(c.Bx),
// 		By:        C.gdouble(c.By),
// 	}
// 	runtime.SetFinalizer(c, func(_ *ColorPrimariesInfo) { C.g_free((C.gpointer)(unsafe.Pointer(i))) })
// 	return i
// }

// Colorimetry is a structure describing the color info.
type Colorimetry struct {
	// The color range. This is the valid range for the samples. It is used to convert the samples to Y'PbPr values.
	Range ColorRange
	// The color matrix. Used to convert between Y'PbPr and non-linear RGB (R'G'B').
	Matrix ColorMatrix
	// The transfer function. used to convert between R'G'B' and RGB.
	Transfer TransferFunction
	// Color primaries. used to convert between R'G'B' and CIE XYZ.
	Primaries ColorPrimaries
}

// func (c *Colorimetry) instance() *C.GstVideoColorimetry {
// 	i := &C.GstVideoColorimetry{
// 		_range:    C.GstVideoColorRange(c.Range),
// 		matrix:    C.GstVideoColorMatrix(c.Matrix),
// 		transfer:  C.GstVideoTransferFunction(c.Transfer),
// 		primaries: C.GstVideoColorPrimaries(c.Primaries),
// 	}
// 	runtime.SetFinalizer(c, func(_ *Colorimetry) { C.g_free((C.gpointer)(unsafe.Pointer(i))) })
// 	return i
// }

func colorimetryFromInstance(c C.GstVideoColorimetry) *Colorimetry {
	return &Colorimetry{
		Range:     ColorRange(c._range),
		Matrix:    ColorMatrix(c.matrix),
		Transfer:  TransferFunction(c.transfer),
		Primaries: ColorPrimaries(c.primaries),
	}
}
