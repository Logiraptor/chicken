package peg

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
)

type Source struct {
	buf []byte
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
func (s *Source) Consume(regex *regexp.Regexp, pos int) []byte {
	loc := regex.FindIndex(s.buf[pos:])
	if loc == nil {
		return nil
	}

	if loc[0] == 0 {
		return s.buf[pos+loc[0] : pos+loc[1]]
	}

	return nil
}

// Consume literal attempts to consume a literal string.
// Returns the consumed text, or nil if there was no match.
func (s *Source) ConsumeLiteral(valid []byte, pos int) []byte {
	if bytes.HasPrefix(s.buf[pos:], valid) {
		pos += len(valid)
		return valid
	}
	return nil
}
