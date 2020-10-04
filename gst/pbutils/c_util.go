package pbutils

// #include <gst/pbutils/pbutils.h>
import "C"
import "errors"

func initPbUtils() { C.gst_pb_utils_init() }

func wrapGerr(gerr *C.GError) error {
	defer C.g_error_free(gerr)
	return errors.New(C.GoString(gerr.message))
}

func gobool(b C.gboolean) bool { return int(b) > 0 }
