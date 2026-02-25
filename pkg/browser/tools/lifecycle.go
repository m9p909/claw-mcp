package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/browser"
	"awesomeProject/pkg/models"
)

// HandleBrowserClose closes the browser
func HandleBrowserClose(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserCloseRequest) (*mcp.CallToolResult, models.BrowserCloseResponse, error) {
	bm := browser.GetInstance()
	start := time.Now()

	// Close browser (idempotent)
	if err := bm.CloseBrowser(ctx); err != nil {
		msg := browser.FormatPlaywrightError(err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserCloseResponse{}, nil
	}

	elapsed := time.Since(start).Seconds()

	resp := models.BrowserCloseResponse{
		Success: true,
		Message: fmt.Sprintf("Browser closed (%.2fs)", elapsed),
	}
	return nil, resp, nil
}
