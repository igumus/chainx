package core

import (
	"errors"
	"sync"
)

type Storage interface {
	Put(*Block) error
	Get(height uint32) (*Block, error)
}

type memoryStorage struct {
	lock    sync.RWMutex
	headers []*Header
	blocks  []*Block
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		headers: []*Header{},
		blocks:  []*Block{},
	}
}

func (ms *memoryStorage) Put(b *Block) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.headers = append(ms.headers, b.Header)
	ms.blocks = append(ms.blocks, b)
	return nil
}

func (ms *memoryStorage) Get(h uint32) (*Block, error) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()
	idx := int(h) + 1
	if len(ms.blocks) <= idx {
		return nil, errors.New("height too high")
	}
	return ms.blocks[idx], nil
}
