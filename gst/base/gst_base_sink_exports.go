package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

//export goGstBaseSinkActivatePull
func goGstBaseSinkActivatePull(sink *C.GstBaseSink, active C.gboolean) C.gboolean {
	var ret C.gboolean
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		ActivatePull(self *GstBaseSink, active bool) bool
	})
	ret = gboolean(iface.ActivatePull(goBaseSink, gobool(active)))

	return ret
}

//export goGstBaseSinkEvent
func goGstBaseSinkEvent(sink *C.GstBaseSink, event *C.GstEvent) C.gboolean {
	var ret C.gboolean
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Event(self *GstBaseSink, event *gst.Event) bool
	})
	ret = gboolean(iface.Event(goBaseSink, gst.ToGstEvent(unsafe.Pointer(event))))

	return ret
}

//export goGstBaseSinkFixate
func goGstBaseSinkFixate(sink *C.GstBaseSink, caps *C.GstCaps) *C.GstCaps {
	var fixated *gst.Caps
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Fixate(self *GstBaseSink, caps *gst.Caps) *gst.Caps
	})
	fixated = iface.Fixate(goBaseSink, gst.ToGstCaps(unsafe.Pointer(caps)))

	if fixated == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(fixated.Instance()))
}

//export goGstBaseSinkGetCaps
func goGstBaseSinkGetCaps(sink *C.GstBaseSink, filter *C.GstCaps) *C.GstCaps {
	var filtered *gst.Caps
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		GetCaps(self *GstBaseSink, filter *gst.Caps) *gst.Caps
	})
	filtered = iface.GetCaps(goBaseSink, gst.ToGstCaps(unsafe.Pointer(filter)))

	if filtered == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(filtered.Instance()))
}

//export goGstBaseSinkGetTimes
func goGstBaseSinkGetTimes(sink *C.GstBaseSink, buf *C.GstBuffer, start, end *C.GstClockTime) {
	var retStart, retEnd time.Duration
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		GetTimes(self *GstBaseSink, buffer *gst.Buffer) (start, end time.Duration) // should this be a ClockTime?
	})
	retStart, retEnd = iface.GetTimes(goBaseSink, gst.ToGstBuffer(unsafe.Pointer(buf)))

	*start = C.GstClockTime(retStart.Nanoseconds())
	*end = C.GstClockTime(retEnd.Nanoseconds())
}

//export goGstBaseSinkPrepare
func goGstBaseSinkPrepare(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Prepare(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	})
	ret = iface.Prepare(goBaseSink, gst.ToGstBuffer(unsafe.Pointer(buf)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkPrepareList
func goGstBaseSinkPrepareList(sink *C.GstBaseSink, list *C.GstBufferList) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		PrepareList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	})
	ret = iface.PrepareList(goBaseSink, gst.ToGstBufferList(unsafe.Pointer(list)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkPreroll
func goGstBaseSinkPreroll(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Preroll(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	})
	ret = iface.Preroll(goBaseSink, gst.ToGstBuffer(unsafe.Pointer(buf)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkProposeAllocation
func goGstBaseSinkProposeAllocation(sink *C.GstBaseSink, query *C.GstQuery) C.gboolean {
	var ret bool
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		ProposeAllocation(self *GstBaseSink, query *gst.Query) bool
	})
	ret = iface.ProposeAllocation(goBaseSink, gst.ToGstQuery(unsafe.Pointer(query)))

	return gboolean(ret)
}

//export goGstBaseSinkQuery
func goGstBaseSinkQuery(sink *C.GstBaseSink, query *C.GstQuery) C.gboolean {
	var ret bool
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Query(self *GstBaseSink, query *gst.Query) bool
	})
	ret = iface.Query(goBaseSink, gst.ToGstQuery(unsafe.Pointer(query)))

	return gboolean(ret)
}

//export goGstBaseSinkRender
func goGstBaseSinkRender(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Render(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	})
	ret = iface.Render(goBaseSink, gst.ToGstBuffer(unsafe.Pointer(buf)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkRenderList
func goGstBaseSinkRenderList(sink *C.GstBaseSink, buf *C.GstBufferList) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		RenderList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	})
	ret = iface.RenderList(goBaseSink, gst.ToGstBufferList(unsafe.Pointer(buf)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkSetCaps
func goGstBaseSinkSetCaps(sink *C.GstBaseSink, caps *C.GstCaps) C.gboolean {
	var ret bool
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		SetCaps(self *GstBaseSink, caps *gst.Caps) bool
	})
	ret = iface.SetCaps(goBaseSink, gst.ToGstCaps(unsafe.Pointer(caps)))

	return gboolean(ret)
}

//export goGstBaseSinkStart
func goGstBaseSinkStart(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Start(self *GstBaseSink) bool
	})
	ret = iface.Start(goBaseSink)

	return gboolean(ret)
}

//export goGstBaseSinkStop
func goGstBaseSinkStop(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Stop(self *GstBaseSink) bool
	})
	ret = iface.Stop(goBaseSink)

	return gboolean(ret)
}

//export goGstBaseSinkUnlock
func goGstBaseSinkUnlock(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		Unlock(self *GstBaseSink) bool
	})
	ret = iface.Unlock(goBaseSink)

	return gboolean(ret)
}

//export goGstBaseSinkUnlockStop
func goGstBaseSinkUnlockStop(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		UnlockStop(self *GstBaseSink) bool
	})
	ret = iface.UnlockStop(goBaseSink)

	return gboolean(ret)
}

//export goGstBaseSinkWaitEvent
func goGstBaseSinkWaitEvent(sink *C.GstBaseSink, event *C.GstEvent) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseSink := ToGstBaseSink(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(sink))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(sink))

	iface := subclass.(interface {
		WaitEvent(self *GstBaseSink, event *gst.Event) gst.FlowReturn
	})
	ret = iface.WaitEvent(goBaseSink, gst.ToGstEvent(unsafe.Pointer(event)))

	return C.GstFlowReturn(ret)
}
