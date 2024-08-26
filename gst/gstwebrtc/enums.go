package gstwebrtc

// #include "gst.go.h"
import "C"

type BundlePolicy C.GstWebRTCBundlePolicy

const (
	BUNDLE_POLICY_NONE       BundlePolicy = C.GST_WEBRTC_BUNDLE_POLICY_NONE       // none
	BUNDLE_POLICY_BALANCED   BundlePolicy = C.GST_WEBRTC_BUNDLE_POLICY_BALANCED   // balanced
	BUNDLE_POLICY_MAX_COMPAT BundlePolicy = C.GST_WEBRTC_BUNDLE_POLICY_MAX_COMPAT // max-compat
	BUNDLE_POLICY_MAX_BUNDLE BundlePolicy = C.GST_WEBRTC_BUNDLE_POLICY_MAX_BUNDLE // max-bundle
)

func (e BundlePolicy) String() string {
	switch e {
	case C.GST_WEBRTC_BUNDLE_POLICY_NONE:
		return "none"
	case C.GST_WEBRTC_BUNDLE_POLICY_BALANCED:
		return "balanced"
	case C.GST_WEBRTC_BUNDLE_POLICY_MAX_COMPAT:
		return "max-compat"
	case C.GST_WEBRTC_BUNDLE_POLICY_MAX_BUNDLE:
		return "max-bundle"
	}
	return "unknown"
}

type DTLSSetup C.GstWebRTCDTLSSetup

const (
	DTLS_SETUP_NONE    DTLSSetup = C.GST_WEBRTC_DTLS_SETUP_NONE    // none
	DTLS_SETUP_ACTPASS DTLSSetup = C.GST_WEBRTC_DTLS_SETUP_ACTPASS // actpass
	DTLS_SETUP_ACTIVE  DTLSSetup = C.GST_WEBRTC_DTLS_SETUP_ACTIVE  // sendonly
	DTLS_SETUP_PASSIVE DTLSSetup = C.GST_WEBRTC_DTLS_SETUP_PASSIVE // recvonly
)

func (e DTLSSetup) String() string {
	switch e {
	case C.GST_WEBRTC_DTLS_SETUP_NONE:
		return "none"
	case C.GST_WEBRTC_DTLS_SETUP_ACTPASS:
		return "actpass"
	case C.GST_WEBRTC_DTLS_SETUP_ACTIVE:
		return "sendonly"
	case C.GST_WEBRTC_DTLS_SETUP_PASSIVE:
		return "recvonly"
	}
	return "unknown"
}

type DTLSTransportState C.GstWebRTCDTLSTransportState

const (
	DTLS_TRANSPORT_STATE_NEW        DTLSTransportState = C.GST_WEBRTC_DTLS_TRANSPORT_STATE_NEW        // new
	DTLS_TRANSPORT_STATE_CLOSED     DTLSTransportState = C.GST_WEBRTC_DTLS_TRANSPORT_STATE_CLOSED     // closed
	DTLS_TRANSPORT_STATE_FAILED     DTLSTransportState = C.GST_WEBRTC_DTLS_TRANSPORT_STATE_FAILED     // failed
	DTLS_TRANSPORT_STATE_CONNECTING DTLSTransportState = C.GST_WEBRTC_DTLS_TRANSPORT_STATE_CONNECTING // connecting
	DTLS_TRANSPORT_STATE_CONNECTED  DTLSTransportState = C.GST_WEBRTC_DTLS_TRANSPORT_STATE_CONNECTED  // connected
)

func (e DTLSTransportState) String() string {
	switch e {
	case C.GST_WEBRTC_DTLS_TRANSPORT_STATE_NEW:
		return "new"
	case C.GST_WEBRTC_DTLS_TRANSPORT_STATE_CLOSED:
		return "closed"
	case C.GST_WEBRTC_DTLS_TRANSPORT_STATE_FAILED:
		return "failed"
	case C.GST_WEBRTC_DTLS_TRANSPORT_STATE_CONNECTING:
		return "connecting"
	case C.GST_WEBRTC_DTLS_TRANSPORT_STATE_CONNECTED:
		return "connected"
	}
	return "unknown"
}

type DataChannelState C.GstWebRTCDataChannelState

const (
	DATA_CHANNEL_STATE_CONNECTING DataChannelState = C.GST_WEBRTC_DATA_CHANNEL_STATE_CONNECTING // connecting
	DATA_CHANNEL_STATE_OPEN       DataChannelState = C.GST_WEBRTC_DATA_CHANNEL_STATE_OPEN       // open
	DATA_CHANNEL_STATE_CLOSING    DataChannelState = C.GST_WEBRTC_DATA_CHANNEL_STATE_CLOSING    // closing
	DATA_CHANNEL_STATE_CLOSED     DataChannelState = C.GST_WEBRTC_DATA_CHANNEL_STATE_CLOSED     // closed
)

func (e DataChannelState) String() string {
	switch e {
	case C.GST_WEBRTC_DATA_CHANNEL_STATE_CONNECTING:
		return "connecting"
	case C.GST_WEBRTC_DATA_CHANNEL_STATE_OPEN:
		return "open"
	case C.GST_WEBRTC_DATA_CHANNEL_STATE_CLOSING:
		return "closing"
	case C.GST_WEBRTC_DATA_CHANNEL_STATE_CLOSED:
		return "closed"
	}
	return "unknown"
}

type Error C.GstWebRTCError

const (
	ERROR_DATA_CHANNEL_FAILURE           Error = C.GST_WEBRTC_ERROR_DATA_CHANNEL_FAILURE           // data-channel-failure
	ERROR_DTLS_FAILURE                   Error = C.GST_WEBRTC_ERROR_DTLS_FAILURE                   // dtls-failure
	ERROR_FINGERPRINT_FAILURE            Error = C.GST_WEBRTC_ERROR_FINGERPRINT_FAILURE            // fingerprint-failure
	ERROR_SCTP_FAILURE                   Error = C.GST_WEBRTC_ERROR_SCTP_FAILURE                   // sctp-failure
	ERROR_SDP_SYNTAX_ERROR               Error = C.GST_WEBRTC_ERROR_SDP_SYNTAX_ERROR               // sdp-syntax-error
	ERROR_HARDWARE_ENCODER_NOT_AVAILABLE Error = C.GST_WEBRTC_ERROR_HARDWARE_ENCODER_NOT_AVAILABLE // hardware-encoder-not-available
	ERROR_ENCODER_ERROR                  Error = C.GST_WEBRTC_ERROR_ENCODER_ERROR                  // encoder-error
	ERROR_INVALID_STATE                  Error = C.GST_WEBRTC_ERROR_INVALID_STATE                  // invalid-state
	ERROR_INTERNAL_FAILURE               Error = C.GST_WEBRTC_ERROR_INTERNAL_FAILURE               // internal-failure
	ERROR_INVALID_MODIFICATION           Error = C.GST_WEBRTC_ERROR_INVALID_MODIFICATION           // invalid-modification
	ERROR_TYPE_ERROR                     Error = C.GST_WEBRTC_ERROR_TYPE_ERROR                     // type-error
)

func (e Error) String() string {
	switch e {
	case C.GST_WEBRTC_ERROR_DATA_CHANNEL_FAILURE:
		return "data-channel-failure"
	case C.GST_WEBRTC_ERROR_DTLS_FAILURE:
		return "dtls-failure"
	case C.GST_WEBRTC_ERROR_FINGERPRINT_FAILURE:
		return "fingerprint-failure"
	case C.GST_WEBRTC_ERROR_SCTP_FAILURE:
		return "sctp-failure"
	case C.GST_WEBRTC_ERROR_SDP_SYNTAX_ERROR:
		return "sdp-syntax-error"
	case C.GST_WEBRTC_ERROR_HARDWARE_ENCODER_NOT_AVAILABLE:
		return "hardware-encoder-not-available"
	case C.GST_WEBRTC_ERROR_ENCODER_ERROR:
		return "encoder-error"
	case C.GST_WEBRTC_ERROR_INVALID_STATE:
		return "invalid-state"
	case C.GST_WEBRTC_ERROR_INTERNAL_FAILURE:
		return "internal-failure"
	case C.GST_WEBRTC_ERROR_INVALID_MODIFICATION:
		return "invalid-modification"
	case C.GST_WEBRTC_ERROR_TYPE_ERROR:
		return "type-error"
	}
	return "unknown"
}

type FECType C.GstWebRTCFECType

const (
	FEC_TYPE_NONE    FECType = C.GST_WEBRTC_FEC_TYPE_NONE    // none
	FEC_TYPE_ULP_RED FECType = C.GST_WEBRTC_FEC_TYPE_ULP_RED // ulpfec + red
)

func (e FECType) String() string {
	switch e {
	case C.GST_WEBRTC_FEC_TYPE_NONE:
		return "none"
	case C.GST_WEBRTC_FEC_TYPE_ULP_RED:
		return "ulpfec + red"
	}
	return "unknown"
}

type ICEComponent C.GstWebRTCICEComponent

//  GST_WEBRTC_ICE_COMPONENT_RTP (0)RTP component
// GST_WEBRTC_ICE_COMPONENT_RTCP (1)RTCP component

const (
	ICE_COMPONENT_RTP  ICEComponent = C.GST_WEBRTC_ICE_COMPONENT_RTP  // RTP component
	ICE_COMPONENT_RTCP ICEComponent = C.GST_WEBRTC_ICE_COMPONENT_RTCP // RTCP component
)

func (e ICEComponent) String() string {
	switch e {
	case C.GST_WEBRTC_ICE_COMPONENT_RTP:
		return "RTP component"
	case C.GST_WEBRTC_ICE_COMPONENT_RTCP:
		return "RTCP component"
	}
	return "unknown"
}

type ICEConnectionState C.GstWebRTCICEConnectionState

const (
	ICE_CONNECTION_STATE_NEW          ICEConnectionState = C.GST_WEBRTC_ICE_CONNECTION_STATE_NEW          // new
	ICE_CONNECTION_STATE_CHECKING     ICEConnectionState = C.GST_WEBRTC_ICE_CONNECTION_STATE_CHECKING     // checking
	ICE_CONNECTION_STATE_CONNECTED    ICEConnectionState = C.GST_WEBRTC_ICE_CONNECTION_STATE_CONNECTED    // connected
	ICE_CONNECTION_STATE_COMPLETED    ICEConnectionState = C.GST_WEBRTC_ICE_CONNECTION_STATE_COMPLETED    // completed
	ICE_CONNECTION_STATE_FAILED       ICEConnectionState = C.GST_WEBRTC_ICE_CONNECTION_STATE_FAILED       // failed
	ICE_CONNECTION_STATE_DISCONNECTED ICEConnectionState = C.GST_WEBRTC_ICE_CONNECTION_STATE_DISCONNECTED // disconnected
	ICE_CONNECTION_STATE_CLOSED       ICEConnectionState = C.GST_WEBRTC_ICE_CONNECTION_STATE_CLOSED       // closed
)

func (e ICEConnectionState) String() string {
	switch e {
	case C.GST_WEBRTC_ICE_CONNECTION_STATE_NEW:
		return "new"
	case C.GST_WEBRTC_ICE_CONNECTION_STATE_CHECKING:
		return "checking"
	case C.GST_WEBRTC_ICE_CONNECTION_STATE_CONNECTED:
		return "connected"
	case C.GST_WEBRTC_ICE_CONNECTION_STATE_COMPLETED:
		return "completed"
	case C.GST_WEBRTC_ICE_CONNECTION_STATE_FAILED:
		return "failed"
	case C.GST_WEBRTC_ICE_CONNECTION_STATE_DISCONNECTED:
		return "disconnected"
	case C.GST_WEBRTC_ICE_CONNECTION_STATE_CLOSED:
		return "closed"
	}
	return "unknown"
}

type ICEGatheringState C.GstWebRTCICEGatheringState

const (
	ICE_GATHERING_STATE_NEW       ICEGatheringState = C.GST_WEBRTC_ICE_GATHERING_STATE_NEW       // new
	ICE_GATHERING_STATE_GATHERING ICEGatheringState = C.GST_WEBRTC_ICE_GATHERING_STATE_GATHERING // gathering
	ICE_GATHERING_STATE_COMPLETE  ICEGatheringState = C.GST_WEBRTC_ICE_GATHERING_STATE_COMPLETE  // complete
)

func (e ICEGatheringState) String() string {
	switch e {
	case C.GST_WEBRTC_ICE_GATHERING_STATE_NEW:
		return "new"
	case C.GST_WEBRTC_ICE_GATHERING_STATE_GATHERING:
		return "gathering"
	case C.GST_WEBRTC_ICE_GATHERING_STATE_COMPLETE:
		return "complete"
	}
	return "unknown"
}

type ICERole C.GstWebRTCICERole

const (
	ICE_ROLE_CONTROLLED  ICERole = C.GST_WEBRTC_ICE_ROLE_CONTROLLED  // controlled
	ICE_ROLE_CONTROLLING ICERole = C.GST_WEBRTC_ICE_ROLE_CONTROLLING // controlling
)

func (e ICERole) String() string {
	switch e {
	case C.GST_WEBRTC_ICE_ROLE_CONTROLLED:
		return "controlled"
	case C.GST_WEBRTC_ICE_ROLE_CONTROLLING:
		return "controlling"
	}
	return "unknown"
}

type ICETransportPolicy C.GstWebRTCICETransportPolicy

const (
	ICE_TRANSPORT_POLICY_ALL   ICETransportPolicy = C.GST_WEBRTC_ICE_TRANSPORT_POLICY_ALL   // all
	ICE_TRANSPORT_POLICY_RELAY ICETransportPolicy = C.GST_WEBRTC_ICE_TRANSPORT_POLICY_RELAY // relay
)

func (e ICETransportPolicy) String() string {
	switch e {
	case C.GST_WEBRTC_ICE_TRANSPORT_POLICY_ALL:
		return "all"
	case C.GST_WEBRTC_ICE_TRANSPORT_POLICY_RELAY:
		return "relay"
	}
	return "unknown"
}

type Kind C.GstWebRTCKind

const (
	UNKNOWN Kind = C.GST_WEBRTC_KIND_UNKNOWN // unknown
	AUDIO   Kind = C.GST_WEBRTC_KIND_AUDIO   // audio
	VIDEO   Kind = C.GST_WEBRTC_KIND_VIDEO   // video
)

func (e Kind) String() string {
	switch e {
	case C.GST_WEBRTC_KIND_UNKNOWN:
		return "unknown"
	case C.GST_WEBRTC_KIND_AUDIO:
		return "audio"
	case C.GST_WEBRTC_KIND_VIDEO:
		return "video"
	}
	return "unknown"
}

type PeerConnectionState C.GstWebRTCPeerConnectionState

//  GST_WEBRTC_PEER_CONNECTION_STATE_NEW (0)new
// GST_WEBRTC_PEER_CONNECTION_STATE_CONNECTING (1)connecting
// GST_WEBRTC_PEER_CONNECTION_STATE_CONNECTED (2)connected
// GST_WEBRTC_PEER_CONNECTION_STATE_DISCONNECTED (3)disconnected
// GST_WEBRTC_PEER_CONNECTION_STATE_FAILED (4)failed
// GST_WEBRTC_PEER_CONNECTION_STATE_CLOSED (5)closed

const (
	PEER_CONNECTION_STATE_NEW          PeerConnectionState = C.GST_WEBRTC_PEER_CONNECTION_STATE_NEW          // new
	PEER_CONNECTION_STATE_CONNECTING   PeerConnectionState = C.GST_WEBRTC_PEER_CONNECTION_STATE_CONNECTING   // connecting
	PEER_CONNECTION_STATE_CONNECTED    PeerConnectionState = C.GST_WEBRTC_PEER_CONNECTION_STATE_CONNECTED    // connected
	PEER_CONNECTION_STATE_DISCONNECTED PeerConnectionState = C.GST_WEBRTC_PEER_CONNECTION_STATE_DISCONNECTED // disconnected
	PEER_CONNECTION_STATE_FAILED       PeerConnectionState = C.GST_WEBRTC_PEER_CONNECTION_STATE_FAILED       // failed
	PEER_CONNECTION_STATE_CLOSED       PeerConnectionState = C.GST_WEBRTC_PEER_CONNECTION_STATE_CLOSED       // closed
)

func (e PeerConnectionState) String() string {
	switch e {
	case C.GST_WEBRTC_PEER_CONNECTION_STATE_NEW:
		return "new"
	case C.GST_WEBRTC_PEER_CONNECTION_STATE_CONNECTING:
		return "connecting"
	case C.GST_WEBRTC_PEER_CONNECTION_STATE_CONNECTED:
		return "connected"
	case C.GST_WEBRTC_PEER_CONNECTION_STATE_DISCONNECTED:
		return "disconnected"
	case C.GST_WEBRTC_PEER_CONNECTION_STATE_FAILED:
		return "failed"
	case C.GST_WEBRTC_PEER_CONNECTION_STATE_CLOSED:
		return "closed"
	}
	return "unknown"
}

type PriorityType C.GstWebRTCPriorityType

const (
	PRIORITY_TYPE_VERY_LOW PriorityType = C.GST_WEBRTC_PRIORITY_TYPE_VERY_LOW // very-low
	PRIORITY_TYPE_LOW      PriorityType = C.GST_WEBRTC_PRIORITY_TYPE_LOW      // low
	PRIORITY_TYPE_MEDIUM   PriorityType = C.GST_WEBRTC_PRIORITY_TYPE_MEDIUM   // medium
	PRIORITY_TYPE_HIGH     PriorityType = C.GST_WEBRTC_PRIORITY_TYPE_HIGH     // high
)

func (e PriorityType) String() string {
	switch e {
	case C.GST_WEBRTC_PRIORITY_TYPE_VERY_LOW:
		return "very-low"
	case C.GST_WEBRTC_PRIORITY_TYPE_LOW:
		return "low"
	case C.GST_WEBRTC_PRIORITY_TYPE_MEDIUM:
		return "medium"
	case C.GST_WEBRTC_PRIORITY_TYPE_HIGH:
		return "high"
	}
	return "unknown"
}

type RTPTransceiverDirection C.GstWebRTCRTPTransceiverDirection

const (
	RTP_TRANSCEIVER_DIRECTION_NONE     RTPTransceiverDirection = C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_NONE     // none
	RTP_TRANSCEIVER_DIRECTION_INACTIVE RTPTransceiverDirection = C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_INACTIVE // inactive
	RTP_TRANSCEIVER_DIRECTION_SENDONLY RTPTransceiverDirection = C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_SENDONLY // sendonly
	RTP_TRANSCEIVER_DIRECTION_RECVONLY RTPTransceiverDirection = C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_RECVONLY // recvonly
	RTP_TRANSCEIVER_DIRECTION_SENDRECV RTPTransceiverDirection = C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_SENDRECV // sendrecv
)

func (e RTPTransceiverDirection) String() string {
	switch e {
	case C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_NONE:
		return "none"
	case C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_INACTIVE:
		return "inactive"
	case C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_SENDONLY:
		return "sendonly"
	case C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_RECVONLY:
		return "recvonly"
	case C.GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_SENDRECV:
		return "sendrecv"
	}
	return "unknown"
}

type SCTPTransportState C.GstWebRTCSCTPTransportState

const (
	SCTP_TRANSPORT_STATE_NEW        SCTPTransportState = C.GST_WEBRTC_SCTP_TRANSPORT_STATE_NEW        // new
	SCTP_TRANSPORT_STATE_CONNECTING SCTPTransportState = C.GST_WEBRTC_SCTP_TRANSPORT_STATE_CONNECTING // connecting
	SCTP_TRANSPORT_STATE_CONNECTED  SCTPTransportState = C.GST_WEBRTC_SCTP_TRANSPORT_STATE_CONNECTED  // connected
	SCTP_TRANSPORT_STATE_CLOSED     SCTPTransportState = C.GST_WEBRTC_SCTP_TRANSPORT_STATE_CLOSED     // closed
)

func (e SCTPTransportState) String() string {
	switch e {
	case C.GST_WEBRTC_SCTP_TRANSPORT_STATE_NEW:
		return "new"
	case C.GST_WEBRTC_SCTP_TRANSPORT_STATE_CONNECTING:
		return "connecting"
	case C.GST_WEBRTC_SCTP_TRANSPORT_STATE_CONNECTED:
		return "connected"
	case C.GST_WEBRTC_SCTP_TRANSPORT_STATE_CLOSED:
		return "closed"
	}
	return "unknown"
}

type SDPType C.GstWebRTCSDPType

const (
	SDP_TYPE_OFFER    SDPType = C.GST_WEBRTC_SDP_TYPE_OFFER    // offer
	SDP_TYPE_PRANSWER SDPType = C.GST_WEBRTC_SDP_TYPE_PRANSWER // pranswer
	SDP_TYPE_ANSWER   SDPType = C.GST_WEBRTC_SDP_TYPE_ANSWER   // answer
	SDP_TYPE_ROLLBACK SDPType = C.GST_WEBRTC_SDP_TYPE_ROLLBACK // rollback
)

func (e SDPType) String() string {
	// returned string is const gchar* and must not be freed
	cstring := C.gst_webrtc_sdp_type_to_string(C.GstWebRTCSDPType(e))

	return C.GoString(cstring)
}

func SDPTypeFromString(s string) SDPType {
	switch s {
	case "offer":
		return SDP_TYPE_OFFER
	case "pranswer":
		return SDP_TYPE_PRANSWER
	case "answer":
		return SDP_TYPE_ANSWER
	case "rollback":
		return SDP_TYPE_ROLLBACK
	default:
		panic("Unknown SDPType")
	}
}

type SignalingState C.GstWebRTCSignalingState

const (
	SIGNALING_STATE_STABLE               SignalingState = C.GST_WEBRTC_SIGNALING_STATE_STABLE               // stable
	SIGNALING_STATE_CLOSED               SignalingState = C.GST_WEBRTC_SIGNALING_STATE_CLOSED               // closed
	SIGNALING_STATE_HAVE_LOCAL_OFFER     SignalingState = C.GST_WEBRTC_SIGNALING_STATE_HAVE_LOCAL_OFFER     // have-local-offer
	SIGNALING_STATE_HAVE_REMOTE_OFFER    SignalingState = C.GST_WEBRTC_SIGNALING_STATE_HAVE_REMOTE_OFFER    // have-remote-offer
	SIGNALING_STATE_HAVE_LOCAL_PRANSWER  SignalingState = C.GST_WEBRTC_SIGNALING_STATE_HAVE_LOCAL_PRANSWER  // have-local-pranswer
	SIGNALING_STATE_HAVE_REMOTE_PRANSWER SignalingState = C.GST_WEBRTC_SIGNALING_STATE_HAVE_REMOTE_PRANSWER // have-remote-pranswer
)

func (e SignalingState) String() string {
	switch e {
	case C.GST_WEBRTC_SIGNALING_STATE_STABLE:
		return "stable"
	case C.GST_WEBRTC_SIGNALING_STATE_CLOSED:
		return "closed"
	case C.GST_WEBRTC_SIGNALING_STATE_HAVE_LOCAL_OFFER:
		return "have-local-offer"
	case C.GST_WEBRTC_SIGNALING_STATE_HAVE_REMOTE_OFFER:
		return "have-remote-offer"
	case C.GST_WEBRTC_SIGNALING_STATE_HAVE_LOCAL_PRANSWER:
		return "have-local-pranswer"
	case C.GST_WEBRTC_SIGNALING_STATE_HAVE_REMOTE_PRANSWER:
		return "have-remote-pranswer"
	}
	return "unknown"
}

type StatsType C.GstWebRTCStatsType

const (
	STATS_CODEC               StatsType = C.GST_WEBRTC_STATS_CODEC               // codec
	STATS_INBOUND_RTP         StatsType = C.GST_WEBRTC_STATS_INBOUND_RTP         // inbound-rtp
	STATS_OUTBOUND_RTP        StatsType = C.GST_WEBRTC_STATS_OUTBOUND_RTP        // outbound-rtp
	STATS_REMOTE_INBOUND_RTP  StatsType = C.GST_WEBRTC_STATS_REMOTE_INBOUND_RTP  // remote-inbound-rtp
	STATS_REMOTE_OUTBOUND_RTP StatsType = C.GST_WEBRTC_STATS_REMOTE_OUTBOUND_RTP // remote-outbound-rtp
	STATS_CSRC                StatsType = C.GST_WEBRTC_STATS_CSRC                // csrc
	STATS_PEER_CONNECTION     StatsType = C.GST_WEBRTC_STATS_PEER_CONNECTION     // peer-connection
	STATS_DATA_CHANNEL        StatsType = C.GST_WEBRTC_STATS_DATA_CHANNEL        // data-channel
	STATS_STREAM              StatsType = C.GST_WEBRTC_STATS_STREAM              // stream
	STATS_TRANSPORT           StatsType = C.GST_WEBRTC_STATS_TRANSPORT           // transport
	STATS_CANDIDATE_PAIR      StatsType = C.GST_WEBRTC_STATS_CANDIDATE_PAIR      // candidate-pair
	STATS_LOCAL_CANDIDATE     StatsType = C.GST_WEBRTC_STATS_LOCAL_CANDIDATE     // local-candidate
	STATS_REMOTE_CANDIDATE    StatsType = C.GST_WEBRTC_STATS_REMOTE_CANDIDATE    // remote-candidate
	STATS_CERTIFICATE         StatsType = C.GST_WEBRTC_STATS_CERTIFICATE         // certificate
)

func (e StatsType) String() string {
	switch e {
	case C.GST_WEBRTC_STATS_CODEC:
		return "codec"
	case C.GST_WEBRTC_STATS_INBOUND_RTP:
		return "inbound-rtp"
	case C.GST_WEBRTC_STATS_OUTBOUND_RTP:
		return "outbound-rtp"
	case C.GST_WEBRTC_STATS_REMOTE_INBOUND_RTP:
		return "remote-inbound-rtp"
	case C.GST_WEBRTC_STATS_REMOTE_OUTBOUND_RTP:
		return "remote-outbound-rtp"
	case C.GST_WEBRTC_STATS_CSRC:
		return "csrc"
	case C.GST_WEBRTC_STATS_PEER_CONNECTION:
		return "peer-connection"
	case C.GST_WEBRTC_STATS_DATA_CHANNEL:
		return "data-channel"
	case C.GST_WEBRTC_STATS_STREAM:
		return "stream"
	case C.GST_WEBRTC_STATS_TRANSPORT:
		return "transport"
	case C.GST_WEBRTC_STATS_CANDIDATE_PAIR:
		return "candidate-pair"
	case C.GST_WEBRTC_STATS_LOCAL_CANDIDATE:
		return "local-candidate"
	case C.GST_WEBRTC_STATS_REMOTE_CANDIDATE:
		return "remote-candidate"
	case C.GST_WEBRTC_STATS_CERTIFICATE:
		return "certificate"
	}
	return "unknown"
}
