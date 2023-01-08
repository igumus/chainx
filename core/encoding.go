package core

import (
	"encoding/gob"
	"io"
)

func generateGobDecoder[T Block | Header | Transaction]() func(io.Reader, *T) error {
	return func(r io.Reader, t *T) error {
		return gob.NewDecoder(r).Decode(t)
	}
}

func generateGobEncoder[T Block | Header | Transaction]() func(io.Writer, *T) error {
	return func(w io.Writer, t *T) error {
		return gob.NewEncoder(w).Encode(t)
	}
}

var (
	EncodeTransaction = generateGobEncoder[Transaction]()
	DecodeTransaction = generateGobDecoder[Transaction]()

	EncodeBlock = generateGobEncoder[Block]()
	DecodeBlock = generateGobDecoder[Block]()

	EncodeHeader = generateGobEncoder[Header]()
	DecodeHeader = generateGobDecoder[Header]()
)
