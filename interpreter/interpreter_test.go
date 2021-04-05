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

	results, err := Run(ast)
	res := results[0]

	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "9" {
			t.Errorf("Wrong response value, expected 9 got %s", res.Value())
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

	results, err := Run(ast)
	res := results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "-3" {
			t.Errorf("Wrong response value, expected -3 got %s", res.Value())
		}
	}

	tokens, err = tokenizer.Tokenize("(- 3)")
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	results, err = Run(ast)
	res = results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "-3" {
			t.Errorf("Wrong response value, expected -3 got %s", res.Value())
		}
	}

	tokens, err = tokenizer.Tokenize("(- 10 3)")
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	results, err = Run(ast)
	res = results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "7" {
			t.Errorf("Wrong response value, expected 7 got %s", res.Value())
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

	results, err := Run(ast)
	res := results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "0" {
			t.Errorf("Wrong response value, expected 0 got %s", res.Value())
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

	results, err = Run(ast)
	res = results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "500" {
			t.Errorf("Wrong response value, expected 500 got %s", res.Value())
		}
	}
}

func Test_div(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(/ 1 0)")
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	_, err = Run(ast)
	if err == nil {
		t.Error("Expected divide by zero error, got nothing")
	}

	tokens, err = tokenizer.Tokenize("(/ 20 2)")
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	results, err := Run(ast)
	res := results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Number {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "10" {
			t.Errorf("Wrong response value, expected 10 got %s", res.Value())
		}
	}
}

func Test_eq(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(= 20 2)")
	if err != nil {
		t.Error(err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	results, err := Run(ast)
	res := results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Boolean {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "False" {
			t.Errorf("Wrong response value, expected  got %s", res.Value())
		}
	}

	tokens, err = tokenizer.Tokenize("(= 20 20 20 defn)")
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	_, err = Run(ast)
	if err == nil {
		t.Error("Expected error")
	}

	tokens, err = tokenizer.Tokenize("(= True (= 2 2))")
	if err != nil {
		t.Error(err)
	}

	ast, err = parser.Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	results, err = Run(ast)
	res = results[0]
	if err != nil {
		t.Error("\n", err)
	} else {
		if res.Type() != Boolean {
			t.Errorf("Wrong response type, expected Number got %s", res.Type())
		}

		if res.Value() != "True" {
			t.Errorf("Wrong response value, expected  got %s", res.Value())
		}
	}

}
