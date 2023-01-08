package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyPairGenerate(t *testing.T) {
	size := 20
	privateKey, err := GenerateKeyPair()
	require.Nil(t, err)
	require.NotNil(t, privateKey)

	pubKey := privateKey.publicKey()
	require.NotNil(t, pubKey)
	require.Greater(t, len(pubKey), 0)

	addr := privateKey.Address()
	require.NotNil(t, addr)
	require.Equal(t, size, len(addr))
}

func TestKeyPairSign(t *testing.T) {
	testcases := []struct {
		name       string
		data       []byte
		verify     []byte
		shouldFail bool
	}{
		{
			name:       "valid-signature",
			data:       []byte("hello world"),
			verify:     []byte("hello world"),
			shouldFail: false,
		},
		{
			name:       "tempered-data-verification-fail",
			data:       []byte("hello world"),
			verify:     []byte("hello world."),
			shouldFail: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := GenerateKeyPair()
			require.Nil(t, err)
			signature, err := key.Sign(tc.data)
			require.Nil(t, err)
			require.Equal(t, tc.shouldFail, !signature.Verify(tc.verify))
		})
	}
}
