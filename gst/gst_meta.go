package gst

/*
#include "gst.go.h"

extern void     goMetaFreeFunc       (GstMeta * meta, GstBuffer * buffer);
extern gboolean goMetaInitFunc       (GstMeta *meta, gpointer params, GstBuffer * buffer);
extern gboolean goMetaTransformFunc  (GstBuffer * transBuf, GstMeta * meta, GstBuffer * buffer, GQuark type, gpointer data);

void cgoMetaFreeFunc (GstMeta * meta, GstBuffer * buffer)
{
	goMetaFreeFunc(meta, buffer);
}

gboolean cgoMetaInitFunc (GstMeta * meta, gpointer params, GstBuffer * buffer)
{
	return goMetaInitFunc(meta, params, buffer);
}

gboolean cgoMetaTransformFunc (GstBuffer * transBuf, GstMeta * meta, GstBuffer * buffer, GQuark type, gpointer data)
{
	return goMetaTransformFunc(transBuf, meta, buffer, type, data);
}

*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Meta is a go representation of GstMeta.
type Meta struct {
	ptr *C.GstMeta
}

// FromGstMetaUnsafe wraps the pointer to the given C GstMeta with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstMetaUnsafe(ptr unsafe.Pointer) *Meta { return wrapMeta(C.toGstMeta(ptr)) }

// Instance returns the underlying GstMeta instance.
func (m *Meta) Instance() *C.GstMeta { return C.toGstMeta(unsafe.Pointer(m.ptr)) }

// Flags returns the flags on this Meta instance.
func (m *Meta) Flags() MetaFlags { return MetaFlags(m.Instance().flags) }

// SetFlags sets the flags on this Meta instance.
func (m *Meta) SetFlags(flags MetaFlags) { m.Instance().flags = C.GstMetaFlags(flags) }

// Info returns the extra info with this metadata.
func (m *Meta) Info() *MetaInfo { return wrapMetaInfo(m.Instance().info) }

// SetInfo sets the info on this metadata.
func (m *Meta) SetInfo(info *MetaInfo) { m.Instance().info = info.Instance() }

// MetaInfo is a go representation of GstMetaInfo
type MetaInfo struct {
	ptr *C.GstMetaInfo
}

// FromGstMetaInfoUnsafe wraps the given unsafe pointer in a MetaInfo instance.
func FromGstMetaInfoUnsafe(ptr unsafe.Pointer) *MetaInfo {
	if ptr == nil {
		return nil
	}
	return &MetaInfo{ptr: (*C.GstMetaInfo)(ptr)}
}

// RegisterAPIType registers and returns a GType for the given api name and associates it with tags.
func RegisterAPIType(name string, tags []string) glib.Type {
	cTags := gcharStrings(tags)
	defer C.g_free((C.gpointer)(unsafe.Pointer(cTags)))
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	newType := C.gst_meta_api_type_register((*C.gchar)(cName), cTags)
	return glib.Type(newType)
}

// GetAPIInfo gets the MetaInfo for the given api type.
func GetAPIInfo(name string) *MetaInfo {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return wrapMetaInfo(C.gst_meta_get_info((*C.gchar)(cName)))
}

// GetAPITags retrieves the tags for the given api type.
func GetAPITags(apiType glib.Type) []string {
	tags := C.gst_meta_api_type_get_tags(C.GType(apiType))
	return goStrings(C.sizeOfGCharArray(tags), tags)
}

// APIHasTag returns true if the given api has the given tag.
func APIHasTag(api glib.Type, tag string) bool {
	q := newQuarkFromString(tag)
	return gobool(C.gst_meta_api_type_has_tag(C.GType(api), q))
}

// MetaInitFunc is a function called when meta is initialized in buffer.
type MetaInitFunc func(params interface{}, buffer *Buffer) bool

// MetaFreeFunc is a function called when meta is freed in buffer.
type MetaFreeFunc func(buffer *Buffer)

// MetaTransformFunc is a function called for each meta in buf as a result
// of performing a transformation on transbuf. Additional type specific transform
// data is passed to the function as data.
type MetaTransformFunc func(transBuf, buf *Buffer, mType string, data *MetaTransformCopy) bool

// MetaTransformCopy is extra data passed to a MetaTransformFunc
type MetaTransformCopy struct {
	// true if only region is copied
	Region bool
	// the offset to copy, 0 if region is FALSE, otherwise > 0
	Offset int64
	// the size to copy, -1 or the buffer size when region is FALSE
	Size int64
}

// MetaInfoCallbackFuncs represents callback functions to includ when registering a new
// meta type.
type MetaInfoCallbackFuncs struct {
	InitFunc      MetaInitFunc
	FreeFunc      MetaFreeFunc
	TransformFunc MetaTransformFunc
}

// Register meta callbacks internally as well so we can track them. Not all
// GstMeta callbacks include userdata.
var registeredMetas = make(map[glib.Type]map[string]*MetaInfoCallbackFuncs)

// RegisterMeta registers and returns a new MetaInfo instance denoting the
// given type, name, and size.
func RegisterMeta(api glib.Type, name string, size int64, cbFuncs *MetaInfoCallbackFuncs) *MetaInfo {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	metaInfo := C.gst_meta_register(
		C.GType(api),
		(*C.gchar)(unsafe.Pointer(cName)),
		C.gsize(size),
		C.GstMetaInitFunction(C.cgoMetaInitFunc),
		C.GstMetaFreeFunction(C.cgoMetaFreeFunc),
		C.GstMetaTransformFunction(C.cgoMetaTransformFunc),
	)
	if metaInfo == nil {
		return nil
	}
	wrapped := wrapMetaInfo(metaInfo)
	if registeredMetas[api] == nil {
		registeredMetas[api] = make(map[string]*MetaInfoCallbackFuncs)
	}
	registeredMetas[api][name] = cbFuncs
	return wrapped
}

// Instance returns the underlying GstMetaInfo instance.
func (m *MetaInfo) Instance() *C.GstMetaInfo { return m.ptr }

// API returns the tag identifying the metadata structure and api.
func (m *MetaInfo) API() glib.Type { return glib.Type(m.Instance().api) }

// SetAPI sets the API tag identifying the metadata structure and api.
func (m *MetaInfo) SetAPI(t glib.Type) { m.Instance().api = C.GType(t) }

// Type returns the type identifying the implementor of the api.
func (m *MetaInfo) Type() glib.Type { return glib.Type(m.Instance()._type) }

// SetType sets the type identifying the implementor of the api.
func (m *MetaInfo) SetType(t glib.Type) { m.Instance()._type = C.GType(t) }

// Size returns the size of the metadata.
func (m *MetaInfo) Size() int64 { return int64(m.Instance().size) }

// SetSize sets the size on the metadata.
func (m *MetaInfo) SetSize(size int64) { m.Instance().size = C.gsize(size) }
