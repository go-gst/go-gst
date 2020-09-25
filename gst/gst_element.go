package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Element is a Go wrapper around a GstElement.
type Element struct{ *Object }

// ElementLinkMany is a go implementation of `gst_element_link_many` to compensate for
// no variadic functions in cgo.
func ElementLinkMany(elems ...*Element) error {
	for idx, elem := range elems {
		if idx == 0 {
			// skip the first one as the loop always links previous to current
			continue
		}
		if err := elems[idx-1].Link(elem); err != nil {
			return err
		}
	}
	return nil
}

// Instance returns the underlying GstElement instance.
func (e *Element) Instance() *C.GstElement { return C.toGstElement(e.unsafe()) }

// Link wraps gst_element_link and links this element to the given one.
func (e *Element) Link(elem *Element) error {
	if ok := C.gst_element_link((*C.GstElement)(e.Instance()), (*C.GstElement)(elem.Instance())); !gobool(ok) {
		return fmt.Errorf("Failed to link %s to %s", e.Name(), elem.Name())
	}
	return nil
}

// LinkFiltered wraps gst_element_link_filtered and link this element to the given one
// using the provided sink caps.
func (e *Element) LinkFiltered(elem *Element, caps Caps) error {
	if ok := C.gst_element_link_filtered((*C.GstElement)(e.Instance()), (*C.GstElement)(elem.Instance()), (*C.GstCaps)(caps.ToGstCaps())); !gobool(ok) {
		return fmt.Errorf("Failed to link %s to %s with provider caps", e.Name(), elem.Name())
	}
	return nil
}

// GetBus returns the GstBus for retrieving messages from this element.
func (e *Element) GetBus() (*Bus, error) {
	bus := C.gst_element_get_bus((*C.GstElement)(e.Instance()))
	if bus == nil {
		return nil, errors.New("Could not retrieve bus from element")
	}
	return wrapBus(glib.Take(unsafe.Pointer(bus))), nil
}

// GetState returns the current state of this element.
func (e *Element) GetState() State {
	return State(e.Instance().current_state)
}

// SetState sets the target state for this element.
func (e *Element) SetState(state State) error {
	stateRet := C.gst_element_set_state((*C.GstElement)(e.Instance()), C.GstState(state))
	if stateRet == C.GST_STATE_CHANGE_FAILURE {
		return fmt.Errorf("Failed to change state to %s", state.String())
	}
	return nil
}

// BlockSetState is like SetState except it will block until the transition
// is complete.
func (e *Element) BlockSetState(state State) error {
	stateRet := C.gst_element_set_state((*C.GstElement)(e.Instance()), C.GST_STATE_PLAYING)
	if stateRet == C.GST_STATE_CHANGE_FAILURE {
		return fmt.Errorf("Failed to change state to %s", state.String())
	}
	var curState C.GstState
	C.gst_element_get_state(
		(*C.GstElement)(e.Instance()),
		(*C.GstState)(unsafe.Pointer(&curState)),
		(*C.GstState)(unsafe.Pointer(&state)),
		C.GST_CLOCK_TIME_NONE,
	)
	return nil
}

// GetFactory returns the factory that created this element. No refcounting is needed.
func (e *Element) GetFactory() *ElementFactory {
	factory := C.gst_element_get_factory((*C.GstElement)(e.Instance()))
	if factory == nil {
		return nil
	}
	return wrapElementFactory(glib.Take(unsafe.Pointer(factory)))
}

// GetPads retrieves a list of pads associated with the element.
func (e *Element) GetPads() []*Pad {
	goList := glib.WrapList(uintptr(unsafe.Pointer(e.Instance().pads)))
	out := make([]*Pad, 0)
	goList.Foreach(func(item interface{}) {
		pt := item.(unsafe.Pointer)
		out = append(out, wrapPad(glib.Take(pt)))
	})
	return out
}

// GetPadTemplates retrieves a list of the pad templates associated with this element.
// The list must not be modified by the calling code.
func (e *Element) GetPadTemplates() []*PadTemplate {
	glist := C.gst_element_get_pad_template_list((*C.GstElement)(e.Instance()))
	if glist == nil {
		return nil
	}
	goList := glib.WrapList(uintptr(unsafe.Pointer(glist)))
	out := make([]*PadTemplate, 0)
	goList.Foreach(func(item interface{}) {
		pt := item.(unsafe.Pointer)
		out = append(out, wrapPadTemplate(glib.Take(pt)))
	})
	return out
}

// GetClock returns the clock for this element or nil. Unref after usage.
func (e *Element) GetClock() *Clock {
	clock := C.gst_element_get_clock((*C.GstElement)(e.Instance()))
	if clock == nil {
		return nil
	}
	return wrapClock(glib.Take(unsafe.Pointer(clock)))
}

// Has returns true if this element has the given flags.
func (e *Element) Has(flags ElementFlags) bool {
	return gobool(C.gstObjectFlagIsSet(C.toGstObject(e.unsafe()), C.GstElementFlags(flags)))
}

// IsURIHandler returns true if this element can handle URIs.
func (e *Element) IsURIHandler() bool {
	return gobool(C.gstElementIsURIHandler(e.Instance()))
}

func (e *Element) uriHandler() *C.GstURIHandler { return C.toGstURIHandler(e.unsafe()) }

// GetURIType returns the type of URI this element can handle.
func (e *Element) GetURIType() URIType {
	if !e.IsURIHandler() {
		return URIUnknown
	}
	ty := C.gst_uri_handler_get_uri_type((*C.GstURIHandler)(e.uriHandler()))
	return URIType(ty)
}

// GetURIProtocols returns the protocols this element can handle.
func (e *Element) GetURIProtocols() []string {
	if !e.IsURIHandler() {
		return nil
	}
	protocols := C.gst_uri_handler_get_protocols((*C.GstURIHandler)(e.uriHandler()))
	if protocols == nil {
		return nil
	}
	size := C.sizeOfGCharArray(protocols)
	return goStrings(size, protocols)
}
