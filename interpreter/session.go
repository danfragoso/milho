package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
)

type Session struct {
	Tree *parser.Node

	Objects []Object
}

func createSession(node *parser.Node) (*Session, error) {
	sess := &Session{
		Tree: node,
	}

	err := expandMacros(sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func expandMacros(session *Session) error {

	return nil
}

func expandDefMacro(node *parser.Node) (Object, error) {
	if len(node.Nodes) != 2 {
		return nil, fmt.Errorf("Wrong number of args on def macro, expected 2, got %d", len(node.Nodes))
	}

	if node.Nodes[0].Type != parser.Identifier {
		return nil, fmt.Errorf("Wrong type for def macro first argument, it must be an Identifier")
	}

	r, err := createTypedResult(ResultType(node.Nodes[1].Type), node.Nodes[1].Identifier)
	if err != nil {
		return nil, err
	}

	return &VariableObj{
		objectType: VariableObject,
		identifier: node.Nodes[0].Identifier,
		value:      r,
	}, nil
}
