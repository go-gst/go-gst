package gst

import (
	"fmt"
	"runtime"
	"strings"
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

// Type returns the MessageType of the message.
func (message *Message) Type() MessageType {
	return MessageType(message.message.native._type)
}

// Source returns the source object of the message.
func (message *Message) Source() Object {
	// a ref on the message means we already have a ref on the source object. This means we only
	// need to borrow the source object from the message.
	obj := UnsafeObjectFromGlibBorrow(unsafe.Pointer(message.message.native.src))

	// keep the message alive until the object is finalized.
	obj.BorrowFrom(message.message)

	return obj
}

// String implements a stringer on the message. It iterates over the type of the message
// and applies the correct parser, then dumps a string of the basic contents of the
// message. This function can be expensive and should only be used for debugging purposes
// or in routines where latency is not a concern.
//
// This stringer really just helps in keeping track of making sure all message types are
// accounted for in some way. It's the devil, writing it was the devil, and I hope you
// enjoy being able to `fmt.Println(msg)`.
func (m *Message) String() string {
	msg := fmt.Sprintf("[%s] %s - ", m.Source().GetName(), m.Type().String())
	switch m.Type() {

	case MessageEOS:
		msg += "End-of-stream reached in the pipeline"

	case MessageInfo:
		info, err := m.ParseInfo()
		msg += fmt.Sprintf("Info: %s, err: %v", info, err)

	case MessageWarning:
		info, err := m.ParseWarning()
		msg += fmt.Sprintf("Warning: %s, err: %v", info, err)

	case MessageError:
		info, err := m.ParseError()
		msg += fmt.Sprintf("Error: %s, err: %v", info, err)

	case MessageTag:
		tags := m.ParseTag()

		_ = tags // TODO
		msg += "Tags: TODO"

	case MessageBuffering:
		mode, avgIn, avgOut, bufferingLeft := m.ParseBufferingStats()
		msg += fmt.Sprintf(
			"Buffering %s - %d%% complete (avg in %d/sec, avg out %d/sec, time left %d)",
			mode,
			m.ParseBuffering(),
			avgIn,
			avgOut,
			bufferingLeft,
		)

	case MessageStateChanged:
		oldstate, newstate, pending := m.ParseStateChanged()

		msg += fmt.Sprintf("State changed from %s to %s (pending %s)", oldstate, newstate, pending)

	case MessageStateDirty:
		msg += "(DEPRECATED MESSAGE) An element changed state in a streaming thread"

	case MessageStepDone:
		format, amount, rate, flush, intermediate, duration, eos := m.ParseStepDone()

		msg += fmt.Sprintf(
			"Step done with format %s, amount %d, rate %f, flush %v, intermediate %v, duration %d, eos %v",
			format.String(),
			amount,
			rate,
			flush,
			intermediate,
			duration,
			eos,
		)

	case MessageClockProvide:
		msg += "Element has clock provide capability"

	case MessageClockLost:
		msg += "Lost a clock"

	case MessageNewClock:
		clock := m.ParseNewClock()
		msg += fmt.Sprintf("New clock: %s (%s)", clock.GetName(), clock.GoValueType())

	case MessageStructureChange:
		chgType, elem, busy := m.ParseStructureChange()
		msg += fmt.Sprintf("Structure change of type %s from %s. (in progress: %v)", chgType.String(), elem.GetName(), busy)

	case MessageStreamStatus:
		statusType, elem := m.ParseStreamStatus()
		msg += fmt.Sprintf("Stream status from %s: %s", elem.GetName(), statusType.String())

	case MessageApplication:
		msg += "Message posted by the application, possibly via an application-specific element."

	case MessageElement:
		msg += "Internal element message posted"

	case MessageSegmentStart:
		format, pos := m.ParseSegmentStart()
		msg += fmt.Sprintf("Segment started at %d %s", pos, format.String())

	case MessageSegmentDone:
		format, pos := m.ParseSegmentDone()
		msg += fmt.Sprintf("Segment started at %d %s", pos, format.String())

	case MessageDurationChanged:
		msg += "The duration of the pipeline changed"

	case MessageLatency:
		msg += "Element's latency has changed"

	case MessageAsyncStart:
		msg += "Async task started"

	case MessageAsyncDone:
		msg += "Async task completed"
		if dur := m.ParseAsyncDone(); dur > 0 {
			msg += fmt.Sprintf(" in %s", dur)
		}

	case MessageRequestState:
		msg += fmt.Sprintf("State change request to %s", m.ParseRequestState().String())

	case MessageStepStart:
		active, format, amount, rate, flush, intermediate := m.ParseStepStart()

		msg += fmt.Sprintf("Step started with active %v, format %s, amount %d, rate %f, flush %v, intermediate %v",
			active, format.String(), amount, rate, flush, intermediate,
		)

	case MessageQos:
		format, processed, dropped := m.ParseQosStats()

		msg += fmt.Sprintf("Qos stats: format %s, processed %d, dropped %d", format.String(), processed, dropped)

	case MessageProgress:
		progressType, code, text := m.ParseProgress()
		msg += fmt.Sprintf("%s - %s - %s", strings.ToUpper(progressType.String()), code, text)

	case MessageToc:
		toc, updated := m.ParseToc()
		_ = toc // TODO: also show some info about the toc
		msg += fmt.Sprintf("Message toc updated: %t", updated)

	case MessageResetTime:
		msg += fmt.Sprintf("Running time: %s", m.ParseResetTime())

	case MessageStreamStart:
		msg += "Pipeline stream is starting"

	case MessageNeedContext:
		msg += "Element needs context"

	case MessageHaveContext:
		ctx := m.ParseHaveContext()
		_ = ctx // TODO
		msg += fmt.Sprintf("Received context of type %s", "ctx.GetType()")

	case MessageExtended:
		msg += "Extended message type"

	case MessageDeviceAdded:
		device := m.ParseDeviceAdded()
		msg += fmt.Sprintf("Device %s added", device.GetDisplayName())

	case MessageDeviceRemoved:
		device := m.ParseDeviceRemoved()
		msg += fmt.Sprintf("Device %s removed", device.GetDisplayName())

	case MessageDeviceChanged:
		device, _ := m.ParseDeviceChanged()
		msg += fmt.Sprintf("Device %s had its properties updated", device.GetDisplayName())

	case MessagePropertyNotify:
		obj, propName, propVal := m.ParsePropertyNotify()

		msg += fmt.Sprintf("Object %s had property '%s' changed to %+v", obj.GetName(), propName, propVal)

	case MessageStreamCollection:
		collection := m.ParseStreamCollection()
		msg += fmt.Sprintf("New stream collection with upstream id: %s", collection.GetUpstreamID())

	case MessageStreamsSelected:
		collection := m.ParseStreamsSelected()
		msg += fmt.Sprintf("Stream with upstream id '%s' has selected new streams", collection.GetUpstreamID())

	case MessageRedirect:
		msg += fmt.Sprintf("Received redirect message with %d entries", m.GetNumRedirectEntries())

	case MessageUnknown:
		msg += "Unknown message type"

	case MessageAny:
		msg += "Message did not match any known types"
	}
	return msg
}
