package gst

// #include <gst/gst.h>
import "C"
import (
	"unsafe"
)

// gobool provides an easy type conversion between C.gboolean and a go bool.
func gobool(b C.gboolean) bool {
	return b != 0
}

// gboolean converts a go bool to a C.gboolean.
func gboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

// goStrings returns a string slice for an array of size argc starting at the address argv.
func goStrings(argc C.int, argv **C.gchar) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.gchar)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

// newQuarkFromString creates a new GQuark (or returns an existing one) for the given
// string
func newQuarkFromString(str string) C.uint {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	quark := C.g_quark_from_string(cstr)
	return quark
}
