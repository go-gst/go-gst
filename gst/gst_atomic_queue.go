package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

// AtomicQueue wraps a GstAtomicQueue that can be used from multiple threads
// without performing any blocking operations.
type AtomicQueue struct {
	ptr *C.GstAtomicQueue
}

// NewAtomicQueue creates a new atomic queue with the given size. The size will
// be rounded up to the nearest power of 2 and used as the initial size of the queue.
func NewAtomicQueue(size int) *AtomicQueue {
	return wrapAtomicQueue(C.gst_atomic_queue_new(C.guint(size)))
}

// Instance returns the underlying queue instance.
func (a *AtomicQueue) Instance() *C.GstAtomicQueue { return a.ptr }

// Length returns the amount of items in this queue.
func (a *AtomicQueue) Length() int {
	return int(C.gst_atomic_queue_length(a.Instance()))
}

// Peek looks at the first item in the queue without removing it. This function
// returns nil if the queue is empty.
func (a *AtomicQueue) Peek() interface{} {
	ptr := C.gst_atomic_queue_peek(a.Instance())
	if ptr == nil {
		return nil
	}
	return gopointer.Restore(unsafe.Pointer(ptr))
}

// Pop pops the head element off the queue. This function returns nil if the queue
// is empty.
func (a *AtomicQueue) Pop() interface{} {
	ptr := C.gst_atomic_queue_pop(a.Instance())
	if ptr == nil {
		return nil
	}
	defer gopointer.Unref(unsafe.Pointer(ptr))
	return gopointer.Restore(unsafe.Pointer(ptr))
}

// Push appends the given data to the end of the queue.
func (a *AtomicQueue) Push(data interface{}) {
	ptr := gopointer.Save(data)
	C.gst_atomic_queue_push(a.Instance(), (C.gpointer)(unsafe.Pointer(ptr)))
}

// Ref increases the ref count on the queue by one.
func (a *AtomicQueue) Ref() {
	C.gst_atomic_queue_ref(a.Instance())
}

// Unref decreaes the ref count on the queue by one. Memory is freed when the
// refcount reaches zero.
func (a *AtomicQueue) Unref() { C.gst_atomic_queue_unref(a.Instance()) }
