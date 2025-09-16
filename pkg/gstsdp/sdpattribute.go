package gstsdp

// #cgo pkg-config: gstreamer-sdp-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/sdp/sdp.h>
import "C"

// getters for the fields of GstSDPAttribute:

// GetKey returns the key of the SDPAttribute
func (a *SDPAttribute) GetKey() string {
	return C.GoString(a.native.key)
}

// GetValue returns the value of the SDPAttribute
func (a *SDPAttribute) GetValue() string {
	return C.GoString(a.native.value)
}
