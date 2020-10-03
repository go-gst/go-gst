package gst

// #include "gst.go.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Memory is a go representation of GstMemory. This object is implemented
// in a read-only fashion currently. You can create new memory blocks, but
// there are no methods implemented yet for modifying ones already in existence.
type Memory struct {
	ptr *C.GstMemory
}

// NewMemoryWrapped allocates a new memory block that wraps the given data.
//
// The prefix/padding must be filled with 0 if flags contains MemoryFlagZeroPrefixed
// and MemoryFlagZeroPadded respectively.
func NewMemoryWrapped(flags MemoryFlags, data []byte, maxSize, offset, size int64) *Memory {
	str := string(data)
	dataPtr := unsafe.Pointer(C.CString(str))
	mem := C.gst_memory_new_wrapped(
		C.GstMemoryFlags(flags),
		(C.gpointer)(dataPtr),
		C.gsize(maxSize),
		C.gsize(offset),
		C.gsize(size),
		nil, // TODO: Allow user to set userdata for destroy notify function
		nil, // TODO: Allow user to set destroy notify function
	)
	return wrapMemory(mem)
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
	return wrapMemory(mem)
}

// Map the data inside the memory. This function can return nil if the memory is not readable.
func (m *Memory) Map() *MapInfo {
	var mapInfo C.GstMapInfo
	C.gst_memory_map(
		(*C.GstMemory)(m.Instance()),
		(*C.GstMapInfo)(unsafe.Pointer(&mapInfo)),
		C.GST_MAP_READ,
	)
	return wrapMapInfo(&mapInfo, func() {
		C.gst_memory_unmap(m.Instance(), (*C.GstMapInfo)(unsafe.Pointer(&mapInfo)))
	})
}

// Bytes will return a byte slice of the data inside this memory if it is readable.
func (m *Memory) Bytes() []byte {
	mapInfo := m.Map()
	if mapInfo.ptr == nil {
		return nil
	}
	defer mapInfo.Unmap()
	return mapInfo.Bytes()
}

// MapInfo is a go representation of a GstMapInfo.
type MapInfo struct {
	ptr       *C.GstMapInfo
	unmapFunc func()
	Memory    unsafe.Pointer // A pointer to the GstMemory object
	Flags     MapFlags
	Data      unsafe.Pointer // A pointer to the actual data
	Size      int64
	MaxSize   int64
}

// Unmap will unmap the MapInfo.
func (m *MapInfo) Unmap() {
	if m.unmapFunc == nil {
		fmt.Println("GO-GST-WARNING: Called Unmap() on unwrapped MapInfo")
	}
	m.unmapFunc()
}

// Bytes returns a byte slice of the data inside this map info.
func (m *MapInfo) Bytes() []byte {
	return C.GoBytes(m.Data, (C.int)(m.Size))
}

// AsInt8Slice returns the contents of this map as a slice of signed 8-bit integers.
func (m *MapInfo) AsInt8Slice() []int8 {
	out := make([]int8, m.Size)
	for i, t := range (*[1 << 30]int8)(m.Data)[:m.Size:m.Size] {
		out[i] = int8(t)
	}
	return out
}

// AsInt16Slice returns the contents of this map as a slice of signed 16-bit integers.
func (m *MapInfo) AsInt16Slice() []int16 {
	out := make([]int16, m.Size)
	for i, t := range (*[1 << 30]int16)(m.Data)[:m.Size:m.Size] {
		out[i] = int16(t)
	}
	return out
}

// AsInt32Slice returns the contents of this map as a slice of signed 32-bit integers.
func (m *MapInfo) AsInt32Slice() []int32 {
	out := make([]int32, m.Size)
	for i, t := range (*[1 << 30]int32)(m.Data)[:m.Size:m.Size] {
		out[i] = int32(t)
	}
	return out
}

// AsInt64Slice returns the contents of this map as a slice of signed 64-bit integers.
func (m *MapInfo) AsInt64Slice() []int64 {
	out := make([]int64, m.Size)
	for i, t := range (*[1 << 30]int64)(m.Data)[:m.Size:m.Size] {
		out[i] = int64(t)
	}
	return out
}

// AsUint8Slice returns the contents of this map as a slice of unsigned 8-bit integers.
func (m *MapInfo) AsUint8Slice() []uint8 {
	out := make([]uint8, m.Size)
	for i, t := range (*[1 << 30]uint8)(m.Data)[:m.Size:m.Size] {
		out[i] = uint8(t)
	}
	return out
}

// AsUint16Slice returns the contents of this map as a slice of unsigned 16-bit integers.
func (m *MapInfo) AsUint16Slice() []uint16 {
	out := make([]uint16, m.Size)
	for i, t := range (*[1 << 30]uint16)(m.Data)[:m.Size:m.Size] {
		out[i] = uint16(t)
	}
	return out
}

// AsUint32Slice returns the contents of this map as a slice of unsigned 32-bit integers.
func (m *MapInfo) AsUint32Slice() []uint32 {
	out := make([]uint32, m.Size)
	for i, t := range (*[1 << 30]uint32)(m.Data)[:m.Size:m.Size] {
		out[i] = uint32(t)
	}
	return out
}

// AsUint64Slice returns the contents of this map as a slice of unsigned 64-bit integers.
func (m *MapInfo) AsUint64Slice() []uint64 {
	out := make([]uint64, m.Size)
	for i, t := range (*[1 << 30]uint64)(m.Data)[:m.Size:m.Size] {
		out[i] = uint64(t)
	}
	return out
}

func wrapMapInfo(mapInfo *C.GstMapInfo, unmapFunc func()) *MapInfo {
	return &MapInfo{
		ptr:       mapInfo,
		unmapFunc: unmapFunc,
		Memory:    unsafe.Pointer(mapInfo.memory),
		Flags:     MapFlags(mapInfo.flags),
		Data:      unsafe.Pointer(mapInfo.data),
		Size:      int64(mapInfo.size),
		MaxSize:   int64(mapInfo.maxsize),
	}
}
