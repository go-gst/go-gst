package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// SystemClock wraps GstSystemClock
type SystemClock struct{ *Clock }

// ObtainSystemClock returns the default SystemClock. The refcount of the clock will be
// increased so you need to unref the clock after usage.
func ObtainSystemClock() *SystemClock {
	clock := C.gst_system_clock_obtain()
	return &SystemClock{wrapClock(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(clock))})}
}
