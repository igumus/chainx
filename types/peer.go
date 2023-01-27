package types

type PeerID string

type PeerState byte

const (
	PendingPeer    PeerState = 0x0
	HandshakedPeer PeerState = 0x1
)
