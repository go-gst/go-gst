package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// MiniObject is an opaque struct meant to form the base of gstreamer
// classes extending the GstMiniObject.
// This object is a WIP and is intended primarily for forming the base of extending classes.
type MiniObject struct {
	ptr    unsafe.Pointer
	parent *MiniObject
}

// NewMiniObject initializes a new mini object with the desired flags, types, and callbacks.
// If you don't need any callbacks you can specify nil.
// TODO: This is more for reference and is not fully implemented.
func NewMiniObject(flags MiniObjectFlags, gtype glib.Type) *MiniObject {
	var cMiniObj C.GstMiniObject
	C.gst_mini_object_init(
		C.toGstMiniObject(unsafe.Pointer(&cMiniObj)),
		C.uint(flags),
		C.gsize(gtype),
		nil, nil, nil,
	)
	return wrapMiniObject(&cMiniObj)
}

// native returns the pointer to the underlying object.
func (m *MiniObject) unsafe() unsafe.Pointer { return m.ptr }

// Parent returns the parent of this MiniObject
func (m *MiniObject) Parent() *MiniObject { return m.parent }

// Instance returns the native GstMiniObject instance.
func (m *MiniObject) Instance() *C.GstMiniObject { return C.toGstMiniObject(m.unsafe()) }

// Ref increases the ref count on this object by one.
func (m *MiniObject) Ref() { C.gst_mini_object_ref(m.Instance()) }

// Unref decresaes the ref count on this object by one.
func (m *MiniObject) Unref() { C.gst_mini_object_unref(m.Instance()) }

// AddParent adds the given object as a parent of this object.
// See https://gstreamer.freedesktop.org/documentation/gstreamer/gstminiobject.html?gi-language=c#gst_mini_object_add_parent.
func (m *MiniObject) AddParent(parent *MiniObject) {
	C.gst_mini_object_add_parent(m.Instance(), parent.Instance())
}

// Copy creates a copy of this object.
func (m *MiniObject) Copy() *MiniObject {
	return wrapMiniObject(C.gst_mini_object_copy(m.Instance()))
}

// Type returns the type of this mini object.
func (m *MiniObject) Type() glib.Type {
	return glib.Type(m.Instance()._type)
}

// GetData returns the userdata pointer associated with this object at the given key,
// or nil if none exists.
func (m *MiniObject) GetData(name string) unsafe.Pointer {
	data := C.gst_mini_object_get_qdata(m.Instance(), newQuarkFromString(name))
	if data == nil {
		return nil
	}
	return unsafe.Pointer(data)
}

// SetData sets a userdata pointer associated with this object at the given key,
// Use nil to delete an existing key.
func (m *MiniObject) SetData(name string, ptr unsafe.Pointer) {
	C.gst_mini_object_set_qdata(m.Instance(), newQuarkFromString(name), (C.gpointer)(ptr), nil)
}

func wrapMiniObject(p *C.GstMiniObject) *MiniObject {
	return &MiniObject{ptr: unsafe.Pointer(p)}
}
