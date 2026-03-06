package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	pkglog "awesomeProject/pkg/log"
	"awesomeProject/pkg/models"
)

// HandleGetAgentContext retrieves and returns the agent context documentation.
func HandleGetAgentContext(ctx context.Context, req *mcp.CallToolRequest, input models.GetAgentContextRequest) (*mcp.CallToolResult, models.GetAgentContextResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	logger.Info(ctx, "Getting agent context")

	// Construct path to docs/AGENT_CONTEXT.md
	cwd, err := os.Getwd()
	if err != nil {
		return errorResult(ctx, "AGENT_CONTEXT_ERROR", fmt.Sprintf("failed to get current directory: %v", err)), models.GetAgentContextResponse{}, nil
	}

	docPath := filepath.Join(cwd, "docs", "AGENT_CONTEXT.md")

	// Read documentation file
	content, err := os.ReadFile(docPath)
	if err != nil {
		msg := fmt.Sprintf("failed to read documentation file at %s: %v", docPath, err)
		return errorResult(ctx, "DOC_NOT_FOUND", msg), models.GetAgentContextResponse{}, nil
	}

	resp := models.GetAgentContextResponse{
		Content: string(content),
	}

	logger.Info(ctx, "Retrieved agent context",
		"doc_path", docPath,
		"content_size", len(content),
		pkglog.Duration(time.Since(start)))

	return nil, resp, nil
}
