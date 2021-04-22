package interpreter

import (
	"fmt"
	"strings"
)

func resolveTypedExpression(exprType ExpressionType, expr Expression, session *Session) (Expression, error) {
	if expr.Type() == SymbolExpr {
		var err error
		expr, err = session.FindObject(expr.(*SymbolExpression).Identifier)
		if err != nil {
			return nil, err
		}
	}

	if expr.Type() != exprType {
		return nil, fmt.Errorf("expected resolved expression %s type to be %s, instead got %s", expr.Value(), exprType, expr.Type())
	}

	return expr, nil
}

func resolveExpression(expr Expression, session *Session) (Expression, error) {
	if expr.Type() == SymbolExpr {
		return session.FindObject(expr.(*SymbolExpression).Identifier)
	}

	return expr, nil
}

func __pr(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = resolveExpression(param, session)
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
		param, err = resolveExpression(param, session)
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
		param, err = resolveExpression(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			fmt.Print(strings.Trim(param.Value(), "\"") + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	return createNilExpression()
}

func __println(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = resolveExpression(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			fmt.Print(strings.Trim(param.Value(), "\"") + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	fmt.Println("")
	return createNilExpression()
}
