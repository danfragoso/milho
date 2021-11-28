package interpreter

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
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

func __negate(params []Expression, session *Session) (Expression, error) {
	var err error
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '0' passed to cmp:[!] function")
	}

	param := params[0]
	param, err = evaluate(param, session)
	if err != nil {
		return nil, err
	}

	if param.Type() != BooleanExpr {
		return nil, fmt.Errorf("Wrong type '%s' passed to cmp:[!] function", param.Type())
	}

	return createBooleanExpression(!param.(*BooleanExpression).Val)
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
			return evaluate(params[2], session)
		}

		return createNilExpression()
	}

	return evaluate(params[1], session)
}

func __exec(params []Expression, session *Session) (Expression, error) {
	var err error

	if len(params) < 1 {
		return nil, fmt.Errorf("Too few args '%d' passed to exec function", len(params))
	}

	eParams := []Expression{}
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() != StringExpr {
			return nil, fmt.Errorf("Wrong type '%s' passed to exec function, expected String", param.Type())
		}

		eParams = append(eParams, param)
	}

	cmd := exec.Command(eParams[0].(*StringExpression).Val)
	for _, param := range eParams[1:] {
		cmd.Args = append(cmd.Args, param.(*StringExpression).Val)
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error executing command '%s': %s", cmd.Args, err)
	}

	return createStringExpression(strings.Trim(string(out), "\n"))
}
