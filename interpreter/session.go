package interpreter

import (
	"github.com/danfragoso/milho/mir"
	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

func CreateSession(node *parser.Node) (*mir.Session, error) {
	sess := &mir.Session{
		SyntaxTree: node,
		CallStack:  &mir.CallStack{},
	}
	tokens, _ := tokenizer.Tokenize(builtinInjector + functionInjector)
	nodes, _ := parser.Parse(tokens)

	RunFromSession(nodes, sess)
	return sess, nil
}

func updateSession(session *mir.Session, node *parser.Node) error {
	session.SyntaxTree = node

	return nil
}
