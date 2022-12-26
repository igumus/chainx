package core

import (
	"encoding/gob"
	"errors"
	"io"
)

var ErrNotImplemented = errors.New("not implemented yet")

func NewGobEncoder[T Block | Header](w io.Writer) func(*T) error {
    return func(t *T) error {
        return gob.NewEncoder(w).Encode(t)
    }
}

func NewGobDecoder[T Block | Header](r io.Reader) func(*T) error {
    return func(t *T) error {
        return gob.NewDecoder(r).Decode(t)
    }
}

// TODO: (@igumus) search: isInstanceOf example for golang
func NewBinaryEncoder[T Block | Header](w io.Writer) func(*T) error {
    return func(t *T) error {
        /*
            if instanceOf(t, Block) {
                // do block encoding stuff and return
            } 
            if instanceOf(t, Header) {
                // do header encoding stuff and return
            }
            if instanceOf(t, Transaction) {
                // do transaction encoding stuff and return
            }
        */
        return ErrNotImplemented
    }
}
