package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

type InterpolationControlSource struct{ *Object }

func (cs *InterpolationControlSource) Instance() *C.GstControlSource {
	return C.toGstControlSource(cs.Unsafe())
}

func NewInterpolationControlSource() *InterpolationControlSource {
	cCs := C.gst_interpolation_control_source_new()

	return &InterpolationControlSource{
		Object: wrapObject(glib.TransferNone(unsafe.Pointer(cCs))),
	}
}

func (cs *InterpolationControlSource) SetTimedValue(time ClockTime, value float64) {
	C.gst_timed_value_control_source_set(C.toGstTimedValueControlSource(cs.Unsafe()), C.ulong(time), C.double(value))
}
