package base

//#include "gst.go.h"
import "C"

import (
	"time"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

//export goGstBaseSrcGetCaps
func goGstBaseSrcGetCaps(src *C.GstBaseSrc, filter *C.GstCaps) *C.GstCaps {
	var caps *gst.Caps

	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		GetCaps(*GstBaseSrc, *gst.Caps) *gst.Caps
	})
	caps = iface.GetCaps(goBaseSrc, gst.ToGstCaps(unsafe.Pointer(filter)))

	if caps == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(caps.Instance()))
}

//export goGstBaseSrcNegotiate
func goGstBaseSrcNegotiate(src *C.GstBaseSrc) C.gboolean {
	var ret bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Negotiate(*GstBaseSrc) bool
	})
	ret = iface.Negotiate(goBaseSrc)

	return gboolean(ret)
}

//export goGstBaseSrcFixate
func goGstBaseSrcFixate(src *C.GstBaseSrc, filter *C.GstCaps) *C.GstCaps {
	var caps *gst.Caps
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Fixate(*GstBaseSrc, *gst.Caps) *gst.Caps
	})
	caps = iface.Fixate(goBaseSrc, gst.ToGstCaps(unsafe.Pointer(filter)))

	if caps == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(caps.Instance()))
}

//export goGstBaseSrcSetCaps
func goGstBaseSrcSetCaps(src *C.GstBaseSrc, caps *C.GstCaps) C.gboolean {
	var ret bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		SetCaps(*GstBaseSrc, *gst.Caps) bool
	})
	ret = iface.SetCaps(goBaseSrc, gst.ToGstCaps(unsafe.Pointer(caps)))

	return gboolean(ret)
}

//export goGstBaseSrcDecideAllocation
func goGstBaseSrcDecideAllocation(src *C.GstBaseSrc, query *C.GstQuery) C.gboolean {
	var ret bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		DecideAllocation(*GstBaseSrc, *gst.Query) bool
	})
	ret = iface.DecideAllocation(goBaseSrc, gst.ToGstQuery(unsafe.Pointer(query)))

	return gboolean(ret)
}

//export goGstBaseSrcStart
func goGstBaseSrcStart(src *C.GstBaseSrc) C.gboolean {
	var ret bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Start(*GstBaseSrc) bool
	})
	ret = iface.Start(goBaseSrc)

	return gboolean(ret)
}

//export goGstBaseSrcStop
func goGstBaseSrcStop(src *C.GstBaseSrc) C.gboolean {
	var ret bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Stop(*GstBaseSrc) bool
	})
	ret = iface.Stop(goBaseSrc)

	return gboolean(ret)
}

//export goGstBaseSrcGetTimes
func goGstBaseSrcGetTimes(src *C.GstBaseSrc, buf *C.GstBuffer, start *C.GstClockTime, end *C.GstClockTime) {
	var gostart, goend time.Duration
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		GetTimes(*GstBaseSrc, *gst.Buffer) (start, end time.Duration) // should this be a ClockTime?
	})
	gostart, goend = iface.GetTimes(goBaseSrc, gst.ToGstBuffer(unsafe.Pointer(buf)))

	*start = C.GstClockTime(gostart.Nanoseconds())
	*end = C.GstClockTime(goend.Nanoseconds())
}

//export goGstBaseSrcGetSize
func goGstBaseSrcGetSize(src *C.GstBaseSrc, size *C.guint64) C.gboolean {
	var gosize int64
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		GetSize(*GstBaseSrc) (bool, int64)
	})
	ok, gosize = iface.GetSize(goBaseSrc)

	if ok {
		*size = C.guint64(gosize)
	}
	return gboolean(ok)
}

//export goGstBaseSrcIsSeekable
func goGstBaseSrcIsSeekable(src *C.GstBaseSrc) C.gboolean {
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		IsSeekable(*GstBaseSrc) bool
	})
	ok = iface.IsSeekable(goBaseSrc)

	return gboolean(ok)
}

//export goGstBaseSrcPrepareSeekSegment
func goGstBaseSrcPrepareSeekSegment(src *C.GstBaseSrc, seek *C.GstEvent, segment *C.GstSegment) C.gboolean {
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		PrepareSeekSegment(*GstBaseSrc, *gst.Event, *gst.Segment) bool
	})
	ok = iface.PrepareSeekSegment(goBaseSrc, gst.ToGstEvent(unsafe.Pointer(seek)), gst.FromGstSegmentUnsafe(unsafe.Pointer(segment)))

	return gboolean(ok)
}

//export goGstBaseSrcDoSeek
func goGstBaseSrcDoSeek(src *C.GstBaseSrc, segment *C.GstSegment) C.gboolean {
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		DoSeek(*GstBaseSrc, *gst.Segment) bool
	})
	ok = iface.DoSeek(goBaseSrc, gst.ToGstSegment(unsafe.Pointer(segment)))

	return gboolean(ok)
}

//export goGstBaseSrcUnlock
func goGstBaseSrcUnlock(src *C.GstBaseSrc) C.gboolean {
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Unlock(*GstBaseSrc) bool
	})
	ok = iface.Unlock(goBaseSrc)

	return gboolean(ok)
}

//export goGstBaseSrcUnlockStop
func goGstBaseSrcUnlockStop(src *C.GstBaseSrc) C.gboolean {
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		UnlockStop(*GstBaseSrc) bool
	})
	ok = iface.UnlockStop(goBaseSrc)

	return gboolean(ok)
}

//export goGstBaseSrcQuery
func goGstBaseSrcQuery(src *C.GstBaseSrc, query *C.GstQuery) C.gboolean {
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Query(*GstBaseSrc, *gst.Query) bool
	})
	ok = iface.Query(goBaseSrc, gst.ToGstQuery(unsafe.Pointer(query)))

	return gboolean(ok)
}

//export goGstBaseSrcEvent
func goGstBaseSrcEvent(src *C.GstBaseSrc, event *C.GstEvent) C.gboolean {
	var ok bool
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Event(*GstBaseSrc, *gst.Event) bool
	})
	ok = iface.Event(goBaseSrc, gst.ToGstEvent(unsafe.Pointer(event)))

	return gboolean(ok)
}

//export goGstBaseSrcCreate
func goGstBaseSrcCreate(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var gobuf *gst.Buffer
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Create(*GstBaseSrc, uint64, uint) (gst.FlowReturn, *gst.Buffer)
	})
	ret, gobuf = iface.Create(goBaseSrc, uint64(offset), uint(size))

	if ret == gst.FlowOK {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(gobuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseSrcAlloc
func goGstBaseSrcAlloc(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var gobuf *gst.Buffer
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Alloc(*GstBaseSrc, uint64, uint) (gst.FlowReturn, *gst.Buffer)
	})
	ret, gobuf = iface.Alloc(goBaseSrc, uint64(offset), uint(size))

	if ret == gst.FlowOK {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(gobuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseSrcFill
func goGstBaseSrcFill(src *C.GstBaseSrc, offset C.guint64, size C.guint, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseSrc := ToGstBaseSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Fill(*GstBaseSrc, uint64, uint, *gst.Buffer) gst.FlowReturn
	})
	ret = iface.Fill(goBaseSrc, uint64(offset), uint(size), gst.ToGstBuffer(unsafe.Pointer(buf)))

	return C.GstFlowReturn(ret)
}
