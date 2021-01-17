package gst

// #include "gst.go.h"
import "C"

// Additional GStreamer ParamSpec flags
const (
	ParameterControllable   = C.GST_PARAM_CONTROLLABLE
	ParameterMutablePlaying = C.GST_PARAM_MUTABLE_PLAYING
	ParameterMutablePaused  = C.GST_PARAM_MUTABLE_PAUSED
	ParameterMutableReady   = C.GST_PARAM_MUTABLE_READY
)
