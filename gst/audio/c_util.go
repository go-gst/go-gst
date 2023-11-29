package audio

/*
#include "gst.go.h"

GstClockTime
framesToClockTime (gint frames, gint rate) { return GST_FRAMES_TO_CLOCK_TIME(frames, rate); }

gint
clockTimeToFrames(GstClockTime ct, gint rate) { return GST_CLOCK_TIME_TO_FRAMES(ct, rate); }

GValue *  audioUtilToGValue (guintptr p) { return (GValue*)(p); }
*/
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-gst/gst"
)

// FramesToClockTime calculates the Clocktime
// from the given frames and rate.
func FramesToClockTime(frames, rate int) gst.ClockTime {
	ct := C.framesToClockTime(C.gint(frames), C.gint(rate))
	return gst.ClockTime(ct)
}

// DurationToFrames calculates the number of frames from the given duration and sample rate.
func DurationToFrames(dur gst.ClockTime, rate int) int {
	return int(C.clockTimeToFrames(C.GstClockTime(dur), C.gint(rate)))
}

// gboolean converts a go bool to a C.gboolean.
func gboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

// gobool provides an easy type conversion between C.gboolean and a go bool.
func gobool(b C.gboolean) bool {
	return int(b) > 0
}

func ptrToGVal(p unsafe.Pointer) *C.GValue {
	return (*C.GValue)(C.audioUtilToGValue(C.guintptr(p)))
}
