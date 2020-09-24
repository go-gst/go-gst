package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"

// URIType casts C GstURIType to a go type
type URIType C.GstURIType

// Type cast URI types
const (
	URIUnknown URIType = C.GST_URI_UNKNOWN // (0) – The URI direction is unknown
	URISink            = C.GST_URI_SINK    // (1) – The URI is a consumer.
	URISource          = C.GST_URI_SRC     // (2) - The URI is a producer.
)

func (u URIType) String() string {
	switch u {
	case URIUnknown:
		return "Unknown"
	case URISink:
		return "Sink"
	case URISource:
		return "Source"
	}
	return ""
}
