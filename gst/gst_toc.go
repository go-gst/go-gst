package gst

// #include "gst.go.h"
import "C"

// ToC is a Go representation of a GstToc.
type ToC struct {
	ptr *C.GstToc
}

// Instance returns the underlying GstToc instance.
func (t *ToC) Instance() *C.GstToc { return t.ptr }
