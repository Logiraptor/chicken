package chickenVM

import (
	"errors"
	"strconv"
	"strings"
)

var NoSuchOpError = errors.New("No such operation")

const (
	OPCODE_SHAMT  = 59
	ITYPE_ARGMASK = 0x7ffffffffffffff
)

const (
	OP_PUSH int64 = iota
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_GT
	OP_POP
	OP_PRINT
	OP_RET
	OP_JUMP
	OP_BRANCH
	OP_LOAD
	OP_STORE
	OP_DUMP
)

func ITypeOperation(code int64, arg int64) int64 {
	return OpFromOpcode(code) | arg
}

func OpFromOpcode(code int64) int64 {
	return code << OPCODE_SHAMT
}

// TODO: add dup
func parseOp(line string) (int64, error) {
	parts := strings.SplitN(line, " ", 2)
	switch parts[0] {
	case "push":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_PUSH, i), err
	case "add":
		return OpFromOpcode(OP_ADD), nil
	case "sub":
		return OpFromOpcode(OP_SUB), nil
	case "mul":
		return OpFromOpcode(OP_MUL), nil
	case "div":
		return OpFromOpcode(OP_DIV), nil
	case "gt":
		return OpFromOpcode(OP_GT), nil
	case "pop":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_POP, i), err
	case "print":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_PRINT, i), err
	case "ret":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_RET, i), err
	case "jump":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_JUMP, i), err
	case "branch":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_BRANCH, i), err
	case "load":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_LOAD, i), err
	case "store":
		i, err := strconv.ParseInt(parts[1], 0, 32)
		return ITypeOperation(OP_STORE, i), err
	case "dump":
		return OpFromOpcode(OP_DUMP), nil
	}
	return 0, NoSuchOpError
}

func GreaterThan(vm *VM) []Value {
	a := vm.stack.Pop()
	b := vm.stack.Pop()
	x := a.Compare(b)
	if x > 0 {
		vm.stack.Push(Number(1))
	} else {
		vm.stack.Push(Number(0))
	}
	return nil
}

// // Load reads the heap value at index and pushes it onto the stack.
// func Load(index int) int64 {
// 	return func(vm *VM) []Value {
// 		vm.stack.Push(vm.heap[index])
// 		return nil
// 	}
// }

// // Store pops a value from the stack and stores it in the heap at index.
// func Store(index int) int64 {
// 	return func(vm *VM) []Value {
// 		vm.ensureHeap(index)
// 		vm.heap[index] = vm.stack.Pop()
// 		return nil
// 	}
// }

// // Branch pops one value from the stack. if it is true, the pc is changed to line
// func Branch(line int) int64 {
// 	return func(vm *VM) []Value {
// 		cond := vm.stack.Pop()
// 		if cond.Bool() {
// 			vm.pc = line
// 		}
// 		return nil
// 	}
// }

// func Jump(line int) int64 {
// 	return func(vm *VM) []Value {
// 		vm.pc = line
// 		return nil
// 	}
// }

// func Dump(vm *VM) []Value {
// 	fmt.Fprintln(vm.Stdout, "----Stack----")
// 	for _, val := range vm.stack {
// 		fmt.Fprintln(vm.Stdout, val.String())
// 	}
// 	fmt.Fprintln(vm.Stdout, "----Heap-----")
// 	for _, val := range vm.heap {
// 		fmt.Fprintln(vm.Stdout, val.String())
// 	}
// 	fmt.Fprintln(vm.Stdout, "-------------")
// 	return nil
// }

// // Print pops n values from the stack then pushes them all back
// func Print(n int) int64 {
// 	return func(vm *VM) []Value {
// 		resp := make([]Value, n)
// 		for i := 0; i < n; i++ {
// 			resp[i] = vm.stack.Pop()
// 			fmt.Fprintf(vm.Stdout, "%s ", resp[i].String())
// 		}
// 		fmt.Fprint(vm.Stdout, "\n")

// 		for i := n - 1; i >= 0; i-- {
// 			vm.stack.Push(resp[i])
// 		}
// 		return nil
// 	}
// }

// func Return(n int) int64 {
// 	return func(vm *VM) []Value {
// 		resp := make([]Value, n)
// 		for i := 0; i < n; i++ {
// 			resp[i] = vm.stack.Pop()
// 		}
// 		return resp
// 	}
// }

// func Push(v ...Value) int64 {
// 	return func(vm *VM) []Value {
// 		for _, x := range v {
// 			vm.stack.Push(x)
// 		}
// 		return nil
// 	}
// }

// func Pop(n int) int64 {
// 	return func(vm *VM) []Value {
// 		for i := 0; i < n; i++ {
// 			vm.stack.Pop()
// 		}
// 		return nil
// 	}
// }

// func Add(vm *VM) []Value {
// 	a := vm.stack.Pop()
// 	b := vm.stack.Pop()
// 	vm.stack.Push(a.Plus(b))
// 	return nil
// }

// func Sub(vm *VM) []Value {
// 	a := vm.stack.Pop()
// 	b := vm.stack.Pop()
// 	vm.stack.Push(a.Sub(b))
// 	return nil
// }

// func Mul(vm *VM) []Value {
// 	a := vm.stack.Pop()
// 	b := vm.stack.Pop()
// 	vm.stack.Push(a.Mult(b))
// 	return nil
// }

// func Div(vm *VM) []Value {
// 	a := vm.stack.Pop()
// 	b := vm.stack.Pop()
// 	vm.stack.Push(a.Div(b))
// 	return nil
// }
