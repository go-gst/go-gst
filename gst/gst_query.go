package gst

// #include "gst.go.h"
import "C"
import "unsafe"

// Query is a go wrapper around a GstQuery.
type Query struct {
	ptr *C.GstQuery
}

// Instance returns the underlying GstQuery instance.
func (q *Query) Instance() *C.GstQuery { return C.toGstQuery(unsafe.Pointer(q.ptr)) }
