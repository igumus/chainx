package types

import (
	"bytes"
	"encoding/gob"
)

type MessageType byte

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

const (
	// network message types
	NetworkHandshake      MessageType = 0x1
	NetworkHandshakeReply MessageType = 0x2
	NetworkReserved_2     MessageType = 0x3
	NetworkReserved_3     MessageType = 0x4
	NetworkReserved_4     MessageType = 0x5
	NetworkReserved_5     MessageType = 0x6

	// chain message types

)

type Message struct {
	Header MessageType
	Data   []byte
}

func (m *Message) ToRemoteMessage(from PeerID) RemoteMessage {
	return RemoteMessage{
		From:    from,
		Payload: m.Data,
	}
}

func (m *Message) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
