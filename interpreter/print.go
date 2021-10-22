package interpreter

import (
	"fmt"
)

func resolveTypedExpression(exprType ExpressionType, expr Expression, session *Session) (Expression, error) {
	var err error
	expr, err = evaluate(expr, session)
	if err != nil {
		return nil, err
	}

	if expr.Type() != exprType {
		return nil, fmt.Errorf("expected resolved expression %s type to be %s, instead got %s", expr.Value(), exprType, expr.Type())
	}

	return expr, err
}

func __pr(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		fmt.Print(param.Value() + " ")
	}

	return createNilExpression()
}

func __prn(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		fmt.Print(param.Value() + " ")
	}

	fmt.Println("")
	return createNilExpression()
}

func __print(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			fmt.Print(param.Value()[1:len(param.Value())-1] + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	return createNilExpression()
}

func __println(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			fmt.Print(param.Value()[1:len(param.Value())-1] + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	fmt.Println("")
	return createNilExpression()
}

func __list(params []Expression, session *Session) (Expression, error) {
	for _, object := range session.Objects {
		fmt.Printf("%s [%s:%s]\n", object.Identifier(), object.value.Value(), object.value.Type())
	}

	return createNilExpression()
}
