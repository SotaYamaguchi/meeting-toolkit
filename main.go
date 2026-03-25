// Package main is the entry point for the mtg CLI.
package main

import (
	"fmt"
	"os"

	"github.com/syamaguc/meeting-toolkit/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]

	var err error
	switch subcommand {
	case "prep":
		err = cmd.RunPrep(os.Args[2:])
	case "memo":
		err = cmd.RunMemo(os.Args[2:])
	case "mail":
		err = cmd.RunMail(os.Args[2:])
	case "list", "-list", "--list":
		err = cmd.RunList()
	case "help", "-h", "--help":
		cmd.PrintUsage()
	case "completion":
		err = cmd.RunCompletion(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "エラー: 不明なサブコマンド '%s'\n\n", subcommand)
		cmd.PrintUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
