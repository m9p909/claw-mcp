package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/hash"
	"awesomeProject/pkg/models"
)

func HandleReadFile(ctx context.Context, req *mcp.CallToolRequest, input models.ReadFileRequest) (*mcp.CallToolResult, models.ReadFileResponse, error) {
	if input.Path == "" {
		return errorResult("INVALID_PATH", "path cannot be empty"), models.ReadFileResponse{}, nil
	}

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult("INVALID_PATH", "invalid path: "+err.Error()), models.ReadFileResponse{}, nil
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("FILE_NOT_FOUND", "file not found"), models.ReadFileResponse{}, nil
		}
		return errorResult("READ_FAILED", "read failed: "+err.Error()), models.ReadFileResponse{}, nil
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
	return nil, resp, nil
}

func HandleWriteFile(ctx context.Context, req *mcp.CallToolRequest, input models.WriteFileRequest) (*mcp.CallToolResult, models.WriteFileResponse, error) {

	if input.Path == "" {
		return errorResult("INVALID_PATH", "path cannot be empty"), models.WriteFileResponse{}, nil
	}

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult("INVALID_PATH", "invalid path: "+err.Error()), models.WriteFileResponse{}, nil
	}

	// If content has hashes, validate and strip them
	lines := strings.Split(input.Content, "\n")
	var actualLines []string
	for _, line := range lines {
		if strings.Contains(line, "|") {
			// Extract hash and validate
			extractedHash, err := hash.ExtractHashFromLine(line)
			if err != nil {
				return errorResult("INVALID_REQUEST", "invalid hash format: "+err.Error()), models.WriteFileResponse{}, nil
			}

			parts := strings.SplitN(line, "|", 2)
			if len(parts) != 2 {
				return errorResult("INVALID_REQUEST", "invalid line format"), models.WriteFileResponse{}, nil
			}

			lineContent := parts[1]

			// Validate hash
			if !hash.ValidateHash(lineContent, extractedHash) {
				return errorResult("HASH_MISMATCH", "hash mismatch on line"), models.WriteFileResponse{}, nil
			}

			actualLines = append(actualLines, lineContent)
		} else {
			actualLines = append(actualLines, line)
		}
	}

	content := strings.Join(actualLines, "\n")

	// Write file
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return errorResult("WRITE_FAILED", "write failed: "+err.Error()), models.WriteFileResponse{}, nil
	}

	resp := models.WriteFileResponse{
		Success: true,
		Message: "File written successfully",
	}
	return nil, resp, nil
}

func HandleEditFile(ctx context.Context, req *mcp.CallToolRequest, input models.EditFileRequest) (*mcp.CallToolResult, models.EditFileResponse, error) {

	if input.Path == "" {
		return errorResult("INVALID_PATH", "path cannot be empty"), models.EditFileResponse{}, nil
	}

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult("INVALID_PATH", "invalid path: "+err.Error()), models.EditFileResponse{}, nil
	}

	// Read current file
	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("FILE_NOT_FOUND", "file not found"), models.EditFileResponse{}, nil
		}
		return errorResult("READ_FAILED", "read failed: "+err.Error()), models.EditFileResponse{}, nil
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
		return errorResult("HASH_MISMATCH", "start hash not found"), models.EditFileResponse{}, nil
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
		return errorResult("HASH_MISMATCH", "end hash not found"), models.EditFileResponse{}, nil
	}

	// Replace lines
	newLines := strings.Split(input.NewContent, "\n")
	result := append(lines[:startIdx], newLines...)
	result = append(result, lines[endIdx+1:]...)

	// Write back
	finalContent := strings.Join(result, "\n")
	if err := os.WriteFile(absPath, []byte(finalContent), 0644); err != nil {
		return errorResult("EDIT_FAILED", "edit failed: "+err.Error()), models.EditFileResponse{}, nil
	}

	resp := models.EditFileResponse{
		Success: true,
		Message: "File edited successfully",
	}
	return nil, resp, nil
}
