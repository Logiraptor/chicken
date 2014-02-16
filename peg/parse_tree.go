package peg

type ParseTree struct {
	Type     string
	Data     []byte
	Children []*ParseTree
}
