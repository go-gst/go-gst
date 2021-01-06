package gst

// #include "gst.go.h"
import "C"

import (
	"strings"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// ParameterSpec is a go representation of a C GParamSpec
type ParameterSpec struct {
	paramSpec    *C.GParamSpec
	defaultValue *glib.Value
}

// NewStringParameter returns a new ParameterSpec that will hold a string value.
func NewStringParameter(name, nick, blurb string, defaultValue *string, flags ParameterFlags) *ParameterSpec {
	var cdefault *C.gchar
	var paramDefault *glib.Value
	if defaultValue != nil {
		cdefault = C.CString(*defaultValue)
		var err error
		paramDefault, err = glib.ValueInit(glib.TYPE_STRING)
		if err != nil {
			return nil
		}
		paramDefault.SetString(*defaultValue)
	}
	paramSpec := C.g_param_spec_string(
		(*C.gchar)(C.CString(name)),
		(*C.gchar)(C.CString(nick)),
		(*C.gchar)(C.CString(blurb)),
		(*C.gchar)(cdefault),
		C.GParamFlags(flags),
	)
	return &ParameterSpec{paramSpec: paramSpec, defaultValue: paramDefault}
}

// Name returns the name of this parameter.
func (p *ParameterSpec) Name() string {
	return C.GoString(C.g_param_spec_get_name(p.paramSpec))
}

// Blurb returns the blurb for this parameter.
func (p *ParameterSpec) Blurb() string {
	return C.GoString(C.g_param_spec_get_blurb(p.paramSpec))
}

// Flags returns the flags for this parameter.
func (p *ParameterSpec) Flags() ParameterFlags {
	return ParameterFlags(p.paramSpec.flags)
}

// ValueType returns the GType for the value inside this parameter.
func (p *ParameterSpec) ValueType() glib.Type {
	return glib.Type(p.paramSpec.value_type)
}

// OwnerType returns the Gtype for the owner of this parameter.
func (p *ParameterSpec) OwnerType() glib.Type {
	return glib.Type(p.paramSpec.owner_type)
}

// DefaultValue returns the default value for the parameter if it was included when the object
// was instantiated. Otherwise it returns nil.
func (p *ParameterSpec) DefaultValue() *glib.Value {
	return p.defaultValue
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
func (p *ParameterSpec) IsCaps() bool { return gobool(C.isParamSpecTypeCaps(p.paramSpec)) }

// IsEnum returns true if this parameter contains an enum.
func (p *ParameterSpec) IsEnum() bool { return gobool(C.isParamSpecEnum(p.paramSpec)) }

// IsFlags returns true if this paramater contains flags.
func (p *ParameterSpec) IsFlags() bool { return gobool(C.isParamSpecFlags(p.paramSpec)) }

// IsObject returns true if this parameter contains an object.
func (p *ParameterSpec) IsObject() bool { return gobool(C.isParamSpecObject(p.paramSpec)) }

// IsBoxed returns true if this parameter contains a boxed object.
func (p *ParameterSpec) IsBoxed() bool { return gobool(C.isParamSpecBoxed(p.paramSpec)) }

// IsPointer returns true if this paramater contains a pointer.
func (p *ParameterSpec) IsPointer() bool { return gobool(C.isParamSpecPointer(p.paramSpec)) }

// IsFraction returns true if this parameter contains a fraction.
func (p *ParameterSpec) IsFraction() bool { return gobool(C.isParamSpecFraction(p.paramSpec)) }

// IsGstArray returns true if this parameter contains a Gst array.
func (p *ParameterSpec) IsGstArray() bool { return gobool(C.isParamSpecGstArray(p.paramSpec)) }

// EnumValue is a go representation of a GEnumValue
type EnumValue struct {
	Value                int
	ValueNick, ValueName string
}

// GetEnumValues returns the possible enum values for this parameter.
func (p *ParameterSpec) GetEnumValues() []*EnumValue {
	var gsize C.guint
	gEnumValues := C.getEnumValues(p.paramSpec, &gsize)
	size := int(gsize)
	out := make([]*EnumValue, size)
	for idx, enumVal := range (*[1 << 30]C.GEnumValue)(unsafe.Pointer(gEnumValues))[:size:size] {
		out[idx] = &EnumValue{
			Value:     int(enumVal.value),
			ValueNick: C.GoString(enumVal.value_nick),
			ValueName: C.GoString(enumVal.value_name),
		}
	}
	return out
}

// FlagsValue is a go representation of GFlagsValue
type FlagsValue struct {
	Value                int
	ValueName, ValueNick string
}

// GetDefaultFlags returns the default flags for this parameter spec.
func (p *ParameterSpec) GetDefaultFlags() int {
	if p.DefaultValue() == nil {
		return 0
	}
	return int(C.g_value_get_flags((*C.GValue)(p.DefaultValue().Native())))
}

// GetFlagValues returns the possible flags for this parameter.
func (p *ParameterSpec) GetFlagValues() []*FlagsValue {
	var gSize C.guint
	gFlags := C.getParamSpecFlags(p.paramSpec, &gSize)
	size := int(gSize)
	out := make([]*FlagsValue, size)
	for idx, flag := range (*[1 << 30]C.GFlagsValue)(unsafe.Pointer(gFlags))[:size:size] {
		out[idx] = &FlagsValue{
			Value:     int(flag.value),
			ValueNick: C.GoString(flag.value_nick),
			ValueName: C.GoString(flag.value_name),
		}
	}
	return out
}

// GetCaps returns the caps in this parameter if it is of type GST_TYPE_CAPS.
func (p *ParameterSpec) GetCaps() *Caps {
	if p.DefaultValue() == nil {
		return nil
	}
	caps := C.gst_value_get_caps((*C.GValue)(unsafe.Pointer(p.DefaultValue().Native())))
	if caps == nil {
		return nil
	}
	return wrapCaps(caps)
}

// ParameterFlags is a go cast of GParamFlags.
type ParameterFlags int

// Has returns true if these flags contain the provided ones.
func (p ParameterFlags) Has(b ParameterFlags) bool { return p&b != 0 }

// Type casting of GParamFlags
const (
	ParameterReadable       ParameterFlags = C.G_PARAM_READABLE // the parameter is readable
	ParameterWritable                      = C.G_PARAM_WRITABLE // the parameter is writable
	ParameterReadWrite                     = ParameterReadable | ParameterWritable
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
	"explicitly notify",
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
