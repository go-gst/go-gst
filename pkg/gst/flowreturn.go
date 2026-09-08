package gst

import (
	"errors"
	"fmt"
)

// Sentinel errors returned by [FlowReturn.Error], one per unsuccessful
// [FlowReturn]. Test for them with [errors.Is] rather than comparing error
// strings:
//
//	if err := ret.Error(); errors.Is(err, gst.ErrFlowFlushing) {
//		// the pad is flushing, stop pushing data
//	}
var (
	ErrFlowNotLinked     = errors.New("pad is not linked")
	ErrFlowFlushing      = errors.New("pad is flushing")
	ErrFlowEOS           = errors.New("pad is end of stream")
	ErrFlowNotNegotiated = errors.New("pad is not negotiated")
	ErrFlowError         = errors.New("data flow error")
	ErrFlowNotSupported  = errors.New("operation is not supported")
)

// IsSuccess reports whether f is a success code, i.e. [FlowOK] or one of the
// FlowCustomSuccess values. It mirrors the GST_FLOW_IS_SUCCESS macro.
func (f FlowReturn) IsSuccess() bool {
	return f >= FlowOK
}

// Error returns nil if f is a success code and an error describing f otherwise.
//
// The returned error matches the corresponding Err* sentinel with [errors.Is].
// Values in the custom error range have no sentinel and are reported by name.
//
// Note that this deliberately does not make [FlowReturn] implement the error
// interface: a FlowReturn is not itself an error, since it also carries success
// codes.
func (f FlowReturn) Error() error {
	switch f {
	case FlowNotLinked:
		return ErrFlowNotLinked
	case FlowFlushing:
		return ErrFlowFlushing
	case FlowEOS:
		return ErrFlowEOS
	case FlowNotNegotiated:
		return ErrFlowNotNegotiated
	case FlowError:
		return ErrFlowError
	case FlowNotSupported:
		return ErrFlowNotSupported
	}

	if f.IsSuccess() {
		return nil
	}

	return fmt.Errorf("data flow error: %s", f.String())
}
