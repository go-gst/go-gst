package gst

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// Init binds to the gst_init() function. Argument parsing is not
// supported.
func Init() {
	C.gst_init(nil, nil)
}
