package interpreter

import (
	"errors"
	"fmt"
)

func __cdr(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to cdr")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if ev.Type() != ListExpr {
		return nil, fmt.Errorf("cdr parameter must be a list, instead got %s", ev.Type())
	}

	exprs := ev.(*ListExpression).Expressions
	if len(exprs) <= 1 {
		return createNilExpression()
	}

	return createListExpression(exprs[1:]...)
}

func __car(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to car")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if ev.Type() != ListExpr {
		return nil, fmt.Errorf("car parameter must be a list, instead got %s", ev.Type())
	}

	exprs := ev.(*ListExpression).Expressions
	if len(exprs) <= 1 {
		return createNilExpression()
	}

	return exprs[0], nil
}

func __cons(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong number of args '2' passed to cons")
	}

	if params[1].Type() != ListExpr {
		return nil, fmt.Errorf("cons second parameter must be a list, instead got %s", params[1].Type())
	}

	var exprs []Expression
	exprs = append(exprs, params[0])
	exprs = append(exprs, params[1].(*ListExpression).Expressions...)

	return createListExpression(exprs...)
}

func __map(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong number of args '2' passed to map")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if ev.Type() != ListExpr {
		return nil, fmt.Errorf("map first parameter must be a list, instead got %s", ev.Type())
	}

	fn, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if fn.Type() != BuiltInExpr && fn.Type() != FunctionExpr {
		return nil, fmt.Errorf("map second parameter must be a function, instead got %s", fn.Type())
	}

	exprs := ev.(*ListExpression).Expressions
	var modExprs []Expression
	for _, expr := range exprs {
		var r Expression
		var e error

		if fn.Type() == BuiltInExpr {
			r, e = fn.(*BuiltInExpression).Function([]Expression{expr}, session)
		} else {
			r, e = evaluateUserFunction(fn.(*FunctionExpression), []Expression{expr}, session)
		}

		if e != nil {
			return nil, fmt.Errorf("map function returned error: %s", e)
		}

		modExprs = append(modExprs, r)
	}

	return createListExpression(modExprs...)
}
