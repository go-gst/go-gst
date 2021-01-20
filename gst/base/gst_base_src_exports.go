package base

//#include "gst.go.h"
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

//export goGstBaseSrcGetCaps
func goGstBaseSrcGetCaps(src *C.GstBaseSrc, filter *C.GstCaps) *C.GstCaps {
	var caps *gst.Caps
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			GetCaps(*GstBaseSrc, *gst.Caps) *gst.Caps
		})
		caps = iface.GetCaps(ToGstBaseSrc(gObject), gst.ToGstCaps(unsafe.Pointer(filter)))
	})
	if caps == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(caps.Instance()))
}

//export goGstBaseSrcNegotiate
func goGstBaseSrcNegotiate(src *C.GstBaseSrc) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Negotiate(*GstBaseSrc) bool
		})
		ret = iface.Negotiate(ToGstBaseSrc(gObject))
	})
	return gboolean(ret)
}

//export goGstBaseSrcFixate
func goGstBaseSrcFixate(src *C.GstBaseSrc, filter *C.GstCaps) *C.GstCaps {
	var caps *gst.Caps
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Fixate(*GstBaseSrc, *gst.Caps) *gst.Caps
		})
		caps = iface.Fixate(ToGstBaseSrc(gObject), gst.ToGstCaps(unsafe.Pointer(filter)))
	})
	if caps == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(caps.Instance()))
}

//export goGstBaseSrcSetCaps
func goGstBaseSrcSetCaps(src *C.GstBaseSrc, caps *C.GstCaps) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			SetCaps(*GstBaseSrc, *gst.Caps) bool
		})
		ret = iface.SetCaps(ToGstBaseSrc(gObject), gst.ToGstCaps(unsafe.Pointer(caps)))
	})
	return gboolean(ret)
}

//export goGstBaseSrcDecideAllocation
func goGstBaseSrcDecideAllocation(src *C.GstBaseSrc, query *C.GstQuery) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			DecideAllocation(*GstBaseSrc, *gst.Query) bool
		})
		ret = iface.DecideAllocation(ToGstBaseSrc(gObject), gst.ToGstQuery(unsafe.Pointer(query)))
	})
	return gboolean(ret)
}

//export goGstBaseSrcStart
func goGstBaseSrcStart(src *C.GstBaseSrc) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Start(*GstBaseSrc) bool
		})
		ret = iface.Start(ToGstBaseSrc(gObject))
	})
	return gboolean(ret)
}

//export goGstBaseSrcStop
func goGstBaseSrcStop(src *C.GstBaseSrc) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Stop(*GstBaseSrc) bool
		})
		ret = iface.Stop(ToGstBaseSrc(gObject))
	})
	return gboolean(ret)
}

//export goGstBaseSrcGetTimes
func goGstBaseSrcGetTimes(src *C.GstBaseSrc, buf *C.GstBuffer, start *C.GstClockTime, end *C.GstClockTime) {
	var gostart, goend time.Duration
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			GetTimes(*GstBaseSrc, *gst.Buffer) (start, end time.Duration)
		})
		gostart, goend = iface.GetTimes(ToGstBaseSrc(gObject), gst.ToGstBuffer(unsafe.Pointer(buf)))
	})
	*start = C.GstClockTime(gostart.Nanoseconds())
	*end = C.GstClockTime(goend.Nanoseconds())
}

//export goGstBaseSrcGetSize
func goGstBaseSrcGetSize(src *C.GstBaseSrc, size *C.guint64) C.gboolean {
	var gosize int64
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			GetSize(*GstBaseSrc) (bool, int64)
		})
		ok, gosize = iface.GetSize(ToGstBaseSrc(gObject))
	})
	if ok {
		*size = C.guint64(gosize)
	}
	return gboolean(ok)
}

//export goGstBaseSrcIsSeekable
func goGstBaseSrcIsSeekable(src *C.GstBaseSrc) C.gboolean {
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			IsSeekable(*GstBaseSrc) bool
		})
		ok = iface.IsSeekable(ToGstBaseSrc(gObject))
	})
	return gboolean(ok)
}

//export goGstBaseSrcPrepareSeekSegment
func goGstBaseSrcPrepareSeekSegment(src *C.GstBaseSrc, seek *C.GstEvent, segment *C.GstSegment) C.gboolean {
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			PrepareSeekSegment(*GstBaseSrc, *gst.Event, *gst.Segment) bool
		})
		ok = iface.PrepareSeekSegment(ToGstBaseSrc(gObject), gst.ToGstEvent(unsafe.Pointer(seek)), gst.FromGstSegmentUnsafe(unsafe.Pointer(segment)))
	})
	return gboolean(ok)
}

//export goGstBaseSrcDoSeek
func goGstBaseSrcDoSeek(src *C.GstBaseSrc, segment *C.GstSegment) C.gboolean {
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			DoSeek(*GstBaseSrc, *gst.Segment) bool
		})
		ok = iface.DoSeek(ToGstBaseSrc(gObject), gst.ToGstSegment(unsafe.Pointer(segment)))
	})
	return gboolean(ok)
}

//export goGstBaseSrcUnlock
func goGstBaseSrcUnlock(src *C.GstBaseSrc) C.gboolean {
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Unlock(*GstBaseSrc) bool
		})
		ok = iface.Unlock(ToGstBaseSrc(gObject))
	})
	return gboolean(ok)
}

//export goGstBaseSrcUnlockStop
func goGstBaseSrcUnlockStop(src *C.GstBaseSrc) C.gboolean {
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			UnlockStop(*GstBaseSrc) bool
		})
		ok = iface.UnlockStop(ToGstBaseSrc(gObject))
	})
	return gboolean(ok)
}

//export goGstBaseSrcQuery
func goGstBaseSrcQuery(src *C.GstBaseSrc, query *C.GstQuery) C.gboolean {
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Query(*GstBaseSrc, *gst.Query) bool
		})
		ok = iface.Query(ToGstBaseSrc(gObject), gst.ToGstQuery(unsafe.Pointer(query)))
	})
	return gboolean(ok)
}

//export goGstBaseSrcEvent
func goGstBaseSrcEvent(src *C.GstBaseSrc, event *C.GstEvent) C.gboolean {
	var ok bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Event(*GstBaseSrc, *gst.Event) bool
		})
		ok = iface.Event(ToGstBaseSrc(gObject), gst.ToGstEvent(unsafe.Pointer(event)))
	})
	return gboolean(ok)
}

//export goGstBaseSrcCreate
func goGstBaseSrcCreate(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var gobuf *gst.Buffer
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Create(*GstBaseSrc, uint64, uint) (gst.FlowReturn, *gst.Buffer)
		})
		ret, gobuf = iface.Create(ToGstBaseSrc(gObject), uint64(offset), uint(size))
	})
	if ret == gst.FlowOK {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(gobuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseSrcAlloc
func goGstBaseSrcAlloc(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var gobuf *gst.Buffer
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Alloc(*GstBaseSrc, uint64, uint) (gst.FlowReturn, *gst.Buffer)
		})
		ret, gobuf = iface.Alloc(ToGstBaseSrc(gObject), uint64(offset), uint(size))
	})
	if ret == gst.FlowOK {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(gobuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseSrcFill
func goGstBaseSrcFill(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Fill(*GstBaseSrc, uint64, uint, *gst.Buffer) gst.FlowReturn
		})
		ret = iface.Fill(ToGstBaseSrc(gObject), uint64(offset), uint(size), gst.ToGstBuffer(unsafe.Pointer(buf)))
	})
	return C.GstFlowReturn(ret)
}
