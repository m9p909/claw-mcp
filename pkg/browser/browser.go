package browser

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

var (
	instance *BrowserManager
	once     sync.Once
)

type BrowserManager struct {
	browser      playwright.Browser
	page         playwright.Page
	mu           sync.RWMutex
	idleTimer    *time.Timer
	lastActivity time.Time
	config       BrowserConfig
}

// NewBrowserManager creates or returns the singleton BrowserManager
func NewBrowserManager() *BrowserManager {
	once.Do(func() {
		instance = &BrowserManager{
			lastActivity: time.Now(),
			config:       loadConfig(),
		}
	})
	return instance
}

// GetInstance returns the singleton BrowserManager
func GetInstance() *BrowserManager {
	return NewBrowserManager()
}

// loadConfig reads environment variables for configuration
func loadConfig() BrowserConfig {
	idleTimeout := 300
	toolTimeout := 30

	if env := os.Getenv("PLAYWRIGHT_IDLE_TIMEOUT_SECS"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 {
			idleTimeout = val
		}
	}

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

// ensureBrowser initializes browser and page if not already running
func (bm *BrowserManager) ensureBrowser(ctx context.Context) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.browser != nil {
		log.Println("[Browser] Browser already running, skipping initialization")
		return nil
	}

	log.Println("[Browser] Initializing Chromium browser in headless mode")

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("failed to run playwright: %w", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to launch chromium: %w", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		browser.Close()
		return fmt.Errorf("failed to create page: %w", err)
	}

	bm.browser = browser
	bm.page = page
	bm.lastActivity = time.Now()
	bm.startIdleTimer()

	log.Println("[Browser] Browser initialized successfully")
	return nil
}

// startIdleTimer starts a background timer to close the browser after idle timeout
func (bm *BrowserManager) startIdleTimer() {
	if bm.idleTimer != nil {
		bm.idleTimer.Stop()
	}

	bm.idleTimer = time.AfterFunc(time.Duration(bm.config.IdleTimeoutSecs)*time.Second, func() {
		log.Printf("[Browser] Idle timeout reached (%d sec), closing browser", bm.config.IdleTimeoutSecs)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		bm.closeBrowser(ctx)
	})
}

// resetIdleTimer extends the browser's lifetime
func (bm *BrowserManager) resetIdleTimer() {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.browser == nil {
		return
	}

	bm.lastActivity = time.Now()
	bm.startIdleTimer()
}

// closeBrowser safely closes the browser and page
func (bm *BrowserManager) closeBrowser(ctx context.Context) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.idleTimer != nil {
		bm.idleTimer.Stop()
		bm.idleTimer = nil
	}

	if bm.page != nil {
		if err := bm.page.Close(); err != nil {
			log.Printf("[Browser] Error closing page: %v", err)
		}
		bm.page = nil
	}

	if bm.browser != nil {
		if err := bm.browser.Close(); err != nil {
			log.Printf("[Browser] Error closing browser: %v", err)
		}
		bm.browser = nil
	}

	log.Println("[Browser] Browser closed and resources freed")
	return nil
}

// GetPage returns the current page (for tool implementations)
func (bm *BrowserManager) GetPage() playwright.Page {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.page
}

// GetBrowser returns the current browser (for tool implementations)
func (bm *BrowserManager) GetBrowser() playwright.Browser {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.browser
}

// IsBrowserRunning checks if browser is currently running
func (bm *BrowserManager) IsBrowserRunning() bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.browser != nil
}

// EnsureBrowser is public wrapper for ensureBrowser
func (bm *BrowserManager) EnsureBrowser(ctx context.Context) error {
	return bm.ensureBrowser(ctx)
}

// ResetIdleTimer is public wrapper
func (bm *BrowserManager) ResetIdleTimer() {
	bm.resetIdleTimer()
}

// CloseBrowser is public wrapper
func (bm *BrowserManager) CloseBrowser(ctx context.Context) error {
	return bm.closeBrowser(ctx)
}
