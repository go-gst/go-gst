package gst

/*
#include "gst.go.h"
*/
import "C"

import (
	"encoding/binary"
	"unsafe"
)

// MapInfo is a go representation of a GstMapInfo.
type MapInfo struct {
	ptr *C.GstMapInfo
}

// Instance returns the underlying GstMapInfo instance.
func (m *MapInfo) Instance() *C.GstMapInfo {
	return m.ptr
}

// WriteData writes the given sequence directly to the map's memory.
func (m *MapInfo) WriteData(data []byte) {
	C.memcpy(unsafe.Pointer(m.ptr.data), unsafe.Pointer(&data[0]), C.gsize(len(data)))
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
