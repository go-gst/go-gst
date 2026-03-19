package gst

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"runtime"
	"unsafe"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// MapInfo is a wrapper around the C struct GstMapInfo
// and implements the io.ReaderAt, io.WriterAt, io.Reader, and io.WriteCloser interfaces.
//
// See https://gstreamer.freedesktop.org/documentation/plugin-development/advanced/allocation.html for why this is needed.
//
// There are no unsafe transfer functions for this type. It needs to be freed differently depending how it was created.
type MapInfo struct {
	*mapInfo

	writeOffset int64
	readOffset  int64
}

// clearAutoCleanup clears the finalizer to prevent automatic unmapping
func (m *MapInfo) clearAutoCleanup() {
	runtime.SetFinalizer(m.mapInfo, nil)
}

func (m *MapInfo) autoCleanup() {
	runtime.SetFinalizer(
		m.mapInfo,
		func(intern *mapInfo) {
			fmt.Println("automatically unmapping MapInfo, you should call Unmap()/Close() instead at an appropriate time")
			intern.unmap()
		},
	)
}

// mapInfo is the struct that is finalized
type mapInfo struct {
	buffer *Buffer

	native *C.GstMapInfo
}

var _ io.WriteCloser = (*MapInfo)(nil)
var _ io.Reader = (*MapInfo)(nil)
var _ io.WriterAt = (*MapInfo)(nil)
var _ io.ReaderAt = (*MapInfo)(nil)

var ErrMapInfoNotReadable = fmt.Errorf("MapInfo is not readable")
var ErrMapInfoNotWritable = fmt.Errorf("MapInfo is not writable")
var ErrMapInfoInvalidOffset = fmt.Errorf("MapInfo invalid offset")
var ErrMapInfoInvalid = fmt.Errorf("MapInfo is invalid")

// ReadAt implements io.ReaderAt.
func (m *MapInfo) ReadAt(p []byte, off int64) (n int, err error) {
	// check for valid MapInfo
	if m.mapInfo == nil {
		return 0, ErrMapInfoInvalid
	}
	// check for read access
	if !m.Flags().Has(MapRead) {
		return 0, ErrMapInfoNotReadable
	}

	if off < 0 {
		return 0, ErrMapInfoInvalidOffset
	}
	// check for EOF
	if off >= int64(m.Length()) {
		return 0, io.EOF
	}

	data := m.Data()

	n = copy(p, data[off:])
	if n == 0 {
		return 0, io.EOF
	}

	return n, nil
}

// WriteAt implements io.WriterAt.
func (m *MapInfo) WriteAt(p []byte, off int64) (n int, err error) {
	// check for valid MapInfo
	if m.mapInfo == nil {
		return 0, ErrMapInfoInvalid
	}
	// check for write access
	if !m.Flags().Has(MapWrite) {
		return 0, ErrMapInfoNotWritable
	}

	if off < 0 {
		return 0, ErrMapInfoInvalidOffset
	}
	// check for EOF
	if off >= int64(m.Length()) {
		return 0, ErrMapInfoInvalidOffset
	}

	data := m.Data()

	n = copy(data[off:], p)

	return n, nil
}

// Read implements io.Reader.
func (m *MapInfo) Read(p []byte) (n int, err error) {
	off := m.readOffset
	n, err = m.ReadAt(p, off)
	m.readOffset += int64(n)
	return n, err
}

// Close implements io.WriteCloser. It calls Unmap() to release the memory.
func (m *MapInfo) Close() error {
	m.Unmap()
	return nil
}

// Write implements io.WriteCloser.
func (m *MapInfo) Write(p []byte) (n int, err error) {
	off := m.writeOffset
	n, err = m.WriteAt(p, off)
	m.writeOffset += int64(n)
	return n, err
}

// unmap is the private unmap function used by the finalizer as well as the manual close/unmap
func (info *mapInfo) unmap() {
	// this needs a different function to unmap depending on memory/buffer unmap
	if info.buffer != nil {
		C.gst_buffer_unmap((*C.GstBuffer)(UnsafeBufferToGlibNone(info.buffer)), info.native)
	} else {
		panic("unmap called on an invalid MapInfo")
	}
}

// Reset resets the read/write offsets to 0.
func (info *MapInfo) Reset() {
	info.readOffset = 0
	info.writeOffset = 0
}

// Unmap Releases the memory previously mapped.
// The MapInfo is invalid after this call.
func (info *MapInfo) Unmap() {
	info.clearAutoCleanup()
	info.mapInfo.unmap()
	info.mapInfo = nil
}

// Length returns the length of the mapped memory.
func (info *MapInfo) Length() int {
	return int(info.mapInfo.native.size)
}

// Flags returns the flags of the mapped memory.
func (info *MapInfo) Flags() MapFlags {
	return MapFlags(info.mapInfo.native.flags)
}

// Data returns the mapped memory as a byte slice.
func (info *MapInfo) Data() []byte {
	if info.mapInfo == nil {
		return nil
	}

	return unsafe.Slice((*byte)(info.mapInfo.native.data), info.mapInfo.native.size)
}

// Float32Data returns a copy of the data as a slice of float32 with the given byte order (endianess).
func (info *MapInfo) Float32Data(byteOrder binary.ByteOrder) []float32 {

	floats := make([]float32, info.Length()/4)

	for i := range floats {
		bits := byteOrder.Uint32(info.Data()[i*4 : (i+1)*4])
		floats[i] = math.Float32frombits(bits)
	}

	return floats
}

// Float64Data returns a copy of the data as a slice of float64 with the given byte order (endianess).
func (info *MapInfo) Float64Data(byteOrder binary.ByteOrder) []float64 {

	floats := make([]float64, info.Length()/8)

	for i := range floats {
		bits := byteOrder.Uint64(info.Data()[i*8 : (i+1)*8])
		floats[i] = math.Float64frombits(bits)
	}

	return floats
}

// unsafeData is used to save on a cast to unsafe.Pointer in the ...Data() functions
func (info *MapInfo) unsafeData() unsafe.Pointer {
	return unsafe.Pointer(info.mapInfo.native.data)
}

// Int8Data returns the mapped data as a slice of int8.
func (info *MapInfo) Int8Data() []int8 {
	if info.mapInfo == nil {
		return nil
	}
	if !info.Flags().Has(MapRead) {
		return nil
	}

	return unsafe.Slice((*int8)(info.unsafeData()), info.mapInfo.native.size)
}

// Uint8Data returns the mapped data as a slice of uint8.
func (info *MapInfo) Uint8Data() []uint8 {
	if info.mapInfo == nil {
		return nil
	}
	if !info.Flags().Has(MapRead) {
		return nil
	}

	return unsafe.Slice((*uint8)(info.unsafeData()), info.mapInfo.native.size)
}

// Uint16Data returns a copy of the data as a slice of uint16 with the given byte order (endianess).
func (info *MapInfo) Uint16Data(byteOrder binary.ByteOrder) []uint16 {
	data := info.Data()
	ints := make([]uint16, len(data)/2)

	for i := range ints {
		bits := byteOrder.Uint16(data[i*2 : (i+1)*2])
		ints[i] = bits
	}

	return ints
}

// Int16Data returns a copy of the data as a slice of int16 with the given byte order (endianess).
func (info *MapInfo) Int16Data(byteOrder binary.ByteOrder) []int16 {
	data := info.Data()
	ints := make([]int16, len(data)/2)

	for i := range ints {
		bits := byteOrder.Uint16(data[i*2 : (i+1)*2])
		ints[i] = int16(bits)
	}

	return ints
}

// Uint32Data returns a copy of the data as a slice of uint32 with the given byte order (endianess).
func (info *MapInfo) Uint32Data(byteOrder binary.ByteOrder) []uint32 {
	data := info.Data()
	ints := make([]uint32, len(data)/4)

	for i := range ints {
		bits := byteOrder.Uint32(data[i*4 : (i+1)*4])
		ints[i] = bits
	}

	return ints
}

// Int32Data returns a copy of the data as a slice of int32 with the given byte order (endianess).
func (info *MapInfo) Int32Data(byteOrder binary.ByteOrder) []int32 {
	data := info.Data()
	ints := make([]int32, len(data)/4)

	for i := range ints {
		bits := byteOrder.Uint32(data[i*4 : (i+1)*4])
		ints[i] = int32(bits)
	}

	return ints
}

// Uint64Data returns a copy of the data as a slice of uint64 with the given byte order (endianess).
func (info *MapInfo) Uint64Data(byteOrder binary.ByteOrder) []uint64 {
	data := info.Data()
	ints := make([]uint64, len(data)/8)

	for i := range ints {
		bits := byteOrder.Uint64(data[i*8 : (i+1)*8])
		ints[i] = bits
	}

	return ints
}

// Int64Data returns a copy of the data as a slice of int64 with the given byte order (endianess).
func (info *MapInfo) Int64Data(byteOrder binary.ByteOrder) []int64 {
	data := info.Data()
	ints := make([]int64, len(data)/8)

	for i := range ints {
		bits := byteOrder.Uint64(data[i*8 : (i+1)*8])
		ints[i] = int64(bits)
	}

	return ints
}
