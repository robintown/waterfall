package peer

import (
	"github.com/pion/webrtc/v3"
	"maunium.net/go/mautrix/event"
)

// Due to the limitation of Go, we're using the `interface{}` to be able to use switch the actual
// type of the message on runtime. The underlying types do not necessary need to be structures.
type MessageContent = interface{}

type JoinedTheCall struct{}

type LeftTheCall struct {
	Reason event.CallHangupReason
}

type NewTrackPublished struct {
	Track *webrtc.TrackLocalStaticRTP
}

type PublishedTrackFailed struct {
	Track *webrtc.TrackLocalStaticRTP
}

type NewICECandidate struct {
	Candidate *webrtc.ICECandidate
}

type ICEGatheringComplete struct{}

type RenegotiationRequired struct {
	Offer *webrtc.SessionDescription
}

type DataChannelMessage struct {
	Message string
}

type DataChannelAvailable struct{}

type RTCPReceived struct {
	TrackID string
	Packets []RTCPPacketType
}

type RTCPPacketType int

const (
	PictureLossIndicator RTCPPacketType = iota + 1
	FullIntraRequest
)
