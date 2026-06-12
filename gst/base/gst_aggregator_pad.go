package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

// GstAggregatorPad represents a GstAggregatorPad
type GstAggregatorPad struct {
	*gst.Pad
}

// ToGstAggregatorPad returns a GstAggregatorPad for the an object.
func ToGstAggregatorPad(obj any) *GstAggregatorPad {
	switch obj := obj.(type) {
	case *gst.Pad:
		return &GstAggregatorPad{Pad: obj}
	case *glib.Object:
		return &GstAggregatorPad{Pad: gst.FromGstPadUnsafeNone(unsafe.Pointer(obj.Unsafe()))}
	}
	return nil
}

// Instance returns the underlying C GstAggregatorPad.
func (g *GstAggregatorPad) Instance() *C.GstAggregatorPad {
	return C.toGstAggregatorPad(g.Unsafe())
}

// IsEOS returns true if the pad is EOS, otherwise false.
func (g *GstAggregatorPad) IsEOS() bool {
	return gobool(C.gst_aggregator_pad_is_eos(g.Instance()))
}
