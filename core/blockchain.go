package core

import (
	"errors"
	"sync"

	"github.com/igumus/chainx/crypto"
)

type BlockChain interface {
	CurrentHeader() *Header
	AddBlock(*crypto.KeyPair, []*Transaction) error
}

type chain struct {
	storage    Storage
	lock       sync.RWMutex
	prevHeader *Header
	currHeader *Header
}

func NewBlockChain() (BlockChain, error) {
	bc := &chain{
		storage:    NewMemoryStorage(),
		prevHeader: nil,
		currHeader: nil,
	}

	genesis, err := GenesisBlock()
	if err != nil {
		return nil, err
	}

	return bc, bc.addBlock(genesis)
}

func (bc *chain) CurrentHeader() *Header {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.currHeader
}

func (bc *chain) AddBlock(key *crypto.KeyPair, txs []*Transaction) error {
	b, err := NewBlock(bc.CurrentHeader(), txs)
	if err != nil {
		return err
	}

	err = b.Sign(key)
	if err != nil {
		return err
	}

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
	return nil
}

func (bc *chain) validateBlock(b *Block) error {
	if b.Header.Height <= bc.currHeader.Height {
		return errors.New("error known block")
	}

	if b.Header.Height > bc.currHeader.Height+1 {
		return errors.New("block too high")
	}

	if !b.Header.PrevBlockHash.IsEqual(bc.currHeader.Hash()) {
		return errors.New("hash of prev block is invalid")
	}

	// TODO (@igumus): add check prev header hash with block prev header hash

	return b.Verify()
}
