package gst_test

import (
	"bytes"
	"log"
	"runtime"
	"testing"

	"github.com/go-gst/go-gst/pkg/gst"
)

func TestReferenceTimestampMeta(t *testing.T) {
	gst.Init()

	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(prev)

	reference := gst.CapsFromString("timestamp/x-test")
	buffer := gst.NewBuffer()

	if buffer.AddReferenceTimestampMeta(reference, gst.ClockTime(1234), gst.ClockTime(5)) == nil {
		t.Fatal("AddReferenceTimestampMeta returned nil")
	}

	meta := buffer.GetReferenceTimestampMeta(reference)
	if meta == nil {
		t.Fatal("GetReferenceTimestampMeta returned nil")
	}
	if meta.GetTimestamp() != 1234 || meta.GetDuration() != 5 {
		t.Fatalf("got timestamp %d and duration %d, want 1234 and 5", meta.GetTimestamp(), meta.GetDuration())
	}
	if !meta.GetReference().IsEqual(reference) {
		t.Fatalf("got reference %s, want %s", meta.GetReference().String(), reference.String())
	}

	if buffer.GetReferenceTimestampMeta(gst.CapsFromString("timestamp/x-other")) != nil {
		t.Fatal("got a meta for a reference that was never added")
	}

	if buf.Len() != 0 {
		t.Fatalf("meta accessors logged: %s", buf.String())
	}
}

func TestBufferIsWritable(t *testing.T) {
	gst.Init()

	buffer := gst.NewBuffer()
	if !buffer.IsWritable() {
		t.Fatal("a new buffer is not writable")
	}

	shared := gst.UnsafeBufferFromGlibNone(gst.UnsafeBufferToGlibNone(buffer))
	if buffer.IsWritable() {
		t.Fatal("a buffer with two references is writable")
	}
	runtime.KeepAlive(shared)
}
