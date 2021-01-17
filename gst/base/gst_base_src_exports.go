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
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		GetCaps(*GstBaseSrc, *gst.Caps) *gst.Caps
	})
	res := caller.GetCaps(wrapGstBaseSrc(src), gst.FromGstCapsUnsafe(unsafe.Pointer(filter)))
	return (*C.GstCaps)(unsafe.Pointer(res.Instance()))
}

//export goGstBaseSrcNegotiate
func goGstBaseSrcNegotiate(src *C.GstBaseSrc) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Negotiate(*GstBaseSrc) bool
	})
	return gboolean(caller.Negotiate(wrapGstBaseSrc(src)))
}

//export goGstBaseSrcFixate
func goGstBaseSrcFixate(src *C.GstBaseSrc, caps *C.GstCaps) *C.GstCaps {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Fixate(*GstBaseSrc, *gst.Caps) *gst.Caps
	})
	res := caller.Fixate(wrapGstBaseSrc(src), gst.FromGstCapsUnsafe(unsafe.Pointer(caps)))
	return (*C.GstCaps)(unsafe.Pointer(res.Instance()))
}

//export goGstBaseSrcSetCaps
func goGstBaseSrcSetCaps(src *C.GstBaseSrc, filter *C.GstCaps) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		SetCaps(*GstBaseSrc, *gst.Caps) bool
	})
	return gboolean(caller.SetCaps(wrapGstBaseSrc(src), gst.FromGstCapsUnsafe(unsafe.Pointer(filter))))
}

//export goGstBaseSrcDecideAllocation
func goGstBaseSrcDecideAllocation(src *C.GstBaseSrc, query *C.GstQuery) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		DecideAllocation(*GstBaseSrc, *gst.Query) bool
	})
	return gboolean(caller.DecideAllocation(wrapGstBaseSrc(src), gst.FromGstQueryUnsafe(unsafe.Pointer(query))))
}

//export goGstBaseSrcStart
func goGstBaseSrcStart(src *C.GstBaseSrc) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Start(*GstBaseSrc) bool
	})
	return gboolean(caller.Start(wrapGstBaseSrc(src)))
}

//export goGstBaseSrcStop
func goGstBaseSrcStop(src *C.GstBaseSrc) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Stop(*GstBaseSrc) bool
	})
	return gboolean(caller.Stop(wrapGstBaseSrc(src)))
}

//export goGstBaseSrcGetTimes
func goGstBaseSrcGetTimes(src *C.GstBaseSrc, buf *C.GstBuffer, start *C.GstClockTime, end *C.GstClockTime) {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		GetTimes(*GstBaseSrc, *gst.Buffer) (start, end time.Duration)
	})
	gostart, goend := caller.GetTimes(wrapGstBaseSrc(src), gst.FromGstBufferUnsafe(unsafe.Pointer(buf)))
	*start = C.GstClockTime(gostart.Nanoseconds())
	*end = C.GstClockTime(goend.Nanoseconds())
}

//export goGstBaseSrcGetSize
func goGstBaseSrcGetSize(src *C.GstBaseSrc, size *C.guint64) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		GetSize(*GstBaseSrc) (bool, int64)
	})
	ok, gosize := caller.GetSize(wrapGstBaseSrc(src))
	if !ok {
		return gboolean(ok)
	}
	*size = C.guint64(gosize)
	return gboolean(ok)
}

//export goGstBaseSrcIsSeekable
func goGstBaseSrcIsSeekable(src *C.GstBaseSrc) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		IsSeekable(*GstBaseSrc) bool
	})
	return gboolean(caller.IsSeekable(wrapGstBaseSrc(src)))
}

//export goGstBaseSrcPrepareSeekSegment
func goGstBaseSrcPrepareSeekSegment(src *C.GstBaseSrc, seek *C.GstEvent, segment *C.GstSegment) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		PrepareSeekSegment(*GstBaseSrc, *gst.Event, *gst.Segment) bool
	})
	return gboolean(caller.PrepareSeekSegment(wrapGstBaseSrc(src), gst.FromGstEventUnsafe(unsafe.Pointer(seek)), gst.FromGstSegmentUnsafe(unsafe.Pointer(segment))))
}

//export goGstBaseSrcDoSeek
func goGstBaseSrcDoSeek(src *C.GstBaseSrc, segment *C.GstSegment) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		DoSeek(*GstBaseSrc, *gst.Segment) bool
	})
	return gboolean(caller.DoSeek(wrapGstBaseSrc(src), gst.FromGstSegmentUnsafe(unsafe.Pointer(segment))))
}

//export goGstBaseSrcUnlock
func goGstBaseSrcUnlock(src *C.GstBaseSrc) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Unlock(*GstBaseSrc) bool
	})
	return gboolean(caller.Unlock(wrapGstBaseSrc(src)))
}

//export goGstBaseSrcUnlockStop
func goGstBaseSrcUnlockStop(src *C.GstBaseSrc) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		UnlockStop(*GstBaseSrc) bool
	})
	return gboolean(caller.UnlockStop(wrapGstBaseSrc(src)))
}

//export goGstBaseSrcQuery
func goGstBaseSrcQuery(src *C.GstBaseSrc, query *C.GstQuery) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Query(*GstBaseSrc, *gst.Query) bool
	})
	return gboolean(caller.Query(wrapGstBaseSrc(src), gst.FromGstQueryUnsafe(unsafe.Pointer(query))))
}

//export goGstBaseSrcEvent
func goGstBaseSrcEvent(src *C.GstBaseSrc, event *C.GstEvent) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Event(*GstBaseSrc, *gst.Event) bool
	})
	return gboolean(caller.Event(wrapGstBaseSrc(src), gst.FromGstEventUnsafe(unsafe.Pointer(event))))
}

//export goGstBaseSrcCreate
func goGstBaseSrcCreate(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf **C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Create(*GstBaseSrc, uint64, uint) (gst.FlowReturn, *gst.Buffer)
	})
	ret, gobuf := caller.Create(wrapGstBaseSrc(src), uint64(offset), uint(size))
	if ret == gst.FlowOK {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(gobuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseSrcAlloc
func goGstBaseSrcAlloc(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf **C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Alloc(*GstBaseSrc, uint64, uint) (gst.FlowReturn, *gst.Buffer)
	})
	ret, gobuf := caller.Alloc(wrapGstBaseSrc(src), uint64(offset), uint(size))
	if ret == gst.FlowOK {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(gobuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseSrcFill
func goGstBaseSrcFill(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf *C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))
	caller := elem.(interface {
		Fill(*GstBaseSrc, uint64, uint, *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(caller.Fill(wrapGstBaseSrc(src), uint64(offset), uint(size), gst.FromGstBufferUnsafe(unsafe.Pointer(buf))))
}
