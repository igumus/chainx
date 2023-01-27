package network

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type tcp struct {
	listenAddr string
	listener   net.Listener
	addPeerCh  chan<- net.Conn
}

func newTCPTransport(addr string, addCh chan<- net.Conn) (Transport, error) {
	t := &tcp{
		listenAddr: addr,
		addPeerCh:  addCh,
	}

	return t, nil
}

func (t *tcp) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"type": t.Type(),
				"addr": t.listenAddr,
				"err":  err,
			}).Error("accepting connection failed")
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
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = ln

	logrus.WithFields(logrus.Fields{
		"type": "tcp",
		"addr": t.listenAddr,
	}).Info("transport started")

	go t.acceptLoop()

	return nil
}

func (t *tcp) Dial(addr string) (net.Conn, error) {
	return net.DialTimeout("tcp", addr, 1*time.Second)
}

func (t *tcp) Close() error {
	logrus.WithFields(logrus.Fields{
		"type": "tcp",
		"addr": t.listenAddr,
	}).Info("closing transport")
	return t.listener.Close()
}
