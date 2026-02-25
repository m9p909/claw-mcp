package tools

import (
	"context"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/browser"
	"awesomeProject/pkg/models"
)

// HandleBrowserSnapshot returns accessibility tree snapshot
func HandleBrowserSnapshot(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserSnapshotRequest) (*mcp.CallToolResult, models.BrowserSnapshotResponse, error) {
	bm := browser.GetInstance()
	start := time.Now()

	// Browser must be running
	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserSnapshotResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserSnapshotResponse{}, nil
	}

	// Get page HTML as snapshot (accessibility tree representation)
	content, err := page.Content()
	if err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Failed to get page content: %v", err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserSnapshotResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	log.Printf("[Browser] Snapshot taken (%.2fs)", elapsed)

	resp := models.BrowserSnapshotResponse{
		Snapshot: content,
	}
	return nil, resp, nil
}

// Note: The snapshot returns the page HTML which clients can parse for element references.
// This provides a structured representation of the page state for LLM-driven automation.
