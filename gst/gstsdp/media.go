package gstsdp

// #include "gst.go.h"
import "C"
import (
	"errors"
	"iter"
	"unsafe"

	"github.com/go-gst/go-gst/gst"
)

type Media struct {
	ptr *C.GstSDPMedia
}

func (m *Media) FormatsLen() int {
	return int(C.gst_sdp_media_formats_len(m.ptr))
}

func (m *Media) Format(idx int) string {
	cstr := C.gst_sdp_media_get_format(m.ptr, C.guint(idx))

	return C.GoString(cstr)
}

func (m *Media) Formats() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i := 0; i < m.FormatsLen(); i++ {
			if !yield(i, m.Format(i)) {
				return
			}
		}
	}
}

var ErrCouldNotGetCaps = errors.New("could not get caps")

func (m *Media) GetCaps(pt int) (*gst.Caps, error) {
	ccaps := C.gst_sdp_media_get_caps_from_media(m.ptr, C.gint(pt))

	if ccaps == nil {
		return nil, ErrCouldNotGetCaps
	}

	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(ccaps)), nil
}
