package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/browser"
	"awesomeProject/pkg/models"
)

// HandleBrowserWaitFor waits for text, element, or timeout
func HandleBrowserWaitFor(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserWaitForRequest) (*mcp.CallToolResult, models.BrowserWaitForResponse, error) {
	// At least one condition must be provided
	if input.Text == "" && input.TextGone == "" && input.Time == 0 {
		return errorResult("INVALID_REQUEST", "at least one of text, textGone, or time must be provided"), models.BrowserWaitForResponse{}, nil
	}

	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserWaitForResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserWaitForResponse{}, nil
	}

	// Wait for text to appear
	if input.Text != "" {
		if err := page.WaitForLoadState(); err != nil {
			msg := browser.FormatPlaywrightError(err)
			log.Printf("[Browser] Wait for load state failed: %v", err)
			return errorResult("BROWSER_ERROR", msg), models.BrowserWaitForResponse{}, nil
		}

		// Check if text is present
		if content, err := page.Content(); err != nil || !contains(content, input.Text) {
			return errorResult("BROWSER_ERROR", fmt.Sprintf("text '%s' not found on page", input.Text)), models.BrowserWaitForResponse{}, nil
		}

		bm.ResetIdleTimer()
		elapsed := time.Since(start).Seconds()

		resp := models.BrowserWaitForResponse{
			Success: true,
			Message: fmt.Sprintf("Found text '%s' (%.2fs)", input.Text, elapsed),
		}
		return nil, resp, nil
	}

	// Wait for text to disappear
	if input.TextGone != "" {
		startWait := time.Now()
		timeout := 30 * time.Second

		for time.Since(startWait) < timeout {
			content, err := page.Content()
			if err == nil && !contains(content, input.TextGone) {
				bm.ResetIdleTimer()
				elapsed := time.Since(start).Seconds()

				resp := models.BrowserWaitForResponse{
					Success: true,
					Message: fmt.Sprintf("Text '%s' disappeared (%.2fs)", input.TextGone, elapsed),
				}
				return nil, resp, nil
			}
			time.Sleep(100 * time.Millisecond)
		}

		return errorResult("BROWSER_ERROR", fmt.Sprintf("timeout waiting for text '%s' to disappear", input.TextGone)), models.BrowserWaitForResponse{}, nil
	}

	// Wait for specified time
	if input.Time > 0 {
		time.Sleep(time.Duration(input.Time*1000) * time.Millisecond)
		bm.ResetIdleTimer()
		elapsed := time.Since(start).Seconds()

		resp := models.BrowserWaitForResponse{
			Success: true,
			Message: fmt.Sprintf("Waited %.2fs", elapsed),
		}
		return nil, resp, nil
	}

	return errorResult("BROWSER_ERROR", "wait condition not met"), models.BrowserWaitForResponse{}, nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
