package gstwebrtc_test

import (
	"testing"

	"github.com/go-gst/go-gst/gst/gstsdp"
	"github.com/go-gst/go-gst/gst/gstwebrtc"
)

func TestSessionDescriptionGValueMarshal(t *testing.T) {
	sdp, err := gstsdp.ParseSDPMessage("v=0\nm=audio 4000 RTP/AVP 111\na=rtpmap:111 OPUS/48000/2\nm=video 4000 RTP/AVP 96\na=rtpmap:96 VP8/90000\na=my-sdp-value")

	if err != nil {
		t.Fatal(err)
	}

	sd := gstwebrtc.NewSessionDescription(gstwebrtc.SDP_TYPE_OFFER, sdp)

	gv, err := sd.ToGValue()

	if err != nil {
		t.Fatal(err)
	}

	sdI, err := gv.GoValue()

	if err != nil {
		t.Fatal(err)
	}

	sd, ok := sdI.(*gstwebrtc.SessionDescription)

	if !ok {
		t.Fatal("Failed to convert to SessionDescription")
	}

	_ = sd
}

func TestSessionDescriptionJSONMarshal(t *testing.T) {
	sdp, err := gstsdp.ParseSDPMessage("v=0\nm=audio 4000 RTP/AVP 111\na=rtpmap:111 OPUS/48000/2\nm=video 4000 RTP/AVP 96\na=rtpmap:96 VP8/90000\na=my-sdp-value")

	if err != nil {
		t.Fatal(err)
	}

	sd := gstwebrtc.NewSessionDescription(gstwebrtc.SDP_TYPE_OFFER, sdp)

	w3 := sd.ToW3SDP()

	sd, err = w3.ToGstSDP()

	if err != nil {
		t.Fatal(err)
	}

	_ = sd
}
