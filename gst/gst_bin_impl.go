package gst

/*
#include "gst.go.h"

extern gboolean  goGstBinAddElement          (GstBin * bin, GstElement * element);
extern void      goGstBinDeepElementAdded    (GstBin * bin, GstBin * subbin, GstElement * child);
extern void      goGstBinDeepElementRemoved  (GstBin * bin, GstBin * subbin, GstElement * child);
extern gboolean  goGstBinDoLatency           (GstBin * bin);
extern void      goGstBinElementAdded        (GstBin * bin, GstElement * child);
extern void      goGstBinElementRemoved      (GstBin * bin, GstElement * child);
extern void      goGstBinHandleMessage       (GstBin * bin, GstMessage * message);
extern gboolean  goGstBinRemoveElement       (GstBin * bin, GstElement * element);

void  setGstBinAddElement           (GstBinClass * klass) { klass->add_element = goGstBinAddElement; };
void  setGstBinDeepElementAdded     (GstBinClass * klass) { klass->deep_element_added = goGstBinDeepElementAdded; };
void  setGstBinDeepElementRemoved   (GstBinClass * klass) { klass->deep_element_removed = goGstBinDeepElementRemoved; };
void  setGstBinDoLatency            (GstBinClass * klass) { klass->do_latency = goGstBinDoLatency; };
void  setGstBinElementAdded         (GstBinClass * klass) { klass->element_added = goGstBinElementAdded; };
void  setGstBinElementRemoved       (GstBinClass * klass) { klass->element_removed = goGstBinElementRemoved; };
void  setGstBinHandleMessage        (GstBinClass * klass) { klass->handle_message = goGstBinHandleMessage; };
void  setGstBinRemoveElement        (GstBinClass * klass) { klass->remove_element = goGstBinRemoveElement; };

*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// ExtendsBin implements an Extendable object based on a GstBin.
var ExtendsBin glib.Extendable = &extendsBin{parent: ExtendsElement}

// BinImpl is the reference interface for Go elements extending a Bin. You only need to
// implement the methods that interest you.
type BinImpl interface {
	AddElement(self *Bin, element *Element) bool
	DeepElementAdded(self *Bin, subbin *Bin, child *Element)
	DeepElementRemoved(self *Bin, subbin *Bin, child *Element)
	DoLatency(self *Bin) bool
	ElementAdded(self *Bin, child *Element)
	ElementRemoved(self *Bin, child *Element)
	HandleMessage(self *Bin, msg *Message)
	RemoveElement(self *Bin, element *Element) bool
}

type extendsBin struct{ parent glib.Extendable }

func (e *extendsBin) Type() glib.Type     { return glib.Type(C.gst_bin_get_type()) }
func (e *extendsBin) ClassSize() int64    { return int64(C.sizeof_GstBinClass) }
func (e *extendsBin) InstanceSize() int64 { return int64(C.sizeof_GstBin) }

func (e *extendsBin) InitClass(klass unsafe.Pointer, elem glib.GoObjectSubclass) {
	e.parent.InitClass(klass, elem)

	class := C.toGstBinClass(klass)

	if _, ok := elem.(interface {
		AddElement(self *Bin, element *Element) bool
	}); ok {
		C.setGstBinAddElement(class)
	}

	if _, ok := elem.(interface {
		DeepElementAdded(self *Bin, subbin *Bin, child *Element)
	}); ok {
		C.setGstBinDeepElementAdded(class)
	}

	if _, ok := elem.(interface {
		DeepElementRemoved(self *Bin, subbin *Bin, child *Element)
	}); ok {
		C.setGstBinDeepElementRemoved(class)
	}

	if _, ok := elem.(interface {
		DoLatency(self *Bin) bool
	}); ok {
		C.setGstBinDoLatency(class)
	}

	if _, ok := elem.(interface {
		ElementAdded(self *Bin, child *Element)
	}); ok {
		C.setGstBinElementAdded(class)
	}

	if _, ok := elem.(interface {
		ElementRemoved(self *Bin, child *Element)
	}); ok {
		C.setGstBinElementRemoved(class)
	}

	if _, ok := elem.(interface {
		HandleMessage(self *Bin, msg *Message)
	}); ok {
		C.setGstBinHandleMessage(class)
	}

	if _, ok := elem.(interface {
		RemoveElement(self *Bin, element *Element) bool
	}); ok {
		C.setGstBinRemoveElement(class)
	}
}
