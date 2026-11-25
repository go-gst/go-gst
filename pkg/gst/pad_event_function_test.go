package gst_test

import (
	"runtime"
	"testing"

	"github.com/go-gst/go-gst/pkg/gst"
)

// activePads returns a linked, active src/sink pair. The event function is installed
// before activation so it also sees stream-start and segment.
func activePads(t *testing.T, fn gst.PadEventFunction) (src, sink gst.Pad) {
	t.Helper()
	gst.Init()

	src = gst.NewPad("src", gst.PadSrc)
	sink = gst.NewPad("sink", gst.PadSink)
	sink.SetEventFunction(fn)
	if ret := src.Link(sink); ret != gst.PadLinkOK {
		t.Fatalf("link failed: %v", ret)
	}
	if !src.SetActive(true) || !sink.SetActive(true) {
		t.Fatal("failed to activate pads")
	}
	return src, sink
}

func TestPadSetEventFunctionOwnsEvent(t *testing.T) {
	var kept []*gst.Event
	var gotPad gst.Pad

	src, sink := activePads(t, func(pad gst.Pad, parent gst.Object, event *gst.Event) bool {
		gotPad = pad
		if parent != nil {
			t.Errorf("expected nil parent for unparented pad, got %v", parent)
		}
		// the callback owns event, so it may keep it beyond the call
		kept = append(kept, event)
		return true
	})

	segment := gst.NewSegment()
	segment.Init(gst.FormatTime)
	for _, ev := range []*gst.Event{gst.NewEventStreamStart("test"), gst.NewEventSegment(segment), gst.NewEventEOS()} {
		if !src.PushEvent(ev) {
			t.Fatal("PushEvent returned false")
		}
	}

	runtime.GC()
	runtime.GC()

	want := []gst.EventType{gst.EventStreamStart, gst.EventSegment, gst.EventEOS}
	if len(kept) != len(want) {
		t.Fatalf("got %d events, want %d", len(kept), len(want))
	}
	for i, ev := range kept {
		if ev.GetType() != want[i] {
			t.Errorf("event %d: got %v, want %v", i, ev.GetType(), want[i])
		}
	}
	if gotPad == nil || gotPad.GetName() != sink.GetName() {
		t.Errorf("callback received wrong pad: %v", gotPad)
	}
}

func TestPadSetEventFunctionReturnValue(t *testing.T) {
	src, _ := activePads(t, func(pad gst.Pad, parent gst.Object, event *gst.Event) bool {
		return event.GetType() != gst.EventEOS
	})

	if !src.PushEvent(gst.NewEventStreamStart("test")) {
		t.Fatal("stream-start should be accepted")
	}
	segment := gst.NewSegment()
	segment.Init(gst.FormatTime)
	if !src.PushEvent(gst.NewEventSegment(segment)) {
		t.Fatal("segment should be accepted")
	}
	if src.PushEvent(gst.NewEventEOS()) {
		t.Fatal("EOS should be rejected by the event function")
	}
}

func TestPadSetEventFunctionEventDefault(t *testing.T) {
	calls := 0
	src, _ := activePads(t, func(pad gst.Pad, parent gst.Object, event *gst.Event) bool {
		calls++
		// hand ownership back to GStreamer
		return pad.EventDefault(parent, event)
	})

	segment := gst.NewSegment()
	segment.Init(gst.FormatTime)
	src.PushEvent(gst.NewEventStreamStart("test"))
	src.PushEvent(gst.NewEventSegment(segment))
	src.PushEvent(gst.NewEventEOS())

	runtime.GC()
	runtime.GC()

	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}

func TestPadSetEventFunctionReplaceAndRemove(t *testing.T) {
	first, second := 0, 0
	src, sink := activePads(t, func(pad gst.Pad, parent gst.Object, event *gst.Event) bool {
		first++
		return true
	})

	src.PushEvent(gst.NewEventStreamStart("test"))

	sink.SetEventFunction(func(pad gst.Pad, parent gst.Object, event *gst.Event) bool {
		second++
		return true
	})
	segment := gst.NewSegment()
	segment.Init(gst.FormatTime)
	src.PushEvent(gst.NewEventSegment(segment))

	sink.SetEventFunction(nil)
	src.PushEvent(gst.NewEventEOS())

	if first != 1 || second != 1 {
		t.Fatalf("got first=%d second=%d, want 1 and 1", first, second)
	}
}
