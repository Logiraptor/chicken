package chickenVM

import (
	"strconv"
)

// Value represents and chicken-VM value.
type Value interface {
	// Native returns the go-native datatype stored internally.
	Native() interface{}
	// Plus adds this value to the supplied value
	Plus(Value) Value
	// Div divides this value by the supplied value
	Div(Value) Value
	// Mult multiplies this value by the supplied value
	Mult(Value) Value
	// Sub substracts the provided value from this value
	Sub(Value) Value
	// String returns a human-readable string of the value
	String() string
	// Bool returns a sensible truthiness for this value
	Bool() bool
	// Compare returns 0 if equal, >0 if greater, <0 if lesser
	Compare(Value) int
}

func parseVal(line string) (Value, error) {
	switch line[0] {
	case '\'':
		return String(line[1:]), nil
	default:
		f, err := strconv.ParseFloat(line, 64)
		return Number(f), err
	}
}

// NilValue is a default implementation of the Value interface.
type NilValue struct {
}

func (c *NilValue) Native() interface{} {
	return *c
}

func (c *NilValue) Plus(Value) Value {
	panic("operator '+' not supported.")
}

func (c *NilValue) Sub(Value) Value {
	panic("operator '-' not supported.")
}

func (c *NilValue) Div(Value) Value {
	panic("operator '/' not supported.")
}

func (c *NilValue) Mult(Value) Value {
	panic("operator '*' not supported.")
}

func (c *NilValue) String() string {
	return "nil"
}

func (c *NilValue) Bool() bool {
	return false
}

func (n *NilValue) Float64() float64 {
	panic("Cannot convert nil to float")
}

func (n *NilValue) Compare(Value) int {
	panic("cannot compare nil value")
}
