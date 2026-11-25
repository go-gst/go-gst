package gst

import (
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/userdata"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// GstPadEventFunction has no user_data parameter: the user data given to
// gst_pad_set_event_function_full is stored in the pad itself and read back
// via GST_PAD_EVENTDATA.
//
//export _gogst_gst1_PadEventFunction
func _gogst_gst1_PadEventFunction(carg1 *C.GstPad, carg2 *C.GstObject, carg3 *C.GstEvent) (cret C.gboolean) {
	var fn PadEventFunction
	{
		v := userdata.Load(unsafe.Pointer(carg1.eventdata))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(PadEventFunction)
	}

	var pad Pad       // in, none, converted
	var parent Object // in, none, converted, nullable
	var event *Event  // in, full, converted

	pad = UnsafePadFromGlibNone(unsafe.Pointer(carg1))
	if carg2 != nil {
		parent = UnsafeObjectFromGlibNone(unsafe.Pointer(carg2))
	}
	event = UnsafeEventFromGlibFull(unsafe.Pointer(carg3))

	if fn(pad, parent, event) {
		cret = C.TRUE
	}

	return cret
}
