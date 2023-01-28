package network

import (
	"errors"
	"strings"

	"github.com/igumus/chainx/crypto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// create network options
type NetworkOption func(*netOptions)

type netOptions struct {
	// debug mode
	debug bool
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
	// cryptographic keypair for network
	keypair *crypto.KeyPair
	// zerolog logger instance
	logger zerolog.Logger
	// network identifier
	id string
}

func createOptions(opts ...NetworkOption) (*netOptions, error) {
	cfg := &netOptions{
		tcpTransport: "",
		udpTransport: "",
		wsoTransport: "",
		name:         "",
		nodes:        nil,
		keypair:      nil,
		id:           "",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.id = cfg.keypair.Address().String()

	if len(cfg.name) == 0 {
		cfg.name = cfg.id
	}

	// create zerolog logger instance
	cfg.logger = log.With().
		Str("id", cfg.id).
		Str("name", cfg.name).
		Str("component", "network").
		Logger()

	return cfg, nil
}

func (n *netOptions) validate() error {
	if len(n.tcpTransport) == 0 {
		return errors.New("tcp transport addr not specified")
	}
	if n.keypair == nil {
		return errors.New("cryptographic keypair not specified")
	}
	return nil
}

func WithKeyPair(kp *crypto.KeyPair) NetworkOption {
	return func(no *netOptions) {
		no.keypair = kp
	}
}

func WithTCPTransport(addr string) NetworkOption {
	return func(no *netOptions) {
		no.tcpTransport = strings.TrimSpace(addr)
	}
}

func WithSeedNode(addr string) NetworkOption {
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

func WithName(n string) NetworkOption {
	return func(no *netOptions) {
		no.name = strings.TrimSpace(n)
	}
}

func WithDebugMode(b bool) NetworkOption {
	return func(no *netOptions) {
		no.debug = b
	}
}
