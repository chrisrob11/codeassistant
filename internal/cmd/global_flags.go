// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import (
	"time"

	"github.com/teilomillet/gollm"
	cli "github.com/urfave/cli/v2"
)

// GlobalFlags defines global CLI flags.
func GlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "llm-provider",
			Value:   "openai",
			Usage:   "LLM provider to use (e.g. openai, ollama, anthropic)",
			EnvVars: []string{"CA_LLM_PROVIDER"},
		},
		&cli.StringFlag{
			Name:    "llm-model",
			Value:   "gpt-4",
			Usage:   "LLM model name (e.g. gpt-4, deepseek)",
			EnvVars: []string{"CA_LLM_MODEL"},
		},
		&cli.StringFlag{
			Name:    "llm-api-key",
			Value:   "",
			Usage:   "API key for the LLM provider (if required)",
			EnvVars: []string{"CA_LLM_API_KEY"},
		},
		&cli.StringFlag{
			Name:    "llm-endpoint",
			Value:   "",
			Usage:   "Custom endpoint for the LLM provider (if applicable)",
			EnvVars: []string{"CA_LLM_ENDPOINT"},
		},
		&cli.IntFlag{
			Name:    "llm-max-tokens",
			Value:   200,
			Usage:   "Maximum tokens for completion",
			EnvVars: []string{"CA_LLM_MAX_TOKENS"},
		},
		&cli.IntFlag{
			Name:    "llm-max-retries",
			Value:   3,
			Usage:   "Maximum number of retries for API calls",
			EnvVars: []string{"CA_LLM_MAX_RETRIES"},
		},
		&cli.DurationFlag{
			Name:    "llm-retry-delay",
			Value:   2 * time.Second,
			Usage:   "Delay between retries (e.g. 2s, 500ms)",
			EnvVars: []string{"CA_LLM_RETRY_DELAY"},
		},
		&cli.IntFlag{
			Name:    "llm-log-level",
			Value:   1,
			Usage:   "Log level (0=debug, 1=info, etc.)",
			EnvVars: []string{"CA_LLM_LOG_LEVEL"},
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

// NewLLMConfigFromContext extracts the LLM configuration from the CLI context.
func NewLLMConfigFromContext(c *cli.Context) *LLMConfig {
	return &LLMConfig{
		Provider:   c.String("llm-provider"),
		Model:      c.String("llm-model"),
		APIKey:     c.String("llm-api-key"),
		Endpoint:   c.String("llm-endpoint"),
		MaxTokens:  c.Int("llm-max-tokens"),
		MaxRetries: c.Int("llm-max-retries"),
		RetryDelay: c.Duration("llm-retry-delay"),
		LogLevel:   gollm.LogLevel(c.Int("llm-log-level")),
	}
}
