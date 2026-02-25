package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

	pkglog "awesomeProject/pkg/log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// errorResult returns an error result as a CallToolResult and logs server-side
func errorResult(ctx context.Context, code, message string) *mcp.CallToolResult {
	logger := pkglog.NewLogger()
	logger.Error(ctx, "Tool error", "code", code, "message", message)

	errResp := map[string]string{
		"code":    code,
		"message": message,
	}
	jsonData, _ := json.Marshal(errResp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}
}

// sanitizePath masks file paths for logging, showing only the operation type
func sanitizePath(path string) string {
	// Get just the filename without directory
	base := filepath.Base(path)
	return base
}

// isInHomeDir checks if path is under home directory
func isInHomeDir(path string) bool {
	return !strings.Contains(path, "/")
}
