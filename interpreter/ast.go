package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
)

type Scope struct {
	Variables []*Var
	Functions []*Func
}
type Var struct {
	Identifier string

	Value Result
}

type Func struct {
	Identifier string

	Value Result
}
type Node struct {
	Type ResultType

	Parent   *Node
	Children []*Node

	Scope *Scope
}

func (n *Node) String() string {
	return n.Sprint("", true)
}

func (n *Node) Sprint(tab string, last bool) string {
	str := fmt.Sprintf("\n%s+- %s:[]", tab, n.Type)
	if last {
		tab += "   "
	} else {
		tab += "|  "
	}

	for idx, cN := range n.Children {
		str += cN.Sprint(tab, idx == len(n.Children)-1)
	}

	return str
}

func expand(ast *parser.Node) *Node {
	rootNode := parserNodeToExecNode(ast)
	for _, n := range ast.Nodes {
		rootNode.Children = append(rootNode.Children,
			expand(n),
		)
	}

	return rootNode
}

func parserNodeToExecNode(pNode *parser.Node) *Node {
	node := &Node{
		Type: ResultType(pNode.Type),
	}

	return node
}

func expandMacros(node *Node, pNode *parser.Node) {
	switch pNode.Identifier {
	case "def":

	}
}
