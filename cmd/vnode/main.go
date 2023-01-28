package main

import (
	"flag"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/network"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	name := flag.String("name", "VNODE", "name of network")
	tcpAddr := flag.String("net-addr", ":3000", "listen address of the tcp transport")
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

	network, err := network.NewNetwork(
		network.WithKeyPair(key),
		network.WithTCPTransport(*tcpAddr),
		network.WithName(*name),
		network.WithDebugMode(*debug),
	)
	if err != nil {
		log.Panic().Err(err)
	}

	network.Start()

	select {}
}
