package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danfragoso/milho"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] != "" {
		file, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		runFile(string(file))
	} else {
		initREPL()
	}
}

func runFile(file string) {
	_, e := milho.RunRaw(file)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
