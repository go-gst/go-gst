package base

/*
#include "gst.go.h"

extern GstFlowReturn  goGstPushSrcAlloc   (GstPushSrc * src, GstBuffer ** buf);
extern GstFlowReturn  goGstPushSrcCreate  (GstPushSrc * src, GstBuffer ** buf);
extern GstFlowReturn  goGstPushSrcFill    (GstPushSrc * src, GstBuffer * buf);

void  setGstPushSrcAlloc   (GstPushSrcClass * klass) { klass->alloc = goGstPushSrcAlloc; }
void  setGstPushSrcCreate  (GstPushSrcClass * klass) { klass->create = goGstPushSrcCreate; }
void  setGstPushSrcFill    (GstPushSrcClass * klass) { klass->fill = goGstPushSrcFill; }

*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

var (
	// ExtendsPushSrc is an Extendable for extending a GstPushSrc
	ExtendsPushSrc glib.Extendable = &extendsPushSrc{parent: ExtendsBaseSrc}
)

// GstPushSrcImpl is the documented interface for an element extending a GstPushSrc. It does not have to
// be implemented in it's entirety. Each of the methods it declares will be checked for their presence
// in the initializing object, and if the object declares an override it will replace the default
// implementation in the virtual methods.
type GstPushSrcImpl interface {
	// Asks the subclass to allocate a buffer. The subclass decides which size this buffer should be.
	// The default implementation will create a new buffer from the negotiated allocator.
	Alloc(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	// Asks the subclass to create a buffer. The subclass decides which size this buffer should be. Other
	// then that, refer to GstBaseSrc.create for more details. If this method is not implemented, alloc
	// followed by fill will be called.
	Create(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	// Asks the subclass to fill the buffer with data.
	Fill(*GstPushSrc, *gst.Buffer) gst.FlowReturn
}

type extendsPushSrc struct{ parent glib.Extendable }

func (e *extendsPushSrc) Type() glib.Type     { return glib.Type(C.gst_push_src_get_type()) }
func (e *extendsPushSrc) ClassSize() int64    { return int64(C.sizeof_GstPushSrcClass) }
func (e *extendsPushSrc) InstanceSize() int64 { return int64(C.sizeof_GstPushSrc) }

func (e *extendsPushSrc) InitClass(klass unsafe.Pointer, elem glib.GoObjectSubclass) {
	e.parent.InitClass(klass, elem)

	srcClass := C.toGstPushSrcClass(klass)

	if _, ok := elem.(interface {
		Alloc(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	}); ok {
		C.setGstPushSrcAlloc(srcClass)
	}

	if _, ok := elem.(interface {
		Create(*GstPushSrc) (gst.FlowReturn, *gst.Buffer)
	}); ok {
		C.setGstPushSrcCreate(srcClass)
	}

	if _, ok := elem.(interface {
		Fill(*GstPushSrc, *gst.Buffer) gst.FlowReturn
	}); ok {
		C.setGstPushSrcFill(srcClass)
	}
}
