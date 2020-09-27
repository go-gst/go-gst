package gst

// #include "gst.go.h"
import "C"
import "unsafe"

// Sample is a go wrapper around a GstSample object.
type Sample struct {
	sample *C.GstSample
}

// FromGstSampleUnsafe wraps the pointer to the given C GstSample with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstSampleUnsafe(sample unsafe.Pointer) *Sample { return wrapSample(C.toGstSample(sample)) }

// Instance returns the underlying *GstSample instance.
func (s *Sample) Instance() *C.GstSample { return C.toGstSample(unsafe.Pointer(s.sample)) }

// Unref calls gst_sample_unref on the sample.
func (s *Sample) Unref() { C.gst_sample_unref((*C.GstSample)(s.Instance())) }

// GetBuffer returns the buffer inside this sample.
func (s *Sample) GetBuffer() *Buffer {
	return wrapBuffer(C.gst_sample_get_buffer((*C.GstSample)(s.Instance())))
}
