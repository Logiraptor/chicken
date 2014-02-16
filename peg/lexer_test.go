package peg

import (
	"strings"
	"testing"
)

type LexTest struct {
	input string
	exp   []item
}

var lexTestTable = []LexTest{
	LexTest{
		"prgm <- 'a'",
		[]item{
			item{typ: itemIdentifier, val: "prgm"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemAssignment, val: "<-"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemLiteral, val: "a"},
			item{typ: itemEOF, val: ""},
		},
	},
	LexTest{
		"prgm <- ~'-?\\d+.?\\d*'",
		[]item{
			item{typ: itemIdentifier, val: "prgm"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemAssignment, val: "<-"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemRegexp, val: "-?\\d+.?\\d*"},
			item{typ: itemEOF, val: ""},
		},
	},
	LexTest{
		"prgm <- a b \n a <- 'c' \n b <- ~'\\d+'",
		[]item{
			item{typ: itemIdentifier, val: "prgm"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemAssignment, val: "<-"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemIdentifier, val: "a"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemIdentifier, val: "b"},
			item{typ: itemWhitespace, val: " \n "},
			item{typ: itemIdentifier, val: "a"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemAssignment, val: "<-"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemLiteral, val: "c"},
			item{typ: itemWhitespace, val: " \n "},
			item{typ: itemIdentifier, val: "b"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemAssignment, val: "<-"},
			item{typ: itemWhitespace, val: " "},
			item{typ: itemRegexp, val: "\\d+"},
			item{typ: itemEOF, val: ""},
		},
	},
}

func TestLexerTable(t *testing.T) {
	for _, tc := range lexTestTable {
		l := lex(strings.NewReader(tc.input))
		for i, it := range tc.exp {
			ot, ok := <-l.items
			if !ok {
				t.Errorf("No more items after: %v", tc.exp[:i])
				t.Errorf("Expected %v", tc.exp[i])
				return
			}
			if ot.val != it.val {
				t.Errorf("incorrect val: %q exp: %q", ot.val, it.val)
				return
			}
			if ot.typ != it.typ {
				t.Errorf("incorrect typ: %q exp: %q", ot.typ, it.typ)
				return
			}
		}

		x, ok := <-l.items
		if ok {
			t.Errorf("There are extra items on the channel: %v", x)
		}
	}
}
