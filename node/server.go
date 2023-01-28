package node

import (
	"time"

	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/network"
	"github.com/igumus/chainx/types"
	"github.com/rs/zerolog"
)

type Node interface {
	Start()
}

type node struct {
	id        string
	debug     bool
	validator bool
	blockTime time.Duration
	keypair   *crypto.KeyPair
	txpool    core.TXPool
	chain     core.BlockChain
	network   network.Network
	logger    zerolog.Logger
	messageCh <-chan types.RemoteMessage
	quitCh    chan struct{}
}

func New(opts ...NodeOption) (Node, error) {
	options, err := createOptions(opts...)
	if err != nil {
		return nil, err
	}

	txpool, err := core.NewTXPool()
	if err != nil {
		return nil, err
	}

	chain, err := core.NewBlockChain()
	if err != nil {
		return nil, err
	}

	n := &node{
		id:        options.keypair.Address().String(),
		keypair:   options.keypair,
		logger:    options.logger,
		debug:     options.debugMode,
		validator: options.validatorNode,
		blockTime: options.blockTime,
		txpool:    txpool,
		chain:     chain,
		network:   options.network,
		messageCh: options.network.Consume(),
		quitCh:    make(chan struct{}, 1),
	}

	return n, nil
}

func (n *node) Start() {
	n.network.Start()
	time.Sleep(1 * time.Second)
	n.logger.Info().Msg("network started")

	if n.validator {
		go n.validatorLoop()
		n.logger.Info().Dur("blockTime", n.blockTime).Msg("validator loop started")
	}

free:
	for {
		select {
		case msg := <-n.messageCh:
			n.logger.Info().Any("remoteMessage", msg).Msg("new chain message received")
		case <-n.quitCh:
			break free
		}
	}
	n.shutdown()
}

func (n *node) validatorLoop() {
	ticker := time.NewTicker(n.blockTime)
	for {
		<-ticker.C
		n.logger.Info().Msg("doing operation every tick")
	}
}

func (n *node) shutdown() {
	n.logger.Info().Msg("shutdown process starting")
}
