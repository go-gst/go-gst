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

type TagSetterExtManual interface {
	// AddTagValue wraps gst_tag_setter_add_tag_value
	//
	// The function takes the following parameters:
	//
	// 	- mode TagMergeMode: the mode to use
	// 	- tag string: tag to set
	// 	- value any: GValue to set for the tag
	//
	// Adds the given tag / GValue pair on the setter using the given merge mode.
	AddTagValue(mode TagMergeMode, tag string, value any)
}

// AddTagValue wraps gst_tag_setter_add_tag_value
//
// The function takes the following parameters:
//
//   - mode TagMergeMode: the mode to use
//   - tag string: tag to set
//   - value any: GValue to set for the tag
//
// Adds the given tag / GValue pair on the setter using the given merge mode.
func (setter *TagSetterInstance) AddTagValue(mode TagMergeMode, tag string, value any) {
	var carg0 *C.GstTagSetter   // in, none, converted
	var carg1 C.GstTagMergeMode // in, none, casted
	var carg2 *C.gchar          // in, none, string
	var carg3 *C.GValue         // in, none, converted

	carg0 = (*C.GstTagSetter)(UnsafeTagSetterToGlibNone(setter))
	carg1 = C.GstTagMergeMode(mode)
	carg2 = (*C.gchar)(unsafe.Pointer(C.CString(tag)))
	defer C.free(unsafe.Pointer(carg2))
	carg3 = (*C.GValue)(gobject.UnsafeValueToGlibNone(gobject.NewValue(value)))

	C.gst_tag_setter_add_tag_value(carg0, carg1, carg2, carg3)
	runtime.KeepAlive(setter)
	runtime.KeepAlive(mode)
	runtime.KeepAlive(tag)
	runtime.KeepAlive(value)
}
