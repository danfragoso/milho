package main

import (
	"syscall/js"

	"github.com/danfragoso/milho"
)

func eval(this js.Value, inputs []js.Value) interface{} {
	milhoSrc := inputs[0].String()
	return milho.Run(milhoSrc)
}

func main() {
	c := make(chan bool)
	js.Global().Set("eval", js.FuncOf(eval))
	<-c
}
