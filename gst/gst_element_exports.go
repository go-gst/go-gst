package gst

/*
#include "gst.go.h"
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

//export goGstElementClassChangeState
func goGstElementClassChangeState(elem *C.GstElement, change C.GstStateChange) C.GstStateChangeReturn {
	var ret StateChangeReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface {
			ChangeState(*Element, StateChange) StateChangeReturn
		})
		ret = iface.ChangeState(wrapElement(gobj), StateChange(change))
	})
	return C.GstStateChangeReturn(ret)
}

//export goGstElementClassGetState
func goGstElementClassGetState(elem *C.GstElement, state, pending *C.GstState, timeout C.GstClockTime) C.GstStateChangeReturn {
	var ret StateChangeReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface {
			GetState(*Element, time.Duration) (ret StateChangeReturn, current, pending State)
		})
		var cur, pend State
		ret, cur, pend = iface.GetState(wrapElement(gobj), time.Duration(timeout)*time.Nanosecond)
		if ret != StateChangeFailure {
			*state = C.GstState(cur)
			*pending = C.GstState(pend)
		}
	})
	return C.GstStateChangeReturn(ret)
}

//export goGstElementClassNoMorePads
func goGstElementClassNoMorePads(elem *C.GstElement) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ NoMorePads(*Element) })
		iface.NoMorePads(wrapElement(gobj))
	})
}

//export goGstElementClassPadAdded
func goGstElementClassPadAdded(elem *C.GstElement, pad *C.GstPad) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ PadAdded(*Element, *Pad) })
		iface.PadAdded(wrapElement(gobj), wrapPad(toGObject(unsafe.Pointer(pad))))
	})
}

//export goGstElementClassPadRemoved
func goGstElementClassPadRemoved(elem *C.GstElement, pad *C.GstPad) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ PadRemoved(*Element, *Pad) })
		iface.PadRemoved(wrapElement(gobj), wrapPad(toGObject(unsafe.Pointer(pad))))
	})
}

//export goGstElementClassPostMessage
func goGstElementClassPostMessage(elem *C.GstElement, msg *C.GstMessage) C.gboolean {
	var ret C.gboolean
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ PostMessage(*Element, *Message) bool })
		ret = gboolean(iface.PostMessage(wrapElement(gobj), wrapMessage(msg)))
	})
	return ret
}

//export goGstElementClassProvideClock
func goGstElementClassProvideClock(elem *C.GstElement) *C.GstClock {
	var clock *Clock
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ ProvideClock(*Element) *Clock })
		clock = iface.ProvideClock(wrapElement(gobj))
	})
	if clock == nil {
		return nil
	}
	return clock.Instance()
}

//export goGstElementClassQuery
func goGstElementClassQuery(elem *C.GstElement, query *C.GstQuery) C.gboolean {
	var ret C.gboolean
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ Query(*Element, *Query) bool })
		ret = gboolean(iface.Query(wrapElement(gobj), wrapQuery(query)))
	})
	return ret
}

//export goGstElementClassReleasePad
func goGstElementClassReleasePad(elem *C.GstElement, pad *C.GstPad) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ ReleasePad(*Element, *Pad) })
		iface.ReleasePad(wrapElement(gobj), wrapPad(toGObject(unsafe.Pointer(pad))))
	})
}

//export goGstElementClassRequestNewPad
func goGstElementClassRequestNewPad(elem *C.GstElement, templ *C.GstPadTemplate, name *C.gchar, caps *C.GstCaps) *C.GstPad {
	var pad *Pad
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface {
			RequestNewPad(self *Element, templ *PadTemplate, name string, caps *Caps) *Pad
		})
		pad = iface.RequestNewPad(
			wrapElement(gobj),
			wrapPadTemplate(toGObject(unsafe.Pointer(templ))),
			C.GoString(name),
			wrapCaps(caps),
		)
	})
	if pad == nil {
		return nil
	}
	return pad.Instance()
}

//export goGstElementClassSendEvent
func goGstElementClassSendEvent(elem *C.GstElement, event *C.GstEvent) C.gboolean {
	var ret C.gboolean
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ SendEvent(*Element, *Event) bool })
		ret = gboolean(iface.SendEvent(wrapElement(gobj), wrapEvent(event)))
	})
	return ret
}

//export goGstElementClassSetBus
func goGstElementClassSetBus(elem *C.GstElement, bus *C.GstBus) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ SetBus(*Element, *Bus) })
		iface.SetBus(wrapElement(gobj), wrapBus(toGObject(unsafe.Pointer(bus))))
	})
}

//export goGstElementClassSetClock
func goGstElementClassSetClock(elem *C.GstElement, clock *C.GstClock) C.gboolean {
	var ret C.gboolean
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ SetClock(*Element, *Clock) bool })
		ret = gboolean(iface.SetClock(wrapElement(gobj), wrapClock(toGObject(unsafe.Pointer(clock)))))
	})
	return ret
}

//export goGstElementClassSetContext
func goGstElementClassSetContext(elem *C.GstElement, ctx *C.GstContext) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface{ SetContext(*Element, *Context) })
		iface.SetContext(wrapElement(gobj), wrapContext(ctx))
	})
}

//export goGstElementClassSetState
func goGstElementClassSetState(elem *C.GstElement, state C.GstState) C.GstStateChangeReturn {
	var ret C.GstStateChangeReturn
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface {
			SetState(*Element, State) StateChangeReturn
		})
		ret = C.GstStateChangeReturn(iface.SetState(wrapElement(gobj), State(state)))
	})
	return ret
}

//export goGstElementClassStateChanged
func goGstElementClassStateChanged(elem *C.GstElement, old, new, pending C.GstState) {
	glib.WithPointerTransferOriginal(unsafe.Pointer(elem), func(gobj *glib.Object, obj glib.GoObjectSubclass) {
		iface := obj.(interface {
			StateChanged(self *Element, old, new, pending State)
		})
		iface.StateChanged(wrapElement(gobj), State(old), State(new), State(pending))
	})
}
