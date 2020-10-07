package gst

/*
#include "gst.go.h"

void writeMapData (GstMapInfo * mapInfo, gint idx, guint8 data) { mapInfo->data[idx] = data; }
*/
import "C"

import (
	"encoding/binary"
	"runtime"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Memory is a go representation of GstMemory. This object is implemented in a read-only fashion
// currently primarily for reference, and as such you should not really use it. You can create new
// memory blocks, but there are no methods implemented yet for modifying ones already in existence.
//
// Use the Buffer and its Map methods to interact with memory in both a read and writable way.
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
	return mapInfo.Bytes()
}

// MapInfo is a go representation of a GstMapInfo.
type MapInfo struct {
	ptr *C.GstMapInfo
}

// Memory returns the underlying memory object.
func (m *MapInfo) Memory() *Memory {
	return wrapMemory(m.ptr.memory)
}

// Data returns a pointer to the raw data inside this map.
func (m *MapInfo) Data() unsafe.Pointer {
	return unsafe.Pointer(m.ptr.data)
}

// Flags returns the flags set on this map.
func (m *MapInfo) Flags() MapFlags {
	return MapFlags(m.ptr.flags)
}

// Size returrns the size of this map.
func (m *MapInfo) Size() int64 {
	return int64(m.ptr.size)
}

// MaxSize returns the maximum size of this map.
func (m *MapInfo) MaxSize() int64 {
	return int64(m.ptr.maxsize)
}

// Bytes returns a byte slice of the data inside this map info.
func (m *MapInfo) Bytes() []byte {
	return C.GoBytes(m.Data(), (C.int)(m.Size()))
}

// AsInt8Slice returns the contents of this map as a slice of signed 8-bit integers.
func (m *MapInfo) AsInt8Slice() []int8 {
	uint8sl := m.AsUint8Slice()
	out := make([]int8, m.Size())
	for i := range out {
		out[i] = int8(uint8sl[i])
	}
	return out
}

// AsInt16Slice returns the contents of this map as a slice of signed 16-bit integers.
func (m *MapInfo) AsInt16Slice() []int16 {
	uint8sl := m.AsUint8Slice()
	out := make([]int16, m.Size()/2)
	for i := range out {
		out[i] = int16(binary.LittleEndian.Uint16(uint8sl[i*2 : (i+1)*2]))
	}
	return out
}

// AsInt32Slice returns the contents of this map as a slice of signed 32-bit integers.
func (m *MapInfo) AsInt32Slice() []int32 {
	uint8sl := m.AsUint8Slice()
	out := make([]int32, m.Size()/4)
	for i := range out {
		out[i] = int32(binary.LittleEndian.Uint32(uint8sl[i*4 : (i+1)*4]))
	}
	return out
}

// AsInt64Slice returns the contents of this map as a slice of signed 64-bit integers.
func (m *MapInfo) AsInt64Slice() []int64 {
	uint8sl := m.AsUint8Slice()
	out := make([]int64, m.Size()/8)
	for i := range out {
		out[i] = int64(binary.LittleEndian.Uint64(uint8sl[i*8 : (i+1)*8]))
	}
	return out
}

// AsUint8Slice returns the contents of this map as a slice of unsigned 8-bit integers.
func (m *MapInfo) AsUint8Slice() []uint8 {
	out := make([]uint8, m.Size())
	for i, t := range (*[1 << 30]uint8)(m.Data())[:m.Size():m.Size()] {
		out[i] = t
	}
	return out
}

// AsUint16Slice returns the contents of this map as a slice of unsigned 16-bit integers.
func (m *MapInfo) AsUint16Slice() []uint16 {
	uint8sl := m.AsUint8Slice()
	out := make([]uint16, m.Size()/2)
	for i := range out {
		out[i] = uint16(binary.LittleEndian.Uint16(uint8sl[i*2 : (i+1)*2]))
	}
	return out
}

// AsUint32Slice returns the contents of this map as a slice of unsigned 32-bit integers.
func (m *MapInfo) AsUint32Slice() []uint32 {
	uint8sl := m.AsUint8Slice()
	out := make([]uint32, m.Size()/4)
	for i := range out {
		out[i] = uint32(binary.LittleEndian.Uint32(uint8sl[i*4 : (i+1)*4]))
	}
	return out
}

// AsUint64Slice returns the contents of this map as a slice of unsigned 64-bit integers.
func (m *MapInfo) AsUint64Slice() []uint64 {
	uint8sl := m.AsUint8Slice()
	out := make([]uint64, m.Size()/8)
	for i := range out {
		out[i] = uint64(binary.LittleEndian.Uint64(uint8sl[i*8 : (i+1)*8]))
	}
	return out
}

// WriteData writes the given values directly to the map's memory.
func (m *MapInfo) WriteData(data []uint8) {
	for i, x := range data {
		C.writeMapData(m.ptr, C.gint(i), C.guint8(x))
	}
}

func wrapMapInfo(mapInfo *C.GstMapInfo, unmapFunc func()) *MapInfo {
	info := &MapInfo{ptr: mapInfo}
	runtime.SetFinalizer(info, func(_ *MapInfo) { unmapFunc() })
	return info
}
