package app

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-gst/gst"
)

func getSinkCbsFromPtr(userData C.gpointer) *SinkCallbacks {
	ptr := gopointer.Restore(unsafe.Pointer(userData))
	cbs, ok := ptr.(*SinkCallbacks)
	if !ok {
		gopointer.Unref(unsafe.Pointer(userData))
		return nil
	}
	return cbs
}

func getSrcCbsFromPtr(userData C.gpointer) *SourceCallbacks {
	ptr := gopointer.Restore(unsafe.Pointer(userData))
	cbs, ok := ptr.(*SourceCallbacks)
	if !ok {
		gopointer.Unref(unsafe.Pointer(userData))
		return nil
	}
	return cbs
}

func wrapCSink(sink *C.GstAppSink) *Sink {
	return wrapAppSink(gst.FromGstElementUnsafeNone(unsafe.Pointer(sink)))
}

func wrapCSource(src *C.GstAppSrc) *Source {
	return wrapAppSrc(gst.FromGstElementUnsafeNone(unsafe.Pointer(src)))
}

//export goNeedDataCb
func goNeedDataCb(src *C.GstAppSrc, length C.guint, userData C.gpointer) {
	cbs := getSrcCbsFromPtr(userData)
	if cbs == nil {
		return
	}
	if cbs.NeedDataFunc == nil {
		return
	}
	cbs.NeedDataFunc(wrapCSource(src), uint(length))
}

//export goEnoughDataDb
func goEnoughDataDb(src *C.GstAppSrc, userData C.gpointer) {
	cbs := getSrcCbsFromPtr(userData)
	if cbs == nil {
		return
	}
	if cbs.EnoughDataFunc == nil {
		return
	}
	cbs.EnoughDataFunc(wrapCSource(src))
}

//export goSeekDataCb
func goSeekDataCb(src *C.GstAppSrc, offset C.guint64, userData C.gpointer) C.gboolean {
	cbs := getSrcCbsFromPtr(userData)
	if cbs == nil {
		return gboolean(false)
	}
	if cbs.SeekDataFunc == nil {
		return gboolean(true)
	}
	return gboolean(cbs.SeekDataFunc(wrapCSource(src), uint64(offset)))
}

//export goSinkEOSCb
func goSinkEOSCb(sink *C.GstAppSink, userData C.gpointer) {
	cbs := getSinkCbsFromPtr(userData)
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
	cbs := getSinkCbsFromPtr(userData)
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
	cbs := getSinkCbsFromPtr(userData)
	if cbs == nil {
		return C.GstFlowReturn(gst.FlowError)
	}
	if cbs.NewSampleFunc == nil {
		return C.GstFlowReturn(gst.FlowOK)
	}
	return C.GstFlowReturn(cbs.NewSampleFunc(wrapCSink(sink)))
}

//export goAppGDestroyNotifyFunc
func goAppGDestroyNotifyFunc(ptr C.gpointer) {
	gopointer.Unref(unsafe.Pointer(ptr))
}
