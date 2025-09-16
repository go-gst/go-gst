package gst

import (
	"context"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/core/userdata"
	"github.com/diamondburned/gotk4/pkg/gobject/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
// extern void _gotk4_gst1_PromiseChangeFunc(GstPromise*, gpointer);
// extern void destroyUserdata(gpointer);
import "C"

var TypePromise = gobject.Type(C.gst_promise_get_type())

func init() {
	gobject.RegisterGValueMarshalers([]gobject.TypeMarshaler{
		gobject.TypeMarshaler{T: TypePromise, F: marshalPromise},
	})
}

// Promise wraps GstPromise
//
// The #GstPromise object implements the container for values that may
// be available later. i.e. a Future or a Promise in
// &lt;https://en.wikipedia.org/wiki/Futures_and_promises&gt;.
// As with all Future/Promise-like functionality, there is the concept of the
// producer of the value and the consumer of the value.
//
// A #GstPromise is created with gst_promise_new() by the consumer and passed
// to the producer to avoid thread safety issues with the change callback.
// A #GstPromise can be replied to with a value (or an error) by the producer
// with gst_promise_reply(). The exact value returned is defined by the API
// contract of the producer and %NULL may be a valid reply.
// gst_promise_interrupt() is for the consumer to
// indicate to the producer that the value is not needed anymore and producing
// that value can stop.  The @GST_PROMISE_RESULT_EXPIRED state set by a call
// to gst_promise_expire() indicates to the consumer that a value will never
// be produced and is intended to be called by a third party that implements
// some notion of message handling such as #GstBus.
// A callback can also be installed at #GstPromise creation for
// result changes with gst_promise_new_with_change_func().
// The change callback can be used to chain #GstPromises's together as in the
// following example.
// |[&lt;!-- language="C" --&gt;
// const GstStructure *reply;
// GstPromise *p;
// if (gst_promise_wait (promise) != GST_PROMISE_RESULT_REPLIED)
//
//	return; // interrupted or expired value
//
// reply = gst_promise_get_reply (promise);
// if (error in reply)
//
//	return; // propagate error
//
// p = gst_promise_new_with_change_func (another_promise_change_func, user_data, notify);
// pass p to promise-using API
// ]|
//
// Each #GstPromise starts out with a #GstPromiseResult of
// %GST_PROMISE_RESULT_PENDING and only ever transitions once
// into one of the other #GstPromiseResult's.
//
// In order to support multi-threaded code, gst_promise_reply(),
// gst_promise_interrupt() and gst_promise_expire() may all be from
// different threads with some restrictions and the final result of the promise
// is whichever call is made first.  There are two restrictions on ordering:
//
// 1. That gst_promise_reply() and gst_promise_interrupt() cannot be called
// after gst_promise_expire()
// 2. That gst_promise_reply() and gst_promise_interrupt()
// cannot be called twice.
//
// The change function set with gst_promise_new_with_change_func() is
// called directly from either the gst_promise_reply(),
// gst_promise_interrupt() or gst_promise_expire() and can be called
// from an arbitrary thread.  #GstPromise using APIs can restrict this to
// a single thread or a subset of threads but that is entirely up to the API
// that uses #GstPromise.
type Promise struct {
	*promise

	// done is a close only channel to signal that the promise is done
	done chan struct{}
}

// promise is the struct that's finalized
type promise struct {
	native *C.GstPromise
}

var _ gobject.GoValueInitializer = (*Promise)(nil)

func marshalPromise(p unsafe.Pointer) (interface{}, error) {
	b := gobject.ValueFromNative(p).Boxed()
	return UnsafePromiseFromGlibBorrow(b), nil
}

func (r *Promise) GoValueType() gobject.Type {
	return TypePromise
}

func (r *Promise) SetGoValue(v *gobject.Value) {
	v.SetBoxed(unsafe.Pointer(r.instance()))
}

func (r *Promise) instance() *C.GstPromise {
	if r == nil {
		return nil
	}
	return r.native
}

// UnsafePromiseFromGlibBorrow is used to convert raw C.GstPromise pointers to go. This is used by the bindings internally.
func UnsafePromiseFromGlibBorrow(p unsafe.Pointer) *Promise {
	if p == nil {
		return nil
	}
	return &Promise{
		promise: &promise{(*C.GstPromise)(p)},
		done:    nil, // this will stay nil if the promise is received from C code
	}
}

// UnsafePromiseFromGlibNone is used to convert raw C.GstPromise pointers to go without transferring ownership. This is used by the bindings internally.
func UnsafePromiseFromGlibNone(p unsafe.Pointer) *Promise {
	if p == nil {
		return nil
	}
	C.gst_promise_ref((*C.GstPromise)(p))
	wrapped := UnsafePromiseFromGlibBorrow(p)
	runtime.SetFinalizer(
		wrapped.promise,
		func(intern *promise) {
			C.gst_promise_unref(intern.native)
		},
	)
	return wrapped
}

// UnsafePromiseFromGlibFull is used to convert raw C.GstPromise pointers to go while taking ownership. This is used by the bindings internally.
func UnsafePromiseFromGlibFull(p unsafe.Pointer) *Promise {
	if p == nil {
		return nil
	}
	wrapped := UnsafePromiseFromGlibBorrow(p)
	runtime.SetFinalizer(
		wrapped.promise,
		func(intern *promise) {
			C.gst_promise_unref(intern.native)
		},
	)
	return wrapped
}

// UnsafePromiseToGlibNone returns the underlying C pointer. This is used by the bindings internally.
func UnsafePromiseToGlibNone(p *Promise) unsafe.Pointer {
	if p == nil {
		return nil
	}
	return unsafe.Pointer(p.native)
}

// UnsafePromiseToGlibFull returns the underlying C pointer and gives up ownership.
// This is used by the bindings internally.
func UnsafePromiseToGlibFull(p *Promise) unsafe.Pointer {
	if p == nil {
		return nil
	}
	runtime.SetFinalizer(p.promise, nil)
	_p := unsafe.Pointer(p.native)
	p.native = nil // Promise is invalid from here on
	return _p
}

// NewPromise wraps gst_promise_new_with_change_func / gst_promise_new and allows the await calls to be more go like.
//
// The function returns the following values:
//
//   - goret *Promise
//
// @func will be called exactly once when transitioning out of
// %GST_PROMISE_RESULT_PENDING into any of the other #GstPromiseResult
// states.
func NewPromise() *Promise {
	done := make(chan struct{})

	changefunc := func(p *Promise) {
		close(p.done)
	}

	var carg1 C.GstPromiseChangeFunc = (*[0]byte)(C._gotk4_gst1_PromiseChangeFunc)
	var carg2 C.gpointer = C.gpointer(userdata.Register(changefunc))
	var carg3 C.GDestroyNotify = (C.GDestroyNotify)((*[0]byte)(C.destroyUserdata))

	cret := C.gst_promise_new_with_change_func(carg1, carg2, carg3)

	var goret *Promise

	goret = UnsafePromiseFromGlibFull(unsafe.Pointer(cret))

	goret.done = done

	return goret
}

// UnsafeRef increases the ref count on the promise. Exposed for completeness sake. Should not be called
// by application code.
func (promise *Promise) UnsafeRef() {
	var carg0 *C.GstPromise // in, none, converted

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	C.gst_promise_ref(carg0)
	runtime.KeepAlive(promise)
}

// UnsafeUnref decreases the ref count on the promise. Exposed for completeness sake. Should not be called
// by application code.
func (promise *Promise) UnsafeUnref() {
	var carg0 *C.GstPromise // in, none, converted

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	C.gst_promise_unref(carg0)
	runtime.KeepAlive(promise)
}

// Expire wraps gst_promise_expire
//
// Expire a @promise.  This will wake up any waiters with
// %GST_PROMISE_RESULT_EXPIRED.  Called by a message loop when the parent
// message is handled and/or destroyed (possibly unanswered).
func (promise *Promise) Expire() {
	var carg0 *C.GstPromise // in, none, converted

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	C.gst_promise_expire(carg0)
	runtime.KeepAlive(promise)
}

// GetReply wraps gst_promise_get_reply
//
// The function returns the following values:
//
//   - goret *Structure (nullable)
//
// Retrieve the reply set on @promise.  @promise must be in
// %GST_PROMISE_RESULT_REPLIED and the returned structure is owned by @promise
func (promise *Promise) GetReply() *Structure {
	var carg0 *C.GstPromise  // in, none, converted
	var cret *C.GstStructure // return, none, converted, nullable

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	cret = C.gst_promise_get_reply(carg0)
	runtime.KeepAlive(promise)

	var goret *Structure

	if cret != nil {
		goret = UnsafeStructureFromGlibNone(unsafe.Pointer(cret))
	}

	return goret
}

// Interrupt wraps gst_promise_interrupt
//
// Interrupt waiting for a @promise.  This will wake up any waiters with
// %GST_PROMISE_RESULT_INTERRUPTED.  Called when the consumer does not want
// the value produced anymore.
func (promise *Promise) Interrupt() {
	var carg0 *C.GstPromise // in, none, converted

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	C.gst_promise_interrupt(carg0)
	runtime.KeepAlive(promise)
}

// Reply wraps gst_promise_reply
//
// The function takes the following parameters:
//
//   - s *Structure (nullable): a #GstStructure with the the reply contents
//
// Set a reply on @promise.  This will wake up any waiters with
// %GST_PROMISE_RESULT_REPLIED.  Called by the producer of the value to
// indicate success (or failure).
//
// If @promise has already been interrupted by the consumer, then this reply
// is not visible to the consumer.
func (promise *Promise) Reply(s *Structure) {
	var carg0 *C.GstPromise   // in, none, converted
	var carg1 *C.GstStructure // in, full, converted, nullable

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))
	if s != nil {
		carg1 = (*C.GstStructure)(UnsafeStructureToGlibFull(s))
	}

	C.gst_promise_reply(carg0, carg1)
	runtime.KeepAlive(promise)
	runtime.KeepAlive(s)
}

// WaitBlocking wraps gst_promise_wait. Prefer to use [Promise.Await] over this.
//
// The function returns the following values:
//
//   - goret PromiseResult
//
// Wait for @promise to move out of the %GST_PROMISE_RESULT_PENDING state.
// If @promise is not in %GST_PROMISE_RESULT_PENDING then it will return
// immediately with the current result.
func (promise *Promise) WaitBlocking() PromiseResult {
	var carg0 *C.GstPromise     // in, none, converted
	var cret C.GstPromiseResult // return, none, casted

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	cret = C.gst_promise_wait(carg0)
	runtime.KeepAlive(promise)

	var goret PromiseResult

	goret = PromiseResult(cret)

	return goret
}

// Await awaits the promise without blocking the thread. It returns the reply returned by the GstPromise
//
// its implementation is preferred over the blocking [Promise.WaitBlocking], which would lock a thread until the
// promise has changed state.
func (p *Promise) Await(ctx context.Context) (*Structure, error) {
	if p.done == nil {
		// this can happen if the promise was received from a C function
		panic("cannot await a promise that has no done channel, this is likely a misuse of the promise")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.done:
	}

	// WaitBlocking will not block here, because the promise has already changed state
	result := p.WaitBlocking()

	if result != PromiseResultReplied {
		return nil, fmt.Errorf("promise did not reply: got %s", result)
	}

	structure := p.GetReply()

	return structure, nil
}
