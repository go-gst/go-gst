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

// newCustomEvent returns a non-sticky event, so that the pad holds no reference to it.
func newCustomEvent(name string) *gst.Event {
	return gst.NewEventCustom(gst.EventCustomDownstream, gst.NewStructureEmpty(name))
}

// writeThroughProbe pushes buf through a probe that takes the buffer, copies it if it is not
// writable, sets its first byte to 0xff and puts it back. It returns whether the taken buffer was
// writable, and the first byte the sink pad received.
func writeThroughProbe(t *testing.T, buf *gst.Buffer) (takenWritable bool, received byte) {
	t.Helper()
	src, sink := linkedPads(t)

	src.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		b := info.TakeBuffer()
		if b == nil {
			t.Error("TakeBuffer returned nil")
			return gst.PadProbeOK
		}
		takenWritable = b.IsWritable()
		if !takenWritable {
			b = b.Copy()
		}
		b.Memset(0, 0xff, 1)
		info.SetBuffer(b)
		return gst.PadProbeOK
	})
	var called bool
	sink.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		called = true
		received = firstByte(t, info.GetBuffer())
		return gst.PadProbeDrop
	})

	if ret := src.Push(buf); ret != gst.FlowOK {
		t.Fatalf("push failed: %v", ret)
	}
	if !called {
		t.Fatal("the sink pad did not receive a buffer")
	}
	return takenWritable, received
}

func TestPadProbeInfoGetBufferBorrows(t *testing.T) {
	src, sink := linkedPads(t)

	var called, writable bool
	src.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		called = true
		writable = info.GetBuffer().IsWritable()
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
	if !writable {
		t.Error("buffer from GetBuffer is not writable although the pad holds the only reference")
	}
}

func TestPadProbeInfoGetEventBorrows(t *testing.T) {
	src, _ := linkedPads(t)

	var called, writable bool
	src.AddProbe(gst.PadProbeTypeEventDownstream, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		called = true
		e := info.GetEvent()
		writable = gst.UnsafeMiniObjectFromGlibBorrow(gst.UnsafeEventToGlibNone(e)).IsWritable()
		return gst.PadProbeDrop
	})

	src.PushEvent(newCustomEvent("test"))
	if !called {
		t.Fatal("probe was not called")
	}
	if !writable {
		t.Error("event from GetEvent is not writable although the pad holds the only reference")
	}
}

func TestPadProbeInfoTakeAndSetBuffer(t *testing.T) {
	t.Run("sole owner", func(t *testing.T) {
		writable, received := writeThroughProbe(t, newBuffer(0))
		if !writable {
			t.Error("taken buffer is not writable although the pad held the only reference")
		}
		if received != 0xff {
			t.Errorf("sink pad received byte %#x, want 0xff", received)
		}
	})

	t.Run("shared", func(t *testing.T) {
		shared := newBuffer(0)
		pushed := gst.UnsafeBufferFromGlibNone(gst.UnsafeBufferToGlibNone(shared))

		writable, received := writeThroughProbe(t, pushed)
		if writable {
			t.Error("taken buffer is writable although another reference exists")
		}
		if received != 0xff {
			t.Errorf("sink pad received byte %#x, want 0xff", received)
		}
		if b := firstByte(t, shared); b != 0 {
			t.Errorf("shared buffer was modified: byte %#x, want 0", b)
		}
	})
}

func TestPadProbeInfoTakeBufferDrops(t *testing.T) {
	src, sink := linkedPads(t)

	var taken *gst.Buffer
	var otherTypes bool
	src.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		otherTypes = info.TakeEvent() != nil || info.TakeBufferList() != nil
		taken = info.TakeBuffer()
		return gst.PadProbeOK
	})
	var sinkCalled bool
	sink.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, _ *gst.PadProbeInfo) gst.PadProbeReturn {
		sinkCalled = true
		return gst.PadProbeDrop
	})

	if ret := src.Push(newBuffer(0x42)); ret != gst.FlowOK {
		t.Fatalf("push failed: %v", ret)
	}
	if otherTypes {
		t.Error("TakeEvent or TakeBufferList returned data in a buffer probe")
	}
	if taken == nil {
		t.Fatal("TakeBuffer returned nil")
	}
	if sinkCalled {
		t.Error("the sink pad received a buffer although the probe took it")
	}
	if !taken.IsWritable() {
		t.Error("taken buffer is not writable after the push, so it is not owned by the caller alone")
	}
	if b := firstByte(t, taken); b != 0x42 {
		t.Errorf("taken buffer has byte %#x, want 0x42", b)
	}
}

func TestPadProbeInfoSetBufferWithBorrowedBuffer(t *testing.T) {
	src, sink := linkedPads(t)

	pushed := newBuffer(0x42)
	ptr := gst.UnsafeBufferToGlibNone(pushed)
	src.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		info.SetBuffer(info.GetBuffer())
		return gst.PadProbeOK
	})
	var received unsafe.Pointer
	var writable bool
	sink.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		buf := info.GetBuffer()
		received = gst.UnsafeBufferToGlibNone(buf)
		writable = buf.IsWritable()
		return gst.PadProbeDrop
	})

	if ret := src.Push(pushed); ret != gst.FlowOK {
		t.Fatalf("push failed: %v", ret)
	}
	if received != ptr {
		t.Fatal("the sink pad did not receive the pushed buffer")
	}
	if !writable {
		t.Error("buffer is not writable at the sink pad, so its reference count changed")
	}
}

func TestPadProbeInfoTakeAndSetBufferList(t *testing.T) {
	src, sink := linkedPads(t)

	src.AddProbe(gst.PadProbeTypeBufferList, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		list := info.TakeBufferList()
		if list == nil {
			t.Error("TakeBufferList returned nil")
			return gst.PadProbeOK
		}
		list.Insert(-1, newBuffer(2))
		info.SetBufferList(list)
		return gst.PadProbeOK
	})
	var length uint32
	sink.AddProbe(gst.PadProbeTypeBufferList, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		length = info.GetBufferList().Length()
		return gst.PadProbeDrop
	})

	list := gst.NewBufferList()
	list.Insert(-1, newBuffer(1))
	if ret := src.PushList(list); ret != gst.FlowOK {
		t.Fatalf("push failed: %v", ret)
	}
	if length != 2 {
		t.Errorf("sink pad received a list of %d buffers, want 2", length)
	}
}

func TestPadProbeInfoTakeAndSetEvent(t *testing.T) {
	src, sink := linkedPads(t)

	src.AddProbe(gst.PadProbeTypeEventDownstream, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		if info.TakeEvent() == nil {
			t.Error("TakeEvent returned nil")
		}
		info.SetEvent(newCustomEvent("replaced"))
		return gst.PadProbeOK
	})
	var name string
	sink.AddProbe(gst.PadProbeTypeEventDownstream, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		name = info.GetEvent().GetStructure().GetName()
		return gst.PadProbeDrop
	})

	src.PushEvent(newCustomEvent("original"))
	if name != "replaced" {
		t.Errorf("sink pad received event %q, want %q", name, "replaced")
	}
}

func TestPadProbeInfoSetBufferReleasesOld(t *testing.T) {
	src, sink := linkedPads(t)

	kept := newBuffer(0)
	pushed := gst.UnsafeBufferFromGlibNone(gst.UnsafeBufferToGlibNone(kept))
	src.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		info.SetBuffer(newBuffer(1))
		return gst.PadProbeOK
	})
	sink.AddProbe(gst.PadProbeTypeBuffer, func(_ gst.Pad, _ *gst.PadProbeInfo) gst.PadProbeReturn {
		return gst.PadProbeDrop
	})

	if ret := src.Push(pushed); ret != gst.FlowOK {
		t.Fatalf("push failed: %v", ret)
	}
	if !kept.IsWritable() {
		t.Error("the replaced buffer still has another reference after the push")
	}
}
