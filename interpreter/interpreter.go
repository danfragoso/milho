package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/mir"
	"github.com/danfragoso/milho/parser"
)

func RunFromSession(nodes []*parser.Node, session *mir.Session) ([]mir.Expression, error) {
	var expressions []mir.Expression
	var err error

	for _, node := range nodes {
		err = updateSession(session, node)
		if err != nil {
			return expressions, err
		}

		expr, err := mir.GenerateMIR(node)
		if err != nil {
			return nil, err
		}

		session.ExprTree = expr
		expressions = append(expressions, expr)
	}

	var evaluatedExpressions []mir.Expression
	for _, expr := range expressions {
		evaluated, err := evaluate(expr, session)
		if err != nil {
			stackTrace := session.CallStack.ToString(err)
			session.CallStack = &mir.CallStack{}
			return nil, fmt.Errorf(stackTrace)
		}

		evaluatedExpressions = append(evaluatedExpressions, evaluated)
	}

	return evaluatedExpressions, nil
}

func Run(nodes []*parser.Node) ([]mir.Expression, error) {
	var session *mir.Session
	var expressions []mir.Expression

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

		expr, err := mir.GenerateMIR(node)
		if err != nil {
			return nil, err
		}

		session.ExprTree = expr
		expressions = append(expressions, expr)
	}

	var evaluatedExpressions []mir.Expression
	for _, expr := range expressions {
		evaluated, err := evaluate(expr, session)
		if err != nil {
			stackTrace := session.CallStack.ToString(err)
			session.CallStack = &mir.CallStack{}
			return nil, fmt.Errorf(stackTrace)
		}

		evaluatedExpressions = append(evaluatedExpressions, evaluated)
	}

	return evaluatedExpressions, nil
}
