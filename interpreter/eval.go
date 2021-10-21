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
	if firstExpr.Type() == ListExpr {
		var err error
		firstExpr, err = evaluate(firstExpr, session)
		if err != nil {
			return nil, err
		}
	}

	var obj Expression
	var err error
	if firstExpr.Type() != FunctionExpr {
		obj, err = session.FindObject(firstExpr.Value())
		if err != nil {
			return nil, err
		}
	}

	switch obj.Type() {
	case BuiltInExpr:
		session.CallStack = append(session.CallStack, obj.Value())
		return obj.(*BuiltInExpression).Function(expressions[1:], session)
	case FunctionExpr:
		session.CallStack = append(session.CallStack, obj.Value())
		return evaluateUserFunction(obj.(*FunctionExpression), expressions[1:], session)
	}

	return nil, fmt.Errorf("undefined function '%s'", firstExpr.Value())
}

func evaluateUserFunction(fn *FunctionExpression, params []Expression, session *Session) (Expression, error) {
	arity := fn.Arities[len(params)]
	if arity == nil {
		var sArity *fnArity
		for _, arity := range fn.Arities {
			if arity.hasVariadic {
				sArity = arity
			}
		}

		if sArity == nil {
			return nil, fmt.Errorf("Wrong number of args passed to function %s", fn.Value())
		}

		arity = sArity
	}

	var objs []*Object
	for i, fnParam := range arity.parameters {
		value, err := evaluate(params[i], session)
		if err != nil {
			return nil, err
		}

		objs = append(objs, &Object{
			identifier: fnParam,
			value:      value,
		})
	}

	switch arity.body.Type() {
	case SymbolExpr:
		obj := findObject(objs, arity.body.Value())
		if obj == nil {
			fObj, err := session.FindObject(arity.body.Value())
			if err != nil {
				return nil, err
			}

			return fObj, nil
		}

		return obj.value, nil

	case ListExpr:
		listExpr := arity.body.(*ListExpression)
		listExpr.AddObjects(objs...)
	}

	return evaluate(arity.body, session)
}
