package gst

// #include "gst.go.h"
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// AllocationParams wraps the GstAllocationParams.
type AllocationParams struct {
	ptr *C.GstAllocationParams
}

// FromGstAllocationParamsUnsafe wraps the given unsafe.Pointer in an AllocationParams instance.
func FromGstAllocationParamsUnsafe(alloc unsafe.Pointer) *AllocationParams {
	return &AllocationParams{ptr: (*C.GstAllocationParams)(alloc)}
}

// NewAllocationParams initializes a set of allocation params with the default
// values.
func NewAllocationParams() *AllocationParams {
	params := &AllocationParams{
		ptr: &C.GstAllocationParams{},
	}
	params.Init()
	runtime.SetFinalizer(params, (*AllocationParams).Free)
	return params
}

// Instance returns the underlying GstAllocationParams.
func (a *AllocationParams) Instance() *C.GstAllocationParams { return a.ptr }

// Init initializes these AllocationParams to their original values.
func (a *AllocationParams) Init() { C.gst_allocation_params_init(a.ptr) }

// Copy copies these AllocationParams.
func (a *AllocationParams) Copy() *AllocationParams {
	return wrapAllocationParams(C.gst_allocation_params_copy(a.ptr))
}

// Free frees the underlying AllocationParams
func (a *AllocationParams) Free() { C.gst_allocation_params_free(a.ptr) }

// GetFlags returns the flags on these AllocationParams.
func (a *AllocationParams) GetFlags() MemoryFlags { return MemoryFlags(a.ptr.flags) }

// SetFlags changes the flags on these AllocationParams. this must be used
func (a *AllocationParams) SetFlags(flags MemoryFlags) { a.ptr.flags = C.GstMemoryFlags(flags) }

// GetAlignment returns the desired alignment of the memory.
func (a *AllocationParams) GetAlignment() int64 { return int64(a.ptr.align) }

// SetAlignment sets the desired alignment of the memory.
func (a *AllocationParams) SetAlignment(align int64) { a.ptr.align = C.gsize(align) }

// GetPrefix returns the desired prefix size.
func (a *AllocationParams) GetPrefix() int64 { return int64(a.ptr.prefix) }

// SetPrefix sets the desired prefix size.
func (a *AllocationParams) SetPrefix(prefix int64) { a.ptr.prefix = C.gsize(prefix) }

// GetPadding returns the desired padding size.
func (a *AllocationParams) GetPadding() int64 { return int64(a.ptr.padding) }

// SetPadding sets the desired padding size.
func (a *AllocationParams) SetPadding(padding int64) { a.ptr.padding = C.gsize(padding) }

// Allocator is a go representation of a GstAllocator
type Allocator struct{ *Object }

// FromGstAllocatorUnsafeNone wraps the given unsafe.Pointer in an Allocator instance.
func FromGstAllocatorUnsafeNone(alloc unsafe.Pointer) *Allocator {
	return wrapAllocator(glib.TransferNone(alloc))
}

// FromGstAllocatorUnsafeFull wraps the given unsafe.Pointer in an Allocator instance.
func FromGstAllocatorUnsafeFull(alloc unsafe.Pointer) *Allocator {
	return wrapAllocator(glib.TransferFull(alloc))
}

// DefaultAllocator returns the default GstAllocator.
func DefaultAllocator() *Allocator {
	return wrapAllocator(glib.TransferFull(unsafe.Pointer(C.gst_allocator_find(nil))))
}

// Instance returns the underlying GstAllocator instance.
func (a *Allocator) Instance() *C.GstAllocator { return C.toGstAllocator(a.Unsafe()) }

// MemType returns the memory type for this allocator.
func (a *Allocator) MemType() string { return C.GoString(a.Instance().mem_type) }

// Alloc is used to allocate a new memory block with memory that is at least size big.
// The optional params can specify the prefix and padding for the memory. If nil is passed,
// no flags, no extra prefix/padding and a default alignment is used.
//
// The prefix/padding will be filled with 0 if flags contains MemoryFlagZeroPrefixed and
// MemoryFlagZeroPadded respectively.
//
// The alignment in params is given as a bitmask so that align + 1 equals the amount of bytes to
// align to. For example, to align to 8 bytes, use an alignment of 7.
func (a *Allocator) Alloc(size int64, params *AllocationParams) *Memory {
	mem := C.gst_allocator_alloc(a.Instance(), C.gsize(size), params.ptr)
	return FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// Free memory that was originally allocated with this allocator.
func (a *Allocator) Free(mem *Memory) {
	C.gst_allocator_free(a.Instance(), mem.Instance())
}
