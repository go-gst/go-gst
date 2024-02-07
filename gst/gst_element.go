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

*/
import "C"

import (
	"fmt"
	"path"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	gopointer "github.com/mattn/go-pointer"
)

// Element is a Go wrapper around a GstElement.
type Element struct{ *Object }

// FromGstElementUnsafeNone wraps the given element with a ref and a finalizer.
func FromGstElementUnsafeNone(elem unsafe.Pointer) *Element {
	if elem == nil {
		return nil
	}
	return &Element{Object: &Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: glib.TransferNone(elem)}}}
}

// FromGstElementUnsafeFull wraps the given element with a finalizer.
func FromGstElementUnsafeFull(elem unsafe.Pointer) *Element {
	if elem == nil {
		return nil
	}
	return &Element{Object: &Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: glib.TransferFull(elem)}}}
}

// ToElement returns an Element object for the given Object. It will work
// on either gst.Object or glib.Object interfaces.
func ToElement(obj interface{}) *Element {
	switch obj := obj.(type) {
	case *Object:
		return &Element{Object: obj}
	case *glib.Object:
		return &Element{Object: &Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: obj}}}
	}
	return nil
}

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

// ElementUnlinkMany is a go implementation of `gst_element_unlink_many` to compensate for
// no variadic functions in cgo.
func ElementUnlinkMany(elems ...*Element) {
	for idx, elem := range elems {
		if idx == 0 {
			// skip the first one as the loop always links previous to current
			continue
		}
		elems[idx-1].Unlink(elem)
	}
}

// RegisterElement creates a new elementfactory capable of instantiating objects of the given GoElement
// and adds the factory to the plugin. A higher rank means more importance when autoplugging.
//
// plugin can also be nil to register a static element
func RegisterElement(plugin *Plugin, name string, rank Rank, elem glib.GoObjectSubclass, extends glib.Extendable, interfaces ...glib.Interface) bool {
	var pluginref *C.GstPlugin

	if plugin != nil {
		pluginref = plugin.Instance()
	}

	return gobool(C.gst_element_register(
		pluginref,
		C.CString(name),
		C.guint(rank),
		C.GType(glib.RegisterGoType(name, elem, extends, interfaces...)),
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

// BlockSetState is a convinience wrapper function for calling SetState and an infinitely blocking GetState
func (e *Element) BlockSetState(state State) error {
	if err := e.SetState(state); err != nil {
		return err
	}

	status, _ := e.GetState(state, ClockTimeNone)

	if status != StateChangeSuccess {
		return fmt.Errorf("failed to change state to %s (got %s)", state, status)
	}

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
// for arguments specify the go-gst equivalent (e.g. *gst.Element for almost all GstElement derivatives).
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

// ContinueState commits the state change of the element and proceed to the next pending state if any. This
// function is used by elements that do asynchronous state changes. The core will normally call this method
// automatically when an element returned GST_STATE_CHANGE_SUCCESS from the state change function.
//
// If after calling this method the element still has not reached the pending state, the next state change is performed.
//
// This method is used internally and should normally not be called by plugins or applications.
//
// This function must be called with STATE_LOCK held.
func (e *Element) ContinueState(ret StateChangeReturn) StateChangeReturn {
	return StateChangeReturn(C.gst_element_continue_state(e.Instance(), C.GstStateChangeReturn(ret)))
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

// Error is a convenience wrapper around ErrorMessage to simply post the provided go error on the bus.
// The domain is assumed to be DomainLibrary and the code is assumed to be LibraryErrorFailed.
func (e *Element) Error(msg string, err error) {
	function, file, line, _ := runtime.Caller(1)
	debugMsg := fmt.Sprintf("%s: %s", msg, err.Error())
	e.MessageFull(MessageError, DomainLibrary, LibraryErrorFailed, err.Error(), debugMsg, path.Base(file), runtime.FuncForPC(function).Name(), line)
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
	var cTxt, cDbg unsafe.Pointer = nil, nil
	if text != "" {
		ctxtstr := C.CString(debug)
		defer C.free(unsafe.Pointer(ctxtstr))
		cTxt = unsafe.Pointer(C.g_strdup((*C.gchar)(unsafe.Pointer(ctxtstr))))
	}
	if debug != "" {
		cdbgstr := C.CString(debug)
		defer C.free(unsafe.Pointer(cdbgstr))
		cDbg = unsafe.Pointer(C.g_strdup((*C.gchar)(unsafe.Pointer(cdbgstr))))
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
	return FromGstBusUnsafeFull(unsafe.Pointer(bus))
}

// GetClock returns the Clock for this element. This is the clock as was last set with gst_element_set_clock.
// Elements in a pipeline will only have their clock set when the pipeline is in the PLAYING state.
func (e *Element) GetClock() *Clock {
	cClock := C.gst_element_get_clock((*C.GstElement)(e.Instance()))
	if cClock == nil {
		return nil
	}
	return FromGstClockUnsafeFull(unsafe.Pointer(cClock))
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
func (e *Element) GetPads() ([]*Pad, error) {
	iter := C.gst_element_iterate_pads(e.Instance())
	if iter == nil {
		return nil, nil
	}
	return iteratorToPadSlice(iter)
}

// GetSinkPads retrieves a list of sink pads associated with the element.
func (e *Element) GetSinkPads() ([]*Pad, error) {
	iter := C.gst_element_iterate_sink_pads(e.Instance())
	if iter == nil {
		return nil, nil
	}
	return iteratorToPadSlice(iter)
}

// GetSrcPads retrieves a list of src pads associated with the element.
func (e *Element) GetSrcPads() ([]*Pad, error) {
	iter := C.gst_element_iterate_src_pads(e.Instance())
	if iter == nil {
		return nil, nil
	}
	return iteratorToPadSlice(iter)
}

// GetPadTemplates retrieves a list of the pad templates associated with this element.
// The list must not be modified by the calling code.
func (e *Element) GetPadTemplates() []*PadTemplate {
	glist := C.gst_element_get_pad_template_list((*C.GstElement)(e.Instance()))
	if glist == nil {
		return nil
	}
	goList := glib.WrapList(unsafe.Pointer(glist))
	out := make([]*PadTemplate, 0)
	goList.Foreach(func(item interface{}) {
		pt := item.(unsafe.Pointer)
		out = append(out, wrapPadTemplate(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(pt))}))
	})
	return out
}

// Gets the state of the element.
//
// For elements that performed an ASYNC state change, as reported by gst_element_set_state, this function will block up
// to the specified timeout value for the state change to complete. If the element completes the state change or goes into
// an error, this function returns immediately with a return value of GST_STATE_CHANGE_SUCCESS or GST_STATE_CHANGE_FAILURE respectively.
//
// For elements that did not return GST_STATE_CHANGE_ASYNC, this function returns the current and pending state immediately.
//
// This function returns GST_STATE_CHANGE_NO_PREROLL if the element successfully changed its state but is not able to provide
// data yet. This mostly happens for live sources that only produce data in GST_STATE_PLAYING. While the state change return is
// equivalent to GST_STATE_CHANGE_SUCCESS, it is returned to the application to signal that some sink elements might not be able
// to complete their state change because an element is not producing data to complete the preroll. When setting the element to
// playing, the preroll will complete and playback will start.
func (e *Element) GetState(state State, timeout ClockTime) (StateChangeReturn, State) {
	pending := C.GstState(state)
	var curState C.GstState
	stateChangeStatus := C.gst_element_get_state(
		(*C.GstElement)(e.Instance()),
		(*C.GstState)(unsafe.Pointer(&curState)),
		(*C.GstState)(unsafe.Pointer(&pending)),
		C.GstClockTime(timeout),
	)

	return StateChangeReturn(stateChangeStatus), State(curState)
}

// this returns the value of the `current_state` member of the element:
//
// the current state of an element
//
// see https://gstreamer.freedesktop.org/documentation/gstreamer/gstelement.html?gi-language=c#members
func (e *Element) GetCurrentState() State {
	return State(e.Instance().current_state)
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
	return FromGstPadUnsafeFull(unsafe.Pointer(pad))
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
		return fmt.Errorf("failed to link %s to %s", e.GetName(), elem.GetName())
	}
	return nil
}

func (e *Element) Unlink(elem *Element) {
	C.gst_element_unlink((*C.GstElement)(e.Instance()), (*C.GstElement)(elem.Instance()))
}

// LinkFiltered wraps gst_element_link_filtered and link this element to the given one
// using the provided sink caps.
func (e *Element) LinkFiltered(elem *Element, filter *Caps) error {
	if filter == nil {
		if ok := C.gst_element_link_filtered(e.Instance(), elem.Instance(), nil); !gobool(ok) {
			return fmt.Errorf("failed to link %s to %s with provided caps", e.GetName(), elem.GetName())
		}
		return nil
	}
	if ok := C.gst_element_link_filtered(e.Instance(), elem.Instance(), filter.Instance()); !gobool(ok) {
		return fmt.Errorf("failed to link %s to %s with provided caps", e.GetName(), elem.GetName())
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
	return gobool(C.gst_element_send_event(e.Instance(), ev.Ref().Instance()))
}

// SetState sets the target state for this element.
func (e *Element) SetState(state State) error {
	stateRet := C.gst_element_set_state((*C.GstElement)(e.Instance()), C.GstState(state))
	if stateRet == C.GST_STATE_CHANGE_FAILURE {
		return fmt.Errorf("failed to change state to %s", state.String())
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

// RemovePad removes pad from element. pad will be destroyed if it has not been referenced elsewhere using gst_object_unparent.
//
// This function is used by plugin developers and should not be used by applications. Pads that were dynamically requested from
// elements with gst_element_request_pad should be released with the gst_element_release_request_pad function instead.
//
// Pads are not automatically deactivated so elements should perform the needed steps to deactivate the pad in case this pad is
// removed in the PAUSED or PLAYING state. See gst_pad_set_active for more information about deactivating pads.
//
// The pad and the element should be unlocked when calling this function.
//
// This function will emit the pad-removed signal on the element.
func (e *Element) RemovePad(pad *Pad) bool {
	return gobool(C.gst_element_remove_pad(e.Instance(), pad.Instance()))
}

// GetRequestPad gets a request pad from the element based on the name of the pad template.
// Unlike static pads, request pads are not created automatically but are only created on demand
// For example, audiomixer has sink template, 'sink_%u', which is used for creating multiple sink pads on demand so that it performs mixing of audio streams by linking multiple upstream elements on it's sink pads created on demand.
// This returns the request pad created on demand. Otherwise, it returns null if failed to create.
func (e *Element) GetRequestPad(name string) *Pad {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	pad := C.gst_element_get_request_pad(e.Instance(), (*C.gchar)(unsafe.Pointer(cname)))
	if pad == nil {
		return nil
	}
	return FromGstPadUnsafeFull(unsafe.Pointer(pad))
}

// ReleaseRequestPad releases request pad
func (e *Element) ReleaseRequestPad(pad *Pad) {
	C.gst_element_release_request_pad(e.Instance(), pad.Instance())
}

// Set the start time of an element. The start time of the element is the running time of the element
// when it last went to the PAUSED state. In READY or after a flushing seek, it is set to 0.
//
// Toplevel elements like GstPipeline will manage the start_time and base_time on its children.
// Setting the start_time to GST_CLOCK_TIME_NONE on such a toplevel element will disable the distribution of the base_time
// to the children and can be useful if the application manages the base_time itself, for example if you want to synchronize
// capture from multiple pipelines, and you can also ensure that the pipelines have the same clock.
//
// MT safe.
func (e *Element) SetStartTime(startTime ClockTime) {
	C.gst_element_set_start_time(e.Instance(), C.GstClockTime(startTime))
}

// Returns the start time of the element. The start time is the running time of the clock when this element was last put to PAUSED.
// Usually the start_time is managed by a toplevel element such as GstPipeline.
// MT safe.
func (e *Element) GetStartTime() ClockTime {
	ctime := C.gst_element_get_start_time(e.Instance())

	return ClockTime(ctime)
}

// Set the base time of an element. The base time is the absolute time of the clock
// when this element was last put to PLAYING. Subtracting the base time from the clock time gives the running time of the element.
func (e *Element) SetBaseTime(startTime ClockTime) {
	C.gst_element_set_base_time(e.Instance(), C.GstClockTime(startTime))
}

// Returns the base time of the element. The base time is the absolute time of the clock
// when this element was last put to PLAYING. Subtracting the base time from the clock time gives the running time of the element.
func (e *Element) GetBaseTime() ClockTime {
	ctime := C.gst_element_get_base_time(e.Instance())

	return ClockTime(ctime)
}

// SeekSimple seeks to the given position in the stream. The element / pipeline should be in the PAUSED or PLAYING state and must be a seekable.
func (e *Element) SeekSimple(position int64, format Format, flag SeekFlags) bool {
	result := C.gst_element_seek_simple(e.Instance(), C.GstFormat(format), C.GstSeekFlags(flag), C.gint64(position))
	return gobool(result)
}

// SeekTime seeks to the given position time in the stream. The element / pipeline should be in the PAUSED or PLAYING state and must be a seekable.
//
// For example, to seek to 40th second of the stream, use:
//
//	pos := int64(time.Duration(40 * time.Second))
//	element.SeekTime(pos, gst.FormatTime, gst.SeekFlagFlush)
//
// to perform a flush seek to the nearest keyframe before the given position.
func (e *Element) SeekTime(position time.Duration, flag SeekFlags) bool {
	return e.SeekSimple(position.Nanoseconds(), FormatTime, flag)
}

// SeekDefault seeks to the given position in the stream. The position is the frame number for video, or sample for audio.
// The element should be in the PAUSED or PLAYING state and must be a seekable.
func (e *Element) SeekDefault(position int64, flag SeekFlags) bool {
	return e.SeekSimple(position, FormatDefault, flag)
}

// this prevents go pointers in cgo when setting a gst.Element to a property
// see (https://github.com/go-gst/go-gst/issues/65)
// ToGValue implements glib.ValueTransformer.
func (e *Element) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(glib.Type(C.GST_TYPE_ELEMENT))
	if err != nil {
		return nil, err
	}
	val.SetInstance(unsafe.Pointer(e.Instance()))
	return val, nil
}
