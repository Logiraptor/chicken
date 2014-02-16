package peg

import (
	"errors"
	"io"
)

type parseStateFn func(*Parser) parseStateFn

type Parser struct {
	lex   *lexer
	state parseStateFn
}

func NewParser(input io.Reader) (*Parser, error) {
	l := lex(input)
	p := &Parser{lex: l}
	err := p.prepare()

	return p, err
}

func (p *Parser) prepare() error {
	for p.state = parsePeg; p.state != nil; {
		p.state = p.state(p)
	}
	return nil
}

func (p *Parser) Parse(in io.Reader) (*ParseTree, error) {
	return nil, errors.New("unimplemented")
}

func parsePeg(p *Parser) parseStateFn {
	return nil
}
