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

//export goGstPushSrcAlloc
func goGstPushSrcAlloc(src *C.GstPushSrc, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var outbuf *gst.Buffer
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Alloc(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
		})
		ret, outbuf = iface.Alloc(ToGstPushSrc(gObject))
	})
	if outbuf != nil {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(outbuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstPushSrcCreate
func goGstPushSrcCreate(src *C.GstPushSrc, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var outbuf *gst.Buffer
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Create(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
		})
		ret, outbuf = iface.Create(ToGstPushSrc(gObject))
	})
	if outbuf != nil {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(outbuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstPushSrcFill
func goGstPushSrcFill(src *C.GstPushSrc, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(src), func(gObject *glib.Object, goObject glib.GoObjectSubclass) {
		iface := goObject.(interface {
			Fill(*GstPushSrc, *gst.Buffer) gst.FlowReturn
		})
		ret = iface.Fill(ToGstPushSrc(gObject), gst.ToGstBuffer(unsafe.Pointer(buf)))
	})
	return C.GstFlowReturn(ret)
}
