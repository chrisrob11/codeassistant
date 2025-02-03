// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.
package cmd

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

// EndSessionCommand archives the session
func EndSessionCommand() *cli.Command {
	return &cli.Command{
		Name:  "end-session",
		Usage: "Archive session to historical storage",
		Action: func(c *cli.Context) error {
			fmt.Println("Ending session...")
			return nil
		},
	}
}
