package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

func HandleWriteMemory(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	argsJSON, _ := json.Marshal(args)
	var input models.WriteMemoryRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	if input.Category == "" || input.Content == "" {
		return errorResult("INVALID_REQUEST", "category and content cannot be empty")
	}

	if err := storage.WriteMemory(input.Category, input.Content); err != nil {
		return errorResult("INTERNAL_ERROR", "failed to write memory: "+err.Error())
	}

	resp := models.WriteMemoryResponse{
		Success: true,
		Message: "Memory written successfully",
	}
	return successResult(resp)
}

func HandleQueryMemory(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	argsJSON, _ := json.Marshal(args)
	var input models.QueryMemoryRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	if input.Query == "" {
		return errorResult("INVALID_REQUEST", "query cannot be empty")
	}

	results, err := storage.QueryMemory(input.Query)
	if err != nil {
		return errorResult("QUERY_FAILED", "query failed: "+err.Error())
	}

	resp := models.QueryMemoryResponse{
		Results: results,
		Message: fmt.Sprintf("Found %d results", len(results)),
	}
	return successResult(resp)
}

func HandleMemorySearch(ctx context.Context, req *mcp.CallToolRequest, args interface{}) (*mcp.CallToolResult, any, error) {
	argsJSON, _ := json.Marshal(args)
	var input models.SearchMemoryRequest
	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return errorResult("INVALID_REQUEST", "invalid request: "+err.Error())
	}

	if input.Query == "" {
		return errorResult("INVALID_REQUEST", "query cannot be empty")
	}

	results, err := storage.SearchMemory(input.Query, input.Limit)
	if err != nil {
		return errorResult("SEARCH_FAILED", "search failed: "+err.Error())
	}

	resp := models.SearchMemoryResponse{
		Results: results,
		Message: fmt.Sprintf("Found %d results", len(results)),
	}
	return successResult(resp)
}
