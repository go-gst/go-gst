package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"
)

// MiniObject is an opaque struct meant to form the base of gstreamer
// classes extending the GstMiniObject.
type MiniObject struct {
	ptr unsafe.Pointer
}

// Instance returns the native GstMiniObject instance.
func (m *MiniObject) Instance() *C.GstMiniObject { return C.toGstMiniObject(m.ptr) }

// Ref increases the ref count on this object by one.
func (m *MiniObject) Ref() {}

// Unref decresaes the ref count on this object by one.
func (m *MiniObject) Unref() {}

func wrapMiniObject(p *C.GstMiniObject) *MiniObject {
	return &MiniObject{ptr: unsafe.Pointer(p)}
}
