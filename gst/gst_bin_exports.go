package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

//export goGstBinAddElement
func goGstBinAddElement(bin *C.GstBin, child *C.GstElement) C.gboolean {
	var ret bool

	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	gochild := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))})

	caller := subclass.(interface {
		AddElement(self *Bin, element *Element) bool
	})
	ret = caller.AddElement(goBin, gochild)

	return gboolean(ret)
}

//export goGstBinDeepElementAdded
func goGstBinDeepElementAdded(bin *C.GstBin, subbin *C.GstBin, child *C.GstElement) {

	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	gosubbin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(subbin))})
	gochild := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))})

	caller := subclass.(interface {
		DeepElementAdded(self *Bin, subbin *Bin, child *Element)
	})
	caller.DeepElementAdded(goBin, gosubbin, gochild)
}

//export goGstBinDeepElementRemoved
func goGstBinDeepElementRemoved(bin *C.GstBin, subbin *C.GstBin, child *C.GstElement) {
	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	gosubbin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(subbin))})
	gochild := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))})

	caller := subclass.(interface {
		DeepElementRemoved(self *Bin, subbin *Bin, child *Element)
	})
	caller.DeepElementRemoved(goBin, gosubbin, gochild)
}

//export goGstBinDoLatency
func goGstBinDoLatency(bin *C.GstBin) C.gboolean {
	var ret bool

	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	caller := subclass.(interface {
		DoLatency(self *Bin) bool
	})
	ret = caller.DoLatency(goBin)

	return gboolean(ret)
}

//export goGstBinElementAdded
func goGstBinElementAdded(bin *C.GstBin, child *C.GstElement) {

	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	gochild := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))})

	caller := subclass.(interface {
		ElementAdded(self *Bin, child *Element)
	})
	caller.ElementAdded(goBin, gochild)
}

//export goGstBinElementRemoved
func goGstBinElementRemoved(bin *C.GstBin, child *C.GstElement) {

	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	gochild := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))})

	caller := subclass.(interface {
		ElementRemoved(self *Bin, child *Element)
	})
	caller.ElementRemoved(goBin, gochild)
}

//export goGstBinHandleMessage
func goGstBinHandleMessage(bin *C.GstBin, message *C.GstMessage) {
	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	caller := subclass.(interface {
		HandleMessage(self *Bin, msg *Message)
	})
	caller.HandleMessage(goBin, wrapMessage(message))
}

//export goGstBinRemoveElement
func goGstBinRemoveElement(bin *C.GstBin, child *C.GstElement) C.gboolean {
	var ret bool

	goBin := wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(bin))

	gochild := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))})

	caller := subclass.(interface {
		RemoveElement(self *Bin, child *Element) bool
	})
	ret = caller.RemoveElement(goBin, gochild)

	return gboolean(ret)
}
