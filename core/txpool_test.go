package core

import (
	"fmt"
	"testing"

	"github.com/igumus/chainx/crypto"
	"github.com/stretchr/testify/assert"
)

func TestTransactionPoolCreate(t *testing.T) {
	txpool, err := NewTXPool()
	assert.Nil(t, err)
	assert.Equal(t, 0, txpool.Size())
}

func TestTransactionPoolAdd(t *testing.T) {
	keypair, err := crypto.GenerateKeyPair()
	assert.Nil(t, err)

	txpool, err := NewTXPool()
	assert.Nil(t, err)
	assert.Equal(t, 0, txpool.Size())

	tx := NewTransaction([]byte("foo"))
	err = tx.Sign(keypair)
	assert.Nil(t, err)

	assert.Nil(t, txpool.Add(tx))
	assert.Equal(t, 1, txpool.Size())
	assert.True(t, txpool.Contains(tx))

	// adding same tx should change nothing
	assert.Nil(t, txpool.Add(tx))
	assert.Equal(t, 1, txpool.Size())
	assert.True(t, txpool.Contains(tx))
}

func TestTransactionPoolFlush(t *testing.T) {
	keypair, err := crypto.GenerateKeyPair()
	assert.Nil(t, err)

	txpool, err := NewTXPool()
	assert.Nil(t, err)
	assert.Equal(t, 0, txpool.Size())

	tx := NewTransaction([]byte("foo"))
	err = tx.Sign(keypair)
	assert.Nil(t, err)

	assert.Nil(t, txpool.Add(tx))
	assert.Equal(t, 1, txpool.Size())
	assert.True(t, txpool.Contains(tx))

	txpool.Flush()

	// adding same tx should change nothing
	assert.Equal(t, 0, txpool.Size())
	assert.False(t, txpool.Contains(tx))
}

func TestTransactionPoolList(t *testing.T) {
	keypair, err := crypto.GenerateKeyPair()
	assert.Nil(t, err)

	txpool, err := NewTXPool()
	assert.Nil(t, err)
	assert.Equal(t, 0, txpool.Size())

	size := 10
	for i := 0; i < size; i++ {
		tx := NewTransaction([]byte(fmt.Sprintf("foo_%d", i)))
		err = tx.Sign(keypair)
		assert.Nil(t, err)

		assert.Nil(t, txpool.Add(tx))
		assert.Equal(t, (i + 1), txpool.Size())
		assert.True(t, txpool.Contains(tx))
	}

	txs := txpool.Transactions()
	for i := 0; i < size; i++ {
		assert.Equal(t, []byte(fmt.Sprintf("foo_%d", i)), txs[i].Data)
	}

}
