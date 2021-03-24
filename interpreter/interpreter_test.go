package interpreter

import (
	"testing"

	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

func Test_add(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(+ 1 2 (+ 3) (+ 3))")
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	res, err := Run(ast)
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type)
		}

		if res.Value != "9" {
			t.Errorf("Wrong response value, expected 9 got %s", res.Value)
		}
	}
}

func Test_sub(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(+ (+ 1 2) (- 3) (- 3))")
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	res, err := Run(ast)
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type)
		}

		if res.Value != "-3" {
			t.Errorf("Wrong response value, expected -3 got %s", res.Value)
		}
	}
}

func Test_mul(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(* 0)")
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	res, err := Run(ast)
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type)
		}

		if res.Value != "0" {
			t.Errorf("Wrong response value, expected 0 got %s", res.Value)
		}
	}

	tokens, err = tokenizer.Tokenize("(* 100 5)")
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	res, err = Run(ast)
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type)
		}

		if res.Value != "500" {
			t.Errorf("Wrong response value, expected 500 got %s", res.Value)
		}
	}
}
