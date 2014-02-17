package peg

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Lexeme struct {
	Name         string
	Dependencies []*Lexeme
	// Lexer returns the parse tree, an error and the number of input bytes consumed.
	Lexer func(*Source, int) (*ParseTree, error, int)
}

func (l *Lexeme) dumpTree(indent string) string {
	s := fmt.Sprintln(indent, l.Name, l.Lexer != nil)
	for _, child := range l.Dependencies {
		s += child.dumpTree(indent + " ")
	}
	return s
}

func (l *Lexeme) String() string {
	return l.dumpTree("")
}

// Language defines lexing and parsing capabilities for a peg defined language.
type Language struct {
	root *Lexeme
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
	tree, err, _ := l.root.Lexer(s, 0)
	return tree, err
}

func NewLiteralLexer(typ, valid string) *Lexeme {
	vbytes := []byte(valid)
	return &Lexeme{
		Name: valid,
		Lexer: func(s *Source, pos int) (*ParseTree, error, int) {
			match := s.ConsumeLiteral(vbytes, pos)
			if match == nil {
				neighborhood := pos
				neighborEnd := pos + 10
				if neighborEnd > len(s.buf) {
					neighborEnd = len(s.buf)
				}

				return nil, errors.New(fmt.Sprintf("expected literal: %q at %q", valid, s.buf[neighborhood:neighborEnd])), 0
			} else {
				return &ParseTree{
					Type: typ,
					Data: vbytes,
				}, nil, len(match)
			}
		},
	}
}

func NewRegexpLexer(typ string, valid *regexp.Regexp) *Lexeme {
	return &Lexeme{
		Name: typ,
		Lexer: func(s *Source, pos int) (*ParseTree, error, int) {
			match := s.Consume(valid, pos)
			if match == nil {
				neighborhood := pos
				neighborEnd := pos + 10
				if neighborEnd > len(s.buf) {
					neighborEnd = len(s.buf)
				}

				return nil, errors.New(fmt.Sprintf("expected regex match: %q at %q", valid.String(), s.buf[neighborhood:neighborEnd])), 0
			} else {
				return &ParseTree{
					Type: typ,
					Data: match,
				}, nil, len(match)
			}
		},
	}
}

func NewRuleLexer(rule string) *Lexeme {
	return &Lexeme{
		Name:  "~" + rule,
		Lexer: nil,
	}
}

func NewConcatLexer(name string, deps []*Lexeme) *Lexeme {
	return &Lexeme{
		Name:         name,
		Dependencies: deps,
		Lexer: func(s *Source, pos int) (*ParseTree, error, int) {
			children := make([]*ParseTree, 0, len(deps))
			offset := 0
			for _, dep := range deps {
				tree, err, l := dep.Lexer(s, pos+offset)
				if err != nil {
					return nil, err, 0
				} else {
					children = append(children, tree)
					offset += l
				}
			}
			return &ParseTree{Type: name, Data: nil, Children: children}, nil, offset
		},
	}
}
