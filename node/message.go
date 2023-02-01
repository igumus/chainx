package node

import (
	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/network"
)

type ChainStateMessage struct {
	ID      network.PeerID
	Version uint32
	Height  uint32
}

type FetchBlockMessage struct {
	ID   network.PeerID
	From uint32
	To   uint32
}

type FetchBlockReply struct {
	Blocks []*core.Block
}
