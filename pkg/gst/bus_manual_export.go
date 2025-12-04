package gst

import (
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/userdata"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

//export _gogst_gst1_BusSyncHandler
func _gogst_gst1_BusSyncHandler(carg1 *C.GstBus, carg2 *C.GstMessage, carg3 C.gpointer) (cret C.GstBusSyncReply) {
	var fn BusSyncHandler
	{
		v := userdata.Load(unsafe.Pointer(carg3))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(BusSyncHandler)
	}

	bus := UnsafeBusFromGlibNone(unsafe.Pointer(carg1))
	msg := UnsafeMessageFromGlibNone(unsafe.Pointer(carg2))

	goret := fn(bus, msg)

	// if the user returns BusDrop then we must unref the message an additional time.
	//
	// the finalizer on message will unref it once more
	if goret == BusDrop {
		miniObjectUnref(unsafe.Pointer(msg.native))
	}

	return C.GstBusSyncReply(goret)
}
