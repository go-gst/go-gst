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
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// ObjectClass is a loose binding around the glib GObjectClass.
// It forms the base of a GstElementClass.
type ObjectClass struct {
	ptr *C.GObjectClass
}

// Unsafe is a convenience wrapper to return the unsafe.Pointer of the underlying C instance.
func (o *ObjectClass) Unsafe() unsafe.Pointer { return unsafe.Pointer(o.ptr) }

// Instance returns the underlying C GObjectClass pointer
func (o *ObjectClass) Instance() *C.GObjectClass { return o.ptr }

// InstallProperties will install the given ParameterSpecs to the object class.
// They will be IDed in the order they are provided.
func (o *ObjectClass) InstallProperties(params []*ParamSpec) {
	for idx, prop := range params {
		C.g_object_class_install_property(
			o.Instance(),
			C.guint(idx+1),
			prop.paramSpec,
		)
	}
}

// TypeInstance is a loose binding around the glib GTypeInstance. It exposes methods required
// to register the various capabilities of an element.
type TypeInstance struct {
	gtype  C.GType
	gotype GoElement
}

// AddInterface will add an interface implementation for the type referenced by this object.
func (t *TypeInstance) AddInterface(iface glib.Type) {
	ifaceInfo := C.GInterfaceInfo{
		interface_data:     nil,
		interface_finalize: nil,
	}
	switch iface {
	case InterfaceURIHandler:
		globalURIHdlr = t.gotype.(URIHandler)
		ifaceInfo.interface_init = C.GInterfaceInitFunc(C.uriHandlerInit)
	}
	C.g_type_add_interface_static(
		(C.GType)(t.gtype),
		(C.GType)(iface),
		&ifaceInfo,
	)
}
