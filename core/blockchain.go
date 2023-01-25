package core

import (
	"errors"
	"sync"
)

type BlockChain struct {
	lock       sync.RWMutex
	storage    Storage
	currHeader *Header
}

func NewBlockChain(genesis *Block) (*BlockChain, error) {
	bc := &BlockChain{
		storage:    NewMemoryStorage(),
		currHeader: nil,
	}

	return bc, bc.addBlock(genesis)
}

func (bc *BlockChain) addBlock(b *Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if err := bc.storage.Put(b); err != nil {
		return err
	}
	bc.currHeader = b.Header
	return nil
}

func (bc *BlockChain) validateBlock(b *Block) error {
	if b.Header.Height <= bc.currHeader.Height {
		return errors.New("error known block")
	}

	if b.Header.Height > bc.currHeader.Height+1 {
		return errors.New("block too high")
	}

	if !b.Header.PrevBlockHash.IsEqual(bc.currHeader.Hash()) {
		return errors.New("hash of prev block is invalid")
	}

	return b.Verify()
}

func (bc *BlockChain) CurrentHeader() (*Header, error) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.currHeader, nil
}

func (bc *BlockChain) AddBlock(b *Block) error {
	if err := bc.validateBlock(b); err != nil {
		return err
	}
	return bc.addBlock(b)
}
