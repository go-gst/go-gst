package video

/*
#include <gst/video/video.h>

extern void goVideoGDestroyNotifyFunc (gpointer user_data);
extern void goVideoConvertSampleCb    (GstSample * sample, GError * gerr, gpointer user_data);

void cgoVideoGDestroyNotifyFunc (gpointer user_data)
{
	goVideoGDestroyNotifyFunc(user_data);
}

void cgoVideoConvertSampleCb (GstSample * sample, GError * gerr, gpointer user_data)
{
	goVideoConvertSampleCb(sample, gerr, user_data);
}
*/
import "C"

import (
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-gst/gst"
)

// ConvertSampleCallback represents a callback from a video convert opereration.
// It contains the converted sample or any error that ocurred.
type ConvertSampleCallback func(*gst.Sample, error)

// ConvertSample converts a raw video buffer into the specified output caps.
//
// The output caps can be any raw video formats or any image formats (jpeg, png, ...).
//
// The width, height and pixel-aspect-ratio can also be specified in the output caps.
func ConvertSample(sample *gst.Sample, toCaps *gst.Caps, timeout time.Duration) (*gst.Sample, error) {
	var gerr *C.GError
	ret := C.gst_video_convert_sample(
		fromCoreSample(sample),
		fromCoreCaps(toCaps),
		durationToClockTime(timeout),
		&gerr,
	)
	if gerr != nil {
		return nil, wrapGerr(gerr)
	}
	if ret == nil {
		return nil, nil
	}
	return gst.FromGstSampleUnsafeFull(unsafe.Pointer(ret)), nil
}

// ConvertSampleAsync converts a raw video buffer into the specified output caps.
//
// The output caps can be any raw video formats or any image formats (jpeg, png, ...).
//
// The width, height and pixel-aspect-ratio can also be specified in the output caps.
//
// The callback will be called after conversion, when an error occurred or if conversion
// didn't finish after timeout.
func ConvertSampleAsync(sample *gst.Sample, toCaps *gst.Caps, timeout time.Duration, cb ConvertSampleCallback) {
	ptr := gopointer.Save(cb)
	C.gst_video_convert_sample_async(
		fromCoreSample(sample),
		fromCoreCaps(toCaps),
		durationToClockTime(timeout),
		C.GstVideoConvertSampleCallback(C.cgoVideoConvertSampleCb),
		(C.gpointer)(unsafe.Pointer(ptr)),
		C.GDestroyNotify(C.cgoVideoGDestroyNotifyFunc),
	)
}
