package gst

/*
#include "gst.go.h"
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

//export goGstElementClassChangeState
func goGstElementClassChangeState(elem *C.GstElement, change C.GstStateChange) C.GstStateChangeReturn {
	var ret StateChangeReturn

	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface {
		ChangeState(*Element, StateChange) StateChangeReturn
	})
	ret = iface.ChangeState(goElem, StateChange(change))

	return C.GstStateChangeReturn(ret)
}

//export goGstElementClassGetState
func goGstElementClassGetState(elem *C.GstElement, state, pending *C.GstState, timeout C.GstClockTime) C.GstStateChangeReturn {
	var ret StateChangeReturn
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface {
		GetState(*Element, time.Duration) (ret StateChangeReturn, current, pending State) // should this be a ClockTime?
	})
	var cur, pend State
	ret, cur, pend = iface.GetState(goElem, time.Duration(timeout)*time.Nanosecond)
	if ret != StateChangeFailure {
		*state = C.GstState(cur)
		*pending = C.GstState(pend)
	}

	return C.GstStateChangeReturn(ret)
}

//export goGstElementClassNoMorePads
func goGstElementClassNoMorePads(elem *C.GstElement) {
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ NoMorePads(*Element) })
	iface.NoMorePads(goElem)
}

//export goGstElementClassPadAdded
func goGstElementClassPadAdded(elem *C.GstElement, pad *C.GstPad) {
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ PadAdded(*Element, *Pad) })
	iface.PadAdded(goElem, wrapPad(toGObject(unsafe.Pointer(pad))))
}

//export goGstElementClassPadRemoved
func goGstElementClassPadRemoved(elem *C.GstElement, pad *C.GstPad) {
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ PadRemoved(*Element, *Pad) })
	iface.PadRemoved(goElem, wrapPad(toGObject(unsafe.Pointer(pad))))
}

//export goGstElementClassPostMessage
func goGstElementClassPostMessage(elem *C.GstElement, msg *C.GstMessage) C.gboolean {
	var ret C.gboolean
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ PostMessage(*Element, *Message) bool })
	ret = gboolean(iface.PostMessage(goElem, wrapMessage(msg)))

	return ret
}

//export goGstElementClassProvideClock
func goGstElementClassProvideClock(elem *C.GstElement) *C.GstClock {
	var clock *Clock
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ ProvideClock(*Element) *Clock })
	clock = iface.ProvideClock(goElem)

	if clock == nil {
		return nil
	}
	return clock.Instance()
}

//export goGstElementClassQuery
func goGstElementClassQuery(elem *C.GstElement, query *C.GstQuery) C.gboolean {
	var ret C.gboolean
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ Query(*Element, *Query) bool })
	ret = gboolean(iface.Query(goElem, wrapQuery(query)))

	return ret
}

//export goGstElementClassReleasePad
func goGstElementClassReleasePad(elem *C.GstElement, pad *C.GstPad) {
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ ReleasePad(*Element, *Pad) })
	iface.ReleasePad(goElem, wrapPad(toGObject(unsafe.Pointer(pad))))
}

//export goGstElementClassRequestNewPad
func goGstElementClassRequestNewPad(elem *C.GstElement, templ *C.GstPadTemplate, name *C.gchar, caps *C.GstCaps) *C.GstPad {
	var pad *Pad

	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface {
		RequestNewPad(self *Element, templ *PadTemplate, name string, caps *Caps) *Pad
	})
	pad = iface.RequestNewPad(
		goElem,
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
	var ret C.gboolean

	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ SendEvent(*Element, *Event) bool })
	ret = gboolean(iface.SendEvent(goElem, wrapEvent(event)))

	return ret
}

//export goGstElementClassSetBus
func goGstElementClassSetBus(elem *C.GstElement, bus *C.GstBus) {
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ SetBus(*Element, *Bus) })
	iface.SetBus(goElem, wrapBus(toGObject(unsafe.Pointer(bus))))
}

//export goGstElementClassSetClock
func goGstElementClassSetClock(elem *C.GstElement, clock *C.GstClock) C.gboolean {
	var ret C.gboolean

	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ SetClock(*Element, *Clock) bool })
	ret = gboolean(iface.SetClock(goElem, wrapClock(toGObject(unsafe.Pointer(clock)))))

	return ret
}

//export goGstElementClassSetContext
func goGstElementClassSetContext(elem *C.GstElement, ctx *C.GstContext) {
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface{ SetContext(*Element, *Context) })
	iface.SetContext(goElem, wrapContext(ctx))
}

//export goGstElementClassSetState
func goGstElementClassSetState(elem *C.GstElement, state C.GstState) C.GstStateChangeReturn {
	var ret C.GstStateChangeReturn

	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface {
		SetState(*Element, State) StateChangeReturn
	})
	ret = C.GstStateChangeReturn(iface.SetState(goElem, State(state)))

	return ret
}

//export goGstElementClassStateChanged
func goGstElementClassStateChanged(elem *C.GstElement, old, new, pending C.GstState) {
	goElem := wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))})
	subclass := glib.FromObjectUnsafePrivate(unsafe.Pointer(elem))

	iface := subclass.(interface {
		StateChanged(self *Element, old, new, pending State)
	})
	iface.StateChanged(goElem, State(old), State(new), State(pending))

}
