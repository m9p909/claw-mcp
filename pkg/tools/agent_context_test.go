package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"awesomeProject/pkg/models"
)

func TestHandleGetAgentContext_Success(t *testing.T) {
	// Create temp docs dir with AGENT_CONTEXT.md
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.Mkdir(docsDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}

	content := "# Agent Context Guide\n\nYou are connected to Claw..."
	docPath := filepath.Join(docsDir, "AGENT_CONTEXT.md")
	if err := os.WriteFile(docPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write AGENT_CONTEXT.md: %v", err)
	}

	// Change to temp dir
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	// Execute
	ctx := context.Background()
	result, resp, err := HandleGetAgentContext(ctx, nil, models.GetAgentContextRequest{})

	// Verify
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got error result")
	}
	if resp.Content != content {
		t.Errorf("expected content %q, got %q", content, resp.Content)
	}
}

func TestHandleGetAgentContext_FileNotFound(t *testing.T) {
	// Create temp dir without docs/AGENT_CONTEXT.md
	tmpDir := t.TempDir()

	// Change to temp dir
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	// Execute
	ctx := context.Background()
	result, resp, err := HandleGetAgentContext(ctx, nil, models.GetAgentContextRequest{})

	// Verify error result
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected error result, got nil")
	}
	if resp.Content != "" {
		t.Errorf("expected empty content on error, got %q", resp.Content)
	}

	// Check error contains expected message
	resultJSON := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(resultJSON, "DOC_NOT_FOUND") {
		t.Errorf("expected DOC_NOT_FOUND error code, got: %s", resultJSON)
	}
}

func TestHandleGetAgentContext_ContentSize(t *testing.T) {
	// Create temp docs dir with large content
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.Mkdir(docsDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}

	content := strings.Repeat("Agent context documentation\n", 100)
	docPath := filepath.Join(docsDir, "AGENT_CONTEXT.md")
	if err := os.WriteFile(docPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write AGENT_CONTEXT.md: %v", err)
	}

	// Change to temp dir
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	// Execute
	ctx := context.Background()
	result, resp, err := HandleGetAgentContext(ctx, nil, models.GetAgentContextRequest{})

	// Verify
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got error result")
	}
	if len(resp.Content) != len(content) {
		t.Errorf("expected content size %d, got %d", len(content), len(resp.Content))
	}
}
