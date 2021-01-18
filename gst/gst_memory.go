package gst

/*
#include "gst.go.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Memory is a go representation of GstMemory. This object is implemented in a read-only fashion
// currently primarily for reference, and as such you should not really use it. You can create new
// memory blocks, but there are no methods implemented yet for modifying ones already in existence.
//
// Use the Buffer and its Map methods to interact with memory in both a read and writable way.
type Memory struct {
	ptr     *C.GstMemory
	mapInfo *MapInfo
}

// FromGstMemoryUnsafe wraps the given C GstMemory in the go type. It is meant for internal usage
// and exported for visibility to other packages.
func FromGstMemoryUnsafe(mem unsafe.Pointer) *Memory {
	wrapped := wrapMemory((*C.GstMemory)(mem))
	wrapped.Ref()
	runtime.SetFinalizer(wrapped, (*Memory).Unref)
	return wrapped
}

// FromGstMemoryUnsafeNone is an alias to FromGstMemoryUnsafe.
func FromGstMemoryUnsafeNone(mem unsafe.Pointer) *Memory {
	return FromGstMemoryUnsafe(mem)
}

// FromGstMemoryUnsafeFull wraps the given memory without taking an additional reference.
func FromGstMemoryUnsafeFull(mem unsafe.Pointer) *Memory {
	wrapped := wrapMemory((*C.GstMemory)(mem))
	runtime.SetFinalizer(wrapped, (*Memory).Unref)
	return wrapped
}

// NewMemoryWrapped allocates a new memory block that wraps the given data.
//
// The prefix/padding must be filled with 0 if flags contains MemoryFlagZeroPrefixed
// and MemoryFlagZeroPadded respectively.
func NewMemoryWrapped(flags MemoryFlags, data []byte, maxSize, offset int64) *Memory {
	mem := C.gst_memory_new_wrapped(
		C.GstMemoryFlags(flags),
		(C.gpointer)(unsafe.Pointer(&data[0])),
		C.gsize(maxSize),
		C.gsize(offset),
		C.gsize(len(data)),
		nil, // TODO: Allow user to set userdata for destroy notify function
		nil, // TODO: Allow user to set destroy notify function
	)
	return FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// Instance returns the underlying GstMemory instance.
func (m *Memory) Instance() *C.GstMemory { return C.toGstMemory(unsafe.Pointer(m.ptr)) }

// Ref increases the ref count on this memory block by one.
func (m *Memory) Ref() *Memory {
	return wrapMemory(C.gst_memory_ref(m.Instance()))
}

// Unref decreases the ref count on this memory block by one. When the refcount reaches
// zero the memory is freed.
func (m *Memory) Unref() { C.gst_memory_unref(m.Instance()) }

// Allocator returns the allocator for this memory.
func (m *Memory) Allocator() *Allocator {
	return wrapAllocator(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(m.Instance().allocator))})
}

// Parent returns this memory block's parent.
func (m *Memory) Parent() *Memory { return wrapMemory(m.Instance().parent) }

// MaxSize returns the maximum size allocated for this memory block.
func (m *Memory) MaxSize() int64 { return int64(m.Instance().maxsize) }

// Alignment returns the alignment of the memory.
func (m *Memory) Alignment() int64 { return int64(m.Instance().align) }

// Offset returns the offset where valid data starts.
func (m *Memory) Offset() int64 { return int64(m.Instance().offset) }

// Size returns the size of valid data.
func (m *Memory) Size() int64 { return int64(m.Instance().size) }

// Copy returns a copy of size bytes from mem starting from offset. This copy is
// guaranteed to be writable. size can be set to -1 to return a copy from offset
// to the end of the memory region.
func (m *Memory) Copy(offset, size int64) *Memory {
	mem := C.gst_memory_copy(m.Instance(), C.gssize(offset), C.gssize(size))
	return FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// Map the data inside the memory. This function can return nil if the memory is not read or writable.
// It is safe to call this function multiple times on the same Memory, however it will retain the flags
// used when mapping the first time. To change between read and write access first unmap and then remap the
// memory with the appropriate flags, or map initially with both read/write access.
//
// Unmap the Memory after usage.
func (m *Memory) Map(flags MapFlags) *MapInfo {
	if m.mapInfo != nil {
		return m.mapInfo
	}
	mapInfo := C.malloc(C.sizeof_GstMapInfo)
	C.gst_memory_map(
		(*C.GstMemory)(m.Instance()),
		(*C.GstMapInfo)(mapInfo),
		C.GstMapFlags(flags),
	)
	if mapInfo == C.NULL {
		return nil
	}
	m.mapInfo = wrapMapInfo((*C.GstMapInfo)(mapInfo))
	return m.mapInfo
}

// Unmap will unmap the data inside this memory. Use this after calling Map on the Memory.
func (m *Memory) Unmap() {
	if m.mapInfo == nil {
		return
	}
	C.gst_memory_unmap(m.Instance(), (*C.GstMapInfo)(m.mapInfo.Instance()))
	C.free(unsafe.Pointer(m.mapInfo.Instance()))
}

// Bytes will return a byte slice of the data inside this memory if it is readable.
func (m *Memory) Bytes() []byte {
	mapInfo := m.Map(MapRead)
	if mapInfo == nil {
		return nil
	}
	defer m.Unmap()
	return mapInfo.Bytes()
}
