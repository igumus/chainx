package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"io"
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

type PublicKey []byte

func (pbk PublicKey) String() string {
	return hex.EncodeToString(pbk)
}

type KeyPair struct {
	privKey *ecdsa.PrivateKey
}

func (p *KeyPair) publicKey() PublicKey {
	return elliptic.MarshalCompressed(p.privKey.PublicKey, p.privKey.PublicKey.X, p.privKey.PublicKey.Y)
}

func (p *KeyPair) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, p.privKey, data)
    if err != nil {
        return nil, err
    }
    return &Signature{
        R: r,
        S: s,
        PubKey: p.publicKey(),
    }, nil

}
