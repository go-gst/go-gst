package gst

/*
#include "gst.go.h"

void memcpy_offset (void * dest, guint offset, const void * src, size_t n) { memcpy(dest + offset, src, n); }


GstByteReader * newByteReader (const guint8 * data, guint size)
{
	GstByteReader *ret = g_slice_new0 (GstByteReader);

	ret->data = data;
  	ret->size = size;

  	return ret;
}

void freeByteReader (GstByteReader * reader)
{
  g_return_if_fail (reader != NULL);
  g_slice_free (GstByteReader, reader);
}


*/
import "C"

import (
	"bytes"
	"errors"
	"io"
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

// Reader returns a Reader for the contents of this map's memory.
func (m *MapInfo) Reader() io.Reader {
	return bytes.NewReader(m.Bytes())
}

type mapInfoWriter struct {
	mapInfo *MapInfo
	wsize   int
}

func (m *mapInfoWriter) Write(p []byte) (int, error) {
	if m.wsize+len(p) > int(m.mapInfo.Size()) {
		return 0, errors.New("Attempt to write more data than allocated to MapInfo")
	}
	m.mapInfo.WriteAt(m.wsize, p)
	m.wsize += len(p)
	return len(p), nil
}

// Writer returns a writer that will copy the contents passed to Write directly to the
// map's memory sequentially. An error will be returned if an attempt is made to write
// more data than the map can hold.
func (m *MapInfo) Writer() io.Writer {
	return &mapInfoWriter{
		mapInfo: m,
		wsize:   0,
	}
}

// WriteData writes the given sequence directly to the map's memory.
func (m *MapInfo) WriteData(data []byte) {
	C.memcpy(unsafe.Pointer(m.ptr.data), unsafe.Pointer(&data[0]), C.gsize(len(data)))
}

// WriteAt writes the given data sequence directly to the map's memory at the given offset.
func (m *MapInfo) WriteAt(offset int, data []byte) {
	C.memcpy_offset(unsafe.Pointer(m.ptr.data), C.guint(offset), unsafe.Pointer(&data[0]), C.gsize(len(data)))
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
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]int8, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gint8
		C.gst_byte_reader_get_int8(br, &gint)
		out = append(out, int8(gint))
	}
	return out
}

// AsInt16BESlice returns the contents of this map as a slice of signed 16-bit big-endian integers.
func (m *MapInfo) AsInt16BESlice() []int16 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]int16, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gint16
		C.gst_byte_reader_get_int16_be(br, &gint)
		out = append(out, int16(gint))
	}
	return out
}

// AsInt16LESlice returns the contents of this map as a slice of signed 16-bit little-endian integers.
func (m *MapInfo) AsInt16LESlice() []int16 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]int16, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gint16
		C.gst_byte_reader_get_int16_le(br, &gint)
		out = append(out, int16(gint))
	}
	return out
}

// AsInt32BESlice returns the contents of this map as a slice of signed 32-bit big-endian integers.
func (m *MapInfo) AsInt32BESlice() []int32 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]int32, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gint32
		C.gst_byte_reader_get_int32_be(br, &gint)
		out = append(out, int32(gint))
	}
	return out
}

// AsInt32LESlice returns the contents of this map as a slice of signed 32-bit little-endian integers.
func (m *MapInfo) AsInt32LESlice() []int32 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]int32, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gint32
		C.gst_byte_reader_get_int32_le(br, &gint)
		out = append(out, int32(gint))
	}
	return out
}

// AsInt64BESlice returns the contents of this map as a slice of signed 64-bit big-endian integers.
func (m *MapInfo) AsInt64BESlice() []int64 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]int64, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gint64
		C.gst_byte_reader_get_int64_be(br, &gint)
		out = append(out, int64(gint))
	}
	return out
}

// AsInt64LESlice returns the contents of this map as a slice of signed 64-bit little-endian integers.
func (m *MapInfo) AsInt64LESlice() []int64 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]int64, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gint64
		C.gst_byte_reader_get_int64_le(br, &gint)
		out = append(out, int64(gint))
	}
	return out
}

// AsUint8Slice returns the contents of this map as a slice of unsigned 8-bit integers.
func (m *MapInfo) AsUint8Slice() []uint8 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]uint8, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.guint8
		C.gst_byte_reader_get_uint8(br, &gint)
		out = append(out, uint8(gint))
	}
	return out
}

// AsUint16BESlice returns the contents of this map as a slice of unsigned 16-bit big-endian integers.
func (m *MapInfo) AsUint16BESlice() []uint16 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]uint16, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.guint16
		C.gst_byte_reader_get_uint16_be(br, &gint)
		out = append(out, uint16(gint))
	}
	return out
}

// AsUint16LESlice returns the contents of this map as a slice of unsigned 16-bit little-endian integers.
func (m *MapInfo) AsUint16LESlice() []uint16 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]uint16, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.guint16
		C.gst_byte_reader_get_uint16_le(br, &gint)
		out = append(out, uint16(gint))
	}
	return out
}

// AsUint32BESlice returns the contents of this map as a slice of unsigned 32-bit big-endian integers.
func (m *MapInfo) AsUint32BESlice() []uint32 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]uint32, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.guint32
		C.gst_byte_reader_get_uint32_be(br, &gint)
		out = append(out, uint32(gint))
	}
	return out
}

// AsUint32LESlice returns the contents of this map as a slice of unsigned 32-bit little-endian integers.
func (m *MapInfo) AsUint32LESlice() []uint32 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]uint32, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.guint32
		C.gst_byte_reader_get_uint32_le(br, &gint)
		out = append(out, uint32(gint))
	}
	return out
}

// AsUint64BESlice returns the contents of this map as a slice of unsigned 64-bit big-endian integers.
func (m *MapInfo) AsUint64BESlice() []uint64 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]uint64, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.guint64
		C.gst_byte_reader_get_uint64_be(br, &gint)
		out = append(out, uint64(gint))
	}
	return out
}

// AsUint64LESlice returns the contents of this map as a slice of unsigned 64-bit little-endian integers.
func (m *MapInfo) AsUint64LESlice() []uint64 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]uint64, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.guint64
		C.gst_byte_reader_get_uint64_le(br, &gint)
		out = append(out, uint64(gint))
	}
	return out
}

// AsFloat32BESlice returns the contents of this map as a slice of 32-bit big-endian floats.
func (m *MapInfo) AsFloat32BESlice() []float32 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]float32, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gfloat
		C.gst_byte_reader_get_float32_be(br, &gint)
		out = append(out, float32(gint))
	}
	return out
}

// AsFloat32LESlice returns the contents of this map as a slice of 32-bit little-endian floats.
func (m *MapInfo) AsFloat32LESlice() []float32 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]float32, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gfloat
		C.gst_byte_reader_get_float32_le(br, &gint)
		out = append(out, float32(gint))
	}
	return out
}

// AsFloat64BESlice returns the contents of this map as a slice of 64-bit big-endian floats.
func (m *MapInfo) AsFloat64BESlice() []float64 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]float64, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gdouble
		C.gst_byte_reader_get_float64_be(br, &gint)
		out = append(out, float64(gint))
	}
	return out
}

// AsFloat64LESlice returns the contents of this map as a slice of 64-bit little-endian floats.
func (m *MapInfo) AsFloat64LESlice() []float64 {
	br := C.newByteReader(m.Instance().data, C.guint(m.Instance().size))
	defer C.freeByteReader(br)
	out := make([]float64, 0)
	for C.gst_byte_reader_get_remaining(br) != C.guint(0) {
		var gint C.gdouble
		C.gst_byte_reader_get_float64_le(br, &gint)
		out = append(out, float64(gint))
	}
	return out
}
