package base

/*
#include "gst.go.h"

extern GstCaps *      goGstBaseSrcGetCaps             (GstBaseSrc * src, GstCaps * caps);
extern gboolean       goGstBaseSrcNegotiate           (GstBaseSrc * src);
extern GstCaps *      goGstBaseSrcFixate              (GstBaseSrc * src, GstCaps * caps);
extern gboolean       goGstBaseSrcSetCaps             (GstBaseSrc * src, GstCaps * filter);
extern gboolean       goGstBaseSrcDecideAllocation    (GstBaseSrc * src, GstQuery * query);
extern gboolean       goGstBaseSrcStart               (GstBaseSrc * src);
extern gboolean       goGstBaseSrcStop                (GstBaseSrc * src);
extern void           goGstBaseSrcGetTimes            (GstBaseSrc * src, GstBuffer * buffer, GstClockTime * start, GstClockTime * end);
extern gboolean       goGstBaseSrcGetSize             (GstBaseSrc * src, guint64 * size);
extern gboolean       goGstBaseSrcIsSeekable          (GstBaseSrc * src);
extern gboolean       goGstBaseSrcPrepareSeekSegment  (GstBaseSrc * src, GstEvent * seek, GstSegment * segment);
extern gboolean       goGstBaseSrcDoSeek              (GstBaseSrc * src, GstSegment * segment);
extern gboolean       goGstBaseSrcUnlock              (GstBaseSrc * src);
extern gboolean       goGstBaseSrcUnlockStop          (GstBaseSrc * src);
extern gboolean       goGstBaseSrcQuery               (GstBaseSrc * src, GstQuery * query);
extern gboolean       goGstBaseSrcEvent               (GstBaseSrc * src, GstEvent * event);
extern GstFlowReturn  goGstBaseSrcCreate              (GstBaseSrc * src, guint64 offset, guint size, GstBuffer ** buffer);
extern GstFlowReturn  goGstBaseSrcAlloc               (GstBaseSrc * src, guint64 offset, guint size, GstBuffer ** buffer);
extern GstFlowReturn  goGstBaseSrcFill                (GstBaseSrc * src, guint64 offset, guint size, GstBuffer * buffer);


void setGstBaseSrcGetCaps              (GstBaseSrcClass * klass)  { klass->get_caps = goGstBaseSrcGetCaps; }
void setGstBaseSrcNegotiate            (GstBaseSrcClass * klass)  { klass->negotiate = goGstBaseSrcNegotiate; }
void setGstBaseSrcFixate               (GstBaseSrcClass * klass)  { klass->fixate = goGstBaseSrcFixate; }
void setGstBaseSrcSetCaps              (GstBaseSrcClass * klass)  { klass->set_caps = goGstBaseSrcSetCaps; }
void setGstBaseSrcDecideAllocation     (GstBaseSrcClass * klass)  { klass->decide_allocation = goGstBaseSrcDecideAllocation; }
void setGstBaseSrcStart                (GstBaseSrcClass * klass)  { klass->start = goGstBaseSrcStart; }
void setGstBaseSrcStop                 (GstBaseSrcClass * klass)  { klass->stop = goGstBaseSrcStop; }
void setGstBaseSrcGetTimes             (GstBaseSrcClass * klass)  { klass->get_times = goGstBaseSrcGetTimes; }
void setGstBaseSrcGetSize              (GstBaseSrcClass * klass)  { klass->get_size = goGstBaseSrcGetSize; }
void setGstBaseSrcIsSeekable           (GstBaseSrcClass * klass)  { klass->is_seekable = goGstBaseSrcIsSeekable; }
void setGstBaseSrcPrepareSeekSegment   (GstBaseSrcClass * klass)  { klass->prepare_seek_segment = goGstBaseSrcPrepareSeekSegment; }
void setGstBaseSrcDoSeek               (GstBaseSrcClass * klass)  { klass->do_seek = goGstBaseSrcDoSeek; }
void setGstBaseSrcUnlock               (GstBaseSrcClass * klass)  { klass->unlock = goGstBaseSrcUnlock; }
void setGstBaseSrcUnlockStop           (GstBaseSrcClass * klass)  { klass->unlock_stop = goGstBaseSrcUnlockStop; }
void setGstBaseSrcQuery                (GstBaseSrcClass * klass)  { klass->query = goGstBaseSrcQuery; }
void setGstBaseSrcEvent                (GstBaseSrcClass * klass)  { klass->event = goGstBaseSrcEvent; }
void setGstBaseSrcCreate               (GstBaseSrcClass * klass)  { klass->create = goGstBaseSrcCreate; }
void setGstBaseSrcAlloc                (GstBaseSrcClass * klass)  { klass->alloc = goGstBaseSrcAlloc; }
void setGstBaseSrcFill                 (GstBaseSrcClass * klass)  { klass->fill = goGstBaseSrcFill; }

*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

var (
	// ExtendsBaseSrc is an Extendable for extending a GstBaseSrc
	ExtendsBaseSrc glib.Extendable = &extendsBaseSrc{parent: gst.ExtendsElement}
)

// GstBaseSrcImpl is the documented interface for an element extending a GstBaseSrc. It does not have to
// be implemented in it's entirety. Each of the methods it declares will be checked for their presence
// in the initializing object, and if the object declares an override it will replace the default
// implementation in the virtual methods.
type GstBaseSrcImpl interface {
	// GetCaps retrieves the caps for this class.
	GetCaps(*GstBaseSrc, *gst.Caps) *gst.Caps
	// Negotiate decides on the caps for this source.
	Negotiate(*GstBaseSrc) bool
	// Fixate is called if, during negotiation, caps need fixating.
	Fixate(*GstBaseSrc, *gst.Caps) *gst.Caps
	// SetCaps is used to notify this class of new caps.
	SetCaps(*GstBaseSrc, *gst.Caps) bool
	// DecideAllocation sets up an allocation query.
	DecideAllocation(*GstBaseSrc, *gst.Query) bool
	// Start the source, ideal for opening resources.
	Start(*GstBaseSrc) bool
	// Stop the source, ideal for closing resources.
	Stop(*GstBaseSrc) bool
	// GetTimes should, given a buffer, return start and stop time when it should be pushed.
	// The base class will sync on the clock using these times.
	GetTimes(*GstBaseSrc, *gst.Buffer) (start, end time.Duration)
	// GetSize should get the total size of the resource in bytes.
	GetSize(*GstBaseSrc) (bool, int64)
	// IsSeekable should check if the resource is seekable.
	IsSeekable(*GstBaseSrc) bool
	// PrepareSeekSegment prepares the segment on which to perform DoSeek, converting to the
	// current basesrc format.
	PrepareSeekSegment(*GstBaseSrc, *gst.Event, *gst.Segment) bool
	// DoSeek is used to notify subclasses of a seek.
	DoSeek(*GstBaseSrc, *gst.Segment) bool
	// Unlock should unlock any pending access to the resource. Subclasses should perform the unlock
	// ASAP.
	Unlock(*GstBaseSrc) bool
	// UnlockStop should clear any pending unlock request, as we succeeded in unlocking.
	UnlockStop(*GstBaseSrc) bool
	// Query is used to notify subclasses of a query.
	Query(*GstBaseSrc, *gst.Query) bool
	// Event is used to notify subclasses of an event.
	Event(*GstBaseSrc, *gst.Event) bool
	// Create asks the subclass to create a buffer with offset and size. The default implementation
	// will call alloc and fill.
	Create(self *GstBaseSrc, offset uint64, size uint) (gst.FlowReturn, *gst.Buffer)
	// Alloc asks the subclass to allocate an output buffer. The default implementation will use the negotiated
	// allocator.
	Alloc(self *GstBaseSrc, offset uint64, size uint) (gst.FlowReturn, *gst.Buffer)
	// Fill asks the subclass to fill the buffer with data from offset and size.
	Fill(self *GstBaseSrc, offset uint64, size uint, buffer *gst.Buffer) gst.FlowReturn
}

type extendsBaseSrc struct{ parent glib.Extendable }

func (e *extendsBaseSrc) Type() glib.Type     { return glib.Type(C.gst_base_src_get_type()) }
func (e *extendsBaseSrc) ClassSize() int64    { return int64(C.sizeof_GstBaseSrcClass) }
func (e *extendsBaseSrc) InstanceSize() int64 { return int64(C.sizeof_GstBaseSrc) }

// InitClass iterates the methods provided by the element and overrides any provided
// in the virtual methods.
func (e *extendsBaseSrc) InitClass(klass unsafe.Pointer, elem glib.GoObjectSubclass) {
	e.parent.InitClass(klass, elem)

	class := C.toGstBaseSrcClass(klass)

	if _, ok := elem.(interface {
		GetCaps(*GstBaseSrc, *gst.Caps) *gst.Caps
	}); ok {
		C.setGstBaseSrcGetCaps(class)
	}

	if _, ok := elem.(interface {
		Negotiate(*GstBaseSrc) bool
	}); ok {
		C.setGstBaseSrcNegotiate(class)
	}

	if _, ok := elem.(interface {
		Fixate(*GstBaseSrc, *gst.Caps) *gst.Caps
	}); ok {
		C.setGstBaseSrcFixate(class)
	}

	if _, ok := elem.(interface {
		SetCaps(*GstBaseSrc, *gst.Caps) bool
	}); ok {
		C.setGstBaseSrcSetCaps(class)
	}

	if _, ok := elem.(interface {
		DecideAllocation(*GstBaseSrc, *gst.Query) bool
	}); ok {
		C.setGstBaseSrcDecideAllocation(class)
	}

	if _, ok := elem.(interface {
		Start(*GstBaseSrc) bool
	}); ok {
		C.setGstBaseSrcStart(class)
	}

	if _, ok := elem.(interface {
		Stop(*GstBaseSrc) bool
	}); ok {
		C.setGstBaseSrcStop(class)
	}

	if _, ok := elem.(interface {
		GetTimes(*GstBaseSrc, *gst.Buffer) (start, end time.Duration)
	}); ok {
		C.setGstBaseSrcGetTimes(class)
	}

	if _, ok := elem.(interface {
		GetSize(*GstBaseSrc) (bool, int64)
	}); ok {
		C.setGstBaseSrcGetSize(class)
	}

	if _, ok := elem.(interface {
		IsSeekable(*GstBaseSrc) bool
	}); ok {
		C.setGstBaseSrcIsSeekable(class)
	}

	if _, ok := elem.(interface {
		PrepareSeekSegment(*GstBaseSrc, *gst.Event, *gst.Segment) bool
	}); ok {
		C.setGstBaseSrcPrepareSeekSegment(class)
	}

	if _, ok := elem.(interface {
		DoSeek(*GstBaseSrc, *gst.Segment) bool
	}); ok {
		C.setGstBaseSrcDoSeek(class)
	}

	if _, ok := elem.(interface {
		Unlock(*GstBaseSrc) bool
	}); ok {
		C.setGstBaseSrcUnlock(class)
	}

	if _, ok := elem.(interface {
		UnlockStop(*GstBaseSrc) bool
	}); ok {
		C.setGstBaseSrcUnlockStop(class)
	}

	if _, ok := elem.(interface {
		Query(*GstBaseSrc, *gst.Query) bool
	}); ok {
		C.setGstBaseSrcQuery(class)
	}

	if _, ok := elem.(interface {
		Event(*GstBaseSrc, *gst.Event) bool
	}); ok {
		C.setGstBaseSrcEvent(class)
	}

	if _, ok := elem.(interface {
		Create(self *GstBaseSrc, offset uint64, size uint) (gst.FlowReturn, *gst.Buffer)
	}); ok {
		C.setGstBaseSrcCreate(class)
	}

	if _, ok := elem.(interface {
		Alloc(self *GstBaseSrc, offset uint64, size uint) (gst.FlowReturn, *gst.Buffer)
	}); ok {
		C.setGstBaseSrcAlloc(class)
	}

	if _, ok := elem.(interface {
		Fill(self *GstBaseSrc, offset uint64, size uint, buffer *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstBaseSrcFill(class)
	}
}
