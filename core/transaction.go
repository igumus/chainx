package core

import (
	"bytes"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/hash"
)

type Transaction struct {
	Data      []byte
	Signature *crypto.Signature

	// cache purpose
	txhash hash.Hash

	// transaction order no for locally
	localOrder int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data:   data,
		txhash: []byte{},
	}
}

func (tx *Transaction) Hash() hash.Hash {
	if tx.txhash.IsZero() {
		// calculate hash
		buf := new(bytes.Buffer)
		EncodeTransaction(buf, tx)
		tx.txhash = hash.CreateHash(buf.Bytes())
	}
	return tx.txhash
}

func (t *Transaction) Sign(kp *crypto.KeyPair) error {
	signature, err := kp.Sign(t.Data)
	if err != nil {
		return err
	}
	t.Signature = signature
	return nil
}

func (t *Transaction) Verify() error {
	return t.Signature.Verify(t.Data)
}

func calculateTransactionHash(txs ...*Transaction) (hash.Hash, error) {
	buf := new(bytes.Buffer)
	for _, tx := range txs {
		if err := tx.Verify(); err != nil {
			return hash.ZeroHash, err
		}
		if err := EncodeTransaction(buf, tx); err != nil {
			return hash.ZeroHash, err
		}
	}
	return hash.CreateHash(buf.Bytes()), nil
}
