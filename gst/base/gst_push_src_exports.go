package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

//export goGstPushSrcAlloc
func goGstPushSrcAlloc(src *C.GstPushSrc, buf **C.GstBuffer) C.GstFlowReturn {
	caller := gst.FromObjectUnsafePrivate(unsafe.Pointer(src)).(interface {
		Alloc(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	})
	ret, buffer := caller.Alloc(wrapGstPushSrc(src))
	if ret != gst.FlowOK {
		return C.GstFlowReturn(ret)
	}
	C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(buffer.Instance()), C.sizeof_GstBuffer)
	return C.GstFlowReturn(ret)
}

//export goGstPushSrcCreate
func goGstPushSrcCreate(src *C.GstPushSrc, buf **C.GstBuffer) C.GstFlowReturn {
	caller := gst.FromObjectUnsafePrivate(unsafe.Pointer(src)).(interface {
		Create(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	})
	ret, buffer := caller.Create(wrapGstPushSrc(src))
	if ret != gst.FlowOK {
		return C.GstFlowReturn(ret)
	}
	C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(buffer.Instance()), C.sizeof_GstBuffer)
	return C.GstFlowReturn(ret)
}

//export goGstPushSrcFill
func goGstPushSrcFill(src *C.GstPushSrc, buf *C.GstBuffer) C.GstFlowReturn {
	caller := gst.FromObjectUnsafePrivate(unsafe.Pointer(src)).(interface {
		Fill(*GstPushSrc, *gst.Buffer) gst.FlowReturn
	})
	return C.GstFlowReturn(caller.Fill(wrapGstPushSrc(src), gst.FromGstBufferUnsafe(unsafe.Pointer(buf))))
}
