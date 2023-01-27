package network

import (
	"io"
	"net"
)

type Transport interface {
	io.Closer
	Type() string
	Listen() error
	Addr() string
	Dial(string) (net.Conn, error)
}
