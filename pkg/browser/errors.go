package browser

import (
	"fmt"
	"strings"
)

// FormatPlaywrightError extracts and formats playwright errors without wrapping
func FormatPlaywrightError(err error) string {
	if err == nil {
		return ""
	}

	// Return playwright error message as-is
	return err.Error()
}

// CreateErrorResponse creates a standardized error response for browser tools
func CreateErrorResponse(err error, defaultMsg string) map[string]interface{} {
	msg := FormatPlaywrightError(err)
	if msg == "" {
		msg = defaultMsg
	}

	return map[string]interface{}{
		"error": msg,
		"code":  "BROWSER_ERROR",
	}
}

// IsElementNotFoundError checks if error is element not found
func IsElementNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "no matching element") ||
		strings.Contains(errMsg, "unable to find element")
}

// IsTimeoutError checks if error is a timeout
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "exceeded")
}

// IsConnectionClosedError checks if browser/page connection is closed
func IsConnectionClosedError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "closed") ||
		strings.Contains(errMsg, "context or browser has been closed")
}

// WrapError wraps a playwright error with context
func WrapError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s failed: %w", operation, err)
}
