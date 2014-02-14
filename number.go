package chickenVM

import (
	"fmt"
)

type Number float64

func (n Number) Native() interface{} {
	return float64(n)
}

func (n Number) Plus(v Value) Value {
	return Number(n + v.(Number))
}

func (n Number) Sub(v Value) Value {
	return Number(n - v.(Number))
}

func (n Number) Div(v Value) Value {
	return Number(n / v.(Number))
}

func (n Number) Mult(v Value) Value {
	return Number(n * v.(Number))
}

func (n Number) String() string {
	return fmt.Sprint(float64(n))
}

func (n Number) Bool() bool {
	return n != 0
}

func (n Number) Compare(v Value) int {
	b := v.(Number)
	if n == b {
		return 0
	} else if n > b {
		return 1
	} else {
		return -1
	}
}
