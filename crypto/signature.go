package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"math/big"
)

var (
	ErrNoSignature      = errors.New("signature not found")
	ErrInvalidSignature = errors.New("signature is invalid")
)

type Signature struct {
	S      *big.Int
	R      *big.Int
	PubKey []byte
}

func (s *Signature) Verify(data []byte) error {
	if s == nil {
		return ErrNoSignature
	}
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), s.PubKey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	if !ecdsa.Verify(key, data, s.R, s.S) {
		return ErrInvalidSignature
	}
	return nil
}

func (s *Signature) Bytes() []byte {
	data := bytes.Join([][]byte{
		s.S.Bytes(),
		s.R.Bytes(),
		s.PubKey,
	}, []byte{})

	return data
}

func (s *Signature) String() string {
	return hex.EncodeToString(s.Bytes())
}
