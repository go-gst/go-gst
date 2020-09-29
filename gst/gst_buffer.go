package gst

/*
#include "gst.go.h"

extern void      goGDestroyNotifyFunc  (gpointer data);
extern gboolean  goBufferMetaForEachCb (GstBuffer * buffer, GstMeta ** meta, gpointer user_data);

void cgoDestroyNotifyFunc (gpointer data) {
	goGDestroyNotifyFunc(data);
}

gboolean cgoBufferMetaForEachCb (GstBuffer * buffer, GstMeta ** meta, gpointer user_data)
{
	return goBufferMetaForEachCb(buffer, meta, user_data);
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

// ParentBufferMeta is a go representation of a GstParentBufferMeta
type ParentBufferMeta struct {
	Parent *Meta
	Buffer *Buffer
}

// AddParentMeta adds a ParentBufferMeta to this buffer that holds a parent reference
// on the given buffer until the it is freed.
func (b *Buffer) AddParentMeta(buf *Buffer) *ParentBufferMeta {
	meta := C.gst_buffer_add_parent_buffer_meta(b.Instance(), buf.Instance())
	return &ParentBufferMeta{
		Parent: wrapMeta(&meta.parent),
		Buffer: wrapBuffer(meta.buffer),
	}
}

// AddProtectionMeta attaches ProtectionMeta to this buffer. The structure contains
// cryptographic information relating to the sample contained in the buffer. This
// function takes ownership of the structure.
func (b *Buffer) AddProtectionMeta(info *Structure) *ProtectionMeta {
	meta := C.gst_buffer_add_protection_meta(b.Instance(), info.Instance())
	return &ProtectionMeta{
		Meta: wrapMeta(&meta.meta),
		Info: wrapStructure(meta.info),
	}
}

// ReferenceTimestampMeta is a go representation of a GstReferenceTimestampMeta.
type ReferenceTimestampMeta struct {
	Parent              *Meta
	Reference           *Caps
	Timestamp, Duration time.Duration
}

// AddReferenceTimestampMeta adds a ReferenceTimestampMeta to this buffer that holds a
// timestamp and optional duration (specify -1 to omit) based on a specific timestamp reference.
//
// See the documentation of GstReferenceTimestampMeta for details.
// https://gstreamer.freedesktop.org/documentation/gstreamer/gstbuffer.html?gi-language=c#GstReferenceTimestampMeta
func (b *Buffer) AddReferenceTimestampMeta(ref *Caps, timestamp, duration time.Duration) *ReferenceTimestampMeta {
	durClockTime := C.GstClockTime(C.GST_CLOCK_TIME_NONE)
	if duration > time.Duration(0) {
		durClockTime = C.GstClockTime(duration.Nanoseconds())
	}
	tsClockTime := C.GstClockTime(timestamp.Nanoseconds())
	meta := C.gst_buffer_add_reference_timestamp_meta(b.Instance(), ref.Instance(), tsClockTime, durClockTime)
	return &ReferenceTimestampMeta{
		Parent:    wrapMeta(&meta.parent),
		Reference: wrapCaps(meta.reference),
		Timestamp: clockTimeToDuration(ClockTime(meta.timestamp)),
		Duration:  clockTimeToDuration(ClockTime(meta.duration)),
	}
}

// Append will append all the memory from the given buffer to this one. The result buffer will
// contain a concatenation of the memory of the two buffers.
func (b *Buffer) Append(buf *Buffer) *Buffer {
	return wrapBuffer(C.gst_buffer_append(b.Instance(), buf.Instance()))
}

// AppendMemory append the memory block to this buffer. This function takes ownership of
// the memory and thus doesn't increase its refcount.
//
// This function is identical to InsertMemory with an index of -1.
func (b *Buffer) AppendMemory(mem *Memory) {
	C.gst_buffer_append_memory(b.Instance(), mem.Instance())
}

// AppendRegion will append size bytes at offset from the given buffer to this one. The result
// buffer will contain a concatenation of the memory of this buffer and the requested region of
// the one provided.
func (b *Buffer) AppendRegion(buf *Buffer, offset, size int64) *Buffer {
	newbuf := C.gst_buffer_append_region(b.Instance(), buf.Instance(), C.gssize(offset), C.gssize(size))
	return wrapBuffer(newbuf)
}

// Copy creates a copy of this buffer. This will only copy the buffer's data to a newly allocated
// Memory if needed (if the type of memory requires it), otherwise the underlying data is just referenced.
// Check DeepCopy if you want to force the data to be copied to newly allocated Memory.
func (b *Buffer) Copy() *Buffer { return wrapBuffer(C.gst_buffer_copy(b.Instance())) }

// DeepCopy creates a copy of the given buffer. This will make a newly allocated copy of the data
// the source buffer contains.
func (b *Buffer) DeepCopy() *Buffer { return wrapBuffer(C.gst_buffer_copy_deep(b.Instance())) }

// CopyInto copies the information from this buffer into the given one. If the given buffer already
// contains memory and flags contains BufferCopyMemory, the memory from this one will be appended to
// that provided.
//
// Flags indicate which fields will be copied. Offset and size dictate from where and how much memory
// is copied. If size is -1 then all data is copied. The function returns true if the copy was successful.
func (b *Buffer) CopyInto(dest *Buffer, flags BufferCopyFlags, offset, size int64) bool {
	ok := C.gst_buffer_copy_into(
		dest.Instance(),
		b.Instance(),
		C.GstBufferCopyFlags(flags),
		C.gsize(offset),
		C.gsize(size),
	)
	return gobool(ok)
}

// CopyRegion creates a sub-buffer from this one at offset and size. This sub-buffer uses the actual memory
// space of the parent buffer. This function will copy the offset and timestamp fields when the offset is 0.
// If not, they will be set to ClockTimeNone and BufferOffsetNone.
//
// If offset equals 0 and size equals the total size of buffer, the duration and offset end fields are also
// copied. If not they will be set to ClockTimeNone and BufferOffsetNone.
func (b *Buffer) CopyRegion(flags BufferCopyFlags, offset, size int64) *Buffer {
	newbuf := C.gst_buffer_copy_region(
		b.Instance(),
		C.GstBufferCopyFlags(flags),
		C.gsize(offset),
		C.gsize(size),
	)
	return wrapBuffer(newbuf)
}

// Extract extracts size bytes starting from offset in this buffer. The data extracted may be lower
// than the actual size if the buffer did not contain enough data.
func (b *Buffer) Extract(offset, size int64) []byte {
	dest := C.malloc(C.sizeof_char * C.ulong(size))
	defer C.free(dest)
	C.gst_buffer_extract(b.Instance(), C.gsize(offset), (C.gpointer)(unsafe.Pointer(dest)), C.gsize(size))
	return C.GoBytes(dest, C.int(size))
}

// Fill adds the given byte slice to the buffer at the given offset. The return value reflects the amount
// of data added to the buffer.
func (b *Buffer) Fill(offset int64, data []byte) int64 {
	str := string(data)
	cStr := C.CString(str)
	gsize := C.gst_buffer_fill(b.Instance(), C.gsize(offset), (C.gconstpointer)(unsafe.Pointer(cStr)), C.gsize(len(str)))
	return int64(gsize)
}

// FindMemory looks for the memory blocks that span size bytes starting from offset in buffer. Size can be -1
// to retrieve all the memory blocks.
//
// Index will contain the index of the first memory block where the byte for offset can be found and length
// contains the number of memory blocks containing the size remaining bytes. Skip contains the number of bytes
// to skip in the memory block at index to get to the byte for offset. All values will be 0 if the memory blocks
// could not be read.
func (b *Buffer) FindMemory(offset, size int64) (index, length uint, skip int64) {
	var gindex, glength C.uint
	var gskip C.gsize
	ok := C.gst_buffer_find_memory(b.Instance(), C.gsize(offset), C.gsize(size), &gindex, &glength, &gskip)
	if !gobool(ok) {
		return
	}
	return uint(gindex), uint(glength), int64(gskip)
}

// ForEachMeta calls the given function for each Meta in this buffer.
//
// The function can modify the passed meta pointer or its contents. The return value defines if this function continues
// or if the remaining metadata items in the buffer should be skipped.
func (b *Buffer) ForEachMeta(f func(meta *Meta) bool) bool {
	fPtr := gopointer.Save(f)
	defer gopointer.Unref(fPtr)
	return gobool(C.gst_buffer_foreach_meta(
		b.Instance(),
		C.GstBufferForeachMetaFunc(C.cgoBufferMetaForEachCb),
		(C.gpointer)(unsafe.Pointer(fPtr)),
	))
}
