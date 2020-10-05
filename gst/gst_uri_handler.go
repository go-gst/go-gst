package gst

// #include "gst.go.h"
import "C"
import (
	"errors"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// InterfaceURIHandler represents the GstURIHandler interface GType. Use this when querying bins
// for elements that implement a URIHandler.
var InterfaceURIHandler = glib.Type(C.GST_TYPE_URI_HANDLER)

// URIHandler represents an interface that elements can implement to provide URI handling
// capabilities.
type URIHandler interface {
	// GetURI gets the currently handled URI.
	GetURI() string
	// GetURIType returns the type of URI this element can handle.
	GetURIType() URIType
	// GetURIProtocols returns the protocols this element can handle.
	GetURIProtocols() []string
	// SetURI tries to set the URI of the given handler.
	SetURI(string) (bool, error)
}

// gstURIHandler implements a URIHandler that is backed by an Element from the C runtime.
type gstURIHandler struct {
	ptr *C.GstElement
}

func (g *gstURIHandler) Instance() *C.GstURIHandler {
	return C.toGstURIHandler(unsafe.Pointer(g.ptr))
}

// GetURI gets the currently handled URI.
func (g *gstURIHandler) GetURI() string {
	ret := C.gst_uri_handler_get_uri(g.Instance())
	defer C.g_free((C.gpointer)(unsafe.Pointer(ret)))
	return C.GoString(ret)
}

// GetURIType returns the type of URI this element can handle.
func (g *gstURIHandler) GetURIType() URIType {
	ty := C.gst_uri_handler_get_uri_type((*C.GstURIHandler)(g.Instance()))
	return URIType(ty)
}

// GetURIProtocols returns the protocols this element can handle.
func (g *gstURIHandler) GetURIProtocols() []string {
	protocols := C.gst_uri_handler_get_protocols((*C.GstURIHandler)(g.Instance()))
	if protocols == nil {
		return nil
	}
	size := C.sizeOfGCharArray(protocols)
	return goStrings(size, protocols)
}

// SetURI tries to set the URI of the given handler.
func (g *gstURIHandler) SetURI(uri string) (bool, error) {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	var gerr *C.GError
	ret := C.gst_uri_handler_set_uri(
		g.Instance(),
		(*C.gchar)(unsafe.Pointer(curi)),
		&gerr,
	)
	if gerr != nil {
		defer C.g_error_free(gerr)
		return gobool(ret), errors.New(C.GoString(gerr.message))
	}
	return gobool(ret), nil
}
