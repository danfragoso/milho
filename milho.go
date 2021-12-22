package milho

import (
	"fmt"
	"runtime"

	"github.com/danfragoso/milho/compiler"
	"github.com/danfragoso/milho/interpreter"
	"github.com/danfragoso/milho/mir"
	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

var version string

func Version() string {
	return fmt.Sprintf("%s_%s:%s", version, runtime.GOARCH, runtime.GOOS)
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
		ret += fmt.Sprintf("\nEvaluation error: \n%s", err)
	}

	return ret
}

func RunSession(src string, sess *mir.Session) string {
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

func RunRaw(src string) ([]mir.Expression, error) {
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

func TranspileToJS(src string) (string, error) {
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		return "", fmt.Errorf("Tokenization error: %s\n", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return "", fmt.Errorf("Parsing error: %s\n", err)
	}

	var combinedSrc string
	for _, node := range ast {
		tree, err := mir.GenerateMIR(node)
		if err != nil {
			return "", err
		}

		combinedSrc += compiler.TranspileJS(tree) + "\n"
	}

	return combinedSrc, nil
}

func TranspileToLLVM(src string) (string, error) {
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		return "", fmt.Errorf("Tokenization error: %s\n", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return "", fmt.Errorf("Parsing error: %s\n", err)
	}

	var combinedSrc string
	for _, node := range ast {
		tree, err := mir.GenerateMIR(node)
		if err != nil {
			return "", err
		}

		combinedSrc += compiler.TranspileLLVM(tree) + "\n"
	}

	return combinedSrc, nil
}

func CreateSession() *mir.Session {
	sess, _ := interpreter.CreateSession(&parser.Node{})
	return sess
}
