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

	err := expandMacros(node, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func updateSession(session *Session, node *parser.Node) error {
	session.Tree = node
	err := expandMacros(node, session)
	if err != nil {
		return err
	}

	return nil
}

func expandMacros(node *parser.Node, session *Session) error {
	if node.Type != parser.Macro {
		for _, child := range node.Nodes {
			err := expandMacros(child, session)
			if err != nil {
				return err
			}
		}
	} else {
		switch node.Identifier {
		case "def":
			obj, err := expandDefMacro(node)
			if err != nil {
				return err
			}

			session.Objects = addObjectToSession(session.Objects, obj)
		}
	}

	return nil
}

func addObjectToSession(objects []Object, obj Object) []Object {
	objIdx := -1
	for i, object := range objects {
		if object.Identifier() == obj.Identifier() {
			objIdx = i
		}
	}

	if objIdx == -1 {
		objects = append(objects, obj)
	} else {
		objects[objIdx] = obj
	}

	return objects
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
		result:     r,
	}, nil
}
