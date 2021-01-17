package gst

/*
#include "gst.go.h"
*/
import "C"

// Domain represents the different types of error domains.
type Domain string

// ErrorDomain castings
const (
	DomainCore     Domain = "CORE"
	DomainLibrary  Domain = "LIBRARY"
	DomainResource Domain = "RESOURCE"
	DomainStream   Domain = "STREAM"
)

func (d Domain) toQuark() C.GQuark {
	switch d {
	case DomainCore:
		return C.gst_core_error_quark()
	case DomainLibrary:
		return C.gst_library_error_quark()
	case DomainResource:
		return C.gst_resource_error_quark()
	case DomainStream:
		return C.gst_stream_error_quark()
	default:
		return C.gst_library_error_quark()
	}
}

// ErrorCode represents GstGError codes.
type ErrorCode int

// Type castings of CoreErrors
const (
	CoreErrorFailed         ErrorCode = C.GST_CORE_ERROR_FAILED          // (1) – a general error which doesn't fit in any other category. Make sure you add a custom message to the error call.
	CoreErrorTooLazy        ErrorCode = C.GST_CORE_ERROR_TOO_LAZY        // (2) – do not use this except as a placeholder for deciding where to go while developing code.
	CoreErrorNotImplemented ErrorCode = C.GST_CORE_ERROR_NOT_IMPLEMENTED // (3) – use this when you do not want to implement this functionality yet.
	CoreErrorStateChange    ErrorCode = C.GST_CORE_ERROR_STATE_CHANGE    // (4) – used for state change errors.
	CoreErrorPad            ErrorCode = C.GST_CORE_ERROR_PAD             // (5) – used for pad-related errors.
	CoreErrorThread         ErrorCode = C.GST_CORE_ERROR_THREAD          // (6) – used for thread-related errors.
	CoreErrorNegotiation    ErrorCode = C.GST_CORE_ERROR_NEGOTIATION     // (7) – used for negotiation-related errors.
	CoreErrorEvent          ErrorCode = C.GST_CORE_ERROR_EVENT           // (8) – used for event-related errors.
	CoreErrorSeek           ErrorCode = C.GST_CORE_ERROR_SEEK            // (9) – used for seek-related errors.
	CoreErrorCaps           ErrorCode = C.GST_CORE_ERROR_CAPS            // (10) – used for caps-related errors.
	CoreErrorTag            ErrorCode = C.GST_CORE_ERROR_TAG             // (11) – used for negotiation-related errors.
	CoreErrorMissingPlugin  ErrorCode = C.GST_CORE_ERROR_MISSING_PLUGIN  // (12) – used if a plugin is missing.
	CoreErrorClock          ErrorCode = C.GST_CORE_ERROR_CLOCK           // (13) – used for clock related errors.
	CoreErrorDisabled       ErrorCode = C.GST_CORE_ERROR_DISABLED        // (14) – used if functionality has been disabled at compile time.
)

// Type castings for LibraryErrors
const (
	LibraryErrorFailed   ErrorCode = C.GST_LIBRARY_ERROR_FAILED   // (1) – a general error which doesn't fit in any other category. Make sure you add a custom message to the error call.
	LibraryErrorTooLazy  ErrorCode = C.GST_LIBRARY_ERROR_TOO_LAZY // (2) – do not use this except as a placeholder for deciding where to go while developing code.
	LibraryErrorInit     ErrorCode = C.GST_LIBRARY_ERROR_INIT     // (3) – used when the library could not be opened.
	LibraryErrorShutdown ErrorCode = C.GST_LIBRARY_ERROR_SHUTDOWN // (4) – used when the library could not be closed.
	LibraryErrorSettings ErrorCode = C.GST_LIBRARY_ERROR_SETTINGS // (5) – used when the library doesn't accept settings.
	LibraryErrorEncode   ErrorCode = C.GST_LIBRARY_ERROR_ENCODE   // (6) – used when the library generated an encoding error.
)

// Type castings for ResourceErrors
const (
	ResourceErrorFailed        ErrorCode = C.GST_RESOURCE_ERROR_FAILED          // (1) – a general error which doesn't fit in any other category. Make sure you add a custom message to the error call.
	ResourceErrorTooLazy       ErrorCode = C.GST_RESOURCE_ERROR_TOO_LAZY        // (2) – do not use this except as a placeholder for deciding where to go while developing code.
	ResourceErrorNotFound      ErrorCode = C.GST_RESOURCE_ERROR_NOT_FOUND       // (3) – used when the resource could not be found.
	ResourceErrorBusy          ErrorCode = C.GST_RESOURCE_ERROR_BUSY            // (4) – used when resource is busy.
	ResourceErrorOpenRead      ErrorCode = C.GST_RESOURCE_ERROR_OPEN_READ       // (5) – used when resource fails to open for reading.
	ResourceErrorOpenWrite     ErrorCode = C.GST_RESOURCE_ERROR_OPEN_WRITE      // (6) – used when resource fails to open for writing.
	ResourceErrorOpenReadWrite ErrorCode = C.GST_RESOURCE_ERROR_OPEN_READ_WRITE // (7) – used when resource cannot be opened for both reading and writing, or either (but unspecified which).
	ResourceErrorClose         ErrorCode = C.GST_RESOURCE_ERROR_CLOSE           // (8) – used when the resource can't be closed.
	ResourceErrorRead          ErrorCode = C.GST_RESOURCE_ERROR_READ            // (9) – used when the resource can't be read from.
	ResourceErrorWrite         ErrorCode = C.GST_RESOURCE_ERROR_WRITE           // (10) – used when the resource can't be written to.
	ResourceErrorSeek          ErrorCode = C.GST_RESOURCE_ERROR_SEEK            // (11) – used when a seek on the resource fails.
	ResourceErrorSync          ErrorCode = C.GST_RESOURCE_ERROR_SYNC            // (12) – used when a synchronize on the resource fails.
	ResourceErrorSettings      ErrorCode = C.GST_RESOURCE_ERROR_SETTINGS        // (13) – used when settings can't be manipulated on.
	ResourceErrorNoSpaceLeft   ErrorCode = C.GST_RESOURCE_ERROR_NO_SPACE_LEFT   // (14) – used when the resource has no space left.
	ResourceErrorNotAuthorized ErrorCode = C.GST_RESOURCE_ERROR_NOT_AUTHORIZED  // (15) – used when the resource can't be opened due to missing authorization. (Since: 1.2.4)
)

// Type castings for StreamErrors
const (
	StreamErrorFailed         ErrorCode = C.GST_STREAM_ERROR_FAILED          // (1) – a general error which doesn't fit in any other category. Make sure you add a custom message to the error call.
	StreamErrorTooLazy        ErrorCode = C.GST_STREAM_ERROR_TOO_LAZY        // (2) – do not use this except as a placeholder for deciding where to go while developing code.
	StreamErrorNotImplemented ErrorCode = C.GST_STREAM_ERROR_NOT_IMPLEMENTED // (3) – use this when you do not want to implement this functionality yet.
	StreamErrorTypeNotFound   ErrorCode = C.GST_STREAM_ERROR_TYPE_NOT_FOUND  // (4) – used when the element doesn't know the stream's type.
	StreamErrorWrongType      ErrorCode = C.GST_STREAM_ERROR_WRONG_TYPE      // (5) – used when the element doesn't handle this type of stream.
	StreamErrorCodecNotFound  ErrorCode = C.GST_STREAM_ERROR_CODEC_NOT_FOUND // (6) – used when there's no codec to handle the stream's type.
	StreamErrorDecode         ErrorCode = C.GST_STREAM_ERROR_DECODE          // (7) – used when decoding fails.
	StreamErrorEncode         ErrorCode = C.GST_STREAM_ERROR_ENCODE          // (8) – used when encoding fails.
	StreamErrorDemux          ErrorCode = C.GST_STREAM_ERROR_DEMUX           // (9) – used when demuxing fails.
	StreamErrorMux            ErrorCode = C.GST_STREAM_ERROR_MUX             // (10) – used when muxing fails.
	StreamErrorFormat         ErrorCode = C.GST_STREAM_ERROR_FORMAT          // (11) – used when the stream is of the wrong format (for example, wrong caps).
	StreamErrorDecrypt        ErrorCode = C.GST_STREAM_ERROR_DECRYPT         // (12) – used when the stream is encrypted and can't be decrypted because this is not supported by the element.
	StreamErrorDecryptNoKey   ErrorCode = C.GST_STREAM_ERROR_DECRYPT_NOKEY   // (13) – used when the stream is encrypted and can't be decrypted because no suitable key is available.
)

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
