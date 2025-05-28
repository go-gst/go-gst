package allocators

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-gst/gst"
)

// FdAllocator is a go representation of a GstFdAllocator
type FdAllocator gst.Allocator

// FdAllocator returns a new GstFdAllocator.
func NewFdAllocator() *FdAllocator {
	return (*FdAllocator)(gst.FromGstAllocatorUnsafeFull(unsafe.Pointer(C.gst_fd_allocator_new())))
}

// Instance returns the underlying GstAllocator instance.
func (a *FdAllocator) Instance() *C.GstAllocator {
	return (*C.GstAllocator)(unsafe.Pointer(((*gst.Allocator)(a).Instance())))
}

// AllocDmaBufWithFlags is used to allocate a new memory block with memory that is at least size big.
func (a *FdAllocator) AllocFd(fd int, size int64, flags FdMemoryFlags) *gst.Memory {
	mem := C.gst_fd_allocator_alloc(a.Instance(), C.gint(fd), C.gsize(size), C.GstFdMemoryFlags(flags))
	return gst.FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// DmaBufAllocator is a go representation of a GstDmaBufAllocator
type DmaBufAllocator gst.Allocator

// DmaBufAllocator returns a new GstDmaBufAllocator.
func NewDmaBufAllocator() *DmaBufAllocator {
	return (*DmaBufAllocator)(gst.FromGstAllocatorUnsafeFull(unsafe.Pointer(C.gst_dmabuf_allocator_new())))
}

// Instance returns the underlying GstAllocator instance.
func (a *DmaBufAllocator) Instance() *C.GstAllocator {
	return (*C.GstAllocator)(unsafe.Pointer(((*gst.Allocator)(a).Instance())))
}

// AllocDmaBuf is used to allocate a new memory block with memory that is at least size big.
func (a *DmaBufAllocator) AllocDmaBuf(fd int, size int64) *gst.Memory {
	mem := C.gst_dmabuf_allocator_alloc(a.Instance(), C.gint(fd), C.gsize(size))
	return gst.FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}

// AllocDmaBufWithFlags is used to allocate a new memory block with memory that is at least size big.
func (a *DmaBufAllocator) AllocDmaBufWithFlags(fd int, size int64, flags FdMemoryFlags) *gst.Memory {
	mem := C.gst_dmabuf_allocator_alloc_with_flags(a.Instance(), C.gint(fd), C.gsize(size), C.GstFdMemoryFlags(flags))
	return gst.FromGstMemoryUnsafeFull(unsafe.Pointer(mem))
}
