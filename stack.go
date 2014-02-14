package chickenVM

// Value stack stores a stack of values. its zero value is an empty stack.
// type ValueStack []Value

// func (v *ValueStack) Push(val Value) {
// 	*v = append(*v, val)
// }

// func (v *ValueStack) Pop() Value {
// 	var rVal Value
// 	rVal, *v = (*v)[len(*v)-1], (*v)[:len(*v)-1]
// 	return rVal
// }

// func (v *ValueStack) IsEmpty() bool {
// 	return len(*v) == 0
// }

type ValueStack struct {
	data []Value
	head int
}

func (v *ValueStack) Push(val Value) {
	if len(v.data) <= v.head {
		v.data = append(v.data, val)
		v.head++
	} else {
		v.data[v.head] = val
		v.head++
	}
}

func (v *ValueStack) Pop() Value {
	v.head--
	return v.data[v.head]
}

func (v *ValueStack) IsEmpty() bool {
	return v.head == 0
}
