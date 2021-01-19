package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// ElementClass represents the subclass of an element provided by a plugin.
type ElementClass struct{ *glib.ObjectClass }

// ToElementClass wraps the given ObjectClass in an ElementClass instance.
func ToElementClass(klass *glib.ObjectClass) *ElementClass {
	return &ElementClass{klass}
}

// Instance returns the underlying GstElementClass instance.
func (e *ElementClass) Instance() *C.GstElementClass {
	return C.toGstElementClass(e.Unsafe())
}

// AddMetadata sets key with the given value in the metadata of the class.
func (e *ElementClass) AddMetadata(key, value string) {
	C.gst_element_class_add_static_metadata(
		e.Instance(),
		(*C.gchar)(C.CString(key)),
		(*C.gchar)(C.CString(value)),
	)
}

// AddPadTemplate adds a padtemplate to an element class. This is mainly used in the
// ClassInit functions of ObjectSubclasses. If a pad template with the same name as an
// already existing one is added the old one is replaced by the new one.
//
// templ's reference count will be incremented, and any floating reference will be removed
func (e *ElementClass) AddPadTemplate(templ *PadTemplate) {
	C.gst_element_class_add_pad_template(
		e.Instance(),
		(*C.GstPadTemplate)(templ.Unsafe()),
	)
}

// AddStaticPadTemplate adds a pad template to an element class based on the pad template templ. The template
// is first converted to a static pad template.
//
// This is mainly used in the ClassInit functions of element implementations. If a pad template with the
// same name already exists, the old one is replaced by the new one.
func (e *ElementClass) AddStaticPadTemplate(templ *PadTemplate) {
	staticTmpl := C.GstStaticPadTemplate{
		name_template: templ.Instance().name_template,
		direction:     templ.Instance().direction,
		presence:      templ.Instance().presence,
		static_caps: C.GstStaticCaps{
			caps:   templ.Caps().Ref().Instance(),
			string: C.CString(templ.Name()),
		},
	}
	C.gst_element_class_add_static_pad_template(
		e.Instance(),
		&staticTmpl,
	)
}

// GetMetadata retrieves the metadata associated with key in the class.
func (e *ElementClass) GetMetadata(key string) string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return C.GoString(C.gst_element_class_get_metadata(e.Instance(), (*C.gchar)(ckey)))
}

// GetPadTemplate retrieves the padtemplate with the given name. No unrefing is necessary.
// If no pad template exists with the given name, nil is returned.
func (e *ElementClass) GetPadTemplate(name string) *PadTemplate {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	tmpl := C.gst_element_class_get_pad_template(e.Instance(), (*C.gchar)(cname))
	if tmpl == nil {
		return nil
	}
	return FromGstPadTemplateUnsafeNone(unsafe.Pointer(tmpl))
}

// GetAllPadTemplates retrieves a slice of all the pad templates associated with this class.
// The list must not be modified.
func (e *ElementClass) GetAllPadTemplates() []*PadTemplate {
	glist := C.gst_element_class_get_pad_template_list(e.Instance())
	return glistToPadTemplateSlice(glist)
}

// SetMetadata sets the detailed information for this class.
//
// `longname` - The english long name of the element. E.g "File Sink"
//
// `classification` - A string describing the type of element, as an unordered list separated with slashes ('/'). E.g: "Sink/File"
//
// `description` - Sentence describing the purpose of the element. E.g: "Write stream to a file"
//
// `author` - Name and contact details of the author(s). Use \n to separate multiple author metadata. E.g: "Joe Bloggs <joe.blogs at foo.com>"
func (e *ElementClass) SetMetadata(longname, classification, description, author string) {
	C.gst_element_class_set_static_metadata(
		e.Instance(),
		(*C.gchar)(C.CString(longname)),
		(*C.gchar)(C.CString(classification)),
		(*C.gchar)(C.CString(description)),
		(*C.gchar)(C.CString(author)),
	)
}
