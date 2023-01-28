package core

import (
	"sort"
	"time"
)

type TXPool interface {
	Add(*Transaction) error
	Contains(*Transaction) bool
	Transactions() []*Transaction
	Size() int
	Flush()
}

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

type pool struct {
	items map[string]*Transaction
}

func NewTXPool() (TXPool, error) {
	return &pool{
		items: make(map[string]*Transaction),
	}, nil
}

func (t *pool) Transactions() []*Transaction {
	txs := newTxSorter(t.items)
	sort.Sort(txs)
	return txs.items
}

func (t *pool) Add(tx *Transaction) error {
	if !t.Contains(tx) {
		if err := tx.Verify(); err != nil {
			return err
		}
		txhash := tx.Hash().String()
		tx.localOrder = time.Now().UnixNano()
		t.items[txhash] = tx
	}

	return nil
}

func (t *pool) Contains(tx *Transaction) bool {
	txhash := tx.Hash().String()
	_, ok := t.items[txhash]
	return ok
}

func (t *pool) Size() int {
	return len(t.items)
}

func (t *pool) Flush() {
	t.items = make(map[string]*Transaction)
}
