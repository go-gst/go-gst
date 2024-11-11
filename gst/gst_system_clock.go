package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

// SystemClock wraps GstSystemClock
type SystemClock struct{ *Clock }

var TYPE_SYSTEM_CLOCK = glib.Type(C.GST_TYPE_SYSTEM_CLOCK)

// ClockType represents GstClockType
type ClockType int

const (
	//time since Epoch
	ClockTypeRealtime = C.GST_CLOCK_TYPE_REALTIME
	//monotonic time since some unspecified starting point
	ClockTypeMonotonic = C.GST_CLOCK_TYPE_MONOTONIC
	// some other time source is used (Since: 1.0.5)
	ClockTypeOther = C.GST_CLOCK_TYPE_OTHER
	// time since Epoch, but using International Atomic Time as reference (Since: 1.18)
	ClockTypeTAI = C.GST_CLOCK_TYPE_TAI
)

// ObtainSystemClock returns the default SystemClock.
func ObtainSystemClock() *SystemClock {
	return &SystemClock{FromGstClockUnsafeFull(unsafe.Pointer(C.gst_system_clock_obtain()))}
}

// NewSystemClock creates a new instance of a SystemClock, with the given clock type parameter
//
// This is only a convenience wrapper for glib.NewObjectWithProperties
func NewSystemClock(clockType ClockType) (*SystemClock, error) {
	clockObj, err := glib.NewObjectWithProperties(TYPE_SYSTEM_CLOCK, map[string]any{
		"clock-type": clockType,
	})

	if err != nil {
		return nil, err
	}

	return &SystemClock{ToGstClock(clockObj)}, nil
}
