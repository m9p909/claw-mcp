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
)

func authMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// /health endpoint bypasses auth
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			// Check Authorization header
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

	// Create MCP server
	server, err := internal.NewServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Create HTTP handler for MCP
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

	// Wrap mux with auth middleware
	authedMux := authMiddleware(token)(mux)

	// Start server
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      authedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting Claw MCP Server on %s", addr)

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
