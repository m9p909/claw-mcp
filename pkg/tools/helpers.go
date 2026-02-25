package tools

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// errorResult returns an error result as a CallToolResult
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
