package app

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	gopointer "github.com/go-gst/go-pointer"
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
	return wrapAppSink(&gst.Element{
		Object: &gst.Object{
			InitiallyUnowned: &glib.InitiallyUnowned{
				Object: &glib.Object{
					GObject: glib.ToGObject(unsafe.Pointer(sink)),
				},
			},
		},
	})
}

func wrapCSource(src *C.GstAppSrc) *Source {
	return wrapAppSrc(&gst.Element{
		Object: &gst.Object{
			InitiallyUnowned: &glib.InitiallyUnowned{
				Object: &glib.Object{
					GObject: glib.ToGObject(unsafe.Pointer(src)),
				},
			},
		},
	})
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
	gosrc := wrapCSource(src)

	cbs.NeedDataFunc(gosrc, uint(length))
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
	gosrc := wrapCSource(src)
	cbs.EnoughDataFunc(gosrc)
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
	gosrc := wrapCSource(src)

	return gboolean(cbs.SeekDataFunc(gosrc, uint64(offset)))
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
	gosink := wrapCSink(sink)

	cbs.EOSFunc(gosink)
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
	gosink := wrapCSink(sink)

	ret := C.GstFlowReturn(cbs.NewPrerollFunc(gosink))

	return ret
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
	gosink := wrapCSink(sink)

	ret := C.GstFlowReturn(cbs.NewSampleFunc(gosink))

	return ret
}

//export goAppGDestroyNotifyFunc
func goAppGDestroyNotifyFunc(ptr C.gpointer) {
	gopointer.Unref(unsafe.Pointer(ptr))
}
