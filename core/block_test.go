package core

import (
	"bytes"
	"testing"

	"github.com/igumus/chainx/crypto"
	"github.com/igumus/chainx/hash"
	"github.com/stretchr/testify/require"
)

func createGenesisHeader(t *testing.T) *Header {
	tx := createSignedTransaction(t, []byte("genesis"))
	txHash, err := calculateTransactionHash([]*Transaction{tx})
	require.Nil(t, err)

	return &Header{
		Version:       1,
		Height:        0,
		Timestamp:     0,
		PrevBlockHash: hash.ZeroHash,
		DataHash:      txHash,
	}
}

func TestBlockCreate(t *testing.T) {
	prevHeader := createGenesisHeader(t)

	txs := []*Transaction{
		createSignedTransaction(t, []byte("foo")),
	}
	block, err := NewBlock(prevHeader, txs)
	require.Nil(t, err)

	currHeader := block.Header

	require.Equal(t, prevHeader.Version, currHeader.Version)
	require.Equal(t, uint32(0), prevHeader.Height)
	require.NotEqual(t, prevHeader.Height, currHeader.Height)
	require.Equal(t, uint32(1), currHeader.Height-prevHeader.Height)

}

func TestBlockSign(t *testing.T) {
	kp, err := crypto.GenerateKeyPair()
	require.Nil(t, err)
	require.NotNil(t, kp)

	prevHeader := createGenesisHeader(t)

	txs := []*Transaction{
		createSignedTransaction(t, []byte("foo")),
	}

	block, err := NewBlock(prevHeader, txs)
	require.Nil(t, err)

	require.Nil(t, block.Signature)
	err = block.Sign(kp)
	require.Nil(t, err)
	require.NotNil(t, block.Signature)
	require.Nil(t, block.Verify())
}

func TestBlockSignatureVerify(t *testing.T) {
	kp, err := crypto.GenerateKeyPair()
	require.Nil(t, err)
	require.NotNil(t, kp)

	prevHeader := createGenesisHeader(t)

	txs := []*Transaction{
		createSignedTransaction(t, []byte("foo")),
	}

	block, err := NewBlock(prevHeader, txs)
	require.Nil(t, err)

	require.Nil(t, block.Signature)
	require.Equal(t, crypto.ErrNoSignature, block.Verify())
	err = block.Sign(kp)
	require.Nil(t, err)
	require.Nil(t, block.Verify())

	x := block.Header.Bytes()
	block.Header.Height = uint32(100)
	y := block.Header.Bytes()

	require.False(t, bytes.Equal(x, y))
	require.NotEqual(t, x, y)
	require.Equal(t, crypto.ErrInvalidSignature, block.Verify())
}
