package gstnet

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// NetAddressMeta can be used to store a network address in a GstBuffer so that it network elements
// can track the to and from address of the buffer.
type NetAddressMeta struct{ ptr *C.GstNetAddressMeta }

// AddNetAddressMeta attaches the given address to a NetAddressMeta on the buffer.
func AddNetAddressMeta(buffer *gst.Buffer, address string, port int) *NetAddressMeta {
	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))
	gaddr := C.g_inet_socket_address_new_from_string((*C.gchar)(caddr), C.guint(port))
	meta := C.gst_buffer_add_net_address_meta(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
		gaddr,
	)
	return &NetAddressMeta{meta}
}

// GetNetAddressMeta retrieves the NetAddressMeta from the given buffer.
func GetNetAddressMeta(buffer *gst.Buffer) *NetAddressMeta {
	meta := C.gst_buffer_get_net_address_meta((*C.GstBuffer)(unsafe.Pointer(buffer.Instance())))
	return &NetAddressMeta{meta}
}

// Meta returns the underlying gst.Meta instance.
func (n *NetAddressMeta) Meta() *gst.Meta { return gst.FromGstMetaUnsafe(unsafe.Pointer(&n.ptr.meta)) }

// Addr returns the address included in the meta.
func (n *NetAddressMeta) Addr() string {
	iaddr := C.g_inet_socket_address_get_address((*C.GInetSocketAddress)(unsafe.Pointer(n.ptr.addr)))
	iaddrstr := C.g_inet_address_to_string(iaddr)
	defer C.g_free((C.gpointer)(unsafe.Pointer(iaddrstr)))
	return C.GoString(iaddrstr)
}

// Port returns the port included in the meta.
func (n *NetAddressMeta) Port() int {
	iport := C.g_inet_socket_address_get_port((*C.GInetSocketAddress)(unsafe.Pointer(n.ptr.addr)))
	return int(iport)
}
