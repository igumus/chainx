package transport

import (
	"fmt"
	"sync"

	"github.com/igumus/chainx/types"
)

type localTransport struct {
	addr      types.NetAddr
	consumeCh chan types.RPC
	mu        sync.RWMutex
	peers     map[types.NetAddr]*localTransport
}

func NewLocalTransport(anAddr types.NetAddr) Transport {
	return createLocalTransport(anAddr)
}

func createLocalTransport(anAddr types.NetAddr) *localTransport {
	return &localTransport{
		addr:      anAddr,
		consumeCh: make(chan types.RPC, 1024),
		peers:     make(map[types.NetAddr]*localTransport),
	}
}

func (t *localTransport) Consume() <-chan types.RPC {
	return t.consumeCh
}

func (t *localTransport) Addr() types.NetAddr {
	return t.addr
}

func (t *localTransport) Connect(to Transport) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.peers[to.Addr()] = to.(*localTransport)
	return nil
}

func (t *localTransport) SendMessage(to types.NetAddr, payload []byte) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: unknown peer: %s", t.addr, to)
	}
	peer.consumeCh <- types.RPC{
		From:    t.addr,
		Payload: payload,
	}
	return nil
}
