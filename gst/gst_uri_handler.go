package gst

/*
#include "gst.go.h"

extern GstURIType              goURIHdlrGetURIType        (GType type);
extern const gchar * const *   goURIHdlrGetProtocols      (GType type);
extern gchar *                 goURIHdlrGetURI            (GstURIHandler * handler);
extern gboolean                goURIHdlrSetURI            (GstURIHandler * handler,
                                                           const gchar   * uri,
														   GError       ** error);

void uriHandlerInit (gpointer iface, gpointer iface_data)
{
	((GstURIHandlerInterface*)iface)->get_type = goURIHdlrGetURIType;
	((GstURIHandlerInterface*)iface)->get_protocols = goURIHdlrGetProtocols;
	((GstURIHandlerInterface*)iface)->get_uri = goURIHdlrGetURI;
	((GstURIHandlerInterface*)iface)->set_uri = goURIHdlrSetURI;
}
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

var globalURIHdlr URIHandler

// InterfaceURIHandler represents the GstURIHandler interface GType. Use this when querying bins
// for elements that implement a URIHandler, or when signaling that a GoObjectSubclass provides this
// interface. Note that the way this interface is implemented, it can only be used once per plugin.
var InterfaceURIHandler glib.Interface = &interfaceURIHandler{}

type interfaceURIHandler struct{ glib.Interface }

func (i *interfaceURIHandler) Type() glib.Type { return glib.Type(C.GST_TYPE_URI_HANDLER) }
func (i *interfaceURIHandler) Init(instance *glib.TypeInstance) {
	globalURIHdlr = instance.GoType.(URIHandler)
	C.uriHandlerInit((C.gpointer)(instance.GTypeInstance), nil)
}

// URIHandler represents an interface that elements can implement to provide URI handling
// capabilities.
type URIHandler interface {
	// GetURI gets the currently handled URI.
	GetURI() string
	// GetURIType returns the type of URI this element can handle.
	GetURIType() URIType
	// GetProtocols returns the protocols this element can handle.
	GetProtocols() []string
	// SetURI tries to set the URI of the given handler.
	SetURI(string) (bool, error)
}

// gstURIHandler implements a URIHandler that is backed by an Element from the C API.
type gstURIHandler struct {
	ptr *C.GstElement
}

func (g *gstURIHandler) Instance() *C.GstURIHandler {
	return C.toGstURIHandler(unsafe.Pointer(g.ptr))
}

// GetURI gets the currently handled URI.
func (g *gstURIHandler) GetURI() string {
	ret := C.gst_uri_handler_get_uri(g.Instance())
	if ret == nil {
		return ""
	}
	defer C.g_free((C.gpointer)(unsafe.Pointer(ret)))
	return C.GoString(ret)
}

// GetURIType returns the type of URI this element can handle.
func (g *gstURIHandler) GetURIType() URIType {
	ty := C.gst_uri_handler_get_uri_type((*C.GstURIHandler)(g.Instance()))
	return URIType(ty)
}

// GetProtocols returns the protocols this element can handle.
func (g *gstURIHandler) GetProtocols() []string {
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
