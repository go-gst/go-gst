package gst

// #include "gst.go.h"
import "C"

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Message is a Go wrapper around a GstMessage. It provides convenience methods for
// unref-ing and parsing the underlying messages.
type Message struct {
	msg *C.GstMessage
}

// String implements a stringer on the message. It iterates over the type of the message
// and applies the correct parser, then dumps a string of the basic contents of the
// message. This function can be expensive and should only be used for debugging purposes
// or in routines where latency is not a concern.
func (m *Message) String() string {
	msg := fmt.Sprintf("[%s] %s - ", m.Source(), strings.ToUpper(m.TypeName()))
	switch m.Type() {
	case MessageEOS:
		msg += "End-of-stream reached in the pipeline"
	case MessageInfo:
		msg += m.parseToError().Message()
	case MessageWarning:
		msg += m.parseToError().Message()
	case MessageError:
		msg += m.parseToError().Message()
	case MessageTag:
		tags := m.ParseTags()
		if tags != nil {
			defer tags.Unref()
			msg += tags.String()
		}
	case MessageBuffering:
	case MessageStateChanged:
		oldState, newState := m.ParseStateChanged()
		msg += fmt.Sprintf("State changed from %s to %s", oldState.String(), newState.String())
	case MessageStateDirty:
	case MessageStepDone:
	case MessageClockProvide:
	case MessageClockLost:
		msg += "Lost a clock"
	case MessageNewClock:
		msg += "Received a new clock"
	case MessageStructureChange:
	case MessageStreamStatus:
		statusType, elem := m.ParseStreamStatus()
		msg += fmt.Sprintf("Stream status from %s: %s", elem.Name(), statusType.String())
	case MessageApplication:
	case MessageElement:
	case MessageSegmentStart:
	case MessageSegmentDone:
	case MessageDurationChanged:
	case MessageLatency:
	case MessageAsyncStart:
	case MessageAsyncDone:
		msg += "Async task completed"
		if dur := m.ParseAsyncDone(); dur > 0 {
			msg += fmt.Sprintf(" in %s", dur.String())
		}
	case MessageRequestState:
	case MessageStepStart:
	case MessageQOS:
	case MessageProgress:
	case MessageTOC:
	case MessageResetTime:
	case MessageStreamStart:
		msg += "Pipeline stream is starting"
	case MessageNeedContext:
	case MessageHaveContext:
	case MessageExtended:
	case MessageDeviceAdded:
	case MessageDeviceRemoved:
	case MessagePropertyNotify:
	case MessageStreamCollection:
	case MessageStreamsSelected:
	case MessageRedirect:
	case MessageDeviceChanged:
	case MessageUnknown:
		msg += "Unknown message type"
	case MessageAny:
		msg += "Message did not match any known types"
	}
	return msg
}

// Instance returns the underlying GstMessage object.
func (m *Message) Instance() *C.GstMessage { return C.toGstMessage(unsafe.Pointer(m.msg)) }

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
	return wrapTagList(tagList)
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
	return StreamStatusType(cStatusType), wrapElement(glib.Take(unsafe.Pointer(cElem)))
}

// ParseAsyncDone extracts the running time from the async task done message.
func (m *Message) ParseAsyncDone() time.Duration {
	var clockTime C.GstClockTime
	C.gst_message_parse_async_done((*C.GstMessage)(m.Instance()), &clockTime)
	return nanosecondsToDuration(uint64(clockTime))
}

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
	return wrapMessage(newNative)
}

// GError is a Go wrapper for a C GError. It implements the error interface
// and provides additional functions for retrieving debug strings and details.
type GError struct {
	errMsg, debugStr string
	structure        *Structure
}

// Message is an alias to `Error()`. It's for clarity when this object
// is parsed from a `GST_MESSAGE_INFO` or `GST_MESSAGE_WARNING`.
func (e *GError) Message() string { return e.Error() }

// Error implements the error interface and returns the error message.
func (e *GError) Error() string { return e.errMsg }

// DebugString returns any debug info alongside the error.
func (e *GError) DebugString() string { return e.debugStr }

// Structure returns the structure of the error message which may contain additional metadata.
func (e *GError) Structure() *Structure { return e.structure }
