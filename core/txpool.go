package core

import (
	"sync"

	"github.com/igumus/chainx/hash"
)

type TXPool interface {
	Add(*Transaction) error
	Contains(*Transaction) bool
	Transactions() []*Transaction
	Size() int
	Flush()
}

/*
type txsorter struct {
	items []*Transaction
}

func newTxSorter(txs map[string]*Transaction) *txsorter {
	s := &txsorter{
		items: make([]*Transaction, len(txs)),
	}
	idx := 0
	for _, tx := range txs {
		s.items[idx] = tx
		idx++
	}
	return s
}

func (s *txsorter) Len() int {
	return len(s.items)
}

func (s *txsorter) Less(i, j int) bool {
	return s.items[i].localOrder < s.items[j].localOrder
}

func (s *txsorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}
*/

type pool struct {
	lock   sync.RWMutex
	lookup map[string]hash.Hash
	items  []*Transaction
}

func NewTXPool() (TXPool, error) {
	return &pool{
		lookup: make(map[string]hash.Hash),
		items:  []*Transaction{},
	}, nil
}

func (t *pool) Transactions() []*Transaction {
	return t.items
}

func (t *pool) Add(tx *Transaction) error {
	if !t.Contains(tx) {
		if err := tx.Verify(); err != nil {
			return err
		}

		t.lock.Lock()
		txhash := tx.Hash()
		t.lookup[txhash.String()] = txhash
		t.items = append(t.items, tx)
		t.lock.Unlock()
	}

	return nil
}

func (t *pool) Contains(tx *Transaction) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	txhash := tx.Hash().String()
	_, ok := t.lookup[txhash]
	return ok
}

func (t *pool) Size() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return len(t.items)
}

func (t *pool) Flush() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.lookup = make(map[string]hash.Hash)
	t.items = []*Transaction{}
}
