package interpreter

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/danfragoso/milho/parser"
)

type ExpressionType int

func (e ExpressionType) String() string {
	return [...]string{"Nil", "Number", "Boolean", "Symbol", "Socket", "FunctionExpr", "String", "Byte", "List", "ErrorExpr", "BuiltIn"}[e]
}

const (
	NilExpr ExpressionType = iota
	NumberExpr
	BooleanExpr
	SymbolExpr
	SocketExpr
	FunctionExpr
	StringExpr
	ByteExpr
	ListExpr
	ErrorExpr
	BuiltInExpr
)

type Expression interface {
	Type() ExpressionType
	Value() string

	Parent() Expression
	setParent(Expression)
}

// Nil Expression
func createNilExpression() (*NilExpression, error) {
	return &NilExpression{}, nil
}

type NilExpression struct {
	ParentExpr Expression
}

func (e *NilExpression) Type() ExpressionType {
	return NilExpr
}

func (e *NilExpression) Value() string {
	return "Nil"
}

func (e *NilExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *NilExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// BuiltIn Expression
type BuiltInExpression struct {
	Identifier string
	Function   func([]Expression, *Session) (Expression, error)
}

func (e *BuiltInExpression) Type() ExpressionType {
	return BuiltInExpr
}

func (e *BuiltInExpression) Value() string {
	return "BuiltIn." + e.Identifier
}

func (e *BuiltInExpression) Parent() Expression {
	return nil
}

func (e *BuiltInExpression) setParent(parent Expression) {
}

// Number Expression
func createNumberExpression(numerator, denominator int64) (*NumberExpression, error) {
	return &NumberExpression{
		Numerator:   numerator,
		Denominator: denominator,
	}, nil
}

type NumberExpression struct {
	ParentExpr  Expression
	Numerator   int64
	Denominator int64
}

func (e *NumberExpression) Type() ExpressionType {
	return NumberExpr
}

func (e *NumberExpression) Value() string {
	r := strconv.FormatInt(e.Numerator, 10)
	if e.Denominator != 1 {
		r += "/" + strconv.FormatInt(e.Denominator, 10)
	}

	return r
}

func (e *NumberExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *NumberExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// List Expression
func createListExpression(expressions ...Expression) (*ListExpression, error) {
	listExpression := &ListExpression{
		Objects:     make(map[string]*Object),
		Expressions: expressions,
	}

	for _, expr := range listExpression.Expressions {
		expr.setParent(listExpression)
	}

	return listExpression, nil
}

type ListExpression struct {
	ParentExpr  Expression
	Expressions []Expression
	Objects     map[string]*Object
}

func (e *ListExpression) AddObjects(objects ...*Object) {
	for _, obj := range objects {
		e.Objects[obj.identifier] = obj
	}
}

func (e *ListExpression) FindObject(identifier string) *Object {
	return e.Objects[identifier]
}

func (e *ListExpression) Type() ExpressionType {
	return ListExpr
}

func (e *ListExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *ListExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

func (e *ListExpression) Value() string {
	v := "("
	for i, exp := range e.Expressions {
		v += exp.Value()
		if i+1 < len(e.Expressions) {
			v += " "
		}
	}

	return v + ")"
}

// Byte Expression
func createByteExpression(value byte) (*ByteExpression, error) {
	return &ByteExpression{
		Val: value,
	}, nil
}

type ByteExpression struct {
	ParentExpr Expression
	Val        byte
}

func (e *ByteExpression) Type() ExpressionType {
	return ByteExpr
}

func (e *ByteExpression) Value() string {
	return fmt.Sprintf("0x%X", e.Val)
}

func (e *ByteExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *ByteExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// Error Expression
func createErrorExpression(code, message string) (*ErrorExpression, error) {
	return &ErrorExpression{
		Code:    code,
		Message: message,
	}, nil
}

type ErrorExpression struct {
	ParentExpr Expression
	Code       string
	Message    string
}

func (e *ErrorExpression) Type() ExpressionType {
	return ErrorExpr
}

func (e *ErrorExpression) Value() string {
	return fmt.Sprintf("Error#%s", e.Code)
}

func (e *ErrorExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *ErrorExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// Boolean Expression
func createBooleanExpression(value bool) (*BooleanExpression, error) {
	return &BooleanExpression{
		Val: value,
	}, nil
}

type BooleanExpression struct {
	ParentExpr Expression
	Val        bool
}

func (e *BooleanExpression) Type() ExpressionType {
	return BooleanExpr
}

func (e *BooleanExpression) Value() string {
	if e.Val {
		return "True"
	}

	return "False"
}

func (e *BooleanExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *BooleanExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// Function Expression

func createFunctionExpression(identifier string, arities map[int]*fnArity) (*FunctionExpression, error) {
	fExpr := &FunctionExpression{
		Arities: arities,
		Objects: make(map[string]Expression),
	}

	if identifier == "" {
		identifier = fmt.Sprintf("%p", fExpr)
	}

	fExpr.Identifier = identifier
	return fExpr, nil
}

type FunctionExpression struct {
	ParentExpr Expression
	Identifier string

	Arities map[int]*fnArity
	Objects map[string]Expression
}

type fnArity struct {
	hasVariadic bool
	body        Expression
	parameters  []string
}

func (e *FunctionExpression) Type() ExpressionType {
	return FunctionExpr
}

func (e *FunctionExpression) Value() string {
	return "fn@" + e.Identifier
}

func (e *FunctionExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *FunctionExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// Symbol Expression
func createSymbolExpression(identifier string, expression Expression) (*SymbolExpression, error) {
	return &SymbolExpression{
		Identifier: identifier,
		Expression: expression,
	}, nil
}

type SymbolExpression struct {
	ParentExpr Expression
	Identifier string
	Expression Expression
}

func (e *SymbolExpression) Type() ExpressionType {
	return SymbolExpr
}

func (e *SymbolExpression) Value() string {
	return e.Identifier
}

func (e *SymbolExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *SymbolExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// Symbol Expression
type socketType string

const (
	ClientSocket socketType = "CLIENT"
	ServerSocket socketType = "SERVER"
)

type socketProtocol string

const (
	TCPSocket socketProtocol = "TCP"
	UDPSocket socketProtocol = "UDP"
)

func createSocketExpression(kind socketType, protocol socketProtocol, port int) (*SocketExpression, error) {
	var conn net.Conn

	if kind == ClientSocket {
		var err error
		conn, err = net.Dial(strings.ToLower(string(protocol)), fmt.Sprintf("localhost:%d", port))
		if err != nil {
			return nil, fmt.Errorf("Error creating %s Client socket: %s", protocol, err)
		}
	}

	return &SocketExpression{
		Kind:     kind,
		Protocol: protocol,
		Port:     port,

		conn: conn,
	}, nil
}

type SocketExpression struct {
	ParentExpr Expression
	Kind       socketType
	Protocol   socketProtocol
	conn       net.Conn
	Port       int
}

func (e *SocketExpression) Type() ExpressionType {
	return SocketExpr
}

func (e *SocketExpression) Value() string {
	return fmt.Sprintf("Socket@%s:%s:%d", e.Kind, e.Protocol, e.Port)
}

func (e *SocketExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *SocketExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// String Expression
func createStringExpression(value string) (*StringExpression, error) {
	return &StringExpression{
		Val: value,
	}, nil
}

type StringExpression struct {
	ParentExpr Expression
	Val        string
}

func (e *StringExpression) Type() ExpressionType {
	return StringExpr
}

func (e *StringExpression) Value() string {
	return fmt.Sprintf("\"%s\"", e.Val)
}

func (e *StringExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *StringExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// Expression Tree
func createExpressionTree(node *parser.Node) (Expression, error) {
	var expressions []Expression
	for _, childNode := range node.Nodes {
		expr, err := createExpressionTree(childNode)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	switch node.Type {
	case parser.Identifier:
		return createSymbolExpression(node.Identifier, nil)

	case parser.String:
		return createStringExpression(node.Identifier)

	case parser.Boolean:
		if node.Identifier == "True" {
			return createBooleanExpression(true)
		}

		return createBooleanExpression(false)

	case parser.Byte:
		b, _ := hex.DecodeString(node.Identifier[2:])
		return createByteExpression(b[0])

	case parser.Number:
		numberStr := strings.Split(node.Identifier, "/")

		numerator := int64(0)
		denominator := int64(1)

		numerator, _ = strconv.ParseInt(numberStr[0], 10, 64)
		if len(numberStr) == 2 {
			denominator, _ = strconv.ParseInt(numberStr[1], 10, 64)

			numerator, denominator = simplifyFraction(numerator, denominator)
		}

		return createNumberExpression(numerator, denominator)
	}

	return createListExpression(expressions...)
}

func simplifyFraction(numerator, denominator int64) (int64, int64) {
	divider := gcd(numerator, denominator)
	return numerator / divider, denominator / divider
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func printExpr(e Expression) {
	fmt.Print(sprintExpr(e, "", true), "\n\n")
}

func sprintExpr(e Expression, tab string, last bool) string {
	var str string

	str += fmt.Sprintf("%s*- ", tab)

	str += e.Type().String()
	str += fmt.Sprintf("#[%s]", e.Value())

	if last {
		tab += "   "
	} else {
		tab += "|  "
	}

	switch lExpr := e.(type) {
	case *ListExpression:
		for idx, expr := range lExpr.Expressions {
			str += sprintExpr(expr, tab, idx == len(lExpr.Expressions)-1)
		}
	}

	return str
}
