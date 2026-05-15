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
		since := ""
		aiOnly := false
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--ai-only" {
				aiOnly = true
			} else if os.Args[i] == "--since" && i+1 < len(os.Args) {
				since = os.Args[i+1]
				i++
			}
		}
		cli.List(aiOnly, since)
	case "export":
		cli.Export("json", "-")
	case "export-csv":
		cli.Export("csv", "-")
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
