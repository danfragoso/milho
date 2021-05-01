package parser

import (
	"fmt"
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
	src := "2"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	if tree.Type != Number {
		t.Errorf("Expected node type to be Number, got %s", tree.Type)
	}

	if tree.Identifier != "2" {
		t.Errorf("Expected node identifier to be 2, got %s", tree.Identifier)
	}
}

func Test_number_slash(t *testing.T) {
	src := "4/3"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	if tree.Type != Number {
		t.Errorf("Expected node type to be Number, got %s", tree.Type)
	}

	if tree.Identifier != "4/3" {
		t.Errorf("Expected node identifier to be 4/3, got %s", tree.Identifier)
	}
}

func Test_number_slash2(t *testing.T) {
	src := "10/5"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	if tree.Type != Number {
		t.Errorf("Expected node type to be Number, got %s", tree.Type)
	}

	if tree.Identifier != "10/5" {
		t.Errorf("Expected node identifier to be 10/5, got %s", tree.Identifier)
	}
}

func Test_parens(t *testing.T) {
	src := "(+ 2 (+ 1 3))"
	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	if tree.Type != List {
		t.Errorf("Expected List node type, got %s", tree.Type)
	}

	if tree.Identifier != "" {
		t.Errorf("Expected List to be +, got %s", tree.Identifier)
	}

	if tree.Parent != nil {
		t.Errorf("Expected tree parent to be nil, got %s", tree)
	}

	if len(tree.Nodes) != 3 {
		t.Errorf("Expected tree to have 3 child nodes, got %s", tree.Nodes)
	}
}

func Test_list(t *testing.T) {
	src := "(2 3 4)"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

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
	src := "(append '(2) ''(3 4))"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	expectedNodes := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "append"}, {Type: List},
		{Type: Identifier, Identifier: "quote"}, {Type: List},
		{Type: Number, Identifier: "2"}, {Type: List},
		{Type: Identifier, Identifier: "quote"}, {Type: List},
		{Type: Identifier, Identifier: "quote"}, {Type: List},
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

func Test_defn(t *testing.T) {
	src := "(defn sum (a b) (+ a b))"
	fmt.Println(src)

	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	expectedNodes := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "defn"}, {Type: Identifier, Identifier: "sum"},
		{Type: List}, {Type: Identifier, Identifier: "a"}, {Type: Identifier, Identifier: "b"},
		{Type: List}, {Type: Identifier, Identifier: "+"}, {Type: Identifier, Identifier: "a"},
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
	src := `
	(defn inc-if-gt-zero (a)
		(if (> a 0)
			(+ a 1) (+ a 0)))`

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	expectedNodes := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "defn"}, {Type: Identifier, Identifier: "inc-if-gt-zero"}, {Type: List},
		{Type: Identifier, Identifier: "a"}, {Type: List}, {Type: Identifier, Identifier: "if"},
		{Type: List}, {Type: Identifier, Identifier: ">"}, {Type: Identifier, Identifier: "a"},
		{Type: Number, Identifier: "0"}, {Type: List}, {Type: Identifier, Identifier: "+"},
		{Type: Identifier, Identifier: "a"}, {Type: Number, Identifier: "1"},
		{Type: List}, {Type: Identifier, Identifier: "+"}, {Type: Identifier, Identifier: "a"},
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

func Test_boolean(t *testing.T) {
	src := `
	(defn dumb-negate (a)
		(if (= a True) False True))`

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	if err != nil {
		t.Error(err)
	}

	tree := nodes[0]
	fmt.Println(tree)

	expectedNodes := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "defn"}, {Type: Identifier, Identifier: "dumb-negate"}, {Type: List},
		{Type: Identifier, Identifier: "a"}, {Type: List}, {Type: Identifier, Identifier: "if"},
		{Type: List}, {Type: Identifier, Identifier: "="}, {Type: Identifier, Identifier: "a"},
		{Type: Boolean, Identifier: "True"}, {Type: Boolean, Identifier: "False"},
		{Type: Boolean, Identifier: "True"},
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

func Test_session_def(t *testing.T) {
	src := `
	(def numb 1000)
	(* 2 numb)`

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	fmt.Println(nodes)
	if err != nil {
		t.Error(err)
	}

	expectedNodes1 := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "def"}, {Type: Identifier, Identifier: "numb"},
		{Type: Number, Identifier: "1000"},
	}

	for idx, node := range treeAsList(nodes[0]) {
		if expectedNodes1[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes1[idx].Type, node.Type)
		}

		if expectedNodes1[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes1[idx].Identifier, node.Identifier)
		}
	}

	expectedNodes2 := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "*"}, {Type: Number, Identifier: "2"},
		{Type: Identifier, Identifier: "numb"},
	}

	for idx, node := range treeAsList(nodes[1]) {
		if expectedNodes2[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes2[idx].Type, node.Type)
		}

		if expectedNodes2[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes2[idx].Identifier, node.Identifier)
		}
	}
}

func Test_string(t *testing.T) {
	src := `
	(def lang "milho")
	(def food (str lang " cozido na agua"))
	(prn food)`

	fmt.Println(src)
	tokens, err := tokenizer.Tokenize(src)

	if err != nil {
		t.Error(err)
	}

	nodes, err := Parse(tokens)
	fmt.Println(nodes)
	if err != nil {
		t.Error(err)
	}

	expectedNodes1 := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "def"}, {Type: Identifier, Identifier: "lang"},
		{Type: String, Identifier: "milho"},
	}

	for idx, node := range treeAsList(nodes[0]) {
		if expectedNodes1[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes1[idx].Type, node.Type)
		}

		if expectedNodes1[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes1[idx].Identifier, node.Identifier)
		}
	}

	expectedNodes2 := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "def"}, {Type: Identifier, Identifier: "food"},
		{Type: List}, {Type: Identifier, Identifier: "str"}, {Type: Identifier, Identifier: "lang"},
		{Type: String, Identifier: " cozido na agua"},
	}

	for idx, node := range treeAsList(nodes[1]) {
		if expectedNodes2[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes2[idx].Type, node.Type)
		}

		if expectedNodes2[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes2[idx].Identifier, node.Identifier)
		}
	}

	expectedNodes3 := []*Node{
		{Type: List}, {Type: Identifier, Identifier: "prn"}, {Type: Identifier, Identifier: "food"},
	}

	for idx, node := range treeAsList(nodes[2]) {
		if expectedNodes3[idx].Type != node.Type {
			t.Errorf("Expected node %d type to be %s, got %s", idx, expectedNodes3[idx].Type, node.Type)
		}

		if expectedNodes3[idx].Identifier != node.Identifier {
			t.Errorf("Expected node %d identifier to be '%s', got '%s'", idx, expectedNodes3[idx].Identifier, node.Identifier)
		}
	}
}
