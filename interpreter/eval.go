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
	var expressions []Expression
	if len(expr.(*ListExpression).Expressions) == 2 &&
		expr.(*ListExpression).Expressions[0].Type() == SymbolExpr &&
		expr.(*ListExpression).Expressions[0].Value() == "quote" {
		return evaluateFunction("quote", expr.(*ListExpression).Expressions[1:], session)
	}

	for _, childExpr := range expr.(*ListExpression).Expressions {
		if childExpr.Type() == ListExpr {
			e, err := evaluate(childExpr, session)
			if err != nil {
				return nil, err
			}

			expressions = append(expressions, e)
		} else {
			expressions = append(expressions, childExpr)
		}
	}

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