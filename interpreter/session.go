package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
)

type Session struct {
	Nodes []*parser.Node

	Variables []*Var
	Functions []*Func
}

type Var struct {
	Identifier string
	Value      Result
}

func (v *Var) String() string {
	return fmt.Sprintf("%s:{%s}", v.Identifier, v.Value)
}

type Func struct {
	Identifier string
	Node       *parser.Node
}

func createSession(nodes []*parser.Node) (*Session, error) {
	sess := &Session{
		Nodes: nodes,
	}

	err := expandMacros(sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func expandMacros(session *Session) error {
	var macrosToRemove []int
	for i, node := range session.Nodes {
		if node.Type == parser.Macro {
			macrosToRemove = append(macrosToRemove, i)

			switch node.Identifier {
			case "def":
				variable, err := expandDefMacro(node)
				if err != nil {
					return err
				}

				session.Variables = append(session.Variables, variable)
			}
		}
	}

	for _, m := range macrosToRemove {
		session.Nodes = append(session.Nodes[:m], session.Nodes[m+1:]...)
	}

	return nil
}

func expandDefMacro(node *parser.Node) (*Var, error) {
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

	return &Var{
		Identifier: node.Nodes[0].Identifier,
		Value:      r,
	}, nil
}
