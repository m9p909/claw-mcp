package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/hash"
	"awesomeProject/pkg/models"
)

func HandleReadFile(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	// Parse args to ReadFileRequest
	argsJSON, _ := json.Marshal(args)
	var input models.ReadFileRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	if input.Path == "" {
		return errorResult("INVALID_PATH", "path cannot be empty")
	}

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult("INVALID_PATH", "invalid path: "+err.Error())
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("FILE_NOT_FOUND", "file not found")
		}
		return errorResult("READ_FAILED", "read failed: "+err.Error())
	}

	// Format content with hashes
	lines := strings.Split(string(data), "\n")
	var formattedLines []string
	for i, line := range lines {
		lineNum := i + 1
		formatted := hash.FormatLineWithHash(lineNum, line)
		formattedLines = append(formattedLines, formatted)
	}
	content := strings.Join(formattedLines, "\n")

	resp := models.ReadFileResponse{Content: content}
	return successResult(resp)
}

func HandleWriteFile(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	argsJSON, _ := json.Marshal(args)
	var input models.WriteFileRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	if input.Path == "" {
		return errorResult("INVALID_PATH", "path cannot be empty")
	}

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult("INVALID_PATH", "invalid path: "+err.Error())
	}

	// If content has hashes, validate and strip them
	lines := strings.Split(input.Content, "\n")
	var actualLines []string
	for _, line := range lines {
		if strings.Contains(line, "|") {
			// Extract hash and validate
			extractedHash, err := hash.ExtractHashFromLine(line)
			if err != nil {
				return errorResult("INVALID_REQUEST", "invalid hash format: "+err.Error())
			}

			parts := strings.SplitN(line, "|", 2)
			if len(parts) != 2 {
				return errorResult("INVALID_REQUEST", "invalid line format")
			}

			lineContent := parts[1]

			// Validate hash
			if !hash.ValidateHash(lineContent, extractedHash) {
				return errorResult("HASH_MISMATCH", "hash mismatch on line")
			}

			actualLines = append(actualLines, lineContent)
		} else {
			actualLines = append(actualLines, line)
		}
	}

	content := strings.Join(actualLines, "\n")

	// Write file
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return errorResult("WRITE_FAILED", "write failed: "+err.Error())
	}

	resp := models.WriteFileResponse{
		Success: true,
		Message: "File written successfully",
	}
	return successResult(resp)
}

func HandleEditFile(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	argsJSON, _ := json.Marshal(args)
	var input models.EditFileRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	if input.Path == "" {
		return errorResult("INVALID_PATH", "path cannot be empty")
	}

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult("INVALID_PATH", "invalid path: "+err.Error())
	}

	// Read current file
	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("FILE_NOT_FOUND", "file not found")
		}
		return errorResult("READ_FAILED", "read failed: "+err.Error())
	}

	lines := strings.Split(string(data), "\n")

	// Find start and end line indices by hash
	var startIdx, endIdx int
	found := false
	for i, line := range lines {
		lineHash := hash.HashLine(line)
		if lineHash == input.StartHash {
			startIdx = i
			found = true
			break
		}
	}

	if !found {
		return errorResult("HASH_MISMATCH", "start hash not found")
	}

	found = false
	for i := startIdx; i < len(lines); i++ {
		lineHash := hash.HashLine(lines[i])
		if lineHash == input.EndHash {
			endIdx = i
			found = true
			break
		}
	}

	if !found {
		return errorResult("HASH_MISMATCH", "end hash not found")
	}

	// Replace lines
	newLines := strings.Split(input.NewContent, "\n")
	result := append(lines[:startIdx], newLines...)
	result = append(result, lines[endIdx+1:]...)

	// Write back
	finalContent := strings.Join(result, "\n")
	if err := os.WriteFile(absPath, []byte(finalContent), 0644); err != nil {
		return errorResult("EDIT_FAILED", "edit failed: "+err.Error())
	}

	resp := models.EditFileResponse{
		Success: true,
		Message: "File edited successfully",
	}
	return successResult(resp)
}
