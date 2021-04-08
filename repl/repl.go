package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/danfragoso/milho"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Milho REPL 🌽 v.%s\n", milho.Version())
	fmt.Printf("© Danilo Fragoso <danilo.fragoso@gmail.com> - 2021\n")

	prompt()

	sess := milho.CreateSession()
	for scanner.Scan() {
		cmd := scanner.Text()
		if strings.TrimSpace(cmd) != "" {
			results := milho.RunSession(cmd, sess)

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

func prompt() {
	fmt.Printf("🌽 > ")
}
