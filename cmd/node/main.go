package main

import (
	"bytes"
	"flag"
	"fmt"
	"time"

	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/network"
	"github.com/igumus/chainx/node"
	"github.com/igumus/chainx/types"

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

	go func() {
		time.Sleep(3 * time.Second)
		size := 10
		for i := 0; i < size; i++ {
			sendTransaction(key, network, i)
		}
	}()

	select {}
}

func sendTransaction(k *crypto.KeyPair, n network.Network, idx int) {
	data := []byte(fmt.Sprintf("foo_%d", idx))
	tx := core.NewTransaction(data)
	tx.Sign(k)

	buf := new(bytes.Buffer)
	if err := core.EncodeTransaction(buf, tx); err != nil {
		log.Error().Err(err)
		return
	}

	message := &types.Message{
		Header: types.ChainTx,
		Data:   buf.Bytes(),
	}

	mbuf, err := message.Bytes()
	if err != nil {
		log.Error().Err(err)
		return
	}

	remote := types.RemoteMessage{
		From:    types.PeerID(k.Address().String()),
		Payload: mbuf,
	}
	if err := n.HandleMessage(remote); err != nil {
		log.Error().Err(err)
	}
}
