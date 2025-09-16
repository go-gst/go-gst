package gst

func (r *Event) GetType() EventType {
	return EventType(r.event.native._type)
}

func (r *Event) GetTimestamp() uint64 {
	return uint64(r.event.native.timestamp)
}

func (r *Event) GetSeqNum() uint32 {
	return uint32(r.event.native.seqnum)
}
