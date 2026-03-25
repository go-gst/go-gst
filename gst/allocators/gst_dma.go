package allocators

/*
#include <gst/allocators/gstdmabuf.h>
*/
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-gst/gst"
)


// IsDMABuffer returns true if a Memory is DMA buffer memory.
func IsDMABuffer(m *gst.Memory) bool {
	return int(C.gst_is_dmabuf_memory((*C.GstMemory)(unsafe.Pointer(m.Instance())))) > 0
}

// GetDMABufferFD returns a DMA buffer file descriptor for a Memory.
// If the memory is not a DMA buffer, returns -1.
// The file descriptor is still owned by the Memory.
// Use dup() if you intend to use it beyond the lifetime of the Memory.
func GetDMABufferFD(m *gst.Memory, idx uint) int {
	if !IsDMABuffer(m) {
		return -1
	}
	fd := int(C.gst_dmabuf_memory_get_fd((*C.GstMemory)(unsafe.Pointer(m.Instance()))))
	return fd
}
