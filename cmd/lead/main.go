// Command lead is the single-binary orchestrator CLI (scaffold for issue #35).
package main

import (
	"fmt"
	"os"
)

// Version info is injected at build time via -ldflags "-X main.Version=...".
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func versionString() string {
	return fmt.Sprintf("lead version %s (commit: %s, built: %s)", Version, Commit, Date)
}

func usage() string {
	return "Usage: lead [--version|-v] [--help|-h] <command>\n\nCommands:\n  version     Print version information\n"
}

func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usage())
		return 2
	}
	switch args[0] {
	case "version", "--version", "-v":
		fmt.Println(versionString())
		return 0
	case "--help", "-h", "help":
		fmt.Print(usage())
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n%s", args[0], usage())
		return 1
	}
}

func main() {
	os.Exit(run(os.Args[1:]))
}
