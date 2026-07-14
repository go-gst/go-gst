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

// GetValueIndex wraps gst_tag_list_get_value_index
func (list *TagList) GetValueIndex(tag string, index uint) any {
	var carg0 *C.GstTagList // in, none, converted
	var carg1 *C.gchar      // in, none, string, casted *C.gchar
	var carg2 C.guint       // in, none, casted
	var cret *C.GValue      // return, none, converted

	carg0 = (*C.GstTagList)(UnsafeTagListToGlibNone(list))
	carg1 = (*C.gchar)(unsafe.Pointer(C.CString(tag)))
	defer C.free(unsafe.Pointer(carg1))
	carg2 = C.guint(index)

	cret = C.gst_tag_list_get_value_index(carg0, carg1, carg2)
	runtime.KeepAlive(list)
	runtime.KeepAlive(tag)
	runtime.KeepAlive(index)

	var goret any

	goret = gobject.ValueFromNative(unsafe.Pointer(cret)).GoValue()

	return goret
}
