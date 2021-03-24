package interpreter

import (
	"testing"

	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

func Test_run(t *testing.T) {
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
