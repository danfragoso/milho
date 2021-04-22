package milho

import (
	"fmt"

	"github.com/danfragoso/milho/interpreter"
	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

var version string

func Version() string {
	return version
}

func Run(src string) string {
	var ret string
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		return fmt.Sprintf("Tokenization error: %s\n", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Sprintf("Parsing error: %s\n", err)
	}

	expressions, err := interpreter.Run(ast)
	for _, expr := range expressions {
		ret += "\n" + expr.Value()
	}

	if err != nil {
		ret += fmt.Sprintf("\nEvaluation error: %s", err)
	}

	return ret
}

func RunSession(src string, sess *interpreter.Session) string {
	var ret string
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		return fmt.Sprintf("Tokenization error: %s\n", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Sprintf("Parsing error: %s\n", err)
	}

	expressions, err := interpreter.RunFromSession(ast, sess)
	for _, expr := range expressions {
		ret += "\n" + expr.Value()
	}

	if err != nil {
		ret += fmt.Sprintf("\nEvaluation error: %s", err)
	}

	return ret
}

func RunRaw(src string) ([]interpreter.Expression, error) {
	sess := CreateSession()
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		return nil, fmt.Errorf("Tokenization error: %s\n", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %s\n", err)
	}

	return interpreter.RunFromSession(ast, sess)
}

func CreateSession() *interpreter.Session {
	return &interpreter.Session{}
}
