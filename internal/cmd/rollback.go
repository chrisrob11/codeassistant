// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

// RollbackCommand undoes a specific AI-modified step.
func RollbackCommand() *cli.Command {
	return &cli.Command{
		Name:  "rollback",
		Usage: "Undo a specific AI-modified step",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "step",
				Usage: "Specify the step to roll back",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Rolling back step:", c.Int("step"))
			return nil
		},
	}
}
