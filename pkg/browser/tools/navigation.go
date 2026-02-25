package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/playwright-community/playwright-go"

	"awesomeProject/pkg/browser"
	"awesomeProject/pkg/models"
)

// HandleBrowserNavigate navigates to a URL
func HandleBrowserNavigate(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserNavigateRequest) (*mcp.CallToolResult, models.BrowserNavigateResponse, error) {
	if input.URL == "" {
		return errorResult("INVALID_REQUEST", "url is required"), models.BrowserNavigateResponse{}, nil
	}

	timeout := 30
	if input.Timeout > 0 {
		timeout = input.Timeout
	}

	bm := browser.GetInstance()
	start := time.Now()

	// Ensure browser is running
	if err := bm.EnsureBrowser(ctx); err != nil {
		msg := browser.FormatPlaywrightError(err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserNavigateResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "browser not initialized"), models.BrowserNavigateResponse{}, nil
	}

	// Navigate to URL
	if _, err := page.Goto(input.URL, playwright.PageGotoOptions{
		Timeout: playwright.Float(float64(timeout * 1000)),
	}); err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Navigate to %s failed: %v", input.URL, err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserNavigateResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserNavigateResponse{
		Success: true,
		Message: fmt.Sprintf("Navigated to %s (%.2fs)", input.URL, elapsed),
	}
	return nil, resp, nil
}

// HandleBrowserNavigateBack navigates backward in history
func HandleBrowserNavigateBack(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserNavigateBackRequest) (*mcp.CallToolResult, models.BrowserNavigateBackResponse, error) {
	timeout := 30
	if input.Timeout > 0 {
		timeout = input.Timeout
	}

	bm := browser.GetInstance()
	start := time.Now()

	// Browser must be running
	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserNavigateBackResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserNavigateBackResponse{}, nil
	}

	// Go back
	if _, err := page.GoBack(playwright.PageGoBackOptions{
		Timeout: playwright.Float(float64(timeout * 1000)),
	}); err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Navigate back failed: %v", err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserNavigateBackResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserNavigateBackResponse{
		Success: true,
		Message: fmt.Sprintf("Navigated back (%.2fs)", elapsed),
	}
	return nil, resp, nil
}
