package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

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

	lines := strings.Split(string(data), "\n")
	formattedLines := make([]string, len(lines))
	for i, line := range lines {
		formattedLines[i] = fmt.Sprintf("%d|%s", i+1, line)
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

	// Write file
	if err := os.WriteFile(absPath, []byte(input.Content), 0644); err != nil {
		return errorResult(ctx, "WRITE_FAILED", "write failed: "+err.Error()), models.WriteFileResponse{}, nil
	}

	logger.Info(ctx, "File write completed",
		"bytes_written", len(input.Content),
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

	startIdx := input.StartLine - 1
	endIdx := input.EndLine - 1

	if input.StartLine < 1 || input.EndLine < input.StartLine || input.EndLine > len(lines) {
		return errorResult(ctx, "INVALID_RANGE", fmt.Sprintf("invalid range %d-%d for file with %d lines", input.StartLine, input.EndLine, len(lines))), models.EditFileResponse{}, nil
	}

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
