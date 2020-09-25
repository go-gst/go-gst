package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Pad is a go representation of a GstPad
type Pad struct{ *Object }

// Instance returns the underlying C GstPad.
func (p *Pad) Instance() *C.GstPad { return C.toGstPad(p.unsafe()) }

// Direction returns the direction of this pad.
func (p *Pad) Direction() PadDirection {
	return PadDirection(C.gst_pad_get_direction((*C.GstPad)(p.Instance())))
}

// Template returns the template for this pad or nil.
func (p *Pad) Template() *PadTemplate {
	return wrapPadTemplate(glib.Take(unsafe.Pointer(p.Instance().padtemplate)))
}

// CurrentCaps returns the caps for this Pad or nil.
func (p *Pad) CurrentCaps() *Caps {
	caps := C.gst_pad_get_current_caps((*C.GstPad)(p.Instance()))
	if caps == nil {
		return nil
	}
	return wrapCaps(caps)
}

// PadTemplate is a go representation of a GstPadTemplate
type PadTemplate struct{ *Object }

// Instance returns the underlying C GstPadTemplate.
func (p *PadTemplate) Instance() *C.GstPadTemplate { return C.toGstPadTemplate(p.unsafe()) }

// Name returns the name of the pad template.
func (p *PadTemplate) Name() string { return C.GoString(p.Instance().name_template) }

// Direction returns the direction of the pad template.
func (p *PadTemplate) Direction() PadDirection { return PadDirection(p.Instance().direction) }

// Presence returns the presence of the pad template.
func (p *PadTemplate) Presence() PadPresence { return PadPresence(p.Instance().presence) }

// Caps returns the caps of the pad template.
func (p *PadTemplate) Caps() *Caps { return wrapCaps(p.Instance().caps) }

// GhostPad is a go representation of a GstGhostPad
type GhostPad struct{ *Pad }
