package main

import (
	"fmt"
	"os"

	"envmon/internal/cli"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		if err := cli.ShowCurrentDeployment(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	switch args[0] {
	case "help", "-h", "--help":
		cli.ShowHelp()
	case "configs":
		if err := cli.ShowConfigs(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		if err := cli.SwitchDeployment(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
