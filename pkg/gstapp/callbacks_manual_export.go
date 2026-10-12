package gstapp

import (
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/userdata"
	"github.com/go-gst/go-gst/pkg/gst"
)

// #cgo pkg-config: gstreamer-app-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/app/app.h>
import "C"

func loadAppSinkCallbacks(data C.gpointer) *AppSinkCallbacks {
	v := userdata.Load(unsafe.Pointer(data))
	if v == nil {
		panic(`appsink callbacks not found`)
	}
	return v.(*AppSinkCallbacks)
}

func loadAppSrcCallbacks(data C.gpointer) *AppSrcCallbacks {
	v := userdata.Load(unsafe.Pointer(data))
	if v == nil {
		panic(`appsrc callbacks not found`)
	}
	return v.(*AppSrcCallbacks)
}

//export _gogst_gstapp1_AppSink_eos
func _gogst_gstapp1_AppSink_eos(carg0 *C.GstAppSink, carg1 C.gpointer) {
	cbs := loadAppSinkCallbacks(carg1)

	cbs.EOS(UnsafeAppSinkFromGlibNone(unsafe.Pointer(carg0)))
}

//export _gogst_gstapp1_AppSink_new_preroll
func _gogst_gstapp1_AppSink_new_preroll(carg0 *C.GstAppSink, carg1 C.gpointer) C.GstFlowReturn {
	cbs := loadAppSinkCallbacks(carg1)

	return C.GstFlowReturn(cbs.NewPreroll(UnsafeAppSinkFromGlibNone(unsafe.Pointer(carg0))))
}

//export _gogst_gstapp1_AppSink_new_sample
func _gogst_gstapp1_AppSink_new_sample(carg0 *C.GstAppSink, carg1 C.gpointer) C.GstFlowReturn {
	cbs := loadAppSinkCallbacks(carg1)

	return C.GstFlowReturn(cbs.NewSample(UnsafeAppSinkFromGlibNone(unsafe.Pointer(carg0))))
}

//export _gogst_gstapp1_AppSink_new_event
func _gogst_gstapp1_AppSink_new_event(carg0 *C.GstAppSink, carg1 C.gpointer) C.gboolean {
	cbs := loadAppSinkCallbacks(carg1)

	if cbs.NewEvent(UnsafeAppSinkFromGlibNone(unsafe.Pointer(carg0))) {
		return C.TRUE
	}
	return C.FALSE
}

//export _gogst_gstapp1_AppSink_propose_allocation
func _gogst_gstapp1_AppSink_propose_allocation(carg0 *C.GstAppSink, carg1 *C.GstQuery, carg2 C.gpointer) C.gboolean {
	cbs := loadAppSinkCallbacks(carg2)

	// the query is owned by the caller and must stay writable, so we only borrow it
	query := gst.UnsafeQueryFromGlibBorrow(unsafe.Pointer(carg1))

	if cbs.ProposeAllocation(UnsafeAppSinkFromGlibNone(unsafe.Pointer(carg0)), query) {
		return C.TRUE
	}
	return C.FALSE
}

//export _gogst_gstapp1_AppSrc_need_data
func _gogst_gstapp1_AppSrc_need_data(carg0 *C.GstAppSrc, carg1 C.guint, carg2 C.gpointer) {
	cbs := loadAppSrcCallbacks(carg2)

	cbs.NeedData(UnsafeAppSrcFromGlibNone(unsafe.Pointer(carg0)), uint32(carg1))
}

//export _gogst_gstapp1_AppSrc_enough_data
func _gogst_gstapp1_AppSrc_enough_data(carg0 *C.GstAppSrc, carg1 C.gpointer) {
	cbs := loadAppSrcCallbacks(carg1)

	cbs.EnoughData(UnsafeAppSrcFromGlibNone(unsafe.Pointer(carg0)))
}

//export _gogst_gstapp1_AppSrc_seek_data
func _gogst_gstapp1_AppSrc_seek_data(carg0 *C.GstAppSrc, carg1 C.guint64, carg2 C.gpointer) C.gboolean {
	cbs := loadAppSrcCallbacks(carg2)

	if cbs.SeekData(UnsafeAppSrcFromGlibNone(unsafe.Pointer(carg0)), uint64(carg1)) {
		return C.TRUE
	}
	return C.FALSE
}
