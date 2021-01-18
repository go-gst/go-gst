package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// InterfaceTOCSetter represents the GstTocSetter interface GType. Use this when querying bins
// for elements that implement a TOCSetter.
var InterfaceTOCSetter = glib.Type(C.GST_TYPE_TOC_SETTER)

// TOCSetter is an interface that elements can implement to provide TOC writing capabilities.
type TOCSetter interface {
	// Return current TOC the setter uses. The TOC should not be modified without making it writable first.
	GetTOC() *TOC
	// Set the given TOC on the setter. Previously set TOC will be unreffed before setting a new one.
	SetTOC(*TOC)
	// Reset the internal TOC. Elements should call this from within the state-change handler.
	Reset()
}

// gstTocSetter implements a TOCSetter that is backed by an Element from the C runtime.
type gstTOCSetter struct {
	ptr *C.GstElement
}

func (g *gstTOCSetter) Instance() *C.GstTocSetter {
	return C.toTocSetter(g.ptr)
}

func (g *gstTOCSetter) GetTOC() *TOC {
	toc := C.gst_toc_setter_get_toc(g.Instance())
	if toc == nil {
		return nil
	}
	return FromGstTOCUnsafeFull(unsafe.Pointer(toc))
}

func (g *gstTOCSetter) SetTOC(toc *TOC) { C.gst_toc_setter_set_toc(g.Instance(), toc.Instance()) }
func (g *gstTOCSetter) Reset()          { C.gst_toc_setter_reset(g.Instance()) }
