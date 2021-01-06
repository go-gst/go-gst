package gst

// GError is a Go wrapper for a C GError in the context of GStreamer. It implements the error interface
// and provides additional functions for retrieving debug strings and details.
type GError struct {
	errMsg, debugStr string
	structure        *Structure

	// used for message constructors
	code ErrorCode
}

// Message is an alias to `Error()`. It's for clarity when this object
// is parsed from a `GST_MESSAGE_INFO` or `GST_MESSAGE_WARNING`.
func (e *GError) Message() string { return e.Error() }

// Error implements the error interface and returns the error message.
func (e *GError) Error() string { return e.errMsg }

// DebugString returns any debug info alongside the error.
func (e *GError) DebugString() string { return e.debugStr }

// Structure returns the structure of the error message which may contain additional metadata.
func (e *GError) Structure() *Structure { return e.structure }

// Code returns the error code of the error message.
func (e *GError) Code() ErrorCode { return e.code }

// NewGError wraps the given error inside a GError (to be used with message constructors).
func NewGError(code ErrorCode, err error) *GError {
	return &GError{
		errMsg: err.Error(),
		code:   code,
	}
}
