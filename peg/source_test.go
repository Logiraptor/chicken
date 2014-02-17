package peg

import (
	"regexp"
	"strings"
	"testing"
)

type ConsumeTest struct {
	Body     string
	Regex    string
	Expected string
}

var sourceConsumeTests = []ConsumeTest{
	ConsumeTest{
		"some text to consume",
		".*",
		"some text to consume",
	},
	ConsumeTest{
		"some text to consume",
		"s.?me",
		"some",
	},
	ConsumeTest{
		"123.43",
		"\\d+",
		"123",
	},
	ConsumeTest{
		"123.43",
		"\\d+.?\\d*",
		"123.43",
	},
}

func TestSourceConsume(t *testing.T) {
	for _, ct := range sourceConsumeTests {
		s, err := NewSource(strings.NewReader(ct.Body))
		if err != nil {
			t.Error(err)
		}
		r := regexp.MustCompile(ct.Regex)
		match := s.Consume(r, 0)
		if match == nil || ct.Expected != string(match) {
			t.Errorf("Source failed to consume input: %s re: %s match: %s exp: %s", ct.Body, ct.Regex, match, ct.Expected)
		}
	}
}
