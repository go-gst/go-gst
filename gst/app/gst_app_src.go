package app

/*
#include "gst.go.h"

extern void  goAppGDestroyNotifyFunc (gpointer user_data);

extern void      goNeedDataCb   (GstAppSrc *src, guint length, gpointer user_data);
extern void      goEnoughDataDb (GstAppSrc *src, gpointer user_data);
extern gboolean  goSeekDataCb   (GstAppSrc *src, guint64 offset, gpointer user_data);

void          cgoSrcGDestroyNotifyFunc (gpointer user_data) { goAppGDestroyNotifyFunc(user_data); }

void      cgoNeedDataCb   (GstAppSrc *src, guint length, gpointer user_data) { goNeedDataCb(src, length, user_data); }
void      cgoEnoughDataCb (GstAppSrc *src, gpointer user_data) { goEnoughDataDb(src, user_data); }
gboolean  cgoSeekDataCb   (GstAppSrc *src, guint64 offset, gpointer user_data) { return goSeekDataCb(src, offset, user_data); }

*/
import "C"
import (
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
)

// SourceCallbacks represents callbacks to configure on an AppSource.
type SourceCallbacks struct {
	NeedDataFunc   func(src *Source, length uint)
	EnoughDataFunc func(src *Source)
	SeekDataFunc   func(src *Source, offset uint64) bool
}

// StreamType casts GstAppStreamType
type StreamType int

// Type castings
const (
	AppStreamTypeStream       StreamType = C.GST_APP_STREAM_TYPE_STREAM        // (0) – No seeking is supported in the stream, such as a live stream.
	AppStreamTypeSeekable     StreamType = C.GST_APP_STREAM_TYPE_SEEKABLE      // (1) – The stream is seekable but seeking might not be very fast, such as data from a webserver.
	AppStreamTypeRandomAccess StreamType = C.GST_APP_STREAM_TYPE_RANDOM_ACCESS // (2) – The stream is seekable and seeking is fast, such as in a local file.
)

// Source wraps an Element made with the appsrc plugin with additional methods for pushing samples.
type Source struct{ *base.GstBaseSrc }

// NewAppSrc returns a new AppSrc element.
func NewAppSrc() (*Source, error) {
	elem, err := gst.NewElement("appsrc")
	if err != nil {
		return nil, err
	}
	return wrapAppSrc(elem), nil
}

// SrcFromElement checks if the given element is an appsrc and if so returns
// a Source interace.
func SrcFromElement(elem *gst.Element) *Source {
	if appSrc := C.toGstAppSrc(elem.Unsafe()); appSrc != nil {
		return wrapAppSrc(elem)
	}
	return nil
}

// Instance returns the native GstAppSink instance.
func (a *Source) Instance() *C.GstAppSrc { return C.toGstAppSrc(a.Unsafe()) }

// EndStream signals to the app source that the stream has ended after the last queued
// buffer.
func (a *Source) EndStream() gst.FlowReturn {
	ret := C.gst_app_src_end_of_stream((*C.GstAppSrc)(a.Instance()))
	return gst.FlowReturn(ret)
}

// GetCaps gets the configures caps on the app src.
func (a *Source) GetCaps() *gst.Caps {
	caps := C.gst_app_src_get_caps(a.Instance())
	if caps == nil {
		return nil
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// GetCurrentLevelBytes gets the number of currently queued bytes inside appsrc.
func (a *Source) GetCurrentLevelBytes() uint64 {
	return uint64(C.gst_app_src_get_current_level_bytes(a.Instance()))
}

var gstClockTimeNone C.GstClockTime = 0xffffffffffffffff

// GetDuration gets the duration of the stream in nanoseconds. A negative value means that the duration is not known.
func (a *Source) GetDuration() time.Duration {
	dur := C.gst_app_src_get_duration(a.Instance())
	if dur == gstClockTimeNone {
		return gst.ClockTimeNone
	}
	return time.Duration(uint64(dur)) * time.Nanosecond
}

// GetEmitSignals checks if appsrc will emit the "new-preroll" and "new-buffer" signals.
func (a *Source) GetEmitSignals() bool {
	return gobool(C.gst_app_src_get_emit_signals(a.Instance()))
}

// GetLatency retrieves the min and max latencies in min and max respectively.
func (a *Source) GetLatency() (min, max uint64) {
	var gmin, gmax C.guint64
	C.gst_app_src_get_latency(a.Instance(), &gmin, &gmax)
	return uint64(gmin), uint64(gmax)
}

// GetMaxBytes gets the maximum amount of bytes that can be queued in appsrc.
func (a *Source) GetMaxBytes() uint64 {
	return uint64(C.gst_app_src_get_max_bytes(a.Instance()))
}

// GetSize gets the size of the stream in bytes. A value of -1 means that the size is not known.
func (a *Source) GetSize() int64 {
	return int64(C.gst_app_src_get_size(a.Instance()))
}

// GetStreamType gets the stream type. Control the stream type of appsrc with SetStreamType.
func (a *Source) GetStreamType() StreamType {
	return StreamType(C.gst_app_src_get_stream_type(a.Instance()))
}

// PushBuffer pushes a buffer to the appsrc. Currently by default this will block
// until the source is ready to accept the buffer.
func (a *Source) PushBuffer(buf *gst.Buffer) gst.FlowReturn {
	ret := C.gst_app_src_push_buffer(
		(*C.GstAppSrc)(a.Instance()),
		(*C.GstBuffer)(unsafe.Pointer(buf.Ref().Instance())),
	)
	return gst.FlowReturn(ret)
}

// PushBufferList adds a buffer list to the queue of buffers and buffer lists that the appsrc element will push
// to its source pad. This function takes ownership of buffer_list.
//
// When the block property is TRUE, this function can block until free space becomes available in the queue.
func (a *Source) PushBufferList(bufList *gst.BufferList) gst.FlowReturn {
	return gst.FlowReturn(C.gst_app_src_push_buffer_list(
		a.Instance(), (*C.GstBufferList)(unsafe.Pointer(bufList.Ref().Instance())),
	))
}

// PushSample Extract a buffer from the provided sample and adds it to the queue of buffers that the appsrc element will
// push to its source pad. Any previous caps that were set on appsrc will be replaced by the caps associated with the
// sample if not equal.
//
// This function does not take ownership of the sample so the sample needs to be unreffed after calling this function.
//
// When the block property is TRUE, this function can block until free space becomes available in the queue.
func (a *Source) PushSample(sample *gst.Sample) gst.FlowReturn {
	return gst.FlowReturn(C.gst_app_src_push_sample(
		a.Instance(), (*C.GstSample)(unsafe.Pointer(sample.Instance())),
	))
}

// SetCallbacks sets callbacks which will be executed when data is needed, enough data has been collected or when a seek
// should be performed. This is an alternative to using the signals, it has lower overhead and is thus less expensive,
// but also less flexible.
//
// If callbacks are installed, no signals will be emitted for performance reasons.
//
// Before 1.16.3 it was not possible to change the callbacks in a thread-safe way.
func (a *Source) SetCallbacks(cbs *SourceCallbacks) {
	ptr := gopointer.Save(cbs)
	appSrcCallbacks := &C.GstAppSrcCallbacks{
		need_data:   (*[0]byte)(unsafe.Pointer(C.cgoNeedDataCb)),
		enough_data: (*[0]byte)(unsafe.Pointer(C.cgoEnoughDataCb)),
		seek_data:   (*[0]byte)(unsafe.Pointer(C.cgoSeekDataCb)),
	}
	C.gst_app_src_set_callbacks(
		a.Instance(),
		appSrcCallbacks,
		(C.gpointer)(unsafe.Pointer(ptr)),
		C.GDestroyNotify(C.cgoSrcGDestroyNotifyFunc),
	)
}

// SetCaps sets the capabilities on the appsrc element. This function takes a copy of the caps structure. After calling this method,
// the source will only produce caps that match caps. caps must be fixed and the caps on the buffers must match the caps or left NULL.
func (a *Source) SetCaps(caps *gst.Caps) {
	C.gst_app_src_set_caps(a.Instance(), (*C.GstCaps)(unsafe.Pointer(caps.Instance())))
}

// SetDuration sets the duration of the source stream. You should call
// this if the value is known.
func (a *Source) SetDuration(dur time.Duration) {
	C.gst_app_src_set_duration((*C.GstAppSrc)(a.Instance()), C.GstClockTime(dur.Nanoseconds()))
}

// SetEmitSignals makes appsrc emit the "new-preroll" and "new-buffer" signals. This option is by default disabled because signal emission
// is expensive and unneeded when the application prefers to operate in pull mode.
func (a *Source) SetEmitSignals(emit bool) {
	C.gst_app_src_set_emit_signals(a.Instance(), gboolean(emit))
}

// SetLatency configures the min and max latency in src. If min is set to -1, the default latency calculations for pseudo-live sources
// will be used.
func (a *Source) SetLatency(min, max uint64) {
	C.gst_app_src_set_latency(a.Instance(), C.guint64(min), C.guint64(max))
}

// SetMaxBytes sets the maximum amount of bytes that can be queued in appsrc. After the maximum amount of bytes are queued, appsrc will
// emit the "enough-data" signal.
func (a *Source) SetMaxBytes(max uint64) {
	C.gst_app_src_set_max_bytes(a.Instance(), C.guint64(max))
}

// SetSize sets the size of the source stream in bytes. You should call this for
// streams of fixed length.
func (a *Source) SetSize(size int64) {
	C.gst_app_src_set_size((*C.GstAppSrc)(a.Instance()), (C.gint64)(size))
}

// SetStreamType sets the stream type on appsrc. For seekable streams, the "seek" signal must be connected to.
func (a *Source) SetStreamType(streamType StreamType) {
	C.gst_app_src_set_stream_type(a.Instance(), C.GstAppStreamType(streamType))
}
