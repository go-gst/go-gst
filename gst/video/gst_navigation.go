package video

/*
#include <stdlib.h>
#include <gst/video/video.h>

GstNavigation * toGstNavigation (GstElement * element)
{
	return GST_NAVIGATION(element);
}
*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// NavigationCommand is a set of commands that may be issued to an element providing the
// Navigation interface. The available commands can be queried via the QueryNewCommands
// query.
type NavigationCommand int

// Type castings
const (
	NavigationCommandInvalid   NavigationCommand = C.GST_NAVIGATION_COMMAND_INVALID    // (0) – An invalid command entry
	NavigationCommandMenu1     NavigationCommand = C.GST_NAVIGATION_COMMAND_MENU1      // (1) – Execute navigation menu command 1. For DVD, this enters the DVD root menu, or exits back to the title from the menu.
	NavigationCommandMenu2     NavigationCommand = C.GST_NAVIGATION_COMMAND_MENU2      // (2) – Execute navigation menu command 2. For DVD, this jumps to the DVD title menu.
	NavigationCommandMenu3     NavigationCommand = C.GST_NAVIGATION_COMMAND_MENU3      // (3) – Execute navigation menu command 3. For DVD, this jumps into the DVD root menu.
	NavigationCommandMenu4     NavigationCommand = C.GST_NAVIGATION_COMMAND_MENU4      // (4) – Execute navigation menu command 4. For DVD, this jumps to the Subpicture menu.
	NavigationCommandMenu5     NavigationCommand = C.GST_NAVIGATION_COMMAND_MENU5      // (5) – Execute navigation menu command 5. For DVD, this jumps to the audio menu.
	NavigationCommandMenu6     NavigationCommand = C.GST_NAVIGATION_COMMAND_MENU6      // (6) – Execute navigation menu command 6. For DVD, this jumps to the angles menu.
	NavigationCommandMenu7     NavigationCommand = C.GST_NAVIGATION_COMMAND_MENU7      // (7) – Execute navigation menu command 7. For DVD, this jumps to the chapter menu.
	NavigationCommandLeft      NavigationCommand = C.GST_NAVIGATION_COMMAND_LEFT       // (20) – Select the next button to the left in a menu, if such a button exists.
	NavigationCommandRight     NavigationCommand = C.GST_NAVIGATION_COMMAND_RIGHT      // (21) – Select the next button to the right in a menu, if such a button exists.
	NavigationCommandUp        NavigationCommand = C.GST_NAVIGATION_COMMAND_UP         // (22) – Select the button above the current one in a menu, if such a button exists.
	NavigationCommandDown      NavigationCommand = C.GST_NAVIGATION_COMMAND_DOWN       // (23) – Select the button below the current one in a menu, if such a button exists.
	NavigationCommandActivate  NavigationCommand = C.GST_NAVIGATION_COMMAND_ACTIVATE   // (24) – Activate (click) the currently selected button in a menu, if such a button exists.
	NavigationCommandPrevAngle NavigationCommand = C.GST_NAVIGATION_COMMAND_PREV_ANGLE // (30) – Switch to the previous angle in a multiangle feature.
	NavigationCommandNextAngle NavigationCommand = C.GST_NAVIGATION_COMMAND_NEXT_ANGLE // (31) – Switch to the next angle in a multiangle feature.
)

// Extra aliases for convenience in handling DVD navigation,
const (
	NavigationCommandDVDMenu           NavigationCommand = C.GST_NAVIGATION_COMMAND_DVD_MENU
	NavigationCommandDVDTitleMenu      NavigationCommand = C.GST_NAVIGATION_COMMAND_DVD_TITLE_MENU
	NavigationCommandDVDRootMenu       NavigationCommand = C.GST_NAVIGATION_COMMAND_DVD_ROOT_MENU
	NavigationCommandDVDSubpictureMenu NavigationCommand = C.GST_NAVIGATION_COMMAND_DVD_SUBPICTURE_MENU
	NavigationCommandDVDAudioMenu      NavigationCommand = C.GST_NAVIGATION_COMMAND_DVD_AUDIO_MENU
	NavigationCommandDVDAngleMenu      NavigationCommand = C.GST_NAVIGATION_COMMAND_DVD_ANGLE_MENU
	NavigationCommandDVDChapterMenu    NavigationCommand = C.GST_NAVIGATION_COMMAND_DVD_CHAPTER_MENU
)

// NavigationEventType are enum values for the various events that an element implementing the
// Navigation interface might send up the pipeline.
type NavigationEventType int

// Type castings
const (
	NavigationEventInvalid            NavigationEventType = C.GST_NAVIGATION_EVENT_INVALID              // (0) – Returned from gst_navigation_event_get_type when the passed event is not a navigation event.
	NavigationEventKeyPress           NavigationEventType = C.GST_NAVIGATION_EVENT_KEY_PRESS            // (1) – A key press event. Use gst_navigation_event_parse_key_event to extract the details from the event.
	NavigationEventKeyRelease         NavigationEventType = C.GST_NAVIGATION_EVENT_KEY_RELEASE          // (2) – A key release event. Use gst_navigation_event_parse_key_event to extract the details from the event.
	NavigationEventMouseButtonPress   NavigationEventType = C.GST_NAVIGATION_EVENT_MOUSE_BUTTON_PRESS   // (3) – A mouse button press event. Use gst_navigation_event_parse_mouse_button_event to extract the details from the event.
	NavigationEventMouseButtonRelease NavigationEventType = C.GST_NAVIGATION_EVENT_MOUSE_BUTTON_RELEASE // (4) – A mouse button release event. Use gst_navigation_event_parse_mouse_button_event to extract the details from the event.
	NavigationEventMouseMove          NavigationEventType = C.GST_NAVIGATION_EVENT_MOUSE_MOVE           // (5) – A mouse movement event. Use gst_navigation_event_parse_mouse_move_event to extract the details from the event.
	NavigationEventCommand            NavigationEventType = C.GST_NAVIGATION_EVENT_COMMAND              // (6) – A navigation command event. Use gst_navigation_event_parse_command to extract the details from the event.
	NavigationEventMouseScroll        NavigationEventType = C.GST_NAVIGATION_EVENT_MOUSE_SCROLL         // (7) – A mouse scroll event. Use gst_navigation_event_parse_mouse_scroll_event to extract the details from the event. (Since: 1.18)
)

// NavigationMessageType is a set of notifications that may be received on the bus when navigation
// related status changes.
type NavigationMessageType int

// Type castings
const (
	NavigationMessageInvalid         NavigationMessageType = C.GST_NAVIGATION_MESSAGE_INVALID          // (0) – Returned from gst_navigation_message_get_type when the passed message is not a navigation message.
	NavigationMessageMouseOver       NavigationMessageType = C.GST_NAVIGATION_MESSAGE_MOUSE_OVER       // (1) – Sent when the mouse moves over or leaves a clickable region of the output, such as a DVD menu button.
	NavigationMessageCommandsChanged NavigationMessageType = C.GST_NAVIGATION_MESSAGE_COMMANDS_CHANGED // (2) – Sent when the set of available commands changes and should re-queried by interested applications.
	NavigationMessageAnglesChanged   NavigationMessageType = C.GST_NAVIGATION_MESSAGE_ANGLES_CHANGED   // (3) – Sent when display angles in a multi-angle feature (such as a multiangle DVD) change - either angles have appeared or disappeared.
	NavigationMessageEvent           NavigationMessageType = C.GST_NAVIGATION_MESSAGE_EVENT            // (4) – Sent when a navigation event was not handled by any element in the pipeline
)

// NavigationQueryType represents types of navigation interface queries.
type NavigationQueryType int

// Type castings
const (
	NavigationQueryInvalid  NavigationQueryType = C.GST_NAVIGATION_QUERY_INVALID  // (0) – invalid query
	NavigationQueryCommands NavigationQueryType = C.GST_NAVIGATION_QUERY_COMMANDS // (1) – command query
	NavigationQueryAngles   NavigationQueryType = C.GST_NAVIGATION_QUERY_ANGLES   // (2) – viewing angle query
)

// KeyEvent represents types of key events.
type KeyEvent string

// Enums
const (
	KeyPress   KeyEvent = "key-press"
	KeyRelease KeyEvent = "key-release"
)

// MouseEvent represents types of mouse events.
type MouseEvent string

// Enums
const (
	MouseButtonPress   MouseEvent = "mouse-button-press"
	MouseButtonRelease MouseEvent = "mouse-button-release"
	MouseMove          MouseEvent = "mouse-move"
)

/*
Navigation interface is used for creating and injecting navigation related events such as
mouse button presses, cursor motion and key presses. The associated library also provides
methods for parsing received events, and for sending and receiving navigation related bus
events. One main use-case is DVD menu navigation.


  The main parts of the API are:

	- The Navigation interface, implemented by elements which provide an application with
      the ability to create and inject navigation events into the pipeline.

	- Navigation event handling API. Navigation events are created in response to calls
	  on a Navigation interface implementation, and sent in the pipeline. Upstream elements
      can use the navigation event API functions to parse the contents of received messages.

	- Navigation message handling API. Navigation messages may be sent on the message bus
	  to inform applications of navigation related changes in the pipeline, such as the mouse
      moving over a clickable region, or the set of available angles changing.


The Navigation message functions provide functions for creating and parsing custom bus
messages for signaling GstNavigation changes.
*/
type Navigation interface {
	// Sends the indicated command to the navigation interface.
	SendCommand(NavigationCommand)
	// Sends an event with the given structure.
	SendEvent(*gst.Structure)
	// Sends the given key event. Recognized values for the event are "key-press"
	// and "key-release". The key is the character representation of the key. This is typically
	// as produced by XKeysymToString.
	SendKeyEvent(event KeyEvent, key string)
	// Sends a mouse event to the navigation interface. Mouse event coordinates are sent relative
	// to the display space of the related output area. This is usually the size in pixels of the
	// window associated with the element implementing the Navigation interface. Use 0 for the
	// button when doing mouse move events.
	SendMouseEvent(event MouseEvent, button int, x, y float64)
	// Sends a mouse scroll event to the navigation interface. Mouse event coordinates are sent
	// relative to the display space of the related output area. This is usually the size in pixels
	// of the window associated with the element implementing the Navigation interface.
	SendMouseScrollEvent(x, y, dX, dY float64)
}

// NavigationFromElement checks if the given element implements the Navigation interface. If it does,
// a useable interface is returned. Otherwise, it returns nil.
func NavigationFromElement(element *gst.Element) Navigation {
	if C.toGstNavigation(fromCoreElement(element)) == nil {
		return nil
	}
	return &gstNavigation{fromCoreElement(element)}
}

type gstNavigation struct {
	elem *C.GstElement
}

func (n *gstNavigation) instance() *C.GstNavigation {
	return C.toGstNavigation(n.elem)
}

func (n *gstNavigation) SendCommand(cmd NavigationCommand) {
	C.gst_navigation_send_command(n.instance(), C.GstNavigationCommand(cmd))
}

func (n *gstNavigation) SendEvent(structure *gst.Structure) {
	C.gst_navigation_send_event(n.instance(), fromCoreStructure(structure))
}

func (n *gstNavigation) SendKeyEvent(event KeyEvent, key string) {
	cEvent := C.CString(string(event))
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cEvent))
	defer C.free(unsafe.Pointer(cKey))
	C.gst_navigation_send_key_event(
		n.instance(),
		(*C.gchar)(unsafe.Pointer(cEvent)),
		(*C.gchar)(unsafe.Pointer(cKey)),
	)
}

func (n *gstNavigation) SendMouseEvent(event MouseEvent, button int, x, y float64) {
	cEvent := C.CString(string(event))
	defer C.free(unsafe.Pointer(cEvent))
	C.gst_navigation_send_mouse_event(
		n.instance(),
		(*C.gchar)(unsafe.Pointer(cEvent)),
		C.int(button), C.double(x), C.double(y),
	)
}

func (n *gstNavigation) SendMouseScrollEvent(x, y, dX, dY float64) {
	C.gst_navigation_send_mouse_scroll_event(
		n.instance(),
		C.double(x), C.double(y), C.double(dX), C.double(dY),
	)
}

// NavigationEvent extends the Event from the core library and is used by elements
// implementing the Navigation interface. You can wrap an event in this struct yourself,
// but it is safer to use the ToNavigationEvent method first to check validity.
type NavigationEvent struct{ *gst.Event }

// ToNavigationEvent checks if the given event is a NavigationEvent, and if so, returrns
// a NavigationEvent instance wrapping the event. If the event is not a NavigationEvent
// this function returns nil.
func ToNavigationEvent(event *gst.Event) *NavigationEvent {
	evType := NavigationEventType(C.gst_navigation_event_get_type(fromCoreEvent(event)))
	if evType == NavigationEventInvalid {
		return nil
	}
	return &NavigationEvent{event}
}

// GetType returns the type of this event.
func (e *NavigationEvent) GetType() NavigationEventType {
	return NavigationEventType(C.gst_navigation_event_get_type(e.instance()))
}

// instance returns the underlying GstEvent instance.
func (e *NavigationEvent) instance() *C.GstEvent { return fromCoreEvent(e.Event) }

// NavigationMessage extends the Event from the core library and is used by elements
// implementing the Navigation interface. You can wrap a message in this struct yourself,
// but it is safer to use the ToNavigationMessage method first to check validity.
type NavigationMessage struct{ *gst.Message }

// ToNavigationMessage checks if the given message is a NavigationMessage, and if so,
// returns a NavigatonMessage instance wrapping the message. If the message is not a
// NavigationMessage, this function returns nil.
func ToNavigationMessage(msg *gst.Message) *NavigationMessage {
	msgType := NavigationMessageType(C.gst_navigation_message_get_type(fromCoreMessage(msg)))
	if msgType == NavigationMessageInvalid {
		return nil
	}
	return &NavigationMessage{msg}
}

// instance returns the underlying GstMessage instance.
func (m *NavigationMessage) instance() *C.GstMessage { return fromCoreMessage(m.Message) }

// GetType returns the type of this message.
func (m *NavigationMessage) GetType() NavigationMessageType {
	return NavigationMessageType(C.gst_navigation_message_get_type(m.instance()))
}

// NavigationQuery extends the Query from the core library and is used by elements
// implementing the Navigation interface. You can wrap a query in this struct yourself,
// but it is safer to use the ToNavigationQuery method first to check validity.
type NavigationQuery struct{ *gst.Query }

// ToNavigationQuery checks if the given query is a NavigationQuery, and if so, returns
// a NavigationQuery instance wrapping the query. If the query is not a NavigationQuery,
// this function returns nil.
func ToNavigationQuery(query *gst.Query) *NavigationQuery {
	qType := NavigationQueryType(C.gst_navigation_query_get_type(fromCoreQuery(query)))
	if qType == NavigationQueryInvalid {
		return nil
	}
	return &NavigationQuery{query}
}

// instance returns the underlying GstQuery instance.
func (q *NavigationQuery) instance() *C.GstQuery { return fromCoreQuery(q.Query) }

// GetType returns the type of this query.
func (q *NavigationQuery) GetType() NavigationQueryType {
	return NavigationQueryType(C.gst_navigation_query_get_type(q.instance()))
}
