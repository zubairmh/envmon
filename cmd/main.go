package main

import (
	"fmt"
	"os"
	"envmon/internal/cli"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: envmon <command>")
		os.Exit(1)
	}

	switch args[0] {
	case "help", "-h", "--help":
		cli.ShowHelp()
	case "configs":
		fmt.Println("TODO")
	default:
		fmt.Println(args[0])
	}
}
