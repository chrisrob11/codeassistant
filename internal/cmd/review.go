// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

// ReviewCommand shows the session progress and diffs
func ReviewCommand() *cli.Command {
	return &cli.Command{
		Name:  "review",
		Usage: "Show session progress and diffs",
		Action: func(c *cli.Context) error {
			fmt.Println("Reviewing session changes...")
			return nil
		},
	}
}
