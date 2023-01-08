package hash

import (
	"bytes"
	"encoding/hex"
	"errors"
)


var (
	zero          byte   = byte(0)
	ZeroHash      Hash   = CreateHash([]byte{0})
	zeroHashBytes []byte = ZeroHash.Bytes()

	// Error Definitions
	ErrHashNotVerified      = errors.New("hash not verified")
	ErrMalformedHash        = errors.New("malformed hash")
)

type Hash []byte

func (h Hash) Bytes() []byte {
	return []byte(h)
}

func (h Hash) String() string {
	return hex.EncodeToString([]byte(h))
}

func (h Hash) IsEqual(o Hash) bool {
	return bytes.Equal(h.Bytes(), o.Bytes())
}

func (h Hash) IsZero() bool {
	return h.IsEqual(ZeroHash)
}

func (h Hash) Verify(data []byte) error {
	version, algo, err := decode(h.Bytes())
	if err != nil {
		return err
	}

	other := createHash(version, algo, data)
	if !h.IsEqual(other) {
		return ErrHashNotVerified
	}

	return nil
}
