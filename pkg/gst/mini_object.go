package gst

import (
	"unsafe"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// miniObjectRef calls gst_mini_object_ref on the given pointer. It is used for MiniObject extending
// records, because the generator does not have the hierarchy to be able to figure out the casting.
//
// ptr must be a valid C pointer to a GstMiniObject or extending record.
func miniObjectRef(ptr unsafe.Pointer) {
	C.gst_mini_object_ref((*C.GstMiniObject)(ptr))
}

// miniObjectUnref calls gst_mini_object_unref on the given pointer. It is used for MiniObject extending
// records, because the generator does not have the hierarchy to be able to figure out the casting.
//
// ptr must be a valid C pointer to a GstMiniObject or extending record.
func miniObjectUnref(ptr unsafe.Pointer) {
	C.gst_mini_object_unref((*C.GstMiniObject)(ptr))
}
