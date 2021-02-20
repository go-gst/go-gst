package gst

// #include "gst.go.h"
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// GhostPad is a go representation of a GstGhostPad.
type GhostPad struct{ *ProxyPad }

// FromGstGhostPadUnsafeNone wraps the given GstGhostPad.
func FromGstGhostPadUnsafeNone(pad unsafe.Pointer) *GhostPad {
	return &GhostPad{&ProxyPad{&Pad{wrapObject(glib.TransferNone(pad))}}}
}

// FromGstGhostPadUnsafeFull wraps the given GstGhostPad.
func FromGstGhostPadUnsafeFull(pad unsafe.Pointer) *GhostPad {
	return &GhostPad{&ProxyPad{&Pad{wrapObject(glib.TransferFull(pad))}}}

}

// NewGhostPad create a new ghostpad with target as the target. The direction will be
// taken from the target pad. The target must be unlinked. If name is empty, one will be
// selected.
//
// Will ref the target.
func NewGhostPad(name string, target *Pad) *GhostPad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_ghost_pad_new(
		cName,
		target.Instance(),
	)
	if pad == nil {
		return nil
	}
	return FromGstGhostPadUnsafeNone(unsafe.Pointer(pad))
}

// NewGhostPadFromTemplate creates a new ghostpad with target as the target. The direction will be taken
// from the target pad. The template used on the ghostpad will be template. If name is empty one will be
// selected.
//
// Will ref the target.
func NewGhostPadFromTemplate(name string, target *Pad, tmpl *PadTemplate) *GhostPad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_ghost_pad_new_from_template(
		cName,
		target.Instance(),
		tmpl.Instance(),
	)
	if pad == nil {
		return nil
	}
	return FromGstGhostPadUnsafeNone(unsafe.Pointer(pad))
}

// NewGhostPadNoTarget creates a new ghostpad without a target with the given direction. A target can be set on the
// ghostpad later with the SetTarget function. If name is empty, one will be selected.
//
// The created ghostpad will not have a padtemplate.
func NewGhostPadNoTarget(name string, direction PadDirection) *GhostPad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_ghost_pad_new_no_target(
		cName,
		C.GstPadDirection(direction),
	)
	if pad == nil {
		return nil
	}
	return FromGstGhostPadUnsafeNone(unsafe.Pointer(pad))
}

// NewGhostPadNoTargetFromTemplate creates a new ghostpad based on templ, without setting a target. The direction will be taken
// from the templ.
func NewGhostPadNoTargetFromTemplate(name string, tmpl *PadTemplate) *GhostPad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_ghost_pad_new_no_target_from_template(
		cName,
		tmpl.Instance(),
	)
	if pad == nil {
		return nil
	}
	return FromGstGhostPadUnsafeNone(unsafe.Pointer(pad))
}

// Instance returns the underlying ghost pad instance.
func (g *GhostPad) Instance() *C.GstGhostPad { return C.toGstGhostPad(g.Unsafe()) }

// GetTarget gets the target pad of gpad.
func (g *GhostPad) GetTarget() *Pad {
	pad := C.gst_ghost_pad_get_target(g.Instance())
	if pad == nil {
		return nil
	}
	return FromGstPadUnsafeFull(unsafe.Pointer(pad))
}

// SetTarget sets the new target of the ghostpad gpad. Any existing target is unlinked and links to the new target are
// established. if newtarget is nil the target will be cleared.
func (g *GhostPad) SetTarget(target *Pad) bool {
	return gobool(C.gst_ghost_pad_set_target(
		g.Instance(), target.Instance(),
	))
}

// ActivateModeDefault invokes the default activate mode function of a ghost pad.
func (g *GhostPad) ActivateModeDefault(parent *Object, mode PadMode, active bool) bool {
	if parent == nil {
		return gobool(C.gst_ghost_pad_activate_mode_default(
			C.toGstPad(g.Unsafe()), nil, C.GstPadMode(mode), gboolean(active),
		))
	}
	return gobool(C.gst_ghost_pad_activate_mode_default(
		C.toGstPad(g.Unsafe()), parent.Instance(), C.GstPadMode(mode), gboolean(active),
	))
}

// InternalActivateModeDefault invokes the default activate mode function of a proxy pad that is owned by a ghost pad.
func (g *GhostPad) InternalActivateModeDefault(parent *Object, mode PadMode, active bool) bool {
	if parent == nil {
		return gobool(C.gst_ghost_pad_internal_activate_mode_default(
			C.toGstPad(g.Unsafe()), nil, C.GstPadMode(mode), gboolean(active),
		))
	}
	return gobool(C.gst_ghost_pad_internal_activate_mode_default(
		C.toGstPad(g.Unsafe()), parent.Instance(), C.GstPadMode(mode), gboolean(active),
	))
}

// ProxyPad is a go representation of a GstProxyPad.
type ProxyPad struct{ *Pad }

// toPad returns the underling GstPad of this ProxyPad.
func (p *ProxyPad) toPad() *C.GstPad { return C.toGstPad(p.Unsafe()) }

// Instance returns the underlying GstProxyPad instance.
func (p *ProxyPad) Instance() *C.GstProxyPad { return C.toGstProxyPad(p.Unsafe()) }

// GetInternal gets the internal pad of pad. Unref target pad after usage.
//
// The internal pad of a GhostPad is the internally used pad of opposite direction, which is used to link to the target.
func (p *ProxyPad) GetInternal() *ProxyPad {
	pad := C.gst_proxy_pad_get_internal(p.Instance())
	proxyPad := wrapProxyPad(toGObject(unsafe.Pointer(pad)))
	runtime.SetFinalizer(proxyPad, (*ProxyPad).Unref)
	return proxyPad
}

// ChainDefault invokes the default chain function of the proxy pad.
func (p *ProxyPad) ChainDefault(parent *Object, buffer *Buffer) FlowReturn {
	return FlowReturn(C.gst_proxy_pad_chain_default(p.toPad(), parent.Instance(), buffer.Instance()))
}

// ChainListDefault invokes the default chain list function of the proxy pad.
func (p *ProxyPad) ChainListDefault(parent *Object, bufferList *BufferList) FlowReturn {
	return FlowReturn(C.gst_proxy_pad_chain_list_default(p.toPad(), parent.Instance(), bufferList.Instance()))
}

// GetRangeDefault invokes the default getrange function of the proxy pad.
func (p *ProxyPad) GetRangeDefault(parent *Object, offset uint64, size uint) (FlowReturn, *Buffer) {
	var buf *C.GstBuffer
	ret := FlowReturn(C.gst_proxy_pad_getrange_default(p.toPad(), parent.Instance(), C.guint64(offset), C.guint(size), &buf))
	if ret != FlowError {
		return ret, FromGstBufferUnsafeFull(unsafe.Pointer(buf))
	}
	return ret, nil
}

// GetInternalLinksDefault invokes the default iterate internal links function of the proxy pad.
func (p *ProxyPad) GetInternalLinksDefault(parent *Object) ([]*Pad, error) {
	iterator := C.gst_proxy_pad_iterate_internal_links_default(p.toPad(), parent.Instance())
	return iteratorToPadSlice(iterator)
}
