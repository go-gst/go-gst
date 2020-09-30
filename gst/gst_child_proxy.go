package gst

// #include "gst.go.h"
import "C"

// ChildProxy is a go wrapper around a GstChildProxy.
type ChildProxy struct {
	ptr *C.GstChildProxy
}
