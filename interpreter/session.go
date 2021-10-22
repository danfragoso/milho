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
	CallStack *CallStack
}

type CallStack struct {
	expressions [][]Expression
}

func (s *CallStack) IsEmpty() bool {
	return len(s.expressions) == 0
}

func (s *CallStack) Push(params []Expression) {
	s.expressions = append(s.expressions, params)
}

func (s *CallStack) Pop() {
	if s.IsEmpty() {
		return
	} else {
		index := len(s.expressions) - 1
		s.expressions = s.expressions[:index]
	}
}

func spaceMul(space string, mul int) string {
	fullSpace := space
	for i := 0; i < mul; i++ {
		fullSpace += space
	}

	return fullSpace
}

func callString(call []Expression) string {
	callStr := ""
	for i, expr := range call {
		callStr += expr.Value()
		if len(call) != i+1 {
			callStr += " "
		}
	}

	return callStr
}

func (s *CallStack) ToString(e error) string {
	callStack := "\n"
	if len(s.expressions) == 0 {
		return e.Error()
	}

	for step, call := range s.expressions {
		callStack += fmt.Sprintf("|%s (%s ", spaceMul("  ", step), callString(call))

		if len(s.expressions) == (step + 1) {
			callStack += fmt.Sprintf("\n|%s    Error@%s[%s]%s", spaceMul("  ", step), call[0].Value(), e.Error(), spaceMul(")", step))
		}

		callStack += "\n"
	}

	return callStack + "\n"
}

func CreateSession(node *parser.Node) (*Session, error) {
	sess := &Session{
		SyntaxTree: node,
		CallStack:  &CallStack{},
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
