package parser

import "fmt"

type NodeType int

func (n NodeType) String() string {
	return [...]string{"Nil", "Number", "Boolean", "List", "Identifier", "String"}[n]
}

const (
	Nil NodeType = iota
	Number
	Boolean
	List
	Identifier
	String
)

type Node struct {
	notToeval bool

	Type   NodeType
	Parent *Node

	Identifier string
	Nodes      []*Node
}

func (n *Node) ShouldBeEvaluated() bool {
	return !n.notToeval
}

func (n *Node) String() string {
	return n.Sprint("", true)
}

func (n *Node) Sprint(tab string, last bool) string {
	var str string

	str += fmt.Sprintf("\n%s+- ", tab)
	if n.notToeval {
		str += "'"
	}

	str += n.Type.String()
	if n.Identifier != "" {
		str += fmt.Sprintf("#[%s]", n.Identifier)
	}

	if last {
		tab += "   "
	} else {
		tab += "|  "
	}

	for idx, cN := range n.Nodes {
		str += cN.Sprint(tab, idx == len(n.Nodes)-1)
	}

	return str
}
