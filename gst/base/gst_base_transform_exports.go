package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

//export goGstBaseTransformAcceptCaps
func goGstBaseTransformAcceptCaps(self *C.GstBaseTransform, direction C.GstPadDirection, caps *C.GstCaps) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		AcceptCaps(self *GstBaseTransform, direction gst.PadDirection, caps *gst.Caps) bool
	})
	return gboolean(caller.AcceptCaps(wrapGstBaseTransform(self), gst.PadDirection(direction), gst.FromGstCapsUnsafeNone(unsafe.Pointer(caps))))
}

//export goGstBaseTransformBeforeTransform
func goGstBaseTransformBeforeTransform(self *C.GstBaseTransform, buffer *C.GstBuffer) {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		BeforeTransform(self *GstBaseTransform, buffer *gst.Buffer)
	})
	caller.BeforeTransform(wrapGstBaseTransform(self), gst.FromGstBufferUnsafeNone(unsafe.Pointer(buffer)))
}

//export goGstBaseTransformCopyMetadata
func goGstBaseTransformCopyMetadata(self *C.GstBaseTransform, input, outbuf *C.GstBuffer) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		CopyMetadata(self *GstBaseTransform, input, output *gst.Buffer) bool
	})
	return gboolean(caller.CopyMetadata(wrapGstBaseTransform(self), gst.FromGstBufferUnsafeNone(unsafe.Pointer(input)), gst.FromGstBufferUnsafeNone(unsafe.Pointer(outbuf))))
}

//export goGstBaseTransformDecideAllocation
func goGstBaseTransformDecideAllocation(self *C.GstBaseTransform, query *C.GstQuery) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		DecideAllocation(self *GstBaseTransform, query *gst.Query) bool
	})
	return gboolean(caller.DecideAllocation(wrapGstBaseTransform(self), gst.FromGstQueryUnsafeNone(unsafe.Pointer(query))))
}

//export goGstBaseTransformFilterMeta
func goGstBaseTransformFilterMeta(self *C.GstBaseTransform, query *C.GstQuery, api C.GType, params *C.GstStructure) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		FilterMeta(self *GstBaseTransform, query *gst.Query, api glib.Type, params *gst.Structure) bool
	})
	return gboolean(caller.FilterMeta(
		wrapGstBaseTransform(self),
		gst.FromGstQueryUnsafeNone(unsafe.Pointer(query)),
		glib.Type(api),
		gst.FromGstStructureUnsafe(unsafe.Pointer(params)),
	))
}

//export goGstBaseTransformFixateCaps
func goGstBaseTransformFixateCaps(self *C.GstBaseTransform, direction C.GstPadDirection, caps, othercaps *C.GstCaps) *C.GstCaps {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		FixateCaps(self *GstBaseTransform, directon gst.PadDirection, caps *gst.Caps, othercaps *gst.Caps) *gst.Caps
	})

	wrappedCaps := gst.FromGstCapsUnsafeNone(unsafe.Pointer(caps))
	wrappedOther := gst.FromGstCapsUnsafeNone(unsafe.Pointer(othercaps))
	defer wrappedOther.Unref()

	fixated := caller.FixateCaps(wrapGstBaseTransform(self), gst.PadDirection(direction), wrappedCaps, wrappedOther)
	if fixated != nil {
		return (*C.GstCaps)(unsafe.Pointer(fixated.Instance()))
	}
	return nil
}

//export goGstBaseTransformGenerateOutput
func goGstBaseTransformGenerateOutput(self *C.GstBaseTransform, buf **C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		GenerateOutput(self *GstBaseTransform) (gst.FlowReturn, *gst.Buffer)
	})
	ret, out := caller.GenerateOutput(wrapGstBaseTransform(self))
	if out != nil {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(out.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseTransformGetUnitSize
func goGstBaseTransformGetUnitSize(self *C.GstBaseTransform, caps *C.GstCaps, size *C.gsize) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		GetUnitSize(self *GstBaseTransform, caps *gst.Caps) (ok bool, size int64)
	})
	ok, retsize := caller.GetUnitSize(wrapGstBaseTransform(self), gst.FromGstCapsUnsafeNone(unsafe.Pointer(caps)))
	if ok {
		*size = C.gsize(retsize)
	}
	return gboolean(ok)
}

//export goGstBaseTransformPrepareOutputBuffer
func goGstBaseTransformPrepareOutputBuffer(self *C.GstBaseTransform, input *C.GstBuffer, outbuf **C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		PrepareOutputBuffer(self *GstBaseTransform, input *gst.Buffer) (gst.FlowReturn, *gst.Buffer)
	})
	ret, out := caller.PrepareOutputBuffer(wrapGstBaseTransform(self), gst.FromGstBufferUnsafeNone(unsafe.Pointer(input)))
	if out != nil {
		C.memcpy(unsafe.Pointer(*outbuf), unsafe.Pointer(out.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstBaseTransformProposeAllocation
func goGstBaseTransformProposeAllocation(self *C.GstBaseTransform, decide, query *C.GstQuery) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		ProposeAllocation(self *GstBaseTransform, decideQuery, query *gst.Query) bool
	})
	return gboolean(caller.ProposeAllocation(wrapGstBaseTransform(self), gst.FromGstQueryUnsafeNone(unsafe.Pointer(decide)), gst.FromGstQueryUnsafeNone(unsafe.Pointer(query))))
}

//export goGstBaseTransformQuery
func goGstBaseTransformQuery(self *C.GstBaseTransform, direction C.GstPadDirection, query *C.GstQuery) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		Query(self *GstBaseTransform, direction gst.PadDirection, query *gst.Query) bool
	})
	return gboolean(caller.Query(wrapGstBaseTransform(self), gst.PadDirection(direction), gst.FromGstQueryUnsafeNone(unsafe.Pointer(query))))
}

//export goGstBaseTransformSetCaps
func goGstBaseTransformSetCaps(self *C.GstBaseTransform, incaps, outcaps *C.GstCaps) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		SetCaps(self *GstBaseTransform, incaps, outcaps *gst.Caps) bool
	})
	return gboolean(caller.SetCaps(
		wrapGstBaseTransform(self),
		gst.FromGstCapsUnsafeNone(unsafe.Pointer(incaps)),
		gst.FromGstCapsUnsafeNone(unsafe.Pointer(outcaps)),
	))
}

//export goGstBaseTransformSinkEvent
func goGstBaseTransformSinkEvent(self *C.GstBaseTransform, event *C.GstEvent) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		SinkEvent(self *GstBaseTransform, event *gst.Event) bool
	})
	return gboolean(caller.SinkEvent(wrapGstBaseTransform(self), gst.FromGstEventUnsafeNone(unsafe.Pointer(event))))
}

//export goGstBaseTransformSrcEvent
func goGstBaseTransformSrcEvent(self *C.GstBaseTransform, event *C.GstEvent) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		SrcEvent(self *GstBaseTransform, event *gst.Event) bool
	})
	return gboolean(caller.SrcEvent(wrapGstBaseTransform(self), gst.FromGstEventUnsafeNone(unsafe.Pointer(event))))
}

//export goGstBaseTransformStart
func goGstBaseTransformStart(self *C.GstBaseTransform) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		Start(self *GstBaseTransform) bool
	})
	return gboolean(caller.Start(wrapGstBaseTransform(self)))
}

//export goGstBaseTransformStop
func goGstBaseTransformStop(self *C.GstBaseTransform) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		Stop(self *GstBaseTransform) bool
	})
	return gboolean(caller.Stop(wrapGstBaseTransform(self)))
}

//export goGstBaseTransformSubmitInputBuffer
func goGstBaseTransformSubmitInputBuffer(self *C.GstBaseTransform, isDiscont C.gboolean, input *C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		SubmitInputBuffer(self *GstBaseTransform, isDiscont bool, input *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(caller.SubmitInputBuffer(
		wrapGstBaseTransform(self),
		gobool(isDiscont),
		gst.FromGstBufferUnsafeNone(unsafe.Pointer(input)),
	))
}

//export goGstBaseTransformTransform
func goGstBaseTransformTransform(self *C.GstBaseTransform, inbuf, outbuf *C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		Transform(self *GstBaseTransform, inbuf, outbuf *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(caller.Transform(
		wrapGstBaseTransform(self),
		gst.FromGstBufferUnsafeNone(unsafe.Pointer(inbuf)),
		gst.FromGstBufferUnsafeNone(unsafe.Pointer(outbuf)),
	))
}

//export goGstBaseTransformTransformCaps
func goGstBaseTransformTransformCaps(self *C.GstBaseTransform, direction C.GstPadDirection, caps, filter *C.GstCaps) *C.GstCaps {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		TransformCaps(self *GstBaseTransform, direction gst.PadDirection, caps, filter *gst.Caps) *gst.Caps
	})
	out := caller.TransformCaps(
		wrapGstBaseTransform(self),
		gst.PadDirection(direction),
		gst.FromGstCapsUnsafeNone(unsafe.Pointer(caps)),
		gst.FromGstCapsUnsafeNone(unsafe.Pointer(filter)),
	)
	if out == nil {
		return nil
	}
	return (*C.GstCaps)(unsafe.Pointer(out.Instance()))
}

//export goGstBaseTransformTransformIP
func goGstBaseTransformTransformIP(self *C.GstBaseTransform, buf *C.GstBuffer) C.GstFlowReturn {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		TransformIP(self *GstBaseTransform, buf *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(caller.TransformIP(
		wrapGstBaseTransform(self),
		gst.FromGstBufferUnsafeNone(unsafe.Pointer(buf)),
	))
}

//export goGstBaseTransformTransformMeta
func goGstBaseTransformTransformMeta(self *C.GstBaseTransform, outbuf *C.GstBuffer, meta *C.GstMeta, inbuf *C.GstBuffer) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		TransformMeta(self *GstBaseTransform, outbuf *gst.Buffer, meta *gst.Meta, inbuf *gst.Buffer) bool
	})
	return gboolean(caller.TransformMeta(
		wrapGstBaseTransform(self),
		gst.FromGstBufferUnsafeNone(unsafe.Pointer(outbuf)),
		gst.FromGstMetaUnsafe(unsafe.Pointer(meta)),
		gst.FromGstBufferUnsafeNone(unsafe.Pointer(inbuf)),
	))
}

//export goGstBaseTransformTransformSize
func goGstBaseTransformTransformSize(self *C.GstBaseTransform, direction C.GstPadDirection, caps *C.GstCaps, size C.gsize, othercaps *C.GstCaps, outsize *C.gsize) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(self))
	caller := elem.(interface {
		TransformSize(self *GstBaseTransform, direction gst.PadDirection, caps *gst.Caps, size int64, othercaps *gst.Caps) (ok bool, othersize int64)
	})
	ok, othersize := caller.TransformSize(
		wrapGstBaseTransform(self),
		gst.PadDirection(direction),
		gst.FromGstCapsUnsafeNone(unsafe.Pointer(caps)),
		int64(size),
		gst.FromGstCapsUnsafeNone(unsafe.Pointer(othercaps)),
	)
	if ok {
		*outsize = C.gsize(othersize)
	}
	return gboolean(ok)
}
