package interpreter

import (
	"errors"
	"fmt"
	"time"
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

		foundObj := findObject(objs, key.Value())
		if foundObj != nil {
			return nil, fmt.Errorf("Can't redeclare variables, '%s' is already defined as '%s'", key.Value(), foundObj.value.Value())
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

func __fn(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to fn, needs to be 2", len(params))
	}

	firstParam := params[0]
	if firstParam.Type() != ListExpr {
		return nil, fmt.Errorf("First argument of fn must be a list")
	}

	fLst := firstParam.(*ListExpression)
	var fParams []string
	var hasVariadic bool

	for i, lstItem := range fLst.Expressions {
		if lstItem.Type() != SymbolExpr {
			return nil, fmt.Errorf("Every item of fn first param must be a Symbol")
		}

		isVariadic := lstItem.Value() == "+rest"
		if i+1 == len(fLst.Expressions) && isVariadic {
			hasVariadic = true
		} else if isVariadic {
			return nil, fmt.Errorf("Only the last param can be a variadic")
		}

		fParams = append(fParams, lstItem.Value())
	}

	return createFunctionExpression("", map[int]*fnArity{
		len(fLst.Expressions): {
			hasVariadic: hasVariadic,
			parameters:  fParams,
			body:        params[1],
		},
	})
}

func __defn(params []Expression, session *Session) (Expression, error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to defn, needs to be at least 3", len(params))
	}

	firstParam := params[0]
	if firstParam.Type() != SymbolExpr {
		return nil, fmt.Errorf("First argument of defn must be a symbol")
	}

	secondParam := params[1]
	if secondParam.Type() != ListExpr {
		return nil, fmt.Errorf("Second argument of defn must be a list of arguments")
	}

	fLst := secondParam.(*ListExpression)
	var fParams []string
	for _, lstItem := range fLst.Expressions {
		if lstItem.Type() != SymbolExpr {
			return nil, fmt.Errorf("Every item of defn second param must be a Symbol")
		}

		fParams = append(fParams, lstItem.Value())
	}

	fExp, err := createFunctionExpression(firstParam.Value(), map[int]*fnArity{
		len(fLst.Expressions): {
			parameters: fParams,
			body:       params[2],
		},
	})

	if err != nil {
		return nil, err
	}

	err = session.AddObject(firstParam.Value(), fExp)
	if err != nil {
		return nil, err
	}

	return fExp, nil
}

func __time(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to time, only one is allowed", len(params))
	}

	start := time.Now()
	ev, err := evaluate(params[0], session)

	if err != nil {
		return nil, err
	}

	duration, _ := createStringExpression(fmt.Sprint(time.Since(start)))
	return createListExpression(duration, ev)
}

func __progn(params []Expression, session *Session) (Expression, error) {
	var res Expression
	var err error

	res, err = createNilExpression()
	for _, param := range params {
		res, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func __eval(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to eval, only one is allowed", len(params))
	}

	param, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	return evaluate(param, session)
}
