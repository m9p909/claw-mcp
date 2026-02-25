package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/hash"
	"awesomeProject/pkg/models"
)

// Test ReadFile: Read existing file
func TestHandleReadFile_ExistingFile(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	input := models.ReadFileRequest{Path: tmpFile}
	toolResult, resp, err := HandleReadFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result, got %v", toolResult)
	}

	lines := strings.Split(resp.Content, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	// Validate format: "1:hash|content"
	validateLineFormat(t, lines[0], "line1")
	validateLineFormat(t, lines[1], "line2")
	validateLineFormat(t, lines[2], "line3")
}

// Test ReadFile: Empty file
func TestHandleReadFile_EmptyFile(t *testing.T) {
	tmpFile := createTempFile(t, "")
	defer os.Remove(tmpFile)

	input := models.ReadFileRequest{Path: tmpFile}
	toolResult, resp, err := HandleReadFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result, got %v", toolResult)
	}

	lines := strings.Split(resp.Content, "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line (empty), got %d", len(lines))
	}
}

// Test ReadFile: File with empty lines
func TestHandleReadFile_EmptyLines(t *testing.T) {
	tmpFile := createTempFile(t, "line1\n\nline3")
	defer os.Remove(tmpFile)

	input := models.ReadFileRequest{Path: tmpFile}
	toolResult, resp, err := HandleReadFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result, got %v", toolResult)
	}

	lines := strings.Split(resp.Content, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}

	validateLineFormat(t, lines[0], "line1")
	validateLineFormat(t, lines[1], "")
	validateLineFormat(t, lines[2], "line3")
}

// Test ReadFile: File not found
func TestHandleReadFile_FileNotFound(t *testing.T) {
	input := models.ReadFileRequest{Path: "/nonexistent/path/file.txt"}
	toolResult, resp, err := HandleReadFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "FILE_NOT_FOUND") {
		t.Errorf("expected FILE_NOT_FOUND error, got: %s", errorText)
	}

	if resp.Content != "" {
		t.Errorf("expected empty content in response")
	}
}

// Test ReadFile: Empty path
func TestHandleReadFile_EmptyPath(t *testing.T) {
	input := models.ReadFileRequest{Path: ""}
	toolResult, _, err := HandleReadFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_PATH") {
		t.Errorf("expected INVALID_PATH error")
	}
}

// Test WriteFile: Write new file
func TestHandleWriteFile_NewFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "newfile.txt")

	content := "hello\nworld"
	input := models.WriteFileRequest{Path: tmpFile, Content: content}
	toolResult, resp, err := HandleWriteFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !resp.Success {
		t.Errorf("expected success, got: %s", resp.Message)
	}

	// Verify file was created
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if string(data) != content {
		t.Errorf("expected content %q, got %q", content, string(data))
	}
}

// Test WriteFile: Overwrite existing file
func TestHandleWriteFile_OverwriteFile(t *testing.T) {
	tmpFile := createTempFile(t, "old content")
	defer os.Remove(tmpFile)

	newContent := "new content"
	input := models.WriteFileRequest{Path: tmpFile, Content: newContent}
	toolResult, resp, err := HandleWriteFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !resp.Success {
		t.Errorf("expected success")
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	if string(data) != newContent {
		t.Errorf("expected %q, got %q", newContent, string(data))
	}
}

// Test WriteFile: With valid hashes
func TestHandleWriteFile_WithValidHashes(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2")
	defer os.Remove(tmpFile)

	// Read file first to get hashes
	readInput := models.ReadFileRequest{Path: tmpFile}
	_, readResp, _ := HandleReadFile(context.Background(), nil, readInput)

	// Use the returned content with hashes for write
	writeInput := models.WriteFileRequest{Path: tmpFile, Content: readResp.Content}
	toolResult, resp, err := HandleWriteFile(context.Background(), nil, writeInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !resp.Success {
		t.Errorf("expected success: %s", resp.Message)
	}
}

// Test WriteFile: With invalid hashes
func TestHandleWriteFile_InvalidHash(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2")
	defer os.Remove(tmpFile)

	// Create content with invalid hash
	badHash := "bad"
	content := "1:" + badHash + "|line1\n2:dead|line2"
	writeInput := models.WriteFileRequest{Path: tmpFile, Content: content}
	toolResult, _, err := HandleWriteFile(context.Background(), nil, writeInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "HASH_MISMATCH") {
		t.Errorf("expected HASH_MISMATCH error, got: %s", errorText)
	}
}

// Test WriteFile: Empty path
func TestHandleWriteFile_EmptyPath(t *testing.T) {
	input := models.WriteFileRequest{Path: "", Content: "content"}
	toolResult, _, err := HandleWriteFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_PATH") {
		t.Errorf("expected INVALID_PATH error")
	}
}

// Test EditFile: Edit single line
func TestHandleEditFile_SingleLine(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	// Hash of the actual content (not the formatted version)
	startHash := hash.HashLine("line1")

	// Edit first line
	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartHash:  startHash,
		EndHash:    startHash,
		NewContent: "modified line1",
	}
	toolResult, resp, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !resp.Success {
		t.Errorf("expected success: %s", resp.Message)
	}

	// Verify file was modified
	data, _ := os.ReadFile(tmpFile)
	content := string(data)
	if !strings.HasPrefix(content, "modified line1") {
		t.Errorf("first line not modified correctly")
	}
}

// Test EditFile: Edit multiple lines
func TestHandleEditFile_MultipleLines(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3\nline4")
	defer os.Remove(tmpFile)

	// Hash of the actual content (not the formatted version)
	startHash := hash.HashLine("line1")
	endHash := hash.HashLine("line3")

	// Edit lines 1-3
	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartHash:  startHash,
		EndHash:    endHash,
		NewContent: "new",
	}
	toolResult, resp, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !resp.Success {
		t.Errorf("expected success: %s", resp.Message)
	}

	// Verify
	data, _ := os.ReadFile(tmpFile)
	content := string(data)
	if !strings.Contains(content, "new") || !strings.Contains(content, "line4") {
		t.Errorf("file not edited correctly: %s", content)
	}
}

// Test EditFile: Hash mismatch
func TestHandleEditFile_HashMismatch(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	// Try to edit with wrong hash
	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartHash:  "bad1",
		EndHash:    "bad2",
		NewContent: "new content",
	}
	toolResult, _, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "HASH_MISMATCH") {
		t.Errorf("expected HASH_MISMATCH error, got: %s", errorText)
	}

	// Verify file unchanged
	data, _ := os.ReadFile(tmpFile)
	if !strings.Contains(string(data), "line1") {
		t.Errorf("file was modified despite hash mismatch")
	}
}

// Test EditFile: Start hash not found
func TestHandleEditFile_StartHashNotFound(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartHash:  "nonexistent",
		EndHash:    "alsobad",
		NewContent: "new",
	}
	toolResult, _, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "HASH_MISMATCH") {
		t.Errorf("expected HASH_MISMATCH error")
	}
}

// Test EditFile: File not found
func TestHandleEditFile_FileNotFound(t *testing.T) {
	editInput := models.EditFileRequest{
		Path:       "/nonexistent/file.txt",
		StartHash:  "hash1",
		EndHash:    "hash2",
		NewContent: "content",
	}
	toolResult, _, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "FILE_NOT_FOUND") {
		t.Errorf("expected FILE_NOT_FOUND error")
	}
}

// Test EditFile: Empty path
func TestHandleEditFile_EmptyPath(t *testing.T) {
	editInput := models.EditFileRequest{
		Path:       "",
		StartHash:  "hash1",
		EndHash:    "hash2",
		NewContent: "content",
	}
	toolResult, _, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_PATH") {
		t.Errorf("expected INVALID_PATH error")
	}
}

// Test EditFile: File changed since read (stale hashes)
func TestHandleEditFile_StaleHashes(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	// Hash of original content
	startHash := hash.HashLine("line1")

	// Modify file to invalidate hashes
	os.WriteFile(tmpFile, []byte("changed\nline2\nline3"), 0644)

	// Try to edit with old hash (which no longer exists in file)
	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartHash:  startHash,
		EndHash:    startHash,
		NewContent: "new",
	}
	toolResult, _, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "HASH_MISMATCH") {
		t.Errorf("expected HASH_MISMATCH error for stale hashes")
	}
}

// Integration Test: Write → Edit → Read workflow
func TestIntegration_WriteEditRead(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "integration_test.txt")

	// Step 1: Write initial content
	initialContent := "line1\nline2\nline3"
	writeInput := models.WriteFileRequest{Path: tmpFile, Content: initialContent}
	toolResult, resp, err := HandleWriteFile(context.Background(), nil, writeInput)

	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected no error on write")
	}
	if !resp.Success {
		t.Fatalf("write not successful: %s", resp.Message)
	}

	// Step 2: Read file to get hashes
	readInput := models.ReadFileRequest{Path: tmpFile}
	toolResult, readResp, err := HandleReadFile(context.Background(), nil, readInput)

	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected no error on read")
	}

	readLines := strings.Split(readResp.Content, "\n")
	if len(readLines) != 3 {
		t.Fatalf("expected 3 lines from read, got %d", len(readLines))
	}

	// Extract hashes from read response
	line1Hash := hash.HashLine("line1")
	line2Hash := hash.HashLine("line2")

	// Step 3: Edit line 1 and 2
	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartHash:  line1Hash,
		EndHash:    line2Hash,
		NewContent: "modified1\nmodified2",
	}
	toolResult, editResp, err := HandleEditFile(context.Background(), nil, editInput)

	if err != nil {
		t.Fatalf("edit failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected no error on edit")
	}
	if !editResp.Success {
		t.Fatalf("edit not successful: %s", editResp.Message)
	}

	// Step 4: Read file again and validate changes
	toolResult, finalResp, err := HandleReadFile(context.Background(), nil, readInput)

	if err != nil {
		t.Fatalf("final read failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected no error on final read")
	}

	finalLines := strings.Split(finalResp.Content, "\n")
	if len(finalLines) != 3 {
		t.Fatalf("expected 3 lines in final read, got %d", len(finalLines))
	}

	// Validate the edits took effect
	validateLineFormat(t, finalLines[0], "modified1")
	validateLineFormat(t, finalLines[1], "modified2")
	validateLineFormat(t, finalLines[2], "line3") // line3 should be unchanged

	// Validate the content
	finalData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read file directly: %v", err)
	}

	expectedFinalContent := "modified1\nmodified2\nline3"
	if string(finalData) != expectedFinalContent {
		t.Errorf("expected final content %q, got %q", expectedFinalContent, string(finalData))
	}
}

// Helper: Create temp file with content
func createTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	return tmpFile.Name()
}

// Helper: Validate line format "linenum:hash|content"
func validateLineFormat(t *testing.T, formatted, expectedContent string) {
	// Extract hash and content
	parts := strings.SplitN(formatted, "|", 2)
	if len(parts) != 2 {
		t.Errorf("invalid format: %s", formatted)
		return
	}

	hashPart := parts[0]
	actualContent := parts[1]

	if actualContent != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, actualContent)
	}

	// Validate hash matches content
	expectedHash := hash.HashLine(expectedContent)
	colonIdx := strings.LastIndex(hashPart, ":")
	actualHash := hashPart[colonIdx+1:]

	if actualHash != expectedHash {
		t.Errorf("hash mismatch, expected %q, got %q", expectedHash, actualHash)
	}
}
