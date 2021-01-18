package gst

// #include "gst.go.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// Context wraps a GstContext object.
type Context struct {
	ptr *C.GstContext
}

// FromGstContextUnsafeFull wraps the given context and places a runtime finalizer on it.
func FromGstContextUnsafeFull(ctx unsafe.Pointer) *Context {
	wrapped := wrapContext((*C.GstContext)(ctx))
	runtime.SetFinalizer(wrapped, (*Context).Unref)
	return wrapped
}

// FromGstContextUnsafeNone refs and wraps the given context and places a runtime finalizer on it.
func FromGstContextUnsafeNone(ctx unsafe.Pointer) *Context {
	wrapped := wrapContext((*C.GstContext)(ctx))
	wrapped.Ref()
	runtime.SetFinalizer(wrapped, (*Context).Unref)
	return wrapped
}

// NewContext creates a new context.
//
//   // Example
//
//   ctx := gst.NewContext("test-context", false)
//   fmt.Println(ctx.IsPersistent())
//
//   ctx = gst.NewContext("test-context", true)
//   fmt.Println(ctx.IsPersistent())
//
//   // false
//   // true
//
func NewContext(ctxType string, persistent bool) *Context {
	cStr := C.CString(ctxType)
	defer C.free(unsafe.Pointer(cStr))
	ctx := C.gst_context_new((*C.gchar)(unsafe.Pointer(cStr)), gboolean(persistent))
	if ctx == nil {
		return nil
	}
	return FromGstContextUnsafeFull(unsafe.Pointer(ctx))
}

// Instance returns the underlying GstContext instance.
func (c *Context) Instance() *C.GstContext { return C.toGstContext(unsafe.Pointer(c.ptr)) }

// GetType returns the type of the context.
//
//   // Example
//
//   ctx := gst.NewContext("test-context", false)
//   fmt.Println(ctx.GetType())
//
//   // test-context
//
func (c *Context) GetType() string {
	return C.GoString(C.gst_context_get_context_type(c.Instance()))
}

// GetStructure returns the structure of the context. You should not modify or unref it.
func (c *Context) GetStructure() *Structure {
	st := C.gst_context_get_structure(c.Instance())
	if st == nil {
		return nil
	}
	return wrapStructure(st)
}

// HasContextType checks if the context has the given type.
//
//   // Example
//
//   ctx := gst.NewContext("test-context", false)
//   fmt.Println(ctx.HasContextType("test-context"))
//   fmt.Println(ctx.HasContextType("another-context"))
//
//   // true
//   // false
//
func (c *Context) HasContextType(ctxType string) bool {
	cStr := C.CString(ctxType)
	defer C.free(unsafe.Pointer(cStr))
	return gobool(C.gst_context_has_context_type(
		c.Instance(),
		(*C.gchar)(unsafe.Pointer(cStr)),
	))
}

// IsPersistent checks if the context is persistent.
func (c *Context) IsPersistent() bool {
	return gobool(C.gst_context_is_persistent(c.Instance()))
}

// IsWritable returns true if the context is writable.
func (c *Context) IsWritable() bool {
	return gobool(C.contextIsWritable(c.Instance()))
}

// MakeWritable returns a writable version of the context.
func (c *Context) MakeWritable() *Context {
	return FromGstContextUnsafeFull(unsafe.Pointer(C.makeContextWritable(c.Instance())))
}

// WritableStructure returns a writable version of the structure. You should still not unref it.
func (c *Context) WritableStructure() *Structure {
	st := C.gst_context_writable_structure(c.Instance())
	if st == nil {
		return nil
	}
	return wrapStructure(st)
}

// Ref increases the ref count on the Context.
func (c *Context) Ref() *Context {
	ctx := C.gst_context_ref(c.Instance())
	return &Context{ptr: ctx}
}

// Unref decreases the ref count on the Context.
func (c *Context) Unref() { C.gst_context_unref(c.Instance()) }
