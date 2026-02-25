package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	pkglog "awesomeProject/pkg/log"
	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

func HandleWriteMemory(ctx context.Context, req *mcp.CallToolRequest, input models.WriteMemoryRequest) (*mcp.CallToolResult, models.WriteMemoryResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Category == "" {
		return errorResult(ctx, "INVALID_REQUEST", "category is required"), models.WriteMemoryResponse{}, nil
	}

	if input.Content == "" {
		return errorResult(ctx, "INVALID_REQUEST", "content is required"), models.WriteMemoryResponse{}, nil
	}

	logger.Info(ctx, "Writing memory", "category", input.Category)
	logger.Debug(ctx, "Memory content size", "bytes", len(input.Content))

	if err := storage.WriteMemory(input.Category, input.Content); err != nil {
		return errorResult(ctx, "INTERNAL_ERROR", "failed to write memory: "+err.Error()), models.WriteMemoryResponse{}, nil
	}

	logger.Info(ctx, "Memory write completed",
		"category", input.Category,
		"content_size", len(input.Content),
		pkglog.Duration(time.Since(start)))

	resp := models.WriteMemoryResponse{
		Success: true,
		Message: "Memory written successfully",
	}
	return nil, resp, nil
}

func HandleQueryMemory(ctx context.Context, req *mcp.CallToolRequest, input models.QueryMemoryRequest) (*mcp.CallToolResult, models.QueryMemoryResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Query == "" {
		return errorResult(ctx, "INVALID_REQUEST", "query cannot be empty"), models.QueryMemoryResponse{}, nil
	}

	logger.Info(ctx, "Executing memory query")
	logger.Debug(ctx, "Query string structure", "has_query", input.Query != "")

	results, err := storage.QueryMemory(input.Query)
	if err != nil {
		return errorResult(ctx, "QUERY_FAILED", "query failed: "+err.Error()), models.QueryMemoryResponse{}, nil
	}

	logger.Info(ctx, "Memory query completed",
		"result_count", len(results),
		pkglog.Duration(time.Since(start)))

	resp := models.QueryMemoryResponse{
		Results: results,
		Message: fmt.Sprintf("Found %d results", len(results)),
	}
	return nil, resp, nil
}

func HandleMemorySearch(ctx context.Context, req *mcp.CallToolRequest, input models.SearchMemoryRequest) (*mcp.CallToolResult, models.SearchMemoryResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Query == "" {
		return errorResult(ctx, "INVALID_REQUEST", "query cannot be empty"), models.SearchMemoryResponse{}, nil
	}

	logger.Info(ctx, "Searching memory")

	if input.Limit == 0 {
		logger.Debug(ctx, "Memory search limit", "limit", "unlimited")
	} else {
		logger.Debug(ctx, "Memory search limit", "limit", input.Limit)
	}

	results, err := storage.SearchMemory(input.Query, input.Limit)
	if err != nil {
		return errorResult(ctx, "SEARCH_FAILED", "search failed: "+err.Error()), models.SearchMemoryResponse{}, nil
	}

	logger.Info(ctx, "Memory search completed",
		"result_count", len(results),
		"limit", input.Limit,
		pkglog.Duration(time.Since(start)))

	resp := models.SearchMemoryResponse{
		Results: results,
		Message: fmt.Sprintf("Found %d results", len(results)),
	}
	return nil, resp, nil
}
