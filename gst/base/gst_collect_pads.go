package base

/*
#include "gst.go.h"

extern GstFlowReturn  goGstCollectPadsBufferFunc   (GstCollectPads * pads, GstCollectData * data, GstBuffer * buffer, gpointer user_data);
extern GstFlowReturn  goGstCollectPadsClipFunc     (GstCollectPads * pads, GstCollectData * data, GstBuffer * inbuffer, GstBuffer ** outbuffer, gpointer user_data);
extern gint           goGstCollectPadsCompareFunc  (GstCollectPads * pads, GstCollectData * data1, GstClockTime ts1, GstCollectData * data2, GstClockTime ts2, gpointer user_data);
extern gboolean       goGstCollectPadsEventFunc    (GstCollectPads * pads, GstCollectData * data, GstEvent * event, gpointer user_data);
extern void           goGstCollectPadsFlushFunc    (GstCollectPads * pads, gpointer user_data);
extern GstFlowReturn  goGstCollectPadsFunc         (GstCollectPads * pads, gpointer user_data);
extern gboolean       goGstCollectPadsQueryFunc    (GstCollectPads * pads, GstCollectData * data, GstQuery * query, gpointer user_data);

GstFlowReturn
cgoGstCollectPadsBufferFunc (GstCollectPads * pads, GstCollectData * data, GstBuffer * buffer, gpointer user_data)
{
	return goGstCollectPadsBufferFunc(pads, data, buffer, user_data);
}

GstFlowReturn
cgoGstCollectPadsClipFunc (GstCollectPads * pads, GstCollectData * data, GstBuffer * inbuffer, GstBuffer ** outbuffer, gpointer user_data)
{
	return goGstCollectPadsClipFunc(pads, data, inbuffer, outbuffer, user_data);
}

gint
cgoGstCollectPadsCompareFunc (GstCollectPads * pads, GstCollectData * data1, GstClockTime ts1, GstCollectData * data2, GstClockTime ts2, gpointer user_data)
{
	return goGstCollectPadsCompareFunc(pads, data1, ts1, data2, ts2, user_data);
}

gboolean
cgoGstCollectPadsEventFunc (GstCollectPads * pads, GstCollectData * data, GstEvent * event, gpointer user_data)
{
	return goGstCollectPadsEventFunc(pads, data, event, user_data);
}

void
cgoGstCollectPadsFlushFunc (GstCollectPads * pads, gpointer user_data)
{
	goGstCollectPadsFlushFunc(pads, user_data);
}

GstFlowReturn
cgoGstCollectPadsFunc (GstCollectPads * pads, gpointer user_data)
{
	return goGstCollectPadsFunc(pads, user_data);
}

gboolean
cgoGstCollectPadsQueryFunc (GstCollectPads * pads, GstCollectData * data, GstQuery * query, gpointer user_data)
{
	return goGstCollectPadsQueryFunc(pads, data, query, user_data);
}

*/
import "C"

import (
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"

	"github.com/tinyzimmer/go-gst/gst"
)

// CollectPadsBufferFunc is a function that will be called when a (considered oldest) buffer can be muxed.
// If all pads have reached EOS, this function is called with a nil data and buffer.
type CollectPadsBufferFunc func(self *CollectPads, data *CollectData, buf *gst.Buffer) gst.FlowReturn

// CollectPadsClipFunc is a function that will be called when a buffer is received on the pad managed by data
// in the collectpad object pads.
//
// The function should use the segment of data and the negotiated media type on the pad to perform clipping of
// the buffer.
//
// The function should return a nil buffer if it should be dropped. The bindings will take care of ownership
// of the in-buffer.
type CollectPadsClipFunc func(self *CollectPads, data *CollectData, inbuffer *gst.Buffer) (gst.FlowReturn, *gst.Buffer)

// CollectPadsCompareFunc is a function for comparing two timestamps of buffers or newsegments collected on
// one pad. The function should return an integer less than zero when first timestamp is deemed older than the
// second one. Zero if the timestamps are deemed equally old. Integer greater than zero when second timestamp
// is deemed older than the first one.
type CollectPadsCompareFunc func(self *CollectPads, data1 *CollectData, ts1 time.Duration, data2 *CollectData, ts2 time.Duration) int

// CollectPadsEventFunc is a function that will be called while processing an event. It takes ownership of the
// event and is responsible for chaining up (to EventDefault) or dropping events (such typical cases being handled
// by the default handler). It should return true if the pad could handle the event.
type CollectPadsEventFunc func(self *CollectPads, data *CollectData, event *gst.Event) bool

// CollectPadsFlushFunc is a function that will be called while processing a flushing seek event.
//
// The function should flush any internal state of the element and the state of all the pads. It should clear
// only the state not directly managed by the pads object. It is therefore not necessary to call SetFlushing()
// nor Clear() from this function.
type CollectPadsFlushFunc func(self *CollectPads)

// CollectPadsFunc is a function that will be called when all pads have received data.
type CollectPadsFunc func(self *CollectPads) gst.FlowReturn

// CollectPadsQueryFunc is a function that will be called while processing a query. It takes ownership of the
// query and is responsible for chaining up (to events downstream (with EventDefault()).
//
// The function should return true if the pad could handle the event.
type CollectPadsQueryFunc func(self *CollectPads, data *CollectData, query *gst.Query) bool

// CollectData is a structure used by CollectPads.
type CollectData struct{ ptr *C.GstCollectData }

func wrapCollectData(ptr *C.GstCollectData) *CollectData { return &CollectData{ptr} }

// Instance returns the underly C object
func (c *CollectData) Instance() *C.GstCollectData { return c.ptr }

// Collect returns the owner CollectPads
func (c *CollectData) Collect() *CollectPads { return wrapCollectPadsNone(c.ptr.collect) }

// Pad returns the pad managed by this data.
func (c *CollectData) Pad() *gst.Pad { return gst.FromGstPadUnsafeNone(unsafe.Pointer(c.ptr.pad)) }

// Buffer returns the currently queued buffer.
func (c *CollectData) Buffer() *gst.Buffer {
	return gst.FromGstBufferUnsafeNone(unsafe.Pointer(c.ptr.buffer))
}

// Pos returns the position in the buffer.
func (c *CollectData) Pos() uint { return uint(c.ptr.pos) }

// Segment returns the last segment received.
func (c *CollectData) Segment() *gst.Segment {
	return gst.FromGstSegmentUnsafe(unsafe.Pointer(&c.ptr.segment))
}

// DTS returns the signed version of the DTS converted to running time.
func (c *CollectData) DTS() time.Duration { return time.Duration(C.gstCollectDataDTS(c.ptr)) }

// CollectPads manages a set of pads that operate in collect mode. This means that control is given to the
// manager of this object when all pads have data.
// For more information see:
// https://gstreamer.freedesktop.org/documentation/base/gstcollectpads.html?gi-language=c#gstcollectpads-page
type CollectPads struct {
	*gst.Object
	funcMap *collectPadsFuncMap
	selfPtr unsafe.Pointer
}

type collectPadsFuncMap struct {
	bufferFunc  CollectPadsBufferFunc
	clipFunc    CollectPadsClipFunc
	compareFunc CollectPadsCompareFunc
	eventFunc   CollectPadsEventFunc
	flushFunc   CollectPadsFlushFunc
	funcFunc    CollectPadsFunc
	queryFunc   CollectPadsQueryFunc
}

// NewCollectPads creates a new CollectPads instance.
func NewCollectPads() *CollectPads {
	return wrapCollectPadsFull(C.gst_collect_pads_new())
}

func wrapCollectPadsFull(ptr *C.GstCollectPads) *CollectPads {
	collect := &CollectPads{
		Object:  gst.FromGstObjectUnsafeFull(unsafe.Pointer(ptr)),
		funcMap: &collectPadsFuncMap{},
	}
	collect.selfPtr = gopointer.Save(collect)
	return collect
}

func wrapCollectPadsNone(ptr *C.GstCollectPads) *CollectPads {
	collect := &CollectPads{
		Object:  gst.FromGstObjectUnsafeNone(unsafe.Pointer(ptr)),
		funcMap: &collectPadsFuncMap{},
	}
	collect.selfPtr = gopointer.Save(collect)
	return collect
}

// Instance returns the underlying C object.
func (c *CollectPads) Instance() *C.GstCollectPads { return C.toGstCollectPads(c.Unsafe()) }

// AddPad adds a pad to the collection of collect pads. The pad has to be a sinkpad. The refcount of the pad is
// incremented. Use RemovePad to remove the pad from the collection again.
//
// Keeping a pad locked in waiting state is only relevant when using the default collection algorithm
// (providing the oldest buffer). It ensures a buffer must be available on this pad for a collection to take
// place. This is of typical use to a muxer element where non-subtitle streams should always be in waiting
// state, e.g. to assure that caps information is available on all these streams when initial headers have
// to be written.
//
// The pad will be automatically activated in push mode when pads is started.
//
// This function can return nil if supplied with invalid arguments.
func (c *CollectPads) AddPad(pad *gst.Pad, lock bool) *CollectData {
	data := C.gst_collect_pads_add_pad(
		c.Instance(),
		(*C.GstPad)(pad.Unsafe()),
		C.guint(C.sizeof_GstCollectData),
		nil,
		gboolean(lock),
	)
	if data == nil {
		return nil
	}
	return wrapCollectData(data)
}

// Available queries how much bytes can be read from each queued buffer. This means that the result of
// this call is the maximum number of bytes that can be read from each of the pads.
//
// This function should be called with pads STREAM_LOCK held, such as in the callback.
func (c *CollectPads) Available() uint { return uint(C.gst_collect_pads_available(c.Instance())) }

// InvalidRunningTime is a cast of C_MININT64 to signify a DTS that is invalid.
var InvalidRunningTime = time.Duration(C.G_MININT64)

// ClipRunningTime is a convenience clipping function that converts incoming buffer's timestamp to running
// time, or clips the buffer if outside configured segment.
//
// Since 1.6, this clipping function also sets the DTS parameter of the GstCollectData structure. This version
// of the running time DTS can be negative. InvalidRunningTime is used to indicate invalid value.
//
// data is the CollectData of the cooresponding pad and buf is the buffer being clipped.
func (c *CollectPads) ClipRunningTime(data *CollectData, buf *gst.Buffer) (ret gst.FlowReturn, outbuf *gst.Buffer) {
	var goutbuf *C.GstBuffer
	ret = gst.FlowReturn(C.gst_collect_pads_clip_running_time(
		c.Instance(),
		data.Instance(),
		(*C.GstBuffer)(unsafe.Pointer(buf.Instance())),
		&goutbuf,
		nil,
	))
	if goutbuf != nil {
		outbuf = gst.FromGstBufferUnsafeFull(unsafe.Pointer(goutbuf))
	}
	return
}

// EventDefault is the default GstCollectPads event handling that elements should always chain up to to ensure
// proper operation. Element might however indicate event should not be forwarded downstream.
func (c *CollectPads) EventDefault(data *CollectData, event *gst.Event, discard bool) bool {
	return gobool(C.gst_collect_pads_event_default(
		c.Instance(),
		data.Instance(),
		(*C.GstEvent)(unsafe.Pointer(event.Instance())),
		gboolean(discard),
	))
}

// Flush size bytes from the pad data. Returns the number of bytes actually flushed.
//
// This function should be called with pads STREAM_LOCK held, such as in the callback.
func (c *CollectPads) Flush(data *CollectData, size uint) uint {
	return uint(C.gst_collect_pads_flush(c.Instance(), data.Instance(), C.guint(size)))
}

// Free will free the C references to any go callbacks registered with the CollectPads. This is required
// due to the way the bindings are implemented around this object currently. While it's safe to assume this
// data will be collected whenever a program exits, in the context of a plugin that might get reused in
// a single application, NOT calling this function between starts and stops of your element could lead to
// memory leaks.
func (c *CollectPads) Free() { gopointer.Unref(c.selfPtr) }

// Peek at the buffer currently queued in data. This function should be called with the pads STREAM_LOCK held,
// such as in the callback handler.
func (c *CollectPads) Peek(data *CollectData) *gst.Buffer {
	buf := C.gst_collect_pads_peek(c.Instance(), data.Instance())
	if buf == nil {
		return nil
	}
	return gst.FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

// Pop the buffer currently queued in data. This function should be called with the pads STREAM_LOCK held, such
// as in the callback handler.
func (c *CollectPads) Pop(data *CollectData) *gst.Buffer {
	buf := C.gst_collect_pads_pop(c.Instance(), data.Instance())
	if buf == nil {
		return nil
	}
	return gst.FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

// QueryDefault is the Default GstCollectPads query handling that elements should always chain up to to ensure
// proper operation. Element might however indicate query should not be forwarded downstream.
func (c *CollectPads) QueryDefault(data *CollectData, query *gst.Query, discard bool) bool {
	return gobool(C.gst_collect_pads_query_default(
		c.Instance(), data.Instance(),
		(*C.GstQuery)(unsafe.Pointer(query.Instance())),
		gboolean(discard),
	))
}

// ReadBuffer gets a subbuffer of size bytes from the given pad data.
//
// This function should be called with pads STREAM_LOCK held, such as in the callback.
func (c *CollectPads) ReadBuffer(data *CollectData, size uint) *gst.Buffer {
	buf := C.gst_collect_pads_read_buffer(c.Instance(), data.Instance(), C.guint(size))
	if buf == nil {
		return nil
	}
	return gst.FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

// RemovePad removes a pad from the collection of collect pads. This function will also free the GstCollectData
// and all the resources that were allocated with AddPad.
//
// The pad will be deactivated automatically when CollectPads is stopped.
func (c *CollectPads) RemovePad(pad *gst.Pad) bool {
	return gobool(C.gst_collect_pads_remove_pad(c.Instance(), (*C.GstPad)(unsafe.Pointer(pad.Instance()))))
}

// SetBufferFunction sets the callback that will be called with the oldest buffer when all pads have been collected,
// or nil on EOS.
func (c *CollectPads) SetBufferFunction(f CollectPadsBufferFunc) {
	c.funcMap.bufferFunc = f
	C.gst_collect_pads_set_buffer_function(
		c.Instance(),
		C.GstCollectPadsBufferFunction(C.cgoGstCollectPadsBufferFunc),
		(C.gpointer)(c.selfPtr),
	)
}

// SetClipFunction installs a clipping function that is called after buffers are received on managed pads.
// See CollectPadsClipFunc for more details.
func (c *CollectPads) SetClipFunction(f CollectPadsClipFunc) {
	c.funcMap.clipFunc = f
	C.gst_collect_pads_set_clip_function(
		c.Instance(),
		C.GstCollectPadsClipFunction(C.cgoGstCollectPadsClipFunc),
		(C.gpointer)(c.selfPtr),
	)
}

// SetCompareFunction sets the timestamp comparisson function.
func (c *CollectPads) SetCompareFunction(f CollectPadsCompareFunc) {
	c.funcMap.compareFunc = f
	C.gst_collect_pads_set_compare_function(
		c.Instance(),
		C.GstCollectPadsCompareFunction(C.cgoGstCollectPadsCompareFunc),
		(C.gpointer)(c.selfPtr),
	)
}

// SetEventFunction sets the event callback function that will be called when collectpads has received an event
// originating from one of the collected pads. If the event being processed is a serialized one, this callback
// is called with pads STREAM_LOCK held, otherwise not. As this lock should be held when calling a number of
// CollectPads functions, it should be acquired if so (unusually) needed.
func (c *CollectPads) SetEventFunction(f CollectPadsEventFunc) {
	c.funcMap.eventFunc = f
	C.gst_collect_pads_set_event_function(
		c.Instance(),
		C.GstCollectPadsEventFunction(C.cgoGstCollectPadsEventFunc),
		(C.gpointer)(c.selfPtr),
	)
}

// SetFlushFunction installs a flush function that is called when the internal state of all pads should be
// flushed as part of flushing seek handling. See CollectPadsFlushFunc for more info.
func (c *CollectPads) SetFlushFunction(f CollectPadsFlushFunc) {
	c.funcMap.flushFunc = f
	C.gst_collect_pads_set_flush_function(
		c.Instance(),
		C.GstCollectPadsFlushFunction(C.cgoGstCollectPadsFlushFunc),
		(C.gpointer)(c.selfPtr),
	)
}

// SetFlushing changes the flushing state of all the pads in the collection. No pad is able to accept anymore
// data when flushing is TRUE. Calling this function with flushing FALSE makes pads accept data again. Caller
// must ensure that downstream streaming (thread) is not blocked, e.g. by sending a FLUSH_START downstream.
func (c *CollectPads) SetFlushing(flushing bool) {
	C.gst_collect_pads_set_flushing(c.Instance(), gboolean(flushing))
}

// SetFunction sets a function that overrides the behavior around the BufferFunction.
//
// CollectPads provides a default collection algorithm that will determine the oldest buffer available on all
// of its pads, and then delegate to a configured callback. However, if circumstances are more complicated
// and/or more control is desired, this sets a callback that will be invoked instead when all the pads added
// to the collection have buffers queued. Evidently, this callback is not compatible with SetBufferFunction
// callback. If this callback is set, the former will be unset.
func (c *CollectPads) SetFunction(f CollectPadsFunc) {
	c.funcMap.funcFunc = f
	C.gst_collect_pads_set_function(
		c.Instance(),
		C.GstCollectPadsFunction(C.cgoGstCollectPadsFunc),
		(C.gpointer)(c.selfPtr),
	)
}

// SetQueryFunction sets the query callback function and user data that will be called after collectpads has
// received a query originating from one of the collected pads. If the query being processed is a serialized
// one, this callback is called with pads STREAM_LOCK held, otherwise not. As this lock should be held when
// calling a number of CollectPads functions, it should be acquired if so (unusually) needed.
func (c *CollectPads) SetQueryFunction(f CollectPadsQueryFunc) {
	c.funcMap.queryFunc = f
	C.gst_collect_pads_set_query_function(
		c.Instance(),
		C.GstCollectPadsQueryFunction(C.cgoGstCollectPadsQueryFunc),
		(C.gpointer)(c.selfPtr),
	)
}

// SetWaiting sets a pad to waiting or non-waiting mode, if at least this pad has not been created with locked
// waiting state, in which case nothing happens.
//
// This function should be called with pads STREAM_LOCK held, such as in the callback.
func (c *CollectPads) SetWaiting(data *CollectData, waiting bool) {
	C.gst_collect_pads_set_waiting(c.Instance(), data.Instance(), gboolean(waiting))
}

// SrcEventDefault is the CollectPads event handling for the src pad of elements. Elements can chain up to this
// to let flushing seek event handling be done by CollectPads.
func (c *CollectPads) SrcEventDefault(pad *gst.Pad, event *gst.Event) bool {
	return gobool(C.gst_collect_pads_src_event_default(
		c.Instance(),
		(*C.GstPad)(unsafe.Pointer(pad.Instance())),
		(*C.GstEvent)(unsafe.Pointer(event.Instance())),
	))
}

// Start starts the processing of data in the collectpads.
func (c *CollectPads) Start() { C.gst_collect_pads_start(c.Instance()) }

// Stop stops the processing of data in the collectpads. It also u nblocks any blocking operations.
func (c *CollectPads) Stop() { C.gst_collect_pads_stop(c.Instance()) }

// TakeBuffer gets a subbuffer of size bytes from the given pad data. Flushes the amount of read bytes.
//
// This function should be called with pads STREAM_LOCK held, such as in the callback.
func (c *CollectPads) TakeBuffer(data *CollectData, size uint) *gst.Buffer {
	buf := C.gst_collect_pads_take_buffer(c.Instance(), data.Instance(), C.guint(size))
	if buf == nil {
		return nil
	}
	return gst.FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}
