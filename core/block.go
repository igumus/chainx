package core

import (
	"bytes"
	"errors"
	"time"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/hash"
	"github.com/sirupsen/logrus"
)

var ErrInvalidDataHash = errors.New("block has invalid data hash")

type Header struct {
	Version       uint32
	Height        uint32
	Timestamp     int64
	PrevBlockHash hash.Hash
	DataHash      hash.Hash
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	EncodeHeader(buf, h)
	return buf.Bytes()
}

func (h *Header) Hash() hash.Hash {
	return hash.CreateHash(h.Bytes())
}

type Block struct {
	Header       *Header
	Transactions []*Transaction
	Signature    *crypto.Signature
}

func NewBlock(prevHeader *Header, txs []*Transaction) (*Block, error) {
	dataHash, err := calculateTransactionHash(txs...)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       prevHeader.Version,
		Height:        prevHeader.Height + 1,
		Timestamp:     time.Now().UnixNano(),
		PrevBlockHash: prevHeader.Hash(),
		DataHash:      dataHash,
	}

	block := &Block{
		Header:       header,
		Transactions: txs,
	}

	return block, nil
}

func (b *Block) Sign(kp *crypto.KeyPair) error {
	data := b.Header.Hash()
	signature, err := kp.Sign(data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"block": data,
			"cause": err,
		}).Error("block sign")
		return err
	}
	b.Signature = signature
	return nil
}

func (b *Block) Verify() error {
	data := b.Header.Hash()
	if err := b.Signature.Verify(data); err != nil {
		logrus.WithFields(logrus.Fields{
			"block": data,
			"phase": "signature verify",
			"cause": err,
		}).Error("block verification failed")
		return err
	}
	datahash, err := calculateTransactionHash(b.Transactions...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"block": data,
			"phase": "calculate tx hash",
			"cause": err,
		}).Error("block verification failed")
		return err
	}

	if !b.Header.DataHash.IsEqual(datahash) {
		logrus.WithFields(logrus.Fields{
			"block": data,
			"phase": "datahash verification",
		}).Error("block verification failed")
		return ErrInvalidDataHash
	}
	return nil
}
