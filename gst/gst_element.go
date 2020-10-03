package gst

// #include "gst.go.h"
import "C"

import (
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
func (e *Element) Instance() *C.GstElement { return C.toGstElement(e.Unsafe()) }

// Link wraps gst_element_link and links this element to the given one.
func (e *Element) Link(elem *Element) error {
	if ok := C.gst_element_link((*C.GstElement)(e.Instance()), (*C.GstElement)(elem.Instance())); !gobool(ok) {
		return fmt.Errorf("Failed to link %s to %s", e.Name(), elem.Name())
	}
	return nil
}

// LinkFiltered wraps gst_element_link_filtered and link this element to the given one
// using the provided sink caps.
func (e *Element) LinkFiltered(elem *Element, caps *Caps) error {
	if ok := C.gst_element_link_filtered((*C.GstElement)(e.Instance()), (*C.GstElement)(elem.Instance()), (*C.GstCaps)(caps.Instance())); !gobool(ok) {
		return fmt.Errorf("Failed to link %s to %s with provider caps", e.Name(), elem.Name())
	}
	return nil
}

// GetBus returns the GstBus for retrieving messages from this element. This function returns
// nil unless the element is a Pipeline.
func (e *Element) GetBus() *Bus {
	bus := C.gst_element_get_bus((*C.GstElement)(e.Instance()))
	if bus == nil {
		return nil
	}
	return wrapBus(glib.Take(unsafe.Pointer(bus)))
}

// GetClock returns the Clock for this element. This is the clock as was last set with gst_element_set_clock.
// Elements in a pipeline will only have their clock set when the pipeline is in the PLAYING state.
func (e *Element) GetClock() *Clock {
	cClock := C.gst_element_get_clock((*C.GstElement)(e.Instance()))
	if cClock == nil {
		return nil
	}
	return wrapClock(glib.Take(unsafe.Pointer(cClock)))
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
	if err := e.SetState(state); err != nil {
		return err
	}
	cState := C.GstState(state)
	var curState C.GstState
	C.gst_element_get_state(
		(*C.GstElement)(e.Instance()),
		(*C.GstState)(unsafe.Pointer(&curState)),
		(*C.GstState)(unsafe.Pointer(&cState)),
		C.GstClockTime(ClockTimeNone),
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

// GetStaticPad retrieves a pad from element by name. This version only retrieves
// already-existing (i.e. 'static') pads.
func (e *Element) GetStaticPad(name string) *Pad {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	pad := C.gst_element_get_static_pad(e.Instance(), (*C.gchar)(unsafe.Pointer(cname)))
	if pad == nil {
		return nil
	}
	return wrapPad(toGObject(unsafe.Pointer(pad)))
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

// Has returns true if this element has the given flags.
func (e *Element) Has(flags ElementFlags) bool {
	return gobool(C.gstObjectFlagIsSet(C.toGstObject(e.Unsafe()), C.GstElementFlags(flags)))
}

// IsURIHandler returns true if this element can handle URIs.
func (e *Element) IsURIHandler() bool {
	return gobool(C.gstElementIsURIHandler(e.Instance()))
}

// TODO: Go back over URI and implement as interface

func (e *Element) uriHandler() *C.GstURIHandler { return C.toGstURIHandler(e.Unsafe()) }

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

// TOCSetter returns a TOCSetter interface if implemented by this element. Otherwise it
// returns nil. Currently this only supports elements built through this package, however,
// inner application elements could still use the interface as a reference implementation.
func (e *Element) TOCSetter() TOCSetter {
	if C.toTocSetter(e.Instance()) == nil {
		return nil
	}
	return &gstTOCSetter{ptr: e.Instance()}
}

// TagSetter returns a TagSetter interface if implemented by this element. Otherwise it returns nil.
// This currently only supports elements built through this package's bindings, however, inner application
// elements can still implement the interface themselves if they want.
func (e *Element) TagSetter() TagSetter {
	if C.toTagSetter(e.Instance()) == nil {
		return nil
	}
	return &gstTagSetter{ptr: e.Instance()}
}

// Query performs a query on the given element.
//
// For elements that don't implement a query handler, this function forwards the query to a random srcpad or
// to the peer of a random linked sinkpad of this element.
//
// Please note that some queries might need a running pipeline to work.
func (e *Element) Query(q *Query) bool {
	return gobool(C.gst_element_query(e.Instance(), q.Instance()))
}

// QueryConvert queries an element to convert src_val in src_format to dest_format.
func (e *Element) QueryConvert(srcFormat Format, srcValue int64, destFormat Format) (bool, int64) {
	var out C.gint64
	gok := C.gst_element_query_convert(e.Instance(), C.GstFormat(srcFormat), C.gint64(srcValue), C.GstFormat(destFormat), &out)
	return gobool(gok), int64(out)
}

// QueryDuration queries an element (usually top-level pipeline or playbin element) for the total stream
// duration in nanoseconds. This query will only work once the pipeline is prerolled (i.e. reached PAUSED
// or PLAYING state). The application will receive an ASYNC_DONE message on the pipeline bus when that is
// the case.
//
// If the duration changes for some reason, you will get a DURATION_CHANGED message on the pipeline bus,
// in which case you should re-query the duration using this function.
func (e *Element) QueryDuration(format Format) (bool, int64) {
	var out C.gint64
	gok := C.gst_element_query_duration(e.Instance(), C.GstFormat(format), &out)
	return gobool(gok), int64(out)
}

// QueryPosition queries an element (usually top-level pipeline or playbin element) for the stream position
// in nanoseconds. This will be a value between 0 and the stream duration (if the stream duration is known).
// This query will usually only work once the pipeline is prerolled (i.e. reached PAUSED or PLAYING state).
// The application will receive an ASYNC_DONE message on the pipeline bus when that is the case.
func (e *Element) QueryPosition(format Format) (bool, int64) {
	var out C.gint64
	gok := C.gst_element_query_position(e.Instance(), C.GstFormat(format), &out)
	return gobool(gok), int64(out)
}

// SendEvent sends an event to an element. If the element doesn't implement an event handler, the event will
// be pushed on a random linked sink pad for downstream events or a random linked source pad for upstream events.
//
// This function takes ownership of the provided event so you should gst_event_ref it if you want to reuse the event
// after this call.
func (e *Element) SendEvent(ev *Event) bool {
	return gobool(C.gst_element_send_event(e.Instance(), ev.Instance()))
}
