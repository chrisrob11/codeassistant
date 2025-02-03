// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

// CodeCommand applies AI modifications to code
func CodeCommand() *cli.Command {
	return &cli.Command{
		Name:  "code",
		Usage: "Apply AI modifications to code",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "files",
				Usage: "Specify files to modify",
			},
			&cli.BoolFlag{
				Name:  "per-file",
				Usage: "Apply the prompt to each file individually",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Preview AI-generated changes without modifying files",
			},
			&cli.BoolFlag{
				Name:  "revise",
				Usage: "Modify the last step instead of creating a new one",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("revise") {
				fmt.Println("Revising the last AI-modified step...")
			} else {
				fmt.Println("Processing new AI code modification...")
			}
			return nil
		},
	}
}
