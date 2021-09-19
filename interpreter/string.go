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

func __asList(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("asList: expected 1 parameter, got %d, parameter must be a string", len(params))
	}

	str, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if str.Type() != StringExpr {
		return nil, fmt.Errorf("asList: expected first parameter to be a string, got %s", str.Type())
	}

	exprs := []Expression{}
	strValue := strings.Trim(str.Value(), "\"")
	for _, c := range strValue {
		s, err := createStringExpression(string(c))
		if err != nil {
			return nil, fmt.Errorf("asList: %s", err)
		}

		exprs = append(exprs, s)
	}

	return createListExpression(exprs...)
}
