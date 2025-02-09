package cmd

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/teilomillet/gollm"
	"github.com/teilomillet/gollm/config"
)

// A list of providers that require an API token.
var providersRequireAPIToken = []string{
	"openai",
	"anthropic",
	"azureopenai",
	"cohere",
	"google",
	"huggingface",
	"replicate",
	"mosaic",
	"promptlayer",
}

// LLMConfig is a catch-all config. Some fields only matter for certain providers.
type LLMConfig struct {
	Provider   string         // E.g. "openai", "ollama", "azureopenai", etc.
	Model      string         // E.g. "gpt-4", "deepseek", ...
	APIKey     string         // For openai/anthropic, not needed by ollama
	Endpoint   string         // E.g. for azureopenai or a remote Ollama
	MaxTokens  int            // Common
	MaxRetries int            // Common
	RetryDelay time.Duration  // Common
	LogLevel   gollm.LogLevel // Common

	defaultsSet bool // internal flag to ensure we only set defaults once
}

// setDefaults populates zero-valued fields with a sensible default.
// You can customize these however you like.
func (c *LLMConfig) setDefaults() {
	if c.defaultsSet {
		// Only set defaults once.
		return
	}
	c.defaultsSet = true

	if c.MaxTokens == 0 {
		c.MaxTokens = 200
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = 3
	}
	if c.RetryDelay == 0 {
		c.RetryDelay = 2 * time.Second
	}
	if c.LogLevel == 0 {
		c.LogLevel = gollm.LogLevelInfo
	}
}

// Validate checks that all required fields are present
// for the given provider.
func (c *LLMConfig) Validate() error {
	// Provider must be non-empty
	if c.Provider == "" {
		return errors.New("provider is required")
	}

	// If provider is in providersRequireAPIToken, ensure APIKey is non-empty
	if slices.Contains(providersRequireAPIToken, c.Provider) && c.APIKey == "" {
		return fmt.Errorf("provider %q requires an API token, but none was provided", c.Provider)
	}

	if c.Model == "" {
		return errors.New("model is required")
	}

	return nil
}

// BuildLLM applies the validated fields to construct a gollm.LLM.
func (c *LLMConfig) BuildLLM() (gollm.LLM, error) {
	// 1) Set defaults if needed
	c.setDefaults()

	// 2) Validate the config
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// 3) Build the base options
	opts := []config.ConfigOption{
		gollm.SetProvider(c.Provider),
		gollm.SetModel(c.Model),
		gollm.SetMaxTokens(c.MaxTokens),
		gollm.SetMaxRetries(c.MaxRetries),
		gollm.SetRetryDelay(c.RetryDelay),
		gollm.SetLogLevel(c.LogLevel),
	}

	// If an APIKey was provided, apply it
	if c.APIKey != "" {
		opts = append(opts, gollm.SetAPIKey(c.APIKey))
	}

	// If an Endpoint was provided, we choose which gollm.Option to apply.
	// For example, Ollama -> gollm.SetOllamaEndpoint, Azure -> gollm.SetAzureOpenAIEndpoint, etc.
	switch c.Provider {
	case "ollama":
		if c.Endpoint != "" {
			opts = append(opts, gollm.SetOllamaEndpoint(c.Endpoint))
		}
	default:
		// For other providers, ignore or handle as needed
	}

	// 4) Initialize the LLM
	llm, err := gollm.NewLLM(opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating LLM for provider %q: %w", c.Provider, err)
	}

	return llm, nil
}
