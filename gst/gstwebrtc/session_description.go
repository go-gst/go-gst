package gstwebrtc

// #include "gst.go.h"
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst/gstsdp"
)

func init() {

	tm := []glib.TypeMarshaler{
		{T: glib.Type(C.GST_TYPE_WEBRTC_SESSION_DESCRIPTION), F: marshalSessionDescription},
	}

	glib.RegisterGValueMarshalers(tm)
}

type SessionDescription struct {
	ptr *C.GstWebRTCSessionDescription
}

func NewSessionDescription(t SDPType, sdp *gstsdp.Message) *SessionDescription {
	sd := C.gst_webrtc_session_description_new(
		C.GstWebRTCSDPType(t),
		(*C.GstSDPMessage)(sdp.UnownedCopy().Instance()),
	)

	return wrapSessionDescriptionAndFinalize(sd)
}

func wrapSessionDescriptionAndFinalize(sdp *C.GstWebRTCSessionDescription) *SessionDescription {
	sd := &SessionDescription{
		ptr: sdp,
	}

	// this requires that we copy the SDP message before passing it to any transfer-ownership function
	runtime.SetFinalizer(sd, func(sd *SessionDescription) {
		sd.Free()
	})

	return sd
}

// W3RTCSessionDescription is used to marshal/unmarshal SessionDescription to/from JSON.
//
// We cannot implement the json.(Un-)Marshaler interfaces on SessionDescription directly because
// the finalizer would run and free the memory, because the value would have to be copied.
//
// it complies with the WebRTC spec for SessionDescription, see https://www.w3.org/TR/webrtc/#rtcsessiondescription-class
type W3RTCSessionDescription struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

// ToGstSDP converts a W3RTCSessionDescription to a SessionDescription
func (w3SDP *W3RTCSessionDescription) ToGstSDP() (*SessionDescription, error) {
	sdp, err := gstsdp.ParseSDPMessage(w3SDP.Sdp)
	if err != nil {
		return nil, err
	}

	return NewSessionDescription(SDPTypeFromString(w3SDP.Type), sdp), nil
}

// ToW3SDP returns a W3RTCSessionDescription that can be marshaled to JSON
func (sd *SessionDescription) ToW3SDP() W3RTCSessionDescription {
	jsonSDP := W3RTCSessionDescription{
		Type: SDPType(sd.ptr._type).String(),
		Sdp:  sd.SDP().String(),
	}

	return jsonSDP
}

func (sd *SessionDescription) Free() {
	C.gst_webrtc_session_description_free(sd.ptr)
}

// UnownedCopy creates a new copy of the SessionDescription that will not be finalized
//
// this is needed for passing the SessionDescription to other functions that will take ownership of it.
//
// used in the bindings, should not be called by application code
func (sd *SessionDescription) UnownedCopy() *SessionDescription {
	newSD := C.gst_webrtc_session_description_copy(sd.ptr)

	return &SessionDescription{
		ptr: newSD,
	}
}

// Copy creates a new copy of the SessionDescription
func (sd *SessionDescription) Copy() *SessionDescription {
	return wrapSessionDescriptionAndFinalize(sd.UnownedCopy().ptr)
}

// ToGValue implements glib.ValueTransformer
func (sd *SessionDescription) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(glib.Type(C.GST_TYPE_WEBRTC_SESSION_DESCRIPTION))
	if err != nil {
		return nil, err
	}
	var ptr *C.GstWebRTCSessionDescription
	if sd != nil {
		ptr = sd.ptr
	}
	val.SetBoxed(unsafe.Pointer(ptr))
	return val, nil
}

func marshalSessionDescription(p unsafe.Pointer) (interface{}, error) {
	c := C.g_value_get_boxed((*C.GValue)(p))

	// we don't own this memory, so we need to copy it to prevent other code from freeing it
	ref := &SessionDescription{
		ptr: (*C.GstWebRTCSessionDescription)(c),
	}

	return ref.Copy(), nil
}

// Copy creates a new copy of the SessionDescription
func (sd *SessionDescription) SDP() *gstsdp.Message {
	sdp := gstsdp.NewMessageFromUnsafe(unsafe.Pointer(sd.ptr.sdp))

	runtime.SetFinalizer(sdp, func(sdp *gstsdp.Message) {
		runtime.KeepAlive(sd)
	})

	return sdp
}
