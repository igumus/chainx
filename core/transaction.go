package core

import (
	"bytes"
	"fmt"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/hash"
)

type Transaction struct {
	Data      []byte
	Signature *crypto.Signature
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash() hash.Hash {
	return hash.CreateHash(tx.Data)
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

func calculateTransactionHash(txs []*Transaction) (hash.Hash, error) {
	buf := new(bytes.Buffer)
	for _, tx := range txs {
		if err := tx.Verify(); err != nil {
			fmt.Printf("tx verification failed: %s\n", err)
			return hash.ZeroHash, err
		}
		if _, err := buf.Write(tx.Data); err != nil {
			return hash.ZeroHash, err
		}
	}
	return hash.CreateHash(buf.Bytes()), nil
}
