package slogsentry

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/getsentry/sentry-go"
)

const (
	shortErrKey = "err"
	longErrKey  = "error"
)

var slogDefaultKeys = []string{slog.TimeKey, slog.LevelKey, slog.SourceKey, slog.MessageKey, shortErrKey, longErrKey}

// SentryHandler is a Handler that writes log records to the Sentry.
type SentryHandler struct {
	slog.Handler
	levels []slog.Level
}

// NewSentryHandler creates a SentryHandler that writes to w,
// using the given options.
func NewSentryHandler(
	handler slog.Handler,
	levels []slog.Level,
) *SentryHandler {
	return &SentryHandler{
		Handler: handler,
		levels:  levels,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (s *SentryHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.Handler.Enabled(ctx, level)
}

// Handle intercepts and processes logger messages.
// In our case, send a message to the Sentry.
func (s *SentryHandler) Handle(ctx context.Context, record slog.Record) error {

	if slices.Contains(s.levels, record.Level) {
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub()
		}
		if hub == nil {
			return fmt.Errorf("sentry: hub is nil")
		}
		var err error
		slogContext := map[string]any{}
		record.Attrs(func(attr slog.Attr) bool {
			if !slices.Contains(slogDefaultKeys, attr.Key) {
				slogContext[attr.Key] = attr.Value.String()
			} else if attr.Key == shortErrKey || attr.Key == longErrKey {
				err = attr.Value.Any().(error)
			}
			return true
		})

		hub.WithScope(func(scope *sentry.Scope) {
			if len(slogContext) > 0 {
				scope.SetContext("slog", slogContext)
			}

			switch record.Level {
			case slog.LevelError:
				if err != nil {
					sentry.CaptureException(err)
				}
			case slog.LevelDebug, slog.LevelInfo, slog.LevelWarn:
				sentry.CaptureMessage(record.Message)

			}
		})
	}

	return s.Handler.Handle(ctx, record)
}

// WithAttrs returns a new SentryHandler whose attributes consists.
func (s *SentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewSentryHandler(s.Handler.WithAttrs(attrs), s.levels)
}

// WithGroup returns a new SentryHandler whose group consists.
func (s *SentryHandler) WithGroup(name string) slog.Handler {
	return NewSentryHandler(s.Handler.WithGroup(name), s.levels)
}
