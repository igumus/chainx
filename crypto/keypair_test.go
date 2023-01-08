package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateKeyPair(t *testing.T) {

	privateKey, err := GenerateKeyPair()
	require.Nil(t, err)
	require.NotNil(t, privateKey)

	require.NotNil(t, privateKey.publicKey())
	require.Greater(t, len(privateKey.publicKey()), 0)

}
