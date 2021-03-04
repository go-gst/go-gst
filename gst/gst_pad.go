package gst

/*
#include "gst.go.h"

extern gboolean          goPadStickyEventForEachFunc  (GstPad * pad, GstEvent ** event, gpointer user_data);
extern GstPadProbeReturn goPadProbeFunc               (GstPad * pad, GstPadProbeInfo * info, gpointer user_data);
extern gboolean          goPadForwardFunc             (GstPad * pad, gpointer user_data);
extern void              goGDestroyNotifyFuncNoRun    (gpointer user_data);

extern void              goGstPadFuncDestroyNotify     (PadDestroyNotifyInfo *info);
extern gboolean          goGstPadActivateFunction      (GstPad * pad, GstObject * parent);
extern gboolean          goGstPadActivateModeFunction  (GstPad * pad, GstObject * parent, GstPadMode mode, gboolean active);
extern GstFlowReturn     goGstPadChainFunction         (GstPad * pad, GstObject * parent, GstBuffer * buffer);
extern GstFlowReturn     goGstPadChainListFunction     (GstPad * pad, GstObject * parent, GstBufferList * list);
extern GstFlowReturn     goGstPadEventFullFunction     (GstPad * pad, GstObject * parent, GstEvent * event);
extern gboolean          goGstPadEventFunction         (GstPad * pad, GstObject * parent, GstEvent * event);
extern GstFlowReturn     goGstPadGetRangeFunction      (GstPad * pad, GstObject * parent, guint64 offset, guint length, GstBuffer ** buffer);
extern GstIterator *     goGstPadIterIntLinkFunction   (GstPad * pad, GstObject * parent);
extern GstPadLinkReturn  goGstPadLinkFunction          (GstPad * pad, GstObject * parent, GstPad * peer);
extern gboolean          goGstPadQueryFunction         (GstPad * pad, GstObject * parent, GstQuery * query);
extern void              goGstPadUnlinkFunction        (GstPad * pad, GstObject * parent);

gboolean
padsAreEqual(void * pad1, void * pad2) { return pad1 == pad2; }

void
cgoPadFuncDestroyNotify (gpointer notify_info)
{
	goGstPadFuncDestroyNotify((PadDestroyNotifyInfo *)notify_info);
	free(notify_info);
}

gboolean
cgoGstPadActivateFunction(GstPad * pad, GstObject * parent) { return goGstPadActivateFunction(pad, parent); }

gboolean
cgoGstPadActivateModeFunction (GstPad * pad, GstObject * parent, GstPadMode mode, gboolean active) { return goGstPadActivateModeFunction(pad, parent, mode, active); }

GstFlowReturn
cgoGstPadChainFunction (GstPad * pad, GstObject * parent, GstBuffer * buffer) { return goGstPadChainFunction(pad, parent, buffer); }

GstFlowReturn
cgoGstPadChainListFunction (GstPad * pad, GstObject * parent, GstBufferList * list) { return goGstPadChainListFunction(pad, parent, list); }

gboolean
cgoGstPadEventFunction (GstPad * pad, GstObject * parent, GstEvent * event) { return goGstPadEventFunction(pad, parent, event); }

GstFlowReturn
cgoGstPadEventFullFunction (GstPad * pad, GstObject * parent, GstEvent * event) { return goGstPadEventFullFunction(pad, parent, event); }

GstFlowReturn
cgoGstPadGetRangeFunction (GstPad * pad, GstObject * parent, guint64 offset, guint length, GstBuffer ** buffer) { return goGstPadGetRangeFunction(pad, parent, offset, length, buffer); }

GstIterator *
cgoGstPadIterIntLinkFunction (GstPad * pad, GstObject * parent) { return goGstPadIterIntLinkFunction(pad, parent); }

GstPadLinkReturn
cgoGstPadLinkFunction (GstPad * pad, GstObject * parent, GstPad * peer) { return goGstPadLinkFunction(pad, parent, peer); }

gboolean
cgoGstPadQueryFunction (GstPad * pad, GstObject * parent, GstQuery * query) { return goGstPadQueryFunction(pad, parent, query); }

void
cgoGstPadUnlinkFunction (GstPad * pad, GstObject * parent) { return goGstPadUnlinkFunction(pad, parent); }

GstPadProbeReturn
cgoPadProbeFunc (GstPad * pad, GstPadProbeInfo * info, gpointer user_data) { return goPadProbeFunc(pad, info, user_data); }

gboolean
cgoPadForwardFunc (GstPad * pad, gpointer user_data) { return goPadForwardFunc(pad, user_data); }

gboolean
cgoPadStickyEventForEachFunc (GstPad * pad, GstEvent ** event, gpointer user_data) { return goPadStickyEventForEachFunc(pad, event, user_data); }

void
cgoGDestroyNotifyFuncNoRun (gpointer user_data) { goGDestroyNotifyFuncNoRun(user_data); }

*/
import "C"

import (
	"errors"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// PadFuncMap is an type of map for registering callbacks to a pad instance.
// It compensates for the lack of userdata in pad callbacks by providing a FuncForPad
// method that will return the function matching the Pad at the given address.
type PadFuncMap map[unsafe.Pointer]interface{}

// FuncForPad will return the function in this map for the given pad.
func (p PadFuncMap) FuncForPad(pad unsafe.Pointer) interface{} {
	for k, v := range p {
		if gobool(C.padsAreEqual(k, pad)) {
			return v
		}
	}
	return nil
}

// RemoveFuncForPad will remove the function for the given pad from this map.
func (p PadFuncMap) RemoveFuncForPad(pad unsafe.Pointer) {
	for k := range p {
		if gobool(C.padsAreEqual(k, pad)) {
			delete(p, k)
		}
	}
}

// Pad is a go representation of a GstPad
type Pad struct{ *Object }

// FromGstPadUnsafeFull wraps the given pad in a reference and finalizer.
func FromGstPadUnsafeFull(pad unsafe.Pointer) *Pad {
	return &Pad{wrapObject(glib.TransferFull(pad))}
}

// FromGstPadUnsafeNone wraps the given pad in a finalizer.
func FromGstPadUnsafeNone(pad unsafe.Pointer) *Pad {
	return &Pad{wrapObject(glib.TransferNone(pad))}
}

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
	return FromGstPadUnsafeNone(unsafe.Pointer(pad))
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
	return FromGstPadUnsafeNone(unsafe.Pointer(pad))
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
	return FromGstCapsUnsafeFull(unsafe.Pointer(caps))
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
	return FlowReturn(C.gst_pad_chain(p.Instance(), buffer.Ref().Instance()))
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
	return FlowReturn(C.gst_pad_chain_list(p.Instance(), bufferList.Ref().Instance()))
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
	if parent == nil {
		return gobool(C.gst_pad_event_default(p.Instance(), nil, event.Ref().Instance()))
	}
	return gobool(C.gst_pad_event_default(p.Instance(), parent.Instance(), event.Ref().Instance()))
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
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_pad_get_allowed_caps(
		p.Instance(),
	)))
}

// GetCurrentCaps gets the capabilities currently configured on pad with the last EventCaps event.
func (p *Pad) GetCurrentCaps() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_pad_get_current_caps(
		p.Instance(),
	)))
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
	return FromGstPadTemplateUnsafeFull(unsafe.Pointer(tmpl))
}

// GetPadTemplateCaps gets the capabilities for pad's template.
func (p *Pad) GetPadTemplateCaps() *Caps {
	caps := C.gst_pad_get_pad_template_caps(p.Instance())
	if caps == nil {
		return nil
	}
	return FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// GetParentElement gets the parent of pad, cast to a Element. If a pad has no parent or its
// parent is not an element, return nil.
func (p *Pad) GetParentElement() *Element {
	elem := C.gst_pad_get_parent_element(p.Instance())
	if elem == nil {
		return nil
	}
	return FromGstElementUnsafeFull(unsafe.Pointer(elem))
}

// GetPeer gets the peer of pad. This function refs the peer pad so you need to unref it after use.
func (p *Pad) GetPeer() *Pad {
	peer := C.gst_pad_get_peer(p.Instance())
	if peer == nil {
		return nil
	}
	return FromGstPadUnsafeFull(unsafe.Pointer(peer))
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

// GetStickyEvent returns a new reference of the sticky event of type event_type from the event.
func (p *Pad) GetStickyEvent(eventType EventType, idx uint) *Event {
	ev := C.gst_pad_get_sticky_event(p.Instance(), C.GstEventType(eventType), C.guint(idx))
	if ev == nil {
		return nil
	}
	return FromGstEventUnsafeFull(unsafe.Pointer(ev))
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
	return FromGstStreamUnsafeFull(unsafe.Pointer(st))
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

// HasCurrentCaps checks if pad has caps set on it with a GST_EVENT_CAPS event.
func (p *Pad) HasCurrentCaps() bool {
	return gobool(C.gst_pad_has_current_caps(p.Instance()))
}

// IsActive queries if a pad is active
func (p *Pad) IsActive() bool {
	return gobool(C.gst_pad_is_active(p.Instance()))
}

// IsBlocked checks if the pad is blocked or not. This function returns the last requested state of the pad. It is not certain that
// the pad is actually blocking at this point (see IsBlocking).
func (p *Pad) IsBlocked() bool {
	return gobool(C.gst_pad_is_blocked(p.Instance()))
}

// IsBlocking checks if the pad is blocking or not. This is a guaranteed state of whether the pad is actually blocking on a GstBuffer
// or a GstEvent.
func (p *Pad) IsBlocking() bool {
	return gobool(C.gst_pad_is_blocking(p.Instance()))
}

// IsLinked checks if a pad is linked to another pad or not.
func (p *Pad) IsLinked() bool {
	return gobool(C.gst_pad_is_linked(p.Instance()))
}

// GetInternalLinks gets the pads to which the given pad is linked to inside of the parent element.
//
// Unref each pad after use.
func (p *Pad) GetInternalLinks() ([]*Pad, error) {
	iterator := C.gst_pad_iterate_internal_links(p.Instance())
	if iterator == nil {
		return nil, nil
	}
	return iteratorToPadSlice(iterator)
}

// GetInternalLinksDefault gets the list of pads to which the given pad is linked to inside of the parent element. This is the default
// handler, and thus returns all of the pads inside the parent element with opposite direction.
func (p *Pad) GetInternalLinksDefault(parent *Object) ([]*Pad, error) {
	iterator := C.gst_pad_iterate_internal_links_default(p.Instance(), parent.Instance())
	if iterator == nil {
		return nil, nil
	}
	return iteratorToPadSlice(iterator)
}

// Link links a sink pad to this source pad.
func (p *Pad) Link(sink *Pad) PadLinkReturn {
	return PadLinkReturn(C.gst_pad_link(p.Instance(), sink.Instance()))
}

// LinkFull links this source pad and the sink pad.
//
// This variant of Link provides a more granular control on the checks being done when linking. While providing some considerable speedups
// the caller of this method must be aware that wrong usage of those flags can cause severe issues. Refer to the documentation of GstPadLinkCheck
// for more information.
func (p *Pad) LinkFull(sink *Pad, flags PadLinkCheck) PadLinkReturn {
	return PadLinkReturn(C.gst_pad_link_full(p.Instance(), sink.Instance(), C.GstPadLinkCheck(flags)))
}

// LinkMaybeGhosting links this src to sink, creating any GstGhostPad's in between as necessary.
//
// This is a convenience function to save having to create and add intermediate GstGhostPad's as required for linking across GstBin boundaries.
//
// If src or sink pads don't have parent elements or do not share a common ancestor, the link will fail.
func (p *Pad) LinkMaybeGhosting(sink *Pad) bool {
	return gobool(C.gst_pad_link_maybe_ghosting(p.Instance(), sink.Instance()))
}

// LinkMaybeGhostingFull links this src to sink, creating any GstGhostPad's in between as necessary.
//
// This is a convenience function to save having to create and add intermediate GstGhostPad's as required for linking across GstBin boundaries.
//
// If src or sink pads don't have parent elements or do not share a common ancestor, the link will fail.
//
// Calling LinkMaybeGhostingFull with flags == PadLinkCheckDefault is the recommended way of linking pads with safety checks applied.
func (p *Pad) LinkMaybeGhostingFull(sink *Pad, flags PadLinkCheck) bool {
	return gobool(C.gst_pad_link_maybe_ghosting_full(p.Instance(), sink.Instance(), C.GstPadLinkCheck(flags)))
}

// MarkReconfigure marks this pad for needing reconfiguration. The next call to CheckReconfigure will return TRUE after this call.
func (p *Pad) MarkReconfigure() {
	C.gst_pad_mark_reconfigure(p.Instance())
}

// NeedsReconfigure checks the GST_PAD_FLAG_NEED_RECONFIGURE flag on pad and return TRUE if the flag was set.
func (p *Pad) NeedsReconfigure() bool {
	return gobool(C.gst_pad_needs_reconfigure(p.Instance()))
}

// PauseTask pauses the task of pad. This function will also wait until the function executed by the task is finished if this function is not called
// from the task function.
func (p *Pad) PauseTask() bool {
	return gobool(C.gst_pad_pause_task(p.Instance()))
}

// PeerQuery performs PadQuery on the peer of pad.
//
// The caller is responsible for both the allocation and deallocation of the query structure.
func (p *Pad) PeerQuery(query *Query) bool {
	return gobool(C.gst_pad_peer_query(p.Instance(), query.Instance()))
}

// PeerQueryAcceptCaps checks if the peer of pad accepts caps. If pad has no peer, this function returns TRUE.
func (p *Pad) PeerQueryAcceptCaps(caps *Caps) bool {
	return gobool(C.gst_pad_peer_query_accept_caps(p.Instance(), caps.Instance()))
}

// PeerQueryCaps gets the capabilities of the peer connected to this pad. Similar to QueryCaps.
//
// When called on srcpads filter contains the caps that upstream could produce in the order preferred by upstream.
// When called on sinkpads filter contains the caps accepted by downstream in the preferred order. filter might be nil but if it is not nil
// the returned caps will be a subset of filter.
func (p *Pad) PeerQueryCaps(filter *Caps) *Caps {
	var caps *C.GstCaps
	if filter == nil {
		caps = C.gst_pad_peer_query_caps(p.Instance(), nil)
	} else {
		caps = C.gst_pad_peer_query_caps(p.Instance(), filter.Instance())
	}
	if caps == nil {
		return nil
	}
	return FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// PeerQueryConvert queries the peer pad of a given sink pad to convert src_val in src_format to dest_format.
func (p *Pad) PeerQueryConvert(srcFormat, destFormat Format, srcVal int64) (bool, int64) {
	var out C.gint64
	gok := C.gst_pad_peer_query_convert(p.Instance(), C.GstFormat(srcFormat), C.gint64(srcVal), C.GstFormat(destFormat), &out)
	return gobool(gok), int64(out)
}

// PeerQueryDuration queries the peer pad of a given sink pad for the total stream duration.
func (p *Pad) PeerQueryDuration(format Format) (bool, int64) {
	var out C.gint64
	gok := C.gst_pad_peer_query_duration(p.Instance(), C.GstFormat(format), &out)
	return gobool(gok), int64(out)
}

// PeerQueryPosition queries the peer of a given sink pad for the stream position.
func (p *Pad) PeerQueryPosition(format Format) (bool, int64) {
	var out C.gint64
	gok := C.gst_pad_peer_query_position(p.Instance(), C.GstFormat(format), &out)
	return gobool(gok), int64(out)
}

// ProxyQueryAcceptCaps checks if all internally linked pads of pad accepts the caps in query and returns the intersection of the results.
//
// This function is useful as a default accept caps query function for an element that can handle any stream format, but requires caps that
// are acceptable for all opposite pads.
func (p *Pad) ProxyQueryAcceptCaps(query *Query) bool {
	return gobool(C.gst_pad_proxy_query_accept_caps(p.Instance(), query.Instance()))
}

// ProxyQueryCaps calls QueryCaps for all internally linked pads of pad and returns the intersection of the results.
//
// This function is useful as a default caps query function for an element that can handle any stream format, but requires all its pads to have
// the same caps. Two such elements are tee and adder.
func (p *Pad) ProxyQueryCaps(query *Query) bool {
	return gobool(C.gst_pad_proxy_query_caps(p.Instance(), query.Instance()))
}

// PullRange pulls a buffer from the peer pad or fills up a provided buffer.
//
// This function will first trigger the pad block signal if it was installed.
//
// When pad is not linked GST_FLOW_NOT_LINKED is returned else this function returns the result of GetRange on the peer pad. See GetRange for a list
// of return values and for the semantics of the arguments of this function.
//
// If buffer points to a variable holding nil, a valid new GstBuffer will be placed in buffer when this function returns GST_FLOW_OK. The new buffer
// must be freed with Unref after usage. When this function returns any other result value, buffer will still point to NULL.
//
// When buffer points to a variable that points to a valid GstBuffer, the buffer will be filled with the result data when this function returns GST_FLOW_OK.
// When this function returns any other result value, buffer will be unchanged. If the provided buffer is larger than size, only size bytes will be filled
// in the result buffer and its size will be updated accordingly.
//
// Note that less than size bytes can be returned in buffer when, for example, an EOS condition is near or when buffer is not large enough to hold size bytes.
// The caller should check the result buffer size to get the result size.
func (p *Pad) PullRange(offset uint64, size uint, buffer *Buffer) (FlowReturn, *Buffer) {
	var buf *C.GstBuffer
	if buffer != nil {
		buf = buffer.Instance()
	}
	ret := C.gst_pad_pull_range(p.Instance(), C.guint64(offset), C.guint(size), &buf)
	var newBuf *Buffer
	if buf != nil {
		newBuf = FromGstBufferUnsafeFull(unsafe.Pointer(buf))
	} else {
		newBuf = nil
	}
	return FlowReturn(ret), newBuf
}

// Push pushes a buffer to the peer of pad.
//
// This function will call installed block probes before triggering any installed data probes.
//
// The function proceeds calling Chain on the peer pad and returns the value from that function. If pad has no peer, GST_FLOW_NOT_LINKED will be returned.
//
// In all cases, success or failure, the caller loses its reference to buffer after calling this function.
func (p *Pad) Push(buf *Buffer) FlowReturn {
	return FlowReturn(C.gst_pad_push(p.Instance(), buf.Ref().Instance()))
}

// PushEvent sends the event to the peer of the given pad. This function is mainly used by elements to send events to their peer elements.
//
// This function takes ownership of the provided event so you should Ref it if you want to reuse the event after this call.
func (p *Pad) PushEvent(ev *Event) bool {
	return gobool(C.gst_pad_push_event(p.Instance(), ev.Ref().Instance()))
}

// PushList pushes a buffer list to the peer of pad.
//
// This function will call installed block probes before triggering any installed data probes.
//
// The function proceeds calling the chain function on the peer pad and returns the value from that function. If pad has no peer, GST_FLOW_NOT_LINKED will be
// returned. If the peer pad does not have any installed chainlist function every group buffer of the list will be merged into a normal GstBuffer and chained via
// Chain.
//
// In all cases, success or failure, the caller loses its reference to list after calling this function.
func (p *Pad) PushList(bufList *BufferList) FlowReturn {
	return FlowReturn(C.gst_pad_push_list(p.Instance(), bufList.Ref().Instance()))
}

// Query dispatches a query to a pad. The query should have been allocated by the caller via one of the type-specific allocation functions. The element that the
// pad belongs to is responsible for filling the query with an appropriate response, which should then be parsed with a type-specific query parsing function.
//
// Again, the caller is responsible for both the allocation and deallocation of the query structure.
//
// Please also note that some queries might need a running pipeline to work.
func (p *Pad) Query(query *Query) bool {
	return gobool(C.gst_pad_query(p.Instance(), query.Instance()))
}

// QueryAcceptCaps checks if the given pad accepts the caps.
func (p *Pad) QueryAcceptCaps(caps *Caps) bool {
	return gobool(C.gst_pad_query_accept_caps(p.Instance(), caps.Instance()))
}

// QueryCaps gets the capabilities this pad can produce or consume. Note that this method doesn't necessarily return the caps set by sending a NewCapsEvent - use GetCurrentCaps
// for that instead. QueryCaps returns all possible caps a pad can operate with, using the pad's CAPS query function, If the query fails, this function will return filter, if not
// NULL, otherwise ANY.
//
// When called on sinkpads filter contains the caps that upstream could produce in the order preferred by upstream. When called on srcpads filter contains the caps accepted by
// downstream in the preferred order. filter might be NULL but if it is not NULL the returned caps will be a subset of filter.
//
// Note that this function does not return writable GstCaps, use gst_caps_make_writable before modifying the caps.
func (p *Pad) QueryCaps(filter *Caps) *Caps {
	var caps *C.GstCaps
	if filter == nil {
		caps = C.gst_pad_query_caps(p.Instance(), nil)
	} else {
		caps = C.gst_pad_query_caps(p.Instance(), filter.Instance())
	}
	if caps == nil {
		return nil
	}
	return FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// QueryConvert queries a pad to convert src_val in src_format to dest_format.
func (p *Pad) QueryConvert(srcFormat, destFormat Format, srcVal int64) (bool, int64) {
	var out C.gint64
	gok := C.gst_pad_query_convert(p.Instance(), C.GstFormat(srcFormat), C.gint64(srcVal), C.GstFormat(destFormat), &out)
	return gobool(gok), int64(out)
}

// QueryDefault invokes the default query handler for the given pad. The query is sent to all pads internally linked to pad. Note that if there are many possible sink pads that are
// internally linked to pad, only one will be sent the query. Multi-sinkpad elements should implement custom query handlers.
func (p *Pad) QueryDefault(parent *Object, query *Query) bool {
	return gobool(C.gst_pad_query_default(p.Instance(), parent.Instance(), query.Instance()))
}

// QueryDuration queries a pad for the total stream duration.
func (p *Pad) QueryDuration(format Format) (bool, int64) {
	var out C.gint64
	gok := C.gst_pad_query_duration(p.Instance(), C.GstFormat(format), &out)
	return gobool(gok), int64(out)
}

// QueryPosition queries a pad for the stream position.
func (p *Pad) QueryPosition(format Format) (bool, int64) {
	var out C.gint64
	gok := C.gst_pad_query_position(p.Instance(), C.GstFormat(format), &out)
	return gobool(gok), int64(out)
}

// RemoveProbe removes the probe with id from pad.
func (p *Pad) RemoveProbe(id uint64) {
	C.gst_pad_remove_probe(p.Instance(), C.gulong(id))
}

// SendEvent sends the event to the pad. This function can be used by applications to send events in the pipeline.
//
// If pad is a source pad, event should be an upstream event. If pad is a sink pad, event should be a downstream event.
// For example, you would not send a GST_EVENT_EOS on a src pad; EOS events only propagate downstream. Furthermore, some
// downstream events have to be serialized with data flow, like EOS, while some can travel out-of-band, like GST_EVENT_FLUSH_START.
// If the event needs to be serialized with data flow, this function will take the pad's stream lock while calling its event function.
//
// To find out whether an event type is upstream, downstream, or downstream and serialized, see GstEventTypeFlags, gst_event_type_get_flags,
// GST_EVENT_IS_UPSTREAM, GST_EVENT_IS_DOWNSTREAM, and GST_EVENT_IS_SERIALIZED. Note that in practice that an application or plugin doesn't
// need to bother itself with this information; the core handles all necessary locks and checks.
//
// This function takes ownership of the provided event so you should gst_event_ref it if you want to reuse the event after this call.
func (p *Pad) SendEvent(ev *Event) bool {
	return gobool(C.gst_pad_send_event(p.Instance(), ev.Ref().Instance()))
}

func (p *Pad) registerCallback(f interface{}, fmap PadFuncMap) *C.PadDestroyNotifyInfo {
	selfPtr := p.Unsafe()

	fmap[selfPtr] = f

	notifyInfo := (*C.PadDestroyNotifyInfo)(C.malloc(C.sizeof_PadDestroyNotifyInfo))
	notifyInfo.pad_ptr = (C.gpointer)(gopointer.Save(selfPtr))
	notifyInfo.func_map_ptr = (C.gpointer)(gopointer.Save(fmap))

	return notifyInfo
}

// PadActivateFunc is called when the pad is activated during the element READY to PAUSED state change. By default
// this function will call the activate function that puts the pad in push mode but elements can override this
// function to activate the pad in pull mode if they wish.
type PadActivateFunc func(self *Pad, parent *Object) bool

var padActivateFuncs PadFuncMap = make(PadFuncMap)

// SetActivateFunction sets the given active function on the pad. The activate function will dispatch to ActivateMode to perform
// the actual activation. Only makes sense to set on sink pads.
func (p *Pad) SetActivateFunction(f PadActivateFunc) {
	notifyInfo := p.registerCallback(f, padActivateFuncs)
	C.gst_pad_set_activate_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadActivateFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// PadActivateModeFunc is the prototype of the push and pull activate functions.
type PadActivateModeFunc func(self *Pad, parent *Object, mode PadMode, active bool) bool

var padActivateModeFuncs PadFuncMap = make(PadFuncMap)

// SetActivateModeFunction sets the given activate_mode function for the pad. An activate_mode function prepares
// the element for data passing.
func (p *Pad) SetActivateModeFunction(f PadActivateModeFunc) {
	notifyInfo := p.registerCallback(f, padActivateModeFuncs)
	C.gst_pad_set_activatemode_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadActivateModeFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// SetActive activates or deactivates the given pad. Normally called from within core state change functions.
//
// If active, makes sure the pad is active. If it is already active, either in push or pull mode, just return. Otherwise dispatches to the
// pad's activate function to perform the actual activation.
//
// If not active, calls ActivateMode with the pad's current mode and a FALSE argument.
func (p *Pad) SetActive(active bool) bool {
	return gobool(C.gst_pad_set_active(p.Instance(), gboolean(active)))
}

// PadChainFunc is a function that will be called on sinkpads when chaining buffers. The function typically
// processes the data contained in the buffer and either consumes the data or passes it on to the internally
// linked pad(s).
//
// The GStreamer API will produce a refcount to buffer, and the bindings will automatically call Unref after the
// function returns. If you'd like to hold on to the buffer longer, add an additional ref yourself or create a copy.
//
// When a chain function detects an error in the data stream, it must post an error on the bus and return an
// appropriate FlowReturn value.
type PadChainFunc func(self *Pad, parent *Object, buffer *Buffer) FlowReturn

var padChainFuncs PadFuncMap = make(PadFuncMap)

// SetChainFunction sets the given chain function for the pad. The chain function is called to process an input buffer.
func (p *Pad) SetChainFunction(f PadChainFunc) {
	notifyInfo := p.registerCallback(f, padChainFuncs)
	C.gst_pad_set_chain_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadChainFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// PadChainListFunc is a function that will be called on sinkpads when chaining buffer lists. The function typically
// processes the data contained in the buffer list and either consumes the data or passes it on to the internally
// linked pad(s).
//
// The GStreamer API will produce a refcount to list, and the bindings will automatically call Unref after the
// function returns. If you'd like to hold on to the list longer, add an additional ref yourself or create a copy.
//
// When a chainlist function detects an error in the data stream, it must post an error on the bus and return an
// appropriate FlowReturn value.
type PadChainListFunc func(self *Pad, parent *Object, list *BufferList) FlowReturn

var padChainListFuncs PadFuncMap = make(PadFuncMap)

// SetChainListFunction sets the given chain function for the pad. The chain function is called to process an input buffer list.
func (p *Pad) SetChainListFunction(f PadChainListFunc) {
	notifyInfo := p.registerCallback(f, padChainListFuncs)
	C.gst_pad_set_chain_list_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadChainListFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// SetElementPrivate sets the given private data pointer on the pad. This function can only be used by the element that owns the pad.
// No locking is performed in this function.
func (p *Pad) SetElementPrivate(data interface{}) {
	ptr := gopointer.Save(data)
	C.gst_pad_set_element_private(p.Instance(), (C.gpointer)(unsafe.Pointer(ptr)))
}

// PadEventFullFunc is the function signature to handle an event for the pad. This variant is for specific elements that
// will take into account the last downstream flow return (from a pad push), in which case they can return it.
type PadEventFullFunc func(self *Pad, parent *Object, event *Event) FlowReturn

var padEventFullFuncs PadFuncMap = make(PadFuncMap)

// SetEventFullFunction sets the given event handler for the pad.
func (p *Pad) SetEventFullFunction(f PadEventFullFunc) {
	notifyInfo := p.registerCallback(f, padEventFullFuncs)
	C.gst_pad_set_event_full_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadEventFullFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// PadEventFunc is the function signature to handle an event for the pad
type PadEventFunc func(self *Pad, parent *Object, event *Event) bool

var padEventFuncs PadFuncMap = make(PadFuncMap)

// SetEventFunction sets the given event handler for the pad.
func (p *Pad) SetEventFunction(f PadEventFunc) {
	notifyInfo := p.registerCallback(f, padEventFuncs)
	C.gst_pad_set_event_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadEventFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// PadGetRangeFunc will be called on source pads when a peer element request a buffer at the specified offset and length.
// If this function returns FlowOK it must also return a valid buffer.
//
// This function is installed on a source pad with SetGetRangeFunction and can only be called on source pads after they are
// successfully activated with ActivateMode with the PadModePull.
//
// offset and length are always given in byte units. offset must normally be a value between 0 and the length in bytes of
// the data available on pad. The length (duration in bytes) can be retrieved with a GST_QUERY_DURATION or with a
// GST_QUERY_SEEKING.
//
// Any offset larger or equal than the length will make the function return FlowEOS, which corresponds to EOS. In this case
// buffer does not contain a valid buffer.
//
// The buffer size of buffer will only be smaller than length when offset is near the end of the stream. In all other cases,
// the size of buffer must be exactly the requested size.
//
// It is allowed to call this function with a 0 length and valid offset, in which case buffer will contain a 0-sized buffer
// and the function returns FlowOK.
//
// When this function is called with a -1 offset, the sequentially next buffer of length length in the stream is returned.
//
// When this function is called with a -1 length, a buffer with a default optimal length is returned in buffer. The length
// might depend on the value of offset.
type PadGetRangeFunc func(self *Pad, parent *Object, offset uint64, length uint) (FlowReturn, *Buffer)

var padGetRangeFuncs PadFuncMap = make(PadFuncMap)

// SetGetRangeFunction sets the given getrange function for the pad. The getrange function is called to produce a new Buffer
// to start the processing pipeline. see PadGetRangeFunc for a description of the getrange function.
func (p *Pad) SetGetRangeFunction(f PadGetRangeFunc) {
	notifyInfo := p.registerCallback(f, padGetRangeFuncs)
	C.gst_pad_set_getrange_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadGetRangeFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// PadIterIntLinkFunc is the signature of the internal pad link iterator function. It should return a slice of pads that are
// linked to the given pad on the inside of the parent element. Currently the function only supports returning
// a single pad until GstIterator is implemented fully in the bindings.
type PadIterIntLinkFunc func(self *Pad, parent *Object) []*Pad

var padIterIntLinkFuncs PadFuncMap = make(PadFuncMap)

// SetIterIntLinkFunction sets the given internal link iterator function for the pad.
func (p *Pad) SetIterIntLinkFunction(f PadIterIntLinkFunc) {
	notifyInfo := p.registerCallback(f, padIterIntLinkFuncs)
	C.gst_pad_set_iterate_internal_links_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadIterIntLinkFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// PadLinkFunc is the function signature to handle a new link on a pad.
type PadLinkFunc func(self *Pad, parent *Object, peer *Pad) PadLinkReturn

var padLinkFuncs PadFuncMap = make(PadFuncMap)

// SetLinkFunction sets the given link function for the pad. It will be called when the pad is linked with another pad.
//
// The return value PadLinkOK should be used when the connection can be made.
//
// The return value PadLinkRefused should be used when the connection cannot be made for some reason.
//
// If link is installed on a source pad, it should call the Link of the peer sink pad, if present.
func (p *Pad) SetLinkFunction(f PadLinkFunc) {
	notifyInfo := p.registerCallback(f, padLinkFuncs)
	C.gst_pad_set_link_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadLinkFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// SetOffset sets the offset that will be applied to the running time of pad.
func (p *Pad) SetOffset(offset int64) {
	C.gst_pad_set_offset(p.Instance(), C.gint64(offset))
}

// PadQueryFunc is a function for performing queries on a pad. It should return true if it could handle the query.
type PadQueryFunc func(self *Pad, parent *Object, query *Query) bool

var padQueryFuncs PadFuncMap = make(PadFuncMap)

// SetQueryFunction sets the query handler for the pad.
func (p *Pad) SetQueryFunction(f PadQueryFunc) {
	notifyInfo := p.registerCallback(f, padQueryFuncs)
	C.gst_pad_set_query_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadQueryFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// PadUnlinkFunc is the function signature to handle a unlinking the pad from its peer.
//
// The pad's lock is already held when the unlink function is called, so most pad functions cannot be called from within
// the callback.
type PadUnlinkFunc func(self *Pad, parent *Object)

var padUnlinkFuncs PadFuncMap = make(PadFuncMap)

// SetUnlinkFunction sets the given unlink function for the pad. It will be called when the pad is unlinked.
//
// Note that the pad's lock is already held when the unlink function is called, so most pad functions cannot be called
// from within the callback.
func (p *Pad) SetUnlinkFunction(f PadUnlinkFunc) {
	notifyInfo := p.registerCallback(f, padUnlinkFuncs)
	C.gst_pad_set_unlink_function_full(
		p.Instance(),
		C.GstPadActivateFunction(C.cgoGstPadUnlinkFunction),
		(C.gpointer)(unsafe.Pointer(notifyInfo)),
		C.GDestroyNotify(C.cgoPadFuncDestroyNotify),
	)
}

// StickyEventsForEachFunc is a callback used by StickyEventsForEach. When this function returns TRUE, the next event will be returned.
// When FALSE is returned, gst_pad_sticky_events_foreach will return.
//
// When event is set to NULL, the item will be removed from the list of sticky events. event can be replaced by assigning a new reference to it.
// This function is responsible for unreffing the old event when removing or modifying.
type StickyEventsForEachFunc func(pad *Pad, event *Event) bool

// StickyEventsForEach iterates all sticky events on pad and calls foreach_func for every event. If foreach_func returns FALSE the iteration is
// immediately stopped.
func (p *Pad) StickyEventsForEach(f StickyEventsForEachFunc) {
	ptr := gopointer.Save(f)
	defer gopointer.Unref(ptr)
	C.gst_pad_sticky_events_foreach(
		p.Instance(),
		C.GstPadStickyEventsForeachFunction(C.cgoPadStickyEventForEachFunc),
		(C.gpointer)(unsafe.Pointer(ptr)),
	)
}

// StoreStickyEvent stores the sticky event on pad
func (p *Pad) StoreStickyEvent(ev *Event) FlowReturn {
	return FlowReturn(C.gst_pad_store_sticky_event(p.Instance(), ev.Instance()))
}

// Unlink unlinks this source pad from the sink pad. Will emit the unlinked signal on both pads.
func (p *Pad) Unlink(pad *Pad) bool {
	return gobool(C.gst_pad_unlink(p.Instance(), pad.Instance()))
}

// UseFixedCaps is a helper function you can use that sets the FIXED_CAPS flag This way the default CAPS query will always return the negotiated caps
// or in case the pad is not negotiated, the padtemplate caps.
//
// The negotiated caps are the caps of the last CAPS event that passed on the pad. Use this function on a pad that, once it negotiated to a CAPS, cannot
// be renegotiated to something else.
func (p *Pad) UseFixedCaps() { C.gst_pad_use_fixed_caps(p.Instance()) }

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

func iteratorToPadSlice(iterator *C.GstIterator) ([]*Pad, error) {
	pads := make([]*Pad, 0)
	gval := new(C.GValue)

	for {
		switch C.gst_iterator_next((*C.GstIterator)(iterator), (*C.GValue)(unsafe.Pointer(gval))) {
		case C.GST_ITERATOR_DONE:
			C.gst_iterator_free((*C.GstIterator)(iterator))
			return pads, nil
		case C.GST_ITERATOR_RESYNC:
			C.gst_iterator_resync((*C.GstIterator)(iterator))
		case C.GST_ITERATOR_OK:
			cPadVoid := C.g_value_get_object((*C.GValue)(gval))
			cPad := (*C.GstPad)(cPadVoid)
			pads = append(pads, FromGstPadUnsafeNone(unsafe.Pointer(cPad)))
			C.g_value_unset((*C.GValue)(gval))
		default:
			return nil, errors.New("Pad iterator failed")
		}
	}
}
