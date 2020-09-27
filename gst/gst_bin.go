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
	return wrapElement(glib.Take(unsafe.Pointer(elem))), nil
}

// GetElements returns a list of the elements added to this pipeline.
func (b *Bin) GetElements() ([]*Element, error) {
	iterator := C.gst_bin_iterate_elements((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// GetSourceElements returns a list of all the source elements in this pipeline.
func (b *Bin) GetSourceElements() ([]*Element, error) {
	iterator := C.gst_bin_iterate_sources((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// GetSinkElements returns a list of all the sink elements in this pipeline. Unref
// elements after usage.
func (b *Bin) GetSinkElements() ([]*Element, error) {
	iterator := C.gst_bin_iterate_sinks((*C.GstBin)(b.Instance()))
	return iteratorToElementSlice(iterator)
}

// Add wraps `gst_bin_add`.
func (b *Bin) Add(elem *Element) error {
	if ok := C.gst_bin_add((*C.GstBin)(b.Instance()), (*C.GstElement)(elem.Instance())); gobool(ok) {
		return fmt.Errorf("Failed to add element to pipeline: %s", elem.Name())
	}
	return nil
}

// AddMany is a go implementation of `gst_bin_add_many` to compensate for the inability
// to use variadic functions in cgo.
func (b *Bin) AddMany(elems ...*Element) error {
	for _, elem := range elems {
		if err := b.Add(elem); err != nil {
			return err
		}
	}
	return nil
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
