package base

/*
#include "gst.go.h"

extern gboolean       goGstBaseTransformAcceptCaps           (GstBaseTransform * self, GstPadDirection direction, GstCaps * caps);
extern void           goGstBaseTransformBeforeTransform      (GstBaseTransform * self, GstBuffer * buffer);
extern gboolean       goGstBaseTransformCopyMetadata         (GstBaseTransform * self, GstBuffer * input, GstBuffer * output);
extern gboolean       goGstBaseTransformDecideAllocation     (GstBaseTransform * self, GstQuery * query);
extern gboolean       goGstBaseTransformFilterMeta           (GstBaseTransform * self, GstQuery * query, GType api, const GstStructure * params);
extern GstCaps *      goGstBaseTransformFixateCaps           (GstBaseTransform * self, GstPadDirection direction, GstCaps * caps, GstCaps * othercaps);
extern GstFlowReturn  goGstBaseTransformGenerateOutput       (GstBaseTransform * self, GstBuffer ** buf);
extern gboolean       goGstBaseTransformGetUnitSize          (GstBaseTransform * self, GstCaps * caps, gsize * size);
extern GstFlowReturn  goGstBaseTransformPrepareOutputBuffer  (GstBaseTransform * self, GstBuffer * input, GstBuffer ** output);
extern gboolean       goGstBaseTransformProposeAllocation    (GstBaseTransform * self, GstQuery * decide, GstQuery * query);
extern gboolean       goGstBaseTransformQuery                (GstBaseTransform * self, GstPadDirection direction, GstQuery * query);
extern gboolean       goGstBaseTransformSetCaps              (GstBaseTransform * self, GstCaps * incaps, GstCaps * outcaps);
extern gboolean       goGstBaseTransformSinkEvent            (GstBaseTransform * self, GstEvent * event);
extern gboolean       goGstBaseTransformSrcEvent             (GstBaseTransform * self, GstEvent * event);
extern gboolean       goGstBaseTransformStart                (GstBaseTransform * self);
extern gboolean       goGstBaseTransformStop                 (GstBaseTransform * self);
extern GstFlowReturn  goGstBaseTransformSubmitInputBuffer    (GstBaseTransform * self, gboolean discont, GstBuffer * input);
extern GstFlowReturn  goGstBaseTransformTransform            (GstBaseTransform * self, GstBuffer * inbuf, GstBuffer * outbuf);
extern GstCaps *      goGstBaseTransformTransformCaps        (GstBaseTransform * self, GstPadDirection direction, GstCaps * caps, GstCaps * filter);
extern GstFlowReturn  goGstBaseTransformTransformIP          (GstBaseTransform * self, GstBuffer * buffer);
extern gboolean       goGstBaseTransformTransformMeta        (GstBaseTransform * self, GstBuffer * outbuf, GstMeta * meta, GstBuffer * inbuf);
extern gboolean       goGstBaseTransformTransformSize        (GstBaseTransform * self, GstPadDirection direction, GstCaps * caps, gsize size, GstCaps * othercaps, gsize * othersize);

void  setGstBaseTransformAcceptCaps           (GstBaseTransformClass * klass) { klass->accept_caps = goGstBaseTransformAcceptCaps; }
void  setGstBaseTransformBeforeTransform      (GstBaseTransformClass * klass) { klass->before_transform = goGstBaseTransformBeforeTransform; }
void  setGstBaseTransformCopyMetadata         (GstBaseTransformClass * klass) { klass->copy_metadata = goGstBaseTransformCopyMetadata; }
void  setGstBaseTransformDecideAllocation     (GstBaseTransformClass * klass) { klass->decide_allocation = goGstBaseTransformDecideAllocation; }
void  setGstBaseTransformFilterMeta           (GstBaseTransformClass * klass) { klass->filter_meta = goGstBaseTransformFilterMeta; }
void  setGstBaseTransformFixateCaps           (GstBaseTransformClass * klass) { klass->fixate_caps = goGstBaseTransformFixateCaps; }
void  setGstBaseTransformGenerateOutput       (GstBaseTransformClass * klass) { klass->generate_output = goGstBaseTransformGenerateOutput; }
void  setGstBaseTransformGetUnitSize          (GstBaseTransformClass * klass) { klass->get_unit_size = goGstBaseTransformGetUnitSize; }
void  setGstBaseTransformPrepareOutputBuffer  (GstBaseTransformClass * klass) { klass->prepare_output_buffer = goGstBaseTransformPrepareOutputBuffer; }
void  setGstBaseTransformProposeAllocation    (GstBaseTransformClass * klass) { klass->propose_allocation = goGstBaseTransformProposeAllocation; }
void  setGstBaseTransformQuery                (GstBaseTransformClass * klass) { klass->query = goGstBaseTransformQuery; }
void  setGstBaseTransformSetCaps              (GstBaseTransformClass * klass) { klass->set_caps = goGstBaseTransformSetCaps; }
void  setGstBaseTransformSinkEvent            (GstBaseTransformClass * klass) { klass->sink_event = goGstBaseTransformSinkEvent; }
void  setGstBaseTransformSrcEvent             (GstBaseTransformClass * klass) { klass->src_event = goGstBaseTransformSrcEvent; }
void  setGstBaseTransformStart                (GstBaseTransformClass * klass) { klass->start = goGstBaseTransformStart; }
void  setGstBaseTransformStop                 (GstBaseTransformClass * klass) { klass->stop = goGstBaseTransformStop; }
void  setGstBaseTransformSubmitInputBuffer    (GstBaseTransformClass * klass) { klass->submit_input_buffer = goGstBaseTransformSubmitInputBuffer; }
void  setGstBaseTransformTransform            (GstBaseTransformClass * klass) { klass->transform = goGstBaseTransformTransform; }
void  setGstBaseTransformTransformCaps        (GstBaseTransformClass * klass) { klass->transform_caps = goGstBaseTransformTransformCaps; }
void  setGstBaseTransformTransformIP          (GstBaseTransformClass * klass) { klass->transform_ip = goGstBaseTransformTransformIP; }
void  setGstBaseTransformTransformMeta        (GstBaseTransformClass * klass) { klass->transform_meta = goGstBaseTransformTransformMeta; }
void  setGstBaseTransformTransformSize        (GstBaseTransformClass * klass) { klass->transform_size = goGstBaseTransformTransformSize; }

*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

var (
	// ExtendsBaseTransform is an Extendable for extending a GstBaseTransform
	ExtendsBaseTransform glib.Extendable = &extendsBaseTransform{parent: gst.ExtendsElement}
)

// GstBaseTransformImpl is the interface for an element extending a GstBaseTransform.
// Subclasses can override any of the available virtual methods or not, as needed. At minimum either
// Transform or TransformIP need to be overridden. If the element can overwrite the input data with
// the results (data is of the same type and quantity) it should provide TransformIP.
//
// For more information:
// https://gstreamer.freedesktop.org/documentation/base/gstbasetransform.html?gi-language=c
type GstBaseTransformImpl interface {
	// Optional. Subclasses can override this method to check if caps can be handled by the element.
	// The default implementation might not be the most optimal way to check this in all cases.
	AcceptCaps(self *GstBaseTransform, direction gst.PadDirection, caps *gst.Caps) bool
	// Optional. This method is called right before the base class will start processing. Dynamic
	// properties or other delayed configuration could be performed in this method.
	BeforeTransform(self *GstBaseTransform, buffer *gst.Buffer)
	// Optional. Copy the metadata from the input buffer to the output buffer. The default implementation
	// will copy the flags, timestamps and offsets of the buffer.
	CopyMetadata(self *GstBaseTransform, input, output *gst.Buffer) bool
	// Setup the allocation parameters for allocating output buffers. The passed in query contains the
	// result of the downstream allocation query. This function is only called when not operating in
	// passthrough mode. The default implementation will remove all memory dependent metadata. If there is
	// a FilterMeta method implementation, it will be called for all metadata API in the downstream query,
	// otherwise the metadata API is removed.
	DecideAllocation(self *GstBaseTransform, query *gst.Query) bool
	// Return TRUE if the metadata API should be proposed in the upstream allocation query. The default
	// implementation is NULL and will cause all metadata to be removed.
	FilterMeta(self *GstBaseTransform, query *gst.Query, api glib.Type, params *gst.Structure) bool
	// Optional. Given the pad in this direction and the given caps, fixate the caps on the other pad.
	// The function returns a fixated version of othercaps. othercaps itself is not guaranteed to be writable
	// and the bindings will Unref them after the callback is complete. So if you want to return othercaps
	// with small modifications, take a copy first with Caps.Copy(), otherwise return a Ref().
	FixateCaps(self *GstBaseTransform, directon gst.PadDirection, caps *gst.Caps, othercaps *gst.Caps) *gst.Caps
	// Called after each new input buffer is submitted repeatedly until it either generates an error or fails to
	// generate an output buffer. The default implementation takes the contents of the queued_buf variable,
	// generates an output buffer if needed by calling the class PrepareOutputBuffer, and then calls either
	// Transform() or TransformIP(). Elements that don't do 1-to-1 transformations of input to output buffers can
	// either return GstBaseTransformFlowDropped or simply not generate an output buffer until they are ready to do
	// so. (Since: 1.6)
	GenerateOutput(self *GstBaseTransform) (gst.FlowReturn, *gst.Buffer)
	// Required if the transform is not in-place. Get the size in bytes of one unit for the given caps.
	GetUnitSize(self *GstBaseTransform, caps *gst.Caps) (ok bool, size int64)
	// Optional. Subclasses can override this to do their own allocation of output buffers. Elements that only
	// do analysis can return a subbuffer or even just return a reference to the input buffer (if in passthrough
	// mode). The default implementation will use the negotiated allocator or bufferpool and TransformSize to
	// allocate an output buffer or it will return the input buffer in passthrough mode.
	PrepareOutputBuffer(self *GstBaseTransform, input *gst.Buffer) (gst.FlowReturn, *gst.Buffer)
	// Propose buffer allocation parameters for upstream elements. This function must be implemented if the element
	// reads or writes the buffer content. The query that was passed to the DecideAllocation is passed in this
	// method (or nil when the element is in passthrough mode). The default implementation will pass the query
	// downstream when in passthrough mode and will copy all the filtered metadata API in non-passthrough mode.
	ProposeAllocation(self *GstBaseTransform, decideQuery, query *gst.Query) bool
	// Optional. Handle a requested query. Subclasses that implement this must chain up to the parent if they
	// didn't handle the query
	Query(self *GstBaseTransform, direction gst.PadDirection, query *gst.Query) bool
	// Allows the subclass to be notified of the actual caps set.
	SetCaps(self *GstBaseTransform, incaps, outcaps *gst.Caps) bool
	// Optional. Event handler on the sink pad. The default implementation handles the event and forwards it
	// downstream.
	SinkEvent(self *GstBaseTransform, event *gst.Event) bool
	// Optional. Event handler on the source pad. The default implementation handles the event and forwards it
	// upstream.
	SrcEvent(self *GstBaseTransform, event *gst.Event) bool
	// Optional. Called when the element starts processing. Allows opening external resources.
	Start(self *GstBaseTransform) bool
	// Optional. Called when the element stops processing. Allows closing external resources.
	Stop(self *GstBaseTransform) bool
	// Function which accepts a new input buffer and pre-processes it. The default implementation performs caps
	// (re)negotiation, then QoS if needed, and places the input buffer into the queued_buf member variable. If
	// the buffer is dropped due to QoS, it returns GstBaseTransformFlowDropped. If this input buffer is not
	// contiguous with any previous input buffer, then isDiscont is set to TRUE. (Since: 1.6)
	SubmitInputBuffer(self *GstBaseTransform, isDiscont bool, input *gst.Buffer) gst.FlowReturn
	// Required if the element does not operate in-place. Transforms one incoming buffer to one outgoing buffer.
	// The function is allowed to change size/timestamp/duration of the outgoing buffer.
	Transform(self *GstBaseTransform, inbuf, outbuf *gst.Buffer) gst.FlowReturn
	// Optional. Given the pad in this direction and the given caps, what caps are allowed on the other pad in
	// this element ?
	TransformCaps(self *GstBaseTransform, direction gst.PadDirection, caps, filter *gst.Caps) *gst.Caps
	// Required if the element operates in-place. Transform the incoming buffer in-place.
	TransformIP(self *GstBaseTransform, buf *gst.Buffer) gst.FlowReturn
	// Optional. Transform the metadata on the input buffer to the output buffer. By default this method copies
	// all meta without tags. Subclasses can implement this method and return TRUE if the metadata is to be copied.
	TransformMeta(self *GstBaseTransform, outbuf *gst.Buffer, meta *gst.Meta, inbuf *gst.Buffer) bool
	// Optional. Given the size of a buffer in the given direction with the given caps, calculate the size in
	// bytes of a buffer on the other pad with the given other caps. The default implementation uses GetUnitSize
	// and keeps the number of units the same.
	TransformSize(self *GstBaseTransform, direction gst.PadDirection, caps *gst.Caps, size int64, othercaps *gst.Caps) (ok bool, othersize int64)
}

type extendsBaseTransform struct{ parent glib.Extendable }

func (e *extendsBaseTransform) Type() glib.Type     { return glib.Type(C.gst_base_transform_get_type()) }
func (e *extendsBaseTransform) ClassSize() int64    { return int64(C.sizeof_GstBaseTransformClass) }
func (e *extendsBaseTransform) InstanceSize() int64 { return int64(C.sizeof_GstBaseTransform) }

func (e *extendsBaseTransform) InitClass(klass unsafe.Pointer, elem glib.GoObjectSubclass) {
	e.parent.InitClass(klass, elem)

	class := C.toGstBaseTransformClass(klass)

	if _, ok := elem.(interface {
		AcceptCaps(self *GstBaseTransform, direction gst.PadDirection, caps *gst.Caps) bool
	}); ok {
		C.setGstBaseTransformAcceptCaps(class)
	}

	if _, ok := elem.(interface {
		BeforeTransform(self *GstBaseTransform, buffer *gst.Buffer)
	}); ok {
		C.setGstBaseTransformBeforeTransform(class)
	}

	if _, ok := elem.(interface {
		CopyMetadata(self *GstBaseTransform, input, output *gst.Buffer) bool
	}); ok {
		C.setGstBaseTransformCopyMetadata(class)
	}

	if _, ok := elem.(interface {
		DecideAllocation(self *GstBaseTransform, query *gst.Query) bool
	}); ok {
		C.setGstBaseTransformDecideAllocation(class)
	}

	if _, ok := elem.(interface {
		FilterMeta(self *GstBaseTransform, query *gst.Query, api glib.Type, params *gst.Structure) bool
	}); ok {
		C.setGstBaseTransformFilterMeta(class)
	}

	if _, ok := elem.(interface {
		FixateCaps(self *GstBaseTransform, directon gst.PadDirection, caps *gst.Caps, othercaps *gst.Caps) *gst.Caps
	}); ok {
		C.setGstBaseTransformFixateCaps(class)
	}

	if _, ok := elem.(interface {
		GenerateOutput(self *GstBaseTransform) (gst.FlowReturn, *gst.Buffer)
	}); ok {
		C.setGstBaseTransformGenerateOutput(class)
	}

	if _, ok := elem.(interface {
		GetUnitSize(self *GstBaseTransform, caps *gst.Caps) (ok bool, size int64)
	}); ok {
		C.setGstBaseTransformGetUnitSize(class)
	}

	if _, ok := elem.(interface {
		PrepareOutputBuffer(self *GstBaseTransform, input *gst.Buffer) (gst.FlowReturn, *gst.Buffer)
	}); ok {
		C.setGstBaseTransformPrepareOutputBuffer(class)
	}

	if _, ok := elem.(interface {
		ProposeAllocation(self *GstBaseTransform, decideQuery, query *gst.Query) bool
	}); ok {
		C.setGstBaseTransformProposeAllocation(class)
	}

	if _, ok := elem.(interface {
		Query(self *GstBaseTransform, direction gst.PadDirection, query *gst.Query) bool
	}); ok {
		C.setGstBaseTransformQuery(class)
	}

	if _, ok := elem.(interface {
		SetCaps(self *GstBaseTransform, incaps, outcaps *gst.Caps) bool
	}); ok {
		C.setGstBaseTransformSetCaps(class)
	}

	if _, ok := elem.(interface {
		SinkEvent(self *GstBaseTransform, event *gst.Event) bool
	}); ok {
		C.setGstBaseTransformSinkEvent(class)
	}

	if _, ok := elem.(interface {
		SrcEvent(self *GstBaseTransform, event *gst.Event) bool
	}); ok {
		C.setGstBaseTransformSrcEvent(class)
	}

	if _, ok := elem.(interface {
		Start(self *GstBaseTransform) bool
	}); ok {
		C.setGstBaseTransformStart(class)
	}

	if _, ok := elem.(interface {
		Stop(self *GstBaseTransform) bool
	}); ok {
		C.setGstBaseTransformStop(class)
	}

	if _, ok := elem.(interface {
		SubmitInputBuffer(self *GstBaseTransform, isDiscont bool, input *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstBaseTransformSubmitInputBuffer(class)
	}

	if _, ok := elem.(interface {
		Transform(self *GstBaseTransform, inbuf, outbuf *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstBaseTransformTransform(class)
	}

	if _, ok := elem.(interface {
		TransformCaps(self *GstBaseTransform, direction gst.PadDirection, caps, filter *gst.Caps) *gst.Caps
	}); ok {
		C.setGstBaseTransformTransformCaps(class)
	}

	if _, ok := elem.(interface {
		TransformIP(self *GstBaseTransform, buf *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstBaseTransformTransformIP(class)
	}

	if _, ok := elem.(interface {
		TransformMeta(self *GstBaseTransform, outbuf *gst.Buffer, meta *gst.Meta, inbuf *gst.Buffer) bool
	}); ok {
		C.setGstBaseTransformTransformMeta(class)
	}

	if _, ok := elem.(interface {
		TransformSize(self *GstBaseTransform, direction gst.PadDirection, caps *gst.Caps, size int64, othercaps *gst.Caps) (ok bool, othersize int64)
	}); ok {
		C.setGstBaseTransformTransformSize(class)
	}

}
