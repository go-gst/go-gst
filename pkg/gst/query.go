package gst

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

func (q *Query) Type() QueryType {
	return QueryType(q.native._type)
}
