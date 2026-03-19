package gst

import (
	"path"
	"runtime"

	"github.com/go-gst/go-glib/pkg/glib/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

type ElementExtManual interface {
	// GetCurrentState returns the value of the current_state member of the struct
	//
	// the current state of an element
	GetCurrentState() State

	// BlockSetState is a convenience wrapper around calling [Element.SetState] and [Element.GetState] to wait for async state changes. See [Element.State] for more info.
	BlockSetState(state State, timeout ClockTime) StateChangeReturn

	// MessageError is a convenience wrapper for posting an error message from inside an element. See [Element.MessageFull] for more info.
	MessageError(domain glib.Quark, code int32, text, debug string)
}

// GetCurrentState returns the value of the current_state member of the struct
//
// the current state of an element
func (el *ElementInstance) GetCurrentState() State {
	native := (*C.GstElement)(UnsafeElementToGlibNone(el))

	return State(native.current_state)
}

// BlockSetState is a convenience wrapper around calling [Element.SetState] and [Element.GetState] to wait for async state changes. See State for more info.
func (el *ElementInstance) BlockSetState(state State, timeout ClockTime) StateChangeReturn {
	ret := el.SetState(state)

	if ret == StateChangeAsync {
		_, _, ret = el.GetState(timeout)
	}

	return ret
}

// MessageError is a convenience wrapper for posting an error message from inside an element. See [Element.MessageFull] for more info.
func (e *ElementInstance) MessageError(domain glib.Quark, code int32, text, debug string) {
	function, file, line, _ := runtime.Caller(1)
	e.MessageFull(MessageError, domain, code, text, debug, path.Base(file), runtime.FuncForPC(function).Name(), int32(line))
}

func LinkMany(elements ...Element) bool {
	if len(elements) < 2 {
		return false
	}

	for i := 0; i < len(elements)-1; i++ {
		if !elements[i].Link(elements[i+1]) {
			return false
		}
	}

	return true
}
