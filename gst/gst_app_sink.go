package gst

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -Wno-unused-function -g
#include <gst/gst.h>
#include <gst/app/gstappsink.h>
#include "gst.go.h"
*/
import "C"

import (
	"bytes"
	"errors"
	"io"
	"unsafe"
)

// AppSink wraps an Element object with additional methods for pulling samples.
type AppSink struct{ *Element }

// NewAppSink returns a new appsink element. Unref after usage.
func NewAppSink() (*AppSink, error) {
	elem, err := NewElement("appsink")
	if err != nil {
		return nil, err
	}
	return wrapAppSink(elem), nil
}

// Instance returns the native GstAppSink instance.
func (a *AppSink) Instance() *C.GstAppSink { return C.toGstAppSink(a.unsafe()) }

// ErrEOS represents that the stream has ended.
var ErrEOS = errors.New("Pipeline has reached end-of-stream")

// IsEOS returns true if this AppSink has reached the end-of-stream.
func (a *AppSink) IsEOS() bool {
	return gobool(C.gst_app_sink_is_eos((*C.GstAppSink)(a.Instance())))
}

// BlockPullSample will block until a sample becomes available or the stream
// is ended.
func (a *AppSink) BlockPullSample() (*Sample, error) {
	for {
		if a.IsEOS() {
			return nil, ErrEOS
		}
		// This function won't block if the entire pipeline is waiting for data
		sample := C.gst_app_sink_pull_sample((*C.GstAppSink)(a.Instance()))
		if sample == nil {
			continue
		}
		return NewSample(sample), nil
	}
}

// PullSample will try to pull a sample or return nil if none is available.
func (a *AppSink) PullSample() (*Sample, error) {
	if a.IsEOS() {
		return nil, ErrEOS
	}
	sample := C.gst_app_sink_try_pull_sample(
		(*C.GstAppSink)(a.Instance()),
		C.GST_SECOND,
	)
	if sample != nil {
		return NewSample(sample), nil
	}
	return nil, nil
}

// Sample is a go wrapper around a GstSample object.
type Sample struct {
	sample *C.GstSample
}

// NewSample creates a new Sample from the given *GstSample.
func NewSample(sample *C.GstSample) *Sample { return &Sample{sample: sample} }

// Instance returns the underlying *GstSample instance.
func (s *Sample) Instance() *C.GstSample { return s.sample }

// Unref calls gst_sample_unref on the sample.
func (s *Sample) Unref() { C.gst_sample_unref((*C.GstSample)(s.Instance())) }

// GetBuffer returns a Reader for the buffer inside this sample.
func (s *Sample) GetBuffer() io.Reader {
	buffer := C.gst_sample_get_buffer((*C.GstSample)(s.Instance()))
	var mapInfo C.GstMapInfo
	C.gst_buffer_map(
		(*C.GstBuffer)(buffer),
		(*C.GstMapInfo)(unsafe.Pointer(&mapInfo)),
		C.GST_MAP_READ,
	)
	defer C.gst_buffer_unmap((*C.GstBuffer)(buffer), (*C.GstMapInfo)(unsafe.Pointer(&mapInfo)))
	return bytes.NewBuffer(C.GoBytes(unsafe.Pointer(mapInfo.data), (C.int)(mapInfo.size)))
}

func wrapAppSink(elem *Element) *AppSink { return &AppSink{elem} }
