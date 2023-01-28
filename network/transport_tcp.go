package network

import (
	"net"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type tcp struct {
	listenAddr string
	listener   net.Listener
	addPeerCh  chan<- net.Conn
	logger     zerolog.Logger
}

func newTCPTransport(id string, addr string, addCh chan<- net.Conn) (Transport, error) {
	t := &tcp{
		listenAddr: addr,
		addPeerCh:  addCh,
		logger: log.With().
			Str("id", id).
			Str("component", "transport").
			Str("type", "tcp").
			Str("addr", addr).
			Logger(),
	}

	return t, nil
}

func (t *tcp) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			t.logger.Error().Err(err).Msg("accepting incoming connection failed")
			continue
		}
		t.addPeerCh <- conn
	}
}

func (t *tcp) Type() string {
	return "tcp"
}

func (t *tcp) Addr() string {
	return t.listenAddr
}

func (t *tcp) Listen() error {
	ln, err := net.Listen(t.Type(), t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = ln

	t.logger.Info().Msg("transport started")

	go t.acceptLoop()

	return nil
}

func (t *tcp) Dial(addr string) (net.Conn, error) {
	return net.DialTimeout(t.Type(), addr, 1*time.Second)
}

func (t *tcp) Close() error {
	t.logger.Info().Msg("closing transport")
	return t.listener.Close()
}
