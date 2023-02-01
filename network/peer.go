package network

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"

	"github.com/rs/zerolog/log"
)

type PeerState byte

const (
	PendingPeer    PeerState = 0x0
	HandshakedPeer PeerState = 0x1
)

type PeerID string

func (pid PeerID) String() string {
	return string(pid)
}

type Peer interface {
	ID() PeerID
	Type() string
	Addr() string
	Send([]byte) error
	IsOutgoing() bool
	io.Closer
}

type peer struct {
	id       string
	state    PeerState
	peerType string
	conn     net.Conn
	incoming bool
}

func (p *peer) readLoop(delCh chan<- *peer, rpcCh chan<- RemoteMessage) {
	var (
		size int64 = 0
		buf  *bytes.Buffer
	)

	for {
		err := binary.Read(p.conn, binary.LittleEndian, &size)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error().Err(err).Str("type", p.Type()).Str("addr", p.Addr()).Msg("reading from peer failed")
			continue
		}
		buf = new(bytes.Buffer)
		n, err := io.CopyN(buf, p.conn, size)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error().Err(err).Str("type", p.Type()).Str("addr", p.Addr()).Msg("copying from peer failed")
			continue
		}

		if e := log.Debug(); e.Enabled() {
			log.Debug().Str("type", p.Type()).Str("addr", p.Addr()).Int64("readBytes", n).Msg("incoming message accepted")
		}

		from := p.ID()
		if p.state == PendingPeer {
			from = PeerID(p.Addr())
		}

		rpcCh <- RemoteMessage{
			From:    from,
			Payload: buf.Bytes(),
		}

	}
	delCh <- p
}

func (p *peer) handshake(id string) {
	p.id = id
	p.state = HandshakedPeer
	log.Info().Str("peer", p.id).Str("addr", p.conn.RemoteAddr().String()).Msg("changed peer state to handshaked")
}

func (p *peer) IsOutgoing() bool {
	return !p.incoming
}

func (p *peer) Type() string {
	return p.peerType
}

func (p *peer) ID() PeerID {
	return PeerID(p.id)
}

func (p *peer) Addr() string {
	return p.conn.RemoteAddr().String()
}

func (p *peer) Close() error {
	return p.conn.Close()
}

func (p *peer) Send(b []byte) error {
	err := binary.Write(p.conn, binary.LittleEndian, int64(len(b)))
	if err != nil {
		return err
	}
	_, err = p.conn.Write(b)
	return err
}
