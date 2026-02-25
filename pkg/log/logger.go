package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	requestIDKey = "request_id"
	sessionIDKey = "session_id"
)

// Logger wraps slog.Logger with context support
type Logger struct {
	inner *slog.Logger
}

// NewLogger creates a logger with DEBUG_LEVEL environment variable support.
// Levels: INFO (default), DEBUG, TRACE
func NewLogger() *Logger {
	debugLevel := os.Getenv("DEBUG_LEVEL")
	debugLevel = strings.ToUpper(debugLevel)

	var level slog.Level
	switch debugLevel {
	case "TRACE":
		level = slog.LevelDebug - 1 // TRACE is below DEBUG
	case "DEBUG":
		level = slog.LevelDebug
	default: // INFO or any unknown value defaults to INFO
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	inner := slog.New(handler)

	return &Logger{inner: inner}
}

// Info logs at INFO level with context support
func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	attrs := l.attrsWithContext(ctx, args)
	l.inner.InfoContext(ctx, msg, attrs...)
}

// Warn logs at WARN level with context support
func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	attrs := l.attrsWithContext(ctx, args)
	l.inner.WarnContext(ctx, msg, attrs...)
}

// Error logs at ERROR level with context support
func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	attrs := l.attrsWithContext(ctx, args)
	l.inner.ErrorContext(ctx, msg, attrs...)
}

// Debug logs at DEBUG level with context support
func (l *Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	attrs := l.attrsWithContext(ctx, args)
	l.inner.DebugContext(ctx, msg, attrs...)
}

// Trace logs below DEBUG level (when DEBUG_LEVEL=TRACE)
func (l *Logger) Trace(ctx context.Context, msg string, args ...interface{}) {
	// Use Debug-1 to get below DEBUG level
	if l.inner.Enabled(ctx, slog.LevelDebug-1) {
		attrs := l.attrsWithContext(ctx, args)
		l.inner.Log(ctx, slog.LevelDebug-1, msg, attrs...)
	}
}

// attrsWithContext builds slog attributes with context values and user args
func (l *Logger) attrsWithContext(ctx context.Context, args []interface{}) []any {
	var attrs []any

	// Add context-based attributes
	if reqID := RequestIDFromContext(ctx); reqID != "" {
		attrs = append(attrs, slog.String(requestIDKey, reqID))
	}

	if sessID := SessionIDFromContext(ctx); sessID != "" {
		attrs = append(attrs, slog.String(sessionIDKey, sessID))
	}

	// Add user-provided key-value pairs
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])
		var val interface{}
		if i+1 < len(args) {
			val = args[i+1]
		}
		attrs = append(attrs, slog.Any(key, val))
	}

	return attrs
}

// WithRequestID adds a request ID to context
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// WithSessionID adds a session ID to context
func WithSessionID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, sessionIDKey, id)
}

// RequestIDFromContext extracts request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// SessionIDFromContext extracts session ID from context
func SessionIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(sessionIDKey).(string); ok {
		return id
	}
	return ""
}

// Duration formats a duration in milliseconds for logging
func Duration(d time.Duration) slog.Attr {
	return slog.Float64("duration_ms", d.Seconds()*1000)
}
