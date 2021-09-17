package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/danfragoso/milho"
)

func initREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Milho REPL 🌽 v.%s\n", milho.Version())
	fmt.Printf("Danilo Fragoso <danilo.fragoso@gmail.com> - 2021\n\n")
	fmt.Printf("Type (list) to see all the available definitions!\n")

	prompt()

	sess := milho.CreateSession()
	cmdBuffer := ""

	for scanner.Scan() {
		cmdBuffer += scanner.Text()
		if validateBuffer(cmdBuffer) {
			results := milho.RunSession(cmdBuffer, sess)
			cmdBuffer = ""

			for _, result := range strings.Split(results, "\n") {
				r := strings.TrimSpace(result)
				if r != "" {
					fmt.Print("🍿 " + r + "\n")
				}
			}

			prompt()
		}
	}

	if scanner.Err() != nil {
		fmt.Printf("\n\nIO Err: %s", scanner.Err())
		os.Exit(1)
	}
}

func validateBuffer(buffer string) bool {
	if strings.TrimSpace(buffer) == "" {
		return false
	}

	var ODelimiter int64
	var CDelimiter int64

	for _, r := range buffer {
		switch r {
		case '(':
			ODelimiter++

		case ')':
			CDelimiter++
		}
	}

	return ODelimiter <= CDelimiter
}

func prompt() {
	fmt.Printf("🌽 > ")
}
