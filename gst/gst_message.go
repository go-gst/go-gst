package gst

// #include "gst.go.h"
import "C"

import (
	"strings"
	"unsafe"
)

// Message is a Go wrapper around a GstMessage. It provides convenience methods for
// unref-ing and parsing the underlying messages.
type Message struct {
	msg *C.GstMessage
}

// Instance returns the underlying GstMessage object.
func (m *Message) Instance() *C.GstMessage { return C.toGstMessage(unsafe.Pointer(m.msg)) }

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

// parseToError returns a new GoGError from this message instance. There are multiple
// message types that parse to this interface.
func (m *Message) parseToError() *GoGError {
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
	return &GoGError{
		errMsg:    C.GoString(gerr.message),
		structure: m.GetStructure(),
		debugStr:  strings.TrimSpace(C.GoString((*C.gchar)(debugInfo))),
	}
}

// ParseInfo is identical to ParseError. The returned types are the same. However,
// this is intended for use with GstMessageType `GST_MESSAGE_INFO`.
func (m *Message) ParseInfo() *GoGError {
	return m.parseToError()
}

// ParseWarning is identical to ParseError. The returned types are the same. However,
// this is intended for use with GstMessageType `GST_MESSAGE_WARNING`.
func (m *Message) ParseWarning() *GoGError {
	return m.parseToError()
}

// ParseError will return a GoGError from the contents of this message. This will only work
// if the GstMessageType is `GST_MESSAGE_ERROR`.
func (m *Message) ParseError() *GoGError {
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

// GoGError is a Go wrapper for a C GError. It implements the error interface
// and provides additional functions for retrieving debug strings and details.
type GoGError struct {
	errMsg, debugStr string
	structure        *Structure
}

// Message is an alias to `Error()`. It's for clarity when this object
// is parsed from a `GST_MESSAGE_INFO` or `GST_MESSAGE_WARNING`.
func (e *GoGError) Message() string { return e.Error() }

// Error implements the error interface and returns the error message.
func (e *GoGError) Error() string { return e.errMsg }

// DebugString returns any debug info alongside the error.
func (e *GoGError) DebugString() string { return e.debugStr }

// Structure returns the structure of the error message which may contain additional metadata.
func (e *GoGError) Structure() *Structure { return e.structure }
