package peg

import (
	"testing"
)

func TestSimpleLanguage(t *testing.T) {
	l := &Language{
		root: NewLiteralLexer("prgm", "source"),
	}

	tree, err := l.ParseString("source")
	if err != nil {
		t.Error(err)
		return
	}

	if tree.Type != "prgm" {
		t.Errorf("Incorrect type parsed: %s", tree.Type)
	}
}
