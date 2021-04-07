package interpreter

import (
	"fmt"

	"github.com/danfragoso/milho/parser"
)

func eval(ast *parser.Node, sess *Session) (Result, error) {
	if ast.Parent != nil && ast.Parent.Type == parser.Macro {
		return createNilResult()
	}

	var results []Result
	for _, childNode := range ast.Nodes {
		childResult, err := eval(childNode, sess)
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

	case parser.Macro:
		result, err := evalMacro(ast, sess.Objects)
		if err != nil {
			return nil, err
		}

		return result, nil

	case parser.Identifier:
		result, err := evalIdentifier(ast.Identifier, sess)
		if err != nil {
			return nil, err
		}

		return result, nil

	default:
		return createTypedResult(ResultType(ast.Type), ast.Identifier)
	}
}

func evalIdentifier(identifier string, sess *Session) (Result, error) {
	for _, obj := range sess.Objects {
		if obj.Identifier() == identifier {
			if obj.Result().Type() != Pending {
				return obj.Result(), nil
			}

			pendingResult := obj.Result().(*PendingResult).Tree
			pendingResult.Parent = nil
			return eval(pendingResult, sess)
		}
	}

	return nil, fmt.Errorf("Identifier %s value couldn’t be resolved", identifier)
}

func evalMacro(node *parser.Node, objs []Object) (Result, error) {
	if len(node.Nodes) < 2 {
		return nil, fmt.Errorf("Malformatted macro %s couldn’t be evaluated", node.Identifier)
	}

	for _, obj := range objs {
		if node.Nodes[0].Identifier == obj.Identifier() {
			return createObjectResult(obj)
		}
	}

	return nil, fmt.Errorf("Macro %s was not correctly expanded", node.Nodes[0].Identifier)
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
