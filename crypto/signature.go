package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"math/big"
)

type Signature struct {
	S      *big.Int
	R      *big.Int
	PubKey PublicKey
}

func (s *Signature) Verify(data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), s.PubKey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.Verify(key, data, s.R, s.S)
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
