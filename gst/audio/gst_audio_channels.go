package audio

/*
#include "gst.go.h"
*/
import "C"
import (
	"math"
	"strings"
	"unsafe"
)

// GetFallbackChannelMask gets the fallback channel-mask for the given number of channels.
func GetFallbackChannelMask(channels int) uint64 {
	return uint64(C.gst_audio_channel_get_fallback_mask(C.gint(channels)))
}

// ChannelPositionsFromMask converts the channels present in mask to a position array (which
// should have at least channels entries ensured by caller). If mask is set to 0, it is considered
// as 'not present' for purpose of conversion. A partially valid mask with less bits set than the
// number of channels is considered valid. This function returns nil if the arguments are invalid.
func ChannelPositionsFromMask(channels int, mask uint64) []ChannelPosition {
	var out C.GstAudioChannelPosition
	ret := C.gst_audio_channel_positions_from_mask(C.gint(channels), C.guint64(mask), &out)
	if !gobool(ret) {
		return nil
	}
	outsl := make([]ChannelPosition, channels)
	tmp := (*[(math.MaxInt32 - 1) / unsafe.Sizeof(C.GST_AUDIO_CHANNEL_POSITION_NONE)]C.GstAudioChannelPosition)(unsafe.Pointer(&out))[:channels:channels]
	for i, s := range tmp {
		outsl[i] = ChannelPosition(s)
	}
	return outsl
}

// ChannelPositionsToMask converts the given channel positions to a bitmask. If forceOrder is true
// it additionally checks the channels in the order required by GStreamer.
func ChannelPositionsToMask(positions []ChannelPosition, forceOrder bool) (mask uint64, ok bool) {
	var out C.guint64
	ok = gobool(C.gst_audio_channel_positions_to_mask(
		(*C.GstAudioChannelPosition)(unsafe.Pointer(&positions[0])),
		C.gint(len(positions)),
		gboolean(forceOrder),
		&out,
	))
	if ok {
		mask = uint64(out)
	}
	return
}

// ChannelPositionsToString converts the given positions into a human-readable string.
func ChannelPositionsToString(positions []ChannelPosition) string {
	ret := C.gst_audio_channel_positions_to_string(
		(*C.GstAudioChannelPosition)(unsafe.Pointer(&positions[0])),
		C.gint(len(positions)),
	)
	defer C.g_free((C.gpointer)(unsafe.Pointer(ret)))
	return C.GoString(ret)
}

// ChannelPositionsToValidOrder reorders the given positions from any order to the GStreamer order.
func ChannelPositionsToValidOrder(positions []ChannelPosition) (ok bool) {
	ok = gobool(C.gst_audio_channel_positions_to_valid_order(
		(*C.GstAudioChannelPosition)(unsafe.Pointer(&positions[0])),
		C.gint(len(positions)),
	))
	return
}

// AreValidChannelPositions checks if the positions are all valid. If forceOrder is true, it also checks
// if they are in the order required by GStreamer.
func AreValidChannelPositions(positions []ChannelPosition, forceOrder bool) bool {
	return gobool(C.gst_audio_check_valid_channel_positions(
		(*C.GstAudioChannelPosition)(unsafe.Pointer(&positions[0])),
		C.gint(len(positions)),
		gboolean(forceOrder),
	))
}

// ChannelPosition represents audio channel positions.
//
// These are the channels defined in SMPTE 2036-2-2008 Table 1 for 22.2 audio systems with the
// Surround and Wide channels from DTS Coherent Acoustics (v.1.3.1) and 10.2 and 7.1 layouts.
// In the caps the actual channel layout is expressed with a channel count and a channel mask,
// which describes the existing channels. The positions in the bit mask correspond to the enum
// values. For negotiation it is allowed to have more bits set in the channel mask than the
// number of channels to specify the allowed channel positions but this is not allowed in negotiated
// caps. It is not allowed in any situation other than the one mentioned below to have less
// bits set in the channel mask than the number of channels.
//
// ChannelPositionMono can only be used with a single mono channel that has no direction information
// and would be mixed into all directional channels. This is expressed in caps by having a single
// channel and no channel mask.
//
// ChannelPositionNone can only be used if all channels have this position. This is expressed in caps
// by having a channel mask with no bits set.
//
// As another special case it is allowed to have two channels without a channel mask. This implicitly
// means that this is a stereo stream with a front left and front right channel.
type ChannelPosition int

// ChannelPosition castings
const (
	ChannelPositionNone               ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_NONE                  // (-3) – used for position-less channels, e.g. from a sound card that records 1024 channels; mutually exclusive with any other channel position
	ChannelPositionMono               ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_MONO                  // (-2) – Mono without direction; can only be used with 1 channel
	ChannelPositionInvalid            ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_INVALID               // (-1) – invalid position
	ChannelPositionFrontLeft          ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_FRONT_LEFT            // (0) – Front left
	ChannelPositionFrontRight         ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_FRONT_RIGHT           // (1) – Front right
	ChannelPositionFrontCenter        ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_FRONT_CENTER          // (2) – Front center
	ChannelPositionLFE1               ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_LFE1                  // (3) – Low-frequency effects 1 (subwoofer)
	ChannelPositionRearLeft           ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_REAR_LEFT             // (4) – Rear left
	ChannelPositionRearRight          ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_REAR_RIGHT            // (5) – Rear right
	ChannelPositionFrontLeftOfCenter  ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_FRONT_LEFT_OF_CENTER  // (6) – Front left of center
	ChannelPositionFrontRightOfCenter ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_FRONT_RIGHT_OF_CENTER // (7) – Front right of center
	ChannelPositionRearCenter         ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_REAR_CENTER           // (8) – Rear center
	ChannelPositionLFE2               ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_LFE2                  // (9) – Low-frequency effects 2 (subwoofer)
	ChannelPositionSideLeft           ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_SIDE_LEFT             // (10) – Side left
	ChannelPositionSideRight          ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_SIDE_RIGHT            // (11) – Side right
	ChannelPositionTopFrontLeft       ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_FRONT_LEFT        // (12) – Top front left
	ChannelPositionTopFrontRight      ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_FRONT_RIGHT       // (13) – Top front right
	ChannelPositionTopFrontCenter     ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_FRONT_CENTER      // (14) – Top front center
	ChannelPositionTopCenter          ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_CENTER            // (15) – Top center
	ChannelPositionTopRearLeft        ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_REAR_LEFT         // (16) – Top rear left
	ChannelPositionTopRearRight       ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_REAR_RIGHT        // (17) – Top rear right
	ChannelPositionTopSideLeft        ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_SIDE_LEFT         // (18) – Top side right
	ChannelPositionTopSideRight       ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_SIDE_RIGHT        // (19) – Top rear right
	ChannelPositionTopRearCenter      ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_TOP_REAR_CENTER       // (20) – Top rear center
	ChannelPositionBottomFrontCenter  ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_BOTTOM_FRONT_CENTER   // (21) – Bottom front center
	ChannelPositionBottomFrontLeft    ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_BOTTOM_FRONT_LEFT     // (22) – Bottom front left
	ChannelPositionBottomFrontRight   ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_BOTTOM_FRONT_RIGHT    // (23) – Bottom front right
	ChannelPositionWideLeft           ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_WIDE_LEFT             // (24) – Wide left (between front left and side left)
	ChannelPositionWideRight          ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_WIDE_RIGHT            // (25) – Wide right (between front right and side right)
	ChannelPositionSurroundLeft       ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_SURROUND_LEFT         // (26) – Surround left (between rear left and side left)
	ChannelPositionSurroundRight      ChannelPosition = C.GST_AUDIO_CHANNEL_POSITION_SURROUND_RIGHT        // (27) – Surround right (between rear right and side right)
)

func (c ChannelPosition) String() string {
	// ugly hack
	return strings.TrimSuffix(strings.TrimPrefix(ChannelPositionsToString([]ChannelPosition{c}), "[ "), " ]")
}
