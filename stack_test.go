package chickenVM

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStack(t *testing.T) {
	var s ValueStack
	Convey("Given a zero-valued stack s", t, func() {

		Convey("It should be empty", func() {
			So(s.IsEmpty(), ShouldBeTrue)
		})

		var val *NilValue
		Convey("Given a value v", func() {
			Convey("When that value is pushed onto the stack", func() {
				s.Push(val)
				Convey("The stack should no longer be empty", func() {
					So(s.IsEmpty(), ShouldBeFalse)
				})
			})

			Convey("When a value is popped", func() {
				v := s.Pop()
				Convey("It should equal v", func() {
					So(v, ShouldEqual, val)
				})
			})

		})
	})
}
