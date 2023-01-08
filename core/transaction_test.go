package core

import (
	"bytes"
	"testing"

	"github.com/igumus/chainx/crypto"
	"github.com/stretchr/testify/require"
)

func createSignedTransaction(t *testing.T, data []byte) *Transaction {
	kp, err := crypto.GenerateKeyPair()
	require.Nil(t, err)
	require.NotNil(t, kp)

	tx := &Transaction{
		Data: data,
	}

	err = tx.Sign(kp)
	require.Nil(t, err)
    return tx
}

func TestTransactionSign(t *testing.T) {
	data := []byte("hello world")
    createSignedTransaction(t, data)
}

func TestTransactionVerifyWithSignature(t *testing.T) {
	data := []byte("hello world")
    tx := createSignedTransaction(t, data)
    err := tx.Verify()
	require.Nil(t, err)
}

func TestTransactionVerifyWithNoSignature(t *testing.T) {
	data := []byte("hello world")
	tx := &Transaction{
		Data: data,
	}

    err := tx.Verify()
	require.Equal(t, err, crypto.ErrNoSignature)
}

func TestTransactionVerifyWithTamperedData(t *testing.T) {
	data := []byte("hello world")
    tx := createSignedTransaction(t, data)
	tx.Data = append(tx.Data, byte(1))

    err := tx.Verify()
	require.Equal(t, crypto.ErrInvalidSignature, err)
}

func TestTransactionEncode(t *testing.T) {
	data := []byte("hello world")
    tx := createSignedTransaction(t, data)

    require.Nil(t, tx.Verify())

    buf := new(bytes.Buffer)
	err := EncodeTransaction(buf, tx)
    require.Nil(t, err)
    require.Greater(t, len(buf.Bytes()), 0)
    
    decoded := new(Transaction)
    err = DecodeTransaction(buf, decoded)
    require.Nil(t, decoded.Verify())
}
