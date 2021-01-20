package gst

/*
#include "gst.go.h"

extern gboolean goBufferListForEachCb (GstBuffer ** buffer, guint idx, gpointer user_data);

gboolean cgoBufferListForEachCb (GstBuffer ** buffer, guint idx, gpointer user_data)
{
	return goBufferListForEachCb(buffer, idx, user_data);
}
*/
import "C"
import (
	"runtime"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

// BufferList is a go wrapper around a GstBufferList for grouping Buffers
type BufferList struct {
	ptr *C.GstBufferList
}

// NewBufferList returns a new BufferList. The given slice can be nil and the returned
// buffer list will be empty.
func NewBufferList(buffers []*Buffer) *BufferList {
	bufList := FromGstBufferListUnsafeFull(unsafe.Pointer(C.gst_buffer_list_new()))
	if buffers == nil {
		return bufList
	}
	for idx, buf := range buffers {
		bufList.Insert(idx, buf)
	}
	return bufList
}

// NewBufferListSized creates a new BufferList with the given size.
func NewBufferListSized(size uint) *BufferList {
	return FromGstBufferListUnsafeFull(unsafe.Pointer(C.gst_buffer_list_new_sized(C.guint(size))))
}

// FromGstBufferListUnsafeNone is used for returns from transfer-none methods.
func FromGstBufferListUnsafeNone(buf unsafe.Pointer) *BufferList {
	wrapped := wrapBufferList((*C.GstBufferList)(buf))
	wrapped.Ref()
	runtime.SetFinalizer(wrapped, (*BufferList).Unref)
	return wrapped
}

// FromGstBufferListUnsafeFull wraps the given buffer without taking an additional reference.
func FromGstBufferListUnsafeFull(buf unsafe.Pointer) *BufferList {
	wrapped := wrapBufferList((*C.GstBufferList)(buf))
	runtime.SetFinalizer(wrapped, (*BufferList).Unref)
	return wrapped
}

// ToGstBufferList converts the given pointer into a BufferList without affecting the ref count or
// placing finalizers.
func ToGstBufferList(buf unsafe.Pointer) *BufferList {
	return wrapBufferList((*C.GstBufferList)(buf))
}

// Instance returns the underlying GstBufferList.
func (b *BufferList) Instance() *C.GstBufferList { return C.toGstBufferList(unsafe.Pointer(b.ptr)) }

// CalculateSize calculates the size of the data contained in this buffer list by adding the size of all buffers.
func (b *BufferList) CalculateSize() int64 {
	return int64(C.gst_buffer_list_calculate_size(b.Instance()))
}

// Copy creates a shallow copy of the given buffer list. This will make a newly allocated copy of the
// source list with copies of buffer pointers. The refcount of buffers pointed to will be increased by one.
func (b *BufferList) Copy() *BufferList {
	return FromGstBufferListUnsafeFull(unsafe.Pointer(C.gst_buffer_list_copy(b.Instance())))
}

// DeepCopy creates a copy of the given buffer list. This will make a newly allocated copy of each buffer
// that the source buffer list contains.
func (b *BufferList) DeepCopy() *BufferList {
	return FromGstBufferListUnsafeFull(unsafe.Pointer(C.gst_buffer_list_copy_deep(b.Instance())))
}

// IsWritable returns true if this BufferList is writable.
func (b *BufferList) IsWritable() bool {
	return gobool(C.bufferListIsWritable(b.Instance()))
}

// MakeWritable makes a writable buffer list from this one. If the source buffer list is already writable,
// this will simply return the same buffer list. A copy will otherwise be made using Copy.
func (b *BufferList) MakeWritable() *BufferList {
	return FromGstBufferListUnsafeFull(unsafe.Pointer(C.makeBufferListWritable(b.Instance())))
}

// ForEach calls the given function for each buffer in list.
//
// The function can modify the passed buffer pointer or its contents. The return value defines if this
// function returns or if the remaining buffers in the list should be skipped.
func (b *BufferList) ForEach(f func(buf *Buffer, idx uint) bool) {
	fPtr := gopointer.Save(f)
	defer gopointer.Unref(fPtr)
	C.gst_buffer_list_foreach(
		b.Instance(),
		C.GstBufferListFunc(C.cgoBufferListForEachCb),
		(C.gpointer)(unsafe.Pointer(fPtr)),
	)
}

// GetBufferAt gets the buffer at idx.
//
// You must make sure that idx does not exceed the number of buffers available.
func (b *BufferList) GetBufferAt(idx uint) *Buffer {
	return FromGstBufferUnsafeNone(unsafe.Pointer(C.gst_buffer_list_get(b.Instance(), C.guint(idx))))
}

// GetWritableBufferAt gets the buffer at idx, ensuring it is a writable buffer.
//
// You must make sure that idx does not exceed the number of buffers available.
func (b *BufferList) GetWritableBufferAt(idx uint) *Buffer {
	return FromGstBufferUnsafeNone(unsafe.Pointer(C.gst_buffer_list_get_writable(b.Instance(), C.guint(idx))))
}

// Insert inserts a buffer at idx in the list. Other buffers are moved to make room for this new buffer.
//
// A -1 value for idx will append the buffer at the end.
func (b *BufferList) Insert(idx int, buf *Buffer) {
	C.gst_buffer_list_insert(b.Instance(), C.gint(idx), buf.Ref().Instance())
}

// Length returns the number of buffers in the list.
func (b *BufferList) Length() uint {
	return uint(C.gst_buffer_list_length(b.Instance()))
}

// Ref increases the refcount of the given buffer list by one.
//
// Note that the refcount affects the writability of list and its data, see MakeWritable.
// It is important to note that keeping additional references to GstBufferList instances
// can potentially increase the number of memcpy operations in a pipeline.
func (b *BufferList) Ref() *BufferList {
	C.gst_buffer_list_ref(b.Instance())
	return b
}

// Remove removes length buffers from the list starting at index. All following buffers
// are moved to close the gap.
func (b *BufferList) Remove(idx, length uint) {
	C.gst_buffer_list_remove(b.Instance(), C.guint(idx), C.guint(length))
}

// Unref decreases the refcount of the buffer list. If the refcount reaches 0, the buffer
// list will be freed.
func (b *BufferList) Unref() {
	C.gst_buffer_list_unref(b.Instance())
}
