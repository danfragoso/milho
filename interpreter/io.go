package interpreter

import "fmt"

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

		fmt.Print(param)
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

		fmt.Print(param.Value())
	}

	fmt.Println("")
	return createNilExpression()
}
