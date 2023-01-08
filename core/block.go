package core

import (
	"bytes"
	"time"

	"github.com/igumus/chainx/hash"
)

type Header struct {
	Version   uint32
	Height    uint32
	Timestamp uint64
	Nonce     uint64
	PrevBlock hash.Hash
}

func EmptyHeader() *Header {
	return &Header{
		Version:   1,
		Height:    0,
		Timestamp: uint64(time.Now().UnixNano()),
		Nonce:     2,
		PrevBlock: hash.ZeroHash,
	}
}

type Block struct {
	Header       *Header
	hash         hash.Hash
	Transactions []*Transaction
}

func NewBlock(h *Header, txs []*Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
		hash:         hash.ZeroHash,
	}
}

func EmptyBlock() *Block {
	return NewBlock(EmptyHeader(), nil)
}

func (b *Block) Hash() hash.Hash {
	buff := new(bytes.Buffer)
	err := NewGobEncoder[Block](buff)(b)
	if err != nil {
		return b.hash
	}
	if b.hash.IsZero() {
		b.hash = hash.CreateHash(buff.Bytes())
	}
	return b.hash
}
