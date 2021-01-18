package app

/*
#include "gst.go.h"

extern void          goAppGDestroyNotifyFunc (gpointer user_data);

extern void          goSinkEOSCb        (GstAppSink * sink, gpointer user_data);
extern GstFlowReturn goSinkNewPrerollCb (GstAppSink * sink, gpointer user_data);
extern GstFlowReturn goSinkNewSampleCb  (GstAppSink * sink, gpointer user_data);

void          cgoSinkGDestroyNotifyFunc (gpointer user_data) { goAppGDestroyNotifyFunc(user_data); }
void          cgoSinkEOSCb        (GstAppSink * sink, gpointer user_data) { return goSinkEOSCb(sink, user_data); }
GstFlowReturn cgoSinkNewPrerollCb (GstAppSink * sink, gpointer user_data) { return goSinkNewPrerollCb(sink, user_data); }
GstFlowReturn cgoSinkNewSampleCb  (GstAppSink * sink, gpointer user_data) { return goSinkNewSampleCb(sink, user_data); }

*/
import "C"

import (
	"errors"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"

	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
)

// SinkCallbacks represents callbacks that can be installed on an app sink when data is available.
type SinkCallbacks struct {
	EOSFunc        func(appSink *Sink)
	NewPrerollFunc func(appSink *Sink) gst.FlowReturn
	NewSampleFunc  func(appSink *Sink) gst.FlowReturn
}

// ErrEOS represents that the stream has ended.
var ErrEOS = errors.New("Pipeline has reached end-of-stream")

// Sink wraps an Element made with the appsink plugin with additional methods for pulling samples.
type Sink struct{ *base.GstBaseSink }

// NewAppSink returns a new appsink element. Unref after usage.
func NewAppSink() (*Sink, error) {
	elem, err := gst.NewElement("appsink")
	if err != nil {
		return nil, err
	}
	return wrapAppSink(elem), nil
}

// SinkFromElement checks if the given element is an appsink and if so returns
// a Sink interace.
func SinkFromElement(elem *gst.Element) *Sink {
	if appSink := C.toGstAppSink(elem.Unsafe()); appSink != nil {
		return wrapAppSink(elem)
	}
	return nil
}

// Instance returns the native GstAppSink instance.
func (a *Sink) Instance() *C.GstAppSink { return C.toGstAppSink(a.Unsafe()) }

// GetBufferListSupport checks if appsink supports buffer lists.
func (a *Sink) GetBufferListSupport() bool {
	return gobool(C.gst_app_sink_get_buffer_list_support(a.Instance()))
}

// GetCaps gets the configured caps on appsink.
func (a *Sink) GetCaps() *gst.Caps {
	caps := C.gst_app_sink_get_caps(a.Instance())
	if caps == nil {
		return nil
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// GetDrop checks if appsink will drop old buffers when the maximum amount of queued buffers is reached.
func (a *Sink) GetDrop() bool {
	return gobool(C.gst_app_sink_get_drop(a.Instance()))
}

// GetEmitSignals checks if appsink will emit the "new-preroll" and "new-sample" signals.
func (a *Sink) GetEmitSignals() bool {
	return gobool(C.gst_app_sink_get_emit_signals(a.Instance()))
}

// GetMaxBuffers gets the maximum amount of buffers that can be queued in appsink.
func (a *Sink) GetMaxBuffers() uint {
	return uint(C.gst_app_sink_get_max_buffers(a.Instance()))
}

// GetWaitOnEOS checks if appsink will wait for all buffers to be consumed when an EOS is received.
func (a *Sink) GetWaitOnEOS() bool {
	return gobool(C.gst_app_sink_get_wait_on_eos(a.Instance()))
}

// IsEOS returns true if this AppSink has reached the end-of-stream.
func (a *Sink) IsEOS() bool {
	return gobool(C.gst_app_sink_is_eos((*C.GstAppSink)(a.Instance())))
}

// PullPreroll gets the last preroll sample in appsink. This was the sample that caused the appsink to preroll in the PAUSED state.
//
// This function is typically used when dealing with a pipeline in the PAUSED state. Calling this function after doing a seek will
// give the sample right after the seek position.
//
// Calling this function will clear the internal reference to the preroll buffer.
//
// Note that the preroll sample will also be returned as the first sample when calling gst_app_sink_pull_sample.
//
// If an EOS event was received before any buffers, this function returns NULL. Use gst_app_sink_is_eos () to check for the EOS condition.
//
// This function blocks until a preroll sample or EOS is received or the appsink element is set to the READY/NULL state.
func (a *Sink) PullPreroll() *gst.Sample {
	smpl := C.gst_app_sink_pull_preroll(a.Instance())
	if smpl == nil {
		return nil
	}
	return gst.FromGstSampleUnsafeFull(unsafe.Pointer(smpl))
}

// PullSample blocks until a sample or EOS becomes available or the appsink element is set to the READY/NULL state.
//
// This function will only return samples when the appsink is in the PLAYING state. All rendered buffers will be put in a queue
// so that the application can pull samples at its own rate. Note that when the application does not pull samples fast enough, the queued
// buffers could consume a lot of memory, especially when dealing with raw video frames.
//
// If an EOS event was received before any buffers, this function returns NULL. Use IsEOS() to check for the EOS condition.
func (a *Sink) PullSample() *gst.Sample {
	smpl := C.gst_app_sink_pull_sample(a.Instance())
	if smpl == nil {
		return nil
	}
	return gst.FromGstSampleUnsafeFull(unsafe.Pointer(smpl))
}

// SetBufferListSupport instructs appsink to enable or disable buffer list support.
//
// For backwards-compatibility reasons applications need to opt in to indicate that they will be able to handle buffer lists.
func (a *Sink) SetBufferListSupport(enabled bool) {
	C.gst_app_sink_set_buffer_list_support(a.Instance(), gboolean(enabled))
}

// SetCallbacks sets callbacks which will be executed for each new preroll, new sample and eos. This is an alternative to using the signals,
// it has lower overhead and is thus less expensive, but also less flexible.
//
// If callbacks are installed, no signals will be emitted for performance reasons.
//
// Before 1.16.3 it was not possible to change the callbacks in a thread-safe way.
func (a *Sink) SetCallbacks(cbs *SinkCallbacks) {
	ptr := gopointer.Save(cbs)
	appSinkCallbacks := &C.GstAppSinkCallbacks{
		eos:         (*[0]byte)(unsafe.Pointer(C.cgoSinkEOSCb)),
		new_preroll: (*[0]byte)(unsafe.Pointer(C.cgoSinkNewPrerollCb)),
		new_sample:  (*[0]byte)(unsafe.Pointer(C.cgoSinkNewSampleCb)),
	}
	C.gst_app_sink_set_callbacks(
		a.Instance(),
		appSinkCallbacks,
		(C.gpointer)(unsafe.Pointer(ptr)),
		C.GDestroyNotify(C.cgoSinkGDestroyNotifyFunc),
	)
}

// SetCaps sets the capabilities on the appsink element. This function takes a copy of the caps structure. After calling this method,
// the sink will only accept caps that match caps. If caps is non-fixed, or incomplete, you must check the caps on the samples to get
// the actual used caps.
func (a *Sink) SetCaps(caps *gst.Caps) {
	C.gst_app_sink_set_caps(a.Instance(), (*C.GstCaps)(unsafe.Pointer(caps.Instance())))
}

// SetDrop instructs appsink to drop old buffers when the maximum amount of queued buffers is reached.
func (a *Sink) SetDrop(drop bool) {
	C.gst_app_sink_set_drop(a.Instance(), gboolean(drop))
}

// SetEmitSignals makes appsink emit the "new-preroll" and "new-sample" signals. This option is by default disabled because signal emission
// is expensive and unneeded when the application prefers to operate in pull mode.
func (a *Sink) SetEmitSignals(emit bool) {
	C.gst_app_sink_set_emit_signals(a.Instance(), gboolean(emit))
}

// SetMaxBuffers sets the maximum amount of buffers that can be queued in appsink. After this amount of buffers are queued in appsink,
// any more buffers will block upstream elements until a sample is pulled from appsink.
func (a *Sink) SetMaxBuffers(max uint) {
	C.gst_app_sink_set_max_buffers(a.Instance(), C.guint(max))
}

// SetWaitOnEOS instructs appsink to wait for all buffers to be consumed when an EOS is received.
func (a *Sink) SetWaitOnEOS(wait bool) {
	C.gst_app_sink_set_wait_on_eos(a.Instance(), gboolean(wait))
}

// TryPullPreroll gets the last preroll sample in appsink. This was the sample that caused the appsink to preroll in the PAUSED state.
//
// This function is typically used when dealing with a pipeline in the PAUSED state. Calling this function after doing a seek will give
// the sample right after the seek position.
//
// Calling this function will clear the internal reference to the preroll buffer.
//
// Note that the preroll sample will also be returned as the first sample when calling PullSample.
//
// If an EOS event was received before any buffers or the timeout expires, this function returns NULL. Use IsEOS () to check for the EOS condition.
//
// This function blocks until a preroll sample or EOS is received, the appsink element is set to the READY/NULL state, or the timeout expires.
func (a *Sink) TryPullPreroll(timeout time.Duration) *gst.Sample {
	tm := C.GstClockTime(timeout.Nanoseconds())
	smpl := C.gst_app_sink_try_pull_preroll(a.Instance(), tm)
	if smpl == nil {
		return nil
	}
	return gst.FromGstSampleUnsafeFull(unsafe.Pointer(smpl))
}

// TryPullSample blocks until a sample or EOS becomes available or the appsink element is set to the READY/NULL state or the timeout expires.
//
// This function will only return samples when the appsink is in the PLAYING state. All rendered buffers will be put in a queue so that the
// application can pull samples at its own rate. Note that when the application does not pull samples fast enough, the queued buffers could
// consume a lot of memory, especially when dealing with raw video frames.
//
// If an EOS event was received before any buffers or the timeout expires, this function returns NULL. Use IsEOS () to check for the EOS condition.
func (a *Sink) TryPullSample(timeout time.Duration) *gst.Sample {
	tm := C.GstClockTime(timeout.Nanoseconds())
	smpl := C.gst_app_sink_try_pull_sample(a.Instance(), tm)
	if smpl == nil {
		return nil
	}
	return gst.FromGstSampleUnsafeFull(unsafe.Pointer(smpl))
}
