package gst

import "github.com/diamondburned/gotk4/pkg/gobject/v2"

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

const (
	// ParamConditionallyAvailable is used on GObject properties of GstObject to indicate that
	// they might not be available depending on environment such as OS, device, etc, so such properties
	// will be installed conditionally only if the GstObject is able to support it.
	ParamConditionallyAvailable gobject.ParamFlags = gobject.ParamFlags(C.GST_PARAM_CONDITIONALLY_AVAILABLE)

	// ParamControllable is used on GObject properties to signal they can make sense to be controlled over time.
	// This hint is used by the GstController.
	ParamControllable gobject.ParamFlags = gobject.ParamFlags(C.GST_PARAM_CONTROLLABLE)

	// ParamDocShowDefault is used on GObject properties of GstObject to indicate that during gst-inspect and friends,
	// the default value should be used as default instead of the current value.
	ParamDocShowDefault gobject.ParamFlags = gobject.ParamFlags(C.GST_PARAM_DOC_SHOW_DEFAULT)

	// ParamMutablePaused is used on GObject properties of GstElements to indicate that they can be changed when the element is in the PAUSED or lower state.
	// This flag implies GST_PARAM_MUTABLE_READY.
	ParamMutablePaused gobject.ParamFlags = gobject.ParamFlags(C.GST_PARAM_MUTABLE_PAUSED)
	// ParamMutablePlaying is used on GObject properties of GstElements to indicate that they can be changed when the element is in the PLAYING or lower state.
	// This flag implies GST_PARAM_MUTABLE_PAUSED.
	ParamMutablePlaying gobject.ParamFlags = gobject.ParamFlags(C.GST_PARAM_MUTABLE_PLAYING)

	// ParamMutableReady is used on GObject properties of GstElements to indicate that they can be changed when the element is in the READY or lower state.
	ParamMutableReady gobject.ParamFlags = gobject.ParamFlags(C.GST_PARAM_MUTABLE_READY)

	// ParamUserShift is used on GObject properties of GstElements to indicate that they can be changed when the element is in the READY or lower state.
	// Bits based on GST_PARAM_USER_SHIFT can be used by 3rd party applications.
	ParamUserShift gobject.ParamFlags = gobject.ParamFlags(C.GST_PARAM_USER_SHIFT)
)
