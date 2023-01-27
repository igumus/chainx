package network

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"

	"github.com/igumus/chainx/types"
	"github.com/sirupsen/logrus"
)

type Peer interface {
	ID() types.PeerID
	Type() string
	Addr() string
	Send([]byte) error
	IsOutgoing() bool
	io.Closer
}

type peer struct {
	id       string
	state    types.PeerState
	peerType string
	conn     net.Conn
	incoming bool
}

func (p *peer) readLoop(delCh chan<- *peer, rpcCh chan<- types.RemoteMessage) {
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
			logrus.WithFields(logrus.Fields{
				"type": p.Type(),
				"addr": p.Addr(),
				"err":  err,
			}).Error("reading from peer failed")
			continue
		}
		buf = new(bytes.Buffer)
		n, err := io.CopyN(buf, p.conn, size)
		if err != nil {
			if err == io.EOF {
				break
			}
			logrus.WithFields(logrus.Fields{
				"type": p.Type(),
				"addr": p.Addr(),
				"err":  err,
			}).Error("copying from peer failed")
			continue
		}

		logrus.WithFields(logrus.Fields{
			"type":      p.Type(),
			"addr":      p.Addr(),
			"readBytes": n,
		}).Debug("incoming message accepted")

		from := p.ID()
		if p.state == types.PendingPeer {
			from = types.PeerID(p.Addr())
		}

		rpcCh <- types.RemoteMessage{
			From:    from,
			Payload: buf.Bytes(),
		}

	}
	delCh <- p
}

func (p *peer) handshake(id string) {
	p.id = id
}

func (p *peer) IsOutgoing() bool {
	return !p.incoming
}

func (p *peer) Type() string {
	return p.peerType
}

func (p *peer) ID() types.PeerID {
	return types.PeerID(p.id)
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
