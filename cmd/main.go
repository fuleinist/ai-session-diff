package main

import (
	"fmt"
	"os"

	"github.com/fuleinist/ai-session-diff/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		cli.ShowHelp()
		os.Exit(0)
	}

	cmd := os.Args[1]
	switch cmd {
	case "install":
		cli.Install()
	case "uninstall":
		cli.Uninstall()
	case "show":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: session ID required")
			fmt.Fprintln(os.Stderr, "Usage: ai-session-diff show <session-id>")
			os.Exit(1)
		}
		cli.Show(os.Args[2])
	case "report":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: session ID required")
			fmt.Fprintln(os.Stderr, "Usage: ai-session-diff report <session-id>")
			os.Exit(1)
		}
		cli.Report(os.Args[2])
	case "list":
		cli.List(len(os.Args) > 2 && os.Args[2] == "--ai-only")
	case "status":
		cli.Status()
	case "help", "--help", "-h":
		cli.ShowHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		cli.ShowHelp()
		os.Exit(1)
	}
}