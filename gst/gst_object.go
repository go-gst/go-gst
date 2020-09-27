package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Object is a go representation of a GstObject. Type casting stops here
// and we do not descend into the glib library.
type Object struct{ *glib.InitiallyUnowned }

// Unsafe returns the unsafe pointer to the underlying object. This method is primarily
// for internal usage and is exposed for visibility in other packages.
func (o *Object) Unsafe() unsafe.Pointer {
	if o == nil || o.GObject == nil {
		return nil
	}
	return unsafe.Pointer(o.GObject)
}

// Instance returns the native C GstObject.
func (o *Object) Instance() *C.GstObject { return C.toGstObject(o.Unsafe()) }

// Unref wraps the underlying Unref from glib, and performs an extra check that the
// object has not already been destroyed.
func (o *Object) Unref() {
	if o.GObject == nil {
		return
	}
	o.Object.Unref()
}

// Class returns the GObjectClass of this instance.
func (o *Object) Class() *C.GObjectClass { return C.getGObjectClass(o.Unsafe()) }

// Name returns the name of this object.
func (o *Object) Name() string {
	cName := C.gst_object_get_name((*C.GstObject)(o.Instance()))
	defer C.free(unsafe.Pointer(cName))
	return C.GoString(cName)
}

// Interfaces returns the interfaces associated with this object.
func (o *Object) Interfaces() []string {
	var size C.guint
	ifaces := C.g_type_interfaces(C.ulong(o.TypeFromInstance()), &size)
	if int(size) == 0 {
		return nil
	}
	defer C.g_free((C.gpointer)(ifaces))
	out := make([]string, int(size))
	for _, t := range (*[1 << 30]int)(unsafe.Pointer(ifaces))[:size:size] {
		out = append(out, glib.Type(t).Name())
	}
	return out
}

// ListProperties returns a list of the properties associated with this object.
// The default values assumed in the parameter spec reflect the values currently
// set in this object, or their defaults.
//
// Unref after usage.
func (o *Object) ListProperties() []*ParameterSpec {
	var size C.guint
	props := C.g_object_class_list_properties((*C.GObjectClass)(o.Class()), &size)
	if props == nil {
		return nil
	}
	defer C.g_free((C.gpointer)(props))
	out := make([]*ParameterSpec, 0)
	for _, prop := range (*[1 << 30]*C.GParamSpec)(unsafe.Pointer(props))[:size:size] {
		var gval C.GValue
		flags := ParameterFlags(prop.flags)
		if flags.Has(ParameterReadable) {
			C.g_object_get_property((*C.GObject)(o.Unsafe()), prop.name, &gval)
		} else {
			C.g_param_value_set_default((*C.GParamSpec)(prop), &gval)
		}
		C.g_param_spec_sink(prop) // steal the ref on the property
		out = append(out, &ParameterSpec{
			paramSpec:    prop,
			Name:         C.GoString(C.g_param_spec_get_name(prop)),
			Blurb:        C.GoString(C.g_param_spec_get_blurb(prop)),
			Flags:        flags,
			ValueType:    glib.Type(prop.value_type),
			OwnerType:    glib.Type(prop.owner_type),
			DefaultValue: glib.ValueFromNative(unsafe.Pointer(&gval)),
		})
	}
	return out
}
