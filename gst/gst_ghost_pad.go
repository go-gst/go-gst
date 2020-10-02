package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// GhostPad is a go representation of a GstGhostPad.
type GhostPad struct{ *ProxyPad }

// NewGhostPad create a new ghostpad with target as the target. The direction will be
// taken from the target pad. The target must be unlinked.
//
// Will ref the target.
func NewGhostPad(name string, target *Pad) *GhostPad {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	pad := C.gst_ghost_pad_new(
		(*C.gchar)(unsafe.Pointer(cName)),
		target.Instance(),
	)
	if pad == nil {
		return nil
	}
	return wrapGhostPad(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(pad))})
}

// ProxyPad is a go representation of a GstProxyPad.
type ProxyPad struct{ *Pad }
