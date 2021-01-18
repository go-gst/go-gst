package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"
)

// SystemClock wraps GstSystemClock
type SystemClock struct{ *Clock }

// ObtainSystemClock returns the default SystemClock. The refcount of the clock will be
// increased so you need to unref the clock after usage.
func ObtainSystemClock() *SystemClock {
	return &SystemClock{FromGstClockUnsafeFull(unsafe.Pointer(C.gst_system_clock_obtain()))}
}
