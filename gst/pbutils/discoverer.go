package pbutils

import (
	"C"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

// #include <gst/pbutils/pbutils.h>
import "C"

import (
	"unsafe"
)

// Discoverer represents a GstDiscoverer
type Discoverer struct{ *glib.Object }

func wrapDiscoverer(d *C.GstDiscoverer) *Discoverer {
	return &Discoverer{Object: &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(d))}}
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
	return wrapDiscoverer(ret), nil
}

// Instance returns the underlying GstDiscoverer instance.
func (d *Discoverer) Instance() *C.GstDiscoverer {
	return (*C.GstDiscoverer)(unsafe.Pointer(d.Native()))
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
	return wrapDiscovererInfo(info), nil
}

// DiscovererInfo represents a GstDiscovererInfo
type DiscovererInfo struct{ *glib.Object }

func wrapDiscovererInfo(d *C.GstDiscovererInfo) *DiscovererInfo {
	return &DiscovererInfo{Object: &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(d))}}
}

// Instance returns the underlying GstDiscovererInfo instance.
func (d *DiscovererInfo) Instance() *C.GstDiscovererInfo {
	return (*C.GstDiscovererInfo)(unsafe.Pointer(d.Native()))
}

// Copy creates a copy of this instance.
func (d *DiscovererInfo) Copy() *DiscovererInfo {
	return wrapDiscovererInfo(C.gst_discoverer_info_copy(d.Instance()))
}

// GetAudioStreams finds all the DiscovererAudioInfo contained in info.
func (d *DiscovererInfo) GetAudioStreams() []*DiscovererAudioInfo {
	gList := C.gst_discoverer_info_get_audio_streams(d.Instance())
	return glistToAudioInfoSlice(gList)
}

// GetContainerStreams finds all the DiscovererContainerInfo contained in info.
func (d *DiscovererInfo) GetContainerStreams() []*DiscovererContainerInfo {
	gList := C.gst_discoverer_info_get_container_streams(d.Instance())
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
	return wrapDiscovererStreamInfo(C.gst_discoverer_info_get_stream_info(d.Instance()))
}

// GetStreamList returns the list of all streams contained in the info.
func (d *DiscovererInfo) GetStreamList() []*DiscovererStreamInfo {
	return glistToStreamInfoSlice(C.gst_discoverer_info_get_stream_list(d.Instance()))
}

// GetSubtitleStreams returns the info about subtitle streams.
func (d *DiscovererInfo) GetSubtitleStreams() []*DiscovererSubtitleInfo {
	gList := C.gst_discoverer_info_get_subtitle_streams(d.Instance())
	return glistToSubtitleInfoSlice(gList)
}

// DiscovererStreamInfo is the base structure for information concerning a media stream.
type DiscovererStreamInfo struct{ *glib.Object }

func wrapDiscovererStreamInfo(d *C.GstDiscovererStreamInfo) *DiscovererStreamInfo {
	return &DiscovererStreamInfo{Object: &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(d))}}
}

// Instance returns the underlying GstDiscovererStreamInfo instance.
func (d *DiscovererStreamInfo) Instance() *C.GstDiscovererStreamInfo {
	return (*C.GstDiscovererStreamInfo)(unsafe.Pointer(d.Native()))
}

// DiscovererAudioInfo contains info specific to audio streams.
type DiscovererAudioInfo struct{ *DiscovererStreamInfo }

func wrapDiscovererAudioInfo(d *C.GstDiscovererAudioInfo) *DiscovererAudioInfo {
	return &DiscovererAudioInfo{&DiscovererStreamInfo{Object: &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(d))}}}
}

// Instance returns the underlying GstDiscovererAudioInfo instance.
func (d *DiscovererAudioInfo) Instance() *C.GstDiscovererAudioInfo {
	return (*C.GstDiscovererAudioInfo)(unsafe.Pointer(d.Native()))
}

// DiscovererVideoInfo contains info specific to video streams
type DiscovererVideoInfo struct{ *DiscovererStreamInfo }

func wrapDiscovererVideoInfo(d *C.GstDiscovererVideoInfo) *DiscovererVideoInfo {
	return &DiscovererVideoInfo{&DiscovererStreamInfo{Object: &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(d))}}}
}

// Instance returns the underlying GstDiscovererVideoInfo instance.
func (d *DiscovererVideoInfo) Instance() *C.GstDiscovererVideoInfo {
	return (*C.GstDiscovererVideoInfo)(unsafe.Pointer(d.Native()))
}

// DiscovererContainerInfo specific to container streams.
type DiscovererContainerInfo struct{ *DiscovererStreamInfo }

func wrapDiscovererContainerInfo(d *C.GstDiscovererContainerInfo) *DiscovererContainerInfo {
	return &DiscovererContainerInfo{&DiscovererStreamInfo{Object: &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(d))}}}
}

// Instance returns the underlying GstDiscovererContainerInfo instance.
func (d *DiscovererContainerInfo) Instance() *C.GstDiscovererContainerInfo {
	return (*C.GstDiscovererContainerInfo)(unsafe.Pointer(d.Native()))
}

// DiscovererSubtitleInfo contains info specific to subtitle streams
type DiscovererSubtitleInfo struct{ *DiscovererStreamInfo }

func wrapDiscovererSubtitleInfo(d *C.GstDiscovererSubtitleInfo) *DiscovererSubtitleInfo {
	return &DiscovererSubtitleInfo{&DiscovererStreamInfo{Object: &glib.Object{GObject: glib.ToGObject(unsafe.Pointer(d))}}}
}

// Instance returns the underlying GstDiscovererSubtitleInfo instance.
func (d *DiscovererSubtitleInfo) Instance() *C.GstDiscovererSubtitleInfo {
	return (*C.GstDiscovererSubtitleInfo)(unsafe.Pointer(d.Native()))
}

func glistToStreamInfoSlice(glist *C.GList) []*DiscovererStreamInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	out := make([]*DiscovererStreamInfo, 0)
	l.Foreach(func(item interface{}) {
		st := item.(*C.GstDiscovererStreamInfo)
		out = append(out, wrapDiscovererStreamInfo(st))
	})
	return out
}

func glistToAudioInfoSlice(glist *C.GList) []*DiscovererAudioInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	out := make([]*DiscovererAudioInfo, 0)
	l.Foreach(func(item interface{}) {
		st := item.(*C.GstDiscovererAudioInfo)
		out = append(out, wrapDiscovererAudioInfo(st))
	})
	return out
}

func glistToVideoInfoSlice(glist *C.GList) []*DiscovererVideoInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	out := make([]*DiscovererVideoInfo, 0)
	l.Foreach(func(item interface{}) {
		st := item.(*C.GstDiscovererVideoInfo)
		out = append(out, wrapDiscovererVideoInfo(st))
	})
	return out
}

func glistToContainerInfoSlice(glist *C.GList) []*DiscovererContainerInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	out := make([]*DiscovererContainerInfo, 0)
	l.Foreach(func(item interface{}) {
		st := item.(*C.GstDiscovererContainerInfo)
		out = append(out, wrapDiscovererContainerInfo(st))
	})
	return out
}

func glistToSubtitleInfoSlice(glist *C.GList) []*DiscovererSubtitleInfo {
	defer C.gst_discoverer_stream_info_list_free(glist)
	l := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	out := make([]*DiscovererSubtitleInfo, 0)
	l.Foreach(func(item interface{}) {
		st := item.(*C.GstDiscovererSubtitleInfo)
		out = append(out, wrapDiscovererSubtitleInfo(st))
	})
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
