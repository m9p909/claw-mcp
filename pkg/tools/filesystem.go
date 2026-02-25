package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/hash"
	pkglog "awesomeProject/pkg/log"
	"awesomeProject/pkg/models"
)

func HandleReadFile(ctx context.Context, req *mcp.CallToolRequest, input models.ReadFileRequest) (*mcp.CallToolResult, models.ReadFileResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Path == "" {
		return errorResult(ctx, "INVALID_PATH", "path cannot be empty"), models.ReadFileResponse{}, nil
	}

	logger.Info(ctx, "Reading file", "path", sanitizePath(input.Path))

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult(ctx, "INVALID_PATH", "invalid path: "+err.Error()), models.ReadFileResponse{}, nil
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult(ctx, "FILE_NOT_FOUND", "file not found"), models.ReadFileResponse{}, nil
		}
		return errorResult(ctx, "READ_FAILED", "read failed: "+err.Error()), models.ReadFileResponse{}, nil
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

	logger.Info(ctx, "File read completed",
		"bytes_read", len(data),
		"lines", len(lines),
		pkglog.Duration(time.Since(start)))

	resp := models.ReadFileResponse{Content: content}
	return nil, resp, nil
}

func HandleWriteFile(ctx context.Context, req *mcp.CallToolRequest, input models.WriteFileRequest) (*mcp.CallToolResult, models.WriteFileResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Path == "" {
		return errorResult(ctx, "INVALID_PATH", "path cannot be empty"), models.WriteFileResponse{}, nil
	}

	logger.Info(ctx, "Writing file", "path", sanitizePath(input.Path))

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult(ctx, "INVALID_PATH", "invalid path: "+err.Error()), models.WriteFileResponse{}, nil
	}

	logger.Debug(ctx, "Path validation", "resolved", "success")

	// If content has hashes, validate and strip them
	lines := strings.Split(input.Content, "\n")
	var actualLines []string
	for _, line := range lines {
		if strings.Contains(line, "|") {
			// Extract hash and validate
			extractedHash, err := hash.ExtractHashFromLine(line)
			if err != nil {
				return errorResult(ctx, "INVALID_REQUEST", "invalid hash format: "+err.Error()), models.WriteFileResponse{}, nil
			}

			parts := strings.SplitN(line, "|", 2)
			if len(parts) != 2 {
				return errorResult(ctx, "INVALID_REQUEST", "invalid line format"), models.WriteFileResponse{}, nil
			}

			lineContent := parts[1]

			// Validate hash
			if !hash.ValidateHash(lineContent, extractedHash) {
				return errorResult(ctx, "HASH_MISMATCH", "hash mismatch on line"), models.WriteFileResponse{}, nil
			}

			actualLines = append(actualLines, lineContent)
		} else {
			actualLines = append(actualLines, line)
		}
	}

	content := strings.Join(actualLines, "\n")

	logger.Debug(ctx, "Hash validation", "result", "success")

	// Write file
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return errorResult(ctx, "WRITE_FAILED", "write failed: "+err.Error()), models.WriteFileResponse{}, nil
	}

	logger.Info(ctx, "File write completed",
		"bytes_written", len(content),
		pkglog.Duration(time.Since(start)))

	resp := models.WriteFileResponse{
		Success: true,
		Message: "File written successfully",
	}
	return nil, resp, nil
}

func HandleEditFile(ctx context.Context, req *mcp.CallToolRequest, input models.EditFileRequest) (*mcp.CallToolResult, models.EditFileResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Path == "" {
		return errorResult(ctx, "INVALID_PATH", "path cannot be empty"), models.EditFileResponse{}, nil
	}

	logger.Info(ctx, "Editing file", "path", sanitizePath(input.Path))

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult(ctx, "INVALID_PATH", "invalid path: "+err.Error()), models.EditFileResponse{}, nil
	}

	// Read current file
	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult(ctx, "FILE_NOT_FOUND", "file not found"), models.EditFileResponse{}, nil
		}
		return errorResult(ctx, "READ_FAILED", "read failed: "+err.Error()), models.EditFileResponse{}, nil
	}

	lines := strings.Split(string(data), "\n")

	// Format lines with hashes to match what was returned to the user
	var formattedLines []string
	for i, line := range lines {
		formatted := hash.FormatLineWithHash(i+1, line)
		formattedLines = append(formattedLines, formatted)
	}

	// Find start and end line indices by hash
	var startIdx, endIdx int
	found := false
	for i, formatted := range formattedLines {
		lineHash := hash.HashLine(formatted)
		if lineHash == input.StartHash {
			startIdx = i
			found = true
			break
		}
	}

	if !found {
		logger.Debug(ctx, "Start hash lookup", "line_count", len(lines), "result", "not_found")
		return errorResult(ctx, "HASH_MISMATCH", "start hash not found"), models.EditFileResponse{}, nil
	}

	logger.Debug(ctx, "Start hash lookup", "result", "found", "line_num", startIdx+1)

	found = false
	for i := startIdx; i < len(formattedLines); i++ {
		lineHash := hash.HashLine(formattedLines[i])
		if lineHash == input.EndHash {
			endIdx = i
			found = true
			break
		}
	}

	if !found {
		logger.Debug(ctx, "End hash lookup", "range_start", startIdx, "range_end", len(lines), "result", "not_found")
		return errorResult(ctx, "HASH_MISMATCH", "end hash not found"), models.EditFileResponse{}, nil
	}

	logger.Debug(ctx, "End hash lookup", "result", "found", "line_num", endIdx+1)

	// Replace lines
	newLines := strings.Split(input.NewContent, "\n")
	linesReplaced := endIdx - startIdx + 1
	result := append(lines[:startIdx], newLines...)
	result = append(result, lines[endIdx+1:]...)

	logger.Debug(ctx, "Lines replacement", "lines_replaced", linesReplaced)

	// Write back
	finalContent := strings.Join(result, "\n")
	if err := os.WriteFile(absPath, []byte(finalContent), 0644); err != nil {
		return errorResult(ctx, "EDIT_FAILED", "edit failed: "+err.Error()), models.EditFileResponse{}, nil
	}

	logger.Info(ctx, "File edit completed",
		"lines_replaced", linesReplaced,
		pkglog.Duration(time.Since(start)))

	resp := models.EditFileResponse{
		Success: true,
		Message: "File edited successfully",
	}
	return nil, resp, nil
}
