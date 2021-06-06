package gst

/*
#include "gst.go.h"

static void* allocArgv(int argc) {
    return malloc(sizeof(char *) * argc);
}
*/
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
		cargc := C.int(len(*args))
		cargv := (*[0xfff]*C.char)(C.allocArgv(cargc))
		defer C.free(unsafe.Pointer(cargv))
		for i, arg := range *args {
			cargv[i] = C.CString(arg)
			defer C.free(unsafe.Pointer(cargv[i]))
		}
		C.gst_init(&cargc, (***C.char)(unsafe.Pointer(&cargv)))
		unhandled := make([]string, cargc)
		for i := 0; i < int(cargc); i++ {
			unhandled[i] = C.GoString(cargv[i])
		}
		*args = unhandled
	} else {
		C.gst_init(nil, nil)
	}
	CAT = NewDebugCategory("GST_GO", DebugColorFgCyan, "GStreamer Go Bindings")
}

// Deinit is a wrapper for gst_deinit Clean up any resources created by GStreamer in gst_init().
// It is normally not needed to call this function in a normal application as the resources will automatically be freed
// when the program terminates. This function is therefore mostly used by testsuites and other memory profiling tools.
func Deinit() {
	C.gst_deinit()
}
