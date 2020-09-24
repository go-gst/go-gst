package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"

import (
	"strings"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Object is a go representation of a GstObject. Type casting stops here
// and we do not descend into the glib library.
type Object struct{ *glib.InitiallyUnowned }

// native returns the pointer to the underlying object.
func (o *Object) unsafe() unsafe.Pointer { return unsafe.Pointer(o.InitiallyUnowned.Native()) }

// Instance returns the native C GstObject.
func (o *Object) Instance() *C.GstObject { return C.toGstObject(o.unsafe()) }

// Class returns the GObjectClass of this instance.
func (o *Object) Class() *C.GObjectClass { return C.getGObjectClass(o.unsafe()) }

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
// set in this object.
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
			C.g_object_get_property((*C.GObject)(o.unsafe()), prop.name, &gval)
		} else {
			C.g_param_value_set_default((*C.GParamSpec)(prop), &gval)
		}
		out = append(out, &ParameterSpec{
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

func wrapObject(o *C.GstObject) *Object {
	obj := &Object{&glib.InitiallyUnowned{Object: glib.Take(unsafe.Pointer(o))}}
	obj.RefSink()
	return obj
}

// ParameterSpec is a go representation of a C GParamSpec
type ParameterSpec struct {
	param        *C.GParamSpec
	Name         string
	Blurb        string
	Flags        ParameterFlags
	ValueType    glib.Type
	OwnerType    glib.Type
	DefaultValue *glib.Value
}

// ParameterFlags is a go cast of GParamFlags.
type ParameterFlags C.GParamFlags

// Has returns true if these flags contain the provided ones.
func (p ParameterFlags) Has(b ParameterFlags) bool { return p&b != 0 }

// Type casting of GParamFlags
const (
	ParameterReadable       ParameterFlags = C.G_PARAM_READABLE        // the parameter is readable
	ParameterWritable                      = C.G_PARAM_WRITABLE        // the parameter is writable
	ParameterConstruct                     = C.G_PARAM_CONSTRUCT       // the parameter will be set upon object construction
	ParameterConstructOnly                 = C.G_PARAM_CONSTRUCT_ONLY  // the parameter can only be set upon object construction
	ParameterLaxValidation                 = C.G_PARAM_LAX_VALIDATION  // upon parameter conversion (see g_param_value_convert()) strict validation is not required
	ParameterStaticName                    = C.G_PARAM_STATIC_NAME     // the string used as name when constructing the parameter is guaranteed to remain valid and unmodified for the lifetime of the parameter. Since 2.8
	ParameterStaticNick                    = C.G_PARAM_STATIC_NICK     // the string used as nick when constructing the parameter is guaranteed to remain valid and unmmodified for the lifetime of the parameter. Since 2.8
	ParameterStaticBlurb                   = C.G_PARAM_STATIC_BLURB    // the string used as blurb when constructing the parameter is guaranteed to remain valid and unmodified for the lifetime of the parameter. Since 2.8
	ParameterExplicitNotify                = C.G_PARAM_EXPLICIT_NOTIFY // calls to g_object_set_property() for this property will not automatically result in a "notify" signal being emitted: the implementation must call g_object_notify() themselves in case the property actually changes. Since: 2.42.
	ParameterDeprecated                    = C.G_PARAM_DEPRECATED      // the parameter is deprecated and will be removed in a future version. A warning will be generated if it is used while running with G_ENABLE_DIAGNOSTIC=1. Since 2.26
	ParameterControllable                  = C.GST_PARAM_CONTROLLABLE
	ParameterMutablePlaying                = C.GST_PARAM_MUTABLE_PLAYING
	ParameterMutablePaused                 = C.GST_PARAM_MUTABLE_PAUSED
	ParameterMutableReady                  = C.GST_PARAM_MUTABLE_READY
)

var allFlags = []ParameterFlags{
	ParameterReadable,
	ParameterWritable,
	ParameterConstruct,
	ParameterConstructOnly,
	ParameterLaxValidation,
	ParameterStaticName,
	ParameterStaticNick,
	ParameterStaticBlurb,
	ParameterExplicitNotify,
	ParameterDeprecated,
	ParameterControllable,
	ParameterMutablePlaying,
	ParameterMutablePaused,
	ParameterMutableReady,
}

var allFlagStrings = []string{
	"readable",
	"writable",
	"construct",
	"construct only",
	"lax validation",
	"static name",
	"static nick",
	"static blurb",
	"explicity notify",
	"deprecated",
	"controllable",
	"changeable in NULL, READY, PAUSED or PLAYING state",
	"changeable only in NULL, READY or PAUSED state",
	"changeable only in NULL or READY state",
}

func (p ParameterFlags) String() string {
	out := make([]string, 0)
	for idx, param := range allFlags {
		if p.Has(param) {
			out = append(out, allFlagStrings[idx])
		}
	}
	return strings.Join(out, ", ")
}

// GstFlagsString returns a string of the flags that are relevant specifically
// to gstreamer.
func (p ParameterFlags) GstFlagsString() string {
	out := make([]string, 0)
	if p.Has(ParameterReadable) {
		out = append(out, "readable")
	}
	if p.Has(ParameterWritable) {
		out = append(out, "writable")
	}
	if p.Has(ParameterControllable) {
		out = append(out, "controllable")
	}
	if p.Has(ParameterMutablePlaying) {
		out = append(out, "changeable in NULL, READY, PAUSED or PLAYING state")
	}
	if p.Has(ParameterMutablePaused) {
		out = append(out, "changeable only in NULL, READY or PAUSED state")
	}
	if p.Has(ParameterMutableReady) {
		out = append(out, "changeable only in NULL or READY state")
	}
	return strings.Join(out, ", ")
}
