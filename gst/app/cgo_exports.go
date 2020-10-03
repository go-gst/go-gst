package app

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-gst/gst"
)

func getCbsFromPtr(userData C.gpointer) *SinkCallbacks {
	ptr := gopointer.Restore(unsafe.Pointer(userData))
	cbs, ok := ptr.(*SinkCallbacks)
	if !ok {
		gopointer.Unref(unsafe.Pointer(userData))
		return nil
	}
	return cbs
}

func wrapCSink(sink *C.GstAppSink) *Sink {
	return wrapAppSink(gst.FromGstElementUnsafe(unsafe.Pointer(sink)))
}

//export goSinkEOSCb
func goSinkEOSCb(sink *C.GstAppSink, userData C.gpointer) {
	cbs := getCbsFromPtr(userData)
	if cbs == nil {
		return
	}
	if cbs.EOSFunc == nil {
		return
	}
	cbs.EOSFunc(wrapCSink(sink))
}

//export goSinkNewPrerollCb
func goSinkNewPrerollCb(sink *C.GstAppSink, userData C.gpointer) C.GstFlowReturn {
	cbs := getCbsFromPtr(userData)
	if cbs == nil {
		return C.GstFlowReturn(gst.FlowError)
	}
	if cbs.NewPrerollFunc == nil {
		return C.GstFlowReturn(gst.FlowOK)
	}
	return C.GstFlowReturn(cbs.NewPrerollFunc(wrapCSink(sink)))
}

//export goSinkNewSampleCb
func goSinkNewSampleCb(sink *C.GstAppSink, userData C.gpointer) C.GstFlowReturn {
	cbs := getCbsFromPtr(userData)
	if cbs == nil {
		return C.GstFlowReturn(gst.FlowError)
	}
	if cbs.NewSampleFunc == nil {
		return C.GstFlowReturn(gst.FlowOK)
	}
	return C.GstFlowReturn(cbs.NewSampleFunc(wrapCSink(sink)))
}

//export goSinkGDestroyNotifyFunc
func goSinkGDestroyNotifyFunc(ptr C.gpointer) {
	gopointer.Unref(unsafe.Pointer(ptr))
}
