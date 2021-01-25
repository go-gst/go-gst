package audio

/*
#include "gst.go.h"

gpointer audioBufferPlaneData(GstAudioBuffer * buffer, gint plane)
{
	return GST_AUDIO_BUFFER_PLANE_DATA(buffer, plane);
}

gint audioBufferPlaneSize(GstAudioBuffer * buffer)
{
	return GST_AUDIO_BUFFER_PLANE_SIZE(buffer);
}

*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// ClipBuffer will return a new buffer clipped to the given segment. The given buffer is no longer valid.
// The returned buffer may be nil if it is completely outside the configured segment.
func ClipBuffer(buffer *gst.Buffer, segment *gst.Segment, rate, bytesPerFrame int) *gst.Buffer {
	buf := C.gst_audio_buffer_clip(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Ref().Instance())),
		(*C.GstSegment)(unsafe.Pointer(segment.Instance())),
		C.gint(rate),
		C.gint(bytesPerFrame),
	)
	if buf == nil {
		return nil
	}
	return gst.FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

// ReorderChannels reorders the buffer against the given positions. The number of channels in each slice
// must be identical.
func ReorderChannels(buffer *gst.Buffer, format Format, from []ChannelPosition, to []ChannelPosition) bool {
	return gobool(C.gst_audio_buffer_reorder_channels(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
		C.GstAudioFormat(format),
		C.gint(len(from)),
		(*C.GstAudioChannelPosition)(unsafe.Pointer(&from[0])),
		(*C.GstAudioChannelPosition)(unsafe.Pointer(&to[0])),
	))
}

// TruncateBuffer truncates the buffer to finally have the given number of samples, removing the necessary
// amount of samples from the end and trim number of samples from the beginning. The original buffer is no
// longer valid. The returned buffer may be nil if the arguments were invalid.
func TruncateBuffer(buffer *gst.Buffer, bytesPerFrame int, trim, samples int64) *gst.Buffer {
	buf := C.gst_audio_buffer_truncate(
		(*C.GstBuffer)(unsafe.Pointer(buffer.Ref().Instance())),
		C.gint(bytesPerFrame),
		C.gsize(trim),
		C.gsize(samples),
	)
	if buf == nil {
		return nil
	}
	return gst.FromGstBufferUnsafeFull(unsafe.Pointer(buf))
}

// MapBuffer maps an audio gst.Buffer so that it can be read or written.
//
// This is especially useful when the gstbuffer is in non-interleaved (planar) layout, in which case this
// function will use the information in the gstbuffer's attached GstAudioMeta in order to map each channel
// in a separate "plane" in GstAudioBuffer. If a GstAudioMeta is not attached on the gstbuffer, then it must
// be in interleaved layout.
//
// If a GstAudioMeta is attached, then the GstAudioInfo on the meta is checked against info. Normally, they
// should be equal, but in case they are not, a g_critical will be printed and the GstAudioInfo from the meta
// will be used.
//
// In non-interleaved buffers, it is possible to have each channel on a separate GstMemory. In this case, each
// memory will be mapped separately to avoid copying their contents in a larger memory area. Do note though
// that it is not supported to have a single channel spanning over two or more different GstMemory objects.
// Although the map operation will likely succeed in this case, it will be highly sub-optimal and it is
// recommended to merge all the memories in the buffer before calling this function.
func MapBuffer(info *Info, buffer *gst.Buffer, flags gst.MapFlags) (*Buffer, bool) {
	var audioBuffer C.GstAudioBuffer
	ret := gobool(C.gst_audio_buffer_map(
		&audioBuffer,
		info.ptr,
		(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
		C.GstMapFlags(flags)),
	)
	if !ret {
		return nil, ret
	}
	return &Buffer{ptr: &audioBuffer}, ret
}

// Buffer is a structure containing the result of a MapBuffer operation. For non-interleaved (planar) buffers,
// the beginning of each channel in the buffer has its own pointer in the planes array. For interleaved
// buffers, the Planes slice only contains one item, which is the pointer to the beginning of the buffer,
// and NumPlanes equals 1.
//
// The different channels in planes are always in the GStreamer channel order.
type Buffer struct {
	ptr *C.GstAudioBuffer
}

// NumSamples returns the size of the buffer in samples.
func (b *Buffer) NumSamples() int64 { return int64(b.ptr.n_samples) }

// NumPlanes returns the number of available planes.
func (b *Buffer) NumPlanes() int { return int(b.ptr.n_planes) }

// Planes returns the planes inside the mapped buffer.
func (b *Buffer) Planes() [][]byte {
	out := make([][]byte, b.NumPlanes())
	for i := 0; i < b.NumPlanes(); i++ {
		out[i] = C.GoBytes(
			unsafe.Pointer(C.audioBufferPlaneData(b.ptr, C.gint(i))),
			C.audioBufferPlaneSize(b.ptr),
		)
	}
	return out
}

// WritePlane writes data to the plane at the given index.
func (b *Buffer) WritePlane(idx int, data []byte) {
	bufdata := C.audioBufferPlaneData(b.ptr, C.gint(idx))
	C.memcpy(unsafe.Pointer(bufdata), unsafe.Pointer(&data[0]), C.gsize(len(data)))
}

// Unmap will unmap the mapped buffer. Use this after calling MapBuffer.
func (b *Buffer) Unmap() { C.gst_audio_buffer_unmap(b.ptr) }
