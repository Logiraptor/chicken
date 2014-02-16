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
		&ParseTree{
			"prgm",
			[]byte("a"),
			nil,
		},
	},
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
			t.Error(err)
			return
		}
	}
}

func treeCompare(a, b *ParseTree) error {
	if a.Type != b.Type {
		return errors.New(fmt.Sprintf("tree type mismatch: %s exp: %s", a.Type, b.Type))
	}
	if !bytes.Equal(a.Data, b.Data) {
		return errors.New(fmt.Sprintf("tree data mismatch: %s exp: %s", string(a.Data), string(b.Data)))
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
