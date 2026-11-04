package gstsdp_test

import (
	"bytes"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstsdp"
)

func init() {
	gst.Init()
}

// waitWritable runs the GC until msg is writable, and fails if that takes more than a few seconds.
func waitWritable(t *testing.T, msg *gstsdp.MIKEYMessage, reason string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !isWritable(msg) {
		if time.Now().After(deadline) {
			t.Fatal(reason)
		}
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
}

func isWritable(msg *gstsdp.MIKEYMessage) bool {
	writable := gst.UnsafeMiniObjectFromGlibBorrow(gstsdp.UnsafeMIKEYMessageToGlibNone(msg)).IsWritable()
	runtime.KeepAlive(msg)
	return writable
}

// A wrapper that does not take ownership must hold its own reference and drop it when collected.
func TestMIKEYMessageReference(t *testing.T) {
	msg := gstsdp.NewMIKEYMessage()
	func() {
		shared := gstsdp.UnsafeMIKEYMessageFromGlibNone(gstsdp.UnsafeMIKEYMessageToGlibNone(msg))
		if isWritable(msg) {
			t.Fatal("message is writable while a second wrapper holds a reference")
		}
		runtime.KeepAlive(shared)
	}()
	waitWritable(t, msg, "message still has an extra reference after the second wrapper was collected")
}

func TestMIKEYMessageGetPayload(t *testing.T) {
	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(prev)

	msg := gstsdp.NewMIKEYMessage()
	payload := gstsdp.NewMIKEYPayload(gstsdp.MikeyPtT)
	if !msg.AddPayload(payload) {
		t.Fatal("AddPayload failed")
	}

	for range 10 {
		if msg.GetPayload(0) == nil {
			t.Fatal("message has no payload")
		}
	}
	runtime.GC()

	if n := msg.GetNPayloads(); n != 1 {
		t.Fatalf("message has %d payloads, want 1", n)
	}
	if buf.Len() != 0 {
		t.Fatalf("GetPayload logged: %s", buf.String())
	}
}
