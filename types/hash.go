package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

var ErrHashFormat = errors.New("given bytes should be with length 32")

const hsize = 32

var (
	ZeroHash      = Hash([hsize]uint8{0})
	zeroHashBytes = ZeroHash.Bytes()
)

type Hash [hsize]uint8

func (h Hash) String() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(h.Bytes()))
}

func (h Hash) Bytes() []byte {
	content := make([]byte, hsize)
	for i := 0; i < hsize; i++ {
		content[i] = h[i]
	}
	return content
}

func (h Hash) IsZeroHash() bool {
	return bytes.Equal(zeroHashBytes[:], h.Bytes())
}

func CreateHash(content []byte) Hash {
	digest := sha256.Sum256(content)
	hash, _ := fromBytes(digest[:])
	return hash
}

func fromBytes(b []byte) (Hash, error) {
	if len(b) != hsize {
		return ZeroHash, ErrHashFormat
	}

	var value [hsize]uint8
	for i := 0; i < hsize; i++ {
		value[i] = b[i]
	}
	return Hash(value), nil
}
