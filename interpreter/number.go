package interpreter

import (
	"errors"
)

func __add(params []Expression, session *Session) (Expression, error) {
	var acc int64

	for _, exp := range params {
		var err error
		exp, err = resolveTypedExpression(NumberExpr, exp, session)
		if err != nil {
			return nil, err
		}

		acc += exp.(*NumberExpression).Numerator
	}

	return createNumberExpression(acc, 1)
}

func __mul(params []Expression, session *Session) (Expression, error) {
	var acc int64 = 1

	for _, exp := range params {
		var err error
		exp, err = resolveTypedExpression(NumberExpr, exp, session)
		if err != nil {
			return nil, err
		}

		acc *= exp.(*NumberExpression).Numerator
	}

	return createNumberExpression(acc, 1)
}

func __sub(params []Expression, session *Session) (Expression, error) {
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

	nB := fExp.(*NumberExpression)
	if len(params) == 1 {
		acc = -nB.Numerator
	} else {
		acc = nB.Numerator
		for _, n := range params[1:] {
			n, err = evaluate(n, session)
			if err != nil {
				return nil, err
			}

			acc -= n.(*NumberExpression).Numerator
		}
	}

	return createNumberExpression(acc, 1)
}

func __div(params []Expression, session *Session) (Expression, error) {
	if len(params) == 0 {
		return nil, errors.New("Wrong number of args '0' passed to Number:[-] function")
	}

	var err error
	fExp := params[0]
	fExp, err = evaluate(fExp, session)
	if err != nil {
		return nil, err
	}

	acc := fExp.(*NumberExpression).Numerator
	if acc == 0 {
		return createNumberExpression(acc, 1)
	}

	for _, nE := range params[1:] {
		nE, err = evaluate(nE, session)
		if err != nil {
			return nil, err
		}

		n := nE.(*NumberExpression).Numerator

		if n == 0 {
			return nil, errors.New("Divide by zero error")
		}

		acc /= n
	}

	return createNumberExpression(acc, 1)
}
