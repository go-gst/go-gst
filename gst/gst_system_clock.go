package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"
)

// SystemClock wraps GstSystemClock
type SystemClock struct{ *Clock }

// ObtainSystemClock returns the default SystemClock.
func ObtainSystemClock() *SystemClock {
	return &SystemClock{FromGstClockUnsafeFull(unsafe.Pointer(C.gst_system_clock_obtain()))}
}
