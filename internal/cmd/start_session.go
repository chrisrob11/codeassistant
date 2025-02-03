// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import (
	"os"

	"github.com/chrisrob11/codeassistant/internal/session"
	cli "github.com/urfave/cli/v2"
)

// NewSessionCommand initializes a new coding session.
func NewSessionCommand() *cli.Command {
	return &cli.Command{
		Name:    "start-session",
		Aliases: []string{"ss"},
		Usage:   "Start a new coding session",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Name of the session to start",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			currentDir, err := os.Getwd()
			if err != nil {
				return err
			}

			return session.StartSession(&session.StartSessionRequest{Dir: currentDir, Name: c.String("name")})
		},
	}
}
