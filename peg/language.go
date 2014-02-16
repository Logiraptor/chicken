package peg

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Lexeme func(*Source) (*ParseTree, error)

// Language defines lexing and parsing capabilities for a peg defined language.
type Language struct {
	root Lexeme
}

// ParseString is identical to Parse, but operates on string input.
func (l *Language) ParseString(source string) (*ParseTree, error) {
	return l.Parse(strings.NewReader(source))
}

// Parse attemps to turn the input reader into a valid parse tree.
func (l *Language) Parse(source io.Reader) (*ParseTree, error) {
	s, err := NewSource(source)
	if err != nil {
		return nil, err
	}
	return l.root(s)
}

func NewLiteralLexer(typ, valid string) Lexeme {
	vbytes := []byte(valid)
	return func(s *Source) (*ParseTree, error) {
		match := s.ConsumeLiteral(vbytes)
		if match == nil {
			return nil, errors.New(fmt.Sprintf("expected literal: %s", valid))
		} else {
			return &ParseTree{
				Type: typ,
				Data: vbytes,
			}, nil
		}
	}
}

func NewRegexpLexer(typ string, valid *regexp.Regexp) Lexeme {
	return func(s *Source) (*ParseTree, error) {
		match := s.Consume(valid)
		if match == nil {
			return nil, errors.New(fmt.Sprintf("expected match: %s", valid.String()))
		} else {
			return &ParseTree{
				Type: typ,
				Data: match,
			}, nil
		}
	}
}
