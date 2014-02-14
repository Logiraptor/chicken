package chickenVM

import (
	"encoding/binary"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"os"
	"strings"
	"testing"
)

var testArithSource = `push 6
push 2
push 4
push 1
push 3
sub
add
div
mul
ret 1
`

func TestSimpleAdd(t *testing.T) {
	var vm = &VM{
		Stdout: os.Stdout,
	}
	Convey("Given an initialized vm.", t, func() {
		Convey("(((4+(3-1))/2)*6)==18", func() {
			err := vm.LoadReader(strings.NewReader(testArithSource))
			if err != nil {
				t.Error(err)
				return
			}
			print(len(vm.code))
			result := vm.Run()
			So(result[0], ShouldEqual, Number(18))
		})
	})
}

func BenchmarkValueFloat(b *testing.B) {
	c := Number(2.345)
	var f Value = Number(1)
	for i := 0; i < b.N; i++ {
		f = f.Mult(c)
	}
}

func BenchmarkByteArrayFloat(b *testing.B) {
	var f uint64 = math.Float64bits(1)
	var mem = make([]byte, 8)
	binary.BigEndian.PutUint64(mem, f)
	for i := 0; i < b.N; i++ {
		x := math.Float64frombits(binary.BigEndian.Uint64(mem))
		x *= 2.345
		binary.BigEndian.PutUint64(mem, math.Float64bits(x))
	}
}

func BenchmarkIntBitsFloat(b *testing.B) {
	var f uint64 = math.Float64bits(1)
	for i := 0; i < b.N; i++ {
		x := math.Float64frombits(f)
		x *= 2.345
		f = math.Float64bits(x)
	}
}

func BenchmarkNativeFloat(b *testing.B) {
	var f float64 = 1
	for i := 0; i < b.N; i++ {
		f *= 2.345
	}
}
