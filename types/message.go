package types

import "fmt"

type NetAddr string

type RPC struct {
	From    NetAddr
	Payload []byte
}

func (r RPC) String() string {
	return fmt.Sprintf("{From: %s, Payload: %v}", r.From, r.Payload)
}
