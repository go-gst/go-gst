package gst

/*
#include "gst.go.h"

extern GstPadProbeReturn goPadProbeFunc             (GstPad * pad, GstPadProbeInfo * info, gpointer user_data);
extern gboolean          goPadForwardFunc           (GstPad * pad, gpointer user_data);
extern void              goGDestroyNotifyFuncNoRun  (gpointer user_data);

GstPadProbeReturn cgoPadProbeFunc (GstPad * pad, GstPadProbeInfo * info, gpointer user_data)
{
	return goPadProbeFunc(pad, info, user_data);
}

gboolean cgoPadForwardFunc (GstPad * pad, gpointer user_data)
{
	return goPadForwardFunc(pad, user_data);
}

void cgoGDestroyNotifyFuncNoRun (gpointer user_data)
{
	goGDestroyNotifyFuncNoRun(user_data);
}

*/
import "C"

import (
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

// Pad is a go representation of a GstPad
type Pad struct{ *Object }

// NewPad returns a new pad with the given direction. If name is empty, one will be generated for you.
func NewPad(name string, direction PadDirection) *Pad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_pad_new(cName, C.GstPadDirection(direction))
	if pad == nil {
		return nil
	}
	return wrapPad(toGObject(unsafe.Pointer(pad)))
}

// NewPadFromTemplate creates a new pad with the given name from the given template. If name is empty, one will
// be generated for you.
func NewPadFromTemplate(tmpl *PadTemplate, name string) *Pad {
	var cName *C.gchar
	if name != "" {
		cStr := C.CString(name)
		defer C.free(unsafe.Pointer(cStr))
		cName = (*C.gchar)(unsafe.Pointer(cStr))
	}
	pad := C.gst_pad_new_from_template(tmpl.Instance(), cName)
	if pad == nil {
		return nil
	}
	return wrapPad(toGObject(unsafe.Pointer(pad)))
}

// Instance returns the underlying C GstPad.
func (p *Pad) Instance() *C.GstPad { return C.toGstPad(p.Unsafe()) }

// Direction returns the direction of this pad.
func (p *Pad) Direction() PadDirection {
	return PadDirection(C.gst_pad_get_direction((*C.GstPad)(p.Instance())))
}

// Template returns the template for this pad or nil.
func (p *Pad) Template() *PadTemplate {
	return wrapPadTemplate(toGObject(unsafe.Pointer(p.Instance().padtemplate)))
}

// CurrentCaps returns the caps for this Pad or nil.
func (p *Pad) CurrentCaps() *Caps {
	caps := C.gst_pad_get_current_caps((*C.GstPad)(p.Instance()))
	if caps == nil {
		return nil
	}
	return wrapCaps(caps)
}

// ActivateMode activates or deactivates the given pad in mode via dispatching to the pad's activatemodefunc.
// For use from within pad activation functions only.
//
// If you don't know what this is, you probably don't want to call it.
func (p *Pad) ActivateMode(mode PadMode, active bool) {
	C.gst_pad_activate_mode(p.Instance(), C.GstPadMode(mode), gboolean(active))
}

// PadProbeCallback is a callback used by Pad AddProbe. It gets called to notify about the current blocking type.
type PadProbeCallback func(*Pad, *PadProbeInfo) PadProbeReturn

// AddProbe adds a callback to be notified of different states of pads. The provided callback is called for every state that matches mask.
//
// Probes are called in groups: First GST_PAD_PROBE_TYPE_BLOCK probes are called, then others, then finally GST_PAD_PROBE_TYPE_IDLE. The only
// exception here are GST_PAD_PROBE_TYPE_IDLE probes that are called immediately if the pad is already idle while calling gst_pad_add_probe.
// In each of the groups, probes are called in the order in which they were added.
//
// A probe ID is returned that can be used to remove the probe.
func (p *Pad) AddProbe(mask PadProbeType, f PadProbeCallback) uint64 {
	ptr := gopointer.Save(f)
	ret := C.gst_pad_add_probe(
		p.Instance(),
		C.GstPadProbeType(mask),
		C.GstPadProbeCallback(C.cgoPadProbeFunc),
		(C.gpointer)(unsafe.Pointer(ptr)),
		C.GDestroyNotify(C.cgoGDestroyNotifyFuncNoRun),
	)
	return uint64(ret)
}

// CanLink checks if this pad is compatible with the given sink pad.
func (p *Pad) CanLink(sink *Pad) bool {
	return gobool(C.gst_pad_can_link(p.Instance(), sink.Instance()))
}

// Chain a buffer to pad.
//
// The function returns FlowFlushing if the pad was flushing.
//
// If the buffer type is not acceptable for pad (as negotiated with a preceding EventCaps event), this function returns FlowNotNegotiated.
//
// The function proceeds calling the chain function installed on pad (see SetChainFunction) and the return value of that function is returned to
// the caller. FlowNotSupported is returned if pad has no chain function.
//
// In all cases, success or failure, the caller loses its reference to buffer after calling this function.
func (p *Pad) Chain(buffer *Buffer) FlowReturn {
	return FlowReturn(C.gst_pad_chain(p.Instance(), buffer.Instance()))
}

// ChainList chains a bufferlist to pad.
//
// The function returns FlowFlushing if the pad was flushing.
//
// If pad was not negotiated properly with a CAPS event, this function returns FlowNotNegotiated.
//
// The function proceeds calling the chainlist function installed on pad (see SetChainListFunction) and the return value of that function is returned
// to the caller. FlowNotSupported is returned if pad has no chainlist function.
//
// In all cases, success or failure, the caller loses its reference to list after calling this function.
func (p *Pad) ChainList(bufferList *BufferList) FlowReturn {
	return FlowReturn(C.gst_pad_chain_list(p.Instance(), bufferList.Instance()))
}

// CheckReconfigure checks and clear the PadFlagNeedReconfigure flag on pad and return TRUE if the flag was set.
func (p *Pad) CheckReconfigure() bool {
	return gobool(C.gst_pad_check_reconfigure(p.Instance()))
}

// CreateStreamID creates a stream-id for the source GstPad pad by combining the upstream information with the optional stream_id of the stream of pad.
// Pad must have a parent GstElement and which must have zero or one sinkpad. stream_id can only be NULL if the parent element of pad has only a single
// source pad.
//
// This function generates an unique stream-id by getting the upstream stream-start event stream ID and appending stream_id to it. If the element has no
// sinkpad it will generate an upstream stream-id by doing an URI query on the element and in the worst case just uses a random number. Source elements
// that don't implement the URI handler interface should ideally generate a unique, deterministic stream-id manually instead.
//
// Since stream IDs are sorted alphabetically, any numbers in the stream ID should be printed with a fixed number of characters, preceded by 0's, such as
// by using the format %03u instead of %u.
func (p *Pad) CreateStreamID(parent *Element, streamID string) string {
	var gstreamID *C.gchar
	if streamID != "" {
		ptr := C.CString(streamID)
		defer C.free(unsafe.Pointer(ptr))
		gstreamID = (*C.gchar)(unsafe.Pointer(ptr))
	}
	ret := C.gst_pad_create_stream_id(p.Instance(), parent.Instance(), gstreamID)
	if ret == nil {
		return ""
	}
	defer C.g_free((C.gpointer)(unsafe.Pointer(ret)))
	return C.GoString(ret)
}

// EventDefault invokes the default event handler for the given pad.
//
// The EOS event will pause the task associated with pad before it is forwarded to all internally linked pads,
//
// The event is sent to all pads internally linked to pad. This function takes ownership of event.
func (p *Pad) EventDefault(parent *Object, event *Event) bool {
	return gobool(C.gst_pad_event_default(p.Instance(), parent.Instance(), event.Instance()))
}

// PadForwardFunc is called for all internally linked pads, see Pad Forward().
// If the function returns true, the procedure is stopped.
type PadForwardFunc func(pad *Pad) bool

// Forward calls the given function for all internally linked pads of pad. This function deals with dynamically changing internal pads and will make sure
// that the forward function is only called once for each pad.
//
// When forward returns TRUE, no further pads will be processed.
func (p *Pad) Forward(f PadForwardFunc) bool {
	ptr := gopointer.Save(f)
	defer gopointer.Unref(ptr)
	return gobool(C.gst_pad_forward(
		p.Instance(),
		C.GstPadForwardFunction(C.cgoPadForwardFunc),
		(C.gpointer)(unsafe.Pointer(ptr)),
	))
}

// GetAllowedCaps getss the capabilities of the allowed media types that can flow through pad and its peer.
//
// The allowed capabilities is calculated as the intersection of the results of calling QueryCaps on pad and its peer. The caller owns a reference on the
// resulting caps.
func (p *Pad) GetAllowedCaps() *Caps {
	return wrapCaps(C.gst_pad_get_allowed_caps(
		p.Instance(),
	))
}

// GetCurrentCaps gets the capabilities currently configured on pad with the last EventCaps event.
func (p *Pad) GetCurrentCaps() *Caps {
	return wrapCaps(C.gst_pad_get_current_caps(
		p.Instance(),
	))
}

// GetDirection gets the direction of the pad. The direction of the pad is decided at construction time so this function does not take the LOCK.
func (p *Pad) GetDirection() PadDirection {
	return PadDirection(C.gst_pad_get_direction(p.Instance()))
}

// GetElementPrivate gets the private data of a pad. No locking is performed in this function.
func (p *Pad) GetElementPrivate() interface{} {
	ptr := C.gst_pad_get_element_private(p.Instance())
	return gopointer.Restore(unsafe.Pointer(ptr))
}

// GetLastFlowReturn gets the FlowReturn return from the last data passed by this pad.
func (p *Pad) GetLastFlowReturn() FlowReturn {
	return FlowReturn(C.gst_pad_get_last_flow_return(
		p.Instance(),
	))
}

// GetOffset gets the offset applied to the running time of pad. pad has to be a source pad.
func (p *Pad) GetOffset() int64 {
	return int64(C.gst_pad_get_offset(p.Instance()))
}

// GetPadTemplate gets the template for this pad.
func (p *Pad) GetPadTemplate() *PadTemplate {
	tmpl := C.gst_pad_get_pad_template(p.Instance())
	if tmpl == nil {
		return nil
	}
	return wrapPadTemplate(toGObject(unsafe.Pointer(tmpl)))
}

// GetPadTemplateCaps gets the capabilities for pad's template.
func (p *Pad) GetPadTemplateCaps() *Caps {
	caps := C.gst_pad_get_pad_template_caps(p.Instance())
	if caps == nil {
		return nil
	}
	return wrapCaps(caps)
}

// GetParentElement gets the parent of pad, cast to a Element. If a pad has no parent or its
// parent is not an element, return nil.
func (p *Pad) GetParentElement() *Element {
	elem := C.gst_pad_get_parent_element(p.Instance())
	if elem == nil {
		return nil
	}
	return wrapElement(toGObject(unsafe.Pointer(elem)))
}

// GetPeer gets the peer of pad. This function refs the peer pad so you need to unref it after use.
func (p *Pad) GetPeer() *Pad {
	peer := C.gst_pad_get_peer(p.Instance())
	if peer == nil {
		return nil
	}
	return wrapPad(toGObject(unsafe.Pointer(peer)))
}

// GetRange calls the getrange function of pad, see PadGetRangeFunc for a description of a getrange function.
// If pad has no getrange function installed (see SetGetRangeFunction) this function returns FlowNotSupported.
//
// If buffer points to a variable holding nil, a valid new GstBuffer will be placed in buffer when this function
// returns FlowOK. The new buffer must be freed with Unref after usage.
//
// When buffer points to a variable that points to a valid Buffer, the buffer will be filled with the result data
// when this function returns FlowOK. If the provided buffer is larger than size, only size bytes will be filled
// in the result buffer and its size will be updated accordingly.
//
// Note that less than size bytes can be returned in buffer when, for example, an EOS condition is near or when
// buffer is not large enough to hold size bytes. The caller should check the result buffer size to get the result
// size.
//
// When this function returns any other result value than FlowOK, buffer will be unchanged.
//
// This is a lowlevel function. Usually PullRange is used.
func (p *Pad) GetRange(offset uint64, size uint, buffer *Buffer) (FlowReturn, *Buffer) {
	var buf *C.GstBuffer
	if buffer != nil {
		buf = buffer.Instance()
	}
	ret := C.gst_pad_get_range(p.Instance(), C.guint64(offset), C.guint(size), &buf)
	var newBuf *Buffer
	if buf != nil {
		newBuf = wrapBuffer(buf)
	} else {
		newBuf = nil
	}
	return FlowReturn(ret), newBuf
}

// GetSingleInternalLink checks if there is a single internal link of the given pad, and returns it. Otherwise, it will
// return nil.
func (p *Pad) GetSingleInternalLink() *Pad {
	pad := C.gst_pad_get_single_internal_link(p.Instance())
	if pad == nil {
		return nil
	}
	return wrapPad(toGObject(unsafe.Pointer(pad)))
}

// GetStickyEvent returns a new reference of the sticky event of type event_type from the event.
func (p *Pad) GetStickyEvent(eventType EventType, idx uint) *Event {
	ev := C.gst_pad_get_sticky_event(p.Instance(), C.GstEventType(eventType), C.guint(idx))
	if ev == nil {
		return nil
	}
	return wrapEvent(ev)
}

// GetStream returns the current Stream for the pad, or nil if none has been set yet, i.e. the pad has not received a
// stream-start event yet.
//
// This is a convenience wrapper around GetStickyEvent and Event ParseStream.
func (p *Pad) GetStream() *Stream {
	st := C.gst_pad_get_stream(p.Instance())
	if st == nil {
		return nil
	}
	return wrapStream(toGObject(unsafe.Pointer(st)))
}

// GetStreamID returns the current stream-id for the pad, or an empty string if none has been set yet, i.e. the pad has not received
// a stream-start event yet.
//
// This is a convenience wrapper around gst_pad_get_sticky_event and gst_event_parse_stream_start.
//
// The returned stream-id string should be treated as an opaque string, its contents should not be interpreted.
func (p *Pad) GetStreamID() string {
	id := C.gst_pad_get_stream_id(p.Instance())
	if id == nil {
		return ""
	}
	defer C.g_free((C.gpointer)(unsafe.Pointer(id)))
	return C.GoString(id)
}

// GetTaskState gets the pad task state. If no task is currently set, TaskStopped is returned.
func (p *Pad) GetTaskState() TaskState {
	return TaskState(C.gst_pad_get_task_state(p.Instance()))
}

// PadProbeInfo represents the info passed to a PadProbeCallback.
type PadProbeInfo struct {
	ptr *C.GstPadProbeInfo
}

// ID returns the id of the probe.
func (p *PadProbeInfo) ID() uint32 { return uint32(p.ptr.id) }

// Type returns the type of the probe. The type indicates the type of data that can be expected
// with the probe.
func (p *PadProbeInfo) Type() PadProbeType { return PadProbeType(p.ptr._type) }

// Offset returns the offset of pull probe, this field is valid when type contains PadProbeTypePull.
func (p *PadProbeInfo) Offset() uint64 { return uint64(p.ptr.offset) }

// Size returns the size of pull probe, this field is valid when type contains PadProbeTypePull.
func (p *PadProbeInfo) Size() uint64 { return uint64(p.ptr.size) }

// GetBuffer returns the buffer, if any, inside this probe info.
func (p *PadProbeInfo) GetBuffer() *Buffer {
	buf := C.gst_pad_probe_info_get_buffer(p.ptr)
	if buf == nil {
		return nil
	}
	return wrapBuffer(buf)
}

// GetBufferList returns the buffer list, if any, inside this probe info.
func (p *PadProbeInfo) GetBufferList() *BufferList {
	bufList := C.gst_pad_probe_info_get_buffer_list(p.ptr)
	if bufList == nil {
		return nil
	}
	return wrapBufferList(bufList)
}

// GetEvent returns the event, if any, inside this probe info.
func (p *PadProbeInfo) GetEvent() *Event {
	ev := C.gst_pad_probe_info_get_event(p.ptr)
	if ev == nil {
		return nil
	}
	return wrapEvent(ev)
}

// GetQuery returns the query, if any, inside this probe info.
func (p *PadProbeInfo) GetQuery() *Query {
	q := C.gst_pad_probe_info_get_query(p.ptr)
	if q == nil {
		return nil
	}
	return wrapQuery(q)
}
