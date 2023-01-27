package network

import (
	"bytes"
	"encoding/gob"

	"github.com/igumus/chainx/types"
)

type networkHandshakeMessage struct {
	Id   string
	Addr string
}

func (n *networkHandshakeMessage) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(n); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n *networkHandshakeMessage) message() ([]byte, error) {
	data, err := n.encode()
	if err != nil {
		return nil, err
	}
	msg := &types.Message{
		Header: types.NetworkHandshake,
		Data:   data,
	}

	return msg.Bytes()
}

func decodeHandshakeMessage(data []byte) (*networkHandshakeMessage, error) {
	msg := &networkHandshakeMessage{}
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type networkHandshakeReplyMessage struct {
	Id   string
	Addr string
}

func (n *networkHandshakeReplyMessage) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(n); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n *networkHandshakeReplyMessage) message() ([]byte, error) {
	data, err := n.encode()
	if err != nil {
		return nil, err
	}
	msg := &types.Message{
		Header: types.NetworkHandshakeReply,
		Data:   data,
	}

	return msg.Bytes()
}

func decodeHandshakeReplyMessage(data []byte) (*networkHandshakeReplyMessage, error) {
	msg := &networkHandshakeReplyMessage{}
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(msg); err != nil {
		return nil, err
	}
	return msg, nil
}
