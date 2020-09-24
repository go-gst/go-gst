package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"
import "unsafe"

// Pad is a go representation of a GstPad
type Pad struct{ *Object }

// Instance returns the underlying C GstPad.
func (p *Pad) Instance() *C.GstPad { return C.toGstPad(p.unsafe()) }

// Direction returns the direction of this pad.
func (p *Pad) Direction() PadDirection {
	return PadDirection(C.gst_pad_get_direction((*C.GstPad)(p.Instance())))
}

// Template returns the template for this pad or nil.
func (p *Pad) Template() *PadTemplate { return wrapPadTemplate(p.Instance().padtemplate) }

// CurrentCaps returns the caps for this Pad or nil.
func (p *Pad) CurrentCaps() Caps {
	caps := C.gst_pad_get_current_caps((*C.GstPad)(p.Instance()))
	if caps == nil {
		return nil
	}
	defer C.gst_caps_unref(caps)
	return FromGstCaps(caps)
}

func wrapPad(p *C.GstPad) *Pad {
	return &Pad{wrapObject(C.toGstObject(unsafe.Pointer(p)))}
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
func (p *PadTemplate) Caps() Caps { return FromGstCaps(p.Instance().caps) }

func wrapPadTemplate(p *C.GstPadTemplate) *PadTemplate {
	return &PadTemplate{wrapObject(C.toGstObject(unsafe.Pointer(p)))}
}

// PadDirection is a cast of GstPadDirection to a go type.
type PadDirection C.GstPadDirection

// Type casting of pad directions
const (
	PadUnknown PadDirection = C.GST_PAD_UNKNOWN // (0) - the direction is unknown
	PadSource               = C.GST_PAD_SRC     // (1) - the pad is a source pad
	PadSink                 = C.GST_PAD_SINK    // (2) - the pad is a sink pad
)

// String implements a Stringer on PadDirection.
func (p PadDirection) String() string {
	switch p {
	case PadUnknown:
		return "Unknown"
	case PadSource:
		return "Src"
	case PadSink:
		return "Sink"
	}
	return ""
}

// PadPresence is a cast of GstPadPresence to a go type.
type PadPresence C.GstPadPresence

// Type casting of pad presences
const (
	PadAlways    PadPresence = C.GST_PAD_ALWAYS    // (0) - the pad is always available
	PadSometimes             = C.GST_PAD_SOMETIMES // (1) - the pad will become available depending on the media stream
	PadRequest               = C.GST_PAD_REQUEST   // (2) - the pad is only available on request with gst_element_request_pad.
)

// String implements a stringer on PadPresence.
func (p PadPresence) String() string {
	switch p {
	case PadAlways:
		return "Always"
	case PadSometimes:
		return "Sometimes"
	case PadRequest:
		return "Request"
	}
	return ""
}
