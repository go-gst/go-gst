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
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		ChangeState(*Element, StateChange) StateChangeReturn
	})
	return C.GstStateChangeReturn(iface.ChangeState(wrapCbElem(elem), StateChange(change)))
}

//export goGstElementClassGetState
func goGstElementClassGetState(elem *C.GstElement, state, pending *C.GstState, timeout C.GstClockTime) C.GstStateChangeReturn {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		GetState(*Element, time.Duration) (ret StateChangeReturn, current, pending State)
	})
	ret, cur, pend := iface.GetState(wrapCbElem(elem), time.Duration(timeout)*time.Nanosecond)
	if ret == StateChangeFailure {
		return C.GstStateChangeReturn(ret)
	}
	*state = C.GstState(cur)
	*pending = C.GstState(pend)
	return C.GstStateChangeReturn(ret)
}

//export goGstElementClassNoMorePads
func goGstElementClassNoMorePads(elem *C.GstElement) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		NoMorePads(*Element)
	})
	iface.NoMorePads(wrapCbElem(elem))
}

//export goGstElementClassPadAdded
func goGstElementClassPadAdded(elem *C.GstElement, pad *C.GstPad) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		PadAdded(*Element, *Pad)
	})
	iface.PadAdded(wrapCbElem(elem), wrapPad(toGObject(unsafe.Pointer(pad))))
}

//export goGstElementClassPadRemoved
func goGstElementClassPadRemoved(elem *C.GstElement, pad *C.GstPad) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		PadRemoved(*Element, *Pad)
	})
	iface.PadRemoved(wrapCbElem(elem), wrapPad(toGObject(unsafe.Pointer(pad))))
}

//export goGstElementClassPostMessage
func goGstElementClassPostMessage(elem *C.GstElement, msg *C.GstMessage) C.gboolean {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		PostMessage(*Element, *Message) bool
	})
	return gboolean(iface.PostMessage(wrapCbElem(elem), wrapMessage(msg)))
}

//export goGstElementClassProvideClock
func goGstElementClassProvideClock(elem *C.GstElement) *C.GstClock {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		ProvideClock(*Element) *Clock
	})
	clock := iface.ProvideClock(wrapCbElem(elem))
	if clock == nil {
		return nil
	}
	return clock.Instance()
}

//export goGstElementClassQuery
func goGstElementClassQuery(elem *C.GstElement, query *C.GstQuery) C.gboolean {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		Query(*Element, *Query) bool
	})
	return gboolean(iface.Query(wrapCbElem(elem), wrapQuery(query)))
}

//export goGstElementClassReleasePad
func goGstElementClassReleasePad(elem *C.GstElement, pad *C.GstPad) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		ReleasePad(*Element, *Pad)
	})
	iface.ReleasePad(wrapCbElem(elem), wrapPad(toGObject(unsafe.Pointer(pad))))
}

//export goGstElementClassRequestNewPad
func goGstElementClassRequestNewPad(elem *C.GstElement, templ *C.GstPadTemplate, name *C.gchar, caps *C.GstCaps) *C.GstPad {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		RequestNewPad(self *Element, templ *PadTemplate, name string, caps *Caps) *Pad
	})
	pad := iface.RequestNewPad(
		wrapCbElem(elem),
		wrapPadTemplate(toGObject(unsafe.Pointer(templ))),
		C.GoString(name),
		wrapCaps(caps),
	)
	if pad == nil {
		return nil
	}
	return pad.Instance()
}

//export goGstElementClassSendEvent
func goGstElementClassSendEvent(elem *C.GstElement, event *C.GstEvent) C.gboolean {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		SendEvent(*Element, *Event) bool
	})
	return gboolean(iface.SendEvent(wrapCbElem(elem), wrapEvent(event)))
}

//export goGstElementClassSetBus
func goGstElementClassSetBus(elem *C.GstElement, bus *C.GstBus) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		SetBus(*Element, *Bus)
	})
	iface.SetBus(wrapCbElem(elem), wrapBus(toGObject(unsafe.Pointer(bus))))
}

//export goGstElementClassSetClock
func goGstElementClassSetClock(elem *C.GstElement, clock *C.GstClock) C.gboolean {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		SetClock(*Element, *Clock) bool
	})
	return gboolean(iface.SetClock(wrapCbElem(elem), wrapClock(toGObject(unsafe.Pointer(clock)))))
}

//export goGstElementClassSetContext
func goGstElementClassSetContext(elem *C.GstElement, ctx *C.GstContext) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		SetContext(*Element, *Context)
	})
	iface.SetContext(wrapCbElem(elem), wrapContext(ctx))
}

//export goGstElementClassSetState
func goGstElementClassSetState(elem *C.GstElement, state C.GstState) C.GstStateChangeReturn {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		SetState(*Element, State) StateChangeReturn
	})
	return C.GstStateChangeReturn(iface.SetState(wrapCbElem(elem), State(state)))
}

//export goGstElementClassStateChanged
func goGstElementClassStateChanged(elem *C.GstElement, old, new, pending C.GstState) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem)).(interface {
		StateChanged(self *Element, old, new, pending State)
	})
	iface.StateChanged(wrapCbElem(elem), State(old), State(new), State(pending))
}

func wrapCbElem(elem *C.GstElement) *Element { return wrapElement(toGObject(unsafe.Pointer(elem))) }
