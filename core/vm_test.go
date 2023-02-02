package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVMStack(t *testing.T) {
	stack := newStack(128)
	stack.push(1)
	stack.push(2)

	value := stack.pop()
	require.Equal(t, value, 2)

	value = stack.pop()
	require.Equal(t, value, 1)

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
			vm := NewVM(tc.contract)

			err := vm.Run()
			require.Nil(t, err)

			result, err := vm.toInt()
			require.Nil(t, err)
			require.Equal(t, result, tc.result)
			require.Equal(t, vm.stack.sp, 0)

		})
	}
}
