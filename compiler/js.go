package compiler

import (
	"fmt"

	"github.com/danfragoso/milho/interpreter"
)

type JSValue interface {
	String() string
}

type JSValue_String struct {
	Value string
}

func (v *JSValue_String) String() string {
	return fmt.Sprintf("\"%s\"", v.Value)
}

type JSValue_Identifier struct {
	Name string
}

func (v *JSValue_Identifier) String() string {
	return v.Name
}

type JSValue_Undefined struct{}

func (v *JSValue_Undefined) String() string {
	return "undefined"
}

type JSValue_Declaration struct {
	Type       string
	Identifier string
	Value      JSValue
}

func (v *JSValue_Declaration) String() string {
	return fmt.Sprintf("%s %s = %s;", v.Type, v.Identifier, v.Value.String())
}

type JSValue_FunctionCall struct {
	Target string
	Params []JSValue
}

func (v *JSValue_FunctionCall) String() string {
	callStr := fmt.Sprintf("%s(", v.Target)
	for i, param := range v.Params {
		pValue := param.String()
		if i+1 != len(v.Params) {
			pValue += ","
		}

		callStr += pValue
	}

	return callStr + ");"
}

var SymbolMap = map[string]func(params []interpreter.Expression) JSValue{
	"def":   objDeclaration,
	"print": logCall,
}

func logCall(params []interpreter.Expression) JSValue {
	call := &JSValue_FunctionCall{}
	call.Target = "console.log"

	for _, param := range params {
		call.Params = append(call.Params, exprToJSValue(param))
	}

	return call
}

func objDeclaration(params []interpreter.Expression) JSValue {
	declaration := &JSValue_Declaration{}
	declaration.Type = "const"

	if len(params) != 2 {
		fmt.Println("wrong params for def JS")
	}

	if params[0].Type() != interpreter.SymbolExpr {
		fmt.Println("errow js def must be symbol")
	}

	declaration.Identifier = params[0].Value()
	declaration.Value = exprToJSValue(params[1])

	return declaration
}

func TranspileJS(expr interpreter.Expression) string {
	return transpileExpr(expr)
}

func transpileExpr(expr interpreter.Expression) string {
	switch expr.Type() {
	case interpreter.ListExpr:
		return transpileListExpr(expr).String()
	}

	return exprToJSValue(expr).String()
}

func transpileListExpr(expr interpreter.Expression) JSValue {
	expressions := expr.(*interpreter.ListExpression).Expressions
	if len(expressions) == 0 {
		return &JSValue_Undefined{}
	}

	firstExpr := expressions[0]
	switch firstExpr.Type() {
	case interpreter.SymbolExpr:
		sym := firstExpr.(*interpreter.SymbolExpression)
		return matchListSymbolExpr(sym, expressions[1:])
	}

	return &JSValue_String{}
}

func matchListSymbolExpr(symbol *interpreter.SymbolExpression, params []interpreter.Expression) JSValue {
	fn, found := SymbolMap[symbol.Identifier]
	if found {
		return fn(params)
	}

	return &JSValue_Undefined{}
}

func exprToJSValue(expr interpreter.Expression) JSValue {
	var returnValue JSValue
	switch expr.Type() {
	case interpreter.StringExpr:
		returnValue = &JSValue_String{
			Value: expr.(*interpreter.StringExpression).Val,
		}

	case interpreter.SymbolExpr:
		returnValue = &JSValue_Identifier{
			Name: expr.(*interpreter.SymbolExpression).Identifier,
		}
	}

	return returnValue
}
