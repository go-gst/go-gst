package gst

/*
#include "gst.go.h"

gint
cgoGstValueCompare (const GValue * value1, const GValue * value2)
{
	return 0;
}

gboolean
cgoGstValueDeserialize (GValue * dest, const gchar * s)
{
	return TRUE;
}

gchar *
cgoGstValueSerialize (const GValue * value1)
{
	return NULL;
}

*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// ValuesCanCompare determines if val1 and val2 can be compared.
func ValuesCanCompare(val1, val2 *glib.Value) bool {
	return gobool(C.gst_value_can_compare(
		(*C.GValue)(unsafe.Pointer(val1.GValue)),
		(*C.GValue)(unsafe.Pointer(val2.GValue)),
	))
}

// ValuesCanIntersect determines if intersecting two values will return
// a valid result. Two values will produce a valid intersection if they
// are the same type.
func ValuesCanIntersect(val1, val2 *glib.Value) bool {
	return gobool(C.gst_value_can_intersect(
		(*C.GValue)(unsafe.Pointer(val1.GValue)),
		(*C.GValue)(unsafe.Pointer(val2.GValue)),
	))
}

// ValuesCanSubtract checks if it's possible to subtract subtrahend from minuend.
func ValuesCanSubtract(minuend, subtrahend *glib.Value) bool {
	return gobool(C.gst_value_can_intersect(
		(*C.GValue)(unsafe.Pointer(minuend.GValue)),
		(*C.GValue)(unsafe.Pointer(subtrahend.GValue)),
	))
}

// ValuesCanUnion determines if val1 and val2 can be non-trivially unioned. Any two
// values can be trivially unioned by adding both of them to a GstValueList. However,
// certain types have the possibility to be unioned in a simpler way. For example, an
// integer range and an integer can be unioned if the integer is a subset of the
// integer range. If there is the possibility that two values can be unioned, this
// function returns TRUE.
func ValuesCanUnion(val1, val2 *glib.Value) bool {
	return gobool(C.gst_value_can_union(
		(*C.GValue)(unsafe.Pointer(val1.GValue)),
		(*C.GValue)(unsafe.Pointer(val2.GValue)),
	))
}

// ValueCmp represents the result of comparing two values.
type ValueCmp int

// ValueCmp castings
const (
	ValueEqual       ValueCmp = C.GST_VALUE_EQUAL        // Indicates that the first value provided to a comparison function (ValueCompare) is equal to the second one.
	ValueGreaterThan ValueCmp = C.GST_VALUE_GREATER_THAN // Indicates that the first value provided to a comparison function (ValueCompare) is greater than the second one.
	ValueLessThan    ValueCmp = C.GST_VALUE_LESS_THAN    // Indicates that the first value provided to a comparison function (ValueCompare) is lesser than the second one.
	ValueUnordered   ValueCmp = C.GST_VALUE_UNORDERED    // Indicates that the comparison function (ValueCompare) can not determine a order for the two provided values.

)

// ValueCompare compares value1 and value2. If value1 and value2 cannot be compared, the function
// returns ValueUnordered. Otherwise, if value1 is greater than value2, ValueGreaterThan is returned.
// If value1 is less than value2, ValueLessThan is returned. If the values are equal, ValueEqual is returned.
func ValueCompare(value1, value2 *glib.Value) ValueCmp {
	return ValueCmp(C.gst_value_compare(
		(*C.GValue)(unsafe.Pointer(value1.GValue)),
		(*C.GValue)(unsafe.Pointer(value2.GValue)),
	))
}

// ValueDeserialize tries to deserialize a string into a glib.Value of the type specified. If the operation succeeds,
// TRUE is returned, FALSE otherwise.
func ValueDeserialize(data string, t glib.Type) (value *glib.Value, ok bool) {
	value, err := glib.ValueInit(t)
	if err != nil {
		return nil, false
	}
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	return value, gobool(C.gst_value_deserialize(
		(*C.GValue)(unsafe.Pointer(value.GValue)),
		(*C.gchar)(cdata),
	))
}

// ValueFixate fixates src into a new value dest. For ranges, the first element is taken. For lists and arrays, the first item
// is fixated and returned. If src is already fixed, this function returns FALSE.
func ValueFixate(src *glib.Value, dest *glib.Value) (ok bool) {
	return gobool(C.gst_value_fixate(
		(*C.GValue)(unsafe.Pointer(dest.GValue)),
		(*C.GValue)(unsafe.Pointer(src.GValue)),
	))
}

// ValueFractionMultiply multiplies the two GValue items containing a TypeFraction and returns the product.
// This function can return false if any error occurs, such as in memory allocation or an integer overflow.
func ValueFractionMultiply(factor1, factor2 *glib.Value) (product *glib.Value, ok bool) {
	out, err := glib.ValueInit(TypeFraction)
	if err != nil {
		return nil, false
	}
	return out, gobool(C.gst_value_fraction_multiply(
		(*C.GValue)(unsafe.Pointer(out.GValue)),
		(*C.GValue)(unsafe.Pointer(factor1.GValue)),
		(*C.GValue)(unsafe.Pointer(factor2.GValue)),
	))
}

// ValueFractionSubtract subtracts the subtrahend from the minuend containing a TypeFraction and returns the result.
// This function can return false if any error occurs, such as in memory allocation or an integer overflow.
func ValueFractionSubtract(minuend, subtrahend *glib.Value) (result *glib.Value, ok bool) {
	out, err := glib.ValueInit(TypeFraction)
	if err != nil {
		return nil, false
	}
	return out, gobool(C.gst_value_fraction_subtract(
		(*C.GValue)(unsafe.Pointer(out.GValue)),
		(*C.GValue)(unsafe.Pointer(minuend.GValue)),
		(*C.GValue)(unsafe.Pointer(subtrahend.GValue)),
	))
}

// ValueGetCaps gets the caps from the given value.
func ValueGetCaps(value *glib.Value) *Caps {
	caps := C.gst_value_get_caps((*C.GValue)(unsafe.Pointer(value.GValue)))
	if caps == nil {
		return nil
	}
	return FromGstCapsUnsafeNone(unsafe.Pointer(caps))
}

// ValueGetCapsFeatures gets the caps features from the given value.
func ValueGetCapsFeatures(value *glib.Value) *CapsFeatures {
	feats := C.gst_value_get_caps_features((*C.GValue)(unsafe.Pointer(value.GValue)))
	if feats == nil {
		return nil
	}
	return &CapsFeatures{native: feats}
}

// ValueGetStructure extracts the GstStructure from a glib.Value, or nil
// if one does not exist.
func ValueGetStructure(gval *glib.Value) *Structure {
	st := C.gst_value_get_structure((*C.GValue)(unsafe.Pointer(gval.GValue)))
	if st == nil {
		return nil
	}
	return wrapStructure(st)
}

// ValueIntersect calculates the intersection of two values. If the values have a non-empty intersection,
// the value representing the intersection isreturned. Otherwise this function returns false. This function
// can also return false for any allocation errors.
func ValueIntersect(value1, value2 *glib.Value) (*glib.Value, bool) {
	out, err := glib.ValueAlloc()
	if err != nil {
		return nil, false
	}
	return out, gobool(C.gst_value_intersect(
		(*C.GValue)(unsafe.Pointer(out.GValue)),
		(*C.GValue)(unsafe.Pointer(value1.GValue)),
		(*C.GValue)(unsafe.Pointer(value2.GValue)),
	))
}

// ValueIsFixed tests if the given GValue, if available in a Structure (or any other container) contains a "fixed"
// (which means: one value) or an "unfixed" (which means: multiple possible values, such as data lists or data ranges) value.
func ValueIsFixed(value *glib.Value) bool {
	return gobool(C.gst_value_is_fixed((*C.GValue)(unsafe.Pointer(value.GValue))))
}

// ValueIsSubset checks that value1 is a subset of value2.
func ValueIsSubset(value1, value2 *glib.Value) bool {
	return gobool(C.gst_value_is_subset(
		(*C.GValue)(unsafe.Pointer(value1.GValue)),
		(*C.GValue)(unsafe.Pointer(value2.GValue)),
	))
}

// ValueSerialize attempts to serialize the given value into a string. An empty string is returned if
// no serializer exists.
func ValueSerialize(value *glib.Value) string {
	str := C.gst_value_serialize(((*C.GValue)(unsafe.Pointer(value.GValue))))
	if str == nil {
		return ""
	}
	defer C.g_free((C.gpointer)(unsafe.Pointer(str)))
	return C.GoString((*C.char)(unsafe.Pointer(str)))
}

// ValueSubtract subtracts subtrahend from minuend and returns the resule. Note that this means subtraction
// as in sets, not as in mathematics. This function can return false if the subtraction is empty or any error
// occurs.
func ValueSubtract(minuend, subtrahend *glib.Value) (*glib.Value, bool) {
	out, err := glib.ValueAlloc()
	if err != nil {
		return nil, false
	}
	return out, gobool(C.gst_value_subtract(
		(*C.GValue)(unsafe.Pointer(out.GValue)),
		(*C.GValue)(unsafe.Pointer(minuend.GValue)),
		(*C.GValue)(unsafe.Pointer(subtrahend.GValue)),
	))
}

// ValueUnion creates a GValue corresponding to the union of value1 and value2.
func ValueUnion(value1, value2 *glib.Value) (*glib.Value, bool) {
	out, err := glib.ValueAlloc()
	if err != nil {
		return nil, false
	}
	return out, gobool(C.gst_value_union(
		(*C.GValue)(unsafe.Pointer(out.GValue)),
		(*C.GValue)(unsafe.Pointer(value1.GValue)),
		(*C.GValue)(unsafe.Pointer(value2.GValue)),
	))
}

// TypeBitmask is the GType for a bitmask value.
var TypeBitmask = glib.Type(C.gst_bitmask_get_type())

// Bitmask represents a bitmask value.
type Bitmask uint64

// ValueGetBitmask gets the bitmask from the given value.
func ValueGetBitmask(value *glib.Value) Bitmask {
	return Bitmask(C.gst_value_get_bitmask((*C.GValue)(unsafe.Pointer(value.GValue))))
}

// ToGValue implements a glib.ValueTransformer
func (b Bitmask) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeBitmask)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_bitmask(
		(*C.GValue)(unsafe.Pointer(val.GValue)),
		C.guint64(b),
	)
	return val, nil
}

// TypeFraction is the GType for a GstFraction
var TypeFraction = glib.Type(C.gst_fraction_get_type())

// FractionValue is a helper structure for building fractions for functions that require them.
type FractionValue struct {
	num, denom int
}

var _ glib.ValueTransformer = &FractionValue{}

// Fraction returns a new GFraction with the given numerator and denominator.
func Fraction(numerator, denominator int) *FractionValue {
	return &FractionValue{num: numerator, denom: denominator}
}

// ValueGetFraction returns the fraction inside the given value.
func ValueGetFraction(value *glib.Value) *FractionValue {
	return &FractionValue{
		num:   int(C.gst_value_get_fraction_numerator((*C.GValue)(unsafe.Pointer(value.GValue)))),
		denom: int(C.gst_value_get_fraction_denominator((*C.GValue)(unsafe.Pointer(value.GValue)))),
	}
}

// Num returns the fraction's numerator.
func (g *FractionValue) Num() int { return g.num }

// Denom returns the fraction's denominator.
func (g *FractionValue) Denom() int { return g.denom }

// String returns a string representation of the fraction.
func (g *FractionValue) String() string { return fmt.Sprintf("%d/%d", g.num, g.denom) }

// ToGValue implements a glib.ValueTransformer.
func (g *FractionValue) ToGValue() (*glib.Value, error) {
	v, err := glib.ValueInit(TypeFraction)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_fraction(
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		C.gint(g.Num()), C.gint(g.Denom()),
	)
	return v, nil
}

// TypeFractionRange is the GType for a GstFractionRange
var TypeFractionRange = glib.Type(C.gst_fraction_range_get_type())

// FractionRangeValue represents a GstFractionRange.
type FractionRangeValue struct {
	start, end *FractionValue
}

var _ glib.ValueTransformer = &FractionRangeValue{}

// FractionRange returns a new GstFractionRange.
func FractionRange(start, end *FractionValue) *FractionRangeValue {
	return &FractionRangeValue{start: start, end: end}
}

// ValueGetFractionRange returns the range inside the given value.
func ValueGetFractionRange(value *glib.Value) *FractionRangeValue {
	start := C.gst_value_get_fraction_range_min((*C.GValue)(unsafe.Pointer(value.GValue)))
	end := C.gst_value_get_fraction_range_max((*C.GValue)(unsafe.Pointer(value.GValue)))
	return &FractionRangeValue{
		start: ValueGetFraction(glib.ValueFromNative(unsafe.Pointer(start))),
		end:   ValueGetFraction(glib.ValueFromNative(unsafe.Pointer(end))),
	}
}

// Start returns the start of the range.
func (g *FractionRangeValue) Start() *FractionValue { return g.start }

// End returns the end of the range.
func (g *FractionRangeValue) End() *FractionValue { return g.end }

// String returns a string representation of the range.
func (g *FractionRangeValue) String() string {
	return fmt.Sprintf("%s - %s", g.Start().String(), g.End().String())
}

// ToGValue implements a glib.ValueTransformer.
func (g *FractionRangeValue) ToGValue() (*glib.Value, error) {
	v, err := glib.ValueInit(TypeFractionRange)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_fraction_range_full(
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		C.gint(g.Start().Num()), C.gint(g.Start().Denom()),
		C.gint(g.End().Num()), C.gint(g.End().Denom()),
	)
	return v, nil
}

// TypeFloat64Range is the GType for a range of 64-bit floating point numbers.
// This is the equivalent of a GstDoubleRange.
var TypeFloat64Range = glib.Type(C.gst_double_range_get_type())

// Float64RangeValue is the go wrapper around a GstDoubleRange value.
type Float64RangeValue struct {
	start, end float64
}

var _ glib.ValueTransformer = &Float64RangeValue{}

// ValueGetFloat64Range returns the range from inside this value.
func ValueGetFloat64Range(value *glib.Value) *Float64RangeValue {
	return &Float64RangeValue{
		start: float64(C.gst_value_get_double_range_min((*C.GValue)(unsafe.Pointer(value.GValue)))),
		end:   float64(C.gst_value_get_double_range_max((*C.GValue)(unsafe.Pointer(value.GValue)))),
	}
}

// Float64Range returns a new Float64RangeValue. This is the equivalent of a double range value.
func Float64Range(start, end float64) *Float64RangeValue {
	return &Float64RangeValue{start: start, end: end}
}

// Start returns the start for this range.
func (f *Float64RangeValue) Start() float64 { return f.start }

// End returns the end for this range.
func (f *Float64RangeValue) End() float64 { return f.end }

// String returns a string representation of the range.
func (f *Float64RangeValue) String() string { return fmt.Sprintf("%v - %v", f.Start(), f.End()) }

// ToGValue implements a glib.ValueTransformer.
func (f *Float64RangeValue) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeFloat64Range)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_double_range(
		(*C.GValue)(unsafe.Pointer(val.GValue)),
		C.gdouble(f.Start()),
		C.gdouble(f.End()),
	)
	return val, nil
}

// TypeFlagset is the GType for a Flagset
var TypeFlagset = glib.Type(C.gst_flagset_get_type())

// FlagsetValue is the go wrapper around a GstFlagSet value.
type FlagsetValue struct {
	flags, mask uint
}

var _ glib.ValueTransformer = &FlagsetValue{}

// ValueGetFlagset returns the flagset inside his value.
func ValueGetFlagset(value *glib.Value) *FlagsetValue {
	return &FlagsetValue{
		flags: uint(C.gst_value_get_flagset_flags((*C.GValue)(unsafe.Pointer(value.GValue)))),
		mask:  uint(C.gst_value_get_flagset_mask((*C.GValue)(unsafe.Pointer(value.GValue)))),
	}
}

// Flagset returns a new FlagsetValue. The flags value indicates the values of flags,
// the mask represents which bits in the flag value have been set, and which are "don't care".
func Flagset(flags, mask uint) *FlagsetValue {
	return &FlagsetValue{flags: flags, mask: mask}
}

// Flags returns the flags for this flagset.
func (f *FlagsetValue) Flags() uint { return f.flags }

// Mask returns the mask for this flagset.
func (f *FlagsetValue) Mask() uint { return f.mask }

// ToGValue implements a glib.ValueTransformer.
func (f *FlagsetValue) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeFlagset)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_flagset(
		(*C.GValue)(unsafe.Pointer(val.GValue)),
		C.guint(f.Flags()),
		C.guint(f.Mask()),
	)
	return val, nil
}

// TypeInt64Range is the GType for a range of 64-bit integers.
var TypeInt64Range = glib.Type(C.gst_int64_range_get_type())

// Int64RangeValue represents a GstInt64Range.
type Int64RangeValue struct {
	start, end, step int64
}

var _ glib.ValueTransformer = &Int64RangeValue{}

// Int64Range returns a new Int64RangeValue.
func Int64Range(start, end, step int64) *Int64RangeValue {
	return &Int64RangeValue{
		start: start, end: end, step: step,
	}
}

// ValueGetInt64Range gets the int64 range from the given value.
func ValueGetInt64Range(value *glib.Value) *Int64RangeValue {
	return &Int64RangeValue{
		start: int64(C.gst_value_get_int64_range_min((*C.GValue)(unsafe.Pointer(value.GValue)))),
		end:   int64(C.gst_value_get_int64_range_max((*C.GValue)(unsafe.Pointer(value.GValue)))),
		step:  int64(C.gst_value_get_int64_range_step((*C.GValue)(unsafe.Pointer(value.GValue)))),
	}
}

// Start returns the start of the range
func (i *Int64RangeValue) Start() int64 { return i.start }

// End returns the end of the range
func (i *Int64RangeValue) End() int64 { return i.end }

// Step returns the step of the range
func (i *Int64RangeValue) Step() int64 { return i.step }

// String implements a Stringer.
func (i *Int64RangeValue) String() string {
	return fmt.Sprintf("%d - %d (%d)", i.Start(), i.End(), i.Step())
}

// ToGValue implements a glib.ValueTransformer.
func (i *Int64RangeValue) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeInt64Range)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_int64_range_step(
		(*C.GValue)(unsafe.Pointer(val.GValue)),
		C.gint64(i.Start()),
		C.gint64(i.End()),
		C.gint64(i.Step()),
	)
	return val, nil
}

// TypeIntRange is the GType for a range of integers.
var TypeIntRange = glib.Type(C.gst_int_range_get_type())

// IntRangeValue represents a GstIntRange.
type IntRangeValue struct {
	start, end, step int
}

var _ glib.ValueTransformer = &Int64RangeValue{}

// IntRange returns a new IntRangeValue.
func IntRange(start, end, step int) *IntRangeValue {
	return &IntRangeValue{
		start: start, end: end, step: step,
	}
}

// ValueGetIntRange gets the int range from the given value.
func ValueGetIntRange(value *glib.Value) *IntRangeValue {
	return &IntRangeValue{
		start: int(C.gst_value_get_int_range_min((*C.GValue)(unsafe.Pointer(value.GValue)))),
		end:   int(C.gst_value_get_int_range_max((*C.GValue)(unsafe.Pointer(value.GValue)))),
		step:  int(C.gst_value_get_int_range_step((*C.GValue)(unsafe.Pointer(value.GValue)))),
	}
}

// Start returns the start of the range
func (i *IntRangeValue) Start() int { return i.start }

// End returns the end of the range
func (i *IntRangeValue) End() int { return i.end }

// Step returns the step of the range
func (i *IntRangeValue) Step() int { return i.step }

// String implements a Stringer.
func (i *IntRangeValue) String() string {
	return fmt.Sprintf("%d - %d (%d)", i.Start(), i.End(), i.Step())
}

// ToGValue implements a glib.ValueTransformer.
func (i *IntRangeValue) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeIntRange)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_int_range_step(
		(*C.GValue)(unsafe.Pointer(val.GValue)),
		C.gint(i.Start()),
		C.gint(i.End()),
		C.gint(i.Step()),
	)
	return val, nil
}

// TypeValueArray is the GType for a GstValueArray
var TypeValueArray = glib.Type(C.gst_value_array_get_type())

// ValueArrayValue represets a GstValueArray.
type ValueArrayValue glib.Value

// ValueArray converts the given slice of Go types into a ValueArrayValue.
// This function can return nil on any conversion or memory allocation errors.
func ValueArray(ss []interface{}) *ValueArrayValue {
	v, err := glib.ValueAlloc()
	if err != nil {
		return nil
	}
	C.gst_value_array_init(
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		C.guint(len(ss)),
	)
	for _, s := range ss {
		val, err := glib.GValue(s)
		if err != nil {
			return nil
		}
		C.gst_value_array_append_value(
			(*C.GValue)(unsafe.Pointer(v.GValue)),
			(*C.GValue)(unsafe.Pointer(val.GValue)),
		)
	}
	out := ValueArrayValue(*v)
	return &out
}

// Size returns the size of the array.
func (v *ValueArrayValue) Size() uint {
	return uint(C.gst_value_array_get_size((*C.GValue)(unsafe.Pointer(v.GValue))))
}

// ValueAt returns the value at the index in the array, or nil on any error.
func (v *ValueArrayValue) ValueAt(idx uint) interface{} {
	gval := C.gst_value_array_get_value(
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		C.guint(idx),
	)
	if gval == nil {
		return nil
	}
	out, err := glib.ValueFromNative(unsafe.Pointer(gval)).GoValue()
	if err != nil {
		return nil
	}
	return out
}

// ToGValue implements a glib.ValueTransformer.
func (v *ValueArrayValue) ToGValue() (*glib.Value, error) {
	out := glib.Value(*v)
	return &out, nil
}

// TypeValueList is the GType for a GstValueList
var TypeValueList = glib.Type(C.gst_value_list_get_type())

// ValueListValue represets a GstValueList.
type ValueListValue glib.Value

// ValueList converts the given slice of Go types into a ValueListValue.
// This function can return nil on any conversion or memory allocation errors.
func ValueList(ss []interface{}) *ValueListValue {
	v, err := glib.ValueAlloc()
	if err != nil {
		return nil
	}
	C.gst_value_list_init(
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		C.guint(len(ss)),
	)
	for _, s := range ss {
		val, err := glib.GValue(s)
		if err != nil {
			return nil
		}
		C.gst_value_list_append_value(
			(*C.GValue)(unsafe.Pointer(v.GValue)),
			(*C.GValue)(unsafe.Pointer(val.GValue)),
		)
	}
	out := ValueListValue(*v)
	return &out
}

// Size returns the size of the list.
func (v *ValueListValue) Size() uint {
	return uint(C.gst_value_list_get_size((*C.GValue)(unsafe.Pointer(v.GValue))))
}

// ValueAt returns the value at the index in the lise, or nil on any error.
func (v *ValueListValue) ValueAt(idx uint) interface{} {
	gval := C.gst_value_list_get_value(
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		C.guint(idx),
	)
	if gval == nil {
		return nil
	}
	out, err := glib.ValueFromNative(unsafe.Pointer(gval)).GoValue()
	if err != nil {
		return nil
	}
	return out
}

// Concat concatenates copies of this list and value into a new list. Values that are not of type
// TypeValueList are treated as if they were lists of length 1. dest will be initialized to the type
// TypeValueList.
func (v *ValueListValue) Concat(value *ValueListValue) *ValueListValue {
	out, err := glib.ValueAlloc()
	if err != nil {
		return nil
	}
	C.gst_value_list_concat(
		(*C.GValue)(unsafe.Pointer(out.GValue)),
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		(*C.GValue)(unsafe.Pointer(value.GValue)),
	)
	o := ValueListValue(*out)
	return &o
}

// Merge merges copies of value into this list. Values that are not of type TypeValueList are treated as
// if they were lists of length 1.
//
// The result will be put into a new value and will either be a list that will not contain any duplicates,
// or a non-list type (if the lists were equal).
func (v *ValueListValue) Merge(value *ValueListValue) *ValueListValue {
	out, err := glib.ValueAlloc()
	if err != nil {
		return nil
	}
	C.gst_value_list_merge(
		(*C.GValue)(unsafe.Pointer(out.GValue)),
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		(*C.GValue)(unsafe.Pointer(value.GValue)),
	)
	o := ValueListValue(*out)
	return &o
}

// ToGValue implements a glib.ValueTransformer.
func (v *ValueListValue) ToGValue() (*glib.Value, error) {
	out := glib.Value(*v)
	return &out, nil
}
