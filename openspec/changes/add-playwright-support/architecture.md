# Architecture Diagram

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    MCP Client (Claude, etc.)                    │
└───────────────────┬─────────────────────────────────────────────┘
                    │ HTTP POST /mcp
                    │ (Bearer token auth)
                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Claw MCP Server (Go)                         │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  main.go - HTTP Handler, Auth Middleware                │  │
│  └──────────────────────────────────────────────────────────┘  │
│                          ▼                                      │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  internal/server.go - MCP Tool Registration             │  │
│  │                                                          │  │
│  │  - Filesystem tools (8)                                │  │
│  │  - Browser tools (12) ◄──── NEW                        │  │
│  └──────────────────────────────────────────────────────────┘  │
│                          ▼                                      │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  pkg/browser/ - Browser Automation (NEW)               │  │
│  │                                                          │  │
│  │  ┌─ browser.go                                          │  │
│  │  │  └─ BrowserManager (singleton)                      │  │
│  │  │     ├─ Browser lifecycle                            │  │
│  │  │     ├─ Idle timeout management                      │  │
│  │  │     └─ Page/Context handling                        │  │
│  │  │                                                       │  │
│  │  ├─ types.go                                            │  │
│  │  │  └─ Request/Response structs (12 tools)             │  │
│  │  │                                                       │  │
│  │  └─ tools/                                              │  │
│  │     ├─ navigation.go (navigate, navigate_back)         │  │
│  │     ├─ snapshot.go (snapshot)                          │  │
│  │     ├─ interaction.go (click, hover, drag)             │  │
│  │     ├─ input.go (type, fill_form, select, press_key) │  │
│  │     ├─ async.go (wait_for)                            │  │
│  │     ├─ dialogs.go (handle_dialog)                      │  │
│  │     └─ lifecycle.go (close)                            │  │
│  │                                                          │  │
│  └──────────────────────────────────────────────────────────┘  │
│                          ▲                                      │
│                          │ Uses                                 │
└──────────────────────────┼──────────────────────────────────────┘
                           │
                           ▼
            ┌──────────────────────────────┐
            │  playwright-go Library       │
            │  (Chromium automation)       │
            └──────────────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │  Chromium Browser    │
                    │  (Headless Mode)     │
                    └──────────────────────┘
```

---

## Tool Call Flow

```
MCP Request: browser_click
    ▼
internal/server.go (route to handler)
    ▼
pkg/browser/tools/interaction.go::HandleBrowserClick()
    ▼
    ├─ Input validation (ref, button, etc.)
    │
    ├─ BrowserManager.ensureBrowser()
    │  └─ Launch Chromium if not running
    │
    ├─ resetIdleTimer()
    │  └─ Extend browser lifetime
    │
    ├─ Get locator by ref from page
    │
    ├─ Call playwright.Click() with options
    │
    └─ Handle errors
       └─ Return playwright error directly
    ▼
MCP Response: { success: bool, message: string }
```

---

## Browser Lifecycle Timeline

```
Time ──────────────────────────────────────────────────────────────>

Tool Call #1: browser_navigate
  ▼
  [Browser Launched - Idle timer starts]
  │
  ├─ Chromium starts
  ├─ Page loaded
  └─ Timer: 5 min (default PLAYWRIGHT_IDLE_TIMEOUT_SECS=300)

Tool Call #2: browser_snapshot (after 30 sec)
  ▼
  [Tool resets idle timer]
  │
  └─ Timer reset: 5 min from now

  ... 60 seconds of inactivity (no tool calls) ...

Tool Call #3: browser_click (after 90 sec total)
  ▼
  [Still running, timer reset again]
  │
  └─ Timer reset: 5 min from now

  ... 5 minutes, 10 seconds of inactivity ...

[Idle timeout fires]
  ▼
  [Browser.Close() called]
  │
  ├─ Page closed
  ├─ Browser closed
  └─ Resources freed

Tool Call #4: browser_navigate (after idle closure)
  ▼
  [Browser reinitialized]
  │
  └─ New Chromium instance launched
```

---

## Thread Safety Model

```
┌─────────────────────────────────────┐
│  BrowserManager (Singleton)         │
│  ├─ mu sync.RWMutex                 │
│  ├─ browser *playwright.Browser     │
│  ├─ page *playwright.Page           │
│  ├─ idleTimer *time.Timer           │
│  └─ lastActivity time.Time          │
└─────────────────────────────────────┘
          ▲
          │
          ├─ ReadLock: All tool handlers
          │  └─ Concurrent reads allowed
          │
          └─ WriteLock: ensureBrowser(), closeBrowser()
             └─ Exclusive access during init/cleanup
```

Every tool:
1. Acquires RWLock (read)
2. Checks if browser initialized
3. If not, upgrades to write lock (atomic)
4. Performs operation
5. Resets idle timer (write lock)
6. Releases lock

---

## Error Handling Strategy

```
Playwright Error (native)
    ▼
    ├─ Caught by tool handler
    │
    └─ Formatted as-is via formatPlaywrightError()
       └─ No wrapping, no simplification
    ▼
{ error: "<original playwright message>", code: "BROWSER_ERROR" }
    ▼
Client receives full context for debugging
```

Examples:
- "Timeout 30000ms exceeded waiting for locator to be visible"
- "Target page, context or browser has been closed"
- "Element not found" → Click on page.locator("ref") failed

---

## Accessibility Snapshot Format

```
Input: browser_snapshot
    ▼
PlaywrightBot.Accessibility.Snapshot() (or similar)
    ▼
Returns Structured HTML + Element References
    ▼
Example output:
{
  "snapshot": "<!DOCTYPE html>\n<html>...",
  "elements": [
    {
      "ref": "elem_1",
      "text": "Submit",
      "selector": "#submit-button",
      "role": "button"
    },
    ...
  ]
}
    ▼
Client parses snapshot, extracts element refs
    ▼
Uses refs in subsequent click, type, etc. calls
```

Refs are opaque identifiers—internal to playwright-go, stable within a page session.

---

## Configuration (Environment Variables)

```
PLAYWRIGHT_IDLE_TIMEOUT_SECS
  │
  └─ How long browser stays alive after last tool call
  └─ Default: 300 seconds (5 minutes)
  └─ Set to 10 for testing, 3600 for long-running workflows

PLAYWRIGHT_TOOL_TIMEOUT_SECS
  │
  └─ Maximum duration for any single tool call
  └─ Default: 30 seconds
  └─ Includes page load time, click time, etc.

CLAW_TOKEN
  │
  └─ Existing auth token (unchanged)
```
