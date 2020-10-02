package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Pad is a go representation of a GstPad
type Pad struct{ *Object }

// NewPad returns a new pad with the given direction. If name is empty, one will be generated for you.
func NewPad(name string, direction PadDirection) *Pad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_pad_new(cName, C.GstPadDirection(direction))
	if pad == nil {
		return nil
	}
	return wrapPad(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(pad))})
}

// NewPadFromTemplate creates a new pad with the given name from the given template. If name is empty, one will
// be generated for you.
func NewPadFromTemplate(tmpl *PadTemplate, name string) *Pad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_pad_new_from_template(tmpl.Instance(), cName)
	if pad == nil {
		return nil
	}
	return wrapPad(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(pad))})
}

// Instance returns the underlying C GstPad.
func (p *Pad) Instance() *C.GstPad { return C.toGstPad(p.Unsafe()) }

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
