// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chrisrob11/codeassistant/internal/session"
	"github.com/teilomillet/gollm"
	cli "github.com/urfave/cli/v2"
)

// Predefined Errors.
var (
	ErrFailedToGetCurrentDir  = errors.New("failed to get current directory")
	ErrMissingPrompt          = errors.New("failed as prompt not specified")
	ErrAIProcessingFailed     = errors.New("AI modification failed")
	ErrFailedToWriteChanges   = errors.New("failed to write changes")
	ErrFailedToLoadSession    = errors.New("failed to load session")
	ErrFailedToSaveSession    = errors.New("failed to save session")
	ErrFileOutsideCurrentDir  = errors.New("file is outside the current directory")
	ErrFailedToResolveAbsPath = errors.New("failed to resolve absolute path")
	ErrFailedToReadFile       = errors.New("failed to read file")
	ErrFilesMustBeSpecified   = errors.New("failed as files not specified")
)

// CodeCommand applies AI modifications to code.
func CodeCommand() *cli.Command {
	return &cli.Command{
		Name:  "code",
		Usage: "Apply AI modifications to code",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "files",
				Aliases: []string{"f"},
				Usage:   "Specify files to modify",
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
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("%w", ErrFailedToGetCurrentDir)
			}

			prompt := c.Args().First()
			if prompt == "" {
				return ErrMissingPrompt
			}

			files := c.StringSlice("files")
			if len(files) == 0 {
				return ErrFilesMustBeSpecified
			}

			absFilePaths := []string{}
			dryRun := c.Bool("dry-run")

			for _, f := range files {
				absPath, err := isValidFilePath(currentDir, f)
				if err != nil {
					return err
				}
				absFilePaths = append(absFilePaths, absPath)
			}

			llmConfig := NewLLMConfigFromContext(c)
			if err := llmConfig.Validate(); err != nil {
				return err
			}

			llm, err := llmConfig.BuildLLM()
			if err != nil {
				return err
			}

			return executeCodeCommand(llm, currentDir, prompt, absFilePaths, dryRun)
		},
	}
}

func executeCodeCommand(llm gollm.LLM, currentDir, prompt string, absFilePaths []string, dryRun bool) error {
	currentSession, err := session.LoadCurrentSession(currentDir)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedToLoadSession, err)
	}

	// Modify code
	modifications, err := modifyCode(llm, prompt, absFilePaths)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAIProcessingFailed, err)
	}

	// Handle dry-run
	if dryRun {
		for file, mod := range modifications {
			fmt.Printf("Changes for %s:\n%s\n", file, mod)
		}

		return nil
	}

	// Write modifications
	for file, mod := range modifications {
		err := os.WriteFile(file, []byte(mod), 0600)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedToWriteChanges, err)
		}
	}

	var lastStep *session.Step
	for _, step := range currentSession.Steps {
		lastStep = step
	}

	stepID := 1
	if lastStep != nil {
		stepID = lastStep.ID + 1
	}

	// Track changes in session
	currentSession.Steps = append(currentSession.Steps, &session.Step{
		ID:        stepID,
		Command:   session.Command{Prompt: prompt, Files: absFilePaths},
		Timestamp: time.Now(),
	})

	err = session.SaveCurrentSession(currentDir, currentSession)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedToSaveSession, err)
	}

	return nil
}

// Function to modify code using AI.
func modifyCode(llm gollm.LLM, prompt string, files []string) (map[string]string, error) {
	modifications := make(map[string]string)

	for _, file := range files {
		// nolint:gosec //Why: files are validated within a specific path
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToReadFile, err)
		}

		// Call AI processing (mocking the LLM call)
		modifiedContent, err := processWithLLM(llm, prompt, string(content))
		if err != nil {
			return nil, err
		}

		modifications[file] = modifiedContent
	}

	return modifications, nil
}

// Ensure that the file is inside the current working directory.
func isValidFilePath(currentDir, filePath string) (string, error) {
	// Get absolute path of the file
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrFailedToResolveAbsPath, err)
	}

	// Ensure the file is within the current directory
	if !strings.HasPrefix(absFilePath, currentDir) {
		return "", fmt.Errorf("%w: %s", ErrFileOutsideCurrentDir, filePath)
	}

	return absFilePath, nil
}

// Use gollm to process the AI request.
func processWithLLM(llm gollm.LLM, prompt, fileContent string) (string, error) {
	ctx := context.Background()

	// Create a basic prompt
	fullPrompt := fmt.Sprintf(
		"Use the following prompt '%s' to modify the file contents and output the update code: \n%s",
		prompt, fileContent)
	promptValue := gollm.NewPrompt(fullPrompt)

	// Generate a response
	response, err := llm.Generate(ctx, promptValue)
	if err != nil {
		return "", err
	}

	return response, nil
}
