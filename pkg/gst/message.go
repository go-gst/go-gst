package gst

import (
	"runtime"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/gobject/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// NewMessagePropertyNotify wraps gst_message_new_property_notify
//
// The function takes the following parameters:
//
//   - src Object: The #GstObject whose property changed (may or may not be a #GstElement)
//   - propertyName string: name of the property that changed
//   - val any (nullable): new property value, or %NULL
//
// The function returns the following values:
//
//   - goret *Message
func NewMessagePropertyNotify(src Object, propertyName string, val any) *Message {
	var carg1 *C.GstObject // in, none, converted
	var carg2 *C.gchar     // in, none, string, casted *C.gchar
	var carg3 *C.GValue    // in, full, converted, nullable
	var cret *C.GstMessage // return, full, converted

	carg1 = (*C.GstObject)(UnsafeObjectToGlibNone(src))
	carg2 = (*C.gchar)(unsafe.Pointer(C.CString(propertyName)))
	defer C.free(unsafe.Pointer(carg2))
	if val != nil {
		carg3 = (*C.GValue)(gobject.UnsafeValueToGlibFull(gobject.NewValue(val)))
	}

	cret = C.gst_message_new_property_notify(carg1, carg2, carg3)
	runtime.KeepAlive(src)
	runtime.KeepAlive(propertyName)
	runtime.KeepAlive(val)

	var goret *Message

	goret = UnsafeMessageFromGlibFull(unsafe.Pointer(cret))

	return goret
}

// ParsePropertyNotify wraps gst_message_parse_property_notify
// The function returns the following values:
//
//   - object Object: location where to store a
//     pointer to the object whose property got changed, or %NULL
//   - propertyName string: return location for
//     the name of the property that got changed, or %NULL
//   - propertyValue *gobject.Value (nullable): return location for
//     the new value of the property that got changed, or %NULL. This will
//     only be set if the property notify watch was told to include the value
//     when it was set up
//
// Parses a property-notify message. These will be posted on the bus only
// when set up with gst_element_add_property_notify_watch() or
// gst_element_add_property_deep_notify_watch().
func (message *Message) ParsePropertyNotify() (Object, string, any) {
	var carg0 *C.GstMessage // in, none, converted
	var carg1 *C.GstObject  // out, none, converted
	var carg2 *C.gchar      // out, none, string, casted *C.gchar
	var carg3 *C.GValue     // out, none, converted, nullable

	carg0 = (*C.GstMessage)(UnsafeMessageToGlibNone(message))

	C.gst_message_parse_property_notify(carg0, &carg1, &carg2, &carg3)
	runtime.KeepAlive(message)

	var object Object
	var propertyName string
	var propertyValue any

	object = UnsafeObjectFromGlibNone(unsafe.Pointer(carg1))
	propertyName = C.GoString((*C.char)(unsafe.Pointer(carg2)))
	if carg3 != nil {
		propertyValue = gobject.ValueFromNative(unsafe.Pointer(carg3)).GoValue()
	}

	return object, propertyName, propertyValue
}

// GetStreamStatusObject wraps gst_message_get_stream_status_object
// The function returns the following values:
//
//   - goret *gobject.Value
//
// Extracts the object managing the streaming thread from @message.
func (message *Message) GetStreamStatusObject() any {
	var carg0 *C.GstMessage // in, none, converted
	var cret *C.GValue      // return, none, converted

	carg0 = (*C.GstMessage)(UnsafeMessageToGlibNone(message))

	cret = C.gst_message_get_stream_status_object(carg0)
	runtime.KeepAlive(message)

	var goret any

	goret = gobject.ValueFromNative(unsafe.Pointer(cret)).GoValue()

	return goret
}
