package gst

// Flags returns the flags of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Flags() SegmentFlags {
	return SegmentFlags(segment.segment.native.flags)
}

// Rate returns the rate of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Rate() float64 {
	return float64(segment.segment.native.rate)
}

// AppliedRate returns the applied rate of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) AppliedRate() float64 {
	return float64(segment.segment.native.applied_rate)
}

// Format returns the format of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Format() Format {
	return Format(segment.segment.native.format)
}

// Base returns the base of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Base() uint64 {
	return uint64(segment.segment.native.base)
}

// Offset returns the offset of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Offset() uint64 {
	return uint64(segment.segment.native.offset)
}

// Start returns the start of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Start() uint64 {
	return uint64(segment.segment.native.start)
}

// Stop returns the stop of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Stop() uint64 {
	return uint64(segment.segment.native.stop)
}

// Time returns the time of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Time() uint64 {
	return uint64(segment.segment.native.time)
}

// Position returns the position of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Position() uint64 {
	return uint64(segment.segment.native.position)
}

// Duration returns the duration of the segment. See https://gstreamer.freedesktop.org/documentation/gstreamer/gstsegment.html?gi-language=c#members
func (segment *Segment) Duration() uint64 {
	return uint64(segment.segment.native.duration)
}
