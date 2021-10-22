package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
)

func RunFromSession(nodes []*parser.Node, session *Session) ([]Expression, error) {
	var expressions []Expression
	var err error

	for _, node := range nodes {
		err = updateSession(session, node)
		if err != nil {
			return expressions, err
		}

		expr, err := CreateExpressionTree(node)
		if err != nil {
			return nil, err
		}

		session.ExprTree = expr
		expressions = append(expressions, expr)
	}

	var evaluatedExpressions []Expression
	for _, expr := range expressions {
		evaluated, err := evaluate(expr, session)
		if err != nil {
			stackTrace := session.CallStack.ToString(err)
			session.CallStack = &CallStack{}
			return nil, fmt.Errorf(stackTrace)
		}

		evaluatedExpressions = append(evaluatedExpressions, evaluated)
	}

	return evaluatedExpressions, nil
}

func Run(nodes []*parser.Node) ([]Expression, error) {
	var session *Session
	var expressions []Expression

	for _, node := range nodes {
		var err error
		if session == nil {
			session, err = CreateSession(node)
		} else {
			err = updateSession(session, node)
		}

		if err != nil {
			return expressions, err
		}

		expr, err := CreateExpressionTree(node)
		if err != nil {
			return nil, err
		}

		session.ExprTree = expr
		expressions = append(expressions, expr)
	}

	var evaluatedExpressions []Expression
	for _, expr := range expressions {
		evaluated, err := evaluate(expr, session)
		if err != nil {
			stackTrace := session.CallStack.ToString(err)
			session.CallStack = &CallStack{}
			return nil, fmt.Errorf(stackTrace)
		}

		evaluatedExpressions = append(evaluatedExpressions, evaluated)
	}

	return evaluatedExpressions, nil
}
