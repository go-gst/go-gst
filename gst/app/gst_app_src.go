package app

// #include "gst.go.h"
import "C"
import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// Source wraps an Element made with the appsrc plugin with additional methods for pushing samples.
type Source struct{ *gst.Element }

// NewAppSrc returns a new AppSrc element.
func NewAppSrc() (*Source, error) {
	elem, err := gst.NewElement("appsrc")
	if err != nil {
		return nil, err
	}
	return wrapAppSrc(elem), nil
}

// Instance returns the native GstAppSink instance.
func (a *Source) Instance() *C.GstAppSrc { return C.toGstAppSrc(a.Unsafe()) }

// SetSize sets the size of the source stream in bytes. You should call this for
// streams of fixed length.
func (a *Source) SetSize(size int64) {
	C.gst_app_src_set_size((*C.GstAppSrc)(a.Instance()), (C.gint64)(size))
}

// SetDuration sets the duration of the source stream. You should call
// this if the value is known.
func (a *Source) SetDuration(dur time.Duration) {
	C.gst_app_src_set_duration((*C.GstAppSrc)(a.Instance()), (C.ulong)(dur.Nanoseconds()))
}

// EndStream signals to the app source that the stream has ended after the last queued
// buffer.
func (a *Source) EndStream() gst.FlowReturn {
	ret := C.gst_app_src_end_of_stream((*C.GstAppSrc)(a.Instance()))
	return gst.FlowReturn(ret)
}

// SetLive sets whether or not this is a live stream.
func (a *Source) SetLive(b bool) error { return a.Set("is-live", b) }

// PushBuffer pushes a buffer to the appsrc. Currently by default this will block
// until the source is ready to accept the buffer.
func (a *Source) PushBuffer(buf *gst.Buffer) gst.FlowReturn {
	ret := C.gst_app_src_push_buffer(
		(*C.GstAppSrc)(a.Instance()),
		(*C.GstBuffer)(unsafe.Pointer(buf.Instance())),
	)
	return gst.FlowReturn(ret)
}
