package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
)

func Run(ast *parser.Node) (Result, error) {
	res, err := eval(ast)
	if err != nil {
		return res, err
	}

	return res, nil
}

func eval(ast *parser.Node) (Result, error) {
	var results []Result
	for _, childNode := range ast.Nodes {
		childResult, err := eval(childNode)
		if err != nil {
			return nil, err
		}

		results = append(results, childResult)
	}

	switch ast.Type {
	case parser.Function:
		result, err := evalFunction(ast.Identifier, results)
		if err != nil {
			return nil, err
		}

		return result, nil

	default:
		return createTypedResult(ResultType(ast.Type), ast.Identifier)
	}
}

func createTypedResult(t ResultType, v string) (Result, error) {
	switch t {
	case Number:
		return createNumberResult(v)
	case Boolean:
		return createBooleanResult(v)
	case Nil:
		return createNilResult()
	}

	return nil, fmt.Errorf("found unresolved %s '%s'", t, v)
}

func evalFunction(identifier string, params []Result) (Result, error) {
	switch identifier {
	case "=":
		return eq(params)
	case "/=":
		return neq(params)

	case "+", "-", "*", "/":
		nParams, err := numberPrepareParams(params)
		if err != nil {
			return nil, err
		}

		switch identifier {
		case "+":
			return numberSum(nParams)
		case "-":
			return numberSub(nParams)
		case "*":
			return numberMul(nParams)
		case "/":
			return numberDiv(nParams)
		}
	}

	return nil, fmt.Errorf("unknown function identifier '%s'", identifier)
}
