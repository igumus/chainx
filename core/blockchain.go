package core

import (
	"errors"
	"sync"

	"github.com/igumus/chainx/crypto"
	"github.com/rs/zerolog/log"
)

type BlockChain interface {
	CurrentHeader() *Header
	GetBlocks(uint32) ([]*Block, error)
	CreateBlock(*crypto.KeyPair, []*Transaction) (*Block, error)
	AddBlock(*Block) error
}

type chain struct {
	storage       Storage
	lock          sync.RWMutex
	prevHeader    *Header
	currHeader    *Header
	contractState *State
}

func NewBlockChain() (BlockChain, error) {
	bc := &chain{
		storage:       NewMemoryStorage(),
		contractState: NewState(),
		prevHeader:    nil,
		currHeader:    nil,
	}

	genesis, err := GenesisBlock()
	if err != nil {
		return nil, err
	}

	return bc, bc.addBlock(genesis)
}
func (bc *chain) GetBlocks(from uint32) ([]*Block, error) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.storage.GetAll(from, bc.currHeader.Height)
}

func (bc *chain) CurrentHeader() *Header {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.currHeader
}

func (bc *chain) CreateBlock(key *crypto.KeyPair, txs []*Transaction) (*Block, error) {
	b, err := NewBlock(bc.CurrentHeader(), txs)
	if err != nil {
		return nil, err
	}

	err = b.Sign(key)
	if err != nil {
		return nil, err
	}

	if err := bc.validateBlock(b); err != nil {
		return nil, err
	}
	return b, bc.addBlock(b)
}

func (bc *chain) AddBlock(b *Block) error {
	if err := bc.validateBlock(b); err != nil {
		return err
	}
	return bc.addBlock(b)
}

func (bc *chain) addBlock(b *Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if err := bc.storage.Put(b); err != nil {
		return err
	}
	bc.prevHeader = bc.currHeader
	bc.currHeader = b.Header

	for id, tx := range b.Transactions {
		vm := NewVM(tx.Data, bc.contractState)

		state, err := vm.Run()
		if err != nil {
			return err
		}
		log.Info().
			Str("blockhash", b.Header.Hash().String()).
			Str("txhash", tx.Hash().String()).
			Int("txSeq", id).
			Msg("executed transaction")

		bc.contractState.Merge(state)

	}

	log.Info().
		Str("blockhash", b.Header.Hash().String()).
		Uint32("height", b.Header.Height).
		Int("transactions", len(b.Transactions)).
		Msg("new block")

	return nil
}

var (
	ErrBlockKnown              = errors.New("block already have")
	ErrBlockTooHigh            = errors.New("block too high")
	ErrBlockPrevHeaderNotValid = errors.New("hash of prev block is invalid")
)

func (bc *chain) validateBlock(b *Block) error {
	if b.Header.Height <= bc.currHeader.Height {
		return ErrBlockKnown
	}

	if b.Header.Height > bc.currHeader.Height+1 {
		return ErrBlockTooHigh
	}

	if !b.Header.PrevBlockHash.IsEqual(bc.currHeader.Hash()) {
		return ErrBlockPrevHeaderNotValid
	}

	// TODO (@igumus): add check prev header hash with block prev header hash

	return b.Verify()
}
