package chickenVM

import (
	"bufio"
	"fmt"
	"io"
)

type VM struct {
	Stdout io.Writer
	stack  ValueStack
	heap   []Value
	code   []int64
	pc     int64
}

func (v *VM) ensureHeap(n int64) {
	if int64(len(v.heap)) > n {
		return
	}
	temp := make([]Value, n+1)
	copy(temp, v.heap)
	v.heap = temp
}

func (v *VM) Run() []Value {
	v.pc = 0
	for {
		rVal := v.Cycle()
		if rVal != nil {
			return rVal
		}
	}
}

func (v *VM) Cycle() []Value {
	op := v.code[v.pc]
	v.pc++
	return v.Execute(op)
}

func (v *VM) Execute(op int64) []Value {
	var iarg int64
	var a, b Value
	switch op >> OPCODE_SHAMT {
	case OP_PUSH:
		iarg = op & ITYPE_ARGMASK
		v.stack.Push(Number(float64(iarg)))
	case OP_ADD:
		a = v.stack.Pop()
		b = v.stack.Pop()
		v.stack.Push(a.Plus(b))
	case OP_SUB:
		a = v.stack.Pop()
		b = v.stack.Pop()
		v.stack.Push(a.Sub(b))
	case OP_MUL:
		a = v.stack.Pop()
		b = v.stack.Pop()
		v.stack.Push(a.Mult(b))
	case OP_DIV:
		a = v.stack.Pop()
		b = v.stack.Pop()
		v.stack.Push(a.Div(b))
	case OP_GT:
		a = v.stack.Pop()
		b = v.stack.Pop()
		x := a.Compare(b)
		if x > 0 {
			v.stack.Push(Number(1))
		} else {
			v.stack.Push(Number(0))
		}
	case OP_POP:
		iarg = op & ITYPE_ARGMASK
		var i int64
		for i = 0; i < iarg; i++ {
			v.stack.Pop()
		}
	case OP_PRINT:
		iarg = op & ITYPE_ARGMASK
		resp := make([]Value, iarg)
		var i int64
		for i = 0; i < iarg; i++ {
			resp[i] = v.stack.Pop()
			fmt.Fprintf(v.Stdout, "%s ", resp[i].String())
		}
		fmt.Fprint(v.Stdout, "\n")

		for i := iarg - 1; i >= 0; i-- {
			v.stack.Push(resp[i])
		}
	case OP_RET:
		iarg = op & ITYPE_ARGMASK
		resp := make([]Value, iarg)
		var i int64
		for i = 0; i < iarg; i++ {
			resp[i] = v.stack.Pop()
		}
		return resp
	case OP_JUMP:
		iarg = op & ITYPE_ARGMASK
		v.pc = iarg
	case OP_BRANCH:
		iarg = op & ITYPE_ARGMASK
		cond := v.stack.Pop()
		if cond.Bool() {
			v.pc = iarg
		}
	case OP_LOAD:
		iarg = op & ITYPE_ARGMASK
		v.stack.Push(v.heap[iarg])
	case OP_STORE:
		iarg = op & ITYPE_ARGMASK
		v.ensureHeap(iarg)
		v.heap[iarg] = v.stack.Pop()
	case OP_DUMP:
		fmt.Fprintln(v.Stdout, "----Stack----")
		for _, val := range v.stack.data {
			fmt.Fprintln(v.Stdout, val.String())
		}
		fmt.Fprintln(v.Stdout, "----Heap-----")
		for _, val := range v.heap {
			fmt.Fprintln(v.Stdout, val.String())
		}
		fmt.Fprintln(v.Stdout, "-------------")
	}

	return nil
}

func (v *VM) LoadReader(in io.Reader) error {
	v.code = []int64{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		op, err := parseOp(line)
		if err != nil {
			return err
		}
		v.code = append(v.code, op)
	}
	return scanner.Err()
}

func (v *VM) Interpret(in io.Reader) ([]Value, error) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		op, err := parseOp(line)
		if err != nil {
			return nil, err
		}
		val := v.Execute(op)
		if val != nil {
			return val, nil
		}
	}
	return nil, scanner.Err()
}
