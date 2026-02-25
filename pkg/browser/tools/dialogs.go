package tools

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/playwright-community/playwright-go"

	"awesomeProject/pkg/browser"
	"awesomeProject/pkg/models"
)

var (
	dialogMu      sync.Mutex
	dialogHandled bool
	lastDialog    playwright.Dialog
)

// HandleBrowserHandleDialog handles JavaScript dialogs (alert, confirm, prompt)
func HandleBrowserHandleDialog(ctx context.Context, req *mcp.CallToolRequest, input models.BrowserHandleDialogRequest) (*mcp.CallToolResult, models.BrowserHandleDialogResponse, error) {
	bm := browser.GetInstance()
	start := time.Now()

	if !bm.IsBrowserRunning() {
		return errorResult("BROWSER_ERROR", "browser not running"), models.BrowserHandleDialogResponse{}, nil
	}

	page := bm.GetPage()
	if page == nil {
		return errorResult("BROWSER_ERROR", "page not available"), models.BrowserHandleDialogResponse{}, nil
	}

	// Register dialog handler
	handled := false
	page.Once("dialog", func(dialog playwright.Dialog) {
		dialogMu.Lock()
		defer dialogMu.Unlock()
		lastDialog = dialog

		if input.Accept {
			if input.PromptText != "" {
				dialog.Accept(input.PromptText)
			} else {
				dialog.Accept("")
			}
		} else {
			dialog.Dismiss()
		}
		handled = true
	})

	// Wait briefly for dialog to appear
	time.Sleep(100 * time.Millisecond)

	bm.ResetIdleTimer()
	elapsed := time.Since(start).Seconds()

	if !handled {
		return errorResult("BROWSER_ERROR", "no dialog appeared"), models.BrowserHandleDialogResponse{}, nil
	}

	action := "accepted"
	if !input.Accept {
		action = "dismissed"
	}

	resp := models.BrowserHandleDialogResponse{
		Success: true,
		Message: fmt.Sprintf("Dialog %s (%.2fs)", action, elapsed),
	}
	return nil, resp, nil
}

// SetupDialogHandler sets up a persistent dialog handler on the page
func SetupDialogHandler(page playwright.Page, accept bool, promptText string) {
	page.On("dialog", func(dialog playwright.Dialog) {
		log.Printf("[Browser] Dialog appeared: %s", dialog.Type())

		if accept {
			if promptText != "" {
				dialog.Accept(promptText)
			} else {
				dialog.Accept("")
			}
		} else {
			dialog.Dismiss()
		}
	})
}
