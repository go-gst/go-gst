package gstsdp

import "iter"

// Getters for SDPMedia properties.
// All simple fields already have getters generated, we need to create iterators for the GArrays.
// Each GArray has a getter to access the underlying length and an index function to access the elements.

// media (gchar *) - the media type
// port (guint) - the transport port to which the media stream will be sent
// num_ports (guint) - the number of ports or -1 if only one port was specified
// proto (gchar *) - the transport protocol
// information (gchar *) - the media title
// key (GstSDPKey) - the encryption key
//
// fmts (GArray *) - an array of gchar formats
// connections (GArray *) - array of GstSDPConnection with media connection information
// bandwidths (GArray *) - array of GstSDPBandwidth with media bandwidth information
// attributes (GArray *) - array of GstSDPAttribute with the additional media attributes

// Formats returns an iterator over the formats in the SDPMedia
func (m *SDPMedia) Formats() iter.Seq2[uint, string] {
	return getIter(m.FormatsLen(), m.GetFormat)
}

// Connections returns an iterator over the connections in the SDPMedia
func (m *SDPMedia) Connections() iter.Seq2[uint, *SDPConnection] {
	return getIter(m.ConnectionsLen(), m.GetConnection)
}

// Bandwidths returns an iterator over the bandwidths in the SDPMedia
func (m *SDPMedia) Bandwidths() iter.Seq2[uint, *SDPBandwidth] {
	return getIter(m.BandwidthsLen(), m.GetBandwidth)
}

// Attributes returns an iterator over the attributes in the SDPMedia
func (m *SDPMedia) Attributes() iter.Seq2[uint, *SDPAttribute] {
	return getIter(m.AttributesLen(), m.GetAttribute)
}
