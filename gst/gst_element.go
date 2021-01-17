package gst

/*
#include "gst.go.h"

extern void goGDestroyNotifyFuncNoRun (gpointer user_data);
extern void goElementCallAsync (GstElement * element, gpointer user_data);

void cgoElementAsyncDestroyNotify (gpointer user_data)
{
	goGDestroyNotifyFuncNoRun(user_data);
}

void cgoElementCallAsync (GstElement * element, gpointer user_data)
{
	goElementCallAsync(element, user_data);
}

gboolean elementParentPostMessage (GstElement * element, GstMessage * message) {
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(element));
	GstElementClass * parent = toGstElementClass(g_type_class_peek_parent(this_class));
	return parent->post_message(element, message);
}

GstStateChangeReturn elementParentChangeState (GstElement * element, GstStateChange transition)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(element));
	GstElementClass * parent = toGstElementClass(g_type_class_peek_parent(this_class));
	return parent->change_state(element, transition);
}

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
	"fmt"
	"path"
	"runtime"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// Element is a Go wrapper around a GstElement.
type Element struct{ *Object }

// ToElement returns an Element object for the given Object.
func ToElement(obj *Object) *Element { return &Element{Object: obj} }

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

// Rank represents a level of importance when autoplugging elements.
type Rank uint

// For now just a single RankNone is provided
const (
	RankNone Rank = 0
)

// RegisterElement creates a new elementfactory capable of instantiating objects of the given GoElement
// and adds the factory to the plugin. A higher rank means more importance when autoplugging.
func RegisterElement(plugin *Plugin, name string, rank Rank, elem glib.GoObjectSubclass, extends glib.Extendable) bool {
	return gobool(C.gst_element_register(
		plugin.Instance(),
		C.CString(name),
		C.guint(rank),
		C.GType(glib.RegisterGoType(name, elem, extends)),
	))
}

// Instance returns the underlying GstElement instance.
func (e *Element) Instance() *C.GstElement { return C.toGstElement(e.Unsafe()) }

// AbortState aborts the state change of the element. This function is used by elements that do asynchronous state changes
// and find out something is wrong.
func (e *Element) AbortState() { C.gst_element_abort_state(e.Instance()) }

// AddPad adds a pad (link point) to element. pad's parent will be set to element
//
// Pads are automatically activated when added in the PAUSED or PLAYING state.
//
// The pad and the element should be unlocked when calling this function.
//
// This function will emit the pad-added signal on the element.
func (e *Element) AddPad(pad *Pad) bool {
	return gobool(C.gst_element_add_pad(e.Instance(), pad.Instance()))
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

// CallAsync calls f from another thread. This is to be used for cases when a state change has to be performed from a streaming
// thread, directly via SetState or indirectly e.g. via SEEK events.
//
// Calling those functions directly from the streaming thread will cause deadlocks in many situations, as they might involve waiting
// for the streaming thread to shut down from this very streaming thread.
func (e *Element) CallAsync(f func()) {
	ptr := gopointer.Save(f)
	C.gst_element_call_async(
		e.Instance(),
		C.GstElementCallAsyncFunc(C.cgoElementCallAsync),
		(C.gpointer)(unsafe.Pointer(ptr)),
		C.GDestroyNotify(C.cgoElementAsyncDestroyNotify),
	)
}

// ChangeState performs the given transition on this element.
func (e *Element) ChangeState(transition StateChange) StateChangeReturn {
	return StateChangeReturn(C.gst_element_change_state(e.Instance(), C.GstStateChange(transition)))
}

// Connect connects to the given signal on this element, and applies f as the callback. The callback must
// match the signature of the expected callback from the documentation. However, instead of specifying C types
// for arguments specify the go-gst equivalent (e.g. *gst.Element for almost all GstElement derivitives).
//
// This and the Emit() method may get moved down the hierarchy to the Object level at some point, since
func (e *Element) Connect(signal string, f interface{}) (glib.SignalHandle, error) {
	// Elements are sometimes their own type unique from TYPE_ELEMENT. So make sure a type marshaler
	// is registered for whatever this type is. Use the built-in element marshaler.
	if e.TypeFromInstance() != glib.Type(C.GST_TYPE_ELEMENT) {
		glib.RegisterGValueMarshalers([]glib.TypeMarshaler{{T: e.TypeFromInstance(), F: marshalElement}})
	}
	return e.Object.Connect(signal, f, nil)
}

// Emit is a wrapper around g_signal_emitv() and emits the signal specified by the string s to an Object. Arguments to
// callback functions connected to this signal must be specified in args. Emit() returns an interface{} which must be
// type asserted as the Go equivalent type to the return value for native C callback.
//
// Note that this code is unsafe in that the types of values in args are not checked against whether they are suitable
// for the callback.
func (e *Element) Emit(signal string, args ...interface{}) (interface{}, error) {
	// We are wrapping this for the same reason as Connect.
	if e.TypeFromInstance() != glib.Type(C.GST_TYPE_ELEMENT) {
		glib.RegisterGValueMarshalers([]glib.TypeMarshaler{{T: e.TypeFromInstance(), F: marshalElement}})
	}
	return e.Object.Emit(signal, args...)
}

// InfoMessage is a convenience wrapper for posting an info message from inside an element. Only to be used from
// plugins.
func (e *Element) InfoMessage(domain Domain, text string) {
	function, file, line, _ := runtime.Caller(1)
	e.MessageFull(MessageInfo, domain, ErrorCode(0), "", text, path.Base(file), runtime.FuncForPC(function).Name(), line)
}

// WarningMessage is a convenience wrapper for posting a warning message from inside an element. Only to be used from
// plugins.
func (e *Element) WarningMessage(domain Domain, text string) {
	function, file, line, _ := runtime.Caller(1)
	e.MessageFull(MessageWarning, domain, ErrorCode(0), "", text, path.Base(file), runtime.FuncForPC(function).Name(), line)
}

// ErrorMessage is a convenience wrapper for posting an error message from inside an element. Only to be used from
// plugins.
func (e *Element) ErrorMessage(domain Domain, code ErrorCode, text, debug string) {
	function, file, line, _ := runtime.Caller(1)
	e.MessageFull(MessageError, domain, code, text, debug, path.Base(file), runtime.FuncForPC(function).Name(), line)
}

// MessageFull will post an error, warning, or info message on the bus from inside an element. Only to be used
// from plugins.
func (e *Element) MessageFull(msgType MessageType, domain Domain, code ErrorCode, text, debug, file, function string, line int) {
	var cTxt, cDbg unsafe.Pointer
	if text != "" {
		cTxt = unsafe.Pointer(C.CString(text))
	}
	if debug != "" {
		cDbg = unsafe.Pointer(C.CString(debug))
	}
	C.gst_element_message_full(
		e.Instance(),
		C.GstMessageType(msgType),
		domain.toQuark(),
		C.gint(code),
		(*C.gchar)(cTxt),
		(*C.gchar)(cDbg),
		C.CString(file),
		C.CString(function),
		C.gint(line),
	)
}

// GetBus returns the GstBus for retrieving messages from this element. This function returns
// nil unless the element is a Pipeline.
func (e *Element) GetBus() *Bus {
	bus := C.gst_element_get_bus((*C.GstElement)(e.Instance()))
	if bus == nil {
		return nil
	}
	return wrapBus(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bus))})
}

// GetClock returns the Clock for this element. This is the clock as was last set with gst_element_set_clock.
// Elements in a pipeline will only have their clock set when the pipeline is in the PLAYING state.
func (e *Element) GetClock() *Clock {
	cClock := C.gst_element_get_clock((*C.GstElement)(e.Instance()))
	if cClock == nil {
		return nil
	}
	return wrapClock(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(cClock))})
}

// GetFactory returns the factory that created this element. No refcounting is needed.
func (e *Element) GetFactory() *ElementFactory {
	factory := C.gst_element_get_factory((*C.GstElement)(e.Instance()))
	if factory == nil {
		return nil
	}
	return wrapElementFactory(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(factory))})
}

// GetPads retrieves a list of pads associated with the element.
func (e *Element) GetPads() []*Pad {
	goList := glib.WrapList(uintptr(unsafe.Pointer(e.Instance().pads)))
	out := make([]*Pad, 0)
	goList.Foreach(func(item interface{}) {
		pt := item.(unsafe.Pointer)
		out = append(out, wrapPad(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(pt))}))
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
		out = append(out, wrapPadTemplate(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(pt))}))
	})
	return out
}

// GetState returns the current state of this element.
func (e *Element) GetState() State {
	return State(e.Instance().current_state)
}

// GetStaticPad retrieves a pad from element by name. This version only retrieves
// already-existing (i.e. 'static') pads. The caller owns a ref on the pad and
// should Unref after usage.
func (e *Element) GetStaticPad(name string) *Pad {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	pad := C.gst_element_get_static_pad(e.Instance(), (*C.gchar)(unsafe.Pointer(cname)))
	if pad == nil {
		return nil
	}
	return wrapPad(toGObject(unsafe.Pointer(pad)))
}

// Has returns true if this element has the given flags.
func (e *Element) Has(flags ElementFlags) bool {
	return gobool(C.gstObjectFlagIsSet(C.toGstObject(e.Unsafe()), C.GstElementFlags(flags)))
}

// IsURIHandler returns true if this element can handle URIs.
func (e *Element) IsURIHandler() bool {
	return gobool(C.gstElementIsURIHandler(e.Instance()))
}

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

// ParentChangeState can be used when extending an element to chain up to the parents ChangeState
// handler.
func (e *Element) ParentChangeState(transition StateChange) StateChangeReturn {
	return StateChangeReturn(C.elementParentChangeState(e.Instance(), C.GstStateChange(transition)))
}

// ParentPostMessage can be used when extending an element. During a PostMessage, use this method
// to have the message posted on the bus after processing.
func (e *Element) ParentPostMessage(msg *Message) bool {
	return gobool(C.elementParentPostMessage(e.Instance(), msg.Instance()))
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

// SetState sets the target state for this element.
func (e *Element) SetState(state State) error {
	stateRet := C.gst_element_set_state((*C.GstElement)(e.Instance()), C.GstState(state))
	if stateRet == C.GST_STATE_CHANGE_FAILURE {
		return fmt.Errorf("Failed to change state to %s", state.String())
	}
	return nil
}

// SyncStateWithParent tries to change the state of the element to the same as its parent. If this function returns
// FALSE, the state of element is undefined.
func (e *Element) SyncStateWithParent() bool {
	return gobool(C.gst_element_sync_state_with_parent(e.Instance()))
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

// URIHandler returns a URIHandler interface if implemented by this element. Otherwise it
// returns nil. Currently this only supports elements built through this package, however,
// inner application elements could still use the interface as a reference implementation.
func (e *Element) URIHandler() URIHandler {
	if C.toGstURIHandler(e.Unsafe()) == nil {
		return nil
	}
	return &gstURIHandler{ptr: e.Instance()}
}

// ExtendsElement signifies a GoElement that extends a GstElement.
var ExtendsElement glib.Extendable = &extendElement{parent: glib.ExtendsObject}

// ElementImpl is an interface containing go quivalents of the virtual methods that can be
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
