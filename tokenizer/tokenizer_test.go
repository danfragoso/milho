package tokenizer

import (
	"fmt"
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

func Test_number_slash(t *testing.T) {
	tokens, err := Tokenize("4/3")
	if err != nil {
		t.Error(err)
	} else {
		if tokens[0].Type != Number {
			t.Error("Wrong token type, expected Number got", tokens[0].Type)
		}

		if tokens[0].Value != "4/3" {
			t.Error("Wrong token value, expected 4/3 got", tokens[0].Value)
		}
	}
}

func Test_number_slash2(t *testing.T) {
	tokens, err := Tokenize("10/5")
	if err != nil {
		t.Error(err)
	} else {
		if tokens[0].Type != Number {
			t.Error("Wrong token type, expected Number got", tokens[0].Type)
		}

		if tokens[0].Value != "10/5" {
			t.Error("Wrong token value, expected 10/5 got", tokens[0].Value)
		}
	}
}

func Test_byte(t *testing.T) {
	tokens, err := Tokenize("0xFF")
	if err != nil {
		t.Error(err)
	} else {
		if tokens[0].Type != Byte {
			t.Error("Wrong token type, expected Byte got", tokens[0].Type)
		}

		if tokens[0].Value != "0xFF" {
			t.Error("Wrong token value, expected 0xFF got", tokens[0].Value)
		}
	}
}

func Test_fakebyte(t *testing.T) {
	_, err := Tokenize("0Xf")
	if err == nil {
		t.Error("Expected tokenizer error for 0Xf!")
	}

	fmt.Println(err)
}

func Test_symbols(t *testing.T) {
	atoms := []string{"c4", "f4", "cd"}
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

func Test_list(t *testing.T) {
	tokens, err := Tokenize("(defn sum (a b) (+ a b))")

	if err != nil {
		t.Error(err)
	} else {
		expectedTokens := []TokenType{
			OParen, Symbol, Symbol, OParen, Symbol, Symbol,
			CParen, OParen, Symbol, Symbol, Symbol, CParen,
			CParen,
		}

		for idx, tok := range expectedTokens {
			if tokens[idx].Type != tok {
				t.Errorf("Wrong token type, expected %s got %s", tok, tokens[idx].Type)
			}
		}
	}
}

func Test_bool(t *testing.T) {
	tokens, err := Tokenize("(= True False)")

	if err != nil {
		t.Error(err)
	} else {
		expectedTokens := []TokenType{
			OParen, Symbol, Boolean, Boolean, CParen,
		}

		for idx, tok := range expectedTokens {
			if tokens[idx].Type != tok {
				t.Errorf("Wrong token type, expected %s got %s", tok, tokens[idx].Type)
			}
		}
	}
}

func Test_session_def(t *testing.T) {
	tokens, err := Tokenize(`
		(def numb 1000)
		(* 2 numb)
	`)

	if err != nil {
		t.Error(err)
	} else {
		expectedTokens := []TokenType{
			OParen, Symbol, Symbol, Number, CParen,
			OParen, Symbol, Number, Symbol, CParen,
		}

		for idx, tok := range expectedTokens {
			if tokens[idx].Type != tok {
				t.Errorf("Wrong token type, expected %s got %s", tok, tokens[idx].Type)
			}
		}
	}
}

func Test_session_fn(t *testing.T) {
	tokens, err := Tokenize(`
		(defn fib-nth (n)
			(if (< n 2) n
			(+ (fib-nth (- n 1)) (fib-nth (- n 2)))))

		(fib-nth 10)
	`)

	if err != nil {
		t.Error(err)
	} else {
		expectedTokens := []TokenType{
			OParen, Symbol, Symbol, OParen, Symbol,
			CParen, OParen, Symbol, OParen, Symbol,
			Symbol, Number, CParen, Symbol, OParen,
			Symbol, OParen, Symbol, OParen, Symbol,
			Symbol, Number, CParen, CParen, OParen,
			Symbol, OParen, Symbol, Symbol, Number,
			CParen, CParen, CParen, CParen, CParen,
			OParen, Symbol, Number, CParen,
		}

		for idx, tok := range expectedTokens {
			if tokens[idx].Type != tok {
				t.Errorf("Wrong token type, expected %s got %s", tok, tokens[idx].Type)
			}
		}
	}
}

func Test_string(t *testing.T) {
	tokens, err := Tokenize(`
		(def lang "milho")
		(def food (str lang " cozido na agua"))
		(prn food)
	`)

	if err != nil {
		t.Error(err)
	} else {
		expectedTokens := []TokenType{
			OParen, Symbol, Symbol, String, CParen,
			OParen, Symbol, Symbol, OParen, Symbol,
			Symbol, String, CParen, CParen, OParen,
			Symbol, Symbol, CParen,
		}

		for idx, tok := range expectedTokens {
			if tokens[idx].Type != tok {
				t.Errorf("Wrong token type, expected %s got %s", tok, tokens[idx].Type)
			}
		}
	}
}
