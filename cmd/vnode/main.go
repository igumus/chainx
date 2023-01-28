package main

import (
	"flag"

	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/network"
	"github.com/igumus/chainx/node"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func panicErr(err error) {
	if err != nil {
		log.Panic().Err(err)
	}
}

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	name := flag.String("name", "VNODE", "name of network")
	tcpAddr := flag.String("net-addr", ":3000", "listen address of the tcp transport")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// creating cryptographic keypair
	key, err := crypto.GenerateKeyPair()
	panicErr(err)

	// creating chain network
	network, err := network.New(
		network.WithKeyPair(key),
		network.WithTCPTransport(*tcpAddr),
		network.WithName(*name),
		network.WithDebugMode(*debug),
	)
	panicErr(err)

	// creating blockchain instance
	bc, err := core.NewBlockChain()
	panicErr(err)

	// creating txpool instance
	txpool, err := core.NewTXPool()
	panicErr(err)

	// creating node instance
	server, err := node.New(
		node.WithDebugMode(*debug),
		node.WithKeypair(key),
		node.WithNetwork(network),
		node.WithChain(bc),
		node.WithTXPool(txpool),
		node.EnableValidator(),
	)
	panicErr(err)

	// starting node server
	server.Start()
}
