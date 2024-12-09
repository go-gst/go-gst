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

//export goGstPushSrcAlloc
func goGstPushSrcAlloc(src *C.GstPushSrc, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var outbuf *gst.Buffer
	goPushSrc := ToGstPushSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Alloc(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	})
	ret, outbuf = iface.Alloc(goPushSrc)

	if outbuf != nil {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(outbuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstPushSrcCreate
func goGstPushSrcCreate(src *C.GstPushSrc, buf **C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	var outbuf *gst.Buffer
	goPushSrc := ToGstPushSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Create(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	})
	ret, outbuf = iface.Create(goPushSrc)

	if outbuf != nil {
		C.memcpy(unsafe.Pointer(*buf), unsafe.Pointer(outbuf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstPushSrcFill
func goGstPushSrcFill(src *C.GstPushSrc, buf *C.GstBuffer) C.GstFlowReturn {
	var ret gst.FlowReturn
	goPushSrc := ToGstPushSrc(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(src))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(src))

	iface := subclass.(interface {
		Fill(*GstPushSrc, *gst.Buffer) gst.FlowReturn
	})
	ret = iface.Fill(goPushSrc, gst.ToGstBuffer(unsafe.Pointer(buf)))

	return C.GstFlowReturn(ret)
}
