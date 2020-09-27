package gst

// #include "gst.go.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// MapInfo is a go representation of a GstMapInfo.
type MapInfo struct {
	ptr       *C.GstMapInfo
	unmapFunc func()
	Memory    unsafe.Pointer // A pointer to the GstMemory object
	Flags     MapFlags
	Data      unsafe.Pointer // A pointer to the actual data
	Size      int64
	MaxSize   int64
}

// MapBuffer will retrieve the map info for the given buffer. Unmap after usage.
func MapBuffer(buffer *Buffer) *MapInfo {
	var mapInfo C.GstMapInfo
	C.gst_buffer_map(
		(*C.GstBuffer)(buffer.Instance()),
		(*C.GstMapInfo)(unsafe.Pointer(&mapInfo)),
		C.GST_MAP_READ,
	)
	return wrapMapInfo(&mapInfo, func() {
		C.gst_buffer_unmap(buffer.Instance(), (*C.GstMapInfo)(unsafe.Pointer(&mapInfo)))
	})
}

// Unmap will unmap the MapInfo.
func (m *MapInfo) Unmap() {
	if m.unmapFunc == nil {
		fmt.Println("GO-GST-WARNING: Called Unmap() on unwrapped MapInfo")
	}
	m.unmapFunc()
}

func wrapMapInfo(mapInfo *C.GstMapInfo, unmapFunc func()) *MapInfo {
	return &MapInfo{
		ptr:       mapInfo,
		unmapFunc: unmapFunc,
		Memory:    unsafe.Pointer(mapInfo.memory),
		Flags:     MapFlags(mapInfo.flags),
		Data:      unsafe.Pointer(mapInfo.data),
		Size:      int64(mapInfo.size),
		MaxSize:   int64(mapInfo.maxsize),
	}
}
