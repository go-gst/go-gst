package gst

// #include "gst.go.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// StreamCollection is a Go representation of a GstStreamCollection.
type StreamCollection struct{ *Object }

// FromGstStreamCollectionUnsafeNone captures a pointer with a ref and finalizer.
func FromGstStreamCollectionUnsafeNone(stream unsafe.Pointer) *StreamCollection {
	return &StreamCollection{wrapObject(glib.TransferNone(stream))}
}

// FromGstStreamCollectionUnsafeFull captures a pointer with just a finalizer.
func FromGstStreamCollectionUnsafeFull(stream unsafe.Pointer) *StreamCollection {
	return &StreamCollection{wrapObject(glib.TransferFull(stream))}
}

// NewStreamCollection returns a new StreamCollection with an upstream parent
// of the given stream ID.
func NewStreamCollection(upstreamID string) *StreamCollection {
	cID := C.CString(upstreamID)
	defer C.free(unsafe.Pointer(cID))
	collection := C.gst_stream_collection_new(cID)
	return FromGstStreamCollectionUnsafeFull(unsafe.Pointer(collection))
}

// Instance returns the underlying GstStreamCollection.
func (s *StreamCollection) Instance() *C.GstStreamCollection {
	return C.toGstStreamCollection(s.Unsafe())
}

// AddStream adds the given stream to this collection.
func (s *StreamCollection) AddStream(stream *Stream) error {
	if ok := gobool(C.gst_stream_collection_add_stream(s.Instance(), stream.Instance())); !ok {
		return fmt.Errorf("Failed to add stream %s to collection", stream.StreamID())
	}
	return nil
}

// GetSize returns the size of this stream collection.
func (s *StreamCollection) GetSize() uint {
	return uint(C.gst_stream_collection_get_size(s.Instance()))
}

// GetStreamAt returns the stream at the given index in this collection.
func (s *StreamCollection) GetStreamAt(idx uint) *Stream {
	stream := C.gst_stream_collection_get_stream(s.Instance(), C.guint(idx))
	return FromGstStreamUnsafeNone(unsafe.Pointer(stream))
}

// GetUpstreamID retrieves the upstream ID for this collection.
func (s *StreamCollection) GetUpstreamID() string {
	return C.GoString(C.gst_stream_collection_get_upstream_id(s.Instance()))
}
