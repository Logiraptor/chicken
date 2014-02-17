package peg

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type ParseTest struct {
	language string
	input    string
	exp      *ParseTree
}

var parseTestTable = []ParseTest{
	ParseTest{
		"prgm <- 'a'",
		"a",
		&ParseTree{"prgm", []byte("a"), nil},
	},
	ParseTest{
		"prgm <- ~'\\d+'",
		"74538",
		&ParseTree{"prgm", []byte("74538"), nil},
	},
	ParseTest{
		"prgm <- 'a'_'b' \n _ <- ~'\\s+'",
		"a b",
		&ParseTree{
			"prgm",
			nil,
			[]*ParseTree{
				&ParseTree{"prgm", []byte("a"), nil},
				&ParseTree{"_", []byte(" "), nil},
				&ParseTree{"prgm", []byte("b"), nil},
			},
		},
	},
	ParseTest{
		"prgm <- name '=' number \n name <- ~'[a-zA-Z]+' \n number <- ~'\\d+'",
		"variableName=432",
		&ParseTree{
			"prgm",
			nil,
			[]*ParseTree{
				&ParseTree{"name", []byte("variableName"), nil},
				&ParseTree{"prgm", []byte("="), nil},
				&ParseTree{"number", []byte("432"), nil},
			},
		},
	},
	ParseTest{
		"prgm <- a+\na <- 'a'",
		"aaa",
		&ParseTree{
			"a+",
			nil,
			[]*ParseTree{
				&ParseTree{"a", []byte("a"), nil},
				&ParseTree{"a", []byte("a"), nil},
				&ParseTree{"a", []byte("a"), nil},
			},
		},
	},
	ParseTest{
		"prgm <- a+\na <- 'a' _?\n_ <- ~'\\s'",
		"aa a",
		&ParseTree{
			"a+",
			nil,
			[]*ParseTree{
				&ParseTree{"a", []byte("a"), nil},
				&ParseTree{"a", nil, []*ParseTree{
					&ParseTree{"a", []byte("a"), nil},
					&ParseTree{"_", []byte(" "), nil},
				}},
				&ParseTree{"a", []byte("a"), nil},
			},
		},
	},
	ParseTest{
		"prgm <- a*\na <- 'a' _?\n_ <- ~'\\s+'",
		"aa \ta",
		&ParseTree{
			"a*",
			nil,
			[]*ParseTree{
				&ParseTree{"a", []byte("a"), nil},
				&ParseTree{"a", nil, []*ParseTree{
					&ParseTree{"a", []byte("a"), nil},
					&ParseTree{"_", []byte(" \t"), nil},
				}},
				&ParseTree{"a", []byte("a"), nil},
			},
		},
	},
}

func TestTemp(t *testing.T) {
	// l := lex(strings.NewReader("prgm <- name '=' number \n name <- ~'[a-zA-Z]+' \n number <- ~'\\d+'"))
	// for item := range l.items {
	// 	fmt.Println(item)
	// }
	// t.Errorf("dummy")
}

func TestParseTable(t *testing.T) {
	for _, tc := range parseTestTable {
		parser, err := NewParser(strings.NewReader(tc.language))
		if err != nil {
			t.Error(err)
			return
		}

		tree, err := parser.Parse(strings.NewReader(tc.input))
		if err != nil {
			t.Error(err)
			return
		}

		if err := treeCompare(tree, tc.exp); err != nil {
			fmt.Println("Got:")
			dumpTree(tree, "")
			fmt.Println("Expected:")
			dumpTree(tc.exp, "")
			t.Error(err)
			return
		}
	}
}

func treeCompare(a, b *ParseTree) error {
	if a == b {
		return nil
	} else if a == nil || b == nil {
		return errors.New(fmt.Sprintf("a or b is nil %v %v", a, b))
	}
	if a.Type != b.Type {
		return errors.New(fmt.Sprintf("tree type mismatch: %q exp: %q", a.Type, b.Type))
	}
	if !bytes.Equal(a.Data, b.Data) {
		return errors.New(fmt.Sprintf("tree data mismatch: %q exp: %q", string(a.Data), string(b.Data)))
	}

	if len(a.Children) != len(b.Children) {
		return errors.New(fmt.Sprintf("trees have different number of chidren: %d exp: %d", len(a.Children), len(b.Children)))
	}

	for i, child := range a.Children {
		if err := treeCompare(child, b.Children[i]); err != nil {
			return err
		}
	}

	return nil
}

func dumpTree(tree *ParseTree, indent string) {
	if tree == nil {
		fmt.Println(indent, nil)
	} else {
		fmt.Println(indent, tree.Type)
		fmt.Printf("%s  %q\n", indent, string(tree.Data))
		for _, child := range tree.Children {
			dumpTree(child, indent+" ")
		}
	}
}
