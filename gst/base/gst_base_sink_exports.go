package base

/*
#include "gst.go.h"
*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

//export goGstBaseSinkActivatePull
func goGstBaseSinkActivatePull(sink *C.GstBaseSink, active C.gboolean) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		ActivatePull(self *GstBaseSink, active bool) bool
	})
	return gboolean(iface.ActivatePull(wrapGstBaseSink(sink), gobool(active)))
}

//export goGstBaseSinkEvent
func goGstBaseSinkEvent(sink *C.GstBaseSink, event *C.GstEvent) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Event(self *GstBaseSink, event *gst.Event) bool
	})
	return gboolean(iface.Event(wrapGstBaseSink(sink), gst.FromGstEventUnsafe(unsafe.Pointer(event))))
}

//export goGstBaseSinkFixate
func goGstBaseSinkFixate(sink *C.GstBaseSink, caps *C.GstCaps) *C.GstCaps {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Fixate(self *GstBaseSink, caps *gst.Caps) *gst.Caps
	})
	fixated := iface.Fixate(wrapGstBaseSink(sink), gst.FromGstCapsUnsafe(unsafe.Pointer(caps)))
	if fixated == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(fixated.Instance()))
}

//export goGstBaseSinkGetCaps
func goGstBaseSinkGetCaps(sink *C.GstBaseSink, filter *C.GstCaps) *C.GstCaps {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		GetCaps(self *GstBaseSink, filter *gst.Caps) *gst.Caps
	})
	filtered := iface.GetCaps(wrapGstBaseSink(sink), gst.FromGstCapsUnsafe(unsafe.Pointer(filter)))
	if filtered == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(filtered.Instance()))
}

//export goGstBaseSinkGetTimes
func goGstBaseSinkGetTimes(sink *C.GstBaseSink, buf *C.GstBuffer, start, end *C.GstClockTime) {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		GetTimes(self *GstBaseSink, buffer *gst.Buffer) (start, end time.Duration)
	})
	retStart, retEnd := iface.GetTimes(wrapGstBaseSink(sink), gst.FromGstBufferUnsafe(unsafe.Pointer(buf)))
	*start = C.GstClockTime(retStart.Nanoseconds())
	*end = C.GstClockTime(retEnd.Nanoseconds())
}

//export goGstBaseSinkPrepare
func goGstBaseSinkPrepare(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Prepare(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(iface.Prepare(wrapGstBaseSink(sink), gst.FromGstBufferUnsafe(unsafe.Pointer(buf))))
}

//export goGstBaseSinkPrepareList
func goGstBaseSinkPrepareList(sink *C.GstBaseSink, list *C.GstBufferList) C.GstFlowReturn {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		PrepareList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	})
	return C.GstFlowReturn(iface.PrepareList(wrapGstBaseSink(sink), gst.FromGstBufferListUnsafe(unsafe.Pointer(list))))
}

//export goGstBaseSinkPreroll
func goGstBaseSinkPreroll(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Preroll(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(iface.Preroll(wrapGstBaseSink(sink), gst.FromGstBufferUnsafe(unsafe.Pointer(buf))))
}

//export goGstBaseSinkProposeAllocation
func goGstBaseSinkProposeAllocation(sink *C.GstBaseSink, query *C.GstQuery) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		ProposeAllocation(self *GstBaseSink, query *gst.Query) bool
	})
	return gboolean(iface.ProposeAllocation(wrapGstBaseSink(sink), gst.FromGstQueryUnsafe(unsafe.Pointer(query))))
}

//export goGstBaseSinkQuery
func goGstBaseSinkQuery(sink *C.GstBaseSink, query *C.GstQuery) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Query(self *GstBaseSink, query *gst.Query) bool
	})
	return gboolean(iface.Query(wrapGstBaseSink(sink), gst.FromGstQueryUnsafe(unsafe.Pointer(query))))
}

//export goGstBaseSinkRender
func goGstBaseSinkRender(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Render(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(iface.Render(wrapGstBaseSink(sink), gst.FromGstBufferUnsafe(unsafe.Pointer(buf))))
}

//export goGstBaseSinkRenderList
func goGstBaseSinkRenderList(sink *C.GstBaseSink, buf *C.GstBufferList) C.GstFlowReturn {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		RenderList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
	})
	return C.GstFlowReturn(iface.RenderList(wrapGstBaseSink(sink), gst.FromGstBufferListUnsafe(unsafe.Pointer(buf))))
}

//export goGstBaseSinkSetCaps
func goGstBaseSinkSetCaps(sink *C.GstBaseSink, caps *C.GstCaps) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		SetCaps(self *GstBaseSink, caps *gst.Caps) bool
	})
	return gboolean(iface.SetCaps(wrapGstBaseSink(sink), gst.FromGstCapsUnsafe(unsafe.Pointer(caps))))
}

//export goGstBaseSinkStart
func goGstBaseSinkStart(sink *C.GstBaseSink) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Start(self *GstBaseSink) bool
	})
	return gboolean(iface.Start(wrapGstBaseSink(sink)))
}

//export goGstBaseSinkStop
func goGstBaseSinkStop(sink *C.GstBaseSink) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Stop(self *GstBaseSink) bool
	})
	return gboolean(iface.Stop(wrapGstBaseSink(sink)))
}

//export goGstBaseSinkUnlock
func goGstBaseSinkUnlock(sink *C.GstBaseSink) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		Unlock(self *GstBaseSink) bool
	})
	return gboolean(iface.Unlock(wrapGstBaseSink(sink)))
}

//export goGstBaseSinkUnlockStop
func goGstBaseSinkUnlockStop(sink *C.GstBaseSink) C.gboolean {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		UnlockStop(self *GstBaseSink) bool
	})
	return gboolean(iface.UnlockStop(wrapGstBaseSink(sink)))
}

//export goGstBaseSinkWaitEvent
func goGstBaseSinkWaitEvent(sink *C.GstBaseSink, event *C.GstEvent) C.GstFlowReturn {
	iface := gst.FromObjectUnsafePrivate(unsafe.Pointer(sink)).(interface {
		WaitEvent(self *GstBaseSink, event *gst.Event) gst.FlowReturn
	})
	return C.GstFlowReturn(iface.WaitEvent(wrapGstBaseSink(sink), gst.FromGstEventUnsafe(unsafe.Pointer(event))))
}
