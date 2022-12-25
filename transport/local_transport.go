package transport

import (
	"fmt"
	"sync"
)

type localTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	mu        sync.RWMutex
	peers     map[NetAddr]*localTransport
}

func NewLocalTransport(anAddr NetAddr) Transport {
	return createLocalTransport(anAddr)
}

func createLocalTransport(anAddr NetAddr) *localTransport {
	return &localTransport{
		addr:      anAddr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*localTransport),
	}
}

func (t *localTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *localTransport) Addr() NetAddr {
	return t.addr
}

func (t *localTransport) Connect(to Transport) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.peers[to.Addr()] = to.(*localTransport)
	return nil
}

func (t *localTransport) SendMessage(to NetAddr, payload []byte) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: unknown peer: %s", t.addr, to)
	}
	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: payload,
	}
	return nil
}
