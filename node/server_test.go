package node

import (
	"bytes"
	"testing"

	"github.com/igumus/chainx/core"
	"github.com/igumus/chainx/types"
	"github.com/stretchr/testify/require"
)

func TestMessageEncoding(t *testing.T) {

	block, err := core.GenesisBlock()
	require.Nil(t, err)
	require.NotNil(t, block)

	buf := new(bytes.Buffer)
	err = core.EncodeBlock(buf, block)
	require.Nil(t, err)

	data := buf.Bytes()

	message := &types.Message{
		Header: types.ChainBlock,
		Data:   data,
	}

	encMsg, err := message.Bytes()
	require.Nil(t, err)

	remote := &types.RemoteMessage{
		From:    "0x00",
		Payload: encMsg,
	}

	dremote, err := remote.Decode()
	require.Nil(t, err)
	require.NotNil(t, dremote)

}
