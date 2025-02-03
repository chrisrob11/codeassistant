// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

// Package session contains structures to track llm calls
package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// File Paths.
const sessionFileName = ".ca_session.json"
const sessionHistoryDirName = ".ca_sessions"

// Custom session errors.
var (
	ErrSessionExists              = errors.New("session already in progress")
	ErrNoActiveSession            = errors.New("no active session found")
	ErrSessionWriteFail           = errors.New("failed to write session file")
	ErrSessionReadFail            = errors.New("failed to read session file")
	ErrSessionParseFail           = errors.New("failed to parse session file")
	ErrSessionArchive             = errors.New("failed to archive session file")
	ErrSessionDelete              = errors.New("failed to delete session file")
	ErrSessionMkdir               = errors.New("failed to create session history directory")
	ErrSessionRootCreate          = errors.New("failed to create session root directory")
	ErrSessionHistoryCreate       = errors.New("failed to create session history directory")
	ErrSessionFileExistanceFailed = errors.New("failed to check if current session file exists")
	ErrSessionDirNotSpecified     = errors.New("failed as the session dir was not specified")
	ErrSessionNameSpecified       = errors.New("failed as the session name was not specified")
)

// Command represents the command details associated with a step.
type Command struct {
	Prompt       string          `json:"prompt"`
	Flags        map[string]bool `json:"flags"`
	AppliedFiles []string        `json:"applied_files"`
}

// FilesDiff represents the differences in files during the step.
type FilesDiff struct {
	Created  []string `json:"created"`
	Modified []string `json:"modified"`
	Deleted  []string `json:"deleted"`
}

// GitState represents the state of the Git repository at a specific point.
type GitState struct {
	Commit            string   `json:"commit"`
	StateBeforeChange string   `json:"state_before_change,omitempty"`
	UncommittedFiles  []string `json:"uncommitted_files,omitempty"`
}

// Git represents the Git information before and after the step.
type Git struct {
	Pre  GitState `json:"pre"`
	Post GitState `json:"post"`
}

// Step represents an individual step within a session.
type Step struct {
	ID        int       `json:"id"`
	Prompt    string    `json:"prompt"`
	Command   Command   `json:"command"`
	Timestamp time.Time `json:"timestamp"`
	FilesDiff FilesDiff `json:"files_diff"`
	Git       Git       `json:"git"`
}

// Session represents a user session with an llm.
type Session struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	Steps       []Step    `json:"steps"`
}

type StartSessionRequest struct {
	Name string
	Dir  string
}

// BuildCurrentSessionFilePath is the path to current session file path.
func BuildCurrentSessionFilePath(path string) string {
	return filepath.Join(path, sessionFileName)
}

// BuildSessionHistoryPath is the path to all history sessions.
func BuildSessionHistoryPath(path string) string {
	return filepath.Join(path, sessionHistoryDirName)
}

// CreateOrCheckSessionDir initializes the session root directory and history dir for storage
// if they don't exist, returns whether the current session file exists.
func CreateOrCheckSessionDir(path string) (currentSessionExists bool, err error) {
	// Ensure the base session directory exists
	if err := os.MkdirAll(path, 0750); err != nil {
		return false, ErrSessionRootCreate
	}

	// Create a history subdirectory inside the session root
	historyDir := BuildSessionHistoryPath(path)
	if err := os.MkdirAll(historyDir, 0750); err != nil {
		return false, fmt.Errorf("%w: %v", ErrSessionHistoryCreate, err)
	}

	sessionFilePath := BuildCurrentSessionFilePath(path)
	// Check if session file exists
	if _, err := os.Stat(sessionFilePath); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("%w: %w", ErrSessionFileExistanceFailed, err)
	}

	return false, nil
}

// StartSession starts a new coding session.
func StartSession(startSession *StartSessionRequest) error {
	if startSession.Dir == "" {
		return ErrSessionDirNotSpecified
	}

	if startSession.Name == "" {
		return ErrSessionNameSpecified
	}

	sessionFileExists, err := CreateOrCheckSessionDir(startSession.Dir)
	if err != nil {
		return err
	}

	// Check if a session is already in progress
	if sessionFileExists {
		return ErrSessionExists
	}

	// Create a new session object
	session := Session{
		ID:        uuid.New().String(),
		Name:      startSession.Name,
		CreatedAt: time.Now(),
		Steps:     []Step{},
	}

	// Convert session to JSON
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: could not serialize session data", ErrSessionWriteFail)
	}

	fullSessionFilePath := filepath.Join(startSession.Dir, sessionFileName)

	// Write session to file with secure permissions
	if err := os.WriteFile(fullSessionFilePath, data, 0600); err != nil {
		return fmt.Errorf("%w: unable to create session file", ErrSessionWriteFail)
	}

	fmt.Printf("✅ Session started: %s\n", startSession.Name)

	return nil
}

// EndSession ends the current session and archives it.
func EndSession(sessionDir string) error {
	// Check if there is an active session
	sessionFilePath := BuildCurrentSessionFilePath(sessionDir)
	if _, err := os.Stat(sessionFilePath); os.IsNotExist(err) {
		return ErrNoActiveSession
	}

	// Read the session file
	// nolint:gosec // Why: session path is not specified by a user but by the system
	data, err := os.ReadFile(sessionFilePath)
	if err != nil {
		return fmt.Errorf("%w: unable to read session file", ErrSessionReadFail)
	}

	// Parse session JSON
	var session Session
	if jsonErr := json.Unmarshal(data, &session); jsonErr != nil {
		return fmt.Errorf("%w: corrupted session data", ErrSessionParseFail)
	}

	// Mark session completion
	session.CompletedAt = time.Now()

	// Log session duration
	duration := session.CompletedAt.Sub(session.CreatedAt)
	fmt.Printf("📅 Session \"%s\" lasted %s\n", session.Name, duration)

	sessionHistoryDirPath := BuildSessionHistoryPath(sessionDir)
	if _, err := os.Stat(sessionHistoryDirPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: history directory does not exist", ErrSessionMkdir)
	} else if err != nil {
		return fmt.Errorf("%w: failed to stat history directory: %v", ErrSessionMkdir, err)
	}

	// Create an archive file with a safe filename
	safeName := strings.ReplaceAll(strings.ToLower(session.Name), " ", "_")
	sessionHistoryPath := filepath.Join(sessionHistoryDirPath,
		fmt.Sprintf("%s_%s.json", time.Now().Format("20060102-150405"), safeName))

	// Convert updated session data to JSON
	updatedData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("%w: unable to serialize session data", ErrSessionWriteFail)
	}

	// Write session archive file
	if err := os.WriteFile(sessionHistoryPath, updatedData, 0600); err != nil {
		return fmt.Errorf("%w: unable to write archive file", ErrSessionArchive)
	}

	// Delete the active session file
	if err := os.Remove(sessionFilePath); err != nil {
		return fmt.Errorf("%w: failed to remove session file", ErrSessionDelete)
	}

	fmt.Println("✅ Session ended and archived.")

	return nil
}
