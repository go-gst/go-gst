package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include <glib-object.h>
#include "gst.go.h"
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Caps is a wrapper around GstCaps. It provides a function for easy type
// conversion.
type Caps []*Structure

// NewRawCaps returns new GstCaps with the given format, sample-rate, and channels.
func NewRawCaps(format string, rate, channels int) Caps {
	return Caps{
		{
			Name: "audio/x-raw",
			Data: map[string]interface{}{
				"format":   format,
				"rate":     rate,
				"channels": channels,
			},
		},
	}
}

// FromGstCaps converts a C GstCaps objects to a go type.
func FromGstCaps(caps *C.GstCaps) Caps {
	out := make(Caps, 0)
	size := int(C.gst_caps_get_size((*C.GstCaps)(caps)))
	for i := 0; i < size-1; i++ {
		s := C.gst_caps_get_structure((*C.GstCaps)(caps), (C.guint(i)))
		out = append(out, FromGstStructure(s))
	}
	return out
}

// ToGstCaps returns the GstCaps representation of this Caps instance.
func (g Caps) ToGstCaps() *C.GstCaps {
	// create a new empty caps object
	caps := C.gst_caps_new_empty()
	if caps == nil {
		// extra nil check but this would only happen when larger issues are present
		return nil
	}
	for _, st := range g {
		// append the structure to the caps
		C.gst_caps_append_structure((*C.GstCaps)(caps), (*C.GstStructure)(st.ToGstStructure()))
	}
	return caps
}

// Structure is a go implementation of a C GstStructure.
type Structure struct {
	Name string
	Data map[string]interface{}
}

// ToGstStructure converts this structure to a C GstStructure.
func (s *Structure) ToGstStructure() *C.GstStructure {
	var structStr string
	structStr = s.Name
	// build a structure string from the data
	if s.Data != nil {
		elems := make([]string, 0)
		for k, v := range s.Data {
			elems = append(elems, fmt.Sprintf("%s=%v", k, v))
		}
		structStr = fmt.Sprintf("%s, %s", s.Name, strings.Join(elems, ", "))
	}
	// convert the structure string to a cstring
	cstr := C.CString(structStr)
	defer C.free(unsafe.Pointer(cstr))
	// a small buffer for garbage
	p := C.malloc(C.size_t(128))
	defer C.free(p)
	// create a structure from the string
	cstruct := C.gst_structure_from_string((*C.gchar)(cstr), (**C.gchar)(p))
	return cstruct
}

// FromGstStructure converts the given C GstStructure into a go structure.
func FromGstStructure(s *C.GstStructure) *Structure {
	v := &Structure{}
	v.Name = C.GoString((*C.char)(C.gst_structure_get_name((*C.GstStructure)(s))))
	n := uint(C.gst_structure_n_fields(s))
	v.Data = make(map[string]interface{})
	for i := uint(0); i < n; i++ {
		fn := C.gst_structure_nth_field_name((*C.GstStructure)(s), C.guint(i))
		fv := glib.ValueFromNative(unsafe.Pointer(C.gst_structure_id_get_value((*C.GstStructure)(s), C.g_quark_from_string(fn))))
		val, _ := fv.GoValue()
		v.Data[C.GoString((*C.char)(fn))] = val
	}
	return v
}
