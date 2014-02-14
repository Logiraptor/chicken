package chickenVM

import (
	"encoding/binary"
	"math"
)

// Dynamic type binding at runtime
type Pointer struct {
	Type    int32
	Address int32
}

const (
	FloatType int32 = iota
	StringType
	// ObjectType
	// ArrayType
)

func (p *Pointer) DecodeNumber(mem []byte) float64 {
	if p.Type != FloatType {
		panic("pointer is not a float")
	}
	return math.Float64frombits(binary.BigEndian.Uint64(mem[p.Address:]))
}

func (p *Pointer) DecodeString(mem []byte) string {
	if p.Type != StringType {
		panic("pointer is not a string")
	}
	return string(mem[p.Address+1 : p.Address+1+int32(mem[p.Address])])
}
