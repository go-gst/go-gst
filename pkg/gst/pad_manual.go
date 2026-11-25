package gst

import (
	"runtime"

	"github.com/go-gst/go-glib/pkg/core/userdata"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
// extern gboolean _gogst_gst1_PadEventFunction(GstPad*, GstObject*, GstEvent*);
// extern void destroyUserdata(gpointer);
import "C"

// PadEventFunction wraps GstPadEventFunction.
//
// The callback owns event (transfer full): it must either pass it on (e.g. via
// [Pad.EventDefault] or [Pad.PushEvent]) or drop it. Dropped events are
// unreffed by the garbage collector. parent may be nil.
//
// see also https://gstreamer.freedesktop.org/documentation/gstreamer/gstpad.html#GstPadEventFunction
type PadEventFunction func(pad Pad, parent Object, event *Event) bool

type PadExtManual interface {
	// SetEventFunction wraps gst_pad_set_event_function_full
	//
	// Passing nil restores the default handler, gst_pad_event_default.
	//
	// see also https://gstreamer.freedesktop.org/documentation/gstreamer/gstpad.html#gst_pad_set_event_function_full
	SetEventFunction(PadEventFunction)
}

// SetEventFunction wraps gst_pad_set_event_function_full
//
// Passing nil restores the default handler, gst_pad_event_default.
//
// see also https://gstreamer.freedesktop.org/documentation/gstreamer/gstpad.html#gst_pad_set_event_function_full
func (pad *PadInstance) SetEventFunction(fn PadEventFunction) {
	var carg0 *C.GstPad             // in, none, converted
	var carg1 C.GstPadEventFunction // callback, scope: notified, closure: carg2, destroy: carg3, nullable
	var carg2 C.gpointer            // implicit
	var carg3 C.GDestroyNotify      // implicit

	carg0 = (*C.GstPad)(UnsafePadToGlibNone(pad))
	if fn != nil {
		carg1 = (*[0]byte)(C._gogst_gst1_PadEventFunction)
		carg2 = C.gpointer(userdata.Register(fn))
		carg3 = (C.GDestroyNotify)((*[0]byte)(C.destroyUserdata))
	} else {
		// restore the handler every pad starts with; a NULL event function
		// makes GStreamer warn and drop every event
		carg1 = (*[0]byte)(C.gst_pad_event_default)
	}

	C.gst_pad_set_event_function_full(carg0, carg1, carg2, carg3)
	runtime.KeepAlive(pad)
	runtime.KeepAlive(fn)
}
