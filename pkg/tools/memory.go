package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

func HandleWriteMemory(ctx context.Context, req *mcp.CallToolRequest, input models.WriteMemoryRequest) (*mcp.CallToolResult, models.WriteMemoryResponse, error) {
	if input.Category == "" || input.Content == "" {
		return errorResult("INVALID_REQUEST", "category and content cannot be empty"), models.WriteMemoryResponse{}, nil
	}

	if err := storage.WriteMemory(input.Category, input.Content); err != nil {
		return errorResult("INTERNAL_ERROR", "failed to write memory: "+err.Error()), models.WriteMemoryResponse{}, nil
	}

	resp := models.WriteMemoryResponse{
		Success: true,
		Message: "Memory written successfully",
	}
	return nil, resp, nil
}

func HandleQueryMemory(ctx context.Context, req *mcp.CallToolRequest, input models.QueryMemoryRequest) (*mcp.CallToolResult, models.QueryMemoryResponse, error) {
	if input.Query == "" {
		return errorResult("INVALID_REQUEST", "query cannot be empty"), models.QueryMemoryResponse{}, nil
	}

	results, err := storage.QueryMemory(input.Query)
	if err != nil {
		return errorResult("QUERY_FAILED", "query failed: "+err.Error()), models.QueryMemoryResponse{}, nil
	}

	resp := models.QueryMemoryResponse{
		Results: results,
		Message: fmt.Sprintf("Found %d results", len(results)),
	}
	return nil, resp, nil
}

func HandleMemorySearch(ctx context.Context, req *mcp.CallToolRequest, input models.SearchMemoryRequest) (*mcp.CallToolResult, models.SearchMemoryResponse, error) {
	if input.Query == "" {
		return errorResult("INVALID_REQUEST", "query cannot be empty"), models.SearchMemoryResponse{}, nil
	}

	results, err := storage.SearchMemory(input.Query, input.Limit)
	if err != nil {
		return errorResult("SEARCH_FAILED", "search failed: "+err.Error()), models.SearchMemoryResponse{}, nil
	}

	resp := models.SearchMemoryResponse{
		Results: results,
		Message: fmt.Sprintf("Found %d results", len(results)),
	}
	return nil, resp, nil
}
