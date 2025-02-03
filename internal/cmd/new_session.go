// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

// NewSessionCommand initializes a new coding session.
func NewSessionCommand() *cli.Command {
	return &cli.Command{
		Name:  "new-session",
		Usage: "Start a new coding session",
		Action: func(c *cli.Context) error {
			fmt.Println("Starting a new session...")
			return nil
		},
	}
}
