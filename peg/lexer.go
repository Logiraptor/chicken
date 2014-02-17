package peg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType
	pos int
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	}
	return fmt.Sprintf("%s:%q", i.typ, i.val)
}

type itemType int

const (
	itemUNKNOWN itemType = iota
	itemError   itemType = iota
	itemAssignment
	itemQuote
	itemLiteral
	itemWhitespace
	itemNewline
	itemIdentifier
	itemRegexp
	itemClosure
	itemPlus
	itemAlternate
	itemOptional
	itemEOF
)

func (i itemType) String() string {
	switch i {
	case itemError:
		return "itemError"
	case itemAssignment:
		return "itemAssignment"
	case itemQuote:
		return "itemQuote"
	case itemLiteral:
		return "itemLiteral"
	case itemWhitespace:
		return "itemWhitespace"
	case itemNewline:
		return "itemNewline"
	case itemIdentifier:
		return "itemIdentifier"
	case itemRegexp:
		return "itemRegexp"
	case itemEOF:
		return "itemEOF"
	case itemClosure:
		return "itemClosure"
	case itemPlus:
		return "itemPlus"
	case itemAlternate:
		return "itemAlternate"
	case itemOptional:
		return "itemOptional"
	}
	return "UNKNOWN"
}

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input  *bufio.Reader
	buffer bytes.Buffer
	state  stateFn
	pos    int
	start  int
	items  chan item
}

func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

func lex(input io.Reader) *lexer {
	l := &lexer{
		input: bufio.NewReader(input),
		items: make(chan item, 1),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexPeg; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) next() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	l.pos += w
	l.buffer.WriteRune(r)
	return r
}

func (l *lexer) peek() rune {
	lead, err := l.input.Peek(1)
	if err == io.EOF {
		return eof
	} else if err != nil {
		l.errorf("peek: %s", err.Error())
		return 0
	}

	p, err := l.input.Peek(runeLen(lead[0]))
	if err == io.EOF {
		return eof
	} else if err != nil {
		l.errorf("peek: %s", err.Error())
		return 0
	}
	r, _ := utf8.DecodeRune(p)
	return r
}

func runeLen(lead byte) int {
	if lead < 0xC0 {
		return 1
	} else if lead < 0xE0 {
		return 2
	} else if lead < 0xF0 {
		return 3
	} else {
		return 4
	}
}

func (l *lexer) emit(t itemType) {
	l.emitInner(t, 0, 0)
}

// emitInner trims left characters from the left,
// right characters from the right side of the token
// and emits that.
func (l *lexer) emitInner(t itemType, left, right int) {
	token := l.buffer.String()
	l.items <- item{t, l.start + left, token[left : len(token)-right]}
	l.start = l.pos
	l.buffer.Truncate(0)
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.peek()) >= 0 {
		l.next()
		return true
	}
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.peek()) >= 0 {
		l.next()
	}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) hasPrefix(prefix string) bool {
	p, err := l.input.Peek(len(prefix))
	if err == io.EOF {
		return false
	} else if err != nil {
		l.errorf("hasPrefix: %s", err.Error())
		return false
	}
	return string(p) == prefix
}

// Accept next count runes. Normally called after hasPrefix().
func (l *lexer) nextRuneCount(count int) {
	for i := 0; i < count; i++ {
		l.next()
	}
}

func isIdentRune(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func lexPeg(l *lexer) stateFn {
	switch r := l.peek(); {
	case isIdentRune(r):
		return lexIdentifier
	case unicode.IsSpace(r) && r != '\n':
		return lexWhitespace
	case r == '\n':
		return lexNewline
	case r == '<':
		return lexAssignment
	case r == '\'':
		return lexLiteral
	case r == '~':
		return lexRegex
	case r == '*':
		return lexClosure
	case r == '+':
		return lexPlus
	case r == '/':
		return lexAlternate
	case r == '?':
		return lexOption
	case r == eof:
		l.emit(itemEOF)
		return nil
	}

	return nil
}

func lexPlus(l *lexer) stateFn {
	l.next()
	l.emit(itemPlus)
	return lexPeg
}

func lexAlternate(l *lexer) stateFn {
	l.next()
	l.emit(itemAlternate)
	return lexPeg
}

func lexOption(l *lexer) stateFn {
	l.next()
	l.emit(itemOptional)
	return lexPeg
}

func lexClosure(l *lexer) stateFn {
	l.next()
	l.emit(itemClosure)
	return lexPeg
}

func lexIdentifier(l *lexer) stateFn {
	for isIdentRune(l.peek()) {
		l.next()
	}
	l.emit(itemIdentifier)
	return lexPeg
}

func lexWhitespace(l *lexer) stateFn {
	for {
		r := l.peek()
		if unicode.IsSpace(r) && r != '\n' {
			l.next()
		} else {
			break
		}
	}
	l.emit(itemWhitespace)
	return lexPeg
}

func lexNewline(l *lexer) stateFn {
	l.next()
	l.emit(itemNewline)
	return lexPeg
}

func lexAssignment(l *lexer) stateFn {
	l.next()
	if l.next() != '-' {
		l.errorf("expected <-")
		return nil
	} else {
		l.emit(itemAssignment)
	}
	return lexPeg
}

func lexLiteral(l *lexer) stateFn {
	l.next() // consume '

	for {
		r := l.next()
		if r == '\\' && l.peek() == '\'' {
			l.next()
		} else if r == '\'' {
			l.emitInner(itemLiteral, 1, 1)
			return lexPeg
		} else if r == eof {
			l.errorf("eof while parsing literal")
			return nil
		}
	}
}

func lexRegex(l *lexer) stateFn {
	l.next() // consume ~

	if l.peek() != '\'' {
		l.errorf("Expected \"'\" after ~")
		return nil
	} else {
		l.next() // consume '
	}

	for {
		r := l.next()
		if r == '\\' && l.peek() == '\'' {
			l.next()
		} else if r == '\'' {
			l.emitInner(itemRegexp, 2, 1)
			return lexPeg
		} else if r == eof {
			l.errorf("eof while parsing regexp")
			return nil
		}
	}
}
