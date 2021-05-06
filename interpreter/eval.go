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

func findExprObject(expr Expression, identifier string) *Object {
	if expr.Parent() == nil {
		return nil
	}

	if expr.Parent().Type() == ListExpr {
		lst := expr.Parent().(*ListExpression)
		obj := lst.FindObject(identifier)
		if obj != nil {
			return obj
		}

		return findExprObject(expr.Parent(), identifier)
	}

	return nil
}

func evaluateSymbol(expr Expression, session *Session) (Expression, error) {
	symbol := expr.(*SymbolExpression)
	obj, err := session.FindObject(symbol.Identifier)
	if err != nil {
		nObj := findExprObject(expr, symbol.Identifier)
		if nObj == nil {
			return nil, err
		}

		obj = nObj.value
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

	obj, err := session.FindObject(firstExpr.Value())
	if err != nil {
		return nil, err
	}

	switch obj.Type() {
	case BuiltInExpr:
		return obj.(*BuiltInExpression).Function(expressions[1:], session)
	}

	return nil, fmt.Errorf("undefined function '%s'", firstExpr.Value())
}
