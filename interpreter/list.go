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
