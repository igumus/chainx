package network

import (
	"context"
	"errors"
	"time"

	"github.com/igumus/chainx/transport"
	"github.com/igumus/chainx/types"
	"github.com/sirupsen/logrus"
)


var ErrNoTransportSpecified = errors.New("server needs at least one transport")

type Server interface {
	Start(context.Context)
}

type server struct {
	transportLayer []transport.Transport
	rpcChan        chan types.RPC
}

func NewServer(transports ...transport.Transport) (Server, error) {
	if len(transports) == 0 {
		return nil, ErrNoTransportSpecified
	}
	return &server{
		transportLayer: transports,
		rpcChan:        make(chan types.RPC, 1024),
	}, nil
}

func (s *server) startTransportLayer() {
    logrus.Info("starting chainx transport layer")
	for _, tr := range s.transportLayer {
		logrus.WithField("addr", tr.Addr()).Info("starting chainx transport item")
		go func(tr transport.Transport) {
			for msg := range tr.Consume() {
				s.rpcChan <- msg
			}
		}(tr)
	}
}

func (s *server) shutdown() {
	logrus.Warn("nothing todo, shutdown process not implemented yet")
}

func (s *server) Start(ctx context.Context) {
    logrus.Info("starting chainx server layer")
	s.startTransportLayer()
	ticker := time.NewTicker(3 * time.Second)
free:
	for {
		select {
		case rpc := <-s.rpcChan:
			logrus.WithField("rpc", rpc).Info("new message arrived")
		case <-ctx.Done():
			logrus.Warn("context cancelled")
			break free
		case <-ticker.C:
			logrus.Info("do stuff every x seconds")
		}
	}
	s.shutdown()
}
