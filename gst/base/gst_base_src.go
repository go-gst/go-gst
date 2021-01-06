package base

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// GstBaseSrc represents a GstBaseSrc.
type GstBaseSrc struct{ *gst.Element }

// ToGstBaseSrc returns a GstBaseSrc object for the given object.
func ToGstBaseSrc(obj *gst.Object) *GstBaseSrc {
	return &GstBaseSrc{&gst.Element{Object: obj}}
}

// wrapGstBaseSrc wraps the given unsafe.Pointer in a GstBaseSrc instance.
func wrapGstBaseSrc(obj *C.GstBaseSrc) *GstBaseSrc {
	return &GstBaseSrc{gst.FromGstElementUnsafe(unsafe.Pointer(obj))}
}

// Instance returns the underlying C GstBaseSrc instance
func (g *GstBaseSrc) Instance() *C.GstBaseSrc {
	return C.toGstBaseSrc(g.Unsafe())
}

// SetFormat sets the default format of the source. This will be the format used for sending
// SEGMENT events and for performing seeks.
//
// If a format of gst.FormatBytes is set, the element will be able to operate in pull mode if the
// IsSeekable returns TRUE.
//
// This function must only be called in when the element is paused.
func (g *GstBaseSrc) SetFormat(format gst.Format) {
	C.gst_base_src_set_format(g.Instance(), C.GstFormat(format))
}

// StartComplete completes an asynchronous start operation. When the subclass overrides the start method,
// it should call StartComplete when the start operation completes either from the same thread or from an
// asynchronous helper thread.
func (g *GstBaseSrc) StartComplete(ret gst.FlowReturn) {
	C.gst_base_src_start_complete(g.Instance(), C.GstFlowReturn(ret))
}
