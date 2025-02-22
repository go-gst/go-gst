package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
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
	//cName could be NULL which needs to be freed using
	defer C.g_free((C.gpointer)(unsafe.Pointer(cName)))
	return C.GoString(cName)
}

// GetParent retrieves the parent of this object.
func (o *Object) GetParent() *Object {
	return wrapObject(glib.Take(unsafe.Pointer(C.gst_object_get_parent(o.Instance()))))
}

// GetValue retrieves the value for the given controlled property at the given timestamp.
func (o *Object) GetValue(property string, timestamp ClockTime) *glib.Value {
	cprop := C.CString(property)
	defer C.free(unsafe.Pointer(cprop))
	gval := C.gst_object_get_value(o.Instance(), (*C.gchar)(cprop), C.GstClockTime(timestamp))
	if gval == nil {
		return nil
	}
	return glib.ValueFromNative(unsafe.Pointer(gval))
}

// SetName sets the name of this object.
func (o *Object) SetName(name string) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return gobool(C.gst_object_set_name((*C.GstObject)(o.Instance()), (*C.gchar)(cName)))
}

// SetParent sets the parent of this object.
func (o *Object) SetParent(parent *Object) bool {
	return gobool(C.gst_object_set_parent(o.Instance(), parent.Instance()))
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

func (o *Object) AddControlBinding(binding *ControlBinding) {
	C.gst_object_add_control_binding(o.Instance(), binding.Instance())
}

func (o *Object) RemoveControlBinding(binding *ControlBinding) {
	C.gst_object_remove_control_binding(o.Instance(), binding.Instance())
}

// TODO: Consider wrapping GstObject GST_OBJECT_LOCK/GST_OBJECT_UNLOCK functionality
// due to following flags related functionality is based on a regular uint32 field
// and is not considered thread safe

// Has returns true if this GstObject has the given flags.
func (o *Object) hasFlags(flags uint32) bool {
	return gobool(C.gstObjectFlagIsSet(o.Instance(), C.guint32(flags)))
}

// SetFlags sets the flags
func (o *Object) setFlags(flags uint32) {
	C.gstObjectFlagSet(o.Instance(), C.guint32(flags))
}

// SetFlags unsets the flags
func (o *Object) unsetFlags(flags uint32) {
	C.gstObjectFlagUnset(o.Instance(), C.guint32(flags))
}
