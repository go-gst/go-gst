package gst

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
// extern void _goglib_gst1_LogFunction(GstDebugCategory*, GstDebugLevel, const gchar*, const gchar*, gint, GObject*, GstDebugMessage*, gpointer);
// extern void destroyUserdata(gpointer);
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/userdata"
)

// LogFunctionHandle is a handle returned by [DebugAddLogFunction] that can be used to
// remove the log function later via [DebugRemoveLogFunction].
type LogFunctionHandle unsafe.Pointer

// DebugAddLogFunction wraps gst_debug_add_log_function.
//
// See https://gstreamer.freedesktop.org/documentation/gstreamer/gstinfo.html#gst_debug_add_log_function
//
// returns a [LogFunctionHandle] that can be passed to [DebugRemoveLogFunction] to remove the log
// function later.
func DebugAddLogFunction(fn LogFunction) LogFunctionHandle {
	ptr := userdata.Register(fn)

	C.gst_debug_add_log_function(
		(*[0]byte)(C._goglib_gst1_LogFunction),
		C.gpointer(ptr),
		(C.GDestroyNotify)((*[0]byte)(C.destroyUserdata)),
	)

	return LogFunctionHandle(ptr)
}

// DebugRemoveLogFunction wraps gst_debug_remove_log_function_by_data.
//
// See https://gstreamer.freedesktop.org/documentation/gstreamer/gstinfo.html?gi-language=c#gst_debug_remove_log_function
func DebugRemoveLogFunction(handle LogFunctionHandle) uint {
	if handle == nil {
		return 0
	}
	ret := C.gst_debug_remove_log_function_by_data(C.gpointer(handle))

	return uint(ret)
}

// DebugRemoveDefaultLogFunction removes GStreamer's default log function (gst_debug_log_default).
// Returns the number of functions removed.
func DebugRemoveDefaultLogFunction() uint {
	return uint(C.gst_debug_remove_log_function((*[0]byte)(C.gst_debug_log_default)))
}

// DebugAddDefaultLogFunction re-adds GStreamer's default log function (gst_debug_log_default).
// Note: calling this multiple times will register the default function multiple times.
func DebugAddDefaultLogFunction() {
	C.gst_debug_add_log_function((*[0]byte)(C.gst_debug_log_default), nil, nil)
}
