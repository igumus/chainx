package main

import (
	"context"
	"time"

	"github.com/igumus/chainx/network"
	"github.com/igumus/chainx/transport"
	"github.com/sirupsen/logrus"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	trLocal := transport.NewLocalTransport("LOCAL")
	trRemote := transport.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	server, err := network.NewServer(trLocal, trRemote)
	if err != nil {
		logrus.WithField("cause", err).Fatal("server instance creation failed")
	}

	go func() {
		for {
			trLocal.SendMessage(trRemote.Addr(), []byte("hello world"))
			time.Sleep(1 * time.Second)
		}
	}()

	server.Start(ctx)
}
