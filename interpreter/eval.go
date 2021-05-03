package interpreter

import (
	"fmt"
)

func evaluate(expr Expression, session *Session) (Expression, error) {
	switch expr.Type() {
	case ListExpr:
		return evaluateList(expr, session)

	case SymbolExpr:
		return evaluateSymbol(expr, session)
	}

	return expr, nil
}

func evaluateSymbol(expr Expression, session *Session) (Expression, error) {
	symbol := expr.(*SymbolExpression)
	obj, err := session.FindObject(symbol.Identifier)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func evaluateList(expr Expression, session *Session) (Expression, error) {
	expressions := expr.(*ListExpression).Expressions
	if len(expressions) == 0 {
		return createNilExpression()
	}

	firstExpr := expressions[0]
	if firstExpr.Type() != SymbolExpr {
		return nil, fmt.Errorf("%s can't be a function", firstExpr.Value())
	}

	return evaluateFunction(firstExpr.Value(), expressions[1:], session)
}

func evaluateFunction(identifier string, params []Expression, session *Session) (Expression, error) {
	switch identifier {
	case "def":
		return __def(params, session)

	case "=":
		return __eq(params, session)

	case "if":
		return __if(params, session)

	case "+":
		return __sum(params, session)
	case "-":
		return __sub(params, session)
	case "*":
		return __mul(params, session)
	case "/":
		return __div(params, session)

	case "pr":
		return __pr(params, session)
	case "prn":
		return __prn(params, session)
	case "print":
		return __print(params, session)
	case "println":
		return __println(params, session)

	case "str":
		return __str(params, session)

	case "quote":
		return __quote(params, session)
	}

	return nil, fmt.Errorf("undefined function '%s'", identifier)
}
