package peg

import (
	"errors"
	"fmt"
	"io"
	"regexp"
)

type parseStateFn func(*parser) parseStateFn

type parser struct {
	lex     *lexer
	state   parseStateFn
	parts   chan *Lexeme
	lastErr error
}

func NewParser(input io.Reader) (*Language, error) {
	l := lex(input)
	p := &parser{lex: l}
	return p.prepare()
}

func (p *parser) Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	p.lastErr = errors.New(s)
}

func (p *parser) prepare() (*Language, error) {
	p.parts = make(chan *Lexeme)
	in := make(chan *Language, 1)
	err := make(chan error, 1)
	go constructLanguage(p.parts, in, err)

	for p.state = parseLexeme; p.state != nil; {
		p.state = p.state(p)
	}

	close(p.parts)

	if p.lastErr != nil {
		return nil, p.lastErr
	}

	select {
	case lang := <-in:
		return lang, nil
	case err := <-err:
		return nil, err
	}
}

func constructLanguage(parts chan *Lexeme, success chan *Language, failure chan error) {
	var lexemes = make(map[string]*Lexeme)
	first, ok := <-parts
	if !ok {
		failure <- errors.New("Parts channel was empty.")
		return
	}
	lexemes[first.Name] = first
	for part := range parts {
		lexemes[part.Name] = part
	}

	lex, err := resolveDependencies(first, lexemes)
	if err != nil {
		failure <- err
		return
	} else {
		success <- &Language{
			root: lex,
		}
		return
	}
}

func resolveDependencies(lex *Lexeme, env map[string]*Lexeme) (*Lexeme, error) {
	if lex.Lexer == nil {
		(*lex) = *env[lex.Name[1:]]
	}

	for i, dep := range lex.Dependencies {
		var err error
		lex.Dependencies[i], err = resolveDependencies(dep, env)
		if err != nil {
			return nil, err
		}
	}

	return lex, nil
}

func parseLexeme(p *parser) parseStateFn {
	next, ok := <-p.lex.items
	if !ok {
		return nil
	}
	switch next.typ {
	case itemIdentifier:
		return parseRule(next.val)
	case itemWhitespace:
		return parseLexeme
	case itemError:
		p.Errorf("lex error: %s", next.String())
	default:
		panic(next.typ.String())
	}
	return nil
}

func parseRule(name string) parseStateFn {
	return func(p *parser) parseStateFn {
		next, ok := <-p.lex.items
		if !ok {
			p.Errorf("item channel drained unexpectedly in parseRule")
			return nil
		}
		switch next.typ {
		case itemWhitespace:
			return parseRule(name)
		case itemAssignment:
			return parseRuleBody(name, nil)
		}
		return nil
	}
}

func parseRuleBody(name string, parts []*Lexeme) parseStateFn {
	return func(p *parser) parseStateFn {
		next, ok := <-p.lex.items
		if !ok {
			p.Errorf("item channel drained unexpectedly in parseRuleBody")
			return nil
		}
		switch next.typ {
		case itemWhitespace:
			return parseRuleBody(name, parts)
		case itemLiteral:
			return parseRuleBody(name, append(parts, NewLiteralLexer(name, next.val)))
		case itemRegexp:
			return parseRuleBody(name, append(parts, NewRegexpLexer(name, regexp.MustCompile(next.val))))
		case itemIdentifier:
			return parseRuleBody(name, append(parts, NewRuleLexer(next.val)))
		case itemPlus:
			if len(parts) == 0 {
				p.Errorf("expected lexeme definition before '+'")
				return nil
			}
			lex := parts[len(parts)-1]
			parts := parts[:len(parts)-1]
			return parseRuleBody(name, append(parts, NewPlusClosure(lex)))
		case itemClosure:
			if len(parts) == 0 {
				p.Errorf("expected lexeme definition before '*'")
				return nil
			}
			lex := parts[len(parts)-1]
			parts := parts[:len(parts)-1]
			return parseRuleBody(name, append(parts, NewStarClosure(lex)))
		case itemOptional:
			if len(parts) == 0 {
				p.Errorf("expected lexeme definition before '?'")
				return nil
			}
			lex := parts[len(parts)-1]
			parts := parts[:len(parts)-1]
			return parseRuleBody(name, append(parts, NewOptionClosure(lex)))
		case itemAlternate:
			return parseAlternateRHS(name, parts)
		case itemNewline, itemEOF:
			if len(parts) == 0 {
				return nil
			} else if len(parts) == 1 { // Prevent single literals from being stuck in an array.
				p.parts <- parts[0]
			} else {
				p.parts <- NewConcatLexer(name, parts)
			}
			return parseLexeme
		default:
			p.Errorf("unexpected token : %v", next)
			return nil
		}
		return nil
	}
}

func parseAlternateRHS(name string, parts []*Lexeme) parseStateFn {
	return func(p *parser) parseStateFn {
		next, ok := <-p.lex.items
		if !ok {
			p.Errorf("expected lexeme after '/'")
			return nil
		}
		var rhs *Lexeme
		switch next.typ {
		case itemWhitespace:
			return parseAlternateRHS(name, parts)
		case itemLiteral:
			rhs = NewLiteralLexer(name, next.val)
		case itemRegexp:
			rhs = NewRegexpLexer(name, regexp.MustCompile(next.val))
		case itemIdentifier:
			rhs = NewRuleLexer(next.val)
		default:
			p.Errorf("unexpected token : %v", next)
			return nil
		}

		lhs := parts[len(parts)-1]
		parts := parts[:len(parts)-1]

		return parseRuleBody(name, append(parts, NewAlternateLexer(name, lhs, rhs)))
	}
}
