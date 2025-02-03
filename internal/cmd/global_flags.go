// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import cli "github.com/urfave/cli/v2"

// GlobalFlags defines global CLI flags.
func GlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "llm-model",
			Usage:   "Set the default LLM model",
			EnvVars: []string{"CA_LLM_MODEL"},
		},
		&cli.Float64Flag{
			Name:    "llm-temperature",
			Usage:   "Set the LLM temperature",
			EnvVars: []string{"CA_LLM_TEMPERATURE"},
		},
		&cli.BoolFlag{
			Name:    "store-summary",
			Usage:   "Enable or disable storing summaries in session",
			EnvVars: []string{"CA_STORE_SUMMARY"},
		},
	}
}
