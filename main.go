package main

import (
	"fmt"
	"os"
)

const (
	exitOK    = 0
	exitErr   = 1
	exitUsage = 2
	exitFatal = 3
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printHelp()
		return exitUsage
	}

	switch args[0] {
	case "serve":
		if err := cmdServe(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "serve: %v\n", err)
			return exitErr
		}
		return exitOK
	case "generate":
		if err := cmdGenerate(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "generate: %v\n", err)
			return exitErr
		}
		return exitOK
	case "help", "-h", "--help":
		printHelp()
		return exitOK
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printHelp()
		return exitUsage
	}
}

// mainWithoutExit returns exit code without calling os.Exit (for tests).
func mainWithoutExit() int {
	return run(os.Args[1:])
}

func printHelp() {
	fmt.Print(`ballo-cli — scaffold Ballo applications

Usage:
  ballo-cli <command> [arguments]

Commands:
  serve [addr]         Start a Ballo development server (default :8080)
  generate <template>  Generate a starter template
  help                 Show this help

Exit codes:
  0  success
  1  runtime error
  2  invalid usage
  3  fatal error
`)
}
