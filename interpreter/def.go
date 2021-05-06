package interpreter

import (
	"errors"
	"fmt"
)

func __def(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong number of args '2' passed to def")
	}

	firstParam := params[0]
	if firstParam.Type() != SymbolExpr {
		return nil, fmt.Errorf("First argument of def should be a symbol")
	}

	value, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	symbol, err := createSymbolExpression(firstParam.(*SymbolExpression).Identifier, value)
	if err != nil {
		return nil, err
	}

	err = session.AddObject(symbol.Identifier, symbol.Expression)
	if err != nil {
		return nil, err
	}

	return symbol, nil
}

func __let(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to def", len(params))
	}

	firstParam := params[0]
	if firstParam.Type() != ListExpr {
		return nil, fmt.Errorf("First argument of let must be a list of key-value pairs")
	} else if len(firstParam.(*ListExpression).Expressions)%2 != 0 {
		return nil, fmt.Errorf("Wrong number of atoms for key-value list, must be even")
	}

	var objs []*Object
	kvPairs := firstParam.(*ListExpression)
	for i := 0; i < len(kvPairs.Expressions); i += 2 {
		key := kvPairs.Expressions[i]
		val := kvPairs.Expressions[i+1]

		if key.Type() != SymbolExpr {
			return nil, fmt.Errorf("Key '%s' for Value '%s' must be a symbol", key.Value(), val.Value())
		}

		val, err := evaluate(val, session)
		if err != nil {
			return nil, err
		}

		objs = append(objs, &Object{
			val, key.Value(),
		})
	}

	secondParam := params[1]
	switch secondParam.Type() {
	case SymbolExpr:
		obj := findObject(objs, secondParam.Value())
		if obj == nil {
			fObj, err := session.FindObject(secondParam.Value())
			if err != nil {
				return nil, err
			}

			return fObj, nil
		}

		return obj.value, nil

	case ListExpr:
		listExpr := secondParam.(*ListExpression)
		listExpr.AddObjects(objs...)
	}

	return evaluate(secondParam, session)
}

func __quote(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to quote")
	}

	return params[0], nil
}

func __type(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to type")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	return createSymbolExpression(ev.Type().String(), &SymbolExpression{})
}
