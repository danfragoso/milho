package interpreter

import (
	"fmt"
	"time"

	"github.com/danfragoso/milho/mir"
)

func resolveTypedExpression(exprType mir.ExpressionType, expr mir.Expression, session *mir.Session) (mir.Expression, error) {
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

func __pr(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		fmt.Print(param.Value() + " ")
	}

	return mir.CreateNilExpression()
}

func __prn(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		fmt.Print(param.Value() + " ")
	}

	fmt.Println("")
	return mir.CreateNilExpression()
}

func __print(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == mir.StringExpr {
			fmt.Print(param.Value()[1:len(param.Value())-1] + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	return mir.CreateNilExpression()
}

func __println(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == mir.StringExpr {
			fmt.Print(param.Value()[1:len(param.Value())-1] + " ")
		} else {
			fmt.Print(param.Value() + " ")
		}
	}

	fmt.Println("")
	return mir.CreateNilExpression()
}

func __list(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	for _, object := range session.Objects {
		fmt.Printf("%s [%s:%s]\n", object.Identifier(), object.Value().Value(), object.Value().Type())
	}

	return mir.CreateNilExpression()
}

func __sleep(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("expected 1 parameter, instead got %d", len(params))
	}

	duration, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if duration.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("expected sleep parameter to be an number, instead got %s", params[0].Type())
	}

	sleepTime := duration.(*mir.NumberExpression).Numerator
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return mir.CreateNilExpression()
}

func __range(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 && params[0].Type() != mir.NumberExpr && params[1].Type() != mir.NumberExpr {
		return nil, fmt.Errorf("expected range parameters to be numbers, instead got %s and %s", params[0].Type(), params[1].Type())
	}

	min := params[0].(*mir.NumberExpression).Numerator
	max := params[1].(*mir.NumberExpression).Numerator

	if min > max {
		return nil, fmt.Errorf("expected range min to be less than max, instead got %d and %d", min, max)
	}

	expressions := []mir.Expression{}
	for i := min; i <= max; i++ {
		expr, _ := mir.CreateNumberExpression(i, 1)
		expressions = append(expressions, expr)
	}

	return mir.CreateListExpression(expressions...)
}
