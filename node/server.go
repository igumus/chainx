package node

import (
	"bytes"
	"encoding/gob"
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
	id        types.PeerID
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
		keypair:   options.keypair,
		logger:    options.logger,
		debug:     options.debugMode,
		validator: options.validatorNode,
		blockTime: options.blockTime,
		txpool:    options.pool,
		chain:     options.chain,
		network:   options.network,
		id:        types.PeerID(options.network.ID()),
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

	n.logger.Info().Int("pendingTXcount", pendingSize).Msg("try to create block with txs")
	block, err := n.chain.CreateBlock(n.keypair, txs)
	if err != nil {
		return err
	}

	n.txpool.Flush()

	if err := n.broadcastBlock("", block); err != nil {
		return err
	}

	return nil
}

func (n *node) fetchBlock(peer types.PeerID, remoteHeight uint32) error {
	n.logger.Info().Str("peer", peer.String()).Uint32("ownHeight", n.chain.CurrentHeader().Height).Uint32("blockHeight", remoteHeight).Msg("fetching blocks")

	nextHeight := n.chain.CurrentHeader().Height + 1
	msg, err := NewFetchBlockMessage(types.PeerID(n.id), nextHeight, remoteHeight)
	if err != nil {
		n.logger.Error().Str("peer", peer.String()).Err(err).Msg("creating fetch message failed")
		return err
	}

	if err := n.network.Send(peer, msg); err != nil {
		n.logger.Error().Str("peer", peer.String()).Err(err).Msg("sending fetch message failed")
		return err
	}

	return nil
}

func (n *node) processBlock(peer types.PeerID, block *core.Block) error {
	n.logger.Info().Str("peer", peer.String()).Str("bHash", block.Header.Hash().String()).Msg("new block arrived")
	if err := n.chain.AddBlock(block); err != nil {
		if err == core.ErrBlockTooHigh {
			go n.fetchBlock(peer, block.Header.Height)
			return nil
		}
		if err == core.ErrBlockKnown {
			return nil
		}
		n.logger.Error().Str("peer", peer.String()).Err(err).Msg("processing block failed")
		return err
	}
	n.logger.Info().Str("peer", peer.String()).Str("bHash", block.Header.Hash().String()).Msg("new block saved")

	if err := n.broadcastBlock(peer, block); err != nil {
		n.logger.Error().Err(err).Msg("broadcasting block failed")
		return err
	}

	return nil
}

func (n *node) broadcastBlock(from types.PeerID, block *core.Block) error {
	buf := new(bytes.Buffer)
	if err := core.EncodeBlock(buf, block); err != nil {
		return err
	}

	message := &types.Message{
		Header: types.ChainBlock,
		Data:   buf.Bytes(),
	}

	if err := n.network.Broadcast(message, from); err != nil {
		n.logger.Error().Err(err).Msg("broadcasting block failed")
		return err
	}
	return nil
}

func (n *node) processTransaction(peer types.PeerID, tx *core.Transaction) error {
	if err := n.txpool.Add(tx); err != nil {
		return err
	}

	go n.broadcastTransaction(peer, tx)

	return nil
}

func (n *node) broadcastTransaction(from types.PeerID, tx *core.Transaction) error {
	buf := new(bytes.Buffer)
	if err := core.EncodeTransaction(buf, tx); err != nil {
		return err
	}

	message := &types.Message{
		Header: types.ChainTx,
		Data:   buf.Bytes(),
	}

	if err := n.network.Broadcast(message, from); err != nil {
		n.logger.Error().Err(err).Msg("broadcasting block failed")
		return err
	}
	return nil
}

func (n *node) processBlockFetch(peer types.PeerID, payload *FetchBlockMessage) error {
	blocks, err := n.chain.GetBlocks(payload.From)
	if err != nil {
		n.logger.Error().Err(err).Str("peer", peer.String()).Uint32("fetchBlockFrom", payload.From).Msg("fetching block from chain failed")
		return err
	}

	reply, err := NewFetchBlockReply(blocks)
	if err != nil {
		n.logger.Error().Err(err).Str("peer", peer.String()).Msg("creating fetch block reply message failed")
		return err
	}

	if err := n.network.Send(peer, reply); err != nil {
		n.logger.Error().Err(err).Str("peer", peer.String()).Msg("sending reply message to peer failed")
		return err
	}

	n.logger.Info().Int("blockCount", len(blocks)).Str("peer", peer.String()).Msg("fetch block reply sent")

	return nil
}

func (n *node) processSyncBlock(peer types.PeerID, payload *FetchBlockReply) error {
	n.logger.Info().Str("from", peer.String()).Int("count", len(payload.Blocks)).Msg("sync blocks arrived")

	for _, block := range payload.Blocks {
		if err := n.chain.AddBlock(block); err != nil {
			n.logger.Error().Err(err).Msg("sync block failed")
			return err
		}
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
	case types.ChainFetchBlock:
		data := &FetchBlockMessage{}
		if err := gob.NewDecoder(payload).Decode(data); err != nil {
			return err
		}
		return n.processBlockFetch(peer, data)
	case types.ChainFetchBlockReply:
		data := &FetchBlockReply{}
		if err := gob.NewDecoder(payload).Decode(data); err != nil {
			return err
		}
		return n.processSyncBlock(peer, data)
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
