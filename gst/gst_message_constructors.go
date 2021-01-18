package gst

// #include "gst.go.h"
import "C"
import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

func getMessageSourceObj(src interface{}) *C.GstObject {
	usafe, ok := src.(interface{ Unsafe() unsafe.Pointer })
	if !ok {
		return nil
	}
	return C.toGstObject(usafe.Unsafe())
}

// NewApplicationMessage creates a new application-typed message. GStreamer will never
// create these messages; they are a gift from them to you. Enjoy.
//
// The source of all message constructors must be a valid Object or descendant, specifically
// one created from the go runtime. If not the message returned will be nil.
func NewApplicationMessage(src interface{}, structure *Structure) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_application(srcObj, structure.Instance())))
}

// NewAsyncDoneMessage builds a message that is posted when elements completed an ASYNC state change.
// RunningTime contains the time of the desired running time when this elements goes to PLAYING.
// A value less than 0 for runningTime means that the element has no clock interaction and thus doesn't
// care about the running time of the pipeline.
func NewAsyncDoneMessage(src interface{}, runningTime time.Duration) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	var cTime C.GstClockTime
	if runningTime.Nanoseconds() < 0 {
		cTime = C.GstClockTime(gstClockTimeNone)
	} else {
		cTime = C.GstClockTime(runningTime.Nanoseconds())
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_async_done(
		srcObj,
		cTime,
	)))
}

// NewAsyncStartMessage returns a message that is posted by elements when they start an ASYNC state change.
func NewAsyncStartMessage(src interface{}) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_async_start(srcObj)))
}

// NewBufferingMessage returns a message that can be posted by an element that needs to buffer data before it
// can continue processing. percent should be a value between 0 and 100. A value of 100 means that the buffering completed.
//
// When percent is < 100 the application should PAUSE a PLAYING pipeline. When percent is 100, the application can set the
// pipeline (back) to PLAYING. The application must be prepared to receive BUFFERING messages in the PREROLLING state and
// may only set the pipeline to PLAYING after receiving a message with percent set to 100, which can happen after the pipeline
// completed prerolling.
func NewBufferingMessage(src interface{}, percent int) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_buffering(srcObj, C.gint(percent))))
}

// NewClockLostMessage creates a clock lost message. This message is posted whenever the clock is not valid anymore.
//
// If this message is posted by the pipeline, the pipeline will select a new clock again when it goes to PLAYING. It might
// therefore be needed to set the pipeline to PAUSED and PLAYING again.
func NewClockLostMessage(src interface{}, clock *Clock) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_clock_lost(srcObj, clock.Instance())))
}

// NewClockProvideMessage creates a clock provide message. This message is posted whenever an element is ready to provide a
// clock or lost its ability to provide a clock (maybe because it paused or became EOS).
//
// This message is mainly used internally to manage the clock selection.
func NewClockProvideMessage(src interface{}, clock *Clock, ready bool) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_clock_provide(srcObj, clock.Instance(), gboolean(ready))))
}

// NewCustomMessage creates a new custom-typed message. This can be used for anything not handled by other message-specific
// functions to pass a message to the app. The structure field can be nil.
func NewCustomMessage(src interface{}, msgType MessageType, structure *Structure) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	if structure == nil {
		return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_custom(C.GstMessageType(msgType), srcObj, nil)))
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_custom(C.GstMessageType(msgType), srcObj, structure.Instance())))
}

// NewDeviceAddedMessage creates a new device-added message. The device-added message is produced by a DeviceProvider or a DeviceMonitor.
// They announce the appearance of monitored devices.
func NewDeviceAddedMessage(src interface{}, device *Device) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_device_added(srcObj, device.Instance())))
}

// NewDeviceChangedMessage creates a new device-changed message. The device-changed message is produced by a DeviceProvider or a DeviceMonitor.
// They announce that a device properties has changed and device represent the new modified version of changed_device.
func NewDeviceChangedMessage(src interface{}, device, changedDevice *Device) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_device_changed(srcObj, device.Instance(), changedDevice.Instance())))
}

// NewDeviceRemovedMessage creates a new device-removed message. The device-removed message is produced by a DeviceProvider or a DeviceMonitor.
// They announce the disappearance of monitored devices.
func NewDeviceRemovedMessage(src interface{}, device *Device) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_device_removed(srcObj, device.Instance())))
}

// NewDurationChangedMessage creates a new duration changed message. This message is posted by elements that know the duration of a
// stream when the duration changes. This message is received by bins and is used to calculate the total duration of a pipeline.
func NewDurationChangedMessage(src interface{}) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_duration_changed(srcObj)))
}

// NewElementMessage creates a new element-specific message. This is meant as a generic way of allowing one-way communication from an
// element to an application, for example "the firewire cable was unplugged". The format of the message should be documented in the
// element's documentation. The structure field can be nil.
func NewElementMessage(src interface{}, structure *Structure) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	if structure == nil {
		return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_element(srcObj, nil)))
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_element(srcObj, structure.Instance())))
}

// NewEOSMessage creates a new eos message. This message is generated and posted in the sink elements of a Bin. The bin will only forward
// the EOS message to the application if all sinks have posted an EOS message.
func NewEOSMessage(src interface{}) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_eos(srcObj)))
}

// NewErrorMessage creates a new error message. The message will copy error and debug. This message is posted by element when a fatal event
// occurred. The pipeline will probably (partially) stop. The application receiving this message should stop the pipeline.
// Structure can be nil to not add a structure to the message.
func NewErrorMessage(src interface{}, err error, debugStr string, structure *Structure) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}

	errmsg := C.CString(err.Error())
	gerr := C.g_error_new_literal(DomainLibrary.toQuark(), C.gint(LibraryErrorFailed), (*C.gchar)(errmsg))
	defer C.free(unsafe.Pointer(errmsg))
	defer C.g_error_free(gerr)

	gdebugStr := C.CString(debugStr)
	defer C.free(unsafe.Pointer(gdebugStr))

	if structure != nil {
		return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_error_with_details(
			srcObj,
			gerr,
			gdebugStr,
			structure.Instance(),
		)))
	}

	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_error(
		srcObj,
		gerr,
		gdebugStr,
	)))
}

// NewHaveContextMessage creates a message that is posted when an element has a new local Context.
func NewHaveContextMessage(src interface{}, ctx *Context) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_have_context(
		srcObj,
		ctx.Ref().Instance(),
	)))
}

// NewInfoMessage creates a new info message. Structure can be nil.
func NewInfoMessage(src interface{}, msg string, debugStr string, structure *Structure) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}

	errmsg := C.CString(msg)
	gerr := C.g_error_new_literal(DomainLibrary.toQuark(), C.gint(0), (*C.gchar)(errmsg))
	defer C.free(unsafe.Pointer(errmsg))
	defer C.g_error_free(gerr)

	gdebugStr := C.CString(debugStr)
	defer C.free(unsafe.Pointer(gdebugStr))

	if structure != nil {
		return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_info_with_details(
			srcObj,
			gerr,
			gdebugStr,
			structure.Instance(),
		)))
	}

	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_info(
		srcObj,
		gerr,
		gdebugStr,
	)))
}

// NewLatencyMessage creates a message that can be posted by elements when their latency requirements have changed.
func NewLatencyMessage(src interface{}) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_latency(srcObj)))
}

// NewNeedContextMessage creates a message that is posted when an element needs a specific Context.
func NewNeedContextMessage(src interface{}, ctxType string) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	cStr := C.CString(ctxType)
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_need_context(srcObj, (*C.gchar)(unsafe.Pointer(cStr)))))
}

// NewNewClockMessage creates a new clock message. This message is posted whenever the pipeline selects a new clock for the pipeline.
func NewNewClockMessage(src interface{}, clock *Clock) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_new_clock(srcObj, clock.Instance())))
}

// NewProgressMessage creates messages that are posted by elements when they use an asynchronous task to perform actions triggered by a state change.
//
// Code contains a well defined string describing the action. Text should contain a user visible string detailing the current action.
func NewProgressMessage(src interface{}, progressType ProgressType, code, text string) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	cCode := (*C.gchar)(unsafe.Pointer(C.CString(code)))
	cText := (*C.gchar)(unsafe.Pointer(C.CString(text)))
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_progress(srcObj, C.GstProgressType(progressType), cCode, cText)))
}

// NewPropertyNotifyMessage creates a new message notifying an object's properties have changed. If the
// source OR the value cannot be coereced to C types, the function will return nil.
func NewPropertyNotifyMessage(src interface{}, propName string, val interface{}) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	gVal, err := glib.GValue(val)
	if err != nil {
		return nil
	}
	cName := (*C.gchar)(unsafe.Pointer(C.CString(propName)))
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_property_notify(
		srcObj,
		cName,
		(*C.GValue)(unsafe.Pointer(gVal.GValue)),
	)))
}

// NewQoSMessage creates a message that is posted on the bus whenever an element decides to drop a buffer because of
// QoS reasons or whenever it changes its processing strategy because of QoS reasons (quality adjustments such as processing at lower accuracy).
//
// This message can be posted by an element that performs synchronisation against the clock (live) or it could be dropped by an element that performs
// QoS because of QOS events received from a downstream element (!live).
//
// running_time, stream_time, timestamp, duration should be set to the respective running-time, stream-time, timestamp and duration of the (dropped) buffer
// that generated the QoS event. Values can be left to less than zero when unknown.
func NewQoSMessage(src interface{}, live bool, runningTime, streamTime, timestamp, duration time.Duration) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_qos(
		srcObj,
		gboolean(live),
		C.guint64((runningTime.Nanoseconds())),
		C.guint64((streamTime.Nanoseconds())),
		C.guint64((timestamp.Nanoseconds())),
		C.guint64((duration.Nanoseconds())),
	)))
}

// NewRedirectMessage creates a new redirect message and adds a new entry to it. Redirect messages are posted when an element detects that the actual
// data has to be retrieved from a different location. This is useful if such a redirection cannot be handled inside a source element, for example when
// HTTP 302/303 redirects return a non-HTTP URL.
//
// The redirect message can hold multiple entries. The first one is added when the redirect message is created, with the given location, tag_list, entry_struct
// arguments. Use AddRedirectEntry to add more entries.
//
// Each entry has a location, a tag list, and a structure. All of these are optional. The tag list and structure are useful for additional metadata, such as
// bitrate statistics for the given location.
//
// By default, message recipients should treat entries in the order they are stored. The recipient should therefore try entry #0 first, and if this entry is not
// acceptable or working, try entry #1 etc. Senders must make sure that they add entries in this order. However, recipients are free to ignore the order and pick
// an entry that is "best" for them. One example would be a recipient that scans the entries for the one with the highest bitrate tag.
//
// The specified location string is copied. However, ownership over the tag list and structure are transferred to the message.
func NewRedirectMessage(src interface{}, location string, tagList *TagList, entryStructure *Structure) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	var loc *C.gchar
	var tl *C.GstTagList
	var st *C.GstStructure
	if location != "" {
		locc := C.CString(location)
		defer C.free(unsafe.Pointer(locc))
		loc = (*C.gchar)(unsafe.Pointer(locc))
	}
	if tagList != nil {
		tl = tagList.Ref().Instance()
	}
	if entryStructure != nil {
		st = entryStructure.Instance()
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_redirect(
		srcObj,
		loc, tl, st,
	)))
}

// AddRedirectEntry creates and appends a new entry to the message.
func (m *Message) AddRedirectEntry(location string, tagList *TagList, entryStructure *Structure) {
	var loc *C.gchar
	var tl *C.GstTagList
	var st *C.GstStructure
	if location != "" {
		locc := C.CString(location)
		defer C.free(unsafe.Pointer(locc))
		loc = (*C.gchar)(unsafe.Pointer(locc))
	}
	if tagList != nil {
		tl = tagList.Ref().Instance()
	}
	if entryStructure != nil {
		st = entryStructure.Instance()
	}
	C.gst_message_add_redirect_entry(m.Instance(), loc, tl, st)
}

// NewRequestStateMessage creates a message that can be posted by elements when they want to have their state changed.
// A typical use case would be an audio server that wants to pause the pipeline because a higher priority stream is being played.
func NewRequestStateMessage(src interface{}, state State) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_request_state(srcObj, C.GstState(state))))
}

// NewResetTimeMessage creates a message that is posted when the pipeline running-time should be reset to running_time, like after a flushing seek.
func NewResetTimeMessage(src interface{}, runningTime time.Duration) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_reset_time(srcObj, C.GstClockTime(runningTime.Nanoseconds()))))
}

// NewSegmentDoneMessage creates a new segment done message. This message is posted by elements that finish playback of a segment as a result of a
// segment seek. This message is received by the application after all elements that posted a segment_start have posted the segment_done.
func NewSegmentDoneMessage(src interface{}, format Format, position int64) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_segment_done(
		srcObj,
		C.GstFormat(format),
		C.gint64(position),
	)))
}

// NewSegmentStartMessage creates a new segment message. This message is posted by elements that start playback of a segment as a result of a segment seek.
// This message is not received by the application but is used for maintenance reasons in container elements.
func NewSegmentStartMessage(src interface{}, format Format, position int64) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_segment_start(
		srcObj,
		C.GstFormat(format),
		C.gint64(position),
	)))
}

// NewStateChangedMessage creates a state change message. This message is posted whenever an element changed its state.
func NewStateChangedMessage(src interface{}, oldState, newState, pendingState State) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_state_changed(
		srcObj,
		C.GstState(oldState), C.GstState(newState), C.GstState(pendingState),
	)))
}

// NewStateDirtyMessage creates a state dirty message. This message is posted whenever an element changed its state asynchronously
// and is used internally to update the states of container objects.
func NewStateDirtyMessage(src interface{}) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_state_dirty(srcObj)))
}

// NewStepDoneMessage creates a message that is posted by elements when they complete a part, when intermediate set to TRUE, or a
// complete step operation.
//
// Duration will contain the amount of time of the stepped amount of media in format format.
func NewStepDoneMessage(src interface{}, format Format, amount uint64, rate float64, flush, intermediate bool, duration time.Duration, eos bool) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_step_done(
		srcObj,
		C.GstFormat(format),
		C.guint64(amount),
		C.gdouble(rate),
		gboolean(flush),
		gboolean(intermediate),
		C.guint64(duration.Nanoseconds()),
		gboolean(eos),
	)))
}

// NewStepStartMessage creates a message that is posted by elements when they accept or activate a new step event for amount in format.
//
// Active is set to FALSE when the element accepted the new step event and has queued it for execution in the streaming threads.
//
// Active is set to TRUE when the element has activated the step operation and is now ready to start executing the step in the streaming thread.
// After this message is emitted, the application can queue a new step operation in the element.
func NewStepStartMessage(src interface{}, active bool, format Format, amount uint64, rate float64, flush, intermediate bool) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_step_start(
		srcObj,
		gboolean(active),
		C.GstFormat(format),
		C.guint64(amount),
		C.gdouble(rate),
		gboolean(flush),
		gboolean(intermediate),
	)))
}

// NewStreamCollectionMessage creates a new stream-collection message. The message is used to announce new StreamCollections.
func NewStreamCollectionMessage(src interface{}, collection *StreamCollection) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_stream_collection(srcObj, collection.Instance())))
}

// NewStreamStartMessage creates a new stream_start message. This message is generated and posted in the sink elements of a Bin.
// The bin will only forward the StreamStart message to the application if all sinks have posted a StreamStart message.
func NewStreamStartMessage(src interface{}) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_stream_start(srcObj)))
}

// NewStreamStatusMessage creates a new stream status message. This message is posted when a streaming thread is created/destroyed or
// when the state changed.
func NewStreamStatusMessage(src interface{}, stType StreamStatusType, owner *Element) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_stream_status(srcObj, C.GstStreamStatusType(stType), owner.Instance())))
}

// NewStreamSelectedMessage creates a new steams-selected message. The message is used to announce that an array of streams has been selected.
// This is generally in response to a GST_EVENT_SELECT_STREAMS event, or when an element (such as decodebin3) makes an initial selection of streams.
//
// The message also contains the StreamCollection to which the various streams belong to.
//
// Users of this constructor can add the selected streams with StreamsSelectedAdd.
func NewStreamSelectedMessage(src interface{}, collection *StreamCollection) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_streams_selected(srcObj, collection.Instance())))
}

// StreamsSelectedAdd adds the stream to the message
func (m *Message) StreamsSelectedAdd(stream *Stream) {
	C.gst_message_streams_selected_add(m.Instance(), stream.Instance())
}

// StreamsSelectedSize returns the number of streams contained in the message.
func (m *Message) StreamsSelectedSize() uint {
	return uint(C.gst_message_streams_selected_get_size(m.Instance()))
}

// StreamsSelectedGetStream retrieves the Stream with index index from the message.
func (m *Message) StreamsSelectedGetStream(index uint) *Stream {
	stream := C.gst_message_streams_selected_get_stream(m.Instance(), C.guint(index))
	if stream == nil {
		return nil
	}
	return wrapStream(glib.TransferFull(unsafe.Pointer(stream)))
}

// NewStructureChangeMessage creates a new structure change message. This message is posted when the structure of a pipeline is in the process
// of being changed, for example when pads are linked or unlinked.
//
// Src should be the sinkpad that unlinked or linked.
func NewStructureChangeMessage(src interface{}, chgType StructureChangeType, owner *Element, busy bool) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_structure_change(
		srcObj,
		C.GstStructureChangeType(chgType),
		owner.Instance(),
		gboolean(busy),
	)))
}

// NewTagMessage creates a new tag message. The message will take ownership of the tag list. The message is posted by elements that discovered a new taglist.
func NewTagMessage(src interface{}, tagList *TagList) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_tag(srcObj, tagList.Ref().Instance())))
}

// NewTOCMessage creates a new TOC message. The message is posted by elements that discovered or updated a TOC.
func NewTOCMessage(src interface{}, toc *TOC, updated bool) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_toc(
		srcObj,
		toc.Instance(),
		gboolean(updated),
	)))
}

// NewWarningMessage creates a new warning message. Structure can be nil.
func NewWarningMessage(src interface{}, msg string, debugStr string, structure *Structure) *Message {
	srcObj := getMessageSourceObj(src)
	if srcObj == nil {
		return nil
	}

	errmsg := C.CString(msg)
	gerr := C.g_error_new_literal(DomainLibrary.toQuark(), C.gint(0), (*C.gchar)(errmsg))
	defer C.free(unsafe.Pointer(errmsg))
	defer C.g_error_free(gerr)

	gdebugStr := C.CString(debugStr)
	defer C.free(unsafe.Pointer(gdebugStr))

	if structure != nil {
		return wrapMessage(C.gst_message_new_warning_with_details(
			srcObj,
			gerr,
			gdebugStr,
			structure.Instance(),
		))
	}

	return FromGstMessageUnsafeFull(unsafe.Pointer(C.gst_message_new_warning(
		srcObj,
		gerr,
		gdebugStr,
	)))
}
