package gst

/*
#include "gst.go.h"

extern void goPromiseChangeFunc (GstPromise*, gpointer user_data);
extern void cgoUnrefGopointerUserData (gpointer);

void cgoPromiseChangeFunc (GstPromise *promise, gpointer data)
{
	goPromiseChangeFunc(promise, data);
}
*/
import "C"
import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	gopointer "github.com/mattn/go-pointer"
)

type PromiseResult int

func (pr PromiseResult) String() string {
	switch pr {
	case PromiseResultPending:
		return "PENDING"
	case PromiseResultInterrupted:
		return "INTERRUPTED"
	case PromiseResultReplied:
		return "REPLIED"
	case PromiseResultExpired:
		return "EXPIRED"
	default:
		return "UNKNOWN"
	}
}

const (
	//Initial state. Waiting for transition to any other state.
	PromiseResultPending = C.GST_PROMISE_RESULT_PENDING
	// Interrupted by the consumer as it doesn't want the value anymore.
	PromiseResultInterrupted = C.GST_PROMISE_RESULT_INTERRUPTED
	// A producer marked a reply
	PromiseResultReplied = C.GST_PROMISE_RESULT_REPLIED
	// The promise expired (the carrying object lost all refs) and the promise will never be fulfilled.
	PromiseResultExpired = C.GST_PROMISE_RESULT_EXPIRED
)

// Promise is a go wrapper around a GstPromise.
// See: https://gstreamer.freedesktop.org/documentation/gstreamer/gstpromise.html
//
// it can be awaited on-blocking using Await, given the promise was constructed in go and not received from FFI.
type Promise struct {
	ptr *C.GstPromise

	// done will be closed when the GstPromise has changed state
	done <-chan struct{}
}

func NewPromise() *Promise {
	done := make(chan struct{})

	fPtr := gopointer.Save(func() {
		close(done)
	})

	cprom := C.gst_promise_new_with_change_func(
		C.GstPromiseChangeFunc(C.cgoPromiseChangeFunc),
		C.gpointer(fPtr),
		C.GDestroyNotify(C.cgoUnrefGopointerUserData),
	)

	prom := &Promise{
		ptr:  cprom,
		done: done,
	}

	runtime.SetFinalizer(prom, func(prom *Promise) {
		prom.Unref()
	})

	return prom
}

func (p *Promise) Instance() *C.GstPromise {
	return p.ptr
}

// Ref increases the ref count on the promise. Exposed for completeness sake. Should not be called
// by application code
func (p *Promise) Ref() {
	C.gst_promise_ref(p.ptr)
}

// Unref decreases the ref count on the promise. Exposed for completeness sake. Should not be called
// by application code
func (p *Promise) Unref() {
	C.gst_promise_unref(p.ptr)
}

var ErrPromiseNotReplied = errors.New("promise was not replied")
var ErrNilPromiseReply = errors.New("promise returned a nil reply")

// ErrCannotAwaitPromise signifies that we do not have a channel that we can await.
//
// this happens if the promise was marshaled from a GValue coming from C
var ErrCannotAwaitPromise = errors.New("promises received from FFI cannot be awaited")

// Await awaits the promise without blocking the thread. It returns the reply returned by the GstPromise
//
// its implementation is preferred over the blocking gst_promise_wait, which would lock a thread until the
// promise has changed state.
func (p *Promise) Await(ctx context.Context) (*Structure, error) {
	if p.done == nil {
		return nil, ErrCannotAwaitPromise
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.done:
	}

	// gst_promise_wait will not block here, because the promise has already changed state
	result := PromiseResult(C.gst_promise_wait(p.ptr))

	if result != PromiseResultReplied {
		return nil, fmt.Errorf("%w: got %s", ErrPromiseNotReplied, result)
	}

	structure := p.GetReply()

	if structure == nil {
		return nil, ErrNilPromiseReply
	}

	return structure, nil
}

// GetReply wraps gst_promise_get_reply and returns the structure, which can be nil.
func (p *Promise) GetReply() *Structure {
	cstruct := C.gst_promise_get_reply(p.ptr)

	if cstruct == nil {
		return nil
	}

	structure := wrapStructure(cstruct)

	// the structure is owned by the promise, so we keep the promise alive
	// until the structure gets GC'ed
	p.Ref()
	runtime.SetFinalizer(structure, func(_ *Structure) {
		p.Unref()
	})

	return structure
}

// Expire wraps gst_promise_expire
func (p *Promise) Expire() {
	C.gst_promise_expire(p.ptr)
}

// Interrupt wraps gst_promise_interrupt
func (p *Promise) Interrupt() {
	C.gst_promise_interrupt(p.ptr)
}

// Reply wraps gst_promise_reply
func (p *Promise) Reply(answer *Structure) {
	C.gst_promise_reply(p.ptr, answer.Instance())
}

var TypePromise = glib.Type(C.GST_TYPE_PROMISE)

// ToGValue implements glib.ValueTransformer
func (p *Promise) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypePromise)
	if err != nil {
		return nil, err
	}
	val.SetInstance(unsafe.Pointer(p.Instance()))
	return val, nil
}

func marshalPromise(p unsafe.Pointer) (interface{}, error) {
	c := C.g_value_get_object(toGValue(p))
	obj := (*C.GstPromise)(unsafe.Pointer(c))

	prom := &Promise{
		ptr:  obj,
		done: nil, // cannot be awaited if received from FFI
	}

	prom.Ref()

	runtime.SetFinalizer(prom, func(p *Promise) {
		p.Unref()
	})

	return prom, nil
}
