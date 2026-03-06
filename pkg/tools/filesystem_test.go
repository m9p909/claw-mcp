package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

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

	validateLineFormat(t, lines[0], 1, "line1")
	validateLineFormat(t, lines[1], 2, "line2")
	validateLineFormat(t, lines[2], 3, "line3")
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
	validateLineFormat(t, lines[0], 1, "")
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

	validateLineFormat(t, lines[0], 1, "line1")
	validateLineFormat(t, lines[1], 2, "")
	validateLineFormat(t, lines[2], 3, "line3")
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

// Test WriteFile: Content with pipe characters written verbatim
func TestHandleWriteFile_PipeCharsVerbatim(t *testing.T) {
	tmpFile := createTempFile(t, "")
	defer os.Remove(tmpFile)

	content := "1|line one\n2|line two"
	input := models.WriteFileRequest{Path: tmpFile, Content: content}
	toolResult, resp, err := HandleWriteFile(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !resp.Success {
		t.Errorf("expected success: %s", resp.Message)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	if string(data) != content {
		t.Errorf("expected verbatim content %q, got %q", content, string(data))
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

	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartLine:  1,
		EndLine:    1,
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

	data, _ := os.ReadFile(tmpFile)
	content := string(data)
	if !strings.HasPrefix(content, "modified line1") {
		t.Errorf("first line not modified correctly: %s", content)
	}
}

// Test EditFile: Edit multiple lines
func TestHandleEditFile_MultipleLines(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3\nline4")
	defer os.Remove(tmpFile)

	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartLine:  1,
		EndLine:    3,
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

	data, _ := os.ReadFile(tmpFile)
	content := string(data)
	if !strings.Contains(content, "new") || !strings.Contains(content, "line4") {
		t.Errorf("file not edited correctly: %s", content)
	}
}

// Test EditFile: Out of range
func TestHandleEditFile_OutOfRange(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartLine:  2,
		EndLine:    10,
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
	if !strings.Contains(errorText, "INVALID_RANGE") {
		t.Errorf("expected INVALID_RANGE error, got: %s", errorText)
	}

	data, _ := os.ReadFile(tmpFile)
	if !strings.Contains(string(data), "line1") {
		t.Errorf("file was modified despite out-of-range")
	}
}

// Test EditFile: start_line greater than end_line
func TestHandleEditFile_StartGreaterThanEnd(t *testing.T) {
	tmpFile := createTempFile(t, "line1\nline2\nline3")
	defer os.Remove(tmpFile)

	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartLine:  3,
		EndLine:    1,
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
	if !strings.Contains(errorText, "INVALID_RANGE") {
		t.Errorf("expected INVALID_RANGE error, got: %s", errorText)
	}
}

// Test EditFile: File not found
func TestHandleEditFile_FileNotFound(t *testing.T) {
	editInput := models.EditFileRequest{
		Path:       "/nonexistent/file.txt",
		StartLine:  1,
		EndLine:    1,
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
		StartLine:  1,
		EndLine:    1,
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

// Integration Test: Write → Edit → Read workflow
func TestIntegration_WriteEditRead(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "integration_test.txt")

	// Step 1: Write initial content
	writeInput := models.WriteFileRequest{Path: tmpFile, Content: "line1\nline2\nline3"}
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

	// Step 2: Edit lines 1-2
	editInput := models.EditFileRequest{
		Path:       tmpFile,
		StartLine:  1,
		EndLine:    2,
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

	// Step 3: Read and validate
	readInput := models.ReadFileRequest{Path: tmpFile}
	toolResult, readResp, err := HandleReadFile(context.Background(), nil, readInput)

	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected no error on read")
	}

	finalLines := strings.Split(readResp.Content, "\n")
	if len(finalLines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(finalLines))
	}

	validateLineFormat(t, finalLines[0], 1, "modified1")
	validateLineFormat(t, finalLines[1], 2, "modified2")
	validateLineFormat(t, finalLines[2], 3, "line3")

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read file directly: %v", err)
	}
	if string(data) != "modified1\nmodified2\nline3" {
		t.Errorf("unexpected final content: %q", string(data))
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

// Helper: Validate line format "linenum|content"
func validateLineFormat(t *testing.T, formatted string, expectedLineNum int, expectedContent string) {
	t.Helper()
	parts := strings.SplitN(formatted, "|", 2)
	if len(parts) != 2 {
		t.Errorf("invalid format (missing |): %q", formatted)
		return
	}

	if parts[0] != fmt.Sprintf("%d", expectedLineNum) {
		t.Errorf("expected line number %d, got %q", expectedLineNum, parts[0])
	}
	if parts[1] != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, parts[1])
	}
}
