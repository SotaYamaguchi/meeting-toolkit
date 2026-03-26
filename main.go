// Package main is the entry point for the mtg CLI.
package main

import (
	"os"

	"github.com/syamaguc/meeting-toolkit/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
