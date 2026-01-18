package main

import (
	"fmt"
	"os"
)


func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: envmon <command>")
		os.Exit(1)
	}

	switch args[0] {
	case "help":
		fmt.Println("TODO")
	}
}