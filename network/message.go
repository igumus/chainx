package network

import (
	"bytes"
	"encoding/gob"
)

type RemoteMessageHandler interface {
	HandleMessage(RemoteMessage) error
}

type RemoteMessage struct {
	From    PeerID
	Payload []byte
}

func (m RemoteMessage) Decode() (*Message, error) {
	msg := &Message{}
	if err := gob.NewDecoder(bytes.NewReader(m.Payload)).Decode(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type MessageType byte

const (
	// network message types
	NetworkHandshake      MessageType = 0x1
	NetworkHandshakeReply MessageType = 0x2
	NetworkReserved_2     MessageType = 0x3
	NetworkReserved_3     MessageType = 0x4
	NetworkReserved_4     MessageType = 0x5
	NetworkReserved_5     MessageType = 0x6

	// chain message types
	ChainState           MessageType = 0x7
	ChainTx              MessageType = 0x8
	ChainBlock           MessageType = 0x9
	ChainFetchBlock      MessageType = 0xa
	ChainFetchBlockReply MessageType = 0xb
)

type Message struct {
	Header MessageType
	Data   []byte
}

func (m *Message) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
