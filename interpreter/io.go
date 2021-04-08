package interpreter

import "fmt"

func pr(params []Result) (Result, error) {
	for _, param := range params {
		fmt.Print(param.Value())
	}

	return &NilResult{}, nil
}

func prn(params []Result) (Result, error) {
	for _, param := range params {
		fmt.Print(param.Value())
	}

	fmt.Println("")
	return &NilResult{}, nil
}
