package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

type Session struct {
	ExprTree   Expression
	SyntaxTree *parser.Node

	Objects   []*Object
	CallStack []string
}

func CreateSession(node *parser.Node) (*Session, error) {
	sess := &Session{
		SyntaxTree: node,
	}
	tokens, _ := tokenizer.Tokenize(builtinInjector + functionInjector)
	nodes, _ := parser.Parse(tokens)

	RunFromSession(nodes, sess)
	return sess, nil
}

func updateSession(session *Session, node *parser.Node) error {
	session.SyntaxTree = node

	return nil
}

func (s *Session) FindObject(identifier string) (Expression, error) {
	builtIn := BuiltIns[identifier]
	if builtIn != nil {
		return builtIn, nil
	}

	for _, obj := range s.Objects {
		if obj.Identifier() == identifier {
			return obj.Value(), nil
		}
	}

	return nil, fmt.Errorf("Symbol %s couldn't be resolved", identifier)
}

func (s *Session) AddObject(identifier string, expr Expression) error {
	s.Objects = addObjectToSession(s.Objects, &Object{
		identifier: identifier,
		value:      expr,
	})

	return nil
}

func addObjectToSession(objects []*Object, obj *Object) []*Object {
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
