package interpreter

import (
	"errors"
	"fmt"
	"time"

	"github.com/danfragoso/milho/mir"
)

func __def(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, errors.New("Wrong number of args '2' passed to def")
	}

	firstParam := params[0]
	if firstParam.Type() != mir.SymbolExpr {
		return nil, fmt.Errorf("First argument of def should be a symbol")
	}

	value, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	symbol, err := mir.CreateSymbolExpression(firstParam.(*mir.SymbolExpression).Identifier)
	if err != nil {
		return nil, err
	}

	err = session.AddObject(symbol.Identifier, value)
	if err != nil {
		return nil, err
	}

	return symbol, nil
}

func __let(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to def", len(params))
	}

	firstParam := params[0]
	if firstParam.Type() != mir.ListExpr {
		return nil, fmt.Errorf("First argument of let must be a list of key-value pairs")
	} else if len(firstParam.(*mir.ListExpression).Expressions)%2 != 0 {
		return nil, fmt.Errorf("Wrong number of atoms for key-value list, must be even")
	}

	var objs []*mir.Object
	kvPairs := firstParam.(*mir.ListExpression)
	for i := 0; i < len(kvPairs.Expressions); i += 2 {
		key := kvPairs.Expressions[i]
		val := kvPairs.Expressions[i+1]

		if key.Type() != mir.SymbolExpr {
			return nil, fmt.Errorf("Key '%s' for Value '%s' must be a symbol", key.Value(), val.Value())
		}

		val, err := evaluate(val, session)
		if err != nil {
			return nil, err
		}

		foundObj := mir.FindObject(objs, key.Value())
		if foundObj != nil {
			return nil, fmt.Errorf("Can't redeclare variables, '%s' is already defined as '%s'", key.Value(), foundObj.Value().Value())
		}

		objs = append(objs, mir.CreateObject(val, key.Value()))
	}

	secondParam := params[1]
	switch secondParam.Type() {
	case mir.SymbolExpr:
		obj := mir.FindObject(objs, secondParam.Value())
		if obj == nil {
			fObj, err := session.FindObject(secondParam.Value())
			if err != nil {
				return nil, err
			}

			return fObj, nil
		}

		return obj.Value(), nil

	case mir.ListExpr:
		listExpr := secondParam.(*mir.ListExpression)
		listExpr.AddObjects(objs...)
	}

	return evaluate(secondParam, session)
}

func __quote(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to quote")
	}

	return params[0], nil
}

func __type(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, errors.New("Wrong number of args '1' passed to type")
	}

	ev, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	return mir.CreateSymbolExpression(ev.Type().String())
}

func __fn(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to fn, needs to be 2", len(params))
	}

	firstParam := params[0]
	if firstParam.Type() != mir.ListExpr {
		return nil, fmt.Errorf("First argument of fn must be a list")
	}

	fLst := firstParam.(*mir.ListExpression)
	var fParams []string
	var hasVariadic bool

	for i, lstItem := range fLst.Expressions {
		if lstItem.Type() != mir.SymbolExpr {
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

	return mir.CreateFunctionExpression("", map[int]*mir.FnArity{
		len(fLst.Expressions): mir.CreateFnArity(hasVariadic, params[1], fParams),
	})
}

func __defn(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to defn, needs to be at least 3", len(params))
	}

	firstParam := params[0]
	if firstParam.Type() != mir.SymbolExpr {
		return nil, fmt.Errorf("First argument of defn must be a symbol")
	}

	secondParam := params[1]
	if secondParam.Type() != mir.ListExpr {
		return nil, fmt.Errorf("Second argument of defn must be a list of arguments")
	}

	fLst := secondParam.(*mir.ListExpression)
	var fParams []string
	for _, lstItem := range fLst.Expressions {
		if lstItem.Type() != mir.SymbolExpr {
			return nil, fmt.Errorf("Every item of defn second param must be a Symbol")
		}

		fParams = append(fParams, lstItem.Value())
	}

	fExp, err := mir.CreateFunctionExpression(firstParam.Value(), map[int]*mir.FnArity{
		len(fLst.Expressions): mir.CreateFnArity(false, params[2], fParams),
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

func __time(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to time, only one is allowed", len(params))
	}

	start := time.Now()
	ev, err := evaluate(params[0], session)

	if err != nil {
		return nil, err
	}

	duration, _ := mir.CreateStringExpression(fmt.Sprint(time.Since(start)))
	return mir.CreateListExpression(duration, ev)
}

func __do(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var res mir.Expression
	var err error

	res, err = mir.CreateNilExpression()
	for _, param := range params {
		res, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func __eval(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to eval, only one is allowed", len(params))
	}

	param, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	return evaluate(param, session)
}

func __mapCreate(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	exprMap := map[string]mir.Expression{}

	if len(params) < 1 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to mapCreate, needs to be at least 1", len(params))
	}

	for _, param := range params {
		if param.Type() != mir.ListExpr {
			return nil, fmt.Errorf("Wrong type of param '%s' passed to mapCreate, only list is allowed", param.Type())
		}

		lst := param.(*mir.ListExpression)
		if len(lst.Expressions) != 2 {
			return nil, fmt.Errorf("Wrong number of items in list passed to mapCreate, only 2 is allowed")
		}

		if lst.Expressions[0].Type() != mir.SymbolExpr {
			return nil, fmt.Errorf("First item of list passed to mapCreate must be a symbol")
		}

		value, _ := evaluate(lst.Expressions[1], session)
		exprMap[lst.Expressions[0].Value()] = value
	}

	return mir.CreateMapExpression(exprMap)
}

func __mapSet(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 3 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to mapSet, needs to be 3", len(params))
	}

	mapExpr, _ := evaluate(params[0], session)
	if mapExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("First argument of mapSet must be a map")
	}

	keyExpr, _ := evaluate(params[1], session)
	if keyExpr.Type() != mir.SymbolExpr {
		return nil, fmt.Errorf("Second argument of mapSet must be a symbol")
	}

	valueExpr, _ := evaluate(params[2], session)

	mapExpr.(*mir.MapExpression).Values[keyExpr.Value()] = valueExpr

	return mapExpr, nil
}

func __mapGet(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to mapGet, needs to be 2", len(params))
	}

	mapExpr, _ := evaluate(params[0], session)
	if mapExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("First argument of mapGet must be a map")
	}

	keyExpr, _ := evaluate(params[1], session)
	if keyExpr.Type() != mir.SymbolExpr {
		return nil, fmt.Errorf("Second argument of mapGet must be a symbol")
	}

	return mapExpr.(*mir.MapExpression).Values[keyExpr.Value()], nil
}

func __mapDelete(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to mapDelete, needs to be 2", len(params))
	}

	mapExpr, _ := evaluate(params[0], session)
	if mapExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("First argument of mapDelete must be a map")
	}

	keyExpr, _ := evaluate(params[1], session)
	if keyExpr.Type() != mir.SymbolExpr {
		return nil, fmt.Errorf("Second argument of mapDelete must be a symbol")
	}

	delete(mapExpr.(*mir.MapExpression).Values, keyExpr.Value())

	return mapExpr, nil
}

func __mapKeys(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("Wrong number of args '%d' passed to mapKeys, needs to be 1", len(params))
	}

	mapExpr, _ := evaluate(params[0], session)
	if mapExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("First argument of mapKeys must be a map")
	}

	var keys []mir.Expression
	for k := range mapExpr.(*mir.MapExpression).Values {
		key, _ := mir.CreateSymbolExpression(k)
		keys = append(keys, key)
	}

	return mir.CreateListExpression(keys...)
}
