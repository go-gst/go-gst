package gstwebrtc

// #include "gst.go.h"
import "C"
import (
	"errors"
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

func init() {

	tm := []glib.TypeMarshaler{
		{T: glib.Type(C.GST_TYPE_WEBRTC_DATA_CHANNEL), F: marshalDataChannel},
	}

	glib.RegisterGValueMarshalers(tm)
}

// DataChannel is a representation of GstWebRTCDataChannel. See https://gstreamer.freedesktop.org/documentation/webrtclib/gstwebrtc-datachannel.html?gi-language=c
//
// there is no constructor for DataChannel, you can get it from webrtcbin signals
type DataChannel struct {
	*glib.Object
}

func (dc *DataChannel) Close() {
	C.gst_webrtc_data_channel_close((*C.GstWebRTCDataChannel)(dc.Native()))
}

func (dc *DataChannel) SendData(data []byte) error {
	var gerr *C.GError

	addr := unsafe.SliceData(data)

	cbytes := C.g_bytes_new(C.gconstpointer(addr), C.gsize(len(data)))
	defer C.g_bytes_unref(cbytes)

	C.gst_webrtc_data_channel_send_data_full((*C.GstWebRTCDataChannel)(dc.Native()), cbytes, &gerr)

	if gerr != nil {
		defer C.g_error_free((*C.GError)(gerr))
		errMsg := C.GoString(gerr.message)
		return errors.New(errMsg)
	}

	return nil
}

// ToGValue implements glib.ValueTransformer
func (dc *DataChannel) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(glib.Type(C.GST_TYPE_WEBRTC_DATA_CHANNEL))
	if err != nil {
		return nil, err
	}
	val.SetInstance(unsafe.Pointer(dc.GObject))
	return val, nil
}

func marshalDataChannel(p unsafe.Pointer) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(p))

	return &DataChannel{
		Object: glib.Take(unsafe.Pointer(c)),
	}, nil
}
