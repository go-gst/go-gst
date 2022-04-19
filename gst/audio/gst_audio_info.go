package audio

/*
#include "gst.go.h"

gint channels(GstAudioInfo * info)
{
	return GST_AUDIO_INFO_CHANNELS(info);
}
*/
import "C"
import (
	"math"
	"runtime"
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// Flags contains extra audio flags
type Flags int

// Flags castings
const (
	FlagNone         Flags = C.GST_AUDIO_FLAG_NONE         // (0) - no valid flag
	FlagUnpositioned Flags = C.GST_AUDIO_FLAG_UNPOSITIONED // (1) – the position array explicitly contains unpositioned channels.
)

func wrapInfoFull(ptr *C.GstAudioInfo) *Info {
	info := &Info{ptr}
	runtime.SetFinalizer(info, (*Info).Free)
	return info
}

// Info is a structure used for describing audio properties. This can be filled in from caps
// or coverted back to caps.
type Info struct {
	ptr *C.GstAudioInfo
}

// NewInfo returns a new Info that is also initialized.
func NewInfo() *Info { return wrapInfoFull(C.gst_audio_info_new()) }

// InfoFromCaps parses the provided caps and creates an info. It returns true if the caps could be parsed.
func InfoFromCaps(caps *gst.Caps) (*Info, bool) {
	info := NewInfo()
	return info, gobool(C.gst_audio_info_from_caps(info.ptr, (*C.GstCaps)(unsafe.Pointer(caps.Instance()))))
}

// FormatInfo returns the format info for the audio.
func (i *Info) FormatInfo() *FormatInfo { return &FormatInfo{i.ptr.finfo} }

// Flags returns additional flags for the audio.
func (i *Info) Flags() Flags { return Flags(i.ptr.flags) }

// Layout returns the audio layout.
func (i *Info) Layout() Layout { return Layout(i.ptr.layout) }

// Rate returns the audio sample rate.
func (i *Info) Rate() int { return int(i.ptr.rate) }

// Channels returns the number of channels.
func (i *Info) Channels() int { return int(C.channels(i.ptr)) }

// BPF returns the number of bytes for one frame. This is the size of one sample * Channels.
func (i *Info) BPF() int { return int(i.ptr.bpf) }

// Positions returns the positions for each channel.
func (i *Info) Positions() []ChannelPosition {
	l := i.Channels()
	out := make([]ChannelPosition, int(l))
	tmpslice := (*[(math.MaxInt32 - 1) / unsafe.Sizeof(C.GST_AUDIO_CHANNEL_POSITION_NONE)]C.GstAudioChannelPosition)(unsafe.Pointer(&i.ptr.position))[:l:l]
	for i, s := range tmpslice {
		out[i] = ChannelPosition(s)
	}
	return out
}

// Init initializes the info with the default values.
func (i *Info) Init() { C.gst_audio_info_init(i.ptr) }

// Free frees the AudioInfo structure. This is usually handled for you by the bindings.
func (i *Info) Free() { C.gst_audio_info_free(i.ptr) }

// Convert converts among various gst.Format types. This function handles gst.FormatBytes, gst.FormatTime,
// and gst.FormatDefault. For raw audio, gst.FormatDefault corresponds to audio frames. This function can
// be used to handle pad queries of the type gst.QueryConvert. To provide a value from a time.Duration, use the
// Nanoseconds() method.
func (i *Info) Convert(srcFmt gst.Format, srcVal int64, destFmt gst.Format) (int64, bool) {
	var out C.gint64
	ret := C.gst_audio_info_convert(
		i.ptr,
		C.GstFormat(srcFmt),
		C.gint64(srcVal),
		C.GstFormat(destFmt),
		&out,
	)
	return int64(out), gobool(ret)
}

// Copy creates a copy of this Info.
func (i *Info) Copy() *Info {
	return wrapInfoFull(C.gst_audio_info_copy(i.ptr))
}

// IsEqual checks if the two infos are equal.
func (i *Info) IsEqual(info *Info) bool {
	return gobool(C.gst_audio_info_is_equal(i.ptr, info.ptr))
}

// SetFormat sets the format for this info. This initializes info first and no values are preserved.
func (i *Info) SetFormat(format Format, rate int, positions []ChannelPosition) {
	C.gst_audio_info_set_format(
		i.ptr,
		C.GstAudioFormat(format),
		C.gint(rate),
		C.gint(len(positions)),
		(*C.GstAudioChannelPosition)(unsafe.Pointer(&positions[0])),
	)
}

// ToCaps returns the caps representation of this info.
func (i *Info) ToCaps() *gst.Caps {
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_audio_info_to_caps(i.ptr)))
}
