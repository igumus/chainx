package node

import (
	"bytes"
	"time"

	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/network"
	"github.com/igumus/chainx/types"
	"github.com/rs/zerolog"
)

type Node interface {
	Start()
	types.RemoteMessageHandler
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

	n := &node{
		id:        options.keypair.Address().String(),
		keypair:   options.keypair,
		logger:    options.logger,
		debug:     options.debugMode,
		validator: options.validatorNode,
		blockTime: options.blockTime,
		txpool:    options.pool,
		chain:     options.chain,
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
			if err := n.HandleMessage(msg); err != nil {
				n.logger.Error().Err(err).Str("from", msg.From.String()).Msg("processing incoming message failed")
			}
		case <-n.quitCh:
			break free
		}
	}
	n.shutdown()
}

func (n *node) createBlock() error {
	txs := n.txpool.Transactions()
	pendingSize := len(txs)
	if pendingSize == 0 {
		if n.debug {
			n.logger.Debug().Int("pendingTXcount", pendingSize).Msg("skipping block creation")
		}
		return nil
	}

	n.logger.Info().Int("pendingTXcount", pendingSize).Msg("try to create block with txs")
	block, err := n.chain.CreateBlock(n.keypair, txs)
	if err != nil {
		return err
	}

	n.txpool.Flush()

	go n.broadcastBlock(block)

	return nil
}

func (n *node) processBlock(peer types.PeerID, block *core.Block) error {
	if err := n.chain.AddBlock(block); err != nil {
		if err == core.ErrBlockTooHigh {
			n.logger.Info().Str("peer", peer.String()).Uint32("ownHeight", n.chain.CurrentHeader().Height).Uint32("blockHeight", block.Header.Height).Msg("should get non existing block(s)")
			return nil
		}
		if err == core.ErrBlockKnown {
			return nil
		}
		n.logger.Error().Str("peer", peer.String()).Err(err).Msg("processing block failed")
		return err
	}

	go n.broadcastBlock(block)

	return nil
}

func (n *node) broadcastBlock(block *core.Block) error {
	buf := new(bytes.Buffer)
	if err := core.EncodeBlock(buf, block); err != nil {
		return err
	}

	message := &types.Message{
		Header: types.ChainBlock,
		Data:   buf.Bytes(),
	}

	if err := n.network.Broadcast(message); err != nil {
		n.logger.Error().Err(err).Msg("broadcasting block failed")
		return err
	}
	return nil
}

// TODO (@igumus): should avoid infinite tx broadcasting
func (n *node) processTransaction(peer types.PeerID, tx *core.Transaction) error {
	if err := n.txpool.Add(tx); err != nil {
		return err
	}

	go n.broadcastTransaction(tx)

	return nil
}

func (n *node) broadcastTransaction(tx *core.Transaction) error {
	buf := new(bytes.Buffer)
	if err := core.EncodeTransaction(buf, tx); err != nil {
		return err
	}

	message := &types.Message{
		Header: types.ChainTx,
		Data:   buf.Bytes(),
	}

	if err := n.network.Broadcast(message); err != nil {
		n.logger.Error().Err(err).Msg("broadcasting block failed")
		return err
	}
	return nil
}

func (n *node) HandleMessage(msg types.RemoteMessage) error {
	decodedMessage, err := msg.Decode()
	if err != nil {
		return err
	}

	peer := msg.From
	payload := bytes.NewReader(decodedMessage.Data)

	switch decodedMessage.Header {
	case types.ChainTx:
		data := &core.Transaction{}
		if err := core.DecodeTransaction(payload, data); err != nil {
			return err
		}
		return n.processTransaction(peer, data)
	case types.ChainBlock:
		data := &core.Block{}
		if err := core.DecodeBlock(payload, data); err != nil {
			return err
		}
		return n.processBlock(peer, data)
	default:
		n.logger.Error().Str("peer", peer.String()).Msg("unknown chain message header")
	}

	return nil
}

func (n *node) validatorLoop() {
	ticker := time.NewTicker(n.blockTime)
	for {
		<-ticker.C
		if err := n.createBlock(); err != nil {
			n.logger.Error().Err(err).Msg("creating block failed")
		}
	}
}

func (n *node) shutdown() {
	n.logger.Info().Msg("shutdown process starting")
}
