package core

import (
	"bytes"
	"time"

	"github.com/igumus/chainx/types"
)

type Header struct {
	Version   uint32
	Height    uint32
	Timestamp uint64
	Nonce     uint64
	PrevBlock types.Hash
}

func EmptyHeader() *Header {
	return &Header{
		Version:   1,
		Height:    0,
		Timestamp: uint64(time.Now().UnixNano()),
		Nonce:     0,
		PrevBlock: types.ZeroHash,
	}
}

type Block struct {
	Header       *Header
	hash         types.Hash
	Transactions []*Transaction
}

func NewBlock(h *Header, txs []*Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
		hash:         types.ZeroHash,
	}
}

func EmptyBlock() *Block {
	return NewBlock(EmptyHeader(), nil)
}

func (b *Block) Hash() types.Hash {
	buff := new(bytes.Buffer)
	err := NewGobEncoder[Block](buff)(b)
	if err != nil {
		return b.hash
	}
	if b.hash.IsZeroHash() {
		b.hash = types.CreateHash(buff.Bytes())
	}
	return b.hash
}
