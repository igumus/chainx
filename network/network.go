package network

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/types"
	"github.com/rs/zerolog"
)

type Network interface {
	ID() string
	Name() string
	Start()
	Dial(addr string) (string, error)
	Consume() <-chan types.RemoteMessage
	Send(peer types.PeerID, msg *types.Message) error
	Broadcast(msg *types.Message) error
	io.Closer
}

type network struct {
	id        string
	name      string
	keypair   *crypto.KeyPair
	logger    zerolog.Logger
	seedNodes []string
	transport Transport
	addPeerCh chan net.Conn
	delPeerCh chan *peer
	rpcPeerCh chan types.RemoteMessage
	messageCh chan types.RemoteMessage

	lock  sync.RWMutex
	peers map[types.PeerID]*peer

	pendingLock  sync.RWMutex
	pendingPeers map[types.PeerID]*peer
}

func NewNetwork(options ...NetworkOption) (Network, error) {
	config, err := createOptions(options...)
	if err != nil {
		return nil, err
	}

	n := &network{
		logger:       config.logger,
		keypair:      config.keypair,
		id:           config.id,
		name:         config.name,
		seedNodes:    config.nodes,
		addPeerCh:    make(chan net.Conn),
		delPeerCh:    make(chan *peer),
		rpcPeerCh:    make(chan types.RemoteMessage, 1024),
		messageCh:    make(chan types.RemoteMessage, 1024),
		peers:        make(map[types.PeerID]*peer),
		pendingPeers: make(map[types.PeerID]*peer),
	}

	tr, err := newTCPTransport(n.id, config.tcpTransport, n.addPeerCh)
	if err != nil {
		return nil, err
	}
	n.transport = tr

	return n, nil
}

func (n *network) bootstrap() {
	for _, addr := range n.seedNodes {
		go func(addr string) {
			_, err := n.Dial(addr)
			if err != nil {
				n.logger.Error().Str("remote", addr).Err(err)
				return
			}
			n.logger.Info().Str("remote", addr).Msg("connected to seed node")
		}(addr)
	}
}

func (n *network) process() {
	for {
		select {
		case conn := <-n.addPeerCh:
			n.processPeerJoin(conn, true)
		case peer := <-n.delPeerCh:
			n.processPeerLeave(peer)
		case rpc := <-n.rpcPeerCh:
			n.processPeerMessage(rpc)
		}
	}
}

func (n *network) processPeerJoin(conn net.Conn, incoming bool) *peer {
	n.logger.Info().Str("peerAddr", conn.RemoteAddr().String()).Msg("peer joined to cluster")

	peer := &peer{
		peerType: conn.LocalAddr().Network(),
		state:    types.PendingPeer,
		conn:     conn,
		incoming: incoming,
	}

	id := types.PeerID(conn.RemoteAddr().String())

	n.pendingLock.Lock()
	n.logger.Info().Str("peerAddr", conn.RemoteAddr().String()).Int("count", len(n.pendingPeers)).Msg("peer joined to cluster")
	n.pendingPeers[id] = peer
	go peer.readLoop(n.delPeerCh, n.rpcPeerCh)
	n.pendingLock.Unlock()
	return peer
}

func (n *network) processPeerLeave(peer Peer) {
	n.lock.Lock()
	err := peer.Close()
	if err != nil {
		n.logger.Error().Str("peerAddr", peer.Addr()).Err(err).Msg("closing peer failed")
	} else {
		n.logger.Info().Str("peerAddr", peer.Addr()).Msg("peer closed")
	}
	delete(n.peers, peer.ID())
	n.lock.Unlock()
}

func (n *network) processPeerMessage(rpc types.RemoteMessage) {
	message, err := rpc.Decode()
	if err != nil {
		n.logger.Error().Str("from", rpc.From.String()).Err(err).Msg("decoding remote message failed")
		return
	}

	switch message.Header {
	case types.NetworkHandshake:
		n.logger.Info().Str("from", rpc.From.String()).Msg("received new handshake message")
		if err := n.processPeerHandshake(rpc.From, message.Data, false); err != nil {
			n.logger.Error().Str("from", rpc.From.String()).Err(err).Msg("processing handshake message failed")
		}
	case types.NetworkHandshakeReply:
		n.logger.Info().Str("from", rpc.From.String()).Msg("received handshake reply message")
		if err := n.processPeerHandshake(rpc.From, message.Data, true); err != nil {
			n.logger.Error().Str("from", rpc.From.String()).Err(err).Msg("processing handshake reply message failed")
		}
	case types.NetworkReserved_2:
		n.logger.Warn().Str("from", rpc.From.String()).Str("type", "NetworkReserved_2").Msg("unhandled network message")
	case types.NetworkReserved_3:
		n.logger.Warn().Str("from", rpc.From.String()).Str("type", "NetworkReserved_3").Msg("unhandled network message")
	case types.NetworkReserved_4:
		n.logger.Warn().Str("from", rpc.From.String()).Str("type", "NetworkReserved_4").Msg("unhandled network message")
	case types.NetworkReserved_5:
		n.logger.Warn().Str("from", rpc.From.String()).Str("type", "NetworkReserved_5").Msg("unhandled network message")
	default:
		n.logger.Info().Str("from", rpc.From.String()).Msg("not a network message forwaring to message channel")
		n.messageCh <- message.ToRemoteMessage(rpc.From)
	}
}

func (n *network) processPeerHandshake(from types.PeerID, rawHandshakeData []byte, reply bool) error {
	msg, err := decodeHandshakeMessage(rawHandshakeData)
	if err != nil {
		return err
	}

	n.pendingLock.RLock()
	peer, ok := n.pendingPeers[from]
	n.pendingLock.RUnlock()
	if !ok {
		return fmt.Errorf("handshaking failed with unknown pending peer: %s", from)
	}
	peer.handshake(msg.Id)

	n.pendingLock.Lock()
	delete(n.pendingPeers, from)
	n.logger.Info().Int("count", len(n.pendingPeers)).Msg("pending peers")
	n.pendingLock.Unlock()

	n.lock.Lock()
	n.peers[peer.ID()] = peer
	n.lock.Unlock()

	if reply {
		n.logger.Info().Str("toNet", peer.ID().String()).Msg("full handshake established")
		return nil
	}

	handshakeReply := &networkHandshakeReplyMessage{
		Id:   n.ID(),
		Addr: n.transport.Addr(),
	}

	hmsg, err := handshakeReply.message()
	if err != nil {
		return err
	}

	return peer.Send(hmsg)
}

func (n *network) Start() {
	n.transport.Listen()
	go n.process()
	time.Sleep(1 * time.Second)
	n.bootstrap()
}

func (n *network) Consume() <-chan types.RemoteMessage {
	return n.messageCh
}

func (n *network) ID() string {
	return n.id
}

func (n *network) Name() string {
	return n.name
}

func (n *network) Dial(addr string) (string, error) {
	conn, err := n.transport.Dial(addr)
	if err != nil {
		return "", err
	}

	peer := n.processPeerJoin(conn, false)

	n.logger.Info().Str("addr", addr).Msg("waiting one second to send handshake")
	time.Sleep(1 * time.Second)

	handshake := &networkHandshakeMessage{
		Id:   n.ID(),
		Addr: n.transport.Addr(),
	}

	msg, err := handshake.message()
	if err != nil {
		return "", err
	}
	if err := peer.Send(msg); err != nil {
		return "", err
	}

	return peer.Addr(), nil
}

func (n *network) Close() error {
	n.logger.Info().Msg("shutdown network")
	n.lock.Lock()
	defer n.lock.Unlock()

	for _, peer := range n.peers {
		n.processPeerLeave(peer)
	}

	return n.transport.Close()
}

func (n *network) Send(to types.PeerID, msg *types.Message) error {
	dmsg, err := msg.Bytes()
	if err != nil {
		return err
	}

	n.lock.RLock()
	defer n.lock.RUnlock()

	peer, ok := n.peers[to]
	if !ok {
		// TODO (@igumus): maybe we should include peer searching mechanism
		return fmt.Errorf("unknown peer: %s", to)
	}
	return peer.Send(dmsg)
}

func (n *network) Broadcast(msg *types.Message) error {
	dmsg, err := msg.Bytes()
	if err != nil {
		return err
	}

	n.lock.RLock()
	defer n.lock.RUnlock()
	for _, peer := range n.peers {
		go func(p Peer) {
			if err := p.Send(dmsg); err != nil {
				n.logger.Error().Str("peer", p.ID().String()).Err(err).Msg("sending broadcast message failed")
			}
		}(peer)
	}
	return nil
}
