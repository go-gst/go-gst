package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Meta is a go representation of GstMeta.
type Meta struct {
	ptr *C.GstMeta
}

// Instance returns the underlying GstMeta instance.
func (m *Meta) Instance() *C.GstMeta { return C.toGstMeta(unsafe.Pointer(m.ptr)) }

// Flags returns the flags on this Meta instance.
func (m *Meta) Flags() MetaFlags { return MetaFlags(m.Instance().flags) }

// Info returns the extra info with this metadata.
func (m *Meta) Info() *MetaInfo { return wrapMetaInfo(m.Instance().info) }

// MetaInfo is a go representation of GstMetaInfo
type MetaInfo struct {
	ptr *C.GstMetaInfo
}

// Instance returns the underlying GstMetaInfo instance.
func (m *MetaInfo) Instance() *C.GstMetaInfo { return m.ptr }

// API returns the tag identifying the metadata structure and api.
func (m *MetaInfo) API() glib.Type { return glib.Type(m.Instance().api) }

// Type returns the type identifying the implementor of the api.
func (m *MetaInfo) Type() glib.Type { return glib.Type(m.Instance()._type) }

// Size returns the size of the metadata.
func (m *MetaInfo) Size() int64 { return int64(m.Instance().size) }
