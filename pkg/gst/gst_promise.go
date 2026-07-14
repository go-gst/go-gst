package gst

import (
	"context"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/userdata"
	"github.com/go-gst/go-glib/pkg/gobject/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
// extern void _goglib_gst1_PromiseChangeFunc(GstPromise*, gpointer);
// extern void destroyUserdata(gpointer);
import "C"

var TypePromise = gobject.Type(C.gst_promise_get_type())

func init() {
	gobject.RegisterGValueMarshalers([]gobject.TypeMarshaler{
		gobject.TypeMarshaler{T: TypePromise, F: marshalPromise},
	})
}

// Promise wraps GstPromise
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
func NewPromise() *Promise {
	done := make(chan struct{})

	var changefunc PromiseChangeFunc = func(_ *Promise) {
		// the promise passed to this function is transferred from C, so we close the done channel
		// directly here
		close(done)
	}

	var carg1 C.GstPromiseChangeFunc = (*[0]byte)(C._goglib_gst1_PromiseChangeFunc)
	var carg2 C.gpointer = C.gpointer(userdata.Register(changefunc))
	var carg3 C.GDestroyNotify = (C.GDestroyNotify)((*[0]byte)(C.destroyUserdata))

	cret := C.gst_promise_new_with_change_func(carg1, carg2, carg3)

	var goret *Promise

	goret = UnsafePromiseFromGlibFull(unsafe.Pointer(cret))

	// save the done channel so [Promise.Await] can use it
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
func (promise *Promise) Expire() {
	var carg0 *C.GstPromise // in, none, converted

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	C.gst_promise_expire(carg0)
	runtime.KeepAlive(promise)
}

// GetReply wraps gst_promise_get_reply
func (promise *Promise) GetReply() *Structure {
	var carg0 *C.GstPromise  // in, none, converted
	var cret *C.GstStructure // return, none, converted, nullable

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	cret = C.gst_promise_get_reply(carg0)
	runtime.KeepAlive(promise)

	var goret *Structure

	if cret != nil {
		goret = UnsafeStructureFromGlibNone(unsafe.Pointer(cret))

		// the returned Structure is borrowed, so keep the promise alive:
		runtime.AddCleanup(goret, func(_ *Promise) {}, promise)
	}

	return goret
}

// Interrupt wraps gst_promise_interrupt
func (promise *Promise) Interrupt() {
	var carg0 *C.GstPromise // in, none, converted

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))

	C.gst_promise_interrupt(carg0)
	runtime.KeepAlive(promise)
}

// Reply wraps gst_promise_reply
func (promise *Promise) Reply(s *Structure) {
	var carg0 *C.GstPromise   // in, none, converted
	var carg1 *C.GstStructure // in, full, converted, nullable

	carg0 = (*C.GstPromise)(UnsafePromiseToGlibNone(promise))
	if s != nil {
		carg1 = (*C.GstStructure)(UnsafeStructureToGlibFull(s.Copy()))
	}

	C.gst_promise_reply(carg0, carg1)
	runtime.KeepAlive(promise)
	runtime.KeepAlive(s)
}

// WaitBlocking wraps gst_promise_wait. Prefer to use [Promise.Await] over this.
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
