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

// HandleBrowserType types text into an element
func HandleBrowserType(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserTypeRequest) (*mcp.CallToolResult, models.BrowserTypeResponse, error) {
	if input.Text == "" {
		return errorResult("INVALID_REQUEST", "text is required"), models.BrowserTypeResponse{}, nil
	}

	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserTypeResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserTypeResponse{}, nil
	}

	// If ref provided, focus first
	if input.Ref != "" {
		if err := page.Locator(input.Ref).Focus(); err != nil {
			msg := browser.FormatPlaywrightError(err)
			log.Printf("[Browser] Focus on %s failed: %v", input.Ref, err)
			return errorResult("BROWSER_ERROR", msg), models.BrowserTypeResponse{}, nil
		}
	}

	// Type text
	keyboard := page.Keyboard()
	if err := keyboard.Type(input.Text); err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Type failed: %v", err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserTypeResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserTypeResponse{
		Success: true,
		Message: fmt.Sprintf("Typed %d characters (%.2fs)", len(input.Text), elapsed),
	}
	return nil, resp, nil
}

// HandleBrowserFillForm fills multiple form fields
func HandleBrowserFillForm(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserFillFormRequest) (*mcp.CallToolResult, models.BrowserFillFormResponse, error) {
	if len(input.Fields) == 0 {
		return errorResult("INVALID_REQUEST", "fields array is required"), models.BrowserFillFormResponse{}, nil
	}

	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserFillFormResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserFillFormResponse{}, nil
	}

	filledCount := 0
	for _, field := range input.Fields {
		if field.Ref == "" {
			continue
		}

		if err := page.Locator(field.Ref).Fill(field.Value); err != nil {
			msg := browser.FormatPlaywrightError(err)
			log.Printf("[Browser] Fill field %s failed: %v", field.Ref, err)
			return errorResult("BROWSER_ERROR", msg), models.BrowserFillFormResponse{}, nil
		}
		filledCount++
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserFillFormResponse{
		Success: true,
		Message: fmt.Sprintf("Filled %d fields (%.2fs)", filledCount, elapsed),
	}
	return nil, resp, nil
}

// HandleBrowserSelectOption selects option from dropdown
func HandleBrowserSelectOption(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserSelectOptionRequest) (*mcp.CallToolResult, models.BrowserSelectOptionResponse, error) {
	if input.Ref == "" {
		return errorResult("INVALID_REQUEST", "ref is required"), models.BrowserSelectOptionResponse{}, nil
	}

	if len(input.Values) == 0 {
		return errorResult("INVALID_REQUEST", "values array is required"), models.BrowserSelectOptionResponse{}, nil
	}

	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserSelectOptionResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserSelectOptionResponse{}, nil
	}

	// Select options
	selectVals := playwright.SelectOptionValues{
		Values: &input.Values,
	}

	if _, err := page.Locator(input.Ref).SelectOption(selectVals); err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Select option on %s failed: %v", input.Ref, err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserSelectOptionResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserSelectOptionResponse{
		Success: true,
		Message: fmt.Sprintf("Selected %d options (%.2fs)", len(input.Values), elapsed),
	}
	return nil, resp, nil
}

// HandleBrowserPressKey presses a keyboard key
func HandleBrowserPressKey(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserPressKeyRequest) (*mcp.CallToolResult, models.BrowserPressKeyResponse, error) {
	if input.Key == "" {
		return errorResult("INVALID_REQUEST", "key is required"), models.BrowserPressKeyResponse{}, nil
	}

	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserPressKeyResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserPressKeyResponse{}, nil
	}

	// If ref provided, focus first
	if input.Ref != "" {
		if err := page.Locator(input.Ref).Focus(); err != nil {
			msg := browser.FormatPlaywrightError(err)
			log.Printf("[Browser] Focus on %s failed: %v", input.Ref, err)
			return errorResult("BROWSER_ERROR", msg), models.BrowserPressKeyResponse{}, nil
		}
	}

	// Press key
	keyboard := page.Keyboard()
	if err := keyboard.Press(input.Key); err != nil {
		msg := browser.FormatPlaywrightError(err)
		log.Printf("[Browser] Press key %s failed: %v", input.Key, err)
		return errorResult("BROWSER_ERROR", msg), models.BrowserPressKeyResponse{}, nil
	}

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	resp := models.BrowserPressKeyResponse{
		Success: true,
		Message: fmt.Sprintf("Pressed key '%s' (%.2fs)", input.Key, elapsed),
	}
	return nil, resp, nil
}
