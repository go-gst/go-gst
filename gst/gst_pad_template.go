package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// PadTemplate is a go representation of a GstPadTemplate
type PadTemplate struct{ *Object }

// FromGstPadTemplateUnsafeNone wraps the given GstPadTemplate in a ref and a finalizer.
func FromGstPadTemplateUnsafeNone(tmpl unsafe.Pointer) *PadTemplate {
	return &PadTemplate{wrapObject(glib.TransferNone(tmpl))}
}

// FromGstPadTemplateUnsafeFull wraps the given GstPadTemplate in a finalizer.
func FromGstPadTemplateUnsafeFull(tmpl unsafe.Pointer) *PadTemplate {
	return &PadTemplate{wrapObject(glib.TransferFull(tmpl))}
}

// NewPadTemplate creates a new pad template with a name according to the given template and with the given arguments.
func NewPadTemplate(nameTemplate string, direction PadDirection, presence PadPresence, caps *Caps) *PadTemplate {
	cName := C.CString(nameTemplate)
	defer C.free(unsafe.Pointer(cName))
	tmpl := C.gst_pad_template_new(
		(*C.gchar)(cName),
		C.GstPadDirection(direction),
		C.GstPadPresence(presence),
		caps.Instance(),
	)
	if tmpl == nil {
		return nil
	}
	return wrapPadTemplate(glib.TransferNone(unsafe.Pointer(tmpl)))
}

// NewPadTemplateWithGType creates a new pad template with a name according to the given template and with the given arguments.
func NewPadTemplateWithGType(nameTemplate string, direction PadDirection, presence PadPresence, caps *Caps, gType glib.Type) *PadTemplate {
	cName := C.CString(nameTemplate)
	defer C.free(unsafe.Pointer(cName))
	tmpl := C.gst_pad_template_new_with_gtype(
		(*C.gchar)(cName),
		C.GstPadDirection(direction),
		C.GstPadPresence(presence),
		caps.Instance(),
		(C.GType)(gType),
	)
	if tmpl == nil {
		return nil
	}
	return wrapPadTemplate(glib.TransferNone(unsafe.Pointer(tmpl)))
}

// Instance returns the underlying C GstPadTemplate.
func (p *PadTemplate) Instance() *C.GstPadTemplate { return C.toGstPadTemplate(p.Unsafe()) }

// Name returns the name of the pad template.
func (p *PadTemplate) Name() string { return C.GoString(p.Instance().name_template) }

// Direction returns the direction of the pad template.
func (p *PadTemplate) Direction() PadDirection { return PadDirection(p.Instance().direction) }

// Presence returns the presence of the pad template.
func (p *PadTemplate) Presence() PadPresence { return PadPresence(p.Instance().presence) }

// Caps returns the caps of the pad template.
func (p *PadTemplate) Caps() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_pad_template_get_caps(p.Instance())))
}

// PadCreated emits the pad-created signal for this template when created by this pad.
func (p *PadTemplate) PadCreated(pad *Pad) {
	C.gst_pad_template_pad_created(p.Instance(), pad.Instance())
}

// // GetDocumentationCaps gets the documentation caps for the template. See SetDocumentationCaps for more information.
// func (p *PadTemplate) GetDocumentationCaps() *Caps {
// 	return wrapCaps(C.gst_pad_template_get_documentation_caps(p.Instance()))
// }

// // SetDocumentationCaps sets caps to be exposed to a user. Certain elements will dynamically construct the caps of
// // their pad templates. In order not to let environment-specific information into the documentation, element authors
// // should use this method to expose "stable" caps to the reader.
// func (p *PadTemplate) SetDocumentationCaps(caps *Caps) {
// 	C.gst_pad_template_set_documentation_caps(p.Instance(), caps.Instance())
// }
