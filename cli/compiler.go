package main

import (
	"fmt"

	"github.com/danfragoso/milho"
)

func compileMilho(milhoSrc string) {
	js, err := milho.TranspileToJS(milhoSrc)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(js)
}
