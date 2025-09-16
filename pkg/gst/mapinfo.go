package gst

import (
	"fmt"
	"io"
	"runtime"
	"unsafe"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// Map wraps gst_buffer_map
//
// The function takes the following parameters:
//
//   - flags MapFlags: flags for the mapping
//
// The function returns the following values:
//
//   - info MapInfo: info about the mapping
//   - goret bool
//
// Fills @info with the #GstMapInfo of all merged memory blocks in @buffer.
//
// @flags describe the desired access of the memory. When @flags is
// #GST_MAP_WRITE, @buffer should be writable (as returned from
// gst_buffer_is_writable()).
//
// When @buffer is writable but the memory isn't, a writable copy will
// automatically be created and returned. The readonly copy of the
// buffer memory will then also be replaced with this writable copy.
//
// The memory in @info should be unmapped with gst_buffer_unmap() after
// usage.
func (buffer *Buffer) Map(flags MapFlags) (*MapInfo, bool) {
	var carg0 *C.GstBuffer  // in, none, converted
	var carg2 C.GstMapFlags // in, none, casted
	var carg1 C.GstMapInfo  // out, transfer: none, C Pointers: 0, Name: MapInfo, caller-allocates
	var cret C.gboolean     // return

	carg0 = (*C.GstBuffer)(UnsafeBufferToGlibNone(buffer))
	carg2 = C.GstMapFlags(flags)

	cret = C.gst_buffer_map(carg0, &carg1, carg2)
	runtime.KeepAlive(buffer)
	runtime.KeepAlive(flags)

	var info *MapInfo
	var goret bool

	info = &MapInfo{
		mapInfo: &mapInfo{
			native: &carg1,
			buffer: buffer,
		},
	}

	info.autoCleanup()

	if cret != 0 {
		goret = true
	}

	return info, goret
}

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

// Unmap Releases the memory previously mapped.
func (info *MapInfo) Unmap() {
	info.mapInfo.unmap()
	info.mapInfo = nil
	runtime.SetFinalizer(info, nil)
}

// Length returns the length of the mapped memory.
func (info *MapInfo) Length() int {
	return int(info.mapInfo.native.size)
}

// Length returns the length of the mapped memory.
func (info *MapInfo) Data() []byte {
	return unsafe.Slice((*byte)(info.mapInfo.native.data), info.mapInfo.native.size)
}

// Flags returns the flags of the mapped memory.
func (info *MapInfo) Flags() MapFlags {
	return MapFlags(info.mapInfo.native.flags)
}
