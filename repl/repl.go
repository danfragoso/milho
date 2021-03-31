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
	fmt.Println("Milho REPL 🌽")
	fmt.Printf("© Danilo Fragoso <danilo.fragoso@gmail.com> - 2021\n\n")

	prompt()

	for scanner.Scan() {
		cmd := scanner.Text()
		if strings.TrimSpace(cmd) != "" {
			fmt.Print("🍿 " + milho.Run(cmd))
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
