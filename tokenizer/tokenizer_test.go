package tokenizer

import (
	"strings"
	"testing"
)

func Test_number(t *testing.T) {
	tokens, err := Tokenize("4")
	if err != nil {
		t.Error(err)
	} else {
		if tokens[0].Type != Number {
			t.Error("Wrong token type, expected Number got", tokens[0].Type)
		}

		if tokens[0].Value != "4" {
			t.Error("Wrong token value, expected 4 got", tokens[0].Value)
		}
	}
}

func Test_symbols(t *testing.T) {
	atoms := []string{"4c", "4f", "cd"}
	tokens, err := Tokenize(strings.Join(atoms, " "))

	if err != nil {
		t.Error(err)
	} else {
		for i, tok := range tokens {
			if tok.Type != Symbol {
				t.Error("Wrong token type, expected Symbol got", tok.Type)
			}

			if tok.Value != atoms[i] {
				t.Errorf("Wrong token value, expected %s got %s", atoms[i], tok.Value)
			}
		}
	}
}

func Test_parens(t *testing.T) {
	tokens, err := Tokenize("(+ 50   5       )")

	if err != nil {
		t.Error(err)
	} else {
		if tokens[0].Type != OParen {
			t.Error("Wrong token type, expected OParen got", tokens[0].Type)
		}

		if tokens[1].Type != Symbol {
			t.Error("Wrong token type, expected Symbol got", tokens[1].Type)
		}

		if tokens[2].Type != Number {
			t.Error("Wrong token type, expected Number got", tokens[2].Type)
		}

		if tokens[3].Type != Number {
			t.Error("Wrong token type, expected Number got", tokens[3].Type)
		}

		if tokens[4].Type != CParen {
			t.Error("Wrong token type, expected Number got", tokens[4].Type)
		}
	}
}

func Test_brackets(t *testing.T) {
	tokens, err := Tokenize("(defn sum [a b] (+ a b))")

	if err != nil {
		t.Error(err)
	} else {
		expectedTokens := []TokenType{
			OParen, Symbol, Symbol, OBrack, Symbol, Symbol,
			CBrack, OParen, Symbol, Symbol, Symbol, CParen,
			CParen,
		}

		for idx, tok := range expectedTokens {
			if tokens[idx].Type != tok {
				t.Errorf("Wrong token type, expected %s got %s", tok, tokens[idx].Type)
			}
		}
	}
}
