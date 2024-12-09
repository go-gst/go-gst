package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

//export goGstBaseTransformAcceptCaps
func goGstBaseTransformAcceptCaps(self *C.GstBaseTransform, direction C.GstPadDirection, caps *C.GstCaps) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		AcceptCaps(*GstBaseTransform, gst.PadDirection, *gst.Caps) bool
	})
	ret = iface.AcceptCaps(goBaseT, gst.PadDirection(direction), gst.ToGstCaps(unsafe.Pointer(caps)))

	return gboolean(ret)
}

//export goGstBaseTransformBeforeTransform
func goGstBaseTransformBeforeTransform(self *C.GstBaseTransform, buffer *C.GstBuffer) {
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		BeforeTransform(*GstBaseTransform, *gst.Buffer)
	})
	iface.BeforeTransform(goBaseT, gst.ToGstBuffer(unsafe.Pointer(buffer)))
}

//export goGstBaseTransformCopyMetadata
func goGstBaseTransformCopyMetadata(self *C.GstBaseTransform, input, output *C.GstBuffer) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		CopyMetadata(self *GstBaseTransform, input, output *gst.Buffer) bool
	})
	ret = iface.CopyMetadata(goBaseT, gst.ToGstBuffer(unsafe.Pointer(input)), gst.ToGstBuffer(unsafe.Pointer(output)))
	return gboolean(ret)
}

//export goGstBaseTransformDecideAllocation
func goGstBaseTransformDecideAllocation(self *C.GstBaseTransform, query *C.GstQuery) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		DecideAllocation(self *GstBaseTransform, query *gst.Query) bool
	})
	ret = iface.DecideAllocation(goBaseT, gst.ToGstQuery(unsafe.Pointer(query)))

	return gboolean(ret)
}

//export goGstBaseTransformFilterMeta
func goGstBaseTransformFilterMeta(self *C.GstBaseTransform, query *C.GstQuery, api C.GType, params *C.GstStructure) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		FilterMeta(self *GstBaseTransform, query *gst.Query, api glib.Type, params *gst.Structure) bool
	})
	ret = iface.FilterMeta(goBaseT, gst.ToGstQuery(unsafe.Pointer(query)), glib.Type(api), gst.FromGstStructureUnsafe(unsafe.Pointer(params)))

	return gboolean(ret)
}

//export goGstBaseTransformFixateCaps
func goGstBaseTransformFixateCaps(self *C.GstBaseTransform, direction C.GstPadDirection, caps, othercaps *C.GstCaps) *C.GstCaps {
	var ret *gst.Caps
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		FixateCaps(self *GstBaseTransform, directon gst.PadDirection, caps *gst.Caps, othercaps *gst.Caps) *gst.Caps
	})
	ret = iface.FixateCaps(goBaseT, gst.PadDirection(direction), gst.ToGstCaps(unsafe.Pointer(caps)), gst.ToGstCaps(unsafe.Pointer(othercaps)))

	if ret == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(ret.Instance()))
}

//export goGstBaseTransformGenerateOutput
func goGstBaseTransformGenerateOutput(self *C.GstBaseTransform, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var out *gst.Buffer
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		GenerateOutput(self *GstBaseTransform) (gst.FlowReturn, *gst.Buffer)
	})
	ret, out = iface.GenerateOutput(goBaseT)

	if out != nil {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(out.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseTransformGetUnitSize
func goGstBaseTransformGetUnitSize(self *C.GstBaseTransform, caps *C.GstCaps, size *C.gsize) C.gboolean {
	var ret bool
	var out int64
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		GetUnitSize(self *GstBaseTransform, caps *gst.Caps) (ok bool, size int64)
	})
	ret, out = iface.GetUnitSize(goBaseT, gst.ToGstCaps(unsafe.Pointer(caps)))

	if ret {
		*size = C.gsize(out)
	}
	return gboolean(ret)
}

//export goGstBaseTransformPrepareOutputBuffer
func goGstBaseTransformPrepareOutputBuffer(self *C.GstBaseTransform, input *C.GstBuffer, outbuf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var out *gst.Buffer
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		PrepareOutputBuffer(self *GstBaseTransform, input *gst.Buffer) (gst.FlowReturn, *gst.Buffer)
	})
	ret, out = iface.PrepareOutputBuffer(goBaseT, gst.ToGstBuffer(unsafe.Pointer(input)))

	if out != nil {
		C.memcpy(unsafe.Pointer(*outbuf), unsafe.Pointer(out.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseTransformProposeAllocation
func goGstBaseTransformProposeAllocation(self *C.GstBaseTransform, decide, query *C.GstQuery) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		ProposeAllocation(self *GstBaseTransform, decideQuery, query *gst.Query) bool
	})
	ret = iface.ProposeAllocation(goBaseT, gst.ToGstQuery(unsafe.Pointer(decide)), gst.ToGstQuery(unsafe.Pointer(query)))

	return gboolean(ret)
}

//export goGstBaseTransformQuery
func goGstBaseTransformQuery(self *C.GstBaseTransform, direction C.GstPadDirection, query *C.GstQuery) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		Query(self *GstBaseTransform, direction gst.PadDirection, query *gst.Query) bool
	})
	ret = iface.Query(goBaseT, gst.PadDirection(direction), gst.ToGstQuery(unsafe.Pointer(query)))

	return gboolean(ret)
}

//export goGstBaseTransformSetCaps
func goGstBaseTransformSetCaps(self *C.GstBaseTransform, incaps, outcaps *C.GstCaps) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		SetCaps(self *GstBaseTransform, incaps, outcaps *gst.Caps) bool
	})
	ret = iface.SetCaps(goBaseT, gst.ToGstCaps(unsafe.Pointer(incaps)), gst.ToGstCaps(unsafe.Pointer(outcaps)))

	return gboolean(ret)
}

//export goGstBaseTransformSinkEvent
func goGstBaseTransformSinkEvent(self *C.GstBaseTransform, event *C.GstEvent) C.gboolean {
	var ret bool

	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		SinkEvent(self *GstBaseTransform, event *gst.Event) bool
	})

	ret = iface.SinkEvent(goBaseT, gst.ToGstEvent(unsafe.Pointer(event)))

	return gboolean(ret)
}

//export goGstBaseTransformSrcEvent
func goGstBaseTransformSrcEvent(self *C.GstBaseTransform, event *C.GstEvent) C.gboolean {
	var ret bool

	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		SrcEvent(self *GstBaseTransform, event *gst.Event) bool
	})
	ret = iface.SrcEvent(goBaseT, gst.ToGstEvent(unsafe.Pointer(event)))

	return gboolean(ret)
}

//export goGstBaseTransformStart
func goGstBaseTransformStart(self *C.GstBaseTransform) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		Start(self *GstBaseTransform) bool
	})
	ret = iface.Start(goBaseT)

	return gboolean(ret)
}

//export goGstBaseTransformStop
func goGstBaseTransformStop(self *C.GstBaseTransform) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		Stop(self *GstBaseTransform) bool
	})
	ret = iface.Stop(goBaseT)

	return gboolean(ret)
}

//export goGstBaseTransformSubmitInputBuffer
func goGstBaseTransformSubmitInputBuffer(self *C.GstBaseTransform, isDiscont C.gboolean, input *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		SubmitInputBuffer(self *GstBaseTransform, isDiscont bool, input *gst.Buffer) gst.FlowReturn
	})
	ret = iface.SubmitInputBuffer(goBaseT, gobool(isDiscont), gst.ToGstBuffer(unsafe.Pointer(input)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseTransformTransform
func goGstBaseTransformTransform(self *C.GstBaseTransform, inbuf, outbuf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		Transform(self *GstBaseTransform, inbuf, outbuf *gst.Buffer) gst.FlowReturn
	})
	ret = iface.Transform(goBaseT, gst.ToGstBuffer(unsafe.Pointer(inbuf)), gst.ToGstBuffer(unsafe.Pointer(outbuf)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseTransformTransformCaps
func goGstBaseTransformTransformCaps(self *C.GstBaseTransform, direction C.GstPadDirection, caps, filter *C.GstCaps) *C.GstCaps {
	var ret *gst.Caps
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		TransformCaps(self *GstBaseTransform, direction gst.PadDirection, caps, filter *gst.Caps) *gst.Caps
	})
	ret = iface.TransformCaps(goBaseT, gst.PadDirection(direction), gst.ToGstCaps(unsafe.Pointer(caps)), gst.ToGstCaps(unsafe.Pointer(filter)))

	if ret == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(ret.Instance()))
}

//export goGstBaseTransformTransformIP
func goGstBaseTransformTransformIP(self *C.GstBaseTransform, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		TransformIP(self *GstBaseTransform, buf *gst.Buffer) gst.FlowReturn
	})
	ret = iface.TransformIP(goBaseT, gst.ToGstBuffer(unsafe.Pointer(buf)))

	return C.GstFlowReturn(ret)
}

//export goGstBaseTransformTransformMeta
func goGstBaseTransformTransformMeta(self *C.GstBaseTransform, outbuf *C.GstBuffer, meta *C.GstMeta, inbuf *C.GstBuffer) C.gboolean {
	var ret bool
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		TransformMeta(self *GstBaseTransform, outbuf *gst.Buffer, meta *gst.Meta, inbuf *gst.Buffer) bool
	})
	ret = iface.TransformMeta(goBaseT, gst.ToGstBuffer(unsafe.Pointer(outbuf)), gst.FromGstMetaUnsafe(unsafe.Pointer(meta)), gst.ToGstBuffer(unsafe.Pointer(inbuf)))

	return gboolean(ret)
}

//export goGstBaseTransformTransformSize
func goGstBaseTransformTransformSize(self *C.GstBaseTransform, direction C.GstPadDirection, caps *C.GstCaps, size C.gsize, othercaps *C.GstCaps, outsize *C.gsize) C.gboolean {
	var ret bool
	var othersize int64
	goBaseT := ToGstBaseTransform(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(self))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))

	iface := subclass.(interface {
		TransformSize(self *GstBaseTransform, direction gst.PadDirection, caps *gst.Caps, size int64, othercaps *gst.Caps) (ok bool, othersize int64)
	})
	ret, othersize = iface.TransformSize(goBaseT, gst.PadDirection(direction), gst.ToGstCaps(unsafe.Pointer(caps)), int64(size), gst.ToGstCaps(unsafe.Pointer(othercaps)))

	if ret {
		*outsize = C.gsize(othersize)
	}
	return gboolean(ret)
}
