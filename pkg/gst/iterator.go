package gst

import (
	"iter"
	"runtime"
	"unsafe"

	"github.com/go-gst/go-glib/pkg/gobject/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// Next wraps gst_iterator_next
// The function returns the following values:
//
//   - elem any: pointer to hold next element
//   - goret IteratorResult
//
// Get the next item from the iterator in @elem.
//
// Only when this function returns %GST_ITERATOR_OK, @elem will contain a valid
// value. @elem must have been initialized to the type of the iterator or
// initialized to zeroes with g_value_unset(). The caller is responsible for
// unsetting or resetting @elem with g_value_unset() or g_value_reset()
// after usage.
//
// When this function returns %GST_ITERATOR_DONE, no more elements can be
// retrieved from @it.
//
// A return value of %GST_ITERATOR_RESYNC indicates that the element list was
// concurrently updated. The user of @it should call gst_iterator_resync() to
// get the newly updated list.
//
// A return value of %GST_ITERATOR_ERROR indicates an unrecoverable fatal error.
func (it *Iterator) Next() (any, IteratorResult) {
	var carg0 *C.GstIterator     // in, none, converted
	var carg1 C.GValue           // out, transfer: none, C Pointers: 0, Name: Value, caller-allocates
	var cret C.GstIteratorResult // return, none, casted

	carg0 = (*C.GstIterator)(UnsafeIteratorToGlibNone(it))

	cret = C.gst_iterator_next(carg0, &carg1)
	runtime.KeepAlive(it)

	var elem any
	var goret IteratorResult

	elem = gobject.ValueFromNative(unsafe.Pointer(&carg1)).GoValue()
	goret = IteratorResult(cret)

	return elem, goret
}

// Values allows you to access the values from the iterator in a go for loop via function iterators
func (it *Iterator) Values() iter.Seq[any] {
	return func(yield func(any) bool) {
		for {
			v, ret := it.Next()
			switch ret {
			case IteratorDone:
				return
			case IteratorResync:
				it.Resync()
			case IteratorOK:
				if !yield(v) {
					return
				}

			case IteratorError:
				panic("iterator values failed")
			default:
				panic("iterator values returned unknown state")
			}
		}
	}
}
