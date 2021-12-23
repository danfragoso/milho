package mir

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
	return [...]string{"Nil", "Number", "Boolean", "Symbol", "Socket", "FunctionExpr", "String", "Byte", "List", "ErrorExpr", "BuiltIn", "StructExpr"}[e]
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
	StructExpr
)

type Expression interface {
	Type() ExpressionType
	Value() string

	Parent() Expression
	setParent(Expression)
}

// Nil Expression
func CreateNilExpression() (*NilExpression, error) {
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
func CreateNumberExpression(numerator, denominator int64) (*NumberExpression, error) {
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
func CreateListExpression(expressions ...Expression) (*ListExpression, error) {
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
func CreateByteExpression(value byte) (*ByteExpression, error) {
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
func CreateErrorExpression(code, message string) (*ErrorExpression, error) {
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
	return fmt.Sprintf("%s (%s)", e.Code, e.Message)
}

func (e *ErrorExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *ErrorExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

// Boolean Expression
func CreateBooleanExpression(value bool) (*BooleanExpression, error) {
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

func CreateFunctionExpression(identifier string, arities map[int]*FnArity) (*FunctionExpression, error) {
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

	Arities map[int]*FnArity
	Objects map[string]Expression
}

type FnArity struct {
	hasVariadic bool
	body        Expression
	parameters  []string
}

func (a *FnArity) HasVariadic() bool {
	return a.hasVariadic
}

func (a *FnArity) Body() Expression {
	return a.body
}

func (a *FnArity) Parameters() []string {
	return a.parameters
}

func CreateFnArity(hasVariadic bool, body Expression, parameters []string) *FnArity {
	return &FnArity{
		hasVariadic: hasVariadic,
		body:        body,
		parameters:  parameters,
	}
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
func CreateSymbolExpression(identifier string, expression Expression) (*SymbolExpression, error) {
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

func CreateSocketExpression(kind socketType, protocol socketProtocol, port int) (*SocketExpression, error) {
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
func CreateStringExpression(value string) (*StringExpression, error) {
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

type StructExpression struct {
	ParentExpr Expression
	Val        *Struct
}

func (e *StructExpression) Type() ExpressionType {
	return StructExpr
}

func (e *StructExpression) Value() string {
	return fmt.Sprintf("\"%s\"", e.Val.String())
}

func (e *StructExpression) Parent() Expression {
	return e.ParentExpr
}

func (e *StructExpression) setParent(parent Expression) {
	e.ParentExpr = parent
}

func GenerateMIR(node *parser.Node) (Expression, error) {
	var expressions []Expression
	for _, childNode := range node.Nodes {
		expr, err := GenerateMIR(childNode)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	switch node.Type {
	case parser.Identifier:
		return CreateSymbolExpression(node.Identifier, nil)

	case parser.String:
		return CreateStringExpression(node.Identifier)

	case parser.Boolean:
		if node.Identifier == "True" {
			return CreateBooleanExpression(true)
		}

		return CreateBooleanExpression(false)

	case parser.Byte:
		b, _ := hex.DecodeString(node.Identifier[2:])
		return CreateByteExpression(b[0])

	case parser.Number:
		numberStr := strings.Split(node.Identifier, "/")

		numerator := int64(0)
		denominator := int64(1)

		numerator, _ = strconv.ParseInt(numberStr[0], 10, 64)
		if len(numberStr) == 2 {
			denominator, _ = strconv.ParseInt(numberStr[1], 10, 64)

			numerator, denominator = simplifyFraction(numerator, denominator)
		}

		return CreateNumberExpression(numerator, denominator)
	}

	return CreateListExpression(expressions...)
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

func PrintExpr(e Expression) {
	fmt.Print(sprintExpr(e, "", true), "\n\n")
}

func SprintExpr(e Expression) string {
	return fmt.Sprint(sprintExpr(e, "", true), "\n\n")
}

func sprintExpr(e Expression, tab string, first bool) string {
	var str string
	tab += "  "
	if e.Type() == ListExpr {
		str += "\n" + tab + "List[" + "\n"

		for idx, expr := range e.(*ListExpression).Expressions {
			str += sprintExpr(expr, tab, idx == 0)
		}
	} else {
		if first {
			str += tab
		}

		str += e.Type().String() + "[" + e.Value() + "] "
	}

	if e.Type() == ListExpr {
		str += "]"
	}

	return str
}
