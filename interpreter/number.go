package interpreter

import (
	"errors"

	"github.com/danfragoso/milho/mir"
)

func __add(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var acc int64

	for _, exp := range params {
		var err error
		exp, err = resolveTypedExpression(mir.NumberExpr, exp, session)
		if err != nil {
			return nil, err
		}

		acc += exp.(*mir.NumberExpression).Numerator
	}

	return mir.CreateNumberExpression(acc, 1)
}

func __mul(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var acc int64 = 1

	for _, exp := range params {
		var err error
		exp, err = resolveTypedExpression(mir.NumberExpr, exp, session)
		if err != nil {
			return nil, err
		}

		acc *= exp.(*mir.NumberExpression).Numerator
	}

	return mir.CreateNumberExpression(acc, 1)
}

func __sub(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) == 0 {
		return nil, errors.New("Wrong number of args '0' passed to Number:[-] function")
	}

	var acc int64
	var err error

	fExp := params[0]
	fExp, err = evaluate(fExp, session)
	if err != nil {
		return nil, err
	}

	nB := fExp.(*mir.NumberExpression)
	if len(params) == 1 {
		acc = -nB.Numerator
	} else {
		acc = nB.Numerator
		for _, n := range params[1:] {
			n, err = evaluate(n, session)
			if err != nil {
				return nil, err
			}

			acc -= n.(*mir.NumberExpression).Numerator
		}
	}

	return mir.CreateNumberExpression(acc, 1)
}

func __div(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) == 0 {
		return nil, errors.New("Wrong number of args '0' passed to Number:[-] function")
	}

	var err error
	fExp := params[0]
	fExp, err = evaluate(fExp, session)
	if err != nil {
		return nil, err
	}

	acc := fExp.(*mir.NumberExpression).Numerator
	if acc == 0 {
		return mir.CreateNumberExpression(acc, 1)
	}

	for _, nE := range params[1:] {
		nE, err = evaluate(nE, session)
		if err != nil {
			return nil, err
		}

		n := nE.(*mir.NumberExpression).Numerator

		if n == 0 {
			return nil, errors.New("Divide by zero error")
		}

		acc /= n
	}

	return mir.CreateNumberExpression(acc, 1)
}
