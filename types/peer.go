package types

type PeerID string

func (pid PeerID) String() string {
	return string(pid)
}

type PeerState byte

const (
	PendingPeer    PeerState = 0x0
	HandshakedPeer PeerState = 0x1
)
