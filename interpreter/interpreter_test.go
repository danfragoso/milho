package interpreter

import (
	"fmt"
	"testing"

	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

func Test_nil(t *testing.T) {
	src := "()"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	printExpr(expressions[0])

	if err != nil {
		t.Error(err)
	}

	if expressions[0].Type() != NilExpr {
		t.Error(err)
	}
}

func Test_number_slash(t *testing.T) {
	src := "4/3"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	printExpr(expressions[0])

	if err != nil {
		t.Error(err)
	}

	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	if expressions[0].Value() != "4/3" {
		t.Error(err)
	}
}

func Test_number_slash2(t *testing.T) {
	src := "10/5"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	printExpr(expressions[0])

	if err != nil {
		t.Error(err)
	}

	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	if expressions[0].Value() != "2" {
		t.Error(err)
	}
}

func Test_add(t *testing.T) {
	src := "(+ 1 2 (+ 1 2))"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	number := expressions[0].(*NumberExpression)

	if number.Numerator != 6 {
		t.Errorf("Wrong expression numerator value, expected 6 got %d", number.Numerator)
	}

	if number.Denominator != 1 {
		t.Errorf("Wrong expression denominator value, expected 1 got %d", number.Denominator)
	}

}

func Test_sub(t *testing.T) {
	src := "(+ (+ 1 2) (- 3) (- 3))"
	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	number := expressions[0].(*NumberExpression)

	if number.Numerator != -3 {
		t.Errorf("Wrong expression numerator value, expected -3 got %d", number.Numerator)
	}

	if number.Denominator != 1 {
		t.Errorf("Wrong expression denominator value, expected 1 got %d", number.Denominator)
	}

	src = "(- 3)"
	fmt.Println(src)
	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err = Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	number = expressions[0].(*NumberExpression)

	if number.Numerator != -3 {
		t.Errorf("Wrong expression numerator value, expected -3 got %d", number.Numerator)
	}

	if number.Denominator != 1 {
		t.Errorf("Wrong expression denominator value, expected 1 got %d", number.Denominator)
	}

	src = "(- 10 3)"
	fmt.Println(src)
	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err = Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	number = expressions[0].(*NumberExpression)

	if number.Numerator != 7 {
		t.Errorf("Wrong expression numerator value, expected 7 got %d", number.Numerator)
	}

	if number.Denominator != 1 {
		t.Errorf("Wrong expression denominator value, expected 1 got %d", number.Denominator)
	}
}

func Test_mul(t *testing.T) {
	src := "(* 0)"
	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	number := expressions[0].(*NumberExpression)

	if number.Numerator != 0 {
		t.Errorf("Wrong expression numerator value, expected 0 got %d", number.Numerator)
	}

	if number.Denominator != 1 {
		t.Errorf("Wrong expression denominator value, expected 1 got %d", number.Denominator)
	}

	src = "(* 10 5)"
	fmt.Println(src)
	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err = Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	number = expressions[0].(*NumberExpression)

	if number.Numerator != 50 {
		t.Errorf("Wrong expression numerator value, expected -3 got %d", number.Numerator)
	}

	if number.Denominator != 1 {
		t.Errorf("Wrong expression denominator value, expected 1 got %d", number.Denominator)
	}
}

func Test_div(t *testing.T) {
	src := "(/ 1 0)"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	_, err = Run(ast)
	fmt.Print(err, "\n\n")
	if err == nil {
		t.Error("Expected divide by zero error, got nothing")
	}

	src = "(/ 20 2)"
	fmt.Println(src)

	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	number := expressions[0].(*NumberExpression)

	if number.Numerator != 10 {
		t.Errorf("Wrong expression numerator value, expected 10 got %d", number.Numerator)
	}

	if number.Denominator != 1 {
		t.Errorf("Wrong expression denominator value, expected 1 got %d", number.Denominator)
	}
}

func Test_cmp(t *testing.T) {
	src := "(= 20 2)"
	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != BooleanExpr {
		t.Error(err)
	}

	r := expressions[0].(*BooleanExpression)

	if r.Val != false {
		t.Errorf("Wrong expression value, expected false got true")
	}

	src = "(= 20 20 20 ddefn)"
	fmt.Println(src)
	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	_, err = Run(ast)
	fmt.Print(err, "\n\n")
	if err == nil {
		t.Error("Expected error")
	}

	src = "(= True (= 2 2))"
	fmt.Println(src)
	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err = Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != BooleanExpr {
		t.Error(err)
	}

	r = expressions[0].(*BooleanExpression)

	if r.Val != true {
		t.Errorf("Wrong expression value, expected true got false")
	}

	src = "(= 20 20 20 vinte)"
	fmt.Println(src)
	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	_, err = Run(ast)
	fmt.Print(err, "\n\n")
	if err == nil {
		t.Error("Expected error")
	}

	src = "(if (= 2 2) 10 2)"
	fmt.Println(src)
	tokens, err = tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err = Run(ast)
	if err != nil {
		t.Error(err)
	}

	printExpr(expressions[0])
	if expressions[0].Type() != NumberExpr {
		t.Error(err)
	}

	n := expressions[0].(*NumberExpression)

	if n.Numerator != 10 {
		t.Errorf("Wrong expression value, expected 10 got %d", n.Numerator)
	}
}

func Test_def(t *testing.T) {
	src := "(def acc 10)\n"
	src += "(+ acc 1)"

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	expectedExpressions := []Expression{
		&SymbolExpression{}, &NumberExpression{},
	}

	for i, expression := range expressions {
		printExpr(expressions[i])

		if expectedExpressions[i].Type() != expression.Type() {
			t.Errorf("Wrong result type found, expected %s got %s", expectedExpressions[i].Type().String(), expression.Type().String())
		}
	}
}

func Test_let(t *testing.T) {
	src := "(let (a 2) (* a a))\n"
	src += "(let (a 2) (let (b (* a a)) (* b a)))\n"

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	expectedExpressions := []Expression{
		&NumberExpression{
			Numerator:   4,
			Denominator: 1,
		},
		&NumberExpression{
			Numerator:   8,
			Denominator: 1,
		},
	}

	for i, expression := range expressions {
		printExpr(expressions[i])

		if expectedExpressions[i].Type() != expression.Type() {
			t.Errorf("Wrong result type found, expected %s got %s", expectedExpressions[i].Type().String(), expression.Type().String())
		}

		if expectedExpressions[i].Value() != expression.Value() {
			t.Errorf("Wrong result value found, expected %s got %s", expectedExpressions[i].Value(), expression.Value())
		}
	}
}

func Test_string(t *testing.T) {
	src := "(def lang \"milho\")\n"
	src += "(def food (str lang \" cozido na agua\")\n"
	src += "(prn food)"

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	expectedExpressions := []Expression{
		&SymbolExpression{}, &SymbolExpression{}, &StringExpression{},
	}

	for i, expression := range expressions {
		printExpr(expressions[i])

		if expectedExpressions[i].Type() != expression.Type() {
			t.Errorf("Wrong result type found, expected %s got %s", expectedExpressions[i].Type().String(), expression.Type().String())
		}
	}
}

func Test_fn(t *testing.T) {
	src := "(fn (a) (* a a))\n"

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	expectedExpressions := []Expression{
		&FunctionExpression{},
	}

	for i, expression := range expressions {
		printExpr(expressions[i])

		if expectedExpressions[i].Type() != expression.Type() {
			t.Errorf("Wrong result type found, expected %s got %s", expectedExpressions[i].Type().String(), expression.Type().String())
		}
	}
}

func Test_defn(t *testing.T) {
	src := "(defn square (a) (* a a))\n"

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expressions, err := Run(ast)
	if err != nil {
		t.Error(err)
	}

	expectedExpressions := []Expression{
		&FunctionExpression{
			Identifier: "square",
		},
	}

	for i, expression := range expressions {
		printExpr(expressions[i])

		if expectedExpressions[i].Type() != expression.Type() {
			t.Errorf("Wrong result type found, expected %s got %s", expectedExpressions[i].Type().String(), expression.Type().String())
		}

		if expectedExpressions[i].Value() != expression.Value() {
			t.Errorf("Wrong result type found, expected %s got %s", expectedExpressions[i].Value(), expression.Value())
		}
	}
}
