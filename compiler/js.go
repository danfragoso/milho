package compiler

import (
	"fmt"
	"strconv"

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

type JSValue_Number struct {
	Value string
}

func (v *JSValue_Number) String() string {
	return fmt.Sprintf("%s", v.Value)
}

type JSValue_Identifier struct {
	Name string
}

func isJSReserverd(identifier string) bool {
	reserved := []string{
		"break", "case", "catch", "continue", "debugger", "default", "delete", "do", "else", "finally", "for", "function", "if", "in", "instanceof", "new", "return", "switch", "this", "throw", "try", "typeof", "var", "void", "while", "with",
		"class", "const", "enum", "export", "extends", "import", "super",
		"implements", "interface", "let", "package", "private", "protected", "public", "static", "yield",
		"null", "true", "false",
	}
	for _, r := range reserved {
		if r == identifier {
			return true
		}
	}

	return false
}

func (v *JSValue_Identifier) String() string {
	if isJSReserverd(v.Name) {
		return "_" + v.Name
	}

	return v.Name
}

type JSValue_Undefined struct{}

func (v *JSValue_Undefined) String() string {
	return "undefined"
}

type JSValue_Lambda struct {
	Params []JSValue
	Body   JSValue
}

func (v *JSValue_Lambda) String() string {
	params := "("
	for _, p := range v.Params {
		params += p.String() + ", "
	}
	params = params[:len(params)-2] + ")"
	return fmt.Sprintf(` %s => {
	return %s
}`, params, v.Body)
}

type JSValue_Declaration struct {
	Type       string
	Identifier string
	Value      JSValue
}

func (v *JSValue_Declaration) String() string {
	return fmt.Sprintf("%s %s = %s;", v.Type, v.Identifier, v.Value.String())
}

type JSValue_FunctionDeclaration struct {
	Name   string
	Params []JSValue
	Body   JSValue
}

func (v *JSValue_FunctionDeclaration) String() string {
	params := "("
	for _, p := range v.Params {
		params += p.String() + ", "
	}
	params = params[:len(params)-2] + ")"
	return fmt.Sprintf(`function %s%s {
	return %s
}`, v.Name, params, v.Body)
}

type JSValue_IfElseBlock struct {
	Condition JSValue
	IfBody    JSValue
	ElseBody  JSValue
}

func (v *JSValue_IfElseBlock) String() string {
	return fmt.Sprintf("%s ? %s : %s;", v.Condition, v.IfBody, v.ElseBody)
	// 	return fmt.Sprintf(`if (%s) {
	// 	return %s
	// } else {
	// 	return %s
	// }`, v.Condition, v.IfBody, v.ElseBody)
}

type JSValue_Operator struct {
	Operator string
	Left     JSValue
	Right    JSValue
}

func (v *JSValue_Operator) String() string {
	return fmt.Sprintf("%s %s %s", v.Left, v.Operator, v.Right)
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

	return callStr + ")"
}

func methodCall(params []interpreter.Expression, methodName string) JSValue {
	call := &JSValue_FunctionCall{}

	if len(params) != 2 {
		fmt.Println("wrong params for method call JS")
	}

	call.Target = exprToJSValue(params[0]).String() + "." + methodName
	call.Params = append(call.Params, exprToJSValue(params[1]))

	return call
}

func lambda(params []interpreter.Expression) JSValue {
	lambda := &JSValue_Lambda{}

	if len(params) != 2 {
		fmt.Println("wrong params for lambda JS")
	}

	if params[0].Type() != interpreter.ListExpr {
		fmt.Println("error js lambda must be symbol")
	}

	for _, p := range params[0].(*interpreter.ListExpression).Expressions {
		lambda.Params = append(lambda.Params, exprToJSValue(p))
	}

	lambda.Body = exprToJSValue(params[1])
	return lambda
}

func ifElseBlock(params []interpreter.Expression) JSValue {
	ifElse := &JSValue_IfElseBlock{}
	if len(params) != 3 {
		fmt.Println("error if else block must have 3 params")
	}

	ifElse.Condition = exprToJSValue(params[0])
	ifElse.IfBody = exprToJSValue(params[1])
	ifElse.ElseBody = exprToJSValue(params[2])

	return ifElse
}

func logCall(params []interpreter.Expression) JSValue {
	call := &JSValue_FunctionCall{}
	call.Target = "console.log"

	for _, param := range params {
		call.Params = append(call.Params, exprToJSValue(param))
	}

	return call
}

func fnDeclaration(params []interpreter.Expression) JSValue {
	declaration := &JSValue_FunctionDeclaration{}

	if len(params) != 3 {
		fmt.Println("wrong params for def JS")
	}

	if params[0].Type() != interpreter.SymbolExpr {
		fmt.Println("error js defn must be symbol")
	}

	if params[1].Type() != interpreter.ListExpr {
		fmt.Println("error js defn params must be list")
	}

	declaration.Name = params[0].Value()
	for _, p := range params[1].(*interpreter.ListExpression).Expressions {
		declaration.Params = append(declaration.Params, exprToJSValue(p))
	}

	declaration.Body = exprToJSValue(params[2])
	return declaration
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
	switch symbol.Identifier {
	case "def":
		return objDeclaration(params)
	case "defn":
		return fnDeclaration(params)
	case "print":
		return logCall(params)
	case "if":
		return ifElseBlock(params)
	case "-":
		return makeOperator("-", params)
	case "*":
		return makeOperator("*", params)
	case "=":
		return makeOperator("===", params)
	case "fn":
		return lambda(params)
	case "split":
		return methodCall(params, "split")
	case "map":
		return methodCall(params, "map")
	case "join":
		return methodCall(params, "join")
	}

	return makeFunctionCall(symbol, params)
}

func makeOperator(operator string, params []interpreter.Expression) JSValue {
	return &JSValue_Operator{
		Operator: operator,
		Left:     exprToJSValue(params[0]),
		Right:    exprToJSValue(params[1]),
	}
}

func makeFunctionCall(symbol *interpreter.SymbolExpression, params []interpreter.Expression) JSValue {
	call := &JSValue_FunctionCall{}
	call.Target = symbol.Identifier

	for _, param := range params {
		call.Params = append(call.Params, exprToJSValue(param))
	}

	return call
}

func exprToJSValue(expr interpreter.Expression) JSValue {
	var returnValue JSValue = &JSValue_Undefined{}
	switch expr.Type() {
	case interpreter.ListExpr:
		returnValue = transpileListExpr(expr)

	case interpreter.StringExpr:
		returnValue = &JSValue_String{
			Value: expr.(*interpreter.StringExpression).Val,
		}

	case interpreter.NumberExpr:
		numerator := expr.(*interpreter.NumberExpression).Numerator
		denominator := expr.(*interpreter.NumberExpression).Denominator

		if denominator == 1 {
			returnValue = &JSValue_Number{
				Value: strconv.FormatInt(numerator, 10),
			}
		} else {
			returnValue = &JSValue_Number{
				Value: strconv.FormatFloat(float64(numerator)/float64(denominator), 'f', 64, 64),
			}
		}

	case interpreter.SymbolExpr:
		returnValue = &JSValue_Identifier{
			Name: expr.(*interpreter.SymbolExpression).Identifier,
		}
	}

	return returnValue
}
