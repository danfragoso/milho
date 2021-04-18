package parser

import "fmt"

type NodeType int

func (n NodeType) String() string {
	return [...]string{"Nil", "Number", "Boolean", "List", "Macro", "Identifier"}[n]
}

const (
	Nil NodeType = iota
	Number
	Boolean
	List
	Macro
	Identifier
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
	if n.notToeval {
		str = fmt.Sprintf("\n%s+- %s:['%s]", tab, n.Type, n.Identifier)
	} else {
		str = fmt.Sprintf("\n%s+- %s:[%s]", tab, n.Type, n.Identifier)
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
