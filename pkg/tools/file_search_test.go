package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/models"
)

func TestHandleSearchFile_PathValidation(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("empty path", func(t *testing.T) {
		input := models.SearchFileRequest{
			Path:  "",
			Query: "test",
		}
		result, _, _ := HandleSearchFile(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for empty path")
		}
	})

	t.Run("empty query", func(t *testing.T) {
		input := models.SearchFileRequest{
			Path:  "/tmp",
			Query: "",
		}
		result, _, _ := HandleSearchFile(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for empty query")
		}
	})

	t.Run("invalid regex", func(t *testing.T) {
		input := models.SearchFileRequest{
			Path:  "/tmp",
			Query: "[invalid(",
			Regex: true,
		}
		result, _, _ := HandleSearchFile(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for invalid regex")
		}
	})
}

func TestHandleFindFiles_PathValidation(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("empty path", func(t *testing.T) {
		input := models.FindFilesRequest{
			Path:    "",
			Pattern: "*.txt",
		}
		result, _, _ := HandleFindFiles(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for empty path")
		}
	})

	t.Run("empty pattern", func(t *testing.T) {
		input := models.FindFilesRequest{
			Path:    "/tmp",
			Pattern: "",
		}
		result, _, _ := HandleFindFiles(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for empty pattern")
		}
	})
}

func TestHandleListDirectory_PathValidation(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("empty path", func(t *testing.T) {
		input := models.ListDirectoryRequest{
			Path: "",
		}
		result, _, _ := HandleListDirectory(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for empty path")
		}
	})

	t.Run("non-existent path", func(t *testing.T) {
		input := models.ListDirectoryRequest{
			Path: "/nonexistent/path/that/does/not/exist",
		}
		result, _, _ := HandleListDirectory(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for non-existent path")
		}
	})
}

func TestHandleTreeDirectory_PathValidation(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("empty path", func(t *testing.T) {
		input := models.TreeDirectoryRequest{
			Path: "",
		}
		result, _, _ := HandleTreeDirectory(ctx, req, input)
		if result == nil {
			t.Error("Should return error result for empty path")
		}
	})
}

func TestHandleSearchFile_SuccessfulSearch(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("line with pattern\nanother line"), 0644); err != nil {
		t.Fatal(err)
	}

	input := models.SearchFileRequest{
		Path:    testFile,
		Query:   "pattern",
		Regex:   false,
		Workers: 1,
	}

	result, resp, err := HandleSearchFile(ctx, req, input)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Error("Should not return error result for successful search")
	}
	if len(resp.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Line != 1 {
		t.Errorf("Expected line 1, got %d", resp.Results[0].Line)
	}
	if resp.Results[0].Hash == "" {
		t.Error("Result should have hash")
	}
}
