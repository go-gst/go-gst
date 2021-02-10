package video

/*
#include <gst/gst.h>

GValue *  videoUtilToGValue (guintptr p) { return (GValue*)(p); }
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

func uintptrToGVal(p uintptr) *C.GValue { return (*C.GValue)(C.videoUtilToGValue(C.guintptr(p))) }

func wrapGerr(gerr *C.GError) error {
	defer C.g_error_free(gerr)
	return errors.New(C.GoString(gerr.message))
}

func durationToClockTime(d time.Duration) C.GstClockTime {
	if d.Nanoseconds() < 0 {
		return C.GstClockTime(gst.ClockTimeNone)
	}
	return C.GstClockTime(d.Nanoseconds())
}

func fromCoreCaps(caps *gst.Caps) *C.GstCaps {
	return (*C.GstCaps)(unsafe.Pointer(caps.Instance()))
}

func fromCoreCapsFeatures(feats *gst.CapsFeatures) *C.GstCapsFeatures {
	return (*C.GstCapsFeatures)(unsafe.Pointer(feats.Instance()))
}

func fromCoreElement(elem *gst.Element) *C.GstElement {
	return (*C.GstElement)(unsafe.Pointer(elem.Instance()))
}

func fromCoreEvent(event *gst.Event) *C.GstEvent {
	return (*C.GstEvent)(unsafe.Pointer(event.Instance()))
}

func fromCoreMessage(msg *gst.Message) *C.GstMessage {
	return (*C.GstMessage)(unsafe.Pointer(msg.Instance()))
}

func fromCoreQuery(query *gst.Query) *C.GstQuery {
	return (*C.GstQuery)(unsafe.Pointer(query.Instance()))
}

func fromCoreSample(sample *gst.Sample) *C.GstSample {
	return (*C.GstSample)(unsafe.Pointer(sample.Instance()))
}

func fromCoreStructure(structure *gst.Structure) *C.GstStructure {
	return (*C.GstStructure)(unsafe.Pointer(structure.Instance()))
}

func gobool(b C.gboolean) bool { return int(b) > 0 }
