package node

import (
	"bytes"
	"encoding/gob"

	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/network"
	"github.com/igumus/chainx/types"
)

type ChainStateMessage struct {
	ID      types.PeerID
	Version uint32
	Height  uint32
}

func NewChainStateMessage(id types.PeerID, version uint32, height uint32) (*network.Message, error) {
	msg := &ChainStateMessage{
		ID:      id,
		Version: version,
		Height:  height,
	}
	return msg.ToMessage()
}

func (m ChainStateMessage) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m ChainStateMessage) ToMessage() (*network.Message, error) {
	data, err := m.Bytes()
	if err != nil {
		return nil, err
	}
	return &network.Message{
		Header: network.ChainState,
		Data:   data,
	}, nil
}

type FetchBlockMessage struct {
	ID   types.PeerID
	From uint32
	To   uint32
}

func NewFetchBlockMessage(id types.PeerID, from uint32, to uint32) (*network.Message, error) {
	msg := &FetchBlockMessage{
		ID:   id,
		From: from,
		To:   to,
	}
	return msg.ToMessage()
}

func (m FetchBlockMessage) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m FetchBlockMessage) ToMessage() (*network.Message, error) {
	data, err := m.Bytes()
	if err != nil {
		return nil, err
	}
	return &network.Message{
		Header: network.ChainFetchBlock,
		Data:   data,
	}, nil
}

type FetchBlockReply struct {
	Blocks []*core.Block
}

func NewFetchBlockReply(blocks []*core.Block) (*network.Message, error) {
	msg := &FetchBlockReply{
		Blocks: blocks,
	}
	return msg.ToMessage()
}

func (m FetchBlockReply) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m FetchBlockReply) ToMessage() (*network.Message, error) {
	data, err := m.Bytes()
	if err != nil {
		return nil, err
	}
	return &network.Message{
		Header: network.ChainFetchBlockReply,
		Data:   data,
	}, nil
}
