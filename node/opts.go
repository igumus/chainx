package node

import (
	"errors"
	"time"

	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/network"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type NodeOption func(*nodeOption)

type nodeOption struct {
	debugMode     bool
	blockTime     time.Duration
	network       network.Network
	validatorNode bool
	keypair       *crypto.KeyPair
	logger        zerolog.Logger
	chain         core.BlockChain
	pool          core.TXPool
}

func (n *nodeOption) validate() error {
	if n.network == nil {
		return errors.New("network not specified")
	}
	if n.keypair == nil {
		return errors.New("keypair not specified")
	}
	return nil
}

func createOptions(opts ...NodeOption) (*nodeOption, error) {
	options := &nodeOption{
		blockTime: 5 * time.Second,
		keypair:   nil,
		chain:     nil,
		pool:      nil,
	}

	for _, opt := range opts {
		opt(options)
	}

	if err := options.validate(); err != nil {
		return nil, err
	}

	// if chain implementation not specified
	// instead of returning error, we will use
	// default chain implementation
	if options.chain == nil {
		c, err := core.NewBlockChain()
		if err != nil {
			return nil, err
		}
		options.chain = c
	}

	// if transaction pool implementation not specified,
	// instead of returning error, we will use
	// default transaction pool implementation
	if options.pool == nil {
		p, err := core.NewTXPool()
		if err != nil {
			return nil, err
		}
		options.pool = p
	}

	// create zerolog logger instance
	options.logger = log.With().
		Str("id", options.keypair.Address().String()).
		Str("component", "node").
		Logger()

	return options, nil
}

func WithBlockTime(t time.Duration) NodeOption {
	return func(nc *nodeOption) {
		nc.blockTime = t
	}
}

func WithNetwork(n network.Network) NodeOption {
	return func(no *nodeOption) {
		no.network = n
	}
}

func EnableValidator() NodeOption {
	return func(no *nodeOption) {
		no.validatorNode = !no.validatorNode
	}
}

func WithKeypair(k *crypto.KeyPair) NodeOption {
	return func(no *nodeOption) {
		no.keypair = k
	}
}

func WithDebugMode(b bool) NodeOption {
	return func(no *nodeOption) {
		no.debugMode = b
	}
}

func WithChain(bc core.BlockChain) NodeOption {
	return func(no *nodeOption) {
		no.chain = bc
	}
}

func WithTXPool(tp core.TXPool) NodeOption {
	return func(no *nodeOption) {
		no.pool = tp
	}
}
