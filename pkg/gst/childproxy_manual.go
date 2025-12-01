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

type ChildProxyExtManual interface {
	// GetProperty wraps gst_child_proxy_get_property
	GetProperty(string) any
	// SetProperty wraps gst_child_proxy_set_property
	SetProperty(string, any)
}

// SetProperty wraps gst_child_proxy_set_property
func (object *ChildProxyInstance) SetProperty(name string, value any) {
	var carg0 *C.GstChildProxy // in, none, converted
	var carg1 *C.gchar         // in, none, string, casted *C.gchar
	var carg2 *C.GValue        // in, none, converted

	carg0 = (*C.GstChildProxy)(UnsafeChildProxyToGlibNone(object))
	carg1 = (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.free(unsafe.Pointer(carg1))
	carg2 = (*C.GValue)(gobject.UnsafeValueToGlibNone(gobject.NewValue(value)))

	C.gst_child_proxy_set_property(carg0, carg1, carg2)
	runtime.KeepAlive(object)
	runtime.KeepAlive(name)
	runtime.KeepAlive(value)
}

// GetProperty wraps gst_child_proxy_get_property
func (object *ChildProxyInstance) GetProperty(name string) any {
	var carg0 *C.GstChildProxy // in, none, converted
	var carg1 *C.gchar         // in, none, string, casted *C.gchar
	var carg2 C.GValue         // out, transfer: none, C Pointers: 0, Name: Value, caller-allocates

	carg0 = (*C.GstChildProxy)(UnsafeChildProxyToGlibNone(object))
	carg1 = (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.free(unsafe.Pointer(carg1))

	C.gst_child_proxy_get_property(carg0, carg1, &carg2)
	runtime.KeepAlive(object)
	runtime.KeepAlive(name)

	var value any

	value = gobject.ValueFromNative(unsafe.Pointer(&carg2)).GoValue()

	return value
}
