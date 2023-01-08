package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZeroHash(t *testing.T) {
	testcases := []struct {
		name   string
		hash   Hash
		isZero bool
	}{
		{
			name:   "check-zero-hash",
			hash:   ZeroHash,
			isZero: true,
		},
		{
			name:   "check-non-zero-hash",
			hash:   createHash(version_1, Sha2_256, []byte("hello world")),
			isZero: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.isZero, tc.hash.IsZero())
		})
	}

}

func TestCreateHashWithAlg(t *testing.T) {
	testcases := []struct {
		name    string
		version hashVersion
		alg     HashAlgorithm
		data    []byte
	}{
		{
			name:    "empty-hash-with-sha512",
			version: version_1,
			alg:     Sha2_512,
			data:    []byte{0},
		},
		{
			name:    "empty-hash-with-sha256",
			version: version_1,
			alg:     Sha2_256,
			data:    []byte{0},
		},
		{
			name:    "create-hash-sha512",
			version: version_1,
			alg:     Sha2_512,
			data:    []byte("hello world"),
		},
		{
			name:    "create-hash-sha256",
			version: version_1,
			alg:     Sha2_256,
			data:    []byte("hello world"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			hash := CreateHashWith(tc.alg, tc.data)
			require.Equal(t, byte(tc.version), hash[0])
			require.Equal(t, byte(tc.alg), hash[1])
			lendigest := hash[2]
			require.Equal(t, int(lendigest), len(hash[3:]))
            require.Nil(t, hash.Verify(tc.data))
		})
	}
}
