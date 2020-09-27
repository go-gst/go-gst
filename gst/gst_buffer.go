package gst

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"
import (
	"bytes"
	"io"
	"io/ioutil"
	"time"
	"unsafe"
)

// Buffer is a go representation of a GstBuffer.
type Buffer struct {
	ptr *C.GstBuffer
}

// NewBufferFromBytes returns a new buffer from the given byte slice.
func NewBufferFromBytes(b []byte) *Buffer {
	str := string(b)
	p := unsafe.Pointer(C.CString(str))
	// memory is freed by gstreamer after building the new buffer
	buf := C.gst_buffer_new_wrapped((C.gpointer)(p), C.ulong(len(str)))
	return wrapBuffer(buf)
}

// NewBufferFromReader returns a new buffer from the given io.Reader.
func NewBufferFromReader(rdr io.Reader) (*Buffer, error) {
	out, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	return NewBufferFromBytes(out), nil
}

// Instance returns the underlying GstBuffer instance.
func (b *Buffer) Instance() *C.GstBuffer { return C.toGstBuffer(unsafe.Pointer(b.ptr)) }

// Reader returns an io.Reader for this buffer.
func (b *Buffer) Reader() io.Reader { return bytes.NewBuffer(b.Bytes()) }

// Bytes returns a byte slice of the data inside this buffer.
func (b *Buffer) Bytes() []byte {
	mapInfo := MapBuffer(b)
	defer mapInfo.Unmap()
	return C.GoBytes(mapInfo.Data, (C.int)(mapInfo.Size))
}

// PresentationTimestamp returns the presentation timestamp of the buffer, or a negative duration
// if not known or relevant. This value contains the timestamp when the media should be
// presented to the user.
func (b *Buffer) PresentationTimestamp() time.Duration {
	pts := b.Instance().pts
	if uint64(pts) == ClockTimeNone {
		return time.Duration(-1)
	}
	return nanosecondsToDuration(uint64(pts))
}

// DecodingTimestamp returns the decoding timestamp of the buffer, or a negative duration if not known
// or relevant. This value contains the timestamp when the media should be processed.
func (b *Buffer) DecodingTimestamp() time.Duration {
	dts := b.Instance().dts
	if uint64(dts) == ClockTimeNone {
		return time.Duration(-1)
	}
	return nanosecondsToDuration(uint64(dts))
}

// Duration returns the length of the data inside this buffer, or a negative duration if not known
// or relevant.
func (b *Buffer) Duration() time.Duration {
	dur := b.Instance().duration
	if uint64(dur) == ClockTimeNone {
		return time.Duration(-1)
	}
	return nanosecondsToDuration(uint64(dur))
}

// Offset returns a media specific offset for the buffer data. For video frames, this is the frame
// number of this buffer. For audio samples, this is the offset of the first sample in this buffer.
// For file data or compressed data this is the byte offset of the first byte in this buffer.
func (b *Buffer) Offset() int64 { return int64(b.Instance().offset) }

// OffsetEnd returns the last offset contained in this buffer. It has the same format as Offset.
func (b *Buffer) OffsetEnd() int64 { return int64(b.Instance().offset_end) }
