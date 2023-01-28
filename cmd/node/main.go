package main

import (
	"flag"
	"fmt"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/network"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	seq := flag.String("seq", "1", "listen address of the http transport")
	//tcpAddr := flag.String("net-addr", ":3001", "listen address of the grpc transport")
	bootstrapnode := flag.String("node", ":3000", "listen address of the grpc transport")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	key, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Panic().Err(err)
	}

	name := fmt.Sprintf("NODE_%s", *seq)
	addr := fmt.Sprintf(":300%s", *seq)

	network, err := network.NewNetwork(key,
		network.WithName(name),
		network.WithTCPTransport(addr),
		network.WithBootstrapNode(*bootstrapnode))
	if err != nil {
		log.Panic().Err(err)
	}

	network.Start()
	select {}
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()

		trLocal := network.NewLocalTransport("LOCAL")
		trRemote := network.NewLocalTransport("REMOTE")

		trLocal.Connect(trRemote)
		trRemote.Connect(trLocal)

		server, err := node.NewServer(trLocal, trRemote)
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
	*/
}
