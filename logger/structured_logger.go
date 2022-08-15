package logger

// Newpackage logger

import (
	"fmt"
	"net/http"
	"order-service/config"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// NewStructuredLogger is a simple, but powerful implementation of a custom structured logger backed on logrus
func NewStructuredLogger() func(next http.Handler) http.Handler {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = config.LogLevel()

	return middleware.RequestLogger(&StructuredLogger{logger})
}

// StructuredLogger ...
type StructuredLogger struct {
	Logger *logrus.Logger
}

// NewLogEntry implementation for interface function
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "http"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields("request started")

	return entry
}

// StructuredLoggerEntry
type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

// Write implementation for interface function
func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Debugln("request complete")
}

// Panic implementation for interface function
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// GetLogEntry helper methods used by hte application to get the request-scoped
// logger entry and set additional fields between handlers.
func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

// LogEntrySetField implementation for interface function
func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

// LogEntrySetFields implementation for interface function
func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
