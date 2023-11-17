package gst

// #include "gst.go.h"
import "C"
import (
	"strconv"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

type InterpolationControlSource struct{ *Object }

type InterpolationMode int

const (
	//steps-like interpolation, default
	InterpolationModeNone InterpolationMode = 0 // GST_INTERPOLATION_MODE_NONE
	//linear interpolation
	InterpolationModeLinear InterpolationMode = 1 // GST_INTERPOLATION_MODE_LINEAR
	//cubic interpolation (natural), may overshoot the min or max values set by the control point, but is more 'curvy'
	InterpolationModeCubic InterpolationMode = 2 // GST_INTERPOLATION_MODE_CUBIC
	//monotonic cubic interpolation, will not produce any values outside of the min-max range set by the control points (Since: 1.8)
	InterpolationModeCubicMonotonic InterpolationMode = 3 // GST_INTERPOLATION_MODE_CUBIC_MONOTONIC
)

func (cs *InterpolationControlSource) Instance() *C.GstControlSource {
	return C.toGstControlSource(cs.Unsafe())
}

func NewInterpolationControlSource() *InterpolationControlSource {
	cCs := C.gst_interpolation_control_source_new()

	return &InterpolationControlSource{
		Object: wrapObject(glib.TransferFull(unsafe.Pointer(cCs))),
	}
}

func (cs *InterpolationControlSource) SetInterpolationMode(mode InterpolationMode) {
	cs.SetArg("mode", strconv.Itoa(int(mode)))
}

func (cs *InterpolationControlSource) SetTimedValue(time ClockTime, value float64) bool {
	return gobool(C.gst_timed_value_control_source_set(C.toGstTimedValueControlSource(cs.Unsafe()), C.GstClockTime(time), C.double(value)))
}

func (cs *InterpolationControlSource) UnsetAll() {
	C.gst_timed_value_control_source_unset_all(C.toGstTimedValueControlSource(cs.Unsafe()))
}
