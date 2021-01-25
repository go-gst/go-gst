package rtp

/*
#include "gst.go.h"
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

// MaxCSRCCount is the maximum number of elements that can be added to a CSRC.
const MaxCSRCCount uint = C.GST_RTP_SOURCE_META_MAX_CSRC_COUNT // 15

// SourceMeta is a wrapper around GstRTPSourceMeta.
type SourceMeta struct {
	ptr *C.GstRTPSourceMeta
}

// SourceMetaAPIType returns the GType for GstRTPSourceMeta.
func SourceMetaAPIType() glib.Type {
	return glib.Type(C.gst_rtp_source_meta_api_get_type())
}

// SourceMetaInfo returns the MetaInfo for GstRTPSourceMeta.
func SourceMetaInfo() *gst.MetaInfo {
	return gst.FromGstMetaInfoUnsafe(unsafe.Pointer(C.gst_rtp_source_meta_get_info()))
}

// AppendCSRC appends the given CSRC to the list of contributing sources in meta.
func (s *SourceMeta) AppendCSRC(csrc []uint32) bool {
	return gobool(C.gst_rtp_source_meta_append_csrc(
		s.ptr,
		(*C.guint32)(unsafe.Pointer(&csrc[0])),
		C.guint(len(csrc)),
	))
}

// GetSourceMeta retrieves the SourceMeta from the given buffer.
func GetSourceMeta(buffer *gst.Buffer) *SourceMeta {
	meta := C.gst_buffer_get_rtp_source_meta((*C.GstBuffer)(unsafe.Pointer(buffer.Instance())))
	if meta == nil {
		return nil
	}
	return &SourceMeta{ptr: meta}
}

// AddSourceMeta attaches the given RTP source information to the buffer.
func AddSourceMeta(buffer *gst.Buffer, ssrc *uint32, csrc []uint32) *SourceMeta {
	var meta *C.GstRTPSourceMeta
	if ssrc == nil && csrc == nil {
		meta = C.gst_buffer_add_rtp_source_meta(
			(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
			nil, nil, 0,
		)
	} else if ssrc == nil {
		meta = C.gst_buffer_add_rtp_source_meta(
			(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
			nil,
			(*C.guint32)(unsafe.Pointer(&csrc[0])), C.guint(len(csrc)),
		)
	} else if csrc == nil {
		meta = C.gst_buffer_add_rtp_source_meta(
			(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
			(*C.guint32)(unsafe.Pointer(ssrc)),
			nil, 0,
		)
	} else {
		meta = C.gst_buffer_add_rtp_source_meta(
			(*C.GstBuffer)(unsafe.Pointer(buffer.Instance())),
			(*C.guint32)(unsafe.Pointer(ssrc)),
			(*C.guint32)(unsafe.Pointer(&csrc[0])), C.guint(len(csrc)),
		)
	}
	if meta == nil {
		return nil
	}
	return &SourceMeta{ptr: meta}
}

// GetSourceCount returns the total number of RTP sources found in this meta, both SSRC and CSRC.
func (s *SourceMeta) GetSourceCount() uint {
	return uint(C.gst_rtp_source_meta_get_source_count(s.ptr))
}

// SetSSRC sets the SSRC on meta. If ssrc is nil, the SSRC of meta will be unset.
func (s *SourceMeta) SetSSRC(ssrc *uint32) bool {
	if ssrc == nil {
		return gobool(C.gst_rtp_source_meta_set_ssrc(s.ptr, nil))
	}
	return gobool(C.gst_rtp_source_meta_set_ssrc(s.ptr, (*C.guint32)(unsafe.Pointer(ssrc))))
}
