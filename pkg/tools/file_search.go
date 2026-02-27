package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	pkglog "awesomeProject/pkg/log"
	"awesomeProject/pkg/models"
	"awesomeProject/pkg/tools/filesearch"
)

func HandleSearchFile(ctx context.Context, req *mcp.CallToolRequest, input models.SearchFileRequest) (*mcp.CallToolResult, models.SearchFileResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Path == "" {
		return errorResult(ctx, "INVALID_PATH", "path cannot be empty"), models.SearchFileResponse{}, nil
	}
	if input.Query == "" {
		return errorResult(ctx, "INVALID_QUERY", "query cannot be empty"), models.SearchFileResponse{}, nil
	}

	logger.Info(ctx, "Searching files", "path", sanitizePath(input.Path), "regex", input.Regex)

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult(ctx, "INVALID_PATH", "invalid path: "+err.Error()), models.SearchFileResponse{}, nil
	}

	var opts *filesearch.SearchOptions
	if input.Regex {
		r, err := regexp.Compile(input.Query)
		if err != nil {
			return errorResult(ctx, "INVALID_REGEX", "bad regex: "+err.Error()), models.SearchFileResponse{}, nil
		}
		opts = &filesearch.SearchOptions{
			Kind:  filesearch.REGEX,
			Regex: r,
		}
	} else {
		opts = &filesearch.SearchOptions{
			Kind:   filesearch.LITERAL,
			Finder: filesearch.MakeStringFinder([]byte(input.Query)),
		}
	}

	resultChan := make(chan filesearch.SearchResult, 100)
	done := make(chan error, 1)

	go func() {
		done <- filesearch.Search([]string{absPath}, opts, input.Workers, resultChan)
	}()

	var results []models.SearchFileResult
	collecting := true
	for collecting {
		select {
		case result, ok := <-resultChan:
			if !ok {
				collecting = false
				break
			}
			results = append(results, models.SearchFileResult{
				File:    result.File,
				Line:    result.Line,
				Hash:    result.Hash,
				Content: result.Content,
			})
		case err := <-done:
			if err != nil {
				return errorResult(ctx, "SEARCH_FAILED", err.Error()), models.SearchFileResponse{}, nil
			}
			// Drain remaining results
			close(resultChan)
			for result := range resultChan {
				results = append(results, models.SearchFileResult{
					File:    result.File,
					Line:    result.Line,
					Hash:    result.Hash,
					Content: result.Content,
				})
			}
			collecting = false
		}
	}

	logger.Info(ctx, "Search completed",
		"matches", len(results),
		pkglog.Duration(time.Since(start)))

	resp := models.SearchFileResponse{
		Results: results,
		Message: fmt.Sprintf("Found %d matches", len(results)),
	}
	return nil, resp, nil
}

func HandleFindFiles(ctx context.Context, req *mcp.CallToolRequest, input models.FindFilesRequest) (*mcp.CallToolResult, models.FindFilesResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Path == "" {
		return errorResult(ctx, "INVALID_PATH", "path cannot be empty"), models.FindFilesResponse{}, nil
	}
	if input.Pattern == "" {
		return errorResult(ctx, "INVALID_PATTERN", "pattern cannot be empty"), models.FindFilesResponse{}, nil
	}

	logger.Info(ctx, "Finding files", "path", sanitizePath(input.Path), "pattern", input.Pattern)

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult(ctx, "INVALID_PATH", "invalid path: "+err.Error()), models.FindFilesResponse{}, nil
	}

	matches, err := filesearch.FindFiles(absPath, input.Pattern)
	if err != nil {
		return errorResult(ctx, "FIND_FAILED", err.Error()), models.FindFilesResponse{}, nil
	}

	var files []models.FindFilesResult
	for _, match := range matches {
		files = append(files, models.FindFilesResult{
			Path:     match.Path,
			Size:     match.Size,
			Modified: match.Modified,
		})
	}

	logger.Info(ctx, "Find completed",
		"matches", len(files),
		pkglog.Duration(time.Since(start)))

	resp := models.FindFilesResponse{
		Files:   files,
		Message: fmt.Sprintf("Found %d files", len(files)),
	}
	return nil, resp, nil
}

func HandleListDirectory(ctx context.Context, req *mcp.CallToolRequest, input models.ListDirectoryRequest) (*mcp.CallToolResult, models.ListDirectoryResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Path == "" {
		return errorResult(ctx, "INVALID_PATH", "path cannot be empty"), models.ListDirectoryResponse{}, nil
	}

	logger.Info(ctx, "Listing directory", "path", sanitizePath(input.Path))

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult(ctx, "INVALID_PATH", "invalid path: "+err.Error()), models.ListDirectoryResponse{}, nil
	}

	entries, err := filesearch.ListDirectory(absPath)
	if err != nil {
		return errorResult(ctx, "LIST_FAILED", err.Error()), models.ListDirectoryResponse{}, nil
	}

	var dirEntries []models.ListDirectoryEntry
	for _, entry := range entries {
		dirEntries = append(dirEntries, models.ListDirectoryEntry{
			Name:        entry.Name,
			Type:        entry.Type,
			Size:        entry.Size,
			Permissions: entry.Permissions,
		})
	}

	logger.Info(ctx, "List completed",
		"entries", len(dirEntries),
		pkglog.Duration(time.Since(start)))

	resp := models.ListDirectoryResponse{
		Entries: dirEntries,
		Message: fmt.Sprintf("Listed %d entries", len(dirEntries)),
	}
	return nil, resp, nil
}

func HandleTreeDirectory(ctx context.Context, req *mcp.CallToolRequest, input models.TreeDirectoryRequest) (*mcp.CallToolResult, models.TreeDirectoryResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Path == "" {
		return errorResult(ctx, "INVALID_PATH", "path cannot be empty"), models.TreeDirectoryResponse{}, nil
	}

	logger.Info(ctx, "Generating tree", "path", sanitizePath(input.Path), "max_depth", input.MaxDepth)

	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return errorResult(ctx, "INVALID_PATH", "invalid path: "+err.Error()), models.TreeDirectoryResponse{}, nil
	}

	tree, err := filesearch.TreeDirectory(absPath, input.MaxDepth)
	if err != nil {
		return errorResult(ctx, "TREE_FAILED", err.Error()), models.TreeDirectoryResponse{}, nil
	}

	logger.Info(ctx, "Tree completed",
		pkglog.Duration(time.Since(start)))

	resp := models.TreeDirectoryResponse{
		Tree:    tree,
		Message: "Tree generated successfully",
	}
	return nil, resp, nil
}
