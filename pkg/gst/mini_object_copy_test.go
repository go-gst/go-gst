package gst_test

import (
	"testing"
	"unsafe"

	"github.com/go-gst/go-gst/pkg/gst"
)

func miniObjectIsWritable(p unsafe.Pointer) bool {
	return gst.UnsafeMiniObjectFromGlibBorrow(p).IsWritable()
}

func TestBufferCopy(t *testing.T) {
	gst.Init()

	orig := newBuffer(0)
	c := orig.Copy()
	if gst.UnsafeBufferToGlibNone(c) == gst.UnsafeBufferToGlibNone(orig) {
		t.Fatal("Copy returned the same buffer")
	}
	if !c.IsWritable() {
		t.Fatal("copy is not writable")
	}
	c.Memset(0, 0xff, 1)
	if b := firstByte(t, orig); b != 0 {
		t.Errorf("writing to the copy changed the original: byte %#x, want 0", b)
	}
}

func TestBufferListCopy(t *testing.T) {
	gst.Init()

	orig := gst.NewBufferList()
	orig.Insert(-1, newBuffer(0))
	c := orig.Copy()
	p := gst.UnsafeBufferListToGlibNone(c)
	if p == gst.UnsafeBufferListToGlibNone(orig) {
		t.Fatal("Copy returned the same list")
	}
	if !miniObjectIsWritable(p) {
		t.Fatal("copy is not writable")
	}
	c.Insert(-1, newBuffer(1))
	if orig.Length() != 1 || c.Length() != 2 {
		t.Errorf("lengths are %d and %d, want 1 and 2", orig.Length(), c.Length())
	}
}

func TestEventCopy(t *testing.T) {
	gst.Init()

	orig := newCustomEvent("test")
	c := orig.Copy()
	p := gst.UnsafeEventToGlibNone(c)
	if p == gst.UnsafeEventToGlibNone(orig) {
		t.Fatal("Copy returned the same event")
	}
	if !miniObjectIsWritable(p) {
		t.Fatal("copy is not writable")
	}
	if name := c.GetStructure().GetName(); name != "test" {
		t.Errorf("copy has structure %q, want %q", name, "test")
	}
}

func TestQueryCopy(t *testing.T) {
	gst.Init()

	orig := gst.NewQueryDuration(gst.FormatTime)
	c := orig.Copy()
	p := gst.UnsafeQueryToGlibNone(c)
	if p == gst.UnsafeQueryToGlibNone(orig) {
		t.Fatal("Copy returned the same query")
	}
	if !miniObjectIsWritable(p) {
		t.Fatal("copy is not writable")
	}
	if format, _ := c.ParseDuration(); format != gst.FormatTime {
		t.Errorf("copy has format %v, want time", format)
	}
}
