# Playwright Browser Tools Specification

## Tool: browser_navigate

Navigate to a URL and wait for page load.

### Request

```json
{
  "url": "string (required)",
  "timeout": "number (optional, seconds, default: 30)"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - Playwright error (page load failed, timeout, invalid URL, etc.)

---

## Tool: browser_snapshot

Get structured accessibility tree snapshot of current page. Returns element references (refs) that can be used in subsequent tool calls.

### Request

```json
{}
```

No parameters.

### Response

```json
{
  "snapshot": "string (accessibility tree representation)"
}
```

### Errors

- `BROWSER_ERROR` - No active browser context

---

## Tool: browser_click

Click element on page using accessibility reference.

### Request

```json
{
  "ref": "string (required, element reference from snapshot)",
  "element": "string (optional, human description for logging)",
  "button": "string (optional, 'left'|'right'|'middle', default: 'left')",
  "doubleClick": "boolean (optional, default: false)",
  "modifiers": "array of strings (optional, 'Alt'|'Control'|'Meta'|'Shift')"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - Element not found, click failed, etc.

---

## Tool: browser_type

Type text into currently focused element.

### Request

```json
{
  "text": "string (required, text to type)",
  "ref": "string (optional, element reference to focus first)"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - Element not found, focus failed, etc.

---

## Tool: browser_fill_form

Fill multiple form fields at once.

### Request

```json
{
  "fields": [
    {
      "ref": "string (element reference)",
      "value": "string (value to fill)",
      "name": "string (optional, human-readable field name)"
    }
  ]
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - One or more fields not found or not fillable

---

## Tool: browser_select_option

Select option from dropdown/select element.

### Request

```json
{
  "ref": "string (required, select element reference)",
  "values": "array of strings (required, option values to select)",
  "element": "string (optional, human description)"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - Element not found, option not available, etc.

---

## Tool: browser_press_key

Press keyboard key(s).

### Request

```json
{
  "key": "string (required, key name: 'Enter', 'Tab', 'Escape', etc., or single character)",
  "ref": "string (optional, element to focus first)"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - Invalid key name, focus failed, etc.

---

## Tool: browser_wait_for

Wait for element, text, or time.

### Request

```json
{
  "text": "string (optional, text to wait for)",
  "textGone": "string (optional, text to wait to disappear)",
  "time": "number (optional, seconds to wait)"
}
```

At least one of `text`, `textGone`, or `time` must be provided.

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - Timeout, no element found, etc.

---

## Tool: browser_handle_dialog

Handle JavaScript dialogs (alert, confirm, prompt).

### Request

```json
{
  "accept": "boolean (required, true to accept/OK, false to cancel)",
  "promptText": "string (optional, text to enter in prompt dialogs)"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - No dialog present, prompt failed, etc.

---

## Tool: browser_navigate_back

Navigate backward in browser history.

### Request

```json
{
  "timeout": "number (optional, seconds, default: 30)"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - No history to go back, timeout, etc.

---

## Tool: browser_hover

Hover over element.

### Request

```json
{
  "ref": "string (required, element reference)",
  "element": "string (optional, human description)"
}
```

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Errors

- `BROWSER_ERROR` - Element not found, hover failed, etc.

---

## Tool: browser_close

Close browser and release all resources.

### Request

```json
{}
```

No parameters.

### Response

```json
{
  "success": "boolean",
  "message": "string"
}
```

### Notes

- Idempotent—safe to call multiple times
- After close, first tool call will reopen browser
- Bypass idle timeout cleanup

---

## Error Code: BROWSER_ERROR

All browser tool failures return this code with Playwright's original error message. Examples:

- "Timeout 30000ms exceeded waiting for locator to be visible"
- "Element not found"
- "Target page, context or browser has been closed"
- "Invalid URL: not a valid URL"

Do NOT wrap or simplify these messages—expose them directly for debugging.

---

## Browser Lifecycle

### Initialization Trigger
First tool call initializes browser with playwright-go defaults:
- Chromium browser (hardcoded)
- Headless mode (always, no GUI)
- Default viewport size

### Idle Timeout
Background ticker closes browser after PLAYWRIGHT_IDLE_TIMEOUT_SECS (env var, default 300 sec = 5 min):
- Resets on every tool call
- Freed resources available for reuse
- Next tool call reinitializes

### Tool Timeout
Each tool call has maximum duration of PLAYWRIGHT_TOOL_TIMEOUT_SECS (env var, default 30 sec):
- If tool exceeds timeout, operation cancels
- Returns `BROWSER_ERROR` with timeout message

---

## Implementation Constraints

1. **Single Browser Instance**: All clients share one browser (singleton pattern like Claw)
2. **Headless Only**: GUI mode explicitly disabled, no environment variable to change this
3. **No Screenshots**: Accessibility snapshots only, visual screenshots out of scope
4. **Chromium Only**: No Firefox/WebKit selection (can add later if needed)
5. **Shared State**: Browser navigation history, cookies, etc. are shared across all clients
6. **No Per-Request Auth**: Browser doesn't authenticate per-request (shares cookies)

---

## Environment Variables

```bash
PLAYWRIGHT_IDLE_TIMEOUT_SECS      # Browser idle close timeout (default: 300)
PLAYWRIGHT_TOOL_TIMEOUT_SECS      # Per-tool timeout (default: 30)
```
