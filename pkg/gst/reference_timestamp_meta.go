package gst

import (
	"runtime"
	"unsafe"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// GetReference returns the reference member of the struct
func (meta *ReferenceTimestampMeta) GetReference() *Caps {
	ret := UnsafeCapsFromGlibNone(unsafe.Pointer(meta.referenceTimestampMeta.native.reference))
	runtime.KeepAlive(meta)
	return ret
}

// GetTimestamp returns the timestamp member of the struct
func (meta *ReferenceTimestampMeta) GetTimestamp() ClockTime {
	ret := ClockTime(meta.referenceTimestampMeta.native.timestamp)
	runtime.KeepAlive(meta)
	return ret
}

// GetDuration returns the duration member of the struct
func (meta *ReferenceTimestampMeta) GetDuration() ClockTime {
	ret := ClockTime(meta.referenceTimestampMeta.native.duration)
	runtime.KeepAlive(meta)
	return ret
}
