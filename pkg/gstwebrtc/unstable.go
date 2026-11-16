package gstwebrtc

// The upstream WebRTC headers emit a compiler warning unless GST_USE_UNSTABLE_API is defined. That
// warning tells the user to define a C macro, which is not actionable for users of the bindings.
// We always define it (cgo applies the CPPFLAGS of every file in a package to all of its files, so
// this also covers the generated code) and emit our own warning in unstable_warning.go instead,
// which is silenced with the gst_use_unstable_api build tag.

// #cgo CPPFLAGS: -DGST_USE_UNSTABLE_API
import "C"
