package gstvideo_test

import (
	"encoding/binary"
	"runtime"
	"testing"
	"unsafe"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstvideo"
)

func isWritable(rect *gstvideo.VideoOverlayRectangle) bool {
	writable := gst.UnsafeMiniObjectFromGlibBorrow(gstvideo.UnsafeVideoOverlayRectangleToGlibNone(rect)).IsWritable()
	runtime.KeepAlive(rect)
	return writable
}

func newRectangle(t *testing.T) *gstvideo.VideoOverlayRectangle {
	gst.Init()

	const width, height = 4, 4
	// GST_VIDEO_OVERLAY_COMPOSITION_FORMAT_RGB depends on byte order
	format := gstvideo.VideoFormatArgb
	if binary.NativeEndian.Uint16([]byte{1, 0}) == 1 {
		format = gstvideo.VideoFormatBgra
	}

	pixels := gst.NewBufferAllocate(nil, width*height*4, nil)
	gstvideo.BufferAddVideoMeta(pixels, gstvideo.VideoFrameFlagNone, format, width, height)

	rect := gstvideo.NewVideoOverlayRectangleRaw(pixels, 0, 0, width, height, gstvideo.VideoOverlayFormatFlagNone)
	if rect == nil {
		t.Fatal("NewVideoOverlayRectangleRaw returned nil")
	}
	return rect
}

// A wrapper that does not take ownership must hold its own reference.
func TestVideoOverlayRectangleReference(t *testing.T) {
	rect := newRectangle(t)

	shared := gstvideo.UnsafeVideoOverlayRectangleFromGlibNone(gstvideo.UnsafeVideoOverlayRectangleToGlibNone(rect))
	if isWritable(rect) {
		t.Fatal("rectangle is writable while a second wrapper holds a reference")
	}

	gstvideo.UnsafeVideoOverlayRectangleUnref(shared)
	if !isWritable(rect) {
		t.Fatal("rectangle still has an extra reference after the second wrapper was released")
	}
}

// The composition must hand out its own rectangle and release it when the composition is released.
func TestVideoOverlayCompositionGetRectangle(t *testing.T) {
	rect := newRectangle(t)

	comp := gstvideo.NewVideoOverlayComposition(rect)
	got := comp.GetRectangle(0)
	if got == nil {
		t.Fatal("composition has no rectangle")
	}
	if unsafe.Pointer(gstvideo.UnsafeVideoOverlayRectangleToGlibNone(got)) != unsafe.Pointer(gstvideo.UnsafeVideoOverlayRectangleToGlibNone(rect)) {
		t.Fatal("GetRectangle returned a different rectangle than the one added")
	}

	gstvideo.UnsafeVideoOverlayCompositionUnref(comp)
	if !isWritable(rect) {
		t.Fatal("rectangle still has an extra reference after the composition was released")
	}
}
