package peg

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
)

type Source struct {
	buf []byte
	pos int
}

func NewSource(in io.Reader) (*Source, error) {
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	return &Source{
		buf: buf,
	}, nil
}

// Consume tries to consume text matching the specified regex
// starting at the current position. Returns the consumed text,
// or nil if there was no match.
func (s *Source) Consume(regex *regexp.Regexp) []byte {
	loc := regex.FindIndex(s.buf[s.pos:])
	if loc == nil {
		return nil
	}

	if loc[0] == s.pos {
		s.pos = loc[1]
		return s.buf[loc[0]:loc[1]]
	}

	return nil
}

// Consume literal attempts to consume a literal string.
// Returns the consumed text, or nil if there was no match.
func (s *Source) ConsumeLiteral(valid []byte) []byte {
	if bytes.HasPrefix(s.buf[s.pos:], valid) {
		s.pos += len(valid)
		return valid
	}
	return nil
}
