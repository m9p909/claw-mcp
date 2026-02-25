package browser

import (
	"os"
	"strconv"
)

// LoadConfig reads environment variables for Playwright configuration
func LoadConfig() BrowserConfig {
	idleTimeout := 300 // 5 minutes
	toolTimeout := 30  // 30 seconds

	// PLAYWRIGHT_IDLE_TIMEOUT_SECS: how long browser stays alive after last tool call
	if env := os.Getenv("PLAYWRIGHT_IDLE_TIMEOUT_SECS"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 {
			idleTimeout = val
		}
	}

	// PLAYWRIGHT_TOOL_TIMEOUT_SECS: maximum duration for any single tool call
	if env := os.Getenv("PLAYWRIGHT_TOOL_TIMEOUT_SECS"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 {
			toolTimeout = val
		}
	}

	return BrowserConfig{
		IdleTimeoutSecs: idleTimeout,
		ToolTimeoutSecs: toolTimeout,
	}
}
