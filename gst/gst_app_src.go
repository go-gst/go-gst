package gst

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include <gst/app/gstappsrc.h>
#include "gst.go.h"
*/
import "C"
import (
	"io"
	"io/ioutil"
	"time"
	"unsafe"
)

// AppSrc wraps an Element object with additional methods for pushing samples.
type AppSrc struct{ *Element }

// NewAppSrc returns a new AppSrc element.
func NewAppSrc() (*AppSrc, error) {
	elem, err := NewElement("appsrc")
	if err != nil {
		return nil, err
	}
	return wrapAppSrc(elem), nil
}

// Instance returns the native GstAppSink instance.
func (a *AppSrc) Instance() *C.GstAppSrc { return C.toGstAppSrc(a.unsafe()) }

// SetSize sets the size of the source stream in bytes. You should call this for
// streams of fixed length.
func (a *AppSrc) SetSize(size int64) {
	C.gst_app_src_set_size((*C.GstAppSrc)(a.Instance()), (C.gint64)(size))
}

// SetDuration sets the duration of the source stream. You should call
// this if the value is known.
func (a *AppSrc) SetDuration(dur time.Duration) {
	C.gst_app_src_set_duration((*C.GstAppSrc)(a.Instance()), (C.ulong)(dur.Nanoseconds()))
}

// EndStream signals to the app source that the stream has ended after the last queued
// buffer.
func (a *AppSrc) EndStream() FlowReturn {
	ret := C.gst_app_src_end_of_stream((*C.GstAppSrc)(a.Instance()))
	return FlowReturn(ret)
}

// SetLive sets whether or not this is a live stream.
func (a *AppSrc) SetLive(b bool) error { return a.Set("is-live", b) }

// PushBuffer pushes a buffer to the appsrc. Currently by default this will block
// until the source is ready to accept the buffer.
func (a *AppSrc) PushBuffer(data io.Reader) FlowReturn {
	out, err := ioutil.ReadAll(data)
	if err != nil {
		return FlowError
	}
	str := string(out)
	p := unsafe.Pointer(C.CString(str))
	defer C.free(p)
	buf := C.gst_buffer_new_wrapped((C.gpointer)(p), C.ulong(len(str)))
	ret := C.gst_app_src_push_buffer((*C.GstAppSrc)(a.Instance()), (*C.GstBuffer)(buf))
	return FlowReturn(ret)
}

func wrapAppSrc(elem *Element) *AppSrc { return &AppSrc{elem} }
