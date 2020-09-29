package gst

// #include "gst.go.h"
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
	mapInfo := b.Map()
	if mapInfo.ptr == nil {
		return nil
	}
	defer mapInfo.Unmap()
	return mapInfo.Bytes()
}

// PresentationTimestamp returns the presentation timestamp of the buffer, or a negative duration
// if not known or relevant. This value contains the timestamp when the media should be
// presented to the user.
func (b *Buffer) PresentationTimestamp() time.Duration {
	pts := b.Instance().pts
	if ClockTime(pts) == ClockTimeNone {
		return time.Duration(-1)
	}
	return guint64ToDuration(pts)
}

// DecodingTimestamp returns the decoding timestamp of the buffer, or a negative duration if not known
// or relevant. This value contains the timestamp when the media should be processed.
func (b *Buffer) DecodingTimestamp() time.Duration {
	dts := b.Instance().dts
	if ClockTime(dts) == ClockTimeNone {
		return time.Duration(-1)
	}
	return guint64ToDuration(dts)
}

// Duration returns the length of the data inside this buffer, or a negative duration if not known
// or relevant.
func (b *Buffer) Duration() time.Duration {
	dur := b.Instance().duration
	if ClockTime(dur) == ClockTimeNone {
		return time.Duration(-1)
	}
	return guint64ToDuration(dur)
}

// Offset returns a media specific offset for the buffer data. For video frames, this is the frame
// number of this buffer. For audio samples, this is the offset of the first sample in this buffer.
// For file data or compressed data this is the byte offset of the first byte in this buffer.
func (b *Buffer) Offset() int64 { return int64(b.Instance().offset) }

// OffsetEnd returns the last offset contained in this buffer. It has the same format as Offset.
func (b *Buffer) OffsetEnd() int64 { return int64(b.Instance().offset_end) }

// Map will map the data inside this buffer.
func (b *Buffer) Map() *MapInfo {
	var mapInfo C.GstMapInfo
	C.gst_buffer_map(
		(*C.GstBuffer)(b.Instance()),
		(*C.GstMapInfo)(unsafe.Pointer(&mapInfo)),
		C.GST_MAP_READ,
	)
	return wrapMapInfo(&mapInfo, func() {
		C.gst_buffer_unmap(b.Instance(), (*C.GstMapInfo)(unsafe.Pointer(&mapInfo)))
	})
}
