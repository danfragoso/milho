package interpreter

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/danfragoso/milho/mir"
)

func __eq(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
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

	return mir.CreateBooleanExpression(result)
}

func __negate(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '0' passed to cmp:[!] function")
	}

	param := params[0]
	param, err = evaluate(param, session)
	if err != nil {
		return nil, err
	}

	if param.Type() != mir.BooleanExpr {
		return nil, fmt.Errorf("Wrong type '%s' passed to cmp:[!] function", param.Type())
	}

	return mir.CreateBooleanExpression(!param.(*mir.BooleanExpression).Val)
}

func __if(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
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

	if fParam.Type() == mir.BooleanExpr &&
		!fParam.(*mir.BooleanExpression).Val {

		if len(params) == 3 {
			return evaluate(params[2], session)
		}

		return mir.CreateNilExpression()
	}

	return evaluate(params[1], session)
}

func __exec(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error

	if len(params) < 1 {
		return nil, fmt.Errorf("Too few args '%d' passed to exec function", len(params))
	}

	eParams := []mir.Expression{}
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() != mir.StringExpr {
			return nil, fmt.Errorf("Wrong type '%s' passed to exec function, expected String", param.Type())
		}

		eParams = append(eParams, param)
	}

	cmd := exec.Command(eParams[0].(*mir.StringExpression).Val)
	for _, param := range eParams[1:] {
		cmd.Args = append(cmd.Args, param.(*mir.StringExpression).Val)
	}

	out, err := cmd.Output()

	stdOut := strings.Trim(string(out), "\n")
	stdErr := ""
	exitCode := 0

	if err != nil {
		fErr := err.(*exec.ExitError)

		stdErr = strings.Trim(string(fErr.Stderr), "\n")
		exitCode = fErr.ExitCode()
	}

	exitCodeExpr, _ := mir.CreateNumberExpression(int64(exitCode), 1)

	stdOutExpr, _ := mir.CreateStringExpression(stdOut)
	stdErrExpr, _ := mir.CreateStringExpression(stdErr)

	return mir.CreateListExpression(exitCodeExpr, stdOutExpr, stdErrExpr)
}

// @TODO: Fix variadics and implement this in milho
func __execCode(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	output, _ := __exec(params, session)
	exprs := output.(*mir.ListExpression).Expressions

	return exprs[0], nil
}

func __execStdout(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	output, _ := __exec(params, session)
	exprs := output.(*mir.ListExpression).Expressions

	return exprs[1], nil
}

func __execStderr(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	output, _ := __exec(params, session)
	exprs := output.(*mir.ListExpression).Expressions

	return exprs[2], nil
}
