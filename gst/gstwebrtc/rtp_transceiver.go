package gstwebrtc

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

func init() {

	tm := []glib.TypeMarshaler{
		{T: glib.Type(C.GST_TYPE_WEBRTC_RTP_TRANSCEIVER), F: marshalRTPTransceiver},
	}

	glib.RegisterGValueMarshalers(tm)
}

type RTPTransceiver struct {
	*gst.Object
}

// ToGValue implements glib.ValueTransformer
func (tc *RTPTransceiver) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(glib.Type(C.GST_TYPE_WEBRTC_RTP_TRANSCEIVER))
	if err != nil {
		return nil, err
	}
	val.SetInstance(unsafe.Pointer(tc.Instance()))
	return val, nil
}

func wrapRTPTransceiver(p unsafe.Pointer) *RTPTransceiver {
	return &RTPTransceiver{
		Object: gst.FromGstObjectUnsafeNone(p),
	}
}

func marshalRTPTransceiver(p unsafe.Pointer) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(p))

	return wrapRTPTransceiver(unsafe.Pointer(c)), nil
}
