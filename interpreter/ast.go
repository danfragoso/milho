package interpreter

import "github.com/danfragoso/milho/parser"

type Node struct {
	Type ResultType

	Parent   *Node
	Children []*Node
}

func expand(ast *parser.Node)
