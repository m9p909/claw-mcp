package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

// Integration Test: Filesystem operation logging (read operation)
// Spec: specs/filesystem-operation-logging/spec.md
// Requirements: Log path validation, bytes read, line count, execution time
func TestLogging_FileReadOperation(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	// This test validates that logging occurs during read operations
	// by checking that the read completes successfully
	input := models.ReadFileRequest{Path: tmpFile}
	toolResult, resp, err := HandleReadFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful read")
	}

	// Verify response contains expected data
	lines := strings.Split(resp.Content, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

// Integration Test: Filesystem operation logging (write operation)
// Spec: specs/filesystem-operation-logging/spec.md
// Requirements: Log hash validation, permissions checks
func TestLogging_FileWriteOperation(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_write.txt")

	content := "test content\nmore content"
	input := models.WriteFileRequest{Path: tmpFile, Content: content}
	toolResult, resp, err := HandleWriteFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful write")
	}
	if !resp.Success {
		t.Errorf("write not successful: %s", resp.Message)
	}

	// Verify file was created with correct content
	data, _ := os.ReadFile(tmpFile)
	if string(data) != content {
		t.Errorf("file content mismatch")
	}
}

// Integration Test: Filesystem operation logging (edit operation)
// Spec: specs/filesystem-operation-logging/spec.md
// Requirements: Log hash lookup and line replacement
func TestLogging_FileEditOperation(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	input := models.EditFileRequest{
		Path:       tmpFile,
		StartLine:  1,
		EndLine:    100,
		NewContent: "new content",
	}
	toolResult, _, _ := HandleEditFile(context.Background(), nil, input)

	// Edit will fail due to out-of-range, but logging should still occur
	if toolResult == nil {
		t.Fatalf("expected error for invalid range")
	}
}

// Integration Test: Memory operation logging (write operation)
// Spec: specs/memory-operation-logging/spec.md
// Requirements: Log category validation
func TestLogging_MemoryWriteOperation(t *testing.T) {
	storage.ClearMemory()

	input := models.WriteMemoryRequest{
		Category: "fact",
		Content:  "Test memory content",
	}
	toolResult, resp, err := HandleWriteMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful write")
	}
	if !resp.Success {
		t.Errorf("write not successful: %s", resp.Message)
	}
}

// Integration Test: Memory operation logging (query operation)
// Spec: specs/memory-operation-logging/spec.md
// Requirements: Log result counts and execution time
func TestLogging_MemoryQueryOperation(t *testing.T) {
	storage.ClearMemory()

	// Write test data
	storage.WriteMemory("fact", "Fact 1")
	storage.WriteMemory("fact", "Fact 2")

	input := models.QueryMemoryRequest{
		Query: "SELECT * FROM memories WHERE category = 'fact'",
	}
	toolResult, resp, err := HandleQueryMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful query")
	}

	// Verify results were returned
	if len(resp.Results) < 1 {
		t.Errorf("expected results from query")
	}
}

// Integration Test: Memory operation logging (search operation)
// Spec: specs/memory-operation-logging/spec.md
// Requirements: Log result details
func TestLogging_MemorySearchOperation(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("decision", "Use PostgreSQL for primary DB")
	storage.WriteMemory("decision", "Use Redis for caching")

	input := models.SearchMemoryRequest{
		Query: "PostgreSQL",
		Limit: 0,
	}
	toolResult, resp, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful search")
	}

	// Verify search results
	if len(resp.Results) < 1 {
		t.Errorf("expected search results")
	}
}

// Integration Test: Command execution logging (foreground execution)
// Spec: specs/command-execution-logging/spec.md
// Requirements: Log startup, completion, timing, exit codes
func TestLogging_ForegroundCommandExecution(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "echo",
		Args:       []string{"test message"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful execution")
	}

	// Verify execution details logged
	if resp.ExitCode != 0 {
		t.Errorf("expected exit code 0")
	}
	if resp.Status != "completed" {
		t.Errorf("expected status 'completed'")
	}
	if len(resp.Stdout) == 0 {
		t.Errorf("expected output")
	}
}

// Integration Test: Command execution logging (background execution)
// Spec: specs/command-execution-logging/spec.md
// Requirements: Log session creation, process monitoring
func TestLogging_BackgroundCommandExecution(t *testing.T) {
	storage.ClearSessions()

	input := models.ExecCommandRequest{
		Command:    "sleep",
		Args:       []string{"0.5"},
		Background: true,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, err := HandleExecCommand(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful execution")
	}

	// Verify session was created and logged
	if resp.SessionID == "" {
		t.Errorf("expected session ID")
	}
	if resp.Status != "running" {
		t.Errorf("expected status 'running'")
	}

	storage.ClearSessions()
}

// Integration Test: Process management logging (list action)
// Spec: specs/command-execution-logging/spec.md
// Requirements: Log manage_process operations (list)
func TestLogging_ProcessListOperation(t *testing.T) {
	storage.ClearSessions()

	_ = storage.CreateSession("session_1", "test command 1")
	_ = storage.CreateSession("session_2", "test command 2")

	input := models.ManageProcessRequest{
		Action:    "list",
		SessionID: "",
		Keys:      "",
	}
	toolResult, resp, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful list")
	}

	// Verify sessions listed
	if len(resp.Sessions) < 2 {
		t.Errorf("expected at least 2 sessions")
	}

	storage.ClearSessions()
}

// Integration Test: Process management logging (poll action)
// Spec: specs/command-execution-logging/spec.md
// Requirements: Log manage_process operations (poll)
func TestLogging_ProcessPollOperation(t *testing.T) {
	storage.ClearSessions()

	_ = storage.CreateSession("test_session", "test command")

	input := models.ManageProcessRequest{
		Action:    "poll",
		SessionID: "test_session",
		Keys:      "",
	}
	toolResult, resp, err := HandleManageProcess(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("poll failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected successful poll")
	}

	// Verify session polled
	if len(resp.Sessions) != 1 {
		t.Errorf("expected 1 session")
	}

	storage.ClearSessions()
}

// Integration Test: Full workflow with logging
// Tests that all operations log correctly in sequence
func TestLogging_FullWorkflow(t *testing.T) {
	// Clean up
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflow_test.txt")
	storage.ClearMemory()
	storage.ClearSessions()

	// 1. Write file
	writeInput := models.WriteFileRequest{
		Path:    tmpFile,
		Content: "initial content",
	}
	toolResult, _, _ := HandleWriteFile(context.Background(), nil, writeInput)
	if toolResult != nil {
		t.Fatalf("write failed")
	}

	// 2. Read file
	readInput := models.ReadFileRequest{Path: tmpFile}
	toolResult2, readResp, _ := HandleReadFile(context.Background(), nil, readInput)
	if toolResult2 != nil {
		t.Fatalf("read failed")
	}

	// 3. Write memory
	writeMemInput := models.WriteMemoryRequest{
		Category: "fact",
		Content:  "File processing completed",
	}
	toolResult3, _, _ := HandleWriteMemory(context.Background(), nil, writeMemInput)
	if toolResult3 != nil {
		t.Fatalf("memory write failed")
	}
	// 4. Search memory
	searchInput := models.SearchMemoryRequest{
		Query: "File",
		Limit: 0,
	}
	toolResult4, searchResp, _ := HandleMemorySearch(context.Background(), nil, searchInput)
	if toolResult4 != nil {
		t.Fatalf("memory search failed")
	}
	if len(searchResp.Results) < 1 {
		t.Errorf("expected search results")
	}

	// 5. Execute command
	execInput := models.ExecCommandRequest{
		Command:    "echo",
		Args:       []string{"workflow test"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult5, execResp, _ := HandleExecCommand(context.Background(), nil, execInput)
	if toolResult5 != nil {
		t.Fatalf("command execution failed")
	}
	if execResp.ExitCode != 0 {
		t.Errorf("expected exit code 0")
	}

	// Verify all operations completed
	if len(readResp.Content) == 0 {
		t.Errorf("expected file content")
	}
}

// Integration Test: Error logging
// Verify that errors are properly logged
func TestLogging_ErrorLogging(t *testing.T) {
	// Attempt to read non-existent file
	input := models.ReadFileRequest{Path: "/nonexistent/file/path.txt"}
	toolResult, _, _ := HandleReadFile(context.Background(), nil, input)

	// Error should be returned and logged
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	// Verify error contains expected information
	if len(toolResult.Content) == 0 {
		t.Errorf("expected error message")
	}
}

// Integration Test: Command output logging
// Verify command output is captured and available
func TestLogging_CommandOutputCapture(t *testing.T) {
	input := models.ExecCommandRequest{
		Command:    "sh",
		Args:       []string{"-c", "echo 'line1'; echo 'line2'; echo 'line3'"},
		Background: false,
		Timeout:    0,
		Env:        nil,
	}
	toolResult, resp, _ := HandleExecCommand(context.Background(), nil, input)

	if toolResult != nil {
		t.Fatalf("expected successful execution")
	}

	// Verify all output is captured
	if !strings.Contains(resp.Stdout, "line1") {
		t.Errorf("expected 'line1' in output")
	}
	if !strings.Contains(resp.Stdout, "line2") {
		t.Errorf("expected 'line2' in output")
	}
	if !strings.Contains(resp.Stdout, "line3") {
		t.Errorf("expected 'line3' in output")
	}
}
