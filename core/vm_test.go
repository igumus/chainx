package core

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVMInstrStore(t *testing.T) {
	state := NewState()
	contract := []byte{
		0x03,
		byte(InstrStrCreate),
		0x66,
		byte(InstrPushByte),
		0x6f,
		byte(InstrPushByte),
		0x6f,
		byte(InstrPushByte),
		byte(InstrStrPack),
		0x01,
		byte(InstrPushInt),
		0x02,
		byte(InstrPushInt),
		byte(InstrSub),
		byte(InstrStore),
	}

	vm := NewVM(contract, state)
	vstate, err := vm.Run()
	require.Nil(t, err)
	require.NotNil(t, vstate)

	bvalue, err := vstate.Get([]byte("foo"))
	require.Nil(t, err)

	value := binary.LittleEndian.Uint64(bvalue)
	require.Equal(t, value, uint64(1))
	require.Equal(t, vm.stack.sp, 0)
}

func TestVMStringCreation(t *testing.T) {
	testcases := []struct {
		name     string
		contract []byte
		result   string
	}{
		{
			name: "create-f",
			contract: []byte{
				0x03,
				byte(InstrStrCreate),
				0x66,
				byte(InstrPushByte),
				0x6f,
				byte(InstrPushByte),
				0x6f,
				byte(InstrPushByte),
				byte(InstrStrPack),
			},
			result: "foo",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			vm := NewVM(tc.contract, NewState())

			_, err := vm.Run()
			require.Nil(t, err)

			result := vm.stack.pop().([]byte)
			require.Nil(t, err)
			require.Equal(t, result, []byte(tc.result))
			require.Equal(t, vm.stack.sp, 0)

		})
	}

}

func TestVMInstrArithmetics(t *testing.T) {
	testcases := []struct {
		name     string
		contract []byte
		result   uint64
	}{
		{
			name: "1+2",
			contract: []byte{
				0x01,
				byte(InstrPushInt),
				0x02,
				byte(InstrPushInt),
				byte(InstrAdd),
			},
			result: 3,
		},
		{
			name: "2-1",
			contract: []byte{
				0x01,
				byte(InstrPushInt),
				0x02,
				byte(InstrPushInt),
				byte(InstrSub),
			},
			result: 1,
		},
		{
			name: "2*3",
			contract: []byte{
				0x03,
				byte(InstrPushInt),
				0x02,
				byte(InstrPushInt),
				byte(InstrMultiply),
			},
			result: 6,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			vm := NewVM(tc.contract, NewState())

			_, err := vm.Run()
			require.Nil(t, err)

			result, err := vm.toInt()
			require.Nil(t, err)
			require.Equal(t, result, tc.result)
			require.Equal(t, vm.stack.sp, 0)

		})
	}
}

func TestVMStack(t *testing.T) {
	stack := newStack(128)
	stack.push(1)
	stack.push(2)

	value := stack.pop()
	require.Equal(t, value, 2)

	value = stack.pop()
	require.Equal(t, value, 1)

}
