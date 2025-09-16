package gst

import (
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/core/userdata"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

//export _gotk4_gst1_BusSyncHandler
func _gotk4_gst1_BusSyncHandler(carg1 *C.GstBus, carg2 *C.GstMessage, carg3 C.gpointer) (cret C.GstBusSyncReply) {
	var fn BusSyncHandler
	{
		v := userdata.Load(unsafe.Pointer(carg3))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(BusSyncHandler)
	}

	bus := UnsafeBusFromGlibNone(unsafe.Pointer(carg1))
	msg := UnsafeMessageFromGlibBorrow(unsafe.Pointer(carg2))

	// if the user returns BusDrop then we must free the message. We pass a copy to the user instead of the original
	// so that they can keep it alive if they want to.
	goret := fn(bus, msg.Copy())

	if goret == BusDrop {
		UnsafeMessageUnref(msg)
	}

	return C.GstBusSyncReply(goret)
}
