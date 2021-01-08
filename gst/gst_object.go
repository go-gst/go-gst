package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Object is a go representation of a GstObject.
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

// BaseObject returns this object for embedding structs.
func (o *Object) BaseObject() *Object { return o }

// GstObject is an alias to Instance on the underlying GstObject of any extending struct.
func (o *Object) GstObject() *C.GstObject { return C.toGstObject(o.Unsafe()) }

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
	ifaces := C.g_type_interfaces(C.gsize(o.TypeFromInstance()), &size)
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
		C.g_param_spec_sink(prop) // steal the ref on the property
		out = append(out, &ParameterSpec{
			paramSpec: prop,
		})
	}
	return out
}

// Log logs a message to the given category from this object using the currently registered
// debugging handlers.
func (o *Object) Log(cat *DebugCategory, level DebugLevel, message string) {
	cat.logDepth(level, message, 2, (*C.GObject)(o.Unsafe()))
}
