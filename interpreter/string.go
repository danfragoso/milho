package interpreter

import (
	"fmt"
	"strings"
)

func __str(params []Expression, session *Session) (Expression, error) {
	var err error
	var resultStr string
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			resultStr += param.(*StringExpression).Val
		} else {
			resultStr += param.Value()
		}
	}

	return createStringExpression(resultStr)
}

func __split(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("split: expected 2 parameter, got %d, parameters must be string and separator", len(params))
	}

	str, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if str.Type() != StringExpr {
		return nil, fmt.Errorf("split: expected first parameter to be a string, got %s", str.Type())
	}

	sep, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if str.Type() != StringExpr {
		return nil, fmt.Errorf("split: expected second parameter to be a string, got %s", str.Type())
	}

	exprs := []Expression{}
	strValue := strings.Trim(str.Value(), "\"")
	sepValue := strings.Trim(sep.Value(), "\"")

	strs := strings.Split(strValue, sepValue)
	for _, st := range strs {
		s, err := createStringExpression(st)
		if err != nil {
			return nil, fmt.Errorf("split: error creating string expression from %s: %s", st, err)
		}

		exprs = append(exprs, s)
	}

	return createListExpression(exprs...)
}

func __join(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("join: expected 2 parameter, got %d, parameters must be list of strings and separator", len(params))
	}

	lst, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if lst.Type() != ListExpr {
		return nil, fmt.Errorf("join: expected first parameter to be a list, got %s", lst.Type())
	}

	var strList []string
	for _, expr := range lst.(*ListExpression).Expressions {
		if expr.Type() != StringExpr {
			return nil, fmt.Errorf("join: expected list to contain only strings, got %s", expr.Type())
		}

		strList = append(strList, expr.(*StringExpression).Val)
	}

	sep, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if sep.Type() != StringExpr {
		return nil, fmt.Errorf("join: expected second parameter to be a string, got %s", sep.Type())
	}

	sepValue := strings.Trim(sep.Value(), "\"")
	return createStringExpression(strings.Join(strList, sepValue))
}
