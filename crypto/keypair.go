package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"

	"github.com/igumus/chainx/hash"
)

func GenerateKeyPair() (*KeyPair, error) {
	return GenerateKeyPairFromReader(rand.Reader)
}

func GenerateKeyPairFromReader(r io.Reader) (*KeyPair, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		privKey: key,
	}, nil
}

type KeyPair struct {
	privKey *ecdsa.PrivateKey
}

func (p *KeyPair) publicKey() []byte {
	pkey := p.privKey.PublicKey
	return elliptic.MarshalCompressed(pkey, pkey.X, pkey.Y)
}

func (p *KeyPair) Address() Address {
	h := hash.CreateHash(p.publicKey())
	return AddressFromBytes(h)
}

func (p *KeyPair) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, p.privKey, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		R:      r,
		S:      s,
		PubKey: p.publicKey(),
	}, nil

}
