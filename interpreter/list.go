package interpreter

import (
	"errors"
	"fmt"

	"github.com/danfragoso/milho/mir"
)

func __cdr(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to cdr")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if ev.Type() != mir.ListExpr {
		return nil, fmt.Errorf("cdr parameter must be a list, instead got %s", ev.Type())
	}

	exprs := ev.(*mir.ListExpression).Expressions
	if len(exprs) <= 1 {
		return mir.CreateNilExpression()
	}

	return mir.CreateListExpression(exprs[1:]...)
}

func __car(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to car")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if ev.Type() != mir.ListExpr {
		return nil, fmt.Errorf("car parameter must be a list, instead got %s", ev.Type())
	}

	exprs := ev.(*mir.ListExpression).Expressions
	if len(exprs) <= 1 {
		return mir.CreateNilExpression()
	}

	return exprs[0], nil
}

func __cons(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong number of args '2' passed to cons")
	}

	if params[1].Type() != mir.ListExpr {
		return nil, fmt.Errorf("cons second parameter must be a list, instead got %s", params[1].Type())
	}

	var exprs []mir.Expression
	exprs = append(exprs, params[0])
	exprs = append(exprs, params[1].(*mir.ListExpression).Expressions...)

	return mir.CreateListExpression(exprs...)
}

func __map(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong number of args '2' passed to map")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if ev.Type() != mir.ListExpr {
		return nil, fmt.Errorf("map first parameter must be a list, instead got %s", ev.Type())
	}

	fn, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if fn.Type() != mir.BuiltInExpr && fn.Type() != mir.FunctionExpr {
		return nil, fmt.Errorf("map second parameter must be a function, instead got %s", fn.Type())
	}

	exprs := ev.(*mir.ListExpression).Expressions
	var modExprs []mir.Expression
	for _, expr := range exprs {
		var r mir.Expression
		var e error

		if fn.Type() == mir.BuiltInExpr {
			r, e = fn.(*mir.BuiltInExpression).Function([]mir.Expression{expr}, session)
		} else {
			r, e = evaluateUserFunction(fn.(*mir.FunctionExpression), []mir.Expression{expr}, session)
		}

		if e != nil {
			return nil, fmt.Errorf("map function returned error: %s", e)
		}

		modExprs = append(modExprs, r)
	}

	return mir.CreateListExpression(modExprs...)
}
