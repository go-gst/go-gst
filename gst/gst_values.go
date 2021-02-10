package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

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

// Num returns the fraction's numerator.
func (g *FractionValue) Num() int { return g.num }

// Denom returns the fraction's denominator.
func (g *FractionValue) Denom() int { return g.denom }

// String returns a string representation of the fraction.
func (g *FractionValue) String() string { return fmt.Sprintf("%d/%d", g.num, g.denom) }

// ToGValue implements a glib.ValueTransformer.
func (g *FractionValue) ToGValue() (*glib.Value, error) {
	v, err := glib.ValueInit(glib.Type(C.GST_TYPE_FRACTION))
	if err != nil {
		return nil, err
	}
	C.gst_value_set_fraction(
		(*C.GValue)(unsafe.Pointer(v.GValue)),
		C.gint(g.Num()), C.gint(g.Denom()),
	)
	return v, nil
}
