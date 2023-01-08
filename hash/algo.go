package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
)

type HashAlgorithm byte

func (a HashAlgorithm) isUnknown() bool {
	return byte(a) == byte(identity)
}

const (
	identity HashAlgorithm = 0
	Sha1     HashAlgorithm = 1
	Sha2_256 HashAlgorithm = 2
	Sha2_512 HashAlgorithm = 3
)

var ErrUnknownHashAlgorithm = errors.New("unknown hash algorithm")

func algorithmFactory(v byte) HashAlgorithm {
	switch v {
	case byte(Sha1):
		return Sha1
	case byte(Sha2_256):
		return Sha2_256
	case byte(Sha2_512):
		return Sha2_512
	default:
		return identity
	}
}

type hashFunc func([]byte) []byte

var defaultHasher hashFunc = hasherFactory(Sha2_256)

func hasherFactory(alg HashAlgorithm) hashFunc {
	switch alg {
	case Sha1:
		return func(b []byte) []byte {
			data := sha1.Sum(b)
			return data[:]
		}
	case Sha2_512:
		return func(b []byte) []byte {
			data := sha512.Sum512(b)
			return data[:]
		}
	case Sha2_256:
		return func(b []byte) []byte {
			data := sha256.Sum256(b)
			return data[:]
		}
	default:
		return nil
	}
}
