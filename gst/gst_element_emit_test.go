package gst_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

func TestSignalEmitSimpleReturnValue(t *testing.T) {
	gst.Init(nil)
	webrtcbin, err := gst.NewElement("webrtcbin")

	if err != nil {
		t.Fatal(err)
	}

	okI, err := webrtcbin.Emit("add-turn-server", "turn://user:password@host:1234")

	if err != nil {
		t.Fatal(err)
	}

	ok := okI.(bool)

	if !ok {
		t.Fatal("Failed to add turn server")
	}
}

func TestSignalEmitVoidReturnValue(t *testing.T) {
	gst.Init(nil)

	elem, err := gst.NewElement("splitmuxsink")
	if err != nil {
		t.Fatal(err)
	}

	result, err := elem.Emit("split-after")
	if err != nil {
		t.Fatal("Result must be nil due to void return type, unless splitmux api changed" +
			err.Error())
	}
	if result != nil {
		t.Fatal("Result must be nil due to void return type, unless splitmux api changed")
	}
}

func TestSignalEmitGObjectReturnValue(t *testing.T) {
	gst.Init(nil)

	elements := []string{
		"rtpbin", "name=rtpbin",

		"videotestsrc", "!", "videoconvert", "!", "queue", "!",
		"x264enc", "bframes=0", "speed-preset=ultrafast", "tune=zerolatency", "name=encoder", "!", "queue", "!", "rtph264pay", "config-interval=1", "!", "rtpbin.send_rtp_sink_0", "rtpbin.send_rtp_src_0", "!",

		"udpsink", "host=127.0.0.1", "port=5510", "sync=false", "async=false", "rtpbin.send_rtcp_src_0", "!",

		"udpsink", "host=127.0.0.1", "port=5511", "sync=false", "async=false",

		"udpsrc", "port=5511", "caps=\"application/x-rtcp\"", "!", "rtpbin.recv_rtcp_sink_0",
	}

	pipeline, err := gst.NewPipelineFromString(strings.Join(elements, " "))
	if err != nil {
		t.Fatal(err)
	}

	rtpbin, err := pipeline.GetElementByName("rtpbin")
	if err != nil {
		t.Fatal(err)
	}

	errchan := make(chan error)

	rtpbin.Connect("on-new-ssrc", func(bin *gst.Element, sessionID uint, ssrc uint32) {
		retI, err := rtpbin.Emit("get-internal-session", sessionID)

		if err != nil {
			errchan <- err
		}

		rtpSession, ok := retI.(*glib.Object)

		if !ok {
			errchan <- errors.New("could not cast return value to *glib.Object")
		}

		bw, err := rtpSession.GetProperty("bandwidth")

		if err != nil {
			errchan <- err
		}

		_ = bw

		close(errchan)
	})

	pipeline.SetState(gst.StatePlaying)

	err = <-errchan

	if err != nil {
		t.Fatal(err)
	}

	pipeline.SetState(gst.StateNull)

}
