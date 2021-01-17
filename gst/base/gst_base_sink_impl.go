package base

/*
#include "gst.go.h"

extern gboolean       goGstBaseSinkActivatePull       (GstBaseSink * sink, gboolean active);
extern gboolean       goGstBaseSinkEvent              (GstBaseSink * sink, GstEvent * event);
extern GstCaps *      goGstBaseSinkFixate             (GstBaseSink * sink, GstCaps * caps);
extern GstCaps *      goGstBaseSinkGetCaps            (GstBaseSink * sink, GstCaps * filter);
extern void           goGstBaseSinkGetTimes           (GstBaseSink * sink, GstBuffer * buffer, GstClockTime * start, GstClockTime * end);
extern GstFlowReturn  goGstBaseSinkPrepare            (GstBaseSink * sink, GstBuffer * buffer);
extern GstFlowReturn  goGstBaseSinkPrepareList        (GstBaseSink * sink, GstBufferList * buffer_list);
extern GstFlowReturn  goGstBaseSinkPreroll            (GstBaseSink * sink, GstBuffer * buffer);
extern gboolean       goGstBaseSinkProposeAllocation  (GstBaseSink * sink, GstQuery * query);
extern gboolean       goGstBaseSinkQuery              (GstBaseSink * sink, GstQuery * query);
extern GstFlowReturn  goGstBaseSinkRender             (GstBaseSink * sink, GstBuffer * buffer);
extern GstFlowReturn  goGstBaseSinkRenderList         (GstBaseSink * sink, GstBufferList * buffer_list);
extern gboolean       goGstBaseSinkSetCaps            (GstBaseSink * sink, GstCaps * caps);
extern gboolean       goGstBaseSinkStart              (GstBaseSink * sink);
extern gboolean       goGstBaseSinkStop               (GstBaseSink * sink);
extern gboolean       goGstBaseSinkUnlock             (GstBaseSink * sink);
extern gboolean       goGstBaseSinkUnlockStop         (GstBaseSink * sink);
extern GstFlowReturn  goGstBaseSinkWaitEvent          (GstBaseSink * sink, GstEvent * event);

GstFlowReturn   do_wait_event  (GstBaseSink * sink, GstEvent * event)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(sink));
	GstBaseSinkClass * parent = toGstBaseSinkClass(g_type_class_peek_parent(this_class));
	GstFlowReturn ret = parent->wait_event(sink, event);
	if (ret == GST_FLOW_ERROR)
		return ret;
	return goGstBaseSinkWaitEvent(sink, event);
}

void setGstBaseSinkActivatePull       (GstBaseSinkClass * klass)  { klass->activate_pull = goGstBaseSinkActivatePull; }
void setGstBaseSinkEvent              (GstBaseSinkClass * klass)  { klass->event = goGstBaseSinkEvent; }
void setGstBaseSinkFixate             (GstBaseSinkClass * klass)  { klass->fixate = goGstBaseSinkFixate; }
void setGstBaseSinkGetCaps            (GstBaseSinkClass * klass)  { klass->get_caps = goGstBaseSinkGetCaps; }
void setGstBaseSinkGetTimes           (GstBaseSinkClass * klass)  { klass->get_times = goGstBaseSinkGetTimes; }
void setGstBaseSinkPrepare            (GstBaseSinkClass * klass)  { klass->prepare = goGstBaseSinkPrepare; }
void setGstBaseSinkPrepareList        (GstBaseSinkClass * klass)  { klass->prepare_list = goGstBaseSinkPrepareList; }
void setGstBaseSinkPreroll            (GstBaseSinkClass * klass)  { klass->preroll = goGstBaseSinkPreroll; }
void setGstBaseSinkProposeAllocation  (GstBaseSinkClass * klass)  { klass->propose_allocation = goGstBaseSinkProposeAllocation; }
void setGstBaseSinkQuery              (GstBaseSinkClass * klass)  { klass->query = goGstBaseSinkQuery; }
void setGstBaseSinkRender             (GstBaseSinkClass * klass)  { klass->render = goGstBaseSinkRender; }
void setGstBaseSinkRenderList         (GstBaseSinkClass * klass)  { klass->render_list = goGstBaseSinkRenderList; }
void setGstBaseSinkSetCaps            (GstBaseSinkClass * klass)  { klass->set_caps = goGstBaseSinkSetCaps; }
void setGstBaseSinkStart              (GstBaseSinkClass * klass)  { klass->start = goGstBaseSinkStart; }
void setGstBaseSinkStop               (GstBaseSinkClass * klass)  { klass->stop = goGstBaseSinkStop; }
void setGstBaseSinkUnlock             (GstBaseSinkClass * klass)  { klass->unlock = goGstBaseSinkUnlock; }
void setGstBaseSinkUnlockStop         (GstBaseSinkClass * klass)  { klass->unlock_stop = goGstBaseSinkUnlockStop; }
void setGstBaseSinkWaitEvent          (GstBaseSinkClass * klass)  { klass->wait_event = do_wait_event; }

*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

var (
	// ExtendsBaseSink is an Extendable for extending a GstBaseSink
	ExtendsBaseSink glib.Extendable = &extendsBaseSink{parent: gst.ExtendsElement}
)

// GstBaseSinkImpl is the documented interface for extending a GstBaseSink. It does not have to
// be implemented in it's entirety. Each of the methods it declares will be checked for their presence
// in the initializing object, and if the object declares an override it will replace the default
// implementation in the virtual methods.
type GstBaseSinkImpl interface {
	// Subclasses should override this when they can provide an alternate method of spawning a thread to
	// drive the pipeline in pull mode. Should start or stop the pulling thread, depending on the value
	// of the "active" argument. Called after actually activating the sink pad in pull mode. The default
	// implementation starts a task on the sink pad.
	ActivatePull(self *GstBaseSink, active bool) bool
	// Override this to handle events arriving on the sink pad
	Event(self *GstBaseSink, event *gst.Event) bool
	// Only useful in pull mode. Implement if you have ideas about what should be the default values for
	// the caps you support.
	Fixate(self *GstBaseSink, caps *gst.Caps) *gst.Caps
	// Called to get sink pad caps from the subclass
	GetCaps(self *GstBaseSink, filter *gst.Caps) *gst.Caps
	// Called to get the start and end times for synchronising the passed buffer to the clock
	GetTimes(self *GstBaseSink, buffer *gst.Buffer) (start, end time.Duration)
	// Called to prepare the buffer for render and preroll. This function is called before synchronization
	// is performed.
	Prepare(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	// Called to prepare the buffer list for render_list. This function is called before synchronization is
	// performed.
	PrepareList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	// Called to present the preroll buffer if desired.
	Preroll(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	// Used to configure the allocation query
	ProposeAllocation(self *GstBaseSink, query *gst.Query) bool
	// Handle queries on the element
	Query(self *GstBaseSink, query *gst.Query) bool
	// Called when a buffer should be presented or output, at the correct moment if the GstBaseSink has been
	// set to sync to the clock.
	Render(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	// Same as render but used with buffer lists instead of buffers.
	RenderList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	// Notify subclass of changed caps
	SetCaps(self *GstBaseSink, caps *gst.Caps) bool
	// Start processing. Ideal for opening resources in the subclass
	Start(self *GstBaseSink) bool
	// Stop processing. Subclasses should use this to close resources.
	Stop(self *GstBaseSink) bool
	// Unlock any pending access to the resource. Subclasses should unblock any blocked function ASAP and call
	// WaitPreroll
	Unlock(self *GstBaseSink) bool
	// Clear the previous unlock request. Subclasses should clear any state they set during Unlock(), and be ready
	// to continue where they left off after WaitPreroll, Wait or WaitClock return or Render() is called again.
	UnlockStop(self *GstBaseSink) bool
	// Override this to implement custom logic to wait for the event time (for events like EOS and GAP). The bindings
	// take care of first chaining up to the parent class.
	WaitEvent(self *GstBaseSink, event *gst.Event) gst.FlowReturn
}

type extendsBaseSink struct{ parent glib.Extendable }

func (e *extendsBaseSink) Type() glib.Type     { return glib.Type(C.gst_base_sink_get_type()) }
func (e *extendsBaseSink) ClassSize() int64    { return int64(C.sizeof_GstBaseSinkClass) }
func (e *extendsBaseSink) InstanceSize() int64 { return int64(C.sizeof_GstBaseSink) }

func (e *extendsBaseSink) InitClass(klass unsafe.Pointer, elem glib.GoObjectSubclass) {
	e.parent.InitClass(klass, elem)

	sinkClass := C.toGstBaseSinkClass(klass)

	if _, ok := elem.(interface {
		ActivatePull(self *GstBaseSink, active bool) bool
	}); ok {
		C.setGstBaseSinkActivatePull(sinkClass)
	}
	if _, ok := elem.(interface {
		Event(self *GstBaseSink, event *gst.Event) bool
	}); ok {
		C.setGstBaseSinkEvent(sinkClass)
	}

	if _, ok := elem.(interface {
		Fixate(self *GstBaseSink, caps *gst.Caps) *gst.Caps
	}); ok {
		C.setGstBaseSinkFixate(sinkClass)
	}

	if _, ok := elem.(interface {
		GetCaps(self *GstBaseSink, filter *gst.Caps) *gst.Caps
	}); ok {
		C.setGstBaseSinkGetCaps(sinkClass)
	}

	if _, ok := elem.(interface {
		GetTimes(self *GstBaseSink, buffer *gst.Buffer) (start, end time.Duration)
	}); ok {
		C.setGstBaseSinkGetTimes(sinkClass)
	}

	if _, ok := elem.(interface {
		Prepare(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstBaseSinkPrepare(sinkClass)
	}

	if _, ok := elem.(interface {
		PrepareList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	}); ok {
		C.setGstBaseSinkPrepareList(sinkClass)
	}

	if _, ok := elem.(interface {
		Preroll(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstBaseSinkPreroll(sinkClass)
	}

	if _, ok := elem.(interface {
		ProposeAllocation(self *GstBaseSink, query *gst.Query) bool
	}); ok {
		C.setGstBaseSinkProposeAllocation(sinkClass)
	}

	if _, ok := elem.(interface {
		Query(self *GstBaseSink, query *gst.Query) bool
	}); ok {
		C.setGstBaseSinkQuery(sinkClass)
	}

	if _, ok := elem.(interface {
		Render(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstBaseSinkRender(sinkClass)
	}

	if _, ok := elem.(interface {
		RenderList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	}); ok {
		C.setGstBaseSinkRenderList(sinkClass)
	}

	if _, ok := elem.(interface {
		SetCaps(self *GstBaseSink, caps *gst.Caps) bool
	}); ok {
		C.setGstBaseSinkSetCaps(sinkClass)
	}

	if _, ok := elem.(interface {
		Start(self *GstBaseSink) bool
	}); ok {
		C.setGstBaseSinkStart(sinkClass)
	}

	if _, ok := elem.(interface {
		Stop(self *GstBaseSink) bool
	}); ok {
		C.setGstBaseSinkStop(sinkClass)
	}

	if _, ok := elem.(interface {
		Unlock(self *GstBaseSink) bool
	}); ok {
		C.setGstBaseSinkUnlock(sinkClass)
	}

	if _, ok := elem.(interface {
		UnlockStop(self *GstBaseSink) bool
	}); ok {
		C.setGstBaseSinkUnlockStop(sinkClass)
	}

	if _, ok := elem.(interface {
		WaitEvent(self *GstBaseSink, event *gst.Event) gst.FlowReturn
	}); ok {
		C.setGstBaseSinkWaitEvent(sinkClass)
	}

}
