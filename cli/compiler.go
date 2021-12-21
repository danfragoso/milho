package main

import (
	"fmt"

	"github.com/danfragoso/milho"
)

func compileMilho(milhoSrc string, target string) error {
	switch target {
	case "JS":
		js, err := milho.TranspileToJS(milhoSrc)
		if err != nil {
			return err
		}

		fmt.Println(js)
		return nil

	case "LLVM":
		llvm, err := milho.TranspileToLLVM(milhoSrc)
		if err != nil {
			return err
		}

		fmt.Println(llvm)
		return nil
	}

	return fmt.Errorf("unknown target: %s", target)
}
