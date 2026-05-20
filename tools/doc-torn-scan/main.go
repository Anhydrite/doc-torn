package main

import (
	"fmt"
	"os"

	"github.com/Anhydrite/doc-torn/tools/doc-torn-scan/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	var err error
	switch subcommand {
	case "tree":
		err = cmd.RunTree()
	case "scaffold":
		err = cmd.RunScaffold(args)
	case "complete":
		err = cmd.RunComplete(args)
	case "meta":
		err = cmd.RunMeta()
	case "status":
		err = cmd.RunStatus()
	case "help", "--help", "-h":
		printUsage()
		return
	default:
		err = fmt.Errorf("unknown subcommand: %s", subcommand)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`doc-torn-scan — Iterative documentation scanner

Usage:
  doc-torn-scan tree              List all project files as JSON
  doc-torn-scan scaffold <name>    Generate markdown skeleton for a feature
  doc-torn-scan complete <name>    Mark a feature as completed
  doc-torn-scan meta               Generate global documentation (L0, arch, matrix, etc.)
  doc-torn-scan status             Show documentation progress`)
}
