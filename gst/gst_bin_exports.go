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
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			AddElement(self *Bin, element *Element) bool
		})
		ret = caller.AddElement(wrapBin(gobj), cbWrapElement(element))
	})
	return gboolean(ret)
}

//export goGstBinDeepElementAdded
func goGstBinDeepElementAdded(bin *C.GstBin, subbin *C.GstBin, child *C.GstElement) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			DeepElementAdded(self *Bin, subbin *Bin, child *Element)
		})
		caller.DeepElementAdded(wrapBin(gobj), cbWrapBin(subbin), cbWrapElement(child))
	})
}

//export goGstBinDeepElementRemoved
func goGstBinDeepElementRemoved(bin *C.GstBin, subbin *C.GstBin, child *C.GstElement) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			DeepElementRemoved(self *Bin, subbin *Bin, child *Element)
		})
		caller.DeepElementRemoved(wrapBin(gobj), cbWrapBin(subbin), cbWrapElement(child))
	})
}

//export goGstBinDoLatency
func goGstBinDoLatency(bin *C.GstBin) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			DoLatency(self *Bin) bool
		})
		ret = caller.DoLatency(wrapBin(gobj))
	})
	return gboolean(ret)
}

//export goGstBinElementAdded
func goGstBinElementAdded(bin *C.GstBin, child *C.GstElement) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			ElementAdded(self *Bin, child *Element)
		})
		caller.ElementAdded(wrapBin(gobj), cbWrapElement(child))
	})
}

//export goGstBinElementRemoved
func goGstBinElementRemoved(bin *C.GstBin, child *C.GstElement) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			ElementRemoved(self *Bin, child *Element)
		})
		caller.ElementRemoved(wrapBin(gobj), cbWrapElement(child))
	})
}

//export goGstBinHandleMessage
func goGstBinHandleMessage(bin *C.GstBin, message *C.GstMessage) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			HandleMessage(self *Bin, msg *Message)
		})
		caller.HandleMessage(wrapBin(gobj), wrapMessage(message))
	})
}

//export goGstBinRemoveElement
func goGstBinRemoveElement(bin *C.GstBin, element *C.GstElement) C.gboolean {
	var ret bool
	glib.WithPointerTransferOriginal(unsafe.Pointer(bin), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		caller := obj.(interface {
			RemoveElement(self *Bin, child *Element) bool
		})
		ret = caller.RemoveElement(wrapBin(gobj), cbWrapElement(element))
	})
	return gboolean(ret)
}
