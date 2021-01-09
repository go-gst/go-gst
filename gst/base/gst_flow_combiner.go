package base

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// FlowCombiner is a helper structure for aggregating flow returns. This struct
// is not thread safe.
// For more information see https://gstreamer.freedesktop.org/documentation/base/gstflowcombiner.html?gi-language=c#GstFlowCombiner
type FlowCombiner struct{ ptr *C.GstFlowCombiner }

func wrapFlowCombiner(ptr *C.GstFlowCombiner) *FlowCombiner {
	return &FlowCombiner{ptr}
}

// NewFlowCombiner creates a new flow combiner. Use Free() to free it.
func NewFlowCombiner() *FlowCombiner {
	return wrapFlowCombiner(C.gst_flow_combiner_new())
}

// Instance returns the underlying GstFlowCombiner instance.
func (f *FlowCombiner) Instance() *C.GstFlowCombiner { return f.ptr }

// AddPad adds a new pad to the FlowCombiner. A reference is taken on the pad.
func (f *FlowCombiner) AddPad(pad *gst.Pad) {
	C.gst_flow_combiner_add_pad(f.Instance(), (*C.GstPad)(unsafe.Pointer(pad.Instance())))
}

// Clear will remove all pads and reset the combiner to its initial state.
func (f *FlowCombiner) Clear() { C.gst_flow_combiner_clear(f.Instance()) }

// Free will free a FlowCombiner and all its internal data.
func (f *FlowCombiner) Free() { C.gst_flow_combiner_free(f.Instance()) }

// Ref will increment the reference count on the FlowCombiner.
func (f *FlowCombiner) Ref() *FlowCombiner {
	return wrapFlowCombiner(C.gst_flow_combiner_ref(f.Instance()))
}

// RemovePad will remove a pad from the FlowCombiner.
func (f *FlowCombiner) RemovePad(pad *gst.Pad) {
	C.gst_flow_combiner_remove_pad(f.Instance(), (*C.GstPad)(unsafe.Pointer(pad.Instance())))
}

// Reset flow combiner and all pads to their initial state without removing pads.
func (f *FlowCombiner) Reset() { C.gst_flow_combiner_reset(f.Instance()) }

// Unref decrements the reference count on the Flow Combiner.
func (f *FlowCombiner) Unref() { C.gst_flow_combiner_unref(f.Instance()) }

// UpdateFlow computes the combined flow return for the pads in it.
//
// The GstFlowReturn parameter should be the last flow return update for a pad in this GstFlowCombiner.
// It will use this value to be able to shortcut some combinations and avoid looking over all pads again.
// e.g. The last combined return is the same as the latest obtained GstFlowReturn.
func (f *FlowCombiner) UpdateFlow(fret gst.FlowReturn) gst.FlowReturn {
	return gst.FlowReturn(C.gst_flow_combiner_update_flow(f.Instance(), C.GstFlowReturn(fret)))
}

// UpdatePadFlow sets the provided pad's last flow return to provided value and computes the combined flow
// return for the pads in it.
//
// The GstFlowReturn parameter should be the last flow return update for a pad in this GstFlowCombiner. It
// will use this value to be able to shortcut some combinations and avoid looking over all pads again. e.g.
// The last combined return is the same as the latest obtained GstFlowReturn.
func (f *FlowCombiner) UpdatePadFlow(pad *gst.Pad, fret gst.FlowReturn) gst.FlowReturn {
	return gst.FlowReturn(C.gst_flow_combiner_update_pad_flow(
		f.Instance(),
		(*C.GstPad)(unsafe.Pointer(pad.Instance())),
		C.GstFlowReturn(fret),
	))
}
