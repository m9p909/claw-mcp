# Implementation Tasks

## Phase 1: Setup & Types (Foundational)

- [x] Task 1.1: Add playwright-go dependency
  - Run `go get github.com/microsoft/playwright-go`
  - Update go.mod and go.sum
  - Verify dependency resolves

- [ ] Task 1.2: Create browser package structure
  - Create `pkg/browser/` directory
  - Create files:
    - `pkg/browser/types.go` (request/response structs)
    - `pkg/browser/browser.go` (BrowserManager singleton)
    - `pkg/browser/errors.go` (error handling)
    - `pkg/browser/tools/` subdirectory for tool implementations

- [ ] Task 1.3: Define all request/response types
  - In `pkg/browser/types.go`, create structs for all 12 tools
  - BrowserNavigateRequest/Response
  - BrowserSnapshotRequest/Response
  - BrowserClickRequest/Response
  - BrowserTypeRequest/Response
  - BrowserFillFormRequest/Response
  - BrowserSelectOptionRequest/Response
  - BrowserPressKeyRequest/Response
  - BrowserWaitForRequest/Response
  - BrowserHandleDialogRequest/Response
  - BrowserNavigateBackRequest/Response
  - BrowserHoverRequest/Response
  - BrowserCloseRequest/Response
  - All with proper `json` and `jsonschema` tags

---

## Phase 2: Core Browser Manager (Lifecycle)

- [ ] Task 2.1: Implement BrowserManager singleton in `pkg/browser/browser.go`
  - Define `type BrowserManager struct` with:
    - `browser *playwright.Browser`
    - `page *playwright.Page`
    - `mu sync.RWMutex`
    - `idleTimer *time.Timer`
    - `lastActivity time.Time`
    - `config BrowserConfig` (timeouts, etc.)
  - `NewBrowserManager() *BrowserManager` singleton creator
  - `GetInstance() *BrowserManager` global access

- [ ] Task 2.2: Implement browser initialization
  - `func (bm *BrowserManager) ensureBrowser(ctx context.Context) error`
  - Launches Chromium via playwright if not running
  - Sets headless=true
  - Stores in bm.browser and bm.page
  - Called by every tool before use

- [ ] Task 2.3: Implement idle timeout mechanism
  - `func (bm *BrowserManager) startIdleTimer()`
  - Runs in background goroutine
  - Closes browser after PLAYWRIGHT_IDLE_TIMEOUT_SECS (env var, default 300s)
  - `func (bm *BrowserManager) resetIdleTimer()`
  - Called by every tool to extend timeout
  - Thread-safe with mutex

- [ ] Task 2.4: Implement resource cleanup
  - `func (bm *BrowserManager) closeBrowser(ctx context.Context) error`
  - Closes page and browser safely
  - Called by idle timeout or browser_close tool
  - Idempotent—safe to call multiple times

---

## Phase 3: Tool Implementations (Navigation)

- [ ] Task 3.1: Implement browser_navigate in `pkg/browser/tools/navigation.go`
  - `HandleBrowserNavigate(ctx, req, input) (*mcp.CallToolResult, BrowserNavigateResponse, error)`
  - Parse URL from input
  - Call `bm.ensureBrowser()` to initialize if needed
  - `page.Goto(url, playwright.PageGotoOptions{Timeout: timeout})`
  - Handle errors, return Playwright's error messages directly
  - Reset idle timer

- [ ] Task 3.2: Implement browser_snapshot in `pkg/browser/tools/snapshot.go`
  - `HandleBrowserSnapshot(ctx, req, input) (*mcp.CallToolResult, BrowserSnapshotResponse, error)`
  - Get accessibility tree from page via playwright
  - Return structured snapshot string with element refs
  - Ensure browser initialized first

- [ ] Task 3.3: Implement browser_navigate_back in `pkg/browser/tools/navigation.go`
  - `HandleBrowserNavigateBack(ctx, req, input) (*mcp.CallToolResult, BrowserNavigateBackResponse, error)`
  - Call `page.GoBack(timeout)`
  - Handle history errors

---

## Phase 4: Tool Implementations (Interaction)

- [ ] Task 4.1: Implement browser_click in `pkg/browser/tools/interaction.go`
  - `HandleBrowserClick(ctx, req, input) (*mcp.CallToolResult, BrowserClickResponse, error)`
  - Get element by ref from snapshot
  - Apply modifiers (Alt, Control, Meta, Shift)
  - Call locator.Click() with button/doubleClick options
  - Handle "element not found" and timeout errors

- [ ] Task 4.2: Implement browser_hover in `pkg/browser/tools/interaction.go`
  - `HandleBrowserHover(ctx, req, input) (*mcp.CallToolResult, BrowserHoverResponse, error)`
  - Get element by ref
  - Call locator.Hover()

---

## Phase 5: Tool Implementations (Input)

- [ ] Task 5.1: Implement browser_type in `pkg/browser/tools/input.go`
  - `HandleBrowserType(ctx, req, input) (*mcp.CallToolResult, BrowserTypeResponse, error)`
  - If ref provided, focus that element first
  - Type text character by character via `locator.Type(text)`
  - Handle focus errors

- [ ] Task 5.2: Implement browser_fill_form in `pkg/browser/tools/input.go`
  - `HandleBrowserFillForm(ctx, req, input) (*mcp.CallToolResult, BrowserFillFormResponse, error)`
  - Iterate fields, fill each by ref
  - `locator.Fill(value)` for each
  - Return success/failure with message

- [ ] Task 5.3: Implement browser_select_option in `pkg/browser/tools/input.go`
  - `HandleBrowserSelectOption(ctx, req, input) (*mcp.CallToolResult, BrowserSelectOptionResponse, error)`
  - Get select element by ref
  - `locator.SelectOption(values)` for each value

- [ ] Task 5.4: Implement browser_press_key in `pkg/browser/tools/input.go`
  - `HandleBrowserPressKey(ctx, req, input) (*mcp.CallToolResult, BrowserPressKeyResponse, error)`
  - If ref, focus first
  - `page.Keyboard.Press(key)`
  - Validate key names (Enter, Tab, Escape, etc.)

---

## Phase 6: Tool Implementations (Async/Dialog)

- [ ] Task 6.1: Implement browser_wait_for in `pkg/browser/tools/async.go`
  - `HandleBrowserWaitFor(ctx, req, input) (*mcp.CallToolResult, BrowserWaitForResponse, error)`
  - One of text/textGone/time must be provided
  - If text: `page.WaitForSelector(text, timeout)`
  - If textGone: `page.WaitForFunction(JS to check element gone)`
  - If time: `time.Sleep(time.Duration(seconds) * time.Second)`

- [ ] Task 6.2: Implement browser_handle_dialog in `pkg/browser/tools/dialogs.go`
  - `HandleBrowserHandleDialog(ctx, req, input) (*mcp.CallToolResult, BrowserHandleDialogResponse, error)`
  - Set dialog handler on page
  - If accept=true and promptText provided: respond with text
  - If accept=true: click OK
  - If accept=false: click Cancel

---

## Phase 7: Tool Implementations (Lifecycle)

- [ ] Task 7.1: Implement browser_close in `pkg/browser/tools/lifecycle.go`
  - `HandleBrowserClose(ctx, req, input) (*mcp.CallToolResult, BrowserCloseResponse, error)`
  - Call `bm.closeBrowser(ctx)`
  - Idempotent—no error if already closed

## Phase 8: Integration & Registration

- [ ] Task 8.1: Register all tools in MCP server in `internal/server.go`
  - In `registerTools()`, add 12 new mcp.AddTool calls for browser tools
  - Use exact names: browser_navigate, browser_snapshot, etc.
  - Wire handlers from pkg/browser/tools to each
  - Update log message: "MCP Server initialized with 20 tools" (was 8)

- [ ] Task 8.2: Add environment variable support in `pkg/browser/config.go`
  - Read PLAYWRIGHT_IDLE_TIMEOUT_SECS (default 300s)
  - Read PLAYWRIGHT_TOOL_TIMEOUT_SECS (default 30s)
  - BrowserConfig struct with these values
  - Used in BrowserManager initialization

## Phase 9: Error Handling & Polish

- [ ] Task 9.1: Implement error handling in `pkg/browser/errors.go`
  - `func formatPlaywrightError(err error) string`
  - Extract playwright library error messages
  - Return as-is without wrapping
  - All tools use this for error responses

- [ ] Task 9.2: Add logging
  - Log browser initialization in browser.go
  - Log idle timeout closure (debug level)
  - Log tool calls with timing

## Phase 10: Testing & Validation

- [ ] Task 10.1: Manual test navigate + snapshot
  - Start server
  - Call browser_navigate to https://example.com
  - Call browser_snapshot
  - Verify snapshot contains element refs
  - Verify browser initialized and stayed open

- [ ] Task 10.2: Manual test click
  - Navigate to simple page with button
  - Call browser_snapshot
  - Extract button ref from snapshot
  - Call browser_click with that ref
  - Verify button was clicked

- [ ] Task 10.3: Manual test idle timeout
  - Set PLAYWRIGHT_IDLE_TIMEOUT_SECS=10
  - Navigate to page
  - Wait 11 seconds without any tool calls
  - Call browser_navigate again
  - Verify browser reinitialized (new instance)

- [ ] Task 10.4: Manual test error handling
  - Call browser_click with invalid ref
  - Verify detailed error returned
  - No custom error messages wrapping

---

## Dependencies

- Tasks 1.x must complete before 2.x
- Task 2.x (browser manager) must complete before 3.x-7.x
- Tasks 3.x-7.x can run in parallel
- Task 8.x (registration) requires 3.x-7.x complete
- Tasks 9.x-10.x last

## Total Effort

~14 Go files created/modified, 12 tool implementations, ~500-700 LOC estimated.
