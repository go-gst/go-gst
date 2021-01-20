package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

//export goGstBaseSinkActivatePull
func goGstBaseSinkActivatePull(sink *C.GstBaseSink, active C.gboolean) C.gboolean {
	var ret C.gboolean
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			ActivatePull(self *GstBaseSink, active bool) bool
		})
		ret = gboolean(iface.ActivatePull(ToGstBaseSink(gObject), gobool(active)))
	})
	return ret
}

//export goGstBaseSinkEvent
func goGstBaseSinkEvent(sink *C.GstBaseSink, event *C.GstEvent) C.gboolean {
	var ret C.gboolean
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Event(self *GstBaseSink, event *gst.Event) bool
		})
		ret = gboolean(iface.Event(ToGstBaseSink(gObject), gst.ToGstEvent(unsafe.Pointer(event))))
	})
	return ret
}

//export goGstBaseSinkFixate
func goGstBaseSinkFixate(sink *C.GstBaseSink, caps *C.GstCaps) *C.GstCaps {
	var fixated *gst.Caps
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Fixate(self *GstBaseSink, caps *gst.Caps) *gst.Caps
		})
		fixated = iface.Fixate(ToGstBaseSink(gObject), gst.ToGstCaps(unsafe.Pointer(caps)))
	})
	if fixated == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(fixated.Instance()))
}

//export goGstBaseSinkGetCaps
func goGstBaseSinkGetCaps(sink *C.GstBaseSink, filter *C.GstCaps) *C.GstCaps {
	var filtered *gst.Caps
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			GetCaps(self *GstBaseSink, filter *gst.Caps) *gst.Caps
		})
		filtered = iface.GetCaps(ToGstBaseSink(gObject), gst.ToGstCaps(unsafe.Pointer(filter)))
	})
	if filtered == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(filtered.Instance()))
}

//export goGstBaseSinkGetTimes
func goGstBaseSinkGetTimes(sink *C.GstBaseSink, buf *C.GstBuffer, start, end *C.GstClockTime) {
	var retStart, retEnd time.Duration
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			GetTimes(self *GstBaseSink, buffer *gst.Buffer) (start, end time.Duration)
		})
		retStart, retEnd = iface.GetTimes(ToGstBaseSink(gObject), gst.ToGstBuffer(unsafe.Pointer(buf)))
	})
	*start = C.GstClockTime(retStart.Nanoseconds())
	*end = C.GstClockTime(retEnd.Nanoseconds())
}

//export goGstBaseSinkPrepare
func goGstBaseSinkPrepare(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Prepare(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
		})
		ret = iface.Prepare(ToGstBaseSink(gObject), gst.ToGstBuffer(unsafe.Pointer(buf)))
	})
	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkPrepareList
func goGstBaseSinkPrepareList(sink *C.GstBaseSink, list *C.GstBufferList) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			PrepareList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
		})
		ret = iface.PrepareList(ToGstBaseSink(gObject), gst.ToGstBufferList(unsafe.Pointer(list)))
	})
	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkPreroll
func goGstBaseSinkPreroll(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Preroll(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
		})
		ret = iface.Preroll(ToGstBaseSink(gObject), gst.ToGstBuffer(unsafe.Pointer(buf)))
	})
	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkProposeAllocation
func goGstBaseSinkProposeAllocation(sink *C.GstBaseSink, query *C.GstQuery) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			ProposeAllocation(self *GstBaseSink, query *gst.Query) bool
		})
		ret = iface.ProposeAllocation(ToGstBaseSink(gObject), gst.ToGstQuery(unsafe.Pointer(query)))
	})
	return gboolean(ret)
}

//export goGstBaseSinkQuery
func goGstBaseSinkQuery(sink *C.GstBaseSink, query *C.GstQuery) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Query(self *GstBaseSink, query *gst.Query) bool
		})
		ret = iface.Query(ToGstBaseSink(gObject), gst.ToGstQuery(unsafe.Pointer(query)))
	})
	return gboolean(ret)
}

//export goGstBaseSinkRender
func goGstBaseSinkRender(sink *C.GstBaseSink, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Render(self *GstBaseSink, buffer *gst.Buffer) gst.FlowReturn
		})
		ret = iface.Render(ToGstBaseSink(gObject), gst.ToGstBuffer(unsafe.Pointer(buf)))
	})
	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkRenderList
func goGstBaseSinkRenderList(sink *C.GstBaseSink, buf *C.GstBufferList) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			RenderList(self *GstBaseSink, bufferList *gst.BufferList) gst.FlowReturn
		})
		ret = iface.RenderList(ToGstBaseSink(gObject), gst.ToGstBufferList(unsafe.Pointer(buf)))
	})
	return C.GstFlowReturn(ret)
}

//export goGstBaseSinkSetCaps
func goGstBaseSinkSetCaps(sink *C.GstBaseSink, caps *C.GstCaps) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			SetCaps(self *GstBaseSink, caps *gst.Caps) bool
		})
		ret = iface.SetCaps(ToGstBaseSink(gObject), gst.ToGstCaps(unsafe.Pointer(caps)))
	})
	return gboolean(ret)
}

//export goGstBaseSinkStart
func goGstBaseSinkStart(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Start(self *GstBaseSink) bool
		})
		ret = iface.Start(ToGstBaseSink(gObject))
	})
	return gboolean(ret)
}

//export goGstBaseSinkStop
func goGstBaseSinkStop(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Stop(self *GstBaseSink) bool
		})
		ret = iface.Stop(ToGstBaseSink(gObject))
	})
	return gboolean(ret)
}

//export goGstBaseSinkUnlock
func goGstBaseSinkUnlock(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Unlock(self *GstBaseSink) bool
		})
		ret = iface.Unlock(ToGstBaseSink(gObject))
	})
	return gboolean(ret)
}

//export goGstBaseSinkUnlockStop
func goGstBaseSinkUnlockStop(sink *C.GstBaseSink) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			UnlockStop(self *GstBaseSink) bool
		})
		ret = iface.UnlockStop(ToGstBaseSink(gObject))
	})
	return gboolean(ret)
}

//export goGstBaseSinkWaitEvent
func goGstBaseSinkWaitEvent(sink *C.GstBaseSink, event *C.GstEvent) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(sink), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			WaitEvent(self *GstBaseSink, event *gst.Event) gst.FlowReturn
		})
		ret = iface.WaitEvent(ToGstBaseSink(gObject), gst.ToGstEvent(unsafe.Pointer(event)))
	})
	return C.GstFlowReturn(ret)
}
