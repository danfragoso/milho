package interpreter

import (
	"fmt"
	"strings"

	"github.com/danfragoso/milho/mir"
)

func evaluate(expr mir.Expression, session *mir.Session) (mir.Expression, error) {
	switch expr.Type() {
	case mir.ListExpr:
		return evaluateList(expr, session)

	case mir.SymbolExpr:
		return evaluateSymbol(expr, session)
	}

	return expr, nil
}

func findExprObject(expr mir.Expression, identifier string) *mir.Object {
	if expr.Parent() == nil {
		return nil
	}

	if expr.Parent().Type() == mir.ListExpr {
		lst := expr.Parent().(*mir.ListExpression)
		obj := lst.FindObject(identifier)
		if obj != nil {
			return obj
		}

		return findExprObject(expr.Parent(), identifier)
	}

	return nil
}

func evaluateSymbol(expr mir.Expression, session *mir.Session) (mir.Expression, error) {
	symbol := expr.(*mir.SymbolExpression)
	if strings.HasPrefix(symbol.Identifier, "#!/") && strings.HasSuffix(symbol.Identifier, "milho") {
		return mir.CreateNilExpression()
	}

	obj, err := session.FindObject(symbol.Identifier)
	if err != nil {
		nObj := findExprObject(expr, symbol.Identifier)
		if nObj == nil {
			return nil, err
		}

		obj = nObj.Value()
	}

	return obj, nil
}

func evaluateList(expr mir.Expression, session *mir.Session) (mir.Expression, error) {
	expressions := expr.(*mir.ListExpression).Expressions
	if len(expressions) == 0 {
		return mir.CreateNilExpression()
	}

	firstExpr := expressions[0]
	if firstExpr.Type() == mir.ListExpr {
		var err error
		firstExpr, err = evaluate(firstExpr, session)
		if err != nil {
			return nil, err
		}
	}

	var obj mir.Expression
	var err error
	if firstExpr.Type() != mir.FunctionExpr {
		obj, err = session.FindObject(firstExpr.Value())
		if err != nil {
			return nil, err
		}
	}

	var result mir.Expression
	session.CallStack.Push(expressions)

	switch obj.Type() {
	case mir.BuiltInExpr:
		result, err = obj.(*mir.BuiltInExpression).Function(expressions[1:], session)
	case mir.FunctionExpr:
		result, err = evaluateUserFunction(obj.(*mir.FunctionExpression), expressions[1:], session)
	}

	if err != nil {
		return nil, err
	}

	session.CallStack.Pop()
	return result, nil
}

func evaluateUserFunction(fn *mir.FunctionExpression, params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	arity := fn.Arities[len(params)]
	if arity == nil {
		var sArity *mir.FnArity
		for _, arity := range fn.Arities {
			if arity.HasVariadic() {
				sArity = arity
			}
		}

		if sArity == nil {
			return nil, fmt.Errorf("Wrong number of args passed to function %s", fn.Value())
		}

		arity = sArity
	}

	var objs []*mir.Object
	for i, fnParam := range arity.Parameters() {
		if arity.HasVariadic() && i+1 == len(arity.Parameters()) {
			paramList, err := mir.CreateListExpression(params[i:]...)
			if err != nil {
				return nil, err
			}

			objs = append(objs, mir.CreateObject(paramList, fnParam))
		} else {
			value, err := evaluate(params[i], session)
			if err != nil {
				return nil, err
			}

			objs = append(objs, mir.CreateObject(value, fnParam))
		}
	}

	switch arity.Body().Type() {
	case mir.SymbolExpr:
		obj := mir.FindObject(objs, arity.Body().Value())
		if obj == nil {
			fObj, err := session.FindObject(arity.Body().Value())
			if err != nil {
				return nil, err
			}

			return fObj, nil
		}

		return obj.Value(), nil

	case mir.ListExpr:
		listExpr := arity.Body().(*mir.ListExpression)
		listExpr.AddObjects(objs...)
	}

	return evaluate(arity.Body(), session)
}
