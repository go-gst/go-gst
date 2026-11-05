package gst_test

import (
	"testing"
	"unsafe"

	"github.com/go-gst/go-gst/pkg/gst"
)

// linkedPads returns an active source pad linked to an active sink pad, ready for buffers.
func linkedPads(t *testing.T) (src, sink gst.Pad) {
	t.Helper()
	gst.Init()

	src = gst.NewPad("src", gst.PadSrc)
	sink = gst.NewPad("sink", gst.PadSink)
	if ret := src.Link(sink); ret != gst.PadLinkOK {
		t.Fatalf("link failed: %v", ret)
	}
	if !src.SetActive(true) || !sink.SetActive(true) {
		t.Fatal("failed to activate pads")
	}
	segment := gst.NewSegment()
	segment.Init(gst.FormatTime)
	if !src.PushEvent(gst.NewEventStreamStart("test")) || !src.PushEvent(gst.NewEventSegment(segment)) {
		t.Fatal("failed to push stream-start and segment")
	}
	return src, sink
}

func isWritable(buf *gst.Buffer) bool {
	return gst.UnsafeMiniObjectFromGlibBorrow(gst.UnsafeBufferToGlibNone(buf)).IsWritable()
}

func newBuffer(fill uint8) *gst.Buffer {
	buf := gst.NewBufferAllocate(nil, 4, nil)
	buf.Memset(0, fill, 4)
	return buf
}

func firstByte(t *testing.T, buf *gst.Buffer) byte {
	t.Helper()
	m, ok := buf.Map(gst.MapRead)
	if !ok {
		t.Fatal("failed to map buffer")
	}
	defer m.Close()
	p := make([]byte, 1)
	if _, err := m.ReadAt(p, 0); err != nil {
		t.Fatal(err)
	}
	return p[0]
}

func TestPadProbeInfoBorrowBuffer(t *testing.T) {
	src, sink := linkedPads(t)

	var called, borrowedWritable bool
	var event *gst.Event
	src.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		called = true
		if buf := info.BorrowBuffer(); buf != nil {
			borrowedWritable = isWritable(buf)
		}
		event = info.BorrowEvent()
		return gst.PadProbeOK
	})
	sink.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, _ *gst.PadProbeInfo) gst.PadProbeReturn {
		return gst.PadProbeDrop
	})

	if ret := src.Push(newBuffer(0)); ret != gst.FlowOK {
		t.Fatalf("push failed: %v", ret)
	}
	if !called {
		t.Fatal("probe was not called")
	}
	if !borrowedWritable {
		t.Error("borrowed buffer is not writable although the pad holds the only reference")
	}
	if event != nil {
		t.Error("BorrowEvent returned an event in a buffer probe")
	}
}

func TestPadProbeInfoMakeBufferWritable(t *testing.T) {
	src, sink := linkedPads(t)

	shared := newBuffer(0)
	pushed := gst.UnsafeBufferFromGlibNone(gst.UnsafeBufferToGlibNone(shared))

	// Compared by address only, since the copy is freed once the sink drops it.
	var writable, received unsafe.Pointer
	var wasWritable bool
	var receivedByte byte
	src.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		buf := info.MakeBufferWritable()
		if buf == nil {
			return gst.PadProbeOK
		}
		writable = gst.UnsafeBufferToGlibNone(buf)
		wasWritable = isWritable(buf)
		if wasWritable {
			buf.Memset(0, 0xff, 4)
		}
		return gst.PadProbeOK
	})
	sink.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		buf := info.BorrowBuffer()
		received = gst.UnsafeBufferToGlibNone(buf)
		receivedByte = firstByte(t, buf)
		return gst.PadProbeDrop
	})

	if ret := src.Push(pushed); ret != gst.FlowOK {
		t.Fatalf("push failed: %v", ret)
	}
	if writable == nil {
		t.Fatal("MakeBufferWritable returned nil")
	}
	if !wasWritable {
		t.Fatal("MakeBufferWritable returned a buffer that is not writable")
	}
	if writable == gst.UnsafeBufferToGlibNone(shared) {
		t.Fatal("MakeBufferWritable returned the shared buffer instead of a copy")
	}
	if received != writable {
		t.Fatal("the sink pad did not receive the buffer returned by MakeBufferWritable")
	}
	if receivedByte != 0xff {
		t.Errorf("sink pad received byte %#x, want 0xff", receivedByte)
	}
	if b := firstByte(t, shared); b != 0 {
		t.Errorf("shared buffer was modified: byte %#x, want 0", b)
	}
}

func TestPadProbeInfoBorrowEvent(t *testing.T) {
	src, _ := linkedPads(t)

	var called bool
	var eventType gst.EventType
	var eventWritable bool
	var buffer, writable *gst.Buffer
	src.AddProbe(gst.PadProbeTypeEventDownstream, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		called = true
		if e := info.BorrowEvent(); e != nil {
			eventType = e.GetType()
			eventWritable = gst.UnsafeMiniObjectFromGlibBorrow(gst.UnsafeEventToGlibNone(e)).IsWritable()
		}
		buffer = info.BorrowBuffer()
		writable = info.MakeBufferWritable()
		return gst.PadProbeDrop
	})

	// Sticky events such as EOS are also stored on the pad, so use a non-sticky one.
	src.PushEvent(gst.NewEventCustom(gst.EventCustomDownstream, gst.NewStructureEmpty("test")))
	if !called {
		t.Fatal("probe was not called")
	}
	if eventType != gst.EventCustomDownstream {
		t.Errorf("BorrowEvent returned event type %v, want custom downstream", eventType)
	}
	if !eventWritable {
		t.Error("borrowed event is not writable although the pad holds the only reference")
	}
	if buffer != nil || writable != nil {
		t.Error("buffer accessors returned a buffer in an event probe")
	}
}
