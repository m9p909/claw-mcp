package tools

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/playwright-community/playwright-go"

	"awesomeProject/pkg/browser"
	"awesomeProject/pkg/models"
)

// HandleBrowserClick clicks an element
func HandleBrowserClick(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserClickRequest) (*mcp.CallToolResult, models.BrowserClickResponse, error) {
	if input.Ref == "" {
		return errorResult("INVALID_REQUEST", "ref is required"), models.BrowserClickResponse{}, nil
	}

	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserClickResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserClickResponse{}, nil
	}

	// Build click options
	opts := playwright.LocatorClickOptions{}

	if input.Button != "" && input.Button != "left" {
		btn := playwright.MouseButton(input.Button)
		opts.Button = &btn
	}

	if input.DoubleClick {
		opts.ClickCount = playwright.Int(2)
	}

	if len(input.Modifiers) > 0 {
		modifiers := make([]playwright.KeyboardModifier, len(input.Modifiers))
		for i, mod := range input.Modifiers {
			modifiers[i] = playwright.KeyboardModifier(mod)
		}
		opts.Modifiers = modifiers
	}

	// Click the element using ref as selector
	if err := page.Locator(input.Ref).Click(opts); err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Click on %s failed: %v", input.Ref, err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserClickResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserClickResponse{
		Success: true,
		Message: fmt.Sprintf("Clicked element (%.2fs)", elapsed),
	}
	return nil, resp, nil
}

// HandleBrowserHover hovers over an element
func HandleBrowserHover(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserHoverRequest) (*mcp.CallToolResult, models.BrowserHoverResponse, error) {
	if input.Ref == "" {
		return errorResult("INVALID_REQUEST", "ref is required"), models.BrowserHoverResponse{}, nil
	}

	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserHoverResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserHoverResponse{}, nil
	}

	// Hover over element using ref as selector
	if err := page.Locator(input.Ref).Hover(); err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Hover on %s failed: %v", input.Ref, err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserHoverResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserHoverResponse{
		Success: true,
		Message: fmt.Sprintf("Hovered over element (%.2fs)", elapsed),
	}
	return nil, resp, nil
}

// Helper to convert modifier names to playwright format
func normalizeModifiers(input []string) []string {
	var result []string
	for _, mod := range input {
		// Playwright expects: Alt, Control, Meta, Shift
		normalized := strings.Title(strings.ToLower(mod))
		if normalized == "Cmd" {
			normalized = "Meta"
		}
		result = append(result, normalized)
	}
	return result
}
