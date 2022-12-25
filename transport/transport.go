package transport

import "github.com/igumus/chainx/types"

type Transport interface {
	Addr() types.NetAddr
	Consume() <-chan types.RPC
	Connect(to Transport) error
	SendMessage(to types.NetAddr, payload []byte) error
}
