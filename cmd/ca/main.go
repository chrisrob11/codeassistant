// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

// Package main is the main package for the production of the ca binary
package main

import (
	"fmt"
	"os"

	"github.com/chrisrob11/codeassistant/internal/cmd"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ca",
		Usage: "AI-powered coding assistant",
		Flags: cmd.GlobalFlags(),
		Commands: []*cli.Command{
			cmd.NewSessionCommand(),
			cmd.CodeCommand(),
			cmd.ReviewCommand(),
			cmd.RollbackCommand(),
			cmd.EndSessionCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}
}
