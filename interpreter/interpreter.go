package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
)

func Run(ast *parser.Node) (*Result, error) {
	res, err := eval(ast)
	if err != nil {
		return res, err
	}

	return res, nil
}

func eval(ast *parser.Node) (*Result, error) {
	res := &Result{}

	var results []*Result
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

	case parser.Number:
		return &Result{
			Type:  Number,
			Value: ast.Identifier,
		}, nil
	}

	return res, nil
}

func evalFunction(identifier string, params []*Result) (*Result, error) {
	switch identifier {
	case "+", "-":
		nParams, err := numberPrepareParams(params)
		if err != nil {
			return nil, err
		}

		switch identifier {
		case "+":
			return numberSum(nParams)

		case "-":
			return numberSub(nParams)
		}
	}

	return nil, fmt.Errorf("unknown function identifier '%s'", identifier)
}
