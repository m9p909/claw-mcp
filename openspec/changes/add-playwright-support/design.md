# Playwright Browser Integration Design

## Overview

Add native Playwright browser automation to Claw MCP server using `playwright-go` library. This provides LLM-driven web browser automation through structured accessibility snapshots, with exact tool signatures matching Microsoft's `playwright-mcp`.

## Architecture

### Package Structure

```
pkg/browser/
├── browser.go          # BrowserManager singleton, session lifecycle
├── context.go          # BrowserContext wrapper around playwright.Browser
├── types.go            # Request/Response models for all tools
├── store.go            # Thread-safe session storage
└── tools/
    ├── navigation.go   # browser_navigate, browser_navigate_back, browser_wait_for
    ├── interaction.go  # browser_click, browser_hover, browser_drag
    ├── input.go        # browser_type, browser_fill_form, browser_select_option, browser_press_key
    ├── snapshot.go     # browser_snapshot
    ├── dialogs.go      # browser_handle_dialog
    └── scripting.go    # browser_evaluate, browser_run_code
```

### Session Management

- **Single shared browser instance** across all MCP clients (singleton like Claw itself)
- Browser initialized on first tool use, maintained in memory
- **Idle timeout mechanism**: Browser closes after no tool calls for configurable duration (default 5 min)
- Thread-safe operations via mutexes on BrowserManager
- Each tool call resets the idle timeout

### Core Tools (12 essential)

1. **browser_navigate** - Open URL
2. **browser_snapshot** - Get accessibility tree (structured, LLM-friendly)
3. **browser_click** - Click element by accessibility ref
4. **browser_type** - Type text into focused field
5. **browser_fill_form** - Fill multiple form fields
6. **browser_select_option** - Select dropdown option
7. **browser_press_key** - Press keyboard keys
8. **browser_wait_for** - Wait for text/element/timeout
9. **browser_handle_dialog** - Accept/dismiss alerts/prompts
10. **browser_navigate_back** - Go back in history
11. **browser_hover** - Hover over element
12. **browser_close** - Close browser, release resources

### Tool Signatures (Matching playwright-mcp)

All tools follow the exact structure from Microsoft's implementation.

Example: browser_click
```go
type BrowserClickRequest struct {
    Ref      string   `json:"ref" jsonschema:"description,Exact target element reference from snapshot"`
    Element  string   `json:"element" jsonschema:"description,Human-readable element description"`
    Button   string   `json:"button" jsonschema:"description,left|right|middle (default: left)"`
    DoubleClick bool  `json:"doubleClick" jsonschema:"description,Whether to double-click"`
    Modifiers []string `json:"modifiers" jsonschema:"description,Modifier keys: Alt, Control, Meta, Shift"`
}

type BrowserClickResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

### Accessibility Snapshots

`browser_snapshot` returns a structured representation of the page accessible to LLMs. Uses Playwright's built-in accessibility tree, not screenshots.

```json
{
  "snapshot": "<accessibility tree as string>"
}
```

The snapshot includes element references (refs) that are used in subsequent tool calls.

## Error Handling

**Philosophy**: Expose playwright errors directly without wrapping/simplifying.

- Playwright library exceptions caught and returned as detailed error strings
- Includes playwright's own error messages (element not found, timeout, etc.)
- Format: `{"error": "<playwright error message>", "code": "BROWSER_ERROR"}`
- NO custom error messages that hide the true cause

## Browser Lifecycle

### Initialization
1. First tool call triggers browser.Launch()
2. Stored in BrowserManager singleton
3. Context initialized with timeout

### Active Use
- Each tool call resets idle timeout to 0
- Timer runs in background goroutine

### Idle Timeout
- After N seconds with no tool calls, browser.Close()
- Resources freed, next tool call reinitializes
- Configurable via environment variable (PLAYWRIGHT_IDLE_TIMEOUT_SECS, default 300)

### Explicit Close
- `browser_close` tool immediately closes browser
- Idempotent—safe to call multiple times

## Implementation Notes

### Dependencies
- `github.com/microsoft/playwright-go` - Playwright Go library
- No Node.js dependency (unlike playwright-mcp)
- Runs entirely in same Go process

### Thread Safety
- BrowserManager uses RWMutex for session access
- Playwright operations are already thread-safe within a context
- Timeout ticker managed safely

### Timeout Handling
- Tool calls timeout after 30 seconds (configurable)
- Browser idle timeout separate from tool call timeout

## Integration with Existing Claw Tools

New tools added to `internal/server.go` in registerTools():

```go
// Browser automation tools
mcp.AddTool(s.mcpServer,
    &mcp.Tool{Name: "browser_navigate", Description: "Navigate to URL"},
    tools.HandleBrowserNavigate)
mcp.AddTool(s.mcpServer,
    &mcp.Tool{Name: "browser_snapshot", Description: "Get accessibility snapshot"},
    tools.HandleBrowserSnapshot)
// ... etc for other tools
```

Total tools: 8 (existing) + 12 (browser) = 20 tools

## No Screenshots

Explicitly out of scope for now:
- Visual screenshots not implemented
- Uses accessibility snapshots instead (LLM-native, efficient, lower bandwidth)
- Can add later if needed

## Testing Strategy

Manual testing via MCP client:
1. Start server with PLAYWRIGHT_IDLE_TIMEOUT_SECS=10 (short timeout for testing)
2. Call browser_navigate to https://example.com
3. Call browser_snapshot, verify structure
4. Call browser_click with ref from snapshot
5. Verify idle timeout closes browser after inactivity

## Future Considerations

- Per-session browser contexts (if multi-tenant needed later)
- Screenshot support with custom sizing
- Video recording
- Network interception
- Cookie management
- Local storage persistence
