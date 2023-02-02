package core

import (
	"encoding/binary"
	"fmt"
)

type Instruction byte

const (
	InstrPushInt  Instruction = 0x0a
	InstrMultiply Instruction = 0x0d
	InstrSub      Instruction = 0x0c
	InstrAdd      Instruction = 0x0b
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
	data  []byte // vm data
	ip    int    // instruction pointer
	stack *stack // stack ds
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: newStack(1024),
	}
}

func (vm *VM) Run() error {
	for {
		instr := vm.data[vm.ip]

		if err := vm.exec(Instruction(instr)); err != nil {
			return err
		}

		vm.ip++
		if vm.ip > len(vm.data)-1 {
			break
		}
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

func (vm *VM) exec(instr Instruction) error {
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
	}
	return nil
}
