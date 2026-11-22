//go:build !gst_use_unstable_api

package gstwebrtc

// #warning "go-gst: the GStreamer WebRTC library is unstable API and may change in future. Build with -tags gst_use_unstable_api to silence this warning."
import "C"
