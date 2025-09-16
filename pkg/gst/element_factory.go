package gst

import (
	"runtime"
	"unsafe"

	"github.com/go-gst/go-glib/pkg/gobject/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// ElementFactoryInstanceMakeWithProperties wraps gst_element_factory_make_with_properties
//
// The function takes the following parameters:
//
//   - factoryname string: a named factory to instantiate
//   - properties map[string]any: a map of properties to set on the element
//
// The function returns the following values:
//
//   - goret Element
//
// Create a new element of the type defined by the given elementfactory.
// The supplied list of properties, will be passed at object construction.
func ElementFactoryMakeWithProperties(factoryname string, properties map[string]any) Element {
	var cname *C.gchar      // out
	var _cret *C.GstElement // in

	cname = (*C.gchar)(unsafe.Pointer(C.CString(factoryname)))
	defer C.free(unsafe.Pointer(cname))

	var cnames **C.gchar
	var cvalues *C.GValue

	if len(properties) > 0 {
		names := make([]*C.char, 0, len(properties))
		values := make([]C.GValue, 0, len(properties))

		for name, value := range properties {
			cname := (*C.char)(C.CString(name))
			defer C.free(unsafe.Pointer(cname))

			gvalue := gobject.NewValue(value)
			defer runtime.KeepAlive(gvalue)

			names = append(names, cname)
			values = append(values, *(*C.GValue)(gobject.UnsafeValueToGlibNone(gvalue)))
		}

		cnames = unsafe.SliceData(names)
		cvalues = unsafe.SliceData(values)
	}

	n_params := C.guint(len(properties))

	_cret = C.gst_element_factory_make_with_properties(cname, n_params, cnames, cvalues)
	runtime.KeepAlive(factoryname)

	var _element Element // out

	if _cret != nil {
		_element = UnsafeElementFromGlibNone(unsafe.Pointer(_cret))
	}

	return _element
}
