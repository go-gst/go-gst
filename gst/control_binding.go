package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

type ControlBinding struct{ *Object }

func (cb *ControlBinding) Instance() *C.GstControlBinding {
	return C.toGstControlBinding(cb.Unsafe())
}

type DirectControlBinding struct{ ControlBinding }

func NewDirectControlBinding(obj *Object, prop string, csource *InterpolationControlSource) *DirectControlBinding {
	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))

	cbinding := C.gst_direct_control_binding_new(obj.Instance(), cprop, csource.Instance())

	return &DirectControlBinding{
		ControlBinding: ControlBinding{
			Object: wrapObject(glib.Take(unsafe.Pointer(cbinding))),
		},
	}
}
