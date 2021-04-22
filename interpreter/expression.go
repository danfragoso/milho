package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danfragoso/milho/parser"
)

type ExpressionType int

func (e ExpressionType) String() string {
	return [...]string{"Nil", "Number", "Boolean", "Symbol", "List"}[e]
}

const (
	NilExpr ExpressionType = iota
	NumberExpr
	BooleanExpr
	SymbolExpr
	ListExpr
)

type Expression interface {
	Type() ExpressionType
	Value() string
}

// Symbol Expression
func createNilExpression() (*NilExpression, error) {
	return &NilExpression{}, nil
}

type NilExpression struct{}

func (e *NilExpression) Type() ExpressionType {
	return NilExpr
}

func (e *NilExpression) Value() string {
	return "Nil"
}

// Number Expression
func createNumberExpression(numerator, denominator int64) (*NumberExpression, error) {
	return &NumberExpression{
		Numerator:   numerator,
		Denominator: denominator,
	}, nil
}

type NumberExpression struct {
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

// List Expression
func createListExpression(expressions ...Expression) (*ListExpression, error) {
	return &ListExpression{
		Expressions: expressions,
	}, nil
}

type ListExpression struct {
	Expressions []Expression
}

func (e *ListExpression) Type() ExpressionType {
	return ListExpr
}

func (e *ListExpression) Value() string {
	return "List"
}

// Boolean Expression
func createBooleanExpression(value bool) (*BooleanExpression, error) {
	return &BooleanExpression{
		Val: value,
	}, nil
}

type BooleanExpression struct {
	Val bool
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

// Symbol Expression
func createSymbolExpression(identifier string, expression Expression) (*SymbolExpression, error) {
	return &SymbolExpression{
		Identifier: identifier,
		Expression: expression,
	}, nil
}

type SymbolExpression struct {
	Identifier string
	Expression Expression
}

func (e *SymbolExpression) Type() ExpressionType {
	return SymbolExpr
}

func (e *SymbolExpression) Value() string {
	return e.Identifier
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

	case parser.Boolean:
		if node.Identifier == "True" {
			return createBooleanExpression(true)
		}

		return createBooleanExpression(false)

	case parser.Number:
		numberStr := strings.Split(node.Identifier, "/")

		numerator := int64(0)
		denominator := int64(1)

		var err error
		if len(numberStr) > 1 {
			denominator, err = strconv.ParseInt(numberStr[1], 10, 64)
			if err != nil {
				return nil, err
			}
		}

		numerator, err = strconv.ParseInt(numberStr[0], 10, 64)
		if err != nil {
			return nil, err
		}

		return createNumberExpression(numerator, denominator)
	}

	return createListExpression(expressions...)
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
