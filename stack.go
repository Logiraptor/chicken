package chickenVM

// Program stack
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
