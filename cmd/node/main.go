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

	network, err := network.NewNetwork(
		network.WithKeyPair(key),
		network.WithTCPTransport(addr),
		network.WithSeedNode(*bootstrapnode),
		network.WithName(name),
		network.WithDebugMode(*debug),
	)
	if err != nil {
		log.Panic().Err(err)
	}

	network.Start()
	select {}
}
