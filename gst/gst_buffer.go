package gst

/*
#include "gst.go.h"

extern void goGDestroyNotifyFunc (gpointer data);

void cgoDestroyNotifyFunc (gpointer data) {
	goGDestroyNotifyFunc(data);
}
*/
import "C"

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	gopointer "github.com/mattn/go-pointer"
)

// Buffer is a go representation of a GstBuffer.
type Buffer struct {
	ptr *C.GstBuffer
}

// NewEmptyBuffer returns a new empty buffer.
func NewEmptyBuffer() *Buffer {
	return wrapBuffer(C.gst_buffer_new())
}

// NewBufferAllocate tries to create a newly allocated buffer with data of the given size
// and extra parameters from allocator. If the requested amount of memory can't be allocated,
// nil will be returned. The allocated buffer memory is not cleared.
//
// When allocator is nil, the default memory allocator will be used.
//
// Note that when size == 0, the buffer will not have memory associated with it.
func NewBufferAllocate(alloc *Allocator, params *AllocationParams, size int64) *Buffer {
	var gstalloc *C.GstAllocator
	if alloc != nil {
		gstalloc = alloc.Instance()
	}
	buf := C.gst_buffer_new_allocate(gstalloc, C.gsize(size), params.Instance())
	if buf == nil {
		return nil
	}
	return wrapBuffer(buf)
}

// NewBufferFromBytes returns a new buffer from the given byte slice.
func NewBufferFromBytes(b []byte) *Buffer {
	str := string(b)
	p := unsafe.Pointer(C.CString(str))
	// memory is freed by gstreamer after building the new buffer
	buf := C.gst_buffer_new_wrapped((C.gpointer)(p), C.ulong(len(str)))
	return wrapBuffer(buf)
}

// NewBufferFromReader returns a new buffer from the given io.Reader.
func NewBufferFromReader(rdr io.Reader) (*Buffer, error) {
	out, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	return NewBufferFromBytes(out), nil
}

// NewBufferFull allocates a new buffer that wraps the given data. The wrapped buffer will
// have the region from offset and size visible. The maxsize must be at least the size of the
// data provided.
//
// When the buffer is destroyed, notifyFunc will be called if it is not nil.
//
// The prefix/padding must be filled with 0 if flags contains MemoryFlagZeroPrefixed and MemoryFlagZeroPadded respectively.
//
// Example
//
//     buf := gst.NewBufferFull(0, []byte("hello-world"), 1024, 0, 1024, func() {
// 	       fmt.Println("buffer was destroyed")
//     })
//     if buf != nil {
// 	       buf.Unref()
//     }
//
//     // > buffer was destroyed
func NewBufferFull(flags MemoryFlags, data []byte, maxSize, offset, size int64, notifyFunc func()) *Buffer {
	var notifyData unsafe.Pointer
	var gnotifyFunc C.GDestroyNotify
	if notifyFunc != nil {
		notifyData = gopointer.Save(notifyFunc)
		gnotifyFunc = C.GDestroyNotify(C.cgoDestroyNotifyFunc)
	}
	dataStr := string(data)
	dataPtr := unsafe.Pointer(C.CString(dataStr))
	buf := C.gst_buffer_new_wrapped_full(
		C.GstMemoryFlags(flags),
		(C.gpointer)(dataPtr),
		C.gsize(maxSize), C.gsize(offset), C.gsize(size),
		(C.gpointer)(notifyData), gnotifyFunc,
	)
	if buf == nil {
		return nil
	}
	return wrapBuffer(buf)
}

// Instance returns the underlying GstBuffer instance.
func (b *Buffer) Instance() *C.GstBuffer { return C.toGstBuffer(unsafe.Pointer(b.ptr)) }

// Ref increases the ref count on the buffer by one.
func (b *Buffer) Ref() *Buffer { return wrapBuffer(C.gst_buffer_ref(b.Instance())) }

// Unref decreaes the ref count on the buffer by one. When the refcount reaches zero, the memory is freed.
func (b *Buffer) Unref() { C.gst_buffer_unref(b.Instance()) }

// Reader returns an io.Reader for this buffer.
func (b *Buffer) Reader() io.Reader { return bytes.NewBuffer(b.Bytes()) }

// Bytes returns a byte slice of the data inside this buffer.
func (b *Buffer) Bytes() []byte {
	mapInfo := b.Map()
	if mapInfo.ptr == nil {
		return nil
	}
	defer mapInfo.Unmap()
	return mapInfo.Bytes()
}

// PresentationTimestamp returns the presentation timestamp of the buffer, or a negative duration
// if not known or relevant. This value contains the timestamp when the media should be
// presented to the user.
func (b *Buffer) PresentationTimestamp() time.Duration {
	pts := b.Instance().pts
	if ClockTime(pts) == ClockTimeNone {
		return time.Duration(-1)
	}
	return guint64ToDuration(pts)
}

// DecodingTimestamp returns the decoding timestamp of the buffer, or a negative duration if not known
// or relevant. This value contains the timestamp when the media should be processed.
func (b *Buffer) DecodingTimestamp() time.Duration {
	dts := b.Instance().dts
	if ClockTime(dts) == ClockTimeNone {
		return time.Duration(-1)
	}
	return guint64ToDuration(dts)
}

// Duration returns the length of the data inside this buffer, or a negative duration if not known
// or relevant.
func (b *Buffer) Duration() time.Duration {
	dur := b.Instance().duration
	if ClockTime(dur) == ClockTimeNone {
		return time.Duration(-1)
	}
	return guint64ToDuration(dur)
}

// Offset returns a media specific offset for the buffer data. For video frames, this is the frame
// number of this buffer. For audio samples, this is the offset of the first sample in this buffer.
// For file data or compressed data this is the byte offset of the first byte in this buffer.
func (b *Buffer) Offset() int64 { return int64(b.Instance().offset) }

// OffsetEnd returns the last offset contained in this buffer. It has the same format as Offset.
func (b *Buffer) OffsetEnd() int64 { return int64(b.Instance().offset_end) }

// Map will map the data inside this buffer.
func (b *Buffer) Map() *MapInfo {
	var mapInfo C.GstMapInfo
	C.gst_buffer_map(
		(*C.GstBuffer)(b.Instance()),
		(*C.GstMapInfo)(unsafe.Pointer(&mapInfo)),
		C.GST_MAP_READ,
	)
	return wrapMapInfo(&mapInfo, func() {
		C.gst_buffer_unmap(b.Instance(), (*C.GstMapInfo)(unsafe.Pointer(&mapInfo)))
	})
}

// AddMeta adds metadata for info to the buffer using the parameters in params. The given
// parameters are passed to the MetaInfo's init function, and as such will only work
// for MetaInfo objects created from the go runtime.
//
// Example
//
//     metaInfo := gst.RegisterMeta(glib.TypeFromName("MyObjectType"), "my-meta", 1024, &gst.MetaInfoCallbackFuncs{
//         InitFunc: func(params interface{}, buffer *gst.Buffer) bool {
//             paramStr := params.(string)
//             fmt.Println("Buffer initialized with params:", paramStr)
//             return true
// 	       },
//         FreeFunc: func(buffer *gst.Buffer) {
// 		       fmt.Println("Buffer was destroyed")
// 	       },
//     })
//
//     buf := gst.NewEmptyBuffer()
//     buf.AddMeta(metaInfo, "hello world")
//
//     buf.Unref()
//
//     // > Buffer initialized with params: hello world
//     // > Buffer was destroyed
//
func (b *Buffer) AddMeta(info *MetaInfo, params interface{}) *Meta {
	meta := C.gst_buffer_add_meta(b.Instance(), info.Instance(), (C.gpointer)(gopointer.Save(params)))
	if meta == nil {
		return nil
	}
	return wrapMeta(meta)
}

// GetMeta retrieves the metadata on the buffer for the given api. If none exists
// then nil is returned.
func (b *Buffer) GetMeta(api glib.Type) *Meta {
	meta := C.gst_buffer_get_meta(b.Instance(), C.GType(api))
	return wrapMeta(meta)
}
