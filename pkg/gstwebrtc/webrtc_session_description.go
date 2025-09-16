package gstwebrtc

import (
	"runtime"
	"unsafe"

	"github.com/go-gst/go-gst/pkg/gstsdp"
)

// Getters for WebRTC session description properties.

// type (GstWebRTCSDPType) - the GstWebRTCSDPType of the description
// sdp (GstSDPMessage *) - the GstSDPMessage of the description

// GetSDPType returns the SDP type of the WebRTC session description.
func (d *WebRTCSessionDescription) GetSDPType() WebRTCSDPType {
	return WebRTCSDPType(d.native._type)
}

// GetSDP returns the SDP message of the WebRTC session description.
func (d *WebRTCSessionDescription) GetSDP() *gstsdp.SDPMessage {
	sdp := gstsdp.UnsafeSDPMessageFromGlibBorrow(unsafe.Pointer(d.native.sdp))

	runtime.AddCleanup(sdp, func(_ struct{}) {
		runtime.KeepAlive(d)
	}, struct{}{})

	return sdp
}
