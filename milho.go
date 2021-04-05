package milho

import (
	"fmt"

	"github.com/danfragoso/milho/interpreter"
	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

func Run(src string) string {
	tokens, err := tokenizer.Tokenize(src)
	if err != nil {
		return fmt.Sprintf("Tokenization error: %s\n", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Sprintf("Parsing error: %s\n", err)
	}

	res, err := interpreter.Run(ast)
	if err != nil {
		return fmt.Sprintf("Evaluation error: %s\n", err)
	} else {
		strResult := "Nil"
		if res.Type() != interpreter.Nil {
			strResult = res.Value()
		}

		return fmt.Sprintln(strResult)
	}
}
