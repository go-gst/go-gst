package audio

// /*
// #include "gst.go.h"
// */
// import "C"
// import (
// 	"runtime"
// 	"unsafe"
// )

// // ChannelMixer is a structure for mixing audio channels.
// type ChannelMixer struct {
// 	ptr *C.GstAudioChannelMixer
// }

// // NewChannelMixer creates a new channel mixer with the given parameters.
// // It returns nil if the format is not supported.
// func NewChannelMixer(flags ChannelMixerFlags, format Format, inPositions, outPositions []ChannelPosition) *ChannelMixer {
// 	mixer := C.gst_audio_channel_mixer_new(
// 		C.GstAudioChannelMixerFlags(flags),
// 		C.GstAudioFormat(format),
// 		C.gint(len(inPositions)),
// 		(*C.GstAudioChannelPosition)(unsafe.Pointer(&inPositions[0])),
// 		C.gint(len(outPositions)),
// 		(*C.GstAudioChannelPosition)(unsafe.Pointer(&outPositions[0])),
// 	)
// 	if mixer == nil {
// 		return nil
// 	}
// 	wrapped := &ChannelMixer{mixer}
// 	runtime.SetFinalizer(wrapped, (*ChannelMixer).Free)
// 	return wrapped
// }

// // MixSamples will mix the given samples.
// //
// // In case the samples are interleaved, in must be a slice with a single element of interleaved samples.
// //
// // If non-interleaved samples are used, in must be a slice of data with an element for each channel

// // Free frees the channel mixer. The bindings will usually take care of this for you.
// func (c *ChannelMixer) Free() { C.gst_audio_channel_mixer_free(c.ptr) }

// // ChannelMixerFlags are flags used to configure a ChannelMixer
// type ChannelMixerFlags int

// // Castings for ChannelMixerFlags
// const (
// 	ChannelMixerFlagsNone              ChannelMixerFlags = C.GST_AUDIO_CHANNEL_MIXER_FLAGS_NONE                // (0) – no flag
// 	ChannelMixerFlagsNonInterleavedIn  ChannelMixerFlags = C.GST_AUDIO_CHANNEL_MIXER_FLAGS_NON_INTERLEAVED_IN  // (1) – input channels are not interleaved
// 	ChannelMixerFlagsNonInterleavedOut ChannelMixerFlags = C.GST_AUDIO_CHANNEL_MIXER_FLAGS_NON_INTERLEAVED_OUT // (2) – output channels are not interleaved
// 	ChannelMixerFlagsUnpositionedIn    ChannelMixerFlags = C.GST_AUDIO_CHANNEL_MIXER_FLAGS_UNPOSITIONED_IN     // (4) – input channels are explicitly unpositioned
// 	ChannelMixerFlagsUnpositionedOut   ChannelMixerFlags = C.GST_AUDIO_CHANNEL_MIXER_FLAGS_UNPOSITIONED_OUT    // (8) – output channels are explicitly unpositioned
// )
