package tools

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

// Test ExecCommand: Simple foreground command
func TestHandleExecCommand_ForegroundSimple(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "echo",
		Args:       []string{"hello world"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !strings.Contains(resp.Stdout, "hello world") {
		t.Errorf("expected 'hello world' in stdout, got: %s", resp.Stdout)
	}
	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", resp.ExitCode)
	}
	if resp.Status != "completed" {
		t.Errorf("expected status 'completed', got %s", resp.Status)
	}
}

// Test ExecCommand: Foreground with arguments
func TestHandleExecCommand_ForegroundWithArgs(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "echo",
		Args:       []string{"line1", "line2", "line3"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", resp.ExitCode)
	}
	if len(resp.Stdout) == 0 {
		t.Errorf("expected output in stdout")
	}
}

// Test ExecCommand: Foreground with non-zero exit
func TestHandleExecCommand_ForegroundNonZeroExit(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "sh",
		Args:       []string{"-c", "exit 42"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if resp.ExitCode != 42 {
		t.Errorf("expected exit code 42, got %d", resp.ExitCode)
	}
	if resp.Status != "completed" {
		t.Errorf("expected status 'completed', got %s", resp.Status)
	}
}

// Test ExecCommand: Empty command
func TestHandleExecCommand_EmptyCommand(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "",
		Args:       []string{},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, _, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for empty command")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_REQUEST") {
		t.Errorf("expected INVALID_REQUEST error")
	}
}

// Test ExecCommand: Invalid command
func TestHandleExecCommand_InvalidCommand(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "nonexistentcommand12345",
		Args:       []string{},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, _, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for invalid command")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "EXEC_FAILED") {
		t.Errorf("expected EXEC_FAILED error")
	}
}

// Test ExecCommand: Background execution
func TestHandleExecCommand_BackgroundExecution(t *testing.T) {
	storage.ClearSessions()

	input := models.ExecCommandRequest{
		Command:    "sleep",
		Args:       []string{"1"},
		Background: true,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if resp.SessionID == "" {
		t.Fatalf("expected session ID for background process")
	}
	if resp.Status != "running" {
		t.Errorf("expected status 'running', got %s", resp.Status)
	}

	// Clean up
	storage.ClearSessions()
}

// Test ExecCommand: Environment variables
func TestHandleExecCommand_WithEnvironment(t *testing.T) {
	input := models.ExecCommandRequest{
		Command: "sh",
		Args:    []string{"-c", "echo $TEST_VAR"},
		Env: map[string]string{
			"TEST_VAR": "test_value",
		},
		Background: false,
		Timeout:    0,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !strings.Contains(resp.Stdout, "test_value") {
		t.Errorf("expected TEST_VAR value in output")
	}
}

// Test ManageProcess: List sessions
func TestHandleManageProcess_ListSessions(t *testing.T) {
	storage.ClearSessions()

	// Create some sessions
	_ = storage.CreateSession("session_1", "echo test1")
	_ = storage.CreateSession("session_2", "echo test2")

	input := models.ManageProcessRequest{
		Action:    "list",
		SessionID: "",
		Keys:      "",
	}
	toolResult, resp, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Sessions) < 2 {
		t.Errorf("expected at least 2 sessions, got %d", len(resp.Sessions))
	}

	storage.ClearSessions()
}

// Test ManageProcess: Poll session
func TestHandleManageProcess_PollSession(t *testing.T) {
	storage.ClearSessions()

	_ = storage.CreateSession("test_session", "echo running")

	input := models.ManageProcessRequest{
		Action:    "poll",
		SessionID: "test_session",
		Keys:      "",
	}
	toolResult, resp, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(resp.Sessions))
	}
	if resp.Sessions[0].SessionID != "test_session" {
		t.Errorf("expected session_id 'test_session', got %s", resp.Sessions[0].SessionID)
	}

	storage.ClearSessions()
}

// Test ManageProcess: Poll non-existent session
func TestHandleManageProcess_PollNonExistent(t *testing.T) {
	storage.ClearSessions()

	input := models.ManageProcessRequest{
		Action:    "poll",
		SessionID: "nonexistent_session",
		Keys:      "",
	}
	toolResult, _, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for non-existent session")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "PROCESS_NOT_FOUND") {
		t.Errorf("expected PROCESS_NOT_FOUND error")
	}

	storage.ClearSessions()
}

// Test ManageProcess: Poll without session ID
func TestHandleManageProcess_PollMissingSessionID(t *testing.T) {
	input := models.ManageProcessRequest{
		Action:    "poll",
		SessionID: "",
		Keys:      "",
	}
	toolResult, _, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_REQUEST") {
		t.Errorf("expected INVALID_REQUEST error")
	}
}

// Test ManageProcess: Unknown action
func TestHandleManageProcess_UnknownAction(t *testing.T) {
	input := models.ManageProcessRequest{
		Action:    "unknown_action",
		SessionID: "",
		Keys:      "",
	}
	toolResult, _, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for unknown action")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_REQUEST") {
		t.Errorf("expected INVALID_REQUEST error")
	}
}

// Test ManageProcess: Send keys (not implemented)
func TestHandleManageProcess_SendKeys(t *testing.T) {
	input := models.ManageProcessRequest{
		Action:    "send_keys",
		SessionID: "test_session",
		Keys:      "test_input",
	}
	toolResult, _, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result (not implemented)")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INTERNAL_ERROR") {
		t.Errorf("expected INTERNAL_ERROR for unimplemented send_keys")
	}
}

// Test ManageProcess: Kill (not implemented)
func TestHandleManageProcess_Kill(t *testing.T) {
	input := models.ManageProcessRequest{
		Action:    "kill",
		SessionID: "test_session",
		Keys:      "",
	}
	toolResult, _, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result (not implemented)")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INTERNAL_ERROR") {
		t.Errorf("expected INTERNAL_ERROR for unimplemented kill")
	}
}

// Integration Test: Foreground command with output capture
func TestIntegration_ForegroundCommandWithOutput(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "sh",
		Args:       []string{"-c", "echo 'line1'; echo 'line2'; echo 'line3'"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}

	// Validate output
	if !strings.Contains(resp.Stdout, "line1") {
		t.Errorf("expected 'line1' in output")
	}
	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", resp.ExitCode)
	}
	if resp.Status != "completed" {
		t.Errorf("expected status 'completed'")
	}
}

// Integration Test: Background command → Poll workflow
func TestIntegration_BackgroundCommandWithPolling(t *testing.T) {
	storage.ClearSessions()

	// Start background command
	execInput := models.ExecCommandRequest{
		Command:    "echo",
		Args:       []string{"test output"},
		Background: true,
		Timeout:    0,
		Env:        nil,
	}
	_, execResp, _ := HandleExecCommand(context.Background(), nil, execInput)
	sessionID := execResp.SessionID

	if sessionID == "" {
		t.Fatalf("expected session ID")
	}

	// Poll the session
	pollInput := models.ManageProcessRequest{
		Action:    "poll",
		SessionID: sessionID,
		Keys:      "",
	}
	toolResult, pollResp, _ := HandleManageProcess(context.Background(), nil, pollInput)

	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(pollResp.Sessions) != 1 {
		t.Errorf("expected 1 session")
	}

	// Wait a bit for command to complete
	time.Sleep(500 * time.Millisecond)

	// Poll again and check status
	_, pollResp2, _ := HandleManageProcess(context.Background(), nil, pollInput)
	if len(pollResp2.Sessions) > 0 {
		session := pollResp2.Sessions[0]
		if session.SessionID != sessionID {
			t.Errorf("expected correct session ID")
		}
	}

	storage.ClearSessions()
}

// Test ExecCommand: Command with stderr
func TestHandleExecCommand_WithStderr(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "sh",
		Args:       []string{"-c", "echo 'error message' >&2; echo 'output'"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0")
	}
}
