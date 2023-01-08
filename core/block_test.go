package core

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeader_Encode_Decode(t *testing.T) {
	header := EmptyHeader()
	buff := new(bytes.Buffer)

	err := NewGobEncoder[Header](buff)(header)
	require.Nil(t, err)

	decoded := &Header{}
	err = NewGobDecoder[Header](buff)(decoded)
	require.Nil(t, err)

	require.Equal(t, header, decoded)
}

func TestBlock_Encode_Decode(t *testing.T) {
	header := EmptyHeader()
	testcases := []struct {
		name string
		txs  []*Transaction
		cnt  int
	}{
		{
			name: "empty-transaction-list",
			txs:  make([]*Transaction, 0),
			cnt:  0,
		},
		{
			name: "nil-transaction-list",
			txs:  nil,
			cnt:  0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			block := NewBlock(header, tc.txs)

			buff := new(bytes.Buffer)
			err := NewGobEncoder[Block](buff)(block)
			require.Nil(t, err)

			decoded := &Block{}
			err = NewGobDecoder[Block](buff)(decoded)
			require.Nil(t, err)
			require.Equal(t, header, block.Header)
			require.Equal(t, tc.cnt, len(block.Transactions))
			require.Equal(t, tc.txs, block.Transactions)
		})
	}
}

func TestBlock_Hash(t *testing.T) {
	block := NewBlock(EmptyHeader(), nil)
	hash := block.Hash()
	require.False(t, hash.IsZero())
}
