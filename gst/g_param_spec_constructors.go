package gst

// #include "gst.go.h"
import "C"

import "github.com/tinyzimmer/go-glib/glib"

// NewStringParam returns a new ParamSpec that will hold a string value.
func NewStringParam(name, nick, blurb string, defaultValue *string, flags ParameterFlags) *ParamSpec {
	var cdefault *C.gchar
	if defaultValue != nil {
		cdefault = C.CString(*defaultValue)
	}
	paramSpec := C.g_param_spec_string(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		(*C.gchar)(cdefault),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// NewBoolParam creates a new ParamSpec that will hold a boolean value.
func NewBoolParam(name, nick, blurb string, defaultValue bool, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_boolean(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		gboolean(defaultValue),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// NewIntParam creates a new ParamSpec that will hold a signed integer value.
func NewIntParam(name, nick, blurb string, min, max, defaultValue int, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_int(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		C.gint(min),
		C.gint(max),
		C.gint(defaultValue),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// NewUintParam creates a new ParamSpec that will hold an unsigned integer value.
func NewUintParam(name, nick, blurb string, min, max, defaultValue uint, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_uint(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		C.guint(min),
		C.guint(max),
		C.guint(defaultValue),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// NewInt64Param creates a new ParamSpec that will hold a signed 64-bit integer value.
func NewInt64Param(name, nick, blurb string, min, max, defaultValue int64, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_int64(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		C.gint64(min),
		C.gint64(max),
		C.gint64(defaultValue),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// NewUint64Param creates a new ParamSpec that will hold an unsigned 64-bit integer value.
func NewUint64Param(name, nick, blurb string, min, max, defaultValue uint64, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_uint64(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		C.guint64(min),
		C.guint64(max),
		C.guint64(defaultValue),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// NewFloat32Param creates a new ParamSpec that will hold a 32-bit float value.
func NewFloat32Param(name, nick, blurb string, min, max, defaultValue float32, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_float(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		C.gfloat(min),
		C.gfloat(max),
		C.gfloat(defaultValue),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// NewFloat64Param creates a new ParamSpec that will hold a 64-bit float value.
func NewFloat64Param(name, nick, blurb string, min, max, defaultValue float64, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_double(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		C.gdouble(min),
		C.gdouble(max),
		C.gdouble(defaultValue),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}

// TypeCaps is the static Glib Type for a GstCaps.
var TypeCaps = glib.Type(C.gst_caps_get_type())

// NewBoxedParam creates a new ParamSpec containing a boxed type. Some helper type castings are included
// in these bindings.
func NewBoxedParam(name, nick, blurb string, boxedType glib.Type, flags ParameterFlags) *ParamSpec {
	paramSpec := C.g_param_spec_boxed(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		C.GType(boxedType),
		C.GParamFlags(flags),
	)
	return &ParamSpec{paramSpec: paramSpec}
}
