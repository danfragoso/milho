package interpreter

import (
	"fmt"
	"strings"

	"github.com/danfragoso/milho/mir"
)

func __str(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error
	var resultStr string
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == mir.StringExpr {
			resultStr += param.(*mir.StringExpression).Val
		} else {
			resultStr += param.Value()
		}
	}

	return mir.CreateStringExpression(resultStr)
}

func __split(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("split: expected 2 parameter, got %d, parameters must be string and separator", len(params))
	}

	str, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if str.Type() != mir.StringExpr {
		return nil, fmt.Errorf("split: expected first parameter to be a string, got %s", str.Type())
	}

	sep, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if str.Type() != mir.StringExpr {
		return nil, fmt.Errorf("split: expected second parameter to be a string, got %s", str.Type())
	}

	exprs := []mir.Expression{}
	strValue := strings.Trim(str.Value(), "\"")
	sepValue := strings.Trim(sep.Value(), "\"")

	strs := strings.Split(strValue, sepValue)
	for _, st := range strs {
		s, err := mir.CreateStringExpression(st)
		if err != nil {
			return nil, fmt.Errorf("split: error creating string expression from %s: %s", st, err)
		}

		exprs = append(exprs, s)
	}

	return mir.CreateListExpression(exprs...)
}

func __join(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("join: expected 2 parameter, got %d, parameters must be list of strings and separator", len(params))
	}

	lst, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if lst.Type() != mir.ListExpr {
		return nil, fmt.Errorf("join: expected first parameter to be a list, got %s", lst.Type())
	}

	var strList []string
	for _, expr := range lst.(*mir.ListExpression).Expressions {
		if expr.Type() != mir.StringExpr {
			return nil, fmt.Errorf("join: expected list to contain only strings, got %s", expr.Type())
		}

		strList = append(strList, expr.(*mir.StringExpression).Val)
	}

	sep, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if sep.Type() != mir.StringExpr {
		return nil, fmt.Errorf("join: expected second parameter to be a string, got %s", sep.Type())
	}

	sepValue := strings.Trim(sep.Value(), "\"")
	return mir.CreateStringExpression(strings.Join(strList, sepValue))
}
