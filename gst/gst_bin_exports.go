package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

func cbWrapBin(bin *C.GstBin) *Bin {
	return wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
}

func cbWrapElement(elem *C.GstElement) *Element {
	return wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
}

//export goGstBinAddElement
func goGstBinAddElement(bin *C.GstBin, element *C.GstElement) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		AddElement(self *Bin, element *Element) bool
	})
	return gboolean(caller.AddElement(cbWrapBin(bin), cbWrapElement(element)))
}

//export goGstBinDeepElementAdded
func goGstBinDeepElementAdded(bin *C.GstBin, subbin *C.GstBin, child *C.GstElement) {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		DeepElementAdded(self *Bin, subbin *Bin, child *Element)
	})
	caller.DeepElementAdded(cbWrapBin(bin), cbWrapBin(subbin), cbWrapElement(child))
}

//export goGstBinDeepElementRemoved
func goGstBinDeepElementRemoved(bin *C.GstBin, subbin *C.GstBin, child *C.GstElement) {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		DeepElementRemoved(self *Bin, subbin *Bin, child *Element)
	})
	caller.DeepElementRemoved(cbWrapBin(bin), cbWrapBin(subbin), cbWrapElement(child))
}

//export goGstBinDoLatency
func goGstBinDoLatency(bin *C.GstBin) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		DoLatency(self *Bin) bool
	})
	return gboolean(caller.DoLatency(cbWrapBin(bin)))
}

//export goGstBinElementAdded
func goGstBinElementAdded(bin *C.GstBin, child *C.GstElement) {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		ElementAdded(self *Bin, child *Element)
	})
	caller.ElementAdded(cbWrapBin(bin), cbWrapElement(child))
}

//export goGstBinElementRemoved
func goGstBinElementRemoved(bin *C.GstBin, child *C.GstElement) {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		ElementRemoved(self *Bin, child *Element)
	})
	caller.ElementRemoved(cbWrapBin(bin), cbWrapElement(child))
}

//export goGstBinHandleMessage
func goGstBinHandleMessage(bin *C.GstBin, message *C.GstMessage) {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		HandleMessage(self *Bin, msg *Message)
	})
	caller.HandleMessage(cbWrapBin(bin), wrapMessage(message))
}

//export goGstBinRemoveElement
func goGstBinRemoveElement(bin *C.GstBin, element *C.GstElement) C.gboolean {
	elem := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))
	caller := elem.(interface {
		RemoveElement(self *Bin, element *Element) bool
	})
	return gboolean(caller.RemoveElement(cbWrapBin(bin), cbWrapElement(element)))
}
