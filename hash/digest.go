package hash

import (
	"bytes"
	"encoding/hex"

	"github.com/rs/zerolog/log"
)

func createHash(v hashVersion, alg HashAlgorithm, data []byte) Hash {
	var hasher hashFunc = nil

	hasher = hasherFactory(alg)
	if hasher == nil {
		if e := log.Logger.Debug(); e.Enabled() {
			log.Debug().Msg("unknown hash algorithm using default sha256")
		}
		hasher = defaultHasher
	}

	digest := hasher(data)
	digestLen := len(digest)

	result := bytes.Join([][]byte{
		{byte(v), byte(alg), byte(digestLen)},
		digest,
	}, []byte{})

	return result
}

func CreateHash(data []byte) Hash {
	return CreateHashWith(Sha2_256, data)
}

func CreateHashWith(alg HashAlgorithm, data []byte) Hash {
	return createHash(version_1, alg, data)
}

func decode(data []byte) (hashVersion, HashAlgorithm, error) {
	version := versionFactory(data[0])
	if version.isUnknown() {
		return version, identity, ErrUnknownHashVersion
	}

	if version == version_1 {
		algorithm := algorithmFactory(data[1])
		if algorithm.isUnknown() {
			return version, identity, ErrUnknownHashAlgorithm
		}
		return version, algorithm, nil
	}
	return version, identity, ErrUnknownHashVersion
}

func FromHexString(s string) (Hash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Error().Str("hex", s).Err(err).Msg("decoding hex failed")
		return ZeroHash, ErrMalformedHash
	}

	_, _, err = decode(b)
	if err != nil {
		return ZeroHash, err
	}
	return Hash(b), nil
}

func FromBytes(b []byte) (Hash, error) {
	_, _, err := decode(b)
	if err != nil {
		return ZeroHash, err
	}
	return Hash(b), nil
}
