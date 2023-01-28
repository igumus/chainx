package core

import (
	"bytes"
	"testing"

	"github.com/igumus/chainx/hash"
	"github.com/stretchr/testify/require"
)

func TestEncodingBlock(t *testing.T) {

	block, err := GenesisBlock()
	require.Nil(t, err)
	require.NotNil(t, block)
	buf := new(bytes.Buffer)
	err = EncodeBlock(buf, block)
	require.Nil(t, err)
	data := buf.Bytes()
	require.Greater(t, len(data), 0)

	dblock := &Block{Header: &Header{
		PrevBlockHash: hash.ZeroHash,
	}}
	err = DecodeBlock(buf, dblock)
	require.Nil(t, err)
	require.Equal(t, block.Transactions, dblock.Transactions)
	require.Equal(t, block.Header.Version, dblock.Header.Version)
	require.Equal(t, block.Header.Height, dblock.Header.Height)
	require.Equal(t, block.Header.Timestamp, dblock.Header.Timestamp)
	require.Equal(t, block.Header.PrevBlockHash, dblock.Header.PrevBlockHash)
	require.Equal(t, block.Header.DataHash, dblock.Header.DataHash)
}
