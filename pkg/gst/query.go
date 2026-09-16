package gst

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/pkg/gobject/v2"
)

var (
	TypeQuery = gobject.Type(C.gst_query_get_type())
)

var _ gobject.GoValueInitializer = (*Query)(nil)

func marshalQuery(p unsafe.Pointer) (interface{}, error) {
	b := gobject.ValueFromNative(p).Boxed()
	return UnsafeQueryFromGlibBorrow(b), nil
}

func (r *Query) GoValueType() gobject.Type {
	return TypeQuery
}

func (r *Query) SetGoValue(v *gobject.Value) {
	v.SetBoxed(unsafe.Pointer(r.instance()))
}

func init() {
	gobject.RegisterGValueMarshalers([]gobject.TypeMarshaler{
		{T: TypeQuery, F: marshalQuery},
	})
}

func (q *Query) Type() QueryType {
	return QueryType(q.native._type)
}
