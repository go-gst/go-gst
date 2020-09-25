package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"

static gboolean
isTypeCaps(GParamSpec * p)
{
	return p->value_type == GST_TYPE_CAPS;
}
*/
import "C"

import (
	"strings"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// ParameterSpec is a go representation of a C GParamSpec
type ParameterSpec struct {
	paramSpec    *C.GParamSpec
	Name         string
	Blurb        string
	Flags        ParameterFlags
	ValueType    glib.Type
	OwnerType    glib.Type
	DefaultValue *glib.Value
}

// Unref the underlying paramater spec.
func (p *ParameterSpec) Unref() { C.g_param_spec_unref(p.paramSpec) }

// UIntRange returns the range of the Uint stored in this parameter spec.
func (p *ParameterSpec) UIntRange() (uint, uint) {
	paramUint := C.getParamUInt(p.paramSpec)
	return uint(paramUint.minimum), uint(paramUint.maximum)
}

// IntRange returns the range of the Int stored in this parameter spec.
func (p *ParameterSpec) IntRange() (int, int) {
	paramUint := C.getParamInt(p.paramSpec)
	return int(paramUint.minimum), int(paramUint.maximum)
}

// UInt64Range returns the range of the Uint64 stored in this parameter spec.
func (p *ParameterSpec) UInt64Range() (uint64, uint64) {
	paramUint := C.getParamUInt64(p.paramSpec)
	return uint64(paramUint.minimum), uint64(paramUint.maximum)
}

// Int64Range returns the range of the Int64 stored in this parameter spec.
func (p *ParameterSpec) Int64Range() (int64, int64) {
	paramUint := C.getParamInt64(p.paramSpec)
	return int64(paramUint.minimum), int64(paramUint.maximum)
}

// FloatRange returns the range of the Float stored in this parameter spec.
func (p *ParameterSpec) FloatRange() (float64, float64) {
	paramUint := C.getParamFloat(p.paramSpec)
	return float64(paramUint.minimum), float64(paramUint.maximum)
}

// DoubleRange returns the range of the Double stored in this parameter spec.
func (p *ParameterSpec) DoubleRange() (float64, float64) {
	paramUint := C.getParamDouble(p.paramSpec)
	return float64(paramUint.minimum), float64(paramUint.maximum)
}

// IsCaps returns true if this parameter contains a caps object.
func (p *ParameterSpec) IsCaps() bool {
	return gobool(C.isTypeCaps(p.paramSpec))
}

// GetCaps returns the caps in this parameter if it is of type GST_TYPE_CAPS.
func (p *ParameterSpec) GetCaps() Caps {
	caps := C.gst_value_get_caps((*C.GValue)(unsafe.Pointer(p.DefaultValue.Native())))
	if caps == nil {
		return nil
	}
	return FromGstCaps(caps)
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
	"deprecated",
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
	if p.Has(ParameterDeprecated) {
		out = append(out, "deprecated")
	}
	return strings.Join(out, ", ")
}
