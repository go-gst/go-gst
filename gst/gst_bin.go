package gst

/*
#include "gst.go.h"

gboolean
binParentAddElement (GstBin * bin, GstElement * element)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(bin));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	return parent->add_element(bin, element);
}

void
binParentDeepElementAdded (GstBin * bin, GstBin * subbin, GstElement * element)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(bin));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	parent->deep_element_added(bin, subbin, element);
}

void
binParentDeepElementRemoved (GstBin * bin, GstBin * subbin, GstElement * element)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(element));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	parent->deep_element_removed(bin, subbin, element);
}

gboolean
binParentDoLatency (GstBin * bin)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(bin));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	return parent->do_latency(bin);
}

void
binParentElementAdded (GstBin * bin, GstElement * element)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(bin));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	parent->element_added(bin, element);
}

void
binParentElementRemoved (GstBin * bin, GstElement * element)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(bin));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	parent->element_removed(bin, element);
}

void
binParentHandleMessage (GstBin * bin, GstMessage * message)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(bin));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	parent->handle_message(bin, message);
}

gboolean
binParentRemoveElement (GstBin * bin, GstElement * element)
{
	GObjectClass * this_class = G_OBJECT_GET_CLASS(G_OBJECT(bin));
	GstBinClass * parent = toGstBinClass(g_type_class_peek_parent(this_class));
	return parent->remove_element(bin, element);
}

*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Bin is a go wrapper arounds a GstBin.
type Bin struct{ *Element }

// NewBin returns a new Bin with the given name.
func NewBin(name string) *Bin {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	bin := C.gst_bin_new((*C.gchar)(unsafe.Pointer(cName)))
	return wrapBin(glib.TransferNone(unsafe.Pointer(bin)))
}

// NewBinFromString constructs a bin from a string description.
// description - command line describing the bin
// ghostUnlinkedPads - whether to automatically create ghost pads for unlinked source or sink pads within the bin
func NewBinFromString(description string, ghostUnlinkedPads bool) (*Bin, error) {
	cDescription := C.CString(description)
	defer C.free(unsafe.Pointer(cDescription))
	var gerr *C.GError
	bin := C.gst_parse_bin_from_description((*C.gchar)(cDescription), gboolean(ghostUnlinkedPads), (**C.GError)(&gerr))
	if gerr != nil {
		defer C.g_error_free((*C.GError)(gerr))
		errMsg := C.GoString(gerr.message)
		return nil, errors.New(errMsg)
	}
	return &Bin{&Element{wrapObject(glib.TransferNone(unsafe.Pointer(bin)))}}, nil
}

// ToGstBin wraps the given glib.Object, gst.Object, or gst.Element in a Bin instance. Only
// works for objects that implement their own Bin.
func ToGstBin(obj interface{}) *Bin {
	switch obj := obj.(type) {
	case *Object:
		return &Bin{&Element{Object: obj}}
	case *Element:
		return &Bin{obj}
	case *glib.Object:
		return &Bin{&Element{Object: &Object{InitiallyUnowned: &glib.InitiallyUnowned{Object: obj}}}}
	}
	return nil
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
	return wrapElement(glib.TransferFull(unsafe.Pointer(elem))), nil
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
	return wrapElement(glib.TransferFull(unsafe.Pointer(elem))), nil
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

// GetByInterface looks for an element inside the bin that implements the given interface. If such an
// element is found, it returns the element. You can cast this element to the given interface afterwards.
// If you want all elements that implement the interface, use GetAllByInterface. This function recurses
// into child bins.
func (b *Bin) GetByInterface(iface glib.Interface) (*Element, error) {
	elem := C.gst_bin_get_by_interface(b.Instance(), C.GType(iface.Type()))
	if elem == nil {
		return nil, fmt.Errorf("Could not find any elements implementing %s", iface.Type().Name())
	}
	return wrapElement(glib.TransferFull(unsafe.Pointer(elem))), nil
}

// GetAllByInterface looks for all elements inside the bin that implements the given interface. You can
// safely cast all returned elements to the given interface. The function recurses inside child bins.
// The function will return a series of Elements that should be unreffed after use.
func (b *Bin) GetAllByInterface(iface glib.Interface) ([]*Element, error) {
	iterator := C.gst_bin_iterate_all_by_interface(b.Instance(), C.GType(iface.Type()))
	return iteratorToElementSlice(iterator)
}

// // GetElementsByFactoryName returns a list of the elements in this bin created from the given factory
// // name.
// func (b *Bin) GetElementsByFactoryName(name string) ([]*Element, error) {
// 	cname := C.CString(name)
// 	defer C.free(unsafe.Pointer(cname))
// 	iterator := C.gst_bin_iterate_all_by_element_factory_name(b.Instance(), (*C.gchar)(unsafe.Pointer(cname)))
// 	return iteratorToElementSlice(iterator)
// }

// Add adds an element to the bin.
func (b *Bin) Add(elem *Element) error {
	if ok := C.gst_bin_add((*C.GstBin)(b.Instance()), (*C.GstElement)(elem.Instance())); !gobool(ok) {
		return fmt.Errorf("Failed to add element to pipeline: %s", elem.GetName())
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
		return fmt.Errorf("Failed to remove element from pipeline: %s", elem.GetName())
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
	return wrapPad(glib.TransferFull(unsafe.Pointer(pad)))
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
			elems = append(elems, wrapElement(glib.TransferNone(unsafe.Pointer(cElem))))
			C.g_value_unset((*C.GValue)(gval))
		default:
			return nil, errors.New("Element iterator failed")
		}
	}
}

// ParentAddElement chains up to the parent AddElement handler.
func (b *Bin) ParentAddElement(element *Element) bool {
	return gobool(C.binParentAddElement(b.Instance(), element.Instance()))
}

// ParentDeepElementAdded chains up to the parent DeepElementAdded handler.
func (b *Bin) ParentDeepElementAdded(subbin *Bin, element *Element) {
	C.binParentDeepElementAdded(b.Instance(), subbin.Instance(), element.Instance())
}

// ParentDeepElementRemoved chains up to the parent DeepElementRemoved handler.
func (b *Bin) ParentDeepElementRemoved(subbin *Bin, element *Element) {
	C.binParentDeepElementRemoved(b.Instance(), subbin.Instance(), element.Instance())
}

// ParentDoLatency chains up to the parent DoLatency handler.
func (b *Bin) ParentDoLatency() bool {
	return gobool(C.binParentDoLatency(b.Instance()))
}

// ParentElementAdded chains up to the parent ElementAdded handler.
func (b *Bin) ParentElementAdded(element *Element) {
	C.binParentElementAdded(b.Instance(), element.Instance())
}

// ParentElementRemoved chains up to the parent ElementRemoved handler.
func (b *Bin) ParentElementRemoved(element *Element) {
	C.binParentElementRemoved(b.Instance(), element.Instance())
}

// ParentHandleMessage chains up to the parent HandleMessage handler.
func (b *Bin) ParentHandleMessage(message *Message) {
	C.binParentHandleMessage(b.Instance(), message.Instance())
}

// ParentRemoveElement chains up to the parent RemoveElement handler.
func (b *Bin) ParentRemoveElement(element *Element) bool {
	return gobool(C.binParentRemoveElement(b.Instance(), element.Instance()))
}
