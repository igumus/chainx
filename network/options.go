package network

import (
	"errors"
	"strings"
)

// create network options
type NetworkOption func(*netOptions)

type netOptions struct {
	// tcp transport
	tcpTransport string
	// udp transport
	udpTransport string
	// websocket transport
	wsoTransport string
	// human readable network name
	name string
	// bootstrap nodes
	nodes []string
}

func createOptions(opts ...NetworkOption) (*netOptions, error) {
	cfg := &netOptions{
		tcpTransport: "",
		udpTransport: "",
		wsoTransport: "",
		name:         "",
		nodes:        nil,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, cfg.validate()
}

func (n *netOptions) validate() error {
	if len(n.tcpTransport) == 0 {
		return errors.New("tcp transport addr not specified")
	}
	return nil
}

func WithName(n string) NetworkOption {
	return func(no *netOptions) {
		no.name = strings.TrimSpace(n)
	}
}

func WithTCPTransport(addr string) NetworkOption {
	return func(no *netOptions) {
		no.tcpTransport = strings.TrimSpace(addr)
	}
}

func WithBootstrapNode(addr string) NetworkOption {
	return func(no *netOptions) {
		naddr := strings.TrimSpace(addr)
		if len(naddr) > 0 {
			if no.nodes == nil {
				no.nodes = make([]string, 0)
			}
			no.nodes = append(no.nodes, naddr)
		}
	}
}
