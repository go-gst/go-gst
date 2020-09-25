package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Init is a wrapper around gst_init() and must be called before any
// other gstreamer calls and is used to initialize everything necessary.
// In addition to setting up gstreamer for usage, a pointer to a slice of
// strings may be passed in to parse standard gst command line arguments.
// args will be modified to remove any flags that were handled.
// Alternatively, nil may be passed in to not perform any command line
// parsing.
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
}

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

// structureToGoMap converts a GstStructure into a Go map of strings.
func structureToGoMap(st *C.GstStructure) map[string]string {
	goDetails := make(map[string]string)
	numFields := int(C.gst_structure_n_fields((*C.GstStructure)(st)))
	for i := 0; i < numFields-1; i++ {
		fieldName := C.gst_structure_nth_field_name((*C.GstStructure)(st), (C.guint)(i))
		fieldValue := C.gst_structure_get_value((*C.GstStructure)(st), (*C.gchar)(fieldName))
		strValueDup := C.g_strdup_value_contents((*C.GValue)(fieldValue))
		goDetails[C.GoString(fieldName)] = C.GoString(strValueDup)
	}
	return goDetails
}

func iteratorToElementSlice(iterator *C.GstIterator) ([]*Element, error) {
	elems := make([]*Element, 0)
	gval := new(C.GValue)

	for {
		switch C.gst_iterator_next((*C.GstIterator)(iterator), (*C.GValue)(unsafe.Pointer(gval))) {
		case C.GST_ITERATOR_DONE:
			C.gst_iterator_free((*C.GstIterator)(iterator))
			return elems, nil
		case C.GST_ITERATOR_RESYNC:
			C.gst_iterator_resync((*C.GstIterator)(iterator))
		case C.GST_ITERATOR_OK:
			cElemVoid := C.g_value_get_object((*C.GValue)(gval))
			cElem := (*C.GstElement)(cElemVoid)
			elems = append(elems, wrapElement(glib.Take(unsafe.Pointer(cElem))))
			C.g_value_reset((*C.GValue)(gval))
		default:
			return nil, errors.New("Element iterator failed")
		}
	}
}

func goStrings(argc C.int, argv **C.gchar) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.gchar)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}
