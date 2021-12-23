package main

import (
	"fmt"
	"syscall/js"

	"github.com/danfragoso/milho"
	"github.com/danfragoso/milho/mir"
	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

var session *mir.Session

func createREPLSession(this js.Value, inputs []js.Value) interface{} {
	session = milho.CreateSession()
	return true
}

func eval(this js.Value, inputs []js.Value) interface{} {
	milhoSrc := inputs[0].String()
	return milho.RunSession(milhoSrc, session)
}

func ast(this js.Value, inputs []js.Value) interface{} {
	milhoSrc := inputs[0].String()
	tokens, err := tokenizer.Tokenize(milhoSrc)
	if err != nil {
		return fmt.Sprintf("+- Nil:[%s]\n", err)
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Sprintf("+- Nil:[%s]\n", err)
	}

	resultAST := ""
	for _, node := range nodes {
		resultAST += node.String()
	}

	return resultAST
}

func showJS(this js.Value, inputs []js.Value) interface{} {
	milhoSrc := inputs[0].String()
	out, err := milho.TranspileToJS(milhoSrc)
	if err != nil {
		return fmt.Sprintf("Compilation Error [JavaScript Backend] @ [%s]\n", err)
	}

	return out
}

func showMir(this js.Value, inputs []js.Value) interface{} {
	milhoSrc := inputs[0].String()
	tokens, err := tokenizer.Tokenize(milhoSrc)
	if err != nil {
		return fmt.Sprintf("+- Nil:[%s]\n", err)
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Sprintf("+- Nil:[%s]\n", err)
	}

	resultMIR := ""
	for _, node := range nodes {
		ir, err := mir.GenerateMIR(node)
		if err != nil {
			return fmt.Sprintf("+- Nil:[%s]\n", err)
		}

		resultMIR += mir.SprintExpr(ir)
	}

	return resultMIR
}

func main() {
	c := make(chan bool)
	js.Global().Set("replVersion", milho.Version())
	js.Global().Set("createREPLSession", js.FuncOf(createREPLSession))

	js.Global().Set("evalMilho", js.FuncOf(eval))
	js.Global().Set("ast", js.FuncOf(ast))
	js.Global().Set("mir", js.FuncOf(showMir))
	js.Global().Set("js", js.FuncOf(showJS))
	<-c
}
