package peg

import (
	"fmt"
)

type ParseTree struct {
	Type     string
	Data     []byte
	Children []*ParseTree
}

func (p *ParseTree) prettyPrint(indent string) string {
	resp := fmt.Sprintln(indent, p.Type)
	resp += fmt.Sprintf("%s %q\n", indent, string(p.Data))
	for _, child := range p.Children {
		resp += child.prettyPrint(indent + " |")
	}
	return resp
}

func (p *ParseTree) String() string {
	return p.prettyPrint("")
}
