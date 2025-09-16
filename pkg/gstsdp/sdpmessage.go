package gstsdp

// #cgo pkg-config: gstreamer-sdp-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/sdp/sdp.h>
import "C"
import (
	"iter"
)

// getters for the fields of GstSDPMessage:
// All simple fields already have getters generated, we need to create iterators for the GArrays.
// Each GArray has a getter to access the underlying length and an index function to access the elements.

// version (gchar *) - the protocol version
// origin (GstSDPOrigin) - owner/creator and session identifier
// session_name (gchar *) - session name
// information (gchar *) - session information
// uri (gchar *) - URI of description
// connection (GstSDPConnection) - connection information for the session
// key (GstSDPKey) - encryption key
//
// emails (GArray *) - array of gchar with email addresses
// phones (GArray *) - array of gchar with phone numbers
// bandwidths (GArray *) - array of GstSDPBandwidth with bandwidth information
// times (GArray *) - array of GstSDPTime with time descriptions
// zones (GArray *) - array of GstSDPZone with time zone adjustments
// attributes (GArray *) - array of GstSDPAttribute with session attributes
// medias (GArray *) - array of GstSDPMedia with media descriptions

// Emails returns an iterator over the email addresses in the SDPMessage
func (s *SDPMessage) Emails() iter.Seq2[uint, string] {
	return getIter(s.EmailsLen(), s.GetEmail)
}

// Phones returns an iterator over the phone numbers in the SDPMessage
func (s *SDPMessage) Phones() iter.Seq2[uint, string] {
	return getIter(s.PhonesLen(), s.GetPhone)
}

// Bandwidths returns an iterator over the bandwidths in the SDPMessage
func (s *SDPMessage) Bandwidths() iter.Seq2[uint, *SDPBandwidth] {
	return getIter(s.BandwidthsLen(), s.GetBandwidth)
}

// Times returns an iterator over the times in the SDPMessage
func (s *SDPMessage) Times() iter.Seq2[uint, *SDPTime] {
	return getIter(s.TimesLen(), s.GetTime)
}

// Zones returns an iterator over the zones in the SDPMessage
func (s *SDPMessage) Zones() iter.Seq2[uint, *SDPZone] {
	return getIter(s.ZonesLen(), s.GetZone)
}

// Attributes returns an iterator over the attributes in the SDPMessage
func (s *SDPMessage) Attributes() iter.Seq2[uint, *SDPAttribute] {
	return getIter(s.AttributesLen(), s.GetAttribute)
}

// Medias returns an iterator over the medias in the SDPMessage
func (s *SDPMessage) Medias() iter.Seq2[uint, *SDPMedia] {
	return getIter(s.MediasLen(), s.GetMedia)
}
