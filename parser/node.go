package parser

import "fmt"

type NodeType int

func (n NodeType) String() string {
	return [...]string{"Nil", "Number", "Function"}[n]
}

const (
	Nil NodeType = iota
	Number
	Function
)

type Node struct {
	Type   NodeType
	Parent *Node

	Identifier string
	Nodes      []*Node
}

func (n *Node) String() string {
	return n.Sprint("", true)
}

func (n *Node) Sprint(tab string, last bool) string {
	str := fmt.Sprintf("\n%s+- %s:[%s]", tab, n.Type, n.Identifier)
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
