package app

// #include "gst.go.h"
import "C"
import (
	"errors"
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// Sink wraps an Element made with the appsink plugin with additional methods for pulling samples.
type Sink struct{ *gst.Element }

// NewAppSink returns a new appsink element. Unref after usage.
func NewAppSink() (*Sink, error) {
	elem, err := gst.NewElement("appsink")
	if err != nil {
		return nil, err
	}
	return wrapAppSink(elem), nil
}

// Instance returns the native GstAppSink instance.
func (a *Sink) Instance() *C.GstAppSink { return C.toGstAppSink(a.Unsafe()) }

// ErrEOS represents that the stream has ended.
var ErrEOS = errors.New("Pipeline has reached end-of-stream")

// IsEOS returns true if this AppSink has reached the end-of-stream.
func (a *Sink) IsEOS() bool {
	return gobool(C.gst_app_sink_is_eos((*C.GstAppSink)(a.Instance())))
}

// BlockPullSample will block until a sample becomes available or the stream
// is ended.
func (a *Sink) BlockPullSample() (*gst.Sample, error) {
	for {
		if a.IsEOS() {
			return nil, ErrEOS
		}
		// This function won't block if the entire pipeline is waiting for data
		sample := C.gst_app_sink_pull_sample((*C.GstAppSink)(a.Instance()))
		if sample == nil {
			continue
		}
		return gst.FromGstSampleUnsafe(unsafe.Pointer(sample)), nil
	}
}

// PullSample will try to pull a sample or return nil if none is available.
func (a *Sink) PullSample() (*gst.Sample, error) {
	if a.IsEOS() {
		return nil, ErrEOS
	}
	sample := C.gst_app_sink_try_pull_sample(
		(*C.GstAppSink)(a.Instance()),
		C.GST_SECOND,
	)
	if sample != nil {
		return gst.FromGstSampleUnsafe(unsafe.Pointer(sample)), nil
	}
	return nil, nil
}
