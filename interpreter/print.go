package interpreter

import (
	"fmt"
	"time"
)

func resolveTypedExpression(exprType ExpressionType, expr Expression, session *Session) (Expression, error) {
	var err error
	expr, err = evaluate(expr, session)
	if err != nil {
		return nil, err
	}

	if expr.Type() != exprType {
		return nil, fmt.Errorf("expected resolved expression %s type to be %s, instead got %s", expr.Value(), exprType, expr.Type())
	}

	return expr, err
}

func __pr(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		fmt.Print(param.Value() + " ")
	}

	return createNilExpression()
}

func __prn(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		fmt.Print(param.Value() + " ")
	}

	fmt.Println("")
	return createNilExpression()
}

func __print(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			fmt.Print(param.Value()[1:len(param.Value())-1] + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	return createNilExpression()
}

func __println(params []Expression, session *Session) (Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			fmt.Print(param.Value()[1:len(param.Value())-1] + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	fmt.Println("")
	return createNilExpression()
}

func __list(params []Expression, session *Session) (Expression, error) {
	for _, object := range session.Objects {
		fmt.Printf("%s [%s:%s]\n", object.Identifier(), object.value.Value(), object.value.Type())
	}

	return createNilExpression()
}

func __sleep(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("expected 1 parameter, instead got %d", len(params))
	}

	duration, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if duration.Type() != NumberExpr {
		return nil, fmt.Errorf("expected sleep parameter to be an number, instead got %s", params[0].Type())
	}

	sleepTime := duration.(*NumberExpression).Numerator
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return createNilExpression()
}

func __range(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 && params[0].Type() != NumberExpr && params[1].Type() != NumberExpr {
		return nil, fmt.Errorf("expected range parameters to be numbers, instead got %s and %s", params[0].Type(), params[1].Type())
	}

	min := params[0].(*NumberExpression).Numerator
	max := params[1].(*NumberExpression).Numerator

	if min > max {
		return nil, fmt.Errorf("expected range min to be less than max, instead got %d and %d", min, max)
	}

	expressions := []Expression{}
	for i := min; i <= max; i++ {
		expr, _ := createNumberExpression(i, 1)
		expressions = append(expressions, expr)
	}

	return createListExpression(expressions...)
}
