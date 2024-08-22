package gstsdp

// #include "gst.go.h"
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type SDPResult C.GstSDPResult

const (
	SDPResultOk SDPResult = C.GST_SDP_OK
	SDPEinval   SDPResult = C.GST_SDP_EINVAL
)

type Message struct {
	ptr *C.GstSDPMessage
}

func wrapSDPMessageAndFinalize(sdp *C.GstSDPMessage) *Message {
	msg := &Message{
		ptr: sdp,
	}

	// this requires that we copy the SDP message before passing it to any transfer-ownership function
	runtime.SetFinalizer(msg, func(msg *Message) {
		msg.Free()
	})

	return msg
}

// NewMessageFromUnsafe creates a new SDP message from a pointer and does not finalize it
func NewMessageFromUnsafe(ptr unsafe.Pointer) *Message {
	return &Message{
		ptr: (*C.GstSDPMessage)(ptr),
	}
}

var ErrSDPInvalid = errors.New("invalid SDP")

func ParseSDPMessage(sdp string) (*Message, error) {
	cstr := C.CString(sdp)
	defer C.free(unsafe.Pointer(cstr))

	var msg *C.GstSDPMessage

	res := SDPResult(C.gst_sdp_message_new_from_text(cstr, &msg))

	if res != SDPResultOk || msg == nil {
		return nil, ErrSDPInvalid
	}

	return wrapSDPMessageAndFinalize(msg), nil
}

func (msg *Message) String() string {
	cstr := C.gst_sdp_message_as_text(msg.ptr)
	defer C.free(unsafe.Pointer(cstr))

	return C.GoString(cstr)
}

// UnownedCopy creates a new copy of the SDP message that will not be finalized
//
// this is needed to pass the message back to C where C takes ownership of the message
//
// the returned SDP message will leak memory if not freed manually
func (msg *Message) UnownedCopy() *Message {
	var newMsg *C.GstSDPMessage
	res := C.gst_sdp_message_copy(msg.ptr, &newMsg)

	if res != C.GST_SDP_OK || newMsg == nil {
		return nil
	}

	return &Message{
		ptr: newMsg,
	}
}

// Free frees the SDP message.
//
// This is called automatically when the object is garbage collected.
func (msg *Message) Free() {
	C.gst_sdp_message_free(msg.ptr)
	msg.ptr = nil
}

func (msg *Message) Instance() unsafe.Pointer {
	return unsafe.Pointer(msg.ptr)
}
