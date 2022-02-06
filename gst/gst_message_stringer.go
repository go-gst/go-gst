package gst

import (
	"encoding/json"
	"fmt"
	"strings"
)

// String implements a stringer on the message. It iterates over the type of the message
// and applies the correct parser, then dumps a string of the basic contents of the
// message. This function can be expensive and should only be used for debugging purposes
// or in routines where latency is not a concern.
//
// This stringer really just helps in keeping track of making sure all message types are
// accounted for in some way. It's the devil, writing it was the devil, and I hope you
// enjoy being able to `fmt.Println(msg)`.
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
			msg += tags.String()
		}

	case MessageBuffering:
		stats := m.ParseBufferingStats()
		msg += fmt.Sprintf(
			"Buffering %s - %d%% complete (avg in %d/sec, avg out %d/sec, time left %s)",
			stats.BufferingMode.String(),
			m.ParseBuffering(),
			stats.AverageIn,
			stats.AverageOut,
			stats.BufferingLeft.String(),
		)

	case MessageStateChanged:
		oldState, newState := m.ParseStateChanged()
		msg += fmt.Sprintf("State changed from %s to %s", oldState.String(), newState.String())

	case MessageStateDirty:
		msg += "(DEPRECATED MESSAGE) An element changed state in a streaming thread"

	case MessageStepDone:
		out, err := json.Marshal(m.ParseStepDone())
		if err == nil {
			msg += string(out)
		}

	case MessageClockProvide:
		msg += "Element has clock provide capability"

	case MessageClockLost:
		msg += "Lost a clock"

	case MessageNewClock:
		msg += "Got a new clock"

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
			msg += fmt.Sprintf(" in %s", dur.String())
		}

	case MessageRequestState:
		msg += fmt.Sprintf("State chnage request to %s", m.ParseRequestState().String())

	case MessageStepStart:
		out, err := json.Marshal(m.ParseStepStart())
		if err == nil {
			msg += string(out)
		}

	case MessageQoS:
		out, err := json.Marshal(m.ParseQoS())
		if err == nil {
			msg += string(out)
		}

	case MessageProgress:
		progressType, code, text := m.ParseProgress()
		msg += fmt.Sprintf("%s - %s - %s", strings.ToUpper(progressType.String()), code, text)

	case MessageTOC:
		// TODO

	case MessageResetTime:
		msg += fmt.Sprintf("Running time: %s", m.ParseResetTime().String())

	case MessageStreamStart:
		msg += "Pipeline stream is starting"

	case MessageNeedContext:
		msg += "Element needs context"

	case MessageHaveContext:
		ctx := m.ParseHaveContext()
		msg += fmt.Sprintf("Received context of type %s", ctx.GetType())

	case MessageExtended:
		msg += "Extended message type"

	case MessageDeviceAdded:
		if device := m.ParseDeviceAdded(); device != nil {
			msg += fmt.Sprintf("Device %s added", device.GetDisplayName())
		}

	case MessageDeviceRemoved:
		if device := m.ParseDeviceRemoved(); device != nil {
			msg += fmt.Sprintf("Device %s removed", device.GetDisplayName())
		}

	case MessageDeviceChanged:
		if device, _ := m.ParseDeviceChanged(); device != nil {
			msg += fmt.Sprintf("Device %s had its properties updated", device.GetDisplayName())
		}

	case MessagePropertyNotify:
		obj, propName, propVal := m.ParsePropertyNotify()
		if obj != nil && propVal != nil {
			goval, err := propVal.GoValue()
			if err != nil {
				msg += fmt.Sprintf("Object %s had property '%s' changed to %+v", obj.GetName(), propName, goval)
			}
		}

	case MessageStreamCollection:
		collection := m.ParseStreamCollection()
		msg += fmt.Sprintf("New stream collection with upstream id: %s", collection.GetUpstreamID())

	case MessageStreamsSelected:
		collection := m.ParseStreamsSelected()
		msg += fmt.Sprintf("Stream with upstream id '%s' has selected new streams", collection.GetUpstreamID())

	case MessageRedirect:
		msg += fmt.Sprintf("Received redirect message with %d entries", m.NumRedirectEntries())

	case MessageUnknown:
		msg += "Unknown message type"

	case MessageAny:
		msg += "Message did not match any known types"
	}
	return msg
}
