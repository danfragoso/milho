package parser

import (
	"testing"

	"github.com/danfragoso/milho/tokenizer"
)

func treeAsList(self *Node) []*Node {
	nodeList := []*Node{self}
	for _, child := range self.Nodes {
		nodeList = append(nodeList, treeAsList(child)...)
	}

	return nodeList
}

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

func Test_list(t *testing.T) {
	tokens, err := tokenizer.Tokenize("[2 3 4]")
	if err != nil {
		t.Error(err)
	}

	tree, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expectedNodes := []*Node{
		{Type: List}, {Type: Number, Identifier: "2"},
		{Type: Number, Identifier: "3"}, {Type: Number, Identifier: "4"},
	}

	for idx, node := range treeAsList(tree) {
		if expectedNodes[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes[idx].Type, node.Type)
		}

		if expectedNodes[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes[idx].Identifier, node.Identifier)
		}
	}
}

func Test_list_param(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(++ [2] [3 4]")
	if err != nil {
		t.Error(err)
	}

	tree, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expectedNodes := []*Node{
		{Type: Function, Identifier: "++"}, {Type: List}, {Type: Number, Identifier: "2"},
		{Type: List}, {Type: Number, Identifier: "3"}, {Type: Number, Identifier: "4"},
	}

	for idx, node := range treeAsList(tree) {
		if expectedNodes[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes[idx].Type, node.Type)
		}

		if expectedNodes[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes[idx].Identifier, node.Identifier)
		}
	}
}

func Test_defn(t *testing.T) {
	tokens, err := tokenizer.Tokenize("(defn sum [a b] (+ a b)")
	if err != nil {
		t.Error(err)
	}

	tree, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expectedNodes := []*Node{
		{Type: Macro, Identifier: "defn"}, {Type: Identifier, Identifier: "sum"}, {Type: List},
		{Type: Identifier, Identifier: "a"}, {Type: Identifier, Identifier: "b"},
		{Type: Function, Identifier: "+"}, {Type: Identifier, Identifier: "a"},
		{Type: Identifier, Identifier: "b"},
	}

	for idx, node := range treeAsList(tree) {
		if expectedNodes[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes[idx].Type, node.Type)
		}

		if expectedNodes[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes[idx].Identifier, node.Identifier)
		}
	}
}

func Test_nested(t *testing.T) {
	tokens, err := tokenizer.Tokenize(`
		(defn inc-if-gt-zero [a] 
			(if (> a 0) 
				(+ a 1) (+ a 0)))
	`)

	if err != nil {
		t.Error(err)
	}

	tree, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	expectedNodes := []*Node{
		{Type: Macro, Identifier: "defn"}, {Type: Identifier, Identifier: "inc-if-gt-zero"}, {Type: List},
		{Type: Identifier, Identifier: "a"}, {Type: Function, Identifier: "if"},
		{Type: Function, Identifier: ">"}, {Type: Identifier, Identifier: "a"},
		{Type: Number, Identifier: "0"}, {Type: Function, Identifier: "+"},
		{Type: Identifier, Identifier: "a"}, {Type: Number, Identifier: "1"},
		{Type: Function, Identifier: "+"}, {Type: Identifier, Identifier: "a"},
		{Type: Number, Identifier: "0"},
	}

	for idx, node := range treeAsList(tree) {
		if expectedNodes[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes[idx].Type, node.Type)
		}

		if expectedNodes[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes[idx].Identifier, node.Identifier)
		}
	}
}
