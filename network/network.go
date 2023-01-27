package network

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/types"
	"github.com/sirupsen/logrus"
)

type Network interface {
	ID() string
	Name() string
	Start()
	Dial(addr string) (string, error)
	Consume() <-chan types.RemoteMessage
	io.Closer
}

type network struct {
	id        string
	name      string
	keypair   *crypto.KeyPair
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

func NewNetwork(key *crypto.KeyPair, options ...NetworkOption) (Network, error) {
	n := &network{
		keypair:      key,
		id:           key.Address().String(),
		name:         key.Address().String(),
		addPeerCh:    make(chan net.Conn),
		delPeerCh:    make(chan *peer),
		rpcPeerCh:    make(chan types.RemoteMessage, 1024),
		messageCh:    make(chan types.RemoteMessage, 1024),
		peers:        make(map[types.PeerID]*peer),
		pendingPeers: make(map[types.PeerID]*peer),
	}

	config, err := createOptions(options...)
	if err != nil {
		return nil, err
	}
	n.name = config.name
	n.seedNodes = config.nodes

	if err := n.createTCPTransport(config); err != nil {
		return nil, err
	}

	return n, nil
}

func (n *network) createTCPTransport(config *netOptions) error {
	tr, err := newTCPTransport(config.tcpTransport, n.addPeerCh)
	if err != nil {
		return err
	}

	n.transport = tr

	return nil
}

func (n *network) bootstrap() {
	for _, addr := range n.seedNodes {
		go func(addr string) {
			_, err := n.Dial(addr)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"netID":      n.ID(),
					"netName":    n.Name(),
					"remoteAddr": addr,
					"err":        err,
				}).Error("failed to connect node")
			} else {
				logrus.WithFields(logrus.Fields{
					"netID":      n.ID(),
					"netName":    n.Name(),
					"remoteAddr": addr,
				}).Info("connection success to bootstrap node")
			}
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

func (n *network) processPeerJoin(conn net.Conn, incoming bool) {
	logrus.WithFields(logrus.Fields{
		"netID":    n.ID(),
		"netName":  n.Name(),
		"peerAddr": conn.RemoteAddr(),
	}).Info("peer joined to cluster")

	peer := &peer{
		peerType: conn.LocalAddr().Network(),
		conn:     conn,
		incoming: incoming,
	}

	n.pendingLock.Lock()
	n.pendingPeers[peer.ID()] = peer
	go peer.readLoop(n.delPeerCh, n.rpcPeerCh)
	n.pendingLock.Unlock()
}

func (n *network) processPeerLeave(peer Peer) {
	n.lock.Lock()
	err := peer.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"netID":    n.ID(),
			"netName":  n.Name(),
			"peerAddr": peer.Addr(),
			"err":      err,
		}).Info("peer leaved from cluster failed")
	} else {
		logrus.WithFields(logrus.Fields{
			"netID":    n.ID(),
			"netName":  n.Name(),
			"peerAddr": peer.Addr(),
		}).Info("peer leaved from cluster")
	}
	delete(n.peers, peer.ID())
	n.lock.Unlock()
}

func (n *network) processPeerMessage(rpc types.RemoteMessage) {
	logrus.WithFields(logrus.Fields{
		"netID":   n.ID(),
		"netName": n.Name(),
		"from":    rpc.From,
	}).Info("new message received")
	message, err := rpc.Decode()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"netID":   n.ID(),
			"netName": n.Name(),
			"from":    rpc.From,
			"err":     err,
		}).Error("decoding new message failed")
		return
	}

	switch message.Header {
	case types.NetworkHandshake:
		logrus.WithFields(logrus.Fields{
			"netID":   n.ID(),
			"netName": n.Name(),
			"from":    rpc.From,
		}).Info("received new handshake message")
		if err := n.processPeerHandshake(rpc.From, message.Data); err != nil {
			logrus.WithFields(logrus.Fields{
				"netID":   n.ID(),
				"netName": n.Name(),
				"from":    rpc.From,
				"err":     err,
			}).Error("processing handshake message failed")
		}
	case types.NetworkPeerDiscovered:
		logrus.WithFields(logrus.Fields{
			"netID":   n.ID(),
			"netName": n.Name(),
			"from":    rpc.From,
		}).Info("received new discovery message")
		if err := n.processPeerDiscovery(rpc.From, message.Data); err != nil {
			logrus.WithFields(logrus.Fields{
				"netID":   n.ID(),
				"netName": n.Name(),
				"from":    rpc.From,
				"err":     err,
			}).Error("processing discovery message failed")
		}
	case types.NetworkReserved_3:
	case types.NetworkReserved_4:
	case types.NetworkReserved_5:
		logrus.WithFields(logrus.Fields{
			"netID":   n.ID(),
			"netName": n.Name(),
			"from":    rpc.From,
		}).Warn("received uknown reserved network message")
	default:
		logrus.WithFields(logrus.Fields{
			"netID":   n.ID(),
			"netName": n.Name(),
			"from":    rpc.From,
		}).Warn("received new non-networked type message")
		n.messageCh <- message.ToRemoteMessage(rpc.From)
	}
}

func (n *network) processPeerHandshake(from types.PeerID, rawHandshakeData []byte) error {
	msg, err := decodeHandshakeMessage(rawHandshakeData)
	if err != nil {
		return err
	}

	n.pendingLock.RLock()
	defer n.pendingLock.RUnlock()
	peer, ok := n.pendingPeers[from]
	if !ok {
		return fmt.Errorf("handshaking failed with unknown pending peer: %s", from)
	}
	peer.handshake(msg.Id)
	n.pendingLock.Lock()
	delete(n.pendingPeers, from)
	logrus.WithField("count", len(n.pendingPeers)).Info("pending peer status")
	n.pendingLock.Unlock()

	n.lock.Lock()
	n.peers[peer.ID()] = peer
	n.lock.Unlock()

	discoveryMessage := newDiscoveryMessage(n.ID(), msg)
	discoveryMessageSerialized, err := discoveryMessage.message()
	if err != nil {
		return err
	}

	n.lock.RLock()
	defer n.lock.RUnlock()
	for id, conn := range n.peers {
		if id != peer.ID() {
			go func(pc Peer) {
				err := pc.Send(discoveryMessageSerialized)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"netID":   n.ID(),
						"netName": n.Name(),
						"toPeer":  pc.ID(),
						"err":     err,
					}).Error("sending discovery message failed")
					return
				}
				logrus.WithFields(logrus.Fields{
					"netID":   n.ID(),
					"netName": n.Name(),
					"toPeer":  pc.ID(),
				}).Info("sent discovery message")

			}(conn)
		}
	}
	return nil
}

func (n *network) processPeerDiscovery(from types.PeerID, rawData []byte) error {
	msg, err := decodeDiscoveryMessage(rawData)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"netID":        n.ID(),
		"netName":      n.Name(),
		"from":         from,
		"discoveredBy": msg.From,
		"peerNetID":    msg.Data.Id,
		"peerNetAddr":  msg.Data.Addr,
	}).Info("received discovery message")

	n.lock.RLock()
	defer n.lock.RUnlock()
	_, ok := n.peers[types.PeerID(msg.Data.Id)]
	if ok {
		logrus.WithFields(logrus.Fields{
			"netID":   n.ID(),
			"netName": n.Name(),
			"peer":    msg.Data.Id,
		}).Warn("discovered peer already in known peer list")
		return nil
	}

	conn, err := n.transport.Dial(msg.Data.Addr)
	if err != nil {
		return err
	}

	peer := &peer{
		id:       msg.Data.Id,
		peerType: n.transport.Type(),
		conn:     conn,
		incoming: false,
	}

	n.lock.Lock()
	defer n.lock.Unlock()
	n.peers[types.PeerID(msg.Data.Id)] = peer
	return nil
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

	n.processPeerJoin(conn, false)

	logrus.WithField("addr", addr).Info("waiting one second to send handshake")
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
	logrus.WithFields(logrus.Fields{
		"netID":   n.ID(),
		"netName": n.Name(),
	}).Info("shutting down network")
	n.lock.Lock()
	defer n.lock.Unlock()

	for _, peer := range n.peers {
		n.processPeerLeave(peer)
	}

	return n.transport.Close()
}
