package gst

// #include "gst.go.h"
import "C"
import "unsafe"

// CAT is the global DebugCategory used for logging from the bindings. It is okay to use
// it from applications, but plugins should use their own.
var CAT *DebugCategory

// Init is a wrapper around gst_init() and must be called before any
// other gstreamer calls and is used to initialize everything necessary.
// In addition to setting up gstreamer for usage, a pointer to a slice of
// strings may be passed in to parse standard gst command line arguments.
// args will be modified to remove any flags that were handled.
// Alternatively, nil may be passed in to not perform any command line
// parsing.
//
// The bindings will also set up their own internal DebugCategory for logging
// than can be invoked from applications or plugins as well. However, for
// plugins it is generally better to initialize your own DebugCategory.
func Init(args *[]string) {
	if args != nil {
		argc := C.int(len(*args))
		argv := make([]*C.char, argc)
		for i, arg := range *args {
			argv[i] = C.CString(arg)
		}
		C.gst_init((*C.int)(unsafe.Pointer(&argc)),
			(***C.char)(unsafe.Pointer(&argv)))
		unhandled := make([]string, argc)
		for i := 0; i < int(argc); i++ {
			unhandled[i] = C.GoString(argv[i])
			C.free(unsafe.Pointer(argv[i]))
		}
		*args = unhandled
	} else {
		C.gst_init(nil, nil)
	}
	CAT = NewDebugCategory("GST_GO", DebugColorNone, "GStreamer Go Bindings")
}
