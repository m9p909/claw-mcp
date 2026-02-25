# Add Playwright Browser Support to Claw MCP

## Problem

Claw currently provides file operations, command execution, and memory tools. There's no browser automation capability. Users cannot interact with web applications through the MCP interface.

## Solution

Add native Playwright browser automation using `playwright-go` library. Provide 12 core browser automation tools with exact signatures matching Microsoft's `playwright-mcp`.

## Scope

### Included

- **12 browser tools**: navigate, snapshot, click, type, fill_form, select_option, press_key, wait_for, handle_dialog, navigate_back, hover, close
- **Accessibility snapshots**: Structured page representation (LLM-native, no vision model needed)
- **Single shared browser instance**: Singleton pattern, shared across all MCP clients
- **Idle timeout lifecycle**: Browser closes after N seconds without tool calls, reinitializes on next use
- **Detailed error messages**: Expose playwright errors directly without simplification

### Excluded

- Screenshots (visual automation)
- Per-client browser contexts (complexity not needed yet)
- Firefox/WebKit browsers (Chromium only)
- Cookie/local storage management APIs

## Architecture

New `pkg/browser/` package with:
- BrowserManager singleton (lifecycle, timeouts)
- Tool implementations (navigation, interaction, input, etc.)
- Type definitions matching playwright-mcp

Integrated with existing Claw tool registration in `internal/server.go`.

Total tools: 8 (existing) + 12 (browser) = 20 tools

## Benefits

1. **No subprocess overhead**: Runs entirely in Go process (unlike playwright-mcp which needs Node.js)
2. **Exact signatures**: Identical tool APIs to Microsoft's implementation—clients can switch seamlessly
3. **LLM-friendly**: Accessibility snapshots instead of screenshots = lower bandwidth, deterministic tool calls
4. **Singleton pattern**: Fits Claw's architecture perfectly
5. **Minimal scope**: Only core tools, no bloat

## Risks

1. **Shared browser state**: All clients share cookies, history, etc. (acceptable by design)
2. **Single browser per process**: Can't isolate high-concurrency workloads (acceptable for now)
3. **Resource cleanup**: Idle timeout may prematurely close browser (configurable, 5 min default)

## Timeline

- Small scope change, ~10 core Go files
- Est. implementation manageable once design is locked

## Questions

None—design covers all critical decisions.
