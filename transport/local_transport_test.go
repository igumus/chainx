package transport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	tra := createLocalTransport("A")
	trb := createLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	require.Equal(t, trb, tra.peers[trb.Addr()])
	require.Equal(t, tra, trb.peers[tra.Addr()])
}

func TestSendMessage(t *testing.T) {
	tra := createLocalTransport("A")
	trb := createLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	payload := []byte("hello world")

	err := tra.SendMessage(trb.addr, payload)
	require.Nil(t, err)

	receivedRPC := <-trb.consumeCh

	require.Equal(t, payload, receivedRPC.Payload)

}

