package gst

// CGO exports have to be defined in a separate file from where they are used or else
// there will be double linkage issues.

// #include <gst/gst.h>
import "C"

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	gopointer "github.com/mattn/go-pointer"
)

//export structForEachCb
func structForEachCb(fieldID C.GQuark, val *C.GValue, chPtr C.gpointer) C.gboolean {
	ptr := gopointer.Restore(unsafe.Pointer(chPtr))
	resCh := ptr.(chan interface{})
	fieldName := C.GoString(C.g_quark_to_string(fieldID))

	var resValue interface{}

	gVal := glib.ValueFromNative(unsafe.Pointer(val))
	if resValue, _ = gVal.GoValue(); resValue == nil {
		// serialize the value if we can't do anything else with it
		serialized := C.gst_value_serialize(val)
		defer C.free(unsafe.Pointer(serialized))
		resValue = C.GoString(serialized)
	}

	resCh <- fieldName
	resCh <- resValue
	return gboolean(true)
}

//export goBusFunc
func goBusFunc(bus *C.GstBus, cMsg *C.GstMessage, userData C.gpointer) C.gboolean {
	// wrap the message
	msg := wrapMessage(cMsg)

	// retrieve the ptr to the function
	ptr := unsafe.Pointer(userData)
	funcIface := gopointer.Restore(ptr)
	busFunc, ok := funcIface.(BusWatchFunc)
	if !ok {
		gopointer.Unref(ptr)
		return gboolean(false)
	}

	// run the call back
	if cont := busFunc(msg); !cont {
		gopointer.Unref(ptr)
		return gboolean(false)
	}

	return gboolean(true)
}
