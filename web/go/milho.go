package main

import (
	"fmt"
	"syscall/js"

	"github.com/danfragoso/milho"
	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

func eval(this js.Value, inputs []js.Value) interface{} {
	milhoSrc := inputs[0].String()
	return milho.Run(milhoSrc)
}

func ast(this js.Value, inputs []js.Value) interface{} {
	milhoSrc := inputs[0].String()
	tokens, err := tokenizer.Tokenize(milhoSrc)
	if err != nil {
		return fmt.Sprintf("+- Nil:[%s]\n", err)
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Sprintf("+- Nil:[%s]\n", err)
	}

	return ast.String()
}

func main() {
	c := make(chan bool)
	js.Global().Set("eval", js.FuncOf(eval))
	js.Global().Set("ast", js.FuncOf(ast))
	<-c
}
