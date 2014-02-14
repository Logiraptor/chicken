package chickenVM

import (
	"strings"
)

type String string

func (s String) Native() interface{} {
	return string(s)
}

func (s String) Plus(v Value) Value {
	return String(string(s) + v.String())
}

func (s String) Sub(v Value) Value {
	panic("operator '-' is undefined for type string")
}

func (s String) Div(Value) Value {
	panic("operator '/' is undefined for type string")
}

func (s String) Mult(v Value) Value {
	return String(strings.Repeat(string(s), int(v.(Number))))
}

func (s String) String() string {
	return string(s)
}

func (s String) Bool() bool {
	return len(s) > 0
}

func (s String) Compare(Value) int {
	panic("Cannot compare string value")
}
