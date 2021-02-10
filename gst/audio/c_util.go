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
import "time"

// FramesToDuration calculates the Clocktime (which is usually referred to as a time.Duration in the bindings)
// from the given frames and rate.
func FramesToDuration(frames, rate int) time.Duration {
	ct := C.framesToClockTime(C.gint(frames), C.gint(rate))
	return time.Duration(ct)
}

// DurationToFrames calculates the number of frames from the given duration and sample rate.
func DurationToFrames(dur time.Duration, rate int) int {
	return int(C.clockTimeToFrames(C.GstClockTime(dur.Nanoseconds()), C.gint(rate)))
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

func uintptrToGVal(p uintptr) *C.GValue { return (*C.GValue)(C.audioUtilToGValue(C.guintptr(p))) }
