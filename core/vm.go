package core

import (
	"encoding/binary"
	"fmt"
)

type Instruction byte

const (
	InstrPushInt   Instruction = 0x0a
	InstrPushByte  Instruction = 0x0b
	InstrStrCreate Instruction = 0x0c
	InstrStrPack   Instruction = 0x0d
	InstrStore     Instruction = 0x0e
	InstrLoadState Instruction = 0x0f

	InstrMultiply Instruction = 0x10
	InstrSub      Instruction = 0x11
	InstrAdd      Instruction = 0x12
)

type stack struct {
	data []any // stack data
	sp   int   // stack pointer
}

func newStack(size int) *stack {
	return &stack{
		data: make([]any, size),
		sp:   0,
	}
}

func (s *stack) pop() any {
	data := s.data[s.sp-1]
	s.data[s.sp-1] = nil
	s.sp--
	return data
}

func (s *stack) push(b any) {
	s.data[s.sp] = b
	s.sp++
}

type VM struct {
	data          []byte // vm data
	ip            int    // instruction pointer
	stack         *stack // stack ds
	strSize       int    // string length
	contractState *State // current contract state
}

func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		data:          data,
		stack:         newStack(1024),
		contractState: contractState,
		ip:            0,
		strSize:       0,
	}
}

func (vm *VM) Run() (*State, error) {
	state := NewState()
	for {
		instr := vm.data[vm.ip]

		if err := vm.exec(state, Instruction(instr)); err != nil {
			return state, err
		}

		vm.ip++
		if vm.ip > len(vm.data)-1 {
			break
		}
	}
	return state, nil
}

func (vm *VM) exec(state *State, instr Instruction) error {
	switch instr {
	case InstrPushInt:
		vm.stack.push(vm.data[vm.ip-1])
		return nil
	case InstrAdd:
		return vm.add()
	case InstrMultiply:
		return vm.multiply()
	case InstrSub:
		return vm.substract()
	case InstrPushByte:
		vm.stack.push(vm.data[vm.ip-1])
		return nil
	case InstrStrCreate:
		// check size which should be greater equal than 1
		size := int(vm.data[vm.ip-1])
		vm.strSize = size
		return nil
	case InstrStrPack:
		// check size which should be greater equal than 1
		size := vm.strSize
		content := make([]byte, size)
		for i := size - 1; i >= 0; i-- {
			content[i] = vm.stack.pop().(byte)
		}
		vm.stack.push(content)
		vm.strSize = 0
		return nil
	case InstrStore:
		value, err := vm.toInt()
		if err != nil {
			return err
		}
		key := vm.stack.pop().([]byte)
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, value)
		if err := state.Put(key, buf); err != nil {
			return err
		}
		return nil
	case InstrLoadState:
		// check size which should be greater equal than 1
		size := vm.strSize
		content := make([]byte, size)
		for i := size - 1; i >= 0; i-- {
			content[i] = vm.stack.pop().(byte)
		}
		value, err := vm.contractState.Get(content)
		if err != nil {
			return err
		}
		vm.stack.push(value)
		vm.strSize = 0
		return nil
	}
	return nil
}

func (vm *VM) toInt() (uint64, error) {
	a := vm.stack.pop()
	va := uint64(0)
	switch t := a.(type) {
	case uint64:
		return t, nil
	case byte:
		va = uint64(t)
		return va, nil
	case []byte:
		va = binary.LittleEndian.Uint64(t)
		return va, nil
	default:
		return va, fmt.Errorf("unknown variable type: %T", a)
	}
}

func (vm *VM) add() error {
	//fmt.Printf("current stack: %+v\n", vm.stack)
	va, err := vm.toInt()
	if err != nil {
		return err
	}
	vb, err := vm.toInt()
	if err != nil {
		return err
	}
	c := va + vb
	vm.stack.push(c)
	return nil
}

func (vm *VM) multiply() error {
	va, err := vm.toInt()
	if err != nil {
		return err
	}
	vb, err := vm.toInt()
	if err != nil {
		return err
	}
	c := va * vb
	vm.stack.push(c)
	return nil
}

func (vm *VM) substract() error {
	va, err := vm.toInt()
	if err != nil {
		return err
	}
	vb, err := vm.toInt()
	if err != nil {
		return err
	}
	c := va - vb
	vm.stack.push(c)
	return nil
}
