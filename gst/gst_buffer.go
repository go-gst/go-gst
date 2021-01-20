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

gboolean isBuffer (GstBuffer * buffer) { return GST_IS_BUFFER(buffer); }
*/
import "C"

import (
	"bytes"
	"io"
	"io/ioutil"
	"runtime"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// GetMaxBufferMemory returns the maximum amount of memory a buffer can hold.
func GetMaxBufferMemory() uint64 { return uint64(C.gst_buffer_get_max_memory()) }

// Buffer is a go representation of a GstBuffer.
type Buffer struct {
	ptr     *C.GstBuffer
	mapInfo *MapInfo
}

// FromGstBufferUnsafeNone wraps the given buffer, sinking any floating references, and places
// a finalizer on the wrapped Buffer.
func FromGstBufferUnsafeNone(buf unsafe.Pointer) *Buffer {
	wrapped := ToGstBuffer(buf)
	wrapped.Ref()
	runtime.SetFinalizer(wrapped, (*Buffer).Unref)
	return wrapped
}

// FromGstBufferUnsafeFull wraps the given buffer without taking an additional reference.
func FromGstBufferUnsafeFull(buf unsafe.Pointer) *Buffer {
	wrapped := ToGstBuffer(buf)
	runtime.SetFinalizer(wrapped, (*Buffer).Unref)
	return wrapped
}

// ToGstBuffer converts the given pointer into a Buffer without affecting the ref count or
// placing finalizers.
func ToGstBuffer(buf unsafe.Pointer) *Buffer {
	return wrapBuffer((*C.GstBuffer)(buf))
}

// NewEmptyBuffer returns a new empty buffer.
func NewEmptyBuffer() *Buffer {
	return FromGstBufferUnsafeFull(unsafe.Pointer(C.gst_buffer_new()))
}

// NewBufferWithSize is a convenience wrapped for NewBufferrAllocate with the default allocator
// and parameters.
func NewBufferWithSize(size int64) *Buffer {
	return NewBufferAllocate(nil, nil, size)
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
	var gstparams *C.GstAllocationParams
	if params != nil {
		gstparams = params.Instance()
	}
	buf := C.gst_buffer_new_allocate(gstalloc, C.gsize(size), gstparams)
	if buf == nil {
		return nil
	}
	return FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

// NewBufferFromBytes returns a new buffer from the given byte slice.
func NewBufferFromBytes(b []byte) *Buffer {
	gbytes := C.g_bytes_new((C.gconstpointer)(unsafe.Pointer(&b[0])), C.gsize(len(b)))
	defer C.g_bytes_unref(gbytes)
	return FromGstBufferUnsafeFull(unsafe.Pointer(C.gst_buffer_new_wrapped_bytes(gbytes)))
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
//   // Example
//
//   buf := gst.NewBufferFull(0, []byte("hello-world"), 1024, 0, 1024, func() {
//       fmt.Println("buffer was destroyed")
//   })
//   if buf != nil {
//       buf.Unref()
//   }
//
//  // > buffer was destroyed
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
	return FromGstBufferUnsafeFull(unsafe.Pointer(buf))
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
	mapInfo := b.Map(MapRead)
	if mapInfo == nil {
		return nil
	}
	defer b.Unmap()
	return mapInfo.Bytes()
}

// PresentationTimestamp returns the presentation timestamp of the buffer, or a negative duration
// if not known or relevant. This value contains the timestamp when the media should be
// presented to the user.
func (b *Buffer) PresentationTimestamp() time.Duration {
	pts := b.Instance().pts
	if pts == gstClockTimeNone {
		return ClockTimeNone
	}
	return time.Duration(pts)
}

// SetPresentationTimestamp sets the presentation timestamp on the buffer.
func (b *Buffer) SetPresentationTimestamp(dur time.Duration) {
	ins := b.Instance()
	ins.pts = C.GstClockTime(dur.Nanoseconds())
}

// DecodingTimestamp returns the decoding timestamp of the buffer, or a negative duration if not known
// or relevant. This value contains the timestamp when the media should be processed.
func (b *Buffer) DecodingTimestamp() time.Duration {
	dts := b.Instance().dts
	if dts == gstClockTimeNone {
		return ClockTimeNone
	}
	return time.Duration(dts)
}

// Duration returns the length of the data inside this buffer, or a negative duration if not known
// or relevant.
func (b *Buffer) Duration() time.Duration {
	dur := b.Instance().duration
	if dur == gstClockTimeNone {
		return ClockTimeNone
	}
	return time.Duration(dur)
}

// SetDuration sets the duration on the buffer.
func (b *Buffer) SetDuration(dur time.Duration) {
	ins := b.Instance()
	ins.duration = C.GstClockTime(dur.Nanoseconds())
}

// Offset returns a media specific offset for the buffer data. For video frames, this is the frame
// number of this buffer. For audio samples, this is the offset of the first sample in this buffer.
// For file data or compressed data this is the byte offset of the first byte in this buffer.
func (b *Buffer) Offset() int64 { return int64(b.Instance().offset) }

// OffsetEnd returns the last offset contained in this buffer. It has the same format as Offset.
func (b *Buffer) OffsetEnd() int64 { return int64(b.Instance().offset_end) }

// AddMeta adds metadata for info to the buffer using the parameters in params. The given
// parameters are passed to the MetaInfo's init function, and as such will only work
// for MetaInfo objects created from the go runtime.
//
//   // Example
//
//   metaInfo := gst.RegisterMeta(glib.TypeFromName("MyObjectType"), "my-meta", 1024, &gst.MetaInfoCallbackFuncs{
//       InitFunc: func(params interface{}, buffer *gst.Buffer) bool {
//           paramStr := params.(string)
//           fmt.Println("Buffer initialized with params:", paramStr)
//           return true
//       },
//       FreeFunc: func(buffer *gst.Buffer) {
//           fmt.Println("Buffer was destroyed")
//       },
//   })
//
//   buf := gst.NewEmptyBuffer()
//   buf.AddMeta(metaInfo, "hello world")
//
//   buf.Unref()
//
//   // > Buffer initialized with params: hello world
//   // > Buffer was destroyed
//
func (b *Buffer) AddMeta(info *MetaInfo, params interface{}) *Meta {
	meta := C.gst_buffer_add_meta(b.Instance(), info.Instance(), (C.gpointer)(gopointer.Save(params)))
	if meta == nil {
		return nil
	}
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
	durClockTime := C.GstClockTime(ClockTimeNone)
	if duration > time.Duration(0) {
		durClockTime = C.GstClockTime(duration.Nanoseconds())
	}
	tsClockTime := C.GstClockTime(timestamp.Nanoseconds())
	meta := C.gst_buffer_add_reference_timestamp_meta(b.Instance(), ref.Instance(), tsClockTime, durClockTime)
	if meta == nil {
		return nil
	}
	return &ReferenceTimestampMeta{
		Parent:    wrapMeta(&meta.parent),
		Reference: wrapCaps(meta.reference),
		Timestamp: time.Duration(meta.timestamp),
		Duration:  time.Duration(meta.duration),
	}
}

// Append will append all the memory from the given buffer to this one. The result buffer will
// contain a concatenation of the memory of the two buffers.
func (b *Buffer) Append(buf *Buffer) *Buffer {
	return FromGstBufferUnsafeFull(unsafe.Pointer(C.gst_buffer_append(b.Ref().Instance(), buf.Ref().Instance())))
}

// AppendMemory append the memory block to this buffer. This function takes ownership of
// the memory and thus doesn't increase its refcount.
//
// This function is identical to InsertMemory with an index of -1.
func (b *Buffer) AppendMemory(mem *Memory) {
	C.gst_buffer_append_memory(b.Instance(), mem.Ref().Instance())
}

// AppendRegion will append size bytes at offset from the given buffer to this one. The result
// buffer will contain a concatenation of the memory of this buffer and the requested region of
// the one provided.
func (b *Buffer) AppendRegion(buf *Buffer, offset, size int64) *Buffer {
	newbuf := C.gst_buffer_append_region(b.Ref().Instance(), buf.Ref().Instance(), C.gssize(offset), C.gssize(size))
	return FromGstBufferUnsafeFull(unsafe.Pointer(newbuf))
}

// Copy creates a copy of this buffer. This will only copy the buffer's data to a newly allocated
// Memory if needed (if the type of memory requires it), otherwise the underlying data is just referenced.
// Check DeepCopy if you want to force the data to be copied to newly allocated Memory.
func (b *Buffer) Copy() *Buffer {
	buf := C.gst_buffer_copy(b.Instance())
	return FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

// DeepCopy creates a copy of the given buffer. This will make a newly allocated copy of the data
// the source buffer contains.
func (b *Buffer) DeepCopy() *Buffer {
	buf := C.gst_buffer_copy_deep(b.Instance())
	return FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

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
	dest := C.malloc(C.sizeof_char * C.gsize(size))
	defer C.free(dest)
	C.gst_buffer_extract(b.Instance(), C.gsize(offset), (C.gpointer)(unsafe.Pointer(dest)), C.gsize(size))
	return C.GoBytes(dest, C.int(size))
}

// FillBytes adds the given byte slice to the buffer at the given offset. The return value reflects the amount
// of data added to the buffer.
func (b *Buffer) FillBytes(offset int64, data []byte) int64 {
	gbytes := C.g_bytes_new((C.gconstpointer)(unsafe.Pointer(&data[0])), C.gsize(len(data)))
	defer C.g_bytes_unref(gbytes)
	var size C.gsize
	gdata := C.g_bytes_get_data(gbytes, &size)
	if gdata == nil {
		return 0
	}
	return int64(C.gst_buffer_fill(b.Instance(), C.gsize(offset), gdata, size))
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

// GetAllMemory retrieves all the memory inside this buffer.
func (b *Buffer) GetAllMemory() *Memory {
	mem := C.gst_buffer_get_all_memory(b.Instance())
	if mem == nil {
		return nil
	}
	return FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// GetFlags returns the flags on this buffer.
func (b *Buffer) GetFlags() BufferFlags {
	return BufferFlags(C.gst_buffer_get_flags(b.Instance()))
}

// GetMemory retrieves the memory block at the given index in the buffer.
func (b *Buffer) GetMemory(idx uint) *Memory {
	mem := C.gst_buffer_get_memory(b.Instance(), C.guint(idx))
	if mem == nil {
		return nil
	}
	return FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// GetMemoryRange retrieves length memory blocks in buffer starting at idx. The memory blocks
// will be merged into one large Memory. If length is -1, all memory starting from idx is merged.
func (b *Buffer) GetMemoryRange(idx uint, length int) *Memory {
	mem := C.gst_buffer_get_memory_range(b.Instance(), C.guint(idx), C.gint(length))
	if mem == nil {
		return nil
	}
	return FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// GetMeta retrieves the metadata for the given api on buffer. When there is no such metadata,
// nil is returned. If multiple metadata with the given api are attached to this buffer only the
// first one is returned. To handle multiple metadata with a given API use ForEachMeta instead
// and check the type.
func (b *Buffer) GetMeta(api glib.Type) *Meta {
	meta := C.gst_buffer_get_meta(b.Instance(), C.GType(api))
	if meta == nil {
		return nil
	}
	return wrapMeta(meta)
}

// GetNumMetas returns the number of metas for the given api type on the buffer.
func (b *Buffer) GetNumMetas(api glib.Type) uint {
	return uint(C.gst_buffer_get_n_meta(b.Instance(), C.GType(api)))
}

// GetReferenceTimestampMeta finds the first ReferenceTimestampMeta on the buffer that conforms to
// reference. Conformance is tested by checking if the meta's reference is a subset of reference.
//
// Buffers can contain multiple ReferenceTimestampMeta metadata items.
func (b *Buffer) GetReferenceTimestampMeta(caps *Caps) *ReferenceTimestampMeta {
	var meta *C.GstReferenceTimestampMeta
	if caps == nil {
		meta = C.gst_buffer_get_reference_timestamp_meta(b.Instance(), nil)
	} else {
		meta = C.gst_buffer_get_reference_timestamp_meta(b.Instance(), caps.Instance())
	}
	if meta == nil {
		return nil
	}
	refMeta := &ReferenceTimestampMeta{
		Parent:    wrapMeta(&meta.parent),
		Timestamp: time.Duration(meta.timestamp),
		Duration:  time.Duration(meta.duration),
	}
	if meta.reference != nil {
		refMeta.Reference = wrapCaps(meta.reference)
	}
	return refMeta
}

// GetSize retrieves the number of Memory blocks in the bffer.
func (b *Buffer) GetSize() int64 {
	return int64(C.gst_buffer_get_size(b.Instance()))
}

// GetSizes will retrieve the size of the buffer, the offset of the first memory block in the buffer,
// and the sum of the size of the buffer, the offset, and any padding. These values can be used to
// resize the buffer with Resize.
func (b *Buffer) GetSizes() (size, offset, maxsize int64) {
	var goffset, gmaxsize C.gsize
	gsize := C.gst_buffer_get_sizes(b.Instance(), &goffset, &gmaxsize)
	return int64(gsize), int64(goffset), int64(gmaxsize)
}

// GetSizesRange will get the total size of length memory blocks stating from idx in buffer.
//
// Offset will contain the offset of the data in the memory block in buffer at idx and maxsize will
// contain the sum of the size and offset and the amount of extra padding on the memory block at
// idx + length -1. Offset and maxsize can be used to resize the buffer memory blocks with ResizeRange.
func (b *Buffer) GetSizesRange(idx uint, length int) (offset, maxsize int64) {
	var goffset, gmaxsize C.gsize
	C.gst_buffer_get_sizes_range(b.Instance(), C.guint(idx), C.gint(length), &goffset, &gmaxsize)
	return int64(goffset), int64(gmaxsize)
}

// HasFlags returns true if this Buffer has the given BufferFlags.
func (b *Buffer) HasFlags(flags BufferFlags) bool {
	return gobool(C.gst_buffer_has_flags(b.Instance(), C.GstBufferFlags(flags)))
}

// InsertMemory insert the memory block to the buffer at idx. This function takes ownership of the Memory
// and thus doesn't increase its refcount.
//
// Only the value from GetMaxBufferMemory can be added to a buffer. If more memory is added, existing memory
// blocks will automatically be merged to make room for the new memory.
func (b *Buffer) InsertMemory(mem *Memory, idx int) {
	C.gst_buffer_insert_memory(b.Instance(), C.gint(idx), mem.Ref().Instance())
}

// IsAllMemoryWritable checks if all memory blocks in buffer are writable.
//
// Note that this function does not check if buffer is writable, use IsWritable to check that if needed.
func (b *Buffer) IsAllMemoryWritable() bool {
	return gobool(C.gst_buffer_is_all_memory_writable(b.Instance()))
}

// IsMemoryRangeWritable checks if length memory blocks in the buffer starting from idx are writable.
//
// Length can be -1 to check all the memory blocks after idx.
//
// Note that this function does not check if buffer is writable, use IsWritable to check that if needed.
func (b *Buffer) IsMemoryRangeWritable(idx uint, length int) bool {
	return gobool(C.gst_buffer_is_memory_range_writable(b.Instance(), C.guint(idx), C.gint(length)))
}

// IsWritable returns true if this buffer is writable.
func (b *Buffer) IsWritable() bool {
	return gobool(C.bufferIsWritable(b.Instance()))
}

// MakeWritable returns a writable copy of this buffer. If the source buffer is already writable,
// this will simply return the same buffer.
//
// Use this function to ensure that a buffer can be safely modified before making changes to it,
// including changing the metadata such as PTS/DTS.
//
// If the reference count of the source buffer buf is exactly one, the caller is the sole owner and
// this function will return the buffer object unchanged.
//
// If there is more than one reference on the object, a copy will be made using gst_buffer_copy. The passed-in buf
// will be unreffed in that case, and the caller will now own a reference to the new returned buffer object. Note that
// this just copies the buffer structure itself, the underlying memory is not copied if it can be shared amongst
// multiple buffers.
//
// In short, this function unrefs the buf in the argument and refs the buffer that it returns. Don't access the argument
// after calling this function unless you have an additional reference to it.
func (b *Buffer) MakeWritable() *Buffer {
	return wrapBuffer(C.makeBufferWritable(b.Instance()))
}

// IterateMeta retrieves the next Meta after the given one. If state points to nil, the first Meta is returned.
func (b *Buffer) IterateMeta(meta *Meta) *Meta {
	ptr := unsafe.Pointer(meta.Instance())
	return wrapMeta(C.gst_buffer_iterate_meta(b.Instance(), (*C.gpointer)(&ptr)))
}

// IterateMetaFiltered is similar to IterateMeta except it will filter on the api type.
func (b *Buffer) IterateMetaFiltered(meta *Meta, apiType glib.Type) *Meta {
	ptr := unsafe.Pointer(meta.Instance())
	return wrapMeta(C.gst_buffer_iterate_meta_filtered(b.Instance(), (*C.gpointer)(&ptr), C.GType(apiType)))
}

// Map will map the data inside this buffer. This function can return nil if the memory is not read or writable.
// It is safe to call this function multiple times on a single Buffer, however it will retain the flags
// used when mapping the first time. To change between read and write access first unmap and then remap the
// buffer with the appropriate flags, or map initially with both read/write access.
//
// Unmap the Buffer after usage.
func (b *Buffer) Map(flags MapFlags) *MapInfo {
	if b.mapInfo != nil {
		return b.mapInfo
	}
	var mapInfo C.GstMapInfo
	C.gst_buffer_map(
		(*C.GstBuffer)(b.Instance()),
		(*C.GstMapInfo)(&mapInfo),
		C.GstMapFlags(flags),
	)
	b.mapInfo = wrapMapInfo((*C.GstMapInfo)(&mapInfo))
	return b.mapInfo
}

// Unmap will unmap the data inside this memory. Use this after calling Map on the buffer.
func (b *Buffer) Unmap() {
	if b.mapInfo == nil {
		return
	}
	C.gst_buffer_unmap(b.Instance(), (*C.GstMapInfo)(b.mapInfo.Instance()))
	b.mapInfo = nil
}

// MapRange maps the info of length merged memory blocks starting at idx in buffer.
// When length is -1, all memory blocks starting from idx are merged and mapped.
//
// Flags describe the desired access of the memory. When flags is MapWrite, buffer should be writable (as returned from IsWritable).
//
// When the buffer is writable but the memory isn't, a writable copy will automatically be
// created and returned. The readonly copy of the buffer memory will then also be replaced with this writable copy.
//
// Unmap the Buffer after usage.
func (b *Buffer) MapRange(idx uint, length int, flags MapFlags) *MapInfo {
	if b.mapInfo != nil {
		return b.mapInfo
	}
	var mapInfo C.GstMapInfo
	C.gst_buffer_map_range(
		(*C.GstBuffer)(b.Instance()),
		C.guint(idx), C.gint(length),
		(*C.GstMapInfo)(&mapInfo),
		C.GstMapFlags(flags),
	)
	b.mapInfo = wrapMapInfo((*C.GstMapInfo)(&mapInfo))
	return b.mapInfo
}

// Memset fills buf with size bytes with val starting from offset. It returns the
// size written to the buffer.
func (b *Buffer) Memset(offset int64, val uint8, size int64) int64 {
	return int64(C.gst_buffer_memset(b.Instance(), C.gsize(offset), C.guint8(val), C.gsize(size)))
}

// NumMemoryBlocks returns the number of memory blocks this buffer has.
func (b *Buffer) NumMemoryBlocks() uint { return uint(C.gst_buffer_n_memory(b.Instance())) }

// PeekMemory gets the memory block at idx in buffer. The memory block stays valid until the
// memory block is removed, replaced, or merged. Typically with any call that modifies the memory in buffer.
func (b *Buffer) PeekMemory(idx uint) *Memory {
	mem := C.gst_buffer_peek_memory(b.Instance(), C.guint(idx))
	if mem == nil {
		return nil
	}
	return FromGstMemoryUnsafeNone(unsafe.Pointer(mem))
}

// PrependMemory prepends the memory block mem to this buffer. This function takes ownership of
// mem and thus doesn't increase its refcount.
//
// This function is identical to InsertMemory with an index of 0.
func (b *Buffer) PrependMemory(mem *Memory) {
	C.gst_buffer_prepend_memory(b.Instance(), mem.Ref().Instance())
}

// RemoveAllMemory removes all memory blocks in the buffer.
func (b *Buffer) RemoveAllMemory() { C.gst_buffer_remove_all_memory(b.Instance()) }

// RemoveMemoryAt removes the memory block at the given index.
func (b *Buffer) RemoveMemoryAt(idx uint) { C.gst_buffer_remove_memory(b.Instance(), C.guint(idx)) }

// RemoveMemoryRange removes length memory blocks in buffer starting from idx.
//
// Length can be -1, in which case all memory starting from idx is removed.
func (b *Buffer) RemoveMemoryRange(idx uint, length int) {
	C.gst_buffer_remove_memory_range(b.Instance(), C.guint(idx), C.gint(length))
}

// RemoveMeta removes the given metadata from the buffer.
func (b *Buffer) RemoveMeta(meta *Meta) bool {
	return gobool(C.gst_buffer_remove_meta(b.Instance(), meta.Instance()))
}

// ReplaceAllMemory replaces all the memory in this buffer with that provided.
func (b *Buffer) ReplaceAllMemory(mem *Memory) {
	C.gst_buffer_replace_all_memory(b.Instance(), mem.Ref().Instance())
}

// ReplaceMemory replaces the memory at the given index with the given memory.
func (b *Buffer) ReplaceMemory(mem *Memory, idx uint) {
	C.gst_buffer_replace_memory(b.Instance(), C.guint(idx), mem.Ref().Instance())
}

// ReplaceMemoryRange replaces length memory blocks in the buffer starting at idx with
// the given memory.
//
// If length is -1, all memory starting from idx will be removed and replaced.
//
// The buffer should be writable.
func (b *Buffer) ReplaceMemoryRange(idx uint, length int, mem *Memory) {
	C.gst_buffer_replace_memory_range(b.Instance(), C.guint(idx), C.gint(length), mem.Ref().Instance())
}

// Resize sets the offset and total size of the memory blocks in this buffer.
func (b *Buffer) Resize(offset, size int64) {
	C.gst_buffer_resize(b.Instance(), C.gssize(offset), C.gssize(size))
}

// ResizeRange sets the total size of the length memory blocks starting at idx in this buffer.
func (b *Buffer) ResizeRange(idx uint, length int, offset, size int64) bool {
	return gobool(C.gst_buffer_resize_range(
		b.Instance(),
		C.guint(idx),
		C.gint(length),
		C.gssize(offset),
		C.gssize(size),
	))
}

// SetFlags sets one or more buffer flags on the buffer.
func (b *Buffer) SetFlags(flags BufferFlags) bool {
	return gobool(C.gst_buffer_set_flags(b.Instance(), C.GstBufferFlags(flags)))
}

// SetSize sets the total size of the memory blocks in buffer.
func (b *Buffer) SetSize(size int64) {
	C.gst_buffer_set_size(b.Instance(), C.gssize(size))
}

// UnsetFlags removes one or more flags from the buffer.
func (b *Buffer) UnsetFlags(flags BufferFlags) bool {
	return gobool(C.gst_buffer_unset_flags(b.Instance(), C.GstBufferFlags(flags)))
}
