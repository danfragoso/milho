package parser

import (
	"testing"

	"github.com/danfragoso/milho/tokenizer"
)

func Test_number(t *testing.T) {
	tokens, err := tokenizer.Tokenize("2")
	if err != nil {
		t.Error(err)
	}

	tree, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	if tree.Type != Number {
		t.Errorf("Expected node type to be Number, got %s", tree.Type)
	}

	if tree.Identifier != "2" {
		t.Errorf("Expected node identifier to be 2, got %s", tree.Identifier)
	}
}

func Test_parens(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(+ 2 (+ 1 3))")
	if err != nil {
		t.Error(err)
	}

	tree, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	if tree.Type != Function {
		t.Errorf("Expected function node type, got %s", tree.Type)
	}

	if tree.Identifier != "+" {
		t.Errorf("Expected function to be +, got %s", tree.Identifier)
	}

	if tree.Parent != nil {
		t.Errorf("Expected tree parent to be nil, got %s", tree)
	}

	if len(tree.Nodes) != 2 {
		t.Errorf("Expected tree to have two child nodes, got %s", tree.Nodes)
	}
}
