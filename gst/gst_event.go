package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"
)

// Event is a go wrapper around a GstEvent.
type Event struct {
	ptr *C.GstEvent
}

// Instance returns the underlying GstEvent instance.
func (e *Event) Instance() *C.GstEvent { return C.toGstEvent(unsafe.Pointer(e.ptr)) }
