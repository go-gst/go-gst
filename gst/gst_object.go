package gst

// #include "gst.go.h"
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Object is a go representation of a GstObject.
type Object struct{ *glib.InitiallyUnowned }

// FromGstObjectUnsafeNone returns an Object wrapping the given pointer. It meant for internal
// usage and exported for visibility to other packages.
func FromGstObjectUnsafeNone(ptr unsafe.Pointer) *Object { return wrapObject(glib.TransferNone(ptr)) }

// FromGstObjectUnsafeFull returns an Object wrapping the given pointer. It meant for internal
// usage and exported for visibility to other packages.
func FromGstObjectUnsafeFull(ptr unsafe.Pointer) *Object { return wrapObject(glib.TransferFull(ptr)) }

// Instance returns the native C GstObject.
func (o *Object) Instance() *C.GstObject { return C.toGstObject(o.Unsafe()) }

// BaseObject is a convenience method for retrieving this object from embedded structs.
func (o *Object) BaseObject() *Object { return o }

// GstObject is an alias to Instance on the underlying GstObject of any extending struct.
func (o *Object) GstObject() *C.GstObject { return C.toGstObject(o.Unsafe()) }

// GObject returns the underlying GObject instance.
func (o *Object) GObject() *glib.Object { return o.InitiallyUnowned.Object }

// GetName returns the name of this object.
func (o *Object) GetName() string {
	cName := C.gst_object_get_name((*C.GstObject)(o.Instance()))
	defer C.free(unsafe.Pointer(cName))
	return C.GoString(cName)
}

// GetValue retrieves the value for the given controlled property at the given timestamp.
func (o *Object) GetValue(property string, timestamp time.Duration) *glib.Value {
	cprop := C.CString(property)
	defer C.free(unsafe.Pointer(cprop))
	gval := C.gst_object_get_value(o.Instance(), (*C.gchar)(cprop), C.GstClockTime(timestamp.Nanoseconds()))
	if gval == nil {
		return nil
	}
	return glib.ValueFromNative(unsafe.Pointer(gval))
}

// SetArg sets the argument name to value on this object. Note that function silently returns
// if object has no property named name or when value cannot be converted to the type for this
// property.
func (o *Object) SetArg(name, value string) {
	cName := C.CString(name)
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cValue))
	C.gst_util_set_object_arg(
		(*C.GObject)(o.Unsafe()),
		(*C.gchar)(unsafe.Pointer(cName)),
		(*C.gchar)(unsafe.Pointer(cValue)),
	)
}

// Log logs a message to the given category from this object using the currently registered
// debugging handlers.
func (o *Object) Log(cat *DebugCategory, level DebugLevel, message string) {
	cat.logDepth(level, message, 2, (*C.GObject)(o.Unsafe()))
}

// Clear will will clear all references to this object. If the reference is already null
// the the function does nothing. Otherwise the reference count is decreased and the pointer
// set to null.
func (o *Object) Clear() {
	if ptr := o.Unsafe(); ptr != nil {
		C.gst_clear_object((**C.GstObject)(unsafe.Pointer(&ptr)))
	}
}

// Ref increments the reference count on object. This function does not take the lock on object
// because it relies on atomic refcounting. For convenience the same object is returned.
func (o *Object) Ref() *Object {
	C.gst_object_ref((C.gpointer)(o.Unsafe()))
	return o
}

// Unref decrements the reference count on object. If reference count hits zero, destroy object.
// This function does not take the lock on object as it relies on atomic refcounting.
func (o *Object) Unref() {
	C.gst_object_unref((C.gpointer)(o.Unsafe()))
}
