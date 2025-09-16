package gstwebrtc

const (
	// WebrtcSdpTypeUnknown indicates an unknown SDP type. It is the catch all for
	// any SDP type that is not recognized.
	WebrtcSdpTypeUnknown WebRTCSDPType = 0
)

func WebRTCSDPTypeFromString(typ string) WebRTCSDPType {
	switch typ {
	case "offer":
		return WebrtcSdpTypeOffer
	case "answer":
		return WebrtcSdpTypeAnswer
	case "pranswer":
		return WebrtcSdpTypePranswer
	case "rollback":
		return WebrtcSdpTypeRollback
	default:
		return WebrtcSdpTypeUnknown
	}
}
