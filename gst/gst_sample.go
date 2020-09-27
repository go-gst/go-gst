package gst

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -Wno-unused-function -g
#include <gst/gst.h>
#include <gst/app/gstappsink.h>
#include "gst.go.h"
*/
import "C"

// Sample is a go wrapper around a GstSample object.
type Sample struct {
	sample *C.GstSample
}

// Instance returns the underlying *GstSample instance.
func (s *Sample) Instance() *C.GstSample { return s.sample }

// Unref calls gst_sample_unref on the sample.
func (s *Sample) Unref() { C.gst_sample_unref((*C.GstSample)(s.Instance())) }

// GetBuffer returns the buffer inside this sample.
func (s *Sample) GetBuffer() *Buffer {
	return wrapBuffer(C.gst_sample_get_buffer((*C.GstSample)(s.Instance())))
}
