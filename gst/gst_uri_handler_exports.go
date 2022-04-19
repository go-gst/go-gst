package gst

/*
#include "gst.go.h"
*/
import "C"

import (
	"math"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

//export goURIHdlrGetURIType
func goURIHdlrGetURIType(gtype C.GType) C.GstURIType {
	return C.GstURIType(globalURIHdlr.GetURIType())
}

//export goURIHdlrGetProtocols
func goURIHdlrGetProtocols(gtype C.GType) **C.gchar {
	protocols := globalURIHdlr.GetProtocols()
	size := C.size_t(unsafe.Sizeof((*C.gchar)(nil)))
	length := C.size_t(len(protocols))
	arr := (**C.gchar)(C.malloc(length * size))
	view := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.gchar)(nil))]*C.gchar)(unsafe.Pointer(arr))[0:len(protocols):len(protocols)]
	for i, proto := range protocols {
		view[i] = (*C.gchar)(C.CString(proto))
	}
	return arr
}

//export goURIHdlrGetURI
func goURIHdlrGetURI(hdlr *C.GstURIHandler) *C.gchar {
	var uri string
	glib.WithPointerTransferOriginal(unsafe.Pointer(hdlr), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		uri = obj.(URIHandler).GetURI()
	})
	if uri == "" {
		return nil
	}
	return (*C.gchar)(unsafe.Pointer(C.CString(uri)))
}

//export goURIHdlrSetURI
func goURIHdlrSetURI(hdlr *C.GstURIHandler, uri *C.gchar, gerr **C.GError) C.gboolean {
	var ok bool
	var err error
	glib.WithPointerTransferOriginal(unsafe.Pointer(hdlr), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		ok, err = obj.(URIHandler).SetURI(C.GoString(uri))
	})
	if err != nil {
		errMsg := C.CString(err.Error())
		defer C.free(unsafe.Pointer(errMsg))
		C.g_set_error_literal(gerr, DomainLibrary.toQuark(), C.gint(LibraryErrorSettings), errMsg)
	}
	return gboolean(ok)
}
