package interpreter

import (
	"errors"
	"fmt"
)

func __eq(params []Expression, session *Session) (Expression, error) {
	var err error
	if len(params) == 0 {
		return nil, errors.New("Wrong number of args '0' passed to cmp:[=] function")
	}

	lastParam := params[0]
	result := true

	for _, param := range params[1:] {
		lastParam, err = evaluate(lastParam, session)
		if err != nil {
			return nil, err
		}

		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if lastParam.Type() != param.Type() || lastParam.Value() != param.Value() {
			result = false
		}

		lastParam = param
	}

	return createBooleanExpression(result)
}

func __if(params []Expression, session *Session) (Expression, error) {
	var err error

	if len(params) < 2 {
		return nil, fmt.Errorf("Too few args '%d' passed to cmp:[if] function", len(params))
	} else if len(params) > 3 {
		return nil, fmt.Errorf("Too many args '%d' passed to cmp:[if] function", len(params))
	}

	fParam := params[0]
	fParam, err = evaluate(fParam, session)
	if err != nil {
		return nil, err
	}

	if fParam.Type() == BooleanExpr &&
		!fParam.(*BooleanExpression).Val {

		if len(params) == 3 {
			return params[2], nil
		}

		return createNilExpression()
	}

	return params[1], nil
}
