package tools

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func successResult(data interface{}) (*mcp.CallToolResult, any, error) {
	jsonData, _ := json.Marshal(data)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

func errorResult(code, message string) (*mcp.CallToolResult, any, error) {
	errResp := map[string]string{
		"code":    code,
		"message": message,
	}
	jsonData, _ := json.Marshal(errResp)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}
