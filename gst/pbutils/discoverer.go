package pbutils

/*
#include <gst/pbutils/pbutils.h>
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

func init() {
	tm := []glib.TypeMarshaler{
		{
			T: glib.Type(C.gst_discoverer_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_object(uintptrToGVal(p))
				return &Discoverer{toGObject(unsafe.Pointer(c))}, nil
			},
		},
		{
			T: glib.Type(C.gst_discoverer_info_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_object(uintptrToGVal(p))
				return &DiscovererInfo{toGObject(unsafe.Pointer(c))}, nil
			},
		},
		{
			T: glib.Type(C.gst_discoverer_stream_info_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_object(uintptrToGVal(p))
				return &DiscovererStreamInfo{toGObject(unsafe.Pointer(c))}, nil
			},
		},
		{
			T: glib.Type(C.gst_discoverer_audio_info_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_object(uintptrToGVal(p))
				return &DiscovererAudioInfo{&DiscovererStreamInfo{toGObject(unsafe.Pointer(c))}}, nil
			},
		},
		{
			T: glib.Type(C.gst_discoverer_video_info_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_object(uintptrToGVal(p))
				return &DiscovererVideoInfo{&DiscovererStreamInfo{toGObject(unsafe.Pointer(c))}}, nil
			},
		},
		{
			T: glib.Type(C.gst_discoverer_subtitle_info_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_object(uintptrToGVal(p))
				return &DiscovererSubtitleInfo{&DiscovererStreamInfo{toGObject(unsafe.Pointer(c))}}, nil
			},
		},
		{
			T: glib.Type(C.gst_discoverer_container_info_get_type()),
			F: func(p uintptr) (interface{}, error) {
				c := C.g_value_get_object(uintptrToGVal(p))
				return &DiscovererContainerInfo{&DiscovererStreamInfo{toGObject(unsafe.Pointer(c))}}, nil
			},
		},
	}
	glib.RegisterGValueMarshalers(tm)
}

func uintptrToGVal(p uintptr) *C.GValue {
	return (*C.GValue)(unsafe.Pointer(p)) // vet thinks this is unsafe and there is no way around it for now.
	// but the given ptr is an address to a C object so go's concerns are misplaced.
}

func toGObject(o unsafe.Pointer) *glib.Object { return &glib.Object{GObject: glib.ToGObject(o)} }

// Discoverer represents a GstDiscoverer
type Discoverer struct{ *glib.Object }

func wrapDiscovererFull(d *C.GstDiscoverer) *Discoverer {
	return &Discoverer{glib.TransferFull(unsafe.Pointer(d))}
}

// NewDiscoverer creates a new Discoverer with the provided timeout.
func NewDiscoverer(timeout time.Duration) (*Discoverer, error) {
	initPbUtils()
	var gerr *C.GError
	var cTime C.GstClockTime
	if timeout < 0 {
		cTime = C.GstClockTime(gst.ClockTimeNone)
	} else {
		cTime = C.GstClockTime(timeout.Nanoseconds())
	}
	ret := C.gst_discoverer_new(C.GstClockTime(cTime), &gerr)
	if gerr != nil {
		return nil, wrapGerr(gerr)
	}
	return wrapDiscovererFull(ret), nil
}

// Instance returns the underlying GstDiscoverer instance.
func (d *Discoverer) Instance() *C.GstDiscoverer {
	return (*C.GstDiscoverer)(unsafe.Pointer(d.GObject))
}

// DiscoverURI synchronously discovers the given uri.
func (d *Discoverer) DiscoverURI(uri string) (*DiscovererInfo, error) {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	var err *C.GError
	info := C.gst_discoverer_discover_uri(d.Instance(), (*C.gchar)(unsafe.Pointer(curi)), &err)
	if err != nil {
		return nil, wrapGerr(err)
	}
	return wrapDiscovererInfoFull(info), nil
}

// DiscovererInfo represents a GstDiscovererInfo
type DiscovererInfo struct{ *glib.Object }

func wrapDiscovererInfoFull(d *C.GstDiscovererInfo) *DiscovererInfo {
	return &DiscovererInfo{glib.TransferFull(unsafe.Pointer(d))}
}

// Instance returns the underlying GstDiscovererInfo instance.
func (d *DiscovererInfo) Instance() *C.GstDiscovererInfo {
	return (*C.GstDiscovererInfo)(unsafe.Pointer(d.GObject))
}

// Copy creates a copy of this instance.
func (d *DiscovererInfo) Copy() *DiscovererInfo {
	return wrapDiscovererInfoFull(C.gst_discoverer_info_copy(d.Instance()))
}

// GetAudioStreams finds all the DiscovererAudioInfo contained in info.
func (d *DiscovererInfo) GetAudioStreams() []*DiscovererAudioInfo {
	gList := C.gst_discoverer_info_get_audio_streams(d.Instance())
	if gList == nil {
		return nil
	}
	return glistToAudioInfoSlice(gList)
}

// GetContainerStreams finds all the DiscovererContainerInfo contained in info.
func (d *DiscovererInfo) GetContainerStreams() []*DiscovererContainerInfo {
	gList := C.gst_discoverer_info_get_container_streams(d.Instance())
	if gList == nil {
		return nil
	}
	return glistToContainerInfoSlice(gList)
}

// GetDuration returns the durartion of the stream.
func (d *DiscovererInfo) GetDuration() time.Duration {
	dur := C.gst_discoverer_info_get_duration(d.Instance())
	return time.Duration(uint64(dur)) * time.Nanosecond
}

// GetLive returns whether this is a live stream.
func (d *DiscovererInfo) GetLive() bool {
	return gobool(C.gst_discoverer_info_get_live(d.Instance()))
}

// GetResult returns the result type.
func (d *DiscovererInfo) GetResult() DiscovererResult {
	return DiscovererResult(C.gst_discoverer_info_get_result(d.Instance()))
}

// GetSeekable returns whether the stream is seekable.
func (d *DiscovererInfo) GetSeekable() bool {
	return gobool(C.gst_discoverer_info_get_seekable(d.Instance()))
}

// GetStreamInfo returns the topology of the URI.
func (d *DiscovererInfo) GetStreamInfo() *DiscovererStreamInfo {
	info := C.gst_discoverer_info_get_stream_info(d.Instance())
	if info == nil {
		return nil
	}
	return wrapDiscovererStreamInfo(info)
}

// GetStreamList returns the list of all streams contained in the info.
func (d *DiscovererInfo) GetStreamList() []*DiscovererStreamInfo {
	gList := C.gst_discoverer_info_get_stream_list(d.Instance())
	if gList == nil {
		return nil
	}
	return glistToStreamInfoSlice(gList)
}

// GetSubtitleStreams returns the info about subtitle streams.
func (d *DiscovererInfo) GetSubtitleStreams() []*DiscovererSubtitleInfo {
	gList := C.gst_discoverer_info_get_subtitle_streams(d.Instance())
	if gList == nil {
		return nil
	}
	return glistToSubtitleInfoSlice(gList)
}

// GetTags retrieves the tag list for the URI stream.
func (d *DiscovererInfo) GetTags() *gst.TagList {
	tagList := C.gst_discoverer_info_get_tags(d.Instance())
	if tagList == nil {
		return nil
	}
	return gst.FromGstTagListUnsafeNone(unsafe.Pointer(tagList))
}

// GetTOC returns the TOC for the URI stream.
func (d *DiscovererInfo) GetTOC() *gst.TOC {
	toc := C.gst_discoverer_info_get_toc(d.Instance())
	if toc == nil {
		return nil
	}
	return gst.FromGstTOCUnsafeNone(unsafe.Pointer(toc))
}

// GetURI returns the URI for this info.
func (d *DiscovererInfo) GetURI() string {
	return C.GoString(C.gst_discoverer_info_get_uri(d.Instance()))
}

// GetVideoStreams finds all the DiscovererVideoInfo contained in info.
func (d *DiscovererInfo) GetVideoStreams() []*DiscovererVideoInfo {
	gList := C.gst_discoverer_info_get_video_streams(d.Instance())
	if gList == nil {
		return nil
	}
	return glistToVideoInfoSlice(gList)
}

// DiscovererStreamInfo is the base structure for information concerning a media stream.
type DiscovererStreamInfo struct{ *glib.Object }

func wrapDiscovererStreamInfo(d *C.GstDiscovererStreamInfo) *DiscovererStreamInfo {
	return &DiscovererStreamInfo{toGObject(unsafe.Pointer(d))}
}

// Instance returns the underlying GstDiscovererStreamInfo instance.
func (d *DiscovererStreamInfo) Instance() *C.GstDiscovererStreamInfo {
	return (*C.GstDiscovererStreamInfo)(unsafe.Pointer(d.GObject))
}

// GetCaps returns the caps from the stream info.
func (d *DiscovererStreamInfo) GetCaps() *gst.Caps {
	caps := C.gst_discoverer_stream_info_get_caps(d.Instance())
	if caps == nil {
		return nil
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// GetStreamID returns the stream ID of this stream.
func (d *DiscovererStreamInfo) GetStreamID() string {
	return C.GoString(C.gst_discoverer_stream_info_get_stream_id(d.Instance()))
}

// GetStreamTypeNick returns a human readable name for the stream type
func (d *DiscovererStreamInfo) GetStreamTypeNick() string {
	return C.GoString(C.gst_discoverer_stream_info_get_stream_type_nick(d.Instance()))
}

// GetTags gets the tags contained in this stream
func (d *DiscovererStreamInfo) GetTags() *gst.TagList {
	tagList := C.gst_discoverer_stream_info_get_tags(d.Instance())
	if tagList == nil {
		return nil
	}
	return gst.FromGstTagListUnsafeNone(unsafe.Pointer(tagList))
}

// GetTOC gets the TOC contained in this stream
func (d *DiscovererStreamInfo) GetTOC() *gst.TOC {
	toc := C.gst_discoverer_stream_info_get_toc(d.Instance())
	if toc == nil {
		return nil
	}
	return gst.FromGstTOCUnsafeNone(unsafe.Pointer(toc))
}

// DiscovererAudioInfo contains info specific to audio streams.
type DiscovererAudioInfo struct{ *DiscovererStreamInfo }

// Instance returns the underlying GstDiscovererAudioInfo instance.
func (d *DiscovererAudioInfo) Instance() *C.GstDiscovererAudioInfo {
	return (*C.GstDiscovererAudioInfo)(unsafe.Pointer(d.GObject))
}

// GetBitate returns the bitrate for the audio stream.
func (d *DiscovererAudioInfo) GetBitate() uint {
	return uint(C.gst_discoverer_audio_info_get_bitrate(d.Instance()))
}

// GetChannelMask returns the channel mask for the audio stream.
func (d *DiscovererAudioInfo) GetChannelMask() uint64 {
	return uint64(C.gst_discoverer_audio_info_get_channel_mask(d.Instance()))
}

// GetChannels returns the number of channels in the stream.
func (d *DiscovererAudioInfo) GetChannels() uint {
	return uint(C.gst_discoverer_audio_info_get_channels(d.Instance()))
}

// GetDepth returns the number of bits used per sample in each channel.
func (d *DiscovererAudioInfo) GetDepth() uint {
	return uint(C.gst_discoverer_audio_info_get_depth(d.Instance()))
}

// GetLanguage returns the language of the stream, or an empty string if unknown.
func (d *DiscovererAudioInfo) GetLanguage() string {
	lang := C.gst_discoverer_audio_info_get_language(d.Instance())
	if lang == nil {
		return ""
	}
	return C.GoString(lang)
}

// GetMaxBitrate returns the maximum bitrate of the stream in bits/second.
func (d *DiscovererAudioInfo) GetMaxBitrate() uint {
	return uint(C.gst_discoverer_audio_info_get_max_bitrate(d.Instance()))
}

// GetSampleRate returns the sample rate of the stream in Hertz.
func (d *DiscovererAudioInfo) GetSampleRate() uint {
	return uint(C.gst_discoverer_audio_info_get_sample_rate(d.Instance()))
}

// DiscovererVideoInfo contains info specific to video streams
type DiscovererVideoInfo struct{ *DiscovererStreamInfo }

// Instance returns the underlying GstDiscovererVideoInfo instance.
func (d *DiscovererVideoInfo) Instance() *C.GstDiscovererVideoInfo {
	return (*C.GstDiscovererVideoInfo)(unsafe.Pointer(d.GObject))
}

// GetBitrate returns the average or nominal bitrate of the video stream in bits/second.
func (d *DiscovererVideoInfo) GetBitrate() uint {
	return uint(C.gst_discoverer_video_info_get_bitrate(d.Instance()))
}

// GetDepth returns the depth in bits of the video stream.
func (d *DiscovererVideoInfo) GetDepth() uint {
	return uint(C.gst_discoverer_video_info_get_depth(d.Instance()))
}

// GetFramerateDenom returns the framerate of the video stream (denominator).
func (d *DiscovererVideoInfo) GetFramerateDenom() uint {
	return uint(C.gst_discoverer_video_info_get_framerate_denom(d.Instance()))
}

// GetFramerateNum returns the framerate of the video stream (numerator).
func (d *DiscovererVideoInfo) GetFramerateNum() uint {
	return uint(C.gst_discoverer_video_info_get_framerate_num(d.Instance()))
}

// GetHeight returns the height of the video stream in pixels.
func (d *DiscovererVideoInfo) GetHeight() uint {
	return uint(C.gst_discoverer_video_info_get_height(d.Instance()))
}

// GetMaxBitrate returns the maximum bitrate of the video stream in bits/second.
func (d *DiscovererVideoInfo) GetMaxBitrate() uint {
	return uint(C.gst_discoverer_video_info_get_max_bitrate(d.Instance()))
}

// GetPARDenom returns the Pixel Aspect Ratio (PAR) of the video stream (denominator).
func (d *DiscovererVideoInfo) GetPARDenom() uint {
	return uint(C.gst_discoverer_video_info_get_par_denom(d.Instance()))
}

// GetPARNum returns the Pixel Aspect Ratio (PAR) of the video stream (numerator).
func (d *DiscovererVideoInfo) GetPARNum() uint {
	return uint(C.gst_discoverer_video_info_get_par_num(d.Instance()))
}

// GetWidth returns the width of the video stream in pixels.
func (d *DiscovererVideoInfo) GetWidth() uint {
	return uint(C.gst_discoverer_video_info_get_width(d.Instance()))
}

// IsImage returns TRUE if the video stream corresponds to an image (i.e. only contains one frame).
func (d *DiscovererVideoInfo) IsImage() bool {
	return gobool(C.gst_discoverer_video_info_is_image(d.Instance()))
}

// IsInterlaced returns TRUE if the stream is interlaced.
func (d *DiscovererVideoInfo) IsInterlaced() bool {
	return gobool(C.gst_discoverer_video_info_is_interlaced(d.Instance()))
}

// DiscovererContainerInfo specific to container streams.
type DiscovererContainerInfo struct{ *DiscovererStreamInfo }

// Instance returns the underlying GstDiscovererContainerInfo instance.
func (d *DiscovererContainerInfo) Instance() *C.GstDiscovererContainerInfo {
	return (*C.GstDiscovererContainerInfo)(unsafe.Pointer(d.GObject))
}

// GetStreams returns the list of streams inside this container.
func (d *DiscovererContainerInfo) GetStreams() []*DiscovererStreamInfo {
	streams := C.gst_discoverer_container_info_get_streams(d.Instance())
	if streams == nil {
		return nil
	}
	return glistToStreamInfoSlice(streams)
}

// DiscovererSubtitleInfo contains info specific to subtitle streams
type DiscovererSubtitleInfo struct{ *DiscovererStreamInfo }

// Instance returns the underlying GstDiscovererSubtitleInfo instance.
func (d *DiscovererSubtitleInfo) Instance() *C.GstDiscovererSubtitleInfo {
	return (*C.GstDiscovererSubtitleInfo)(unsafe.Pointer(d.GObject))
}

// GetLanguage returns the language of the subtitles.
func (d *DiscovererSubtitleInfo) GetLanguage() string {
	lang := C.gst_discoverer_subtitle_info_get_language(d.Instance())
	if lang == nil {
		return ""
	}
	return C.GoString(lang)
}

func glistToStreamInfoSlice(glist *C.GList) []*DiscovererStreamInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := C.g_list_length(glist)
	out := make([]*DiscovererStreamInfo, int(l))
	for i := 0; i < int(l); i++ {
		data := C.g_list_nth_data(glist, C.guint(i))
		if data == nil {
			return out // safety
		}
		out[i] = &DiscovererStreamInfo{glib.TransferFull(unsafe.Pointer(data))}
	}
	return out
}

func glistToAudioInfoSlice(glist *C.GList) []*DiscovererAudioInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := C.g_list_length(glist)
	out := make([]*DiscovererAudioInfo, int(l))
	for i := 0; i < int(l); i++ {
		data := C.g_list_nth_data(glist, C.guint(i))
		if data == nil {
			return out // safety
		}
		out[i] = &DiscovererAudioInfo{&DiscovererStreamInfo{glib.TransferFull(unsafe.Pointer(data))}}
	}
	return out
}

func glistToVideoInfoSlice(glist *C.GList) []*DiscovererVideoInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := C.g_list_length(glist)
	out := make([]*DiscovererVideoInfo, 0)
	for i := 0; i < int(l); i++ {
		data := C.g_list_nth_data(glist, C.guint(i))
		if data == nil {
			return out // safety
		}
		out[i] = &DiscovererVideoInfo{&DiscovererStreamInfo{glib.TransferFull(unsafe.Pointer(data))}}
	}
	return out
}

func glistToContainerInfoSlice(glist *C.GList) []*DiscovererContainerInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := C.g_list_length(glist)
	out := make([]*DiscovererContainerInfo, 0)
	for i := 0; i < int(l); i++ {
		data := C.g_list_nth_data(glist, C.guint(i))
		if data == nil {
			return out // safety
		}
		out[i] = &DiscovererContainerInfo{&DiscovererStreamInfo{glib.TransferFull(unsafe.Pointer(data))}}
	}
	return out
}

func glistToSubtitleInfoSlice(glist *C.GList) []*DiscovererSubtitleInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := C.g_list_length(glist)
	out := make([]*DiscovererSubtitleInfo, 0)
	for i := 0; i < int(l); i++ {
		data := C.g_list_nth_data(glist, C.guint(i))
		if data == nil {
			return out // safety
		}
		out[i] = &DiscovererSubtitleInfo{&DiscovererStreamInfo{glib.TransferFull(unsafe.Pointer(data))}}
	}
	return out
}

// DiscovererResult casts a GstDiscovererResult
type DiscovererResult int

// Type castings
const (
	DiscovererResultOK             DiscovererResult = C.GST_DISCOVERER_OK              // (0) – The discovery was successful
	DiscovererResultURIInvalid     DiscovererResult = C.GST_DISCOVERER_URI_INVALID     // (1) – the URI is invalid
	DiscovererResultError          DiscovererResult = C.GST_DISCOVERER_ERROR           // (2) – an error happened and the GError is set
	DiscovererResultTimeout        DiscovererResult = C.GST_DISCOVERER_TIMEOUT         // (3) – the discovery timed-out
	DiscovererResultBusy           DiscovererResult = C.GST_DISCOVERER_BUSY            // (4) – the discoverer was already discovering a file
	DiscovererResultMissingPlugins DiscovererResult = C.GST_DISCOVERER_MISSING_PLUGINS // (5) – Some plugins are missing for full discovery
)

// DiscovererSerializeFlags casts GstDiscovererSerializeFlags.
type DiscovererSerializeFlags int

// Type castings
const (
	DiscovererSerializeBasic DiscovererSerializeFlags = C.GST_DISCOVERER_SERIALIZE_BASIC // (0) – Serialize only basic information, excluding caps, tags and miscellaneous information
	DiscovererSerializeCaps  DiscovererSerializeFlags = C.GST_DISCOVERER_SERIALIZE_CAPS  // (1) – Serialize the caps for each stream
	DiscovererSerializeTags  DiscovererSerializeFlags = C.GST_DISCOVERER_SERIALIZE_TAGS  // (2) – Serialize the tags for each stream
	DiscovererSerializeMisc  DiscovererSerializeFlags = C.GST_DISCOVERER_SERIALIZE_MISC  // (4) – Serialize miscellaneous information for each stream
	DiscovererSerializeAll   DiscovererSerializeFlags = C.GST_DISCOVERER_SERIALIZE_ALL   // (7) – Serialize all the available info, including caps, tags and miscellaneous
)
