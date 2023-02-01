package crypto

import (
	"encoding/hex"
)

const size = 20

type Address [size]byte

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) String() string {
	return hex.EncodeToString(a.Bytes())
}

func AddressFromBytes(b []byte) Address {
	if len(b) < size {
		panic("given bytes for address creation should be with size 20")
	}

	buf := b[len(b)-size:]
	var data [size]byte
	for i := 0; i < size; i++ {
		data[i] = buf[i]
	}

	return Address(data)
}
