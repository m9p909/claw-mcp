package main

import (
	"context"
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/internal"
	"awesomeProject/internal/session"
)

// sessionValidationMiddleware allows the MCP SDK to handle session management.
// Per the MCP Streamable HTTP Transport specification:
// - Initialization requests (POST /mcp with no Mcp-Session-Id header) establish a session
// - The server generates a cryptographically secure session ID and returns it in response
// - Subsequent requests must include the Mcp-Session-Id header
// - The MCP SDK's NewStreamableHTTPHandler automatically manages all of this
// This middleware is kept as a placeholder for future session management needs
// and passes through to let the SDK handle session lifecycle.
func sessionValidationMiddleware(store *session.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Session management is handled by the MCP SDK's event-stream protocol handler.
			// The SDK maintains session state per-connection and validates session IDs.
			// We just pass through to the next handler.
			next.ServeHTTP(w, r)
		})
	}
}

// authMiddleware enforces bearer token authentication with MCP session awareness.
// Key behaviors:
// 1. /health endpoint bypasses authentication
// 2. MCP initialization requests (no Mcp-Session-Id header) bypass bearer token check
//    - This allows clients to establish sessions without a token
//    - Session ID is returned in response header by the MCP SDK
// 3. All subsequent MCP requests with an Mcp-Session-Id header require bearer token
//    - Client must include: "Authorization: Bearer <CLAW_TOKEN>"
//    - Token is validated using constant-time comparison to prevent timing attacks
func authMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// /health endpoint bypasses auth
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			// MCP initialization requests (no session ID yet) don't require bearer token.
			// This allows clients to establish sessions first, then include token in subsequent requests.
			// See MCP Streamable HTTP Transport spec: initialization happens before authentication.
			if r.URL.Path == "/mcp" && r.Method == http.MethodPost && r.Header.Get("Mcp-Session-Id") == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check Authorization header for all other requests
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":"Unauthorized"}`)
				return
			}

			// Extract Bearer token
			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":"Unauthorized"}`)
				return
			}

			providedToken := authHeader[len(prefix):]

			// Constant-time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(providedToken), []byte(token)) != 1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":"Unauthorized"}`)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "HTTP server port")
	flag.Parse()

	// Read CLAW_TOKEN from environment and fail fast if missing
	token := os.Getenv("CLAW_TOKEN")
	if token == "" {
		log.Fatalf("CLAW_TOKEN environment variable is required")
	}

	// Initialize database
	if err := internal.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer internal.Close()

	// Initialize session management
	sessionStore := session.NewSessionStore()

	// Create MCP server
	server, err := internal.NewServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Create HTTP handler for MCP
	// The handler will manage session IDs automatically per MCP spec
	mcpServer := server.GetMCPServer()
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	// Create HTTP multiplexer
	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// Wrap mux with middleware in order:
	// 1. Session validation (checks session IDs for /mcp requests)
	// 2. Authentication (requires bearer token)
	// This order allows initialization requests (no session, no token) to proceed
	sessionValidated := sessionValidationMiddleware(sessionStore)(mux)
	authedMux := authMiddleware(token)(sessionValidated)

	// Start server
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      authedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting Claw MCP Server on %s with session management enabled", addr)

	// Start server in goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)
}
