package gstapp

import (
	"runtime"

	"github.com/go-gst/go-glib/pkg/core/userdata"
	"github.com/go-gst/go-gst/pkg/gst"
)

// #cgo pkg-config: gstreamer-app-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/app/app.h>
// extern void _gogst_gstapp1_AppSink_eos(GstAppSink*, gpointer);
// extern GstFlowReturn _gogst_gstapp1_AppSink_new_preroll(GstAppSink*, gpointer);
// extern GstFlowReturn _gogst_gstapp1_AppSink_new_sample(GstAppSink*, gpointer);
// extern gboolean _gogst_gstapp1_AppSink_new_event(GstAppSink*, gpointer);
// extern gboolean _gogst_gstapp1_AppSink_propose_allocation(GstAppSink*, GstQuery*, gpointer);
// extern void destroyUserdata(gpointer);
import "C"

// AppSinkCallbacks contains the callbacks that can be installed on an [AppSink] with
// [AppSink.SetCallbacks]. All callbacks are optional, nil callbacks are not installed.
//
// All callbacks are called from the streaming thread, so they must not block for long.
//
// see also https://gstreamer.freedesktop.org/documentation/app/gstappsink.html#GstAppSinkCallbacks
type AppSinkCallbacks struct {
	// EOS gets invoked once the sink saw the end of the stream.
	EOS func(appsink AppSink)

	// NewPreroll gets invoked when a preroll sample is waiting in the sink. Fetch it with
	// [AppSink.PullPreroll], here or on another goroutine.
	NewPreroll func(appsink AppSink) gst.FlowReturn

	// NewSample gets invoked when a sample is waiting in the sink. Fetch it with
	// [AppSink.PullSample], here or on another goroutine.
	NewSample func(appsink AppSink) gst.FlowReturn

	// NewEvent gets invoked when a serialized event is waiting in the sink. Fetch it with
	// [AppSink.PullObject], here or on another goroutine. Return true when the event was handled.
	NewEvent func(appsink AppSink) bool

	// ProposeAllocation gets invoked for allocation queries arriving at the sink. Answer the
	// query by modifying it in place.
	//
	// The query is borrowed: it is only valid while the callback runs and must not be kept
	// around afterwards.
	ProposeAllocation func(appsink AppSink, query *gst.Query) bool
}

// AppSinkExtManual contains the handwritten methods of [AppSink].
type AppSinkExtManual interface {
	// SetCallbacks wraps gst_app_sink_set_callbacks
	//
	// The callbacks replace the signals of the appsink. They are called directly, without
	// marshaling the arguments through GValues. Any previously installed callbacks are replaced.
	//
	// Passing nil removes all previously installed callbacks.
	//
	// see also https://gstreamer.freedesktop.org/documentation/app/gstappsink.html#gst_app_sink_set_callbacks
	SetCallbacks(callbacks *AppSinkCallbacks)
}

var _ AppSinkExtManual = (*AppSinkInstance)(nil)

// SetCallbacks wraps gst_app_sink_set_callbacks
func (appsink *AppSinkInstance) SetCallbacks(callbacks *AppSinkCallbacks) {
	var carg0 *C.GstAppSink         // in, none, converted
	var carg1 C.GstAppSinkCallbacks // in, copied by gstreamer
	var carg2 C.gpointer            // closure
	var carg3 C.GDestroyNotify      // destroy

	carg0 = (*C.GstAppSink)(UnsafeAppSinkToGlibNone(appsink))

	if callbacks != nil {
		// copy the struct so that later modifications by the caller have no effect
		cbs := *callbacks

		if cbs.EOS != nil {
			carg1.eos = (*[0]byte)(C._gogst_gstapp1_AppSink_eos)
		}
		if cbs.NewPreroll != nil {
			carg1.new_preroll = (*[0]byte)(C._gogst_gstapp1_AppSink_new_preroll)
		}
		if cbs.NewSample != nil {
			carg1.new_sample = (*[0]byte)(C._gogst_gstapp1_AppSink_new_sample)
		}
		if cbs.NewEvent != nil {
			carg1.new_event = (*[0]byte)(C._gogst_gstapp1_AppSink_new_event)
		}
		if cbs.ProposeAllocation != nil {
			carg1.propose_allocation = (*[0]byte)(C._gogst_gstapp1_AppSink_propose_allocation)
		}

		carg2 = C.gpointer(userdata.Register(&cbs))
		carg3 = (C.GDestroyNotify)((*[0]byte)(C.destroyUserdata))
	}

	C.gst_app_sink_set_callbacks(carg0, &carg1, carg2, carg3)
	runtime.KeepAlive(appsink)
}
