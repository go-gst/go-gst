package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

// GstPushSrc represents a GstBaseSrc.
type GstPushSrc struct{ *GstBaseSrc }

// ToGstPushSrc returns a GstPushSrc object for the given object.
func ToGstPushSrc(obj interface{}) *GstPushSrc {
	switch obj := obj.(type) {
	case *gst.Object:
		return &GstPushSrc{&GstBaseSrc{&gst.Element{Object: obj}}}
	case *glib.Object:
		return &GstPushSrc{&GstBaseSrc{&gst.Element{Object: &gst.Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: obj}}}}}
	}
	return nil
}

// wrapGstPushSrc wraps the given unsafe.Pointer in a GstPushSrc instance.
func wrapGstPushSrc(obj *C.GstPushSrc) *GstPushSrc {
	return &GstPushSrc{&GstBaseSrc{gst.FromGstElementUnsafeNone(unsafe.Pointer(obj))}}
}

// Instance returns the underlying C GstBaseSrc instance
func (g *GstPushSrc) Instance() *C.GstPushSrc {
	return C.toGstPushSrc(g.Unsafe())
}
