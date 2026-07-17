package gstnet

// #cgo pkg-config: gstreamer-net-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/net/net.h>
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/pkg/gio/v2"
)

// Addr returns the address included in the meta.
func (n *NetAddressMeta) Addr() gio.SocketAddress {
	if n == nil || n.native == nil || n.native.addr == nil {
		return nil
	}
	return gio.UnsafeSocketAddressFromGlibNone(unsafe.Pointer(n.native.addr))
}
