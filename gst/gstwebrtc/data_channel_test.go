package gstwebrtc

import (
	"testing"

	"github.com/go-gst/go-gst/gst"
)

func TestDataChannelMarshal(t *testing.T) {
	gst.Init(nil)

	// hack to get a valid glib.Object
	el, err := gst.NewElement("webrtcbin")

	if err != nil {
		t.Error(err)
	}

	dc := &DataChannel{
		Object: el.Object.Object,
	}

	gv, err := dc.ToGValue()

	if err != nil {
		t.Error(err)
	}

	dcI, err := gv.GoValue()

	if err != nil {
		t.Error(err)
	}

	dc, ok := dcI.(*DataChannel)

	if !ok {
		t.Error("Failed to convert to DataChannel")
	}

	_ = dc
}
