package gst

// #include "gst.go.h"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Bin is a go wrapper arounds a GstBin.
type Bin struct{ *Element }

// NewBin returns a new Bin with the given name.
func NewBin(name string) *Bin {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	bin := C.gst_bin_new((*C.gchar)(unsafe.Pointer(cName)))
	return wrapBin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(bin))})
}

// BinFromElement wraps the given Element in a Bin reference. This only works for elements
// that implement their own Bin, such as playbin. If the provided element does not implement
// a Bin then nil is returned.
func BinFromElement(elem *Element) *Bin {
	if C.toGstBin(elem.Unsafe()) == nil {
		return nil
	}
	return &Bin{elem}
}

// Instance returns the underlying GstBin instance.
func (b *Bin) Instance() *C.GstBin { return C.toGstBin(b.Unsafe()) }

// GetElementByName returns the element with the given name. Unref after usage.
func (b *Bin) GetElementByName(name string) (*Element, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	elem := C.gst_bin_get_by_name((*C.GstBin)(b.Instance()), (*C.gchar)(cName))
	if elem == nil {
		return nil, fmt.Errorf("Could not find element with name %s", name)
	}
	return wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))}), nil
}

// GetElementByNameRecursive returns the element with the given name. If it is not
// found in this Bin, parent Bins are searched recursively. Unref after usage.
func (b *Bin) GetElementByNameRecursive(name string) (*Element, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	elem := C.gst_bin_get_by_name_recurse_up((*C.GstBin)(b.Instance()), (*C.gchar)(cName))
	if elem == nil {
		return nil, fmt.Errorf("Could not find element with name %s", name)
	}
	return wrapElement(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(elem))}), nil
}

// GetElements returns a list of the elements added to this pipeline.
func (b *Bin) GetElements() ([]*Element, error) {
	iterator := C.gst_bin_iterate_elements((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// GetElementsRecursive returns a list of the elements added to this Bin. It recurses
// children Bins.
func (b *Bin) GetElementsRecursive() ([]*Element, error) {
	iterator := C.gst_bin_iterate_recurse((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// GetSourceElements returns a list of all the source elements in this Bin.
func (b *Bin) GetSourceElements() ([]*Element, error) {
	iterator := C.gst_bin_iterate_sources((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// GetSinkElements returns a list of all the sink elements in this Bin. Unref
// elements after usage.
func (b *Bin) GetSinkElements() ([]*Element, error) {
	iterator := C.gst_bin_iterate_sinks((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// GetElementsSorted returns a list of the elements in this bin in topologically sorted order.
// This means that the elements are returned from the most downstream elements (sinks) to the sources.
func (b *Bin) GetElementsSorted() ([]*Element, error) {
	iterator := C.gst_bin_iterate_sorted((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// GetElementsByFactoryName returns a list of the elements in this bin created from the given factory
// name.
func (b *Bin) GetElementsByFactoryName(name string) ([]*Element, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	iterator := C.gst_bin_iterate_all_by_element_factory_name(b.Instance(), (*C.gchar)(unsafe.Pointer(cname)))
	return iteratorToElementSlice(iterator)
}

// Add adds an element to the bin.
func (b *Bin) Add(elem *Element) error {
	if ok := C.gst_bin_add((*C.GstBin)(b.Instance()), (*C.GstElement)(elem.Instance())); !gobool(ok) {
		return fmt.Errorf("Failed to add element to pipeline: %s", elem.Name())
	}
	return nil
}

// AddMany is a go implementation of `gst_bin_add_many`.
func (b *Bin) AddMany(elems ...*Element) error {
	for _, elem := range elems {
		if err := b.Add(elem); err != nil {
			return err
		}
	}
	return nil
}

// Remove removes an element from the Bin.
func (b *Bin) Remove(elem *Element) error {
	if ok := C.gst_bin_remove((*C.GstBin)(b.Instance()), (*C.GstElement)(elem.Instance())); !gobool(ok) {
		return fmt.Errorf("Failed to add element to pipeline: %s", elem.Name())
	}
	return nil
}

// RemoveMany is a go implementation of `gst_bin_remove_many`.
func (b *Bin) RemoveMany(elems ...*Element) error {
	for _, elem := range elems {
		if err := b.Remove(elem); err != nil {
			return err
		}
	}
	return nil
}

// FindUnlinkedPad recursively looks for elements with an unlinked pad of the given direction
// within this bin and returns an unlinked pad if one is found, or NULL otherwise. If a pad is
// found, the caller owns a reference to it and should unref it when it is not needed any longer.
func (b *Bin) FindUnlinkedPad(direction PadDirection) *Pad {
	pad := C.gst_bin_find_unlinked_pad(b.Instance(), C.GstPadDirection(direction))
	if pad == nil {
		return nil
	}
	return wrapPad(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(pad))})
}

// GetSuppressedFlags returns the suppressed flags of the bin.
func (b *Bin) GetSuppressedFlags() ElementFlags {
	return ElementFlags(C.gst_bin_get_suppressed_flags(b.Instance()))
}

// SetSuppressedFlags suppresses the given flags on the bin. ElementFlags of a child element are
// propagated when it is added to the bin. When suppressed flags are set, those specified flags
// will not be propagated to the bin.
func (b *Bin) SetSuppressedFlags(flags ElementFlags) {
	C.gst_bin_set_suppressed_flags(b.Instance(), C.GstElementFlags(flags))
}

// RecalculateLatency queries bin for the current latency and reconfigures this latency to all
// the elements with a LATENCY event.
//
// This method is typically called on the pipeline when a MessageLatency is posted on the bus.
//
// This function simply emits the 'do-latency' signal so any custom latency calculations will be performed.
// It returns true if the latency could be queried and reconfigured.
func (b *Bin) RecalculateLatency() bool {
	return gobool(C.gst_bin_recalculate_latency(b.Instance()))
}

// SyncChildrenStates synchronizes the state of every child with the state of this Bin.
// This function returns true if the operation was successful.
func (b *Bin) SyncChildrenStates() bool {
	return gobool(C.gst_bin_sync_children_states(b.Instance()))
}

// DEBUG OPERATIONS //

// DebugGraphDetails casts GstDebugGraphDetails
type DebugGraphDetails int

// Type castings
const (
	DebugGraphShowMediaType        DebugGraphDetails = C.GST_DEBUG_GRAPH_SHOW_MEDIA_TYPE         // (1) – show caps-name on edges
	DebugGraphShowCapsDetails      DebugGraphDetails = C.GST_DEBUG_GRAPH_SHOW_CAPS_DETAILS       // (2) – show caps-details on edges
	DebugGraphShowNonDefaultParams DebugGraphDetails = C.GST_DEBUG_GRAPH_SHOW_NON_DEFAULT_PARAMS // (4) – show modified parameters on elements
	DebugGraphShowStates           DebugGraphDetails = C.GST_DEBUG_GRAPH_SHOW_STATES             // (8) – show element states
	DebugGraphShowPullParams       DebugGraphDetails = C.GST_DEBUG_GRAPH_SHOW_FULL_PARAMS        // (16) – show full element parameter values even if they are very long
	DebugGraphShowAll              DebugGraphDetails = C.GST_DEBUG_GRAPH_SHOW_ALL                // (15) – show all the typical details that one might want
	DebugGraphShowVerbose          DebugGraphDetails = C.GST_DEBUG_GRAPH_SHOW_VERBOSE            // (4294967295) – show all details regardless of how large or verbose they make the resulting output
)

// DebugBinToDotData will obtain the whole network of gstreamer elements that form the pipeline into a dot file.
// This data can be processed with graphviz to get an image.
func (b *Bin) DebugBinToDotData(details DebugGraphDetails) string {
	ret := C.gst_debug_bin_to_dot_data(b.Instance(), C.GstDebugGraphDetails(details))
	defer C.g_free((C.gpointer)(unsafe.Pointer(ret)))
	return C.GoString(ret)
}

// DebugBinToDotFile is like DebugBinToDotData except it will write the dot data to the filename
// specified.
func (b *Bin) DebugBinToDotFile(details DebugGraphDetails, filename string) {
	cname := C.CString(filename)
	defer C.free(unsafe.Pointer(cname))
	C.gst_debug_bin_to_dot_file(b.Instance(), C.GstDebugGraphDetails(details), (*C.gchar)(unsafe.Pointer(cname)))
}

// DebugBinToDotFileWithTs is like DebugBinToDotFile except it will write the dot data to the filename
// specified, except it will append the current timestamp to the filename.
func (b *Bin) DebugBinToDotFileWithTs(details DebugGraphDetails, filename string) {
	cname := C.CString(filename)
	defer C.free(unsafe.Pointer(cname))
	C.gst_debug_bin_to_dot_file_with_ts(b.Instance(), C.GstDebugGraphDetails(details), (*C.gchar)(unsafe.Pointer(cname)))
}

func iteratorToElementSlice(iterator *C.GstIterator) ([]*Element, error) {
	elems := make([]*Element, 0)
	gval := new(C.GValue)

	for {
		switch C.gst_iterator_next((*C.GstIterator)(iterator), (*C.GValue)(unsafe.Pointer(gval))) {
		case C.GST_ITERATOR_DONE:
			C.gst_iterator_free((*C.GstIterator)(iterator))
			return elems, nil
		case C.GST_ITERATOR_RESYNC:
			C.gst_iterator_resync((*C.GstIterator)(iterator))
		case C.GST_ITERATOR_OK:
			cElemVoid := C.g_value_get_object((*C.GValue)(gval))
			cElem := (*C.GstElement)(cElemVoid)
			elems = append(elems, wrapElement(glib.Take(unsafe.Pointer(cElem))))
			C.g_value_reset((*C.GValue)(gval))
		default:
			return nil, errors.New("Element iterator failed")
		}
	}
}
