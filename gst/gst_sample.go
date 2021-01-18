package gst

// #include "gst.go.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// Sample is a go wrapper around a GstSample object.
type Sample struct {
	sample *C.GstSample
}

// FromGstSampleUnsafeNone wraps the pointer to the given C GstSample with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstSampleUnsafeNone(sample unsafe.Pointer) *Sample {
	s := wrapSample(C.toGstSample(sample))
	s.Ref()
	runtime.SetFinalizer(s, (*Sample).Unref)
	return s
}

// FromGstSampleUnsafeFull wraps the pointer to the given C GstSample with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstSampleUnsafeFull(sample unsafe.Pointer) *Sample {
	s := wrapSample(C.toGstSample(sample))
	runtime.SetFinalizer(s, (*Sample).Unref)
	return s
}

// Instance returns the underlying *GstSample instance.
func (s *Sample) Instance() *C.GstSample { return C.toGstSample(unsafe.Pointer(s.sample)) }

// Ref increases the ref count on the sample.
func (s *Sample) Ref() *Sample {
	return wrapSample(C.gst_sample_ref(s.Instance()))
}

// Copy creates a copy of the given sample. This will also make a newly allocated copy of the data
// the source sample contains.
func (s *Sample) Copy() *Sample {
	return FromGstSampleUnsafeFull(unsafe.Pointer(C.gst_sample_copy(s.Instance())))
}

// GetBuffer returns the buffer inside this sample.
func (s *Sample) GetBuffer() *Buffer {
	return FromGstBufferUnsafeNone(unsafe.Pointer(C.gst_sample_get_buffer((*C.GstSample)(s.Instance()))))
}

// GetBufferList gets the buffer list associated with this sample.
func (s *Sample) GetBufferList() *BufferList {
	return FromGstBufferListUnsafeNone(unsafe.Pointer(C.gst_sample_get_buffer_list(s.Instance())))
}

// GetCaps returns the caps associated with this sample. Take a ref if you need to hold on to them
// longer then the life of the sample.
func (s *Sample) GetCaps() *Caps {
	return FromGstCapsUnsafeNone(unsafe.Pointer(C.gst_sample_get_caps(s.Instance())))
}

// GetInfo gets extra information about this sample. The structure remains valid as long as sample is valid.
func (s *Sample) GetInfo() *Structure { return wrapStructure(C.gst_sample_get_info(s.Instance())) }

// GetSegment gets the segment associated with the sample. The segmenr remains valid as long as sample is valid.
func (s *Sample) GetSegment() *Segment { return wrapSegment(C.gst_sample_get_segment(s.Instance())) }

// SetBuffer sets the buffer inside this sample. The sample must be writable.
func (s *Sample) SetBuffer(buf *Buffer) { C.gst_sample_set_buffer(s.Instance(), buf.Instance()) }

// SetBufferList sets the buffer list for this sample. The sample must be writable.
func (s *Sample) SetBufferList(buf *BufferList) {
	C.gst_sample_set_buffer_list(s.Instance(), buf.Instance())
}

// SetCaps sets the caps on this sample. The sample must be writable.
func (s *Sample) SetCaps(caps *Caps) { C.gst_sample_set_caps(s.Instance(), caps.Instance()) }

// SetInfo sets the info on this sample. The sample must be writable.
func (s *Sample) SetInfo(st *Structure) bool {
	return gobool(C.gst_sample_set_info(s.Instance(), st.Instance()))
}

// SetSegment sets the segment on this sample. The sample must be writable.
func (s *Sample) SetSegment(segment *Segment) {
	C.gst_sample_set_segment(s.Instance(), segment.Instance())
}

// Unref calls gst_sample_unref on the sample.
func (s *Sample) Unref() { C.gst_sample_unref((*C.GstSample)(s.Instance())) }
