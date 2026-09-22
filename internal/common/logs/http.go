package logs

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func NewStructuredLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: logger})
}

type StructuredLogger struct {
	Logger *slog.Logger
}

func (l StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: l.Logger, ctx: r.Context()}

	attrs := []any{
		"http_method", r.Method,
		"remote_addr", r.RemoteAddr,
		"uri", r.RequestURI,
	}
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		attrs = append(attrs, "req_id", reqID)
	}

	entry.Logger = entry.Logger.With(attrs...)
	entry.Logger.InfoContext(r.Context(), "Request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger *slog.Logger
	ctx context.Context
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.With(
		"resp_status", status,
		"resp_bytes_length", bytes,
		"resp_elapsed", elapsed.Round(time.Millisecond/100).String(),
	)

	l.Logger.InfoContext(l.ctx, "Request completed")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		"stack", string(stack),
		"panic", fmt.Sprintf("%+v", v),
	)
}

func GetLogEntry(r *http.Request) *slog.Logger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)

	return entry.Logger
}
