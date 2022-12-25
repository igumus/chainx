package transport

import "fmt"

type NetAddr string

type RPC struct {
	From    NetAddr
	Payload []byte
}

func (r RPC) String() string {
	return fmt.Sprintf("{From: %s, Payload: %v}", r.From, r.Payload)
}

type Transport interface {
	Addr() NetAddr
	Consume() <-chan RPC
	Connect(to Transport) error
	SendMessage(to NetAddr, payload []byte) error
}
