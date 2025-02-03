// Copyright (c) 2025 - Chris Robinson
// Licensed under the BSD 3-Clause License.
// See LICENSE file for details.

package session_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/chrisrob11/codeassistant/internal/session"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

// setupTestEnv sets up a temporary test environment and initializes the session directory.
func setupTestEnv(t *testing.T) string {
	tempDir := t.TempDir() // Create a temporary directory for testing

	// Initialize the session directory
	exists, err := session.CreateOrCheckSessionDir(tempDir)
	assert.NilError(t, err, "Failed to initialize session directories.")
	assert.Assert(t, cmp.Equal(exists, false), "Expected no session file initially.")

	return tempDir
}

// TestCreateOrCheckSessionDir_CreatesDirs verifies that directories are created correctly.
func TestCreateOrCheckSessionDir_CreatesDirs(t *testing.T) {
	sessionDir := setupTestEnv(t)

	// Verify session root directory exists
	_, err := os.Stat(sessionDir)
	assert.NilError(t, err, "Session root directory should exist.")

	sessionHistoryDir := session.BuildSessionHistoryPath(sessionDir)
	_, err = os.Stat(sessionHistoryDir)
	assert.NilError(t, err, "History directory should exist.")
}

// TestCreateOrCheckSessionDir_DetectsSessionFile verifies that an existing session file is detected.
func TestCreateOrCheckSessionDir_DetectsSessionFile(t *testing.T) {
	sessionDir := setupTestEnv(t)

	sessionFilePath := session.BuildCurrentSessionFilePath(sessionDir)
	// Manually create a session file
	err := os.WriteFile(sessionFilePath, []byte("{}"), 0600)
	assert.NilError(t, err, "Failed to create test session file.")

	// Run CreateOrCheckSessionDir again
	exists, err := session.CreateOrCheckSessionDir(sessionDir)
	assert.NilError(t, err, "Error checking session file existence.")
	assert.Assert(t, cmp.Equal(exists, true), "Expected session file to exist.")
}

// TestCreateOrCheckSessionDir_FailsOnPermissions verifies failure when directory creation is not allowed.
func TestCreateOrCheckSessionDir_FailsOnPermissions(t *testing.T) {
	tempDir := setupTestEnv(t)

	// Create a non-writable session root directory
	lockedDir := filepath.Join(tempDir, "locked")
	err := os.MkdirAll(lockedDir, 0000) // No permissions
	assert.NilError(t, err, "Failed to create locked directory.")

	// Try to run the function in the locked directory
	_, err = session.CreateOrCheckSessionDir(lockedDir)
	assert.Assert(t, errors.Is(err, session.ErrSessionHistoryCreate), "Expected ErrSessionRootCreate due to permissions.")

	// Restore permissions so cleanup works
	// nolint:gosec // Why: test code
	err = os.Chmod(lockedDir, 0755)
	assert.NilError(t, err)
}

// TestStartSession_Success verifies that a session starts correctly.
func TestStartSession_Success(t *testing.T) {
	sessionDir := setupTestEnv(t)
	sessionFilePath := session.BuildCurrentSessionFilePath(sessionDir)

	req := &session.StartSessionRequest{Name: "Test Session", Dir: sessionDir}
	err := session.StartSession(req)
	assert.NilError(t, err)

	// Check if session file was created
	_, err = os.Stat(sessionFilePath)
	assert.NilError(t, err)

	// Validate session contents
	// nolint:gosec // Why: test code
	data, err := os.ReadFile(sessionFilePath)
	assert.NilError(t, err)

	var session session.Session
	err = json.Unmarshal(data, &session)
	assert.NilError(t, err)

	assert.Assert(t, cmp.Equal(session.Name, "Test Session"))
	assert.Assert(t, cmp.Equal(len(session.Steps), 0)) // Ensure steps are empty
}

// TestStartSession_AlreadyExists ensures that a session cannot start if one is active.
func TestStartSession_AlreadyExists(t *testing.T) {
	sessionDir := setupTestEnv(t)

	req := &session.StartSessionRequest{Name: "Duplicate Session", Dir: sessionDir}

	// First session should start successfully
	assert.NilError(t, session.StartSession(req))

	// Second attempt should fail with ErrSessionExists
	err := session.StartSession(req)
	assert.Assert(t, errors.Is(err, session.ErrSessionExists), "Expected ErrSessionExists, got: %v", err)
}

// TestEndSession_Success ensures that ending a session archives it correctly.
func TestEndSession_Success(t *testing.T) {
	sessionDir := setupTestEnv(t)
	sessionFilePath := session.BuildCurrentSessionFilePath(sessionDir)
	sessionHistoryPath := session.BuildSessionHistoryPath(sessionDir)

	// Start a session first
	req := &session.StartSessionRequest{Name: "Session to End", Dir: sessionDir}
	assert.NilError(t, session.StartSession(req))

	// End the session
	err := session.EndSession(sessionDir)
	assert.NilError(t, err)

	// Ensure session file is deleted
	_, err = os.Stat(sessionFilePath)
	assert.Assert(t, os.IsNotExist(err), "Expected session file to be deleted.")

	// Ensure archive file exists
	files, err := os.ReadDir(sessionHistoryPath)
	assert.NilError(t, err)
	assert.Assert(t, len(files) > 0, "Expected at least one archived session file.")
}

// TestEndSession_NoActiveSession ensures that ending a session without an active one returns an error.
func TestEndSession_NoActiveSession(t *testing.T) {
	sessionPath := setupTestEnv(t)

	err := session.EndSession(sessionPath)
	assert.Assert(t, errors.Is(err, session.ErrNoActiveSession), "Expected ErrNoActiveSession, got: %v", err)
}

// TestEndSession_ArchiveFailure simulates a failure in archiving the session.
func TestEndSession_ArchiveFailure(t *testing.T) {
	sessionDir := setupTestEnv(t)
	sessionFilePath := session.BuildCurrentSessionFilePath(sessionDir)
	sessionHistoryPath := session.BuildSessionHistoryPath(sessionDir)

	// Start a session first
	req := &session.StartSessionRequest{Name: "Session to Archive Fail", Dir: sessionDir}
	assert.NilError(t, session.StartSession(req))

	// Remove archive directory to force a failure
	// nolint:gosec // Why: test code
	os.RemoveAll(sessionHistoryPath)

	err := session.EndSession(sessionDir)
	assert.Assert(t, errors.Is(err, session.ErrSessionMkdir), "Expected ErrSessionMkdir, got: %v", err)

	// Ensure session file still exists (not deleted on failure)
	_, err = os.Stat(sessionFilePath)
	assert.NilError(t, err, "Session file should still exist after archive failure.")
}
