package gstapp

import (
	"runtime"

	"github.com/go-gst/go-glib/pkg/core/userdata"
)

// #cgo pkg-config: gstreamer-app-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/app/app.h>
// extern void _gogst_gstapp1_AppSrc_need_data(GstAppSrc*, guint, gpointer);
// extern void _gogst_gstapp1_AppSrc_enough_data(GstAppSrc*, gpointer);
// extern gboolean _gogst_gstapp1_AppSrc_seek_data(GstAppSrc*, guint64, gpointer);
// extern void destroyUserdata(gpointer);
import "C"

// AppSrcCallbacks contains the callbacks that can be installed on an [AppSrc] with
// [AppSrc.SetCallbacks]. All callbacks are optional, nil callbacks are not installed.
//
// All callbacks are called from the streaming thread, so they must not block for long.
//
// see also https://gstreamer.freedesktop.org/documentation/app/gstappsrc.html#GstAppSrcCallbacks
type AppSrcCallbacks struct {
	// NeedData gets invoked when the source wants more input. Respond by pushing buffers or
	// signalling end of stream, here or on another goroutine. length is a hint for the
	// amount of bytes wanted, math.MaxUint32 means that any amount is fine.
	NeedData func(appsrc AppSrc, length uint)

	// EnoughData gets invoked when the internal queue of the source is full. Stop pushing
	// until NeedData is invoked again.
	EnoughData func(appsrc AppSrc)

	// SeekData gets invoked when the source must continue producing data from offset. It is
	// only used for seekable stream types. Return true when the seek succeeded.
	SeekData func(appsrc AppSrc, offset uint64) bool
}

// AppSrcExtManual contains the handwritten methods of [AppSrc].
type AppSrcExtManual interface {
	// SetCallbacks wraps gst_app_src_set_callbacks
	//
	// The callbacks replace the signals of the appsrc. They are called directly, without
	// marshaling the arguments through GValues. Any previously installed callbacks are replaced.
	//
	// Passing nil removes all previously installed callbacks.
	//
	// see also https://gstreamer.freedesktop.org/documentation/app/gstappsrc.html#gst_app_src_set_callbacks
	SetCallbacks(callbacks *AppSrcCallbacks)
}

var _ AppSrcExtManual = (*AppSrcInstance)(nil)

// SetCallbacks wraps gst_app_src_set_callbacks
func (appsrc *AppSrcInstance) SetCallbacks(callbacks *AppSrcCallbacks) {
	var carg0 *C.GstAppSrc         // in, none, converted
	var carg1 C.GstAppSrcCallbacks // in, copied by gstreamer
	var carg2 C.gpointer           // closure
	var carg3 C.GDestroyNotify     // destroy

	carg0 = (*C.GstAppSrc)(UnsafeAppSrcToGlibNone(appsrc))

	if callbacks != nil {
		// copy the struct so that later modifications by the caller have no effect
		cbs := *callbacks

		if cbs.NeedData != nil {
			carg1.need_data = (*[0]byte)(C._gogst_gstapp1_AppSrc_need_data)
		}
		if cbs.EnoughData != nil {
			carg1.enough_data = (*[0]byte)(C._gogst_gstapp1_AppSrc_enough_data)
		}
		if cbs.SeekData != nil {
			carg1.seek_data = (*[0]byte)(C._gogst_gstapp1_AppSrc_seek_data)
		}

		carg2 = C.gpointer(userdata.Register(&cbs))
		carg3 = (C.GDestroyNotify)((*[0]byte)(C.destroyUserdata))
	}

	C.gst_app_src_set_callbacks(carg0, &carg1, carg2, carg3)
	runtime.KeepAlive(appsrc)
}
