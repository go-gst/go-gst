package video

/*
#include <gst/gst.h>
*/
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-gst/gst"
	gopointer "github.com/mattn/go-pointer"
)

//export goVideoGDestroyNotifyFunc
func goVideoGDestroyNotifyFunc(ptr C.gpointer) {
	gopointer.Unref(unsafe.Pointer(ptr))
}

//export goVideoConvertSampleCb
func goVideoConvertSampleCb(gsample *C.GstSample, gerr *C.GError, userData C.gpointer) {
	var sample *gst.Sample
	var err error
	if gerr != nil {
		err = wrapGerr(gerr)
	}
	if gsample != nil {
		sample = gst.FromGstSampleUnsafeFull(unsafe.Pointer(gsample))
	}
	iface := gopointer.Restore(unsafe.Pointer(userData))
	if iface == nil {
		return
	}
	cb, ok := iface.(ConvertSampleCallback)
	if !ok {
		return
	}
	cb(sample, err)
}
