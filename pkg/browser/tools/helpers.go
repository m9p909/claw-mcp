package tools

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/browser"
)

// errorResult creates a standardized error response
func errorResult(code, message string) *mcp.CallToolResult {
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

// FormatError creates a formatted error using browser error formatter
func FormatError(err error) string {
	return browser.FormatPlaywrightError(err)
}
