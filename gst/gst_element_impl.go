package gst

/*
#include "gst.go.h"

extern GstStateChangeReturn  goGstElementClassChangeState    (GstElement * element, GstStateChange change);
extern GstStateChangeReturn  goGstElementClassGetState       (GstElement * element, GstState * state, GstState * pending, GstClockTime timeout);
extern void                  goGstElementClassNoMorePads     (GstElement * element);
extern void                  goGstElementClassPadAdded       (GstElement * element, GstPad * pad);
extern void                  goGstElementClassPadRemoved     (GstElement * element, GstPad * pad);
extern gboolean              goGstElementClassPostMessage    (GstElement * element, GstMessage * msg);
extern GstClock *            goGstElementClassProvideClock   (GstElement * element);
extern gboolean              goGstElementClassQuery          (GstElement * element, GstQuery * query);
extern void                  goGstElementClassReleasePad     (GstElement * element, GstPad * pad);
extern GstPad *              goGstElementClassRequestNewPad  (GstElement * element, GstPadTemplate * templ, const gchar * name, const GstCaps * caps);
extern gboolean              goGstElementClassSendEvent      (GstElement * element, GstEvent * event);
extern void                  goGstElementClassSetBus         (GstElement * element, GstBus * bus);
extern gboolean              goGstElementClassSetClock       (GstElement * element, GstClock * clock);
extern void                  goGstElementClassSetContext     (GstElement * element, GstContext * ctx);
extern GstStateChangeReturn  goGstElementClassSetState       (GstElement * element, GstState state);
extern void                  goGstElementClassStateChanged   (GstElement * element, GstState old, GstState new, GstState pending);

void  setGstElementClassChangeState    (GstElementClass * klass) { klass->change_state = goGstElementClassChangeState; }
void  setGstElementClassGetState       (GstElementClass * klass) { klass->get_state = goGstElementClassGetState; }
void  setGstElementClassNoMorePads     (GstElementClass * klass) { klass->no_more_pads = goGstElementClassNoMorePads; }
void  setGstElementClassPadAdded       (GstElementClass * klass) { klass->pad_added = goGstElementClassPadAdded; }
void  setGstElementClassPadRemoved     (GstElementClass * klass) { klass->pad_removed = goGstElementClassPadRemoved; }
void  setGstElementClassPostMessage    (GstElementClass * klass) { klass->post_message = goGstElementClassPostMessage; }
void  setGstElementClassProvideClock   (GstElementClass * klass) { klass->provide_clock = goGstElementClassProvideClock; }
void  setGstElementClassQuery          (GstElementClass * klass) { klass->query = goGstElementClassQuery; }
void  setGstElementClassReleasePad     (GstElementClass * klass) { klass->release_pad = goGstElementClassReleasePad; }
void  setGstElementClassRequestNewPad  (GstElementClass * klass) { klass->request_new_pad = goGstElementClassRequestNewPad; }
void  setGstElementClassSendEvent      (GstElementClass * klass) { klass->send_event = goGstElementClassSendEvent; }
void  setGstElementClassSetBus         (GstElementClass * klass) { klass->set_bus = goGstElementClassSetBus; }
void  setGstElementClassSetClock       (GstElementClass * klass) { klass->set_clock = goGstElementClassSetClock; }
void  setGstElementClassSetContext     (GstElementClass * klass) { klass->set_context = goGstElementClassSetContext; }
void  setGstElementClassSetState       (GstElementClass * klass) { klass->set_state = goGstElementClassSetState; }
void  setGstElementClassStateChanged   (GstElementClass * klass) { klass->state_changed = goGstElementClassStateChanged; }

*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// ExtendsElement implements an Extendable object based on a GstElement.
var ExtendsElement glib.Extendable = &extendElement{parent: glib.ExtendsObject}

// ElementImpl is an interface containing go equivalents of the virtual methods that can be
// overridden by a plugin extending an Element.
type ElementImpl interface {
	// ChangeState is called by SetState to perform an incremental state change.
	ChangeState(*Element, StateChange) StateChangeReturn
	// GetState should return the states of the element
	GetState(self *Element, timeout time.Duration) (ret StateChangeReturn, current, pending State)
	// NoMorePads is called when there are no more pads on the element.
	NoMorePads(*Element)
	// PadAdded is called to add a pad to the element.
	PadAdded(*Element, *Pad)
	// PadRemoved is called to remove a pad from the element.
	PadRemoved(*Element, *Pad)
	// PostMessage is called when a message is posted to the element. Call Element.ParentPostMessage
	// to have it posted on the bus after processing.
	PostMessage(*Element, *Message) bool
	// ProvideClock is called to retrieve the clock provided by the element.
	ProvideClock(*Element) *Clock
	// Query is called to perform a query on the element.
	Query(*Element, *Query) bool
	// ReleasePad is called when a request pad is to be released.
	ReleasePad(*Element, *Pad)
	// RequestNewPad is called when a new pad is requested from the element.
	RequestNewPad(self *Element, templ *PadTemplate, name string, caps *Caps) *Pad
	// SendEvent is called to send an event to the element.
	SendEvent(*Element, *Event) bool
	// SetBus is called to set the Bus on the element.
	SetBus(*Element, *Bus)
	// SetClock is called to set the clock on the element.
	SetClock(*Element, *Clock) bool
	// SetContext is called to set the Context on the element.
	SetContext(*Element, *Context)
	// SetState is called to set a new state on the element.
	SetState(*Element, State) StateChangeReturn
	// StateChanged is called immediately after a new state was set.
	StateChanged(self *Element, old, new, pending State)
}

type extendElement struct{ parent glib.Extendable }

func (e *extendElement) Type() glib.Type     { return glib.Type(C.gst_element_get_type()) }
func (e *extendElement) ClassSize() int64    { return int64(C.sizeof_GstElementClass) }
func (e *extendElement) InstanceSize() int64 { return int64(C.sizeof_GstElement) }

func (e *extendElement) InitClass(klass unsafe.Pointer, elem glib.GoObjectSubclass) {
	e.parent.InitClass(klass, elem)

	elemClass := C.toGstElementClass(klass)

	if _, ok := elem.(interface {
		ChangeState(*Element, StateChange) StateChangeReturn
	}); ok {
		C.setGstElementClassChangeState(elemClass)
	}

	if _, ok := elem.(interface {
		GetState(self *Element, timeout time.Duration) (ret StateChangeReturn, current, pending State)
	}); ok {
		C.setGstElementClassGetState(elemClass)
	}

	if _, ok := elem.(interface {
		NoMorePads(*Element)
	}); ok {
		C.setGstElementClassNoMorePads(elemClass)
	}

	if _, ok := elem.(interface {
		PadAdded(*Element, *Pad)
	}); ok {
		C.setGstElementClassPadAdded(elemClass)
	}

	if _, ok := elem.(interface {
		PadRemoved(*Element, *Pad)
	}); ok {
		C.setGstElementClassPadRemoved(elemClass)
	}

	if _, ok := elem.(interface {
		PostMessage(*Element, *Message) bool
	}); ok {
		C.setGstElementClassPostMessage(elemClass)
	}

	if _, ok := elem.(interface {
		ProvideClock(*Element) *Clock
	}); ok {
		C.setGstElementClassProvideClock(elemClass)
	}

	if _, ok := elem.(interface {
		Query(*Element, *Query) bool
	}); ok {
		C.setGstElementClassQuery(elemClass)
	}

	if _, ok := elem.(interface {
		ReleasePad(*Element, *Pad)
	}); ok {
		C.setGstElementClassReleasePad(elemClass)
	}

	if _, ok := elem.(interface {
		RequestNewPad(self *Element, templ *PadTemplate, name string, caps *Caps) *Pad
	}); ok {
		C.setGstElementClassRequestNewPad(elemClass)
	}

	if _, ok := elem.(interface {
		SendEvent(*Element, *Event) bool
	}); ok {
		C.setGstElementClassSendEvent(elemClass)
	}

	if _, ok := elem.(interface {
		SetBus(*Element, *Bus)
	}); ok {
		C.setGstElementClassSetBus(elemClass)
	}

	if _, ok := elem.(interface {
		SetClock(*Element, *Clock) bool
	}); ok {
		C.setGstElementClassSetClock(elemClass)
	}

	if _, ok := elem.(interface {
		SetContext(*Element, *Context)
	}); ok {
		C.setGstElementClassSetContext(elemClass)
	}

	if _, ok := elem.(interface {
		SetState(*Element, State) StateChangeReturn
	}); ok {
		C.setGstElementClassSetState(elemClass)
	}

	if _, ok := elem.(interface {
		StateChanged(self *Element, old, new, pending State)
	}); ok {
		C.setGstElementClassStateChanged(elemClass)
	}
}
