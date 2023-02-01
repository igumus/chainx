package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"time"

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
	debug := flag.Bool("debug", false, "sets log level to debug")
	seq := flag.String("seq", "1", "sequence number of node")
	//tcpAddr := flag.String("net-addr", ":3001", "listen address of the grpc transport")
	bootstrapnode := flag.String("node", ":3000", "seed node listen addr")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// creating cryptographic keypair
	key, err := crypto.GenerateKeyPair()
	panicErr(err)

	name := fmt.Sprintf("NODE_%s", *seq)
	addr := fmt.Sprintf(":300%s", *seq)

	// creating chain network
	network, err := network.New(
		network.WithDebugMode(*debug),
		network.WithName(name),
		network.WithKeyPair(key),
		network.WithTCPTransport(addr),
		network.WithSeedNode(*bootstrapnode),
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
	)
	panicErr(err)

	// starting node server
	go server.Start()

	txTicker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			sendTransaction(key, network)
			<-txTicker.C
		}
	}()

	select {}
}

func sendTransaction(k *crypto.KeyPair, n network.Network) {
	data := make([]byte, 20)
	rand.Reader.Read(data)
	tx := core.NewTransaction(data)
	tx.Sign(k)

	buf := new(bytes.Buffer)
	if err := core.EncodeTransaction(buf, tx); err != nil {
		log.Error().Err(err)
		return
	}

	message := &network.Message{
		Header: network.ChainTx,
		Data:   buf.Bytes(),
	}

	mbuf, err := message.Bytes()
	if err != nil {
		log.Error().Err(err)
		return
	}

	remote := network.RemoteMessage{
		From:    network.PeerID(k.Address().String()),
		Payload: mbuf,
	}
	if err := n.HandleMessage(remote); err != nil {
		log.Error().Err(err)
	}
}
