package gst

// #include "gst.go.h"
import "C"

import (
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Message is a Go wrapper around a GstMessage. It provides convenience methods for
// unref-ing and parsing the underlying messages.
type Message struct {
	msg *C.GstMessage
}

// FromGstMessageUnsafeNone wraps the given unsafe.Pointer in a message. A ref is taken
// on the message and a runtime finalizer placed on the object.
func FromGstMessageUnsafeNone(msg unsafe.Pointer) *Message {
	gomsg := ToGstMessage(msg)
	gomsg.Ref()
	runtime.SetFinalizer(gomsg, (*Message).Unref)
	return gomsg
}

// FromGstMessageUnsafeFull wraps the given unsafe.Pointer in a message. No ref is taken
// and a finalizer is placed on the resulting object.
func FromGstMessageUnsafeFull(msg unsafe.Pointer) *Message {
	gomsg := ToGstMessage(msg)
	runtime.SetFinalizer(gomsg, (*Message).Unref)
	return gomsg
}

// ToGstMessage converts the given pointer into a Message without affecting the ref count or
// placing finalizers.
func ToGstMessage(msg unsafe.Pointer) *Message { return wrapMessage((*C.GstMessage)(msg)) }

// Instance returns the underlying GstMessage object.
func (m *Message) Instance() *C.GstMessage { return C.toGstMessage(unsafe.Pointer(m.msg)) }

// Unref will call `gst_message_unref` on the underlying GstMessage, freeing it from memory.
func (m *Message) Unref() { C.gst_message_unref((*C.GstMessage)(m.Instance())) }

// Ref will increase the ref count on this message. This increases the total amount of times
// Unref needs to be called before the object is freed from memory. It returns the underlying
// message object for convenience.
func (m *Message) Ref() *Message {
	C.gst_message_ref((*C.GstMessage)(m.Instance()))
	return m
}

// Copy will copy this object into a new Message that can be Unrefed separately.
func (m *Message) Copy() *Message {
	newNative := C.gst_message_copy((*C.GstMessage)(m.Instance()))
	return FromGstMessageUnsafeFull(unsafe.Pointer(newNative))
}

// Source returns the source of the message.
func (m *Message) Source() string { return C.GoString(m.Instance().src.name) }

// Type returns the MessageType of the message.
func (m *Message) Type() MessageType {
	return MessageType(m.Instance()._type)
}

// TypeName returns a Go string of the GstMessageType name.
func (m *Message) TypeName() string {
	return C.GoString(C.gst_message_type_get_name((C.GstMessageType)(m.Type())))
}

// GetStructure returns the GstStructure of this message, using the type of the message
// to determine the method to use. The returned structure must not be freed.
func (m *Message) GetStructure() *Structure {
	var st *C.GstStructure

	switch m.Type() {
	case MessageError:
		C.gst_message_parse_error_details((*C.GstMessage)(m.Instance()), (**C.GstStructure)(unsafe.Pointer(&st)))
	case MessageInfo:
		C.gst_message_parse_info_details((*C.GstMessage)(m.Instance()), (**C.GstStructure)(unsafe.Pointer(&st)))
	case MessageWarning:
		C.gst_message_parse_warning_details((*C.GstMessage)(m.Instance()), (**C.GstStructure)(unsafe.Pointer(&st)))
	}

	// if no structure was returned, immediately return nil
	if st == nil {
		return nil
	}

	// The returned structure must not be freed. Applies to all methods.
	// https://gstreamer.freedesktop.org/documentation/gstreamer/gstmessage.html#gst_message_parse_error_details
	return wrapStructure(st)
}

// parseToError returns a new GError from this message instance. There are multiple
// message types that parse to this interface.
func (m *Message) parseToError() *GError {
	var gerr *C.GError
	var debugInfo *C.gchar

	switch m.Type() {
	case MessageError:
		C.gst_message_parse_error((*C.GstMessage)(m.Instance()), (**C.GError)(unsafe.Pointer(&gerr)), (**C.gchar)(unsafe.Pointer(&debugInfo)))
	case MessageInfo:
		C.gst_message_parse_info((*C.GstMessage)(m.Instance()), (**C.GError)(unsafe.Pointer(&gerr)), (**C.gchar)(unsafe.Pointer(&debugInfo)))
	case MessageWarning:
		C.gst_message_parse_warning((*C.GstMessage)(m.Instance()), (**C.GError)(unsafe.Pointer(&gerr)), (**C.gchar)(unsafe.Pointer(&debugInfo)))
	}

	// if error was nil return immediately
	if gerr == nil {
		return nil
	}

	// cleanup the C error immediately and let the garbage collector
	// take over from here.
	defer C.g_error_free((*C.GError)(gerr))
	defer C.g_free((C.gpointer)(debugInfo))
	return &GError{
		errMsg:    C.GoString(gerr.message),
		structure: m.GetStructure(),
		debugStr:  strings.TrimSpace(C.GoString((*C.gchar)(debugInfo))),
	}
}

// ParseInfo is identical to ParseError. The returned types are the same. However,
// this is intended for use with GstMessageType `GST_MESSAGE_INFO`.
func (m *Message) ParseInfo() *GError {
	return m.parseToError()
}

// ParseWarning is identical to ParseError. The returned types are the same. However,
// this is intended for use with GstMessageType `GST_MESSAGE_WARNING`.
func (m *Message) ParseWarning() *GError {
	return m.parseToError()
}

// ParseError will return a GError from the contents of this message. This will only work
// if the GstMessageType is `GST_MESSAGE_ERROR`.
func (m *Message) ParseError() *GError {
	return m.parseToError()
}

// ParseStateChanged will return the old and new states as Go strings. This will only work
// if the GstMessageType is `GST_MESSAGE_STATE_CHANGED`.
func (m *Message) ParseStateChanged() (oldState, newState State) {
	var gOldState, gNewState C.GstState
	C.gst_message_parse_state_changed((*C.GstMessage)(m.Instance()), (*C.GstState)(unsafe.Pointer(&gOldState)), (*C.GstState)(unsafe.Pointer(&gNewState)), nil)
	oldState = State(gOldState)
	newState = State(gNewState)
	return
}

// ParseTags extracts the tag list from the GstMessage. Tags are copied and should be
// unrefed after usage.
func (m *Message) ParseTags() *TagList {
	var tagList *C.GstTagList
	C.gst_message_parse_tag((*C.GstMessage)(m.Instance()), &tagList)
	if tagList == nil {
		return nil
	}
	return FromGstTagListUnsafeFull(unsafe.Pointer(tagList))
}

// ParseTOC extracts the TOC from the GstMessage. The TOC returned in the output argument is
// a copy; the caller must free it with Unref when done.
func (m *Message) ParseTOC() (toc *TOC, updated bool) {
	var gtoc *C.GstToc
	var gupdated C.gboolean
	C.gst_message_parse_toc(m.Instance(), &gtoc, &gupdated)
	return FromGstTOCUnsafeFull(unsafe.Pointer(gtoc)), gobool(gupdated)
}

// ParseStreamStatus parses the stream status type of the message as well as the element
// that produced it. The element returned should NOT be unrefed.
func (m *Message) ParseStreamStatus() (StreamStatusType, *Element) {
	var cElem *C.GstElement
	var cStatusType C.GstStreamStatusType
	C.gst_message_parse_stream_status(
		(*C.GstMessage)(m.Instance()),
		(*C.GstStreamStatusType)(&cStatusType),
		(**C.GstElement)(&cElem),
	)
	return StreamStatusType(cStatusType), FromGstElementUnsafeNone(unsafe.Pointer(cElem))
}

// ParseAsyncDone extracts the running time from the async task done message.
func (m *Message) ParseAsyncDone() time.Duration {
	var clockTime C.GstClockTime
	C.gst_message_parse_async_done((*C.GstMessage)(m.Instance()), &clockTime)
	return time.Duration(clockTime)
}

// BufferingStats represents the buffering stats as retrieved from a GST_MESSAGE_TYPE_BUFFERING.
type BufferingStats struct {
	// The buffering mode
	BufferingMode BufferingMode
	// The average input rate
	AverageIn int
	// The average output rate
	AverageOut int
	// Amount of time until buffering is complete
	BufferingLeft time.Duration
}

// ParseBuffering extracts the buffering percent from the GstMessage.
func (m *Message) ParseBuffering() int {
	var cInt C.gint
	C.gst_message_parse_buffering((*C.GstMessage)(m.Instance()), &cInt)
	return int(cInt)
}

// ParseBufferingStats extracts the buffering stats values from message.
func (m *Message) ParseBufferingStats() *BufferingStats {
	var mode C.GstBufferingMode
	var avgIn, avgOut C.gint
	var bufLeft C.gint64
	C.gst_message_parse_buffering_stats(
		(*C.GstMessage)(m.Instance()),
		&mode, &avgIn, &avgOut, &bufLeft,
	)
	return &BufferingStats{
		BufferingMode: BufferingMode(mode),
		AverageIn:     int(avgIn),
		AverageOut:    int(avgOut),
		BufferingLeft: time.Duration(int64(bufLeft)) * time.Millisecond,
	}
}

// StepStartValues represents the values inside a StepStart message.
type StepStartValues struct {
	Active       bool
	Format       Format
	Amount       uint64
	Rate         float64
	Flush        bool
	Intermediate bool
}

// ParseStepStart extracts the values for the StepStart message.
func (m *Message) ParseStepStart() *StepStartValues {
	var active, flush, intermediate C.gboolean
	var amount C.guint64
	var rate C.gdouble
	var format C.GstFormat
	C.gst_message_parse_step_start(
		(*C.GstMessage)(m.Instance()),
		&active, &format, &amount, &rate, &flush, &intermediate,
	)
	return &StepStartValues{
		Active:       gobool(active),
		Format:       Format(format),
		Amount:       uint64(amount),
		Rate:         float64(rate),
		Flush:        gobool(flush),
		Intermediate: gobool(intermediate),
	}
}

// StepDoneValues represents the values inside a StepDone message.
type StepDoneValues struct {
	Format       Format
	Amount       uint64
	Rate         float64
	Flush        bool
	Intermediate bool
	Duration     time.Duration
	EOS          bool
}

// ParseStepDone extracts the values for the StepDone message.
func (m *Message) ParseStepDone() *StepDoneValues {
	var format C.GstFormat
	var amount, duration C.guint64
	var rate C.gdouble
	var flush, intermediate, eos C.gboolean
	C.gst_message_parse_step_done(
		(*C.GstMessage)(m.Instance()),
		&format,
		&amount,
		&rate,
		&flush,
		&intermediate,
		&duration,
		&eos,
	)
	return &StepDoneValues{
		Format:       Format(format),
		Amount:       uint64(amount),
		Rate:         float64(rate),
		Flush:        gobool(flush),
		Intermediate: gobool(intermediate),
		Duration:     time.Duration(uint64(duration)) * time.Nanosecond,
		EOS:          gobool(eos),
	}
}

// ParseNewClock parses the new Clock in the message. The clock object returned
// remains valid until the message is freed.
func (m *Message) ParseNewClock() *Clock {
	var clock *C.GstClock
	C.gst_message_parse_new_clock((*C.GstMessage)(m.Instance()), &clock)
	return FromGstClockUnsafeNone(unsafe.Pointer(clock))
}

// ParseClockProvide extracts the clock and ready flag from the GstMessage.
// The clock object returned remains valid until the message is freed.
func (m *Message) ParseClockProvide() (clock *Clock, ready bool) {
	var gclock *C.GstClock
	var gready C.gboolean
	C.gst_message_parse_clock_provide((*C.GstMessage)(m.Instance()), &gclock, &gready)
	return FromGstClockUnsafeNone(unsafe.Pointer(clock)), gobool(gready)
}

// ParseStructureChange extracts the change type and completion status from the GstMessage.
// If the returned bool is true, the change is still in progress.
func (m *Message) ParseStructureChange() (chgType StructureChangeType, owner *Element, busy bool) {
	var gElem *C.GstElement
	var gbusy C.gboolean
	var gchgType C.GstStructureChangeType
	C.gst_message_parse_structure_change(
		(*C.GstMessage)(m.Instance()),
		&gchgType, &gElem, &gbusy,
	)
	return StructureChangeType(gchgType), wrapElement(toGObject(unsafe.Pointer(gElem))), gobool(gbusy)
}

// ParseSegmentStart extracts the position and format of the SegmentStart message.
func (m *Message) ParseSegmentStart() (Format, int64) {
	var format C.GstFormat
	var position C.gint64
	C.gst_message_parse_segment_start((*C.GstMessage)(m.Instance()), &format, &position)
	return Format(format), int64(position)
}

// ParseSegmentDone extracts the position and format of the SegmentDone message.
func (m *Message) ParseSegmentDone() (Format, int64) {
	var format C.GstFormat
	var position C.gint64
	C.gst_message_parse_segment_done((*C.GstMessage)(m.Instance()), &format, &position)
	return Format(format), int64(position)
}

// ParseRequestState parses the requests state from the message.
func (m *Message) ParseRequestState() State {
	var state C.GstState
	C.gst_message_parse_request_state((*C.GstMessage)(m.Instance()), &state)
	return State(state)
}

// QoSValues represents the values inside a QoS message.
type QoSValues struct {
	// If the message was generated by a live element
	Live bool
	// The running time of the buffer that generated the message
	RunningTime time.Duration
	// The stream time of the buffer that generated the message
	StreamTime time.Duration
	//  The timestamps of the buffer that generated the message
	Timestamp time.Duration
	//  The duration of the buffer that generated the message
	Duration time.Duration
}

// ParseQoS extracts the timestamps and live status from the QoS message.
// The values reflect those of the dropped buffer. Values of ClockTimeNone
// or -1 mean unknown values.
func (m *Message) ParseQoS() *QoSValues {
	var live C.gboolean
	var runningTime, streamTime, timestamp, duration C.guint64
	C.gst_message_parse_qos(
		(*C.GstMessage)(m.Instance()),
		&live, &runningTime, &streamTime, &timestamp, &duration,
	)
	return &QoSValues{
		Live:        gobool(live),
		RunningTime: time.Duration(runningTime),
		StreamTime:  time.Duration(streamTime),
		Timestamp:   time.Duration(timestamp),
		Duration:    time.Duration(duration),
	}
}

// ParseProgress parses the progress type, code and text.
func (m *Message) ParseProgress() (progressType ProgressType, code, text string) {
	codePtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(codePtr))
	textPtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(textPtr))
	var gpType C.GstProgressType
	C.gst_message_parse_progress(
		(*C.GstMessage)(m.Instance()),
		&gpType,
		(**C.gchar)(unsafe.Pointer(codePtr)),
		(**C.gchar)(unsafe.Pointer(textPtr)),
	)
	return ProgressType(gpType),
		string(C.GoBytes(codePtr, C.sizeOfGCharArray((**C.gchar)(codePtr)))),
		string(C.GoBytes(textPtr, C.sizeOfGCharArray((**C.gchar)(textPtr))))
}

// ParseResetTime extracts the running-time from the ResetTime message.
func (m *Message) ParseResetTime() time.Duration {
	var clockTime C.GstClockTime
	C.gst_message_parse_reset_time((*C.GstMessage)(m.Instance()), &clockTime)
	return time.Duration(clockTime)
}

// ParseDeviceAdded parses a device-added message. The device-added message is
// produced by GstDeviceProvider or a GstDeviceMonitor. It announces the appearance
// of monitored devices.
func (m *Message) ParseDeviceAdded() *Device {
	var device *C.GstDevice
	C.gst_message_parse_device_added((*C.GstMessage)(m.Instance()), &device)
	return FromGstDeviceUnsafeFull(unsafe.Pointer(device))
}

// ParseDeviceRemoved parses a device-removed message. The device-removed message
// is produced by GstDeviceProvider or a GstDeviceMonitor. It announces the disappearance
// of monitored devices.
func (m *Message) ParseDeviceRemoved() *Device {
	var device *C.GstDevice
	C.gst_message_parse_device_removed((*C.GstMessage)(m.Instance()), &device)
	return FromGstDeviceUnsafeFull(unsafe.Pointer(device))
}

// ParseDeviceChanged Parses a device-changed message. The device-changed message is
// produced by GstDeviceProvider or a GstDeviceMonitor. It announces that a device's properties
// have changed.
// The first device returned is the updated Device, and the second changedDevice represents
// the old state of the device.
func (m *Message) ParseDeviceChanged() (newDevice, oldDevice *Device) {
	var gstNewDevice, gstOldDevice *C.GstDevice
	C.gst_message_parse_device_changed((*C.GstMessage)(m.Instance()), &gstNewDevice, &gstOldDevice)
	return FromGstDeviceUnsafeFull(unsafe.Pointer(gstNewDevice)),
		FromGstDeviceUnsafeFull(unsafe.Pointer(gstOldDevice))
}

// ParsePropertyNotify parses a property-notify message. These will be posted on the bus only
// when set up with Element.AddPropertyNotifyWatch (TODO) or Element.AddPropertyDeepNotifyWatch (TODO).
func (m *Message) ParsePropertyNotify() (obj *Object, propertName string, propertyValue *glib.Value) {
	var gstobj *C.GstObject
	var gval *C.GValue
	namePtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(namePtr))
	C.gst_message_parse_property_notify(
		(*C.GstMessage)(m.Instance()),
		&gstobj, (**C.gchar)(unsafe.Pointer(namePtr)), &gval,
	)
	return wrapObject(toGObject(unsafe.Pointer(gstobj))),
		string(C.GoBytes(namePtr, C.sizeOfGCharArray((**C.gchar)(namePtr)))),
		glib.ValueFromNative(unsafe.Pointer(gval))
}

// ParseStreamCollection parses a stream-collection message.
func (m *Message) ParseStreamCollection() *StreamCollection {
	var collection *C.GstStreamCollection
	C.gst_message_parse_stream_collection(
		(*C.GstMessage)(m.Instance()),
		&collection,
	)
	return FromGstStreamCollectionUnsafeFull(unsafe.Pointer(collection))
}

// ParseStreamsSelected parses a streams-selected message.
func (m *Message) ParseStreamsSelected() *StreamCollection {
	var collection *C.GstStreamCollection
	C.gst_message_parse_streams_selected(
		(*C.GstMessage)(m.Instance()),
		&collection,
	)
	return FromGstStreamCollectionUnsafeFull(unsafe.Pointer(collection))
}

// NumRedirectEntries returns the number of redirect entries in a MessageRedirect.
func (m *Message) NumRedirectEntries() int64 {
	return int64(C.gst_message_get_num_redirect_entries((*C.GstMessage)(m.Instance())))
}

// ParseRedirectEntryAt parses the redirect entry at the given index. Total indices can be retrieved
// with NumRedirectEntries().
func (m *Message) ParseRedirectEntryAt(idx int64) (location string, tags *TagList, structure *Structure) {
	locPtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(locPtr))
	var tagList *C.GstTagList
	var entryStruct *C.GstStructure
	C.gst_message_parse_redirect_entry(
		(*C.GstMessage)(m.Instance()),
		C.gsize(idx),
		(**C.char)(locPtr),
		&tagList,
		&entryStruct,
	)
	return string(C.GoBytes(locPtr, C.sizeOfGCharArray((**C.gchar)(locPtr)))),
		FromGstTagListUnsafeNone(unsafe.Pointer(tagList)), wrapStructure(entryStruct)
}

// ParseHaveContext parses the context from a HaveContext message.
func (m *Message) ParseHaveContext() *Context {
	var ctx *C.GstContext
	C.gst_message_parse_have_context(m.Instance(), &ctx)
	return FromGstContextUnsafeFull(unsafe.Pointer(ctx))
}
