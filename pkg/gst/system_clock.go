package gst

import "github.com/go-gst/go-glib/pkg/gobject/v2"

// NewSystemClock creates a new instance of a SystemClock, with the given clock type parameter
//
// This is only a convenience wrapper for gobject.NewObjectWithProperties
func NewSystemClock(clockType ClockType) SystemClock {
	clockObj := gobject.NewObjectWithProperties(TypeSystemClock, map[string]any{
		"clock-type": clockType,
	})

	return clockObj.(SystemClock)
}
